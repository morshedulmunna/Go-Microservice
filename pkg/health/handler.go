package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response represents the health check response
type Response struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// Handler represents the health check handler
type Handler struct {
	serviceName string
	version     string
}

// NewHandler creates a new health check handler
func NewHandler(serviceName, version string) *Handler {
	return &Handler{
		serviceName: serviceName,
		version:     version,
	}
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:    "UP",
		Service:   h.serviceName,
		Version:   h.version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
