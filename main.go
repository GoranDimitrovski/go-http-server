package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simplesurance/config"
	"simplesurance/handler"
	"simplesurance/persistence"
	"simplesurance/service"
	"simplesurance/store"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)

	// Initialize dependencies
	persister := persistence.NewFilePersistence()
	memoryStore := store.NewMemoryStore(cfg.Filename, persister)
	timestampService := service.NewTimestampService(memoryStore, cfg.Threshold)
	timestampHandler := handler.NewTimestampHandler(timestampService, logger)

	// Initialize the service (load and clean expired timestamps)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := timestampService.Initialize(ctx); err != nil {
		logger.Fatalf("failed to initialize service: %v", err)
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Route, timestampHandler.HandleTimestamp)
	mux.HandleFunc("/health", timestampHandler.HandleHealth)

	server := &http.Server{
		Addr:         cfg.ServerAddr(),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Starting server at http://%s%s", cfg.Address, cfg.ServerAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	}

	// Close store
	if err := memoryStore.Close(); err != nil {
		logger.Printf("Error closing store: %v", err)
	}

	logger.Println("Server exited")
}
