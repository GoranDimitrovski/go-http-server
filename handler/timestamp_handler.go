package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"simplesurance/service"
	"time"
)

// TimestampHandler handles HTTP requests for timestamp operations
type TimestampHandler struct {
	service *service.TimestampService
	logger  *log.Logger
}

// NewTimestampHandler creates a new timestamp handler
func NewTimestampHandler(service *service.TimestampService, logger *log.Logger) *TimestampHandler {
	return &TimestampHandler{
		service: service,
		logger:  logger,
	}
}

// HandleTimestamp handles GET requests to record a timestamp
func (h *TimestampHandler) HandleTimestamp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	count, err := h.service.RecordTimestamp(ctx)
	if err != nil {
		h.logger.Printf("error recording timestamp: %v", err)
		h.respondError(w, http.StatusInternalServerError, "failed to record timestamp")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]int{"count": count})
}

// respondJSON sends a JSON response
func (h *TimestampHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("error encoding JSON response: %v", err)
	}
}

// HandleHealth handles health check requests
func (h *TimestampHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// respondError sends an error response
func (h *TimestampHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
