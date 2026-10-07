package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"simplesurance/internal/application"
)

type TimestampHandler struct {
	service *application.TimestampService
	logger  *log.Logger
}

func NewTimestampHandler(service *application.TimestampService, logger *log.Logger) *TimestampHandler {
	return &TimestampHandler{service: service, logger: logger}
}

// Routes registers the handlers; the "GET" patterns make the mux answer 405 for other methods.
func (h *TimestampHandler) Routes(mux *http.ServeMux, route string) {
	mux.HandleFunc("GET /health", h.HandleHealth)
	mux.HandleFunc("GET "+route, h.HandleTimestamp)
}

func (h *TimestampHandler) HandleTimestamp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	count, err := h.service.RecordTimestamp(ctx)
	if err != nil {
		h.logger.Printf("error recording timestamp: %v", err)
		h.respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to record timestamp"})
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *TimestampHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *TimestampHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("error encoding JSON response: %v", err)
	}
}
