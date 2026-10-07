package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simplesurance/internal/application"
	"simplesurance/internal/config"
	"simplesurance/internal/infrastructure/repository"
	preshttp "simplesurance/internal/presentation/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)
	memoryStore := repository.NewMemoryStore(cfg.Filename)
	timestampService := application.NewTimestampService(memoryStore, cfg.Threshold)
	timestampHandler := preshttp.NewTimestampHandler(timestampService, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := timestampService.Initialize(ctx); err != nil {
		logger.Fatalf("failed to initialize service: %v", err)
	}
	mux := http.NewServeMux()
	timestampHandler.Routes(mux, cfg.Route)

	server := &http.Server{
		Addr:         cfg.ServerAddr(),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		logger.Printf("Listening on %s", cfg.ServerAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failed to start: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	}
	if err := memoryStore.Close(); err != nil {
		logger.Printf("Error closing store: %v", err)
	}

	logger.Println("Server exited")
}
