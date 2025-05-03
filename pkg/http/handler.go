package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/morshedulmunna/go-microservice/pkg/errors"
	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"go.uber.org/zap"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Handler represents a base HTTP handler
type Handler struct {
	logger *zap.Logger
}

// NewHandler creates a new Handler instance
func NewHandler() *Handler {
	return &Handler{
		logger: logger.GetLogger(),
	}
}

// RespondJSON sends a JSON response
func (h *Handler) RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response := Response{
		Success: status >= 200 && status < 400,
		Data:    payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RespondError sends an error response
func (h *Handler) RespondError(w http.ResponseWriter, err error) {
	var status int
	var response Response

	switch e := err.(type) {
	case *errors.AppError:
		status = e.Code
		response = Response{
			Success: false,
			Error: map[string]interface{}{
				"code":    e.Code,
				"message": e.Message,
			},
		}
	default:
		status = http.StatusInternalServerError
		response = Response{
			Success: false,
			Error: map[string]interface{}{
				"code":    status,
				"message": "Internal Server Error",
			},
		}
	}

	h.logger.Error("Request error",
		zap.Error(err),
		zap.Int("status", status),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode error response", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// DecodeJSON decodes JSON from request body
func (h *Handler) DecodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.NewAppError(http.StatusBadRequest, "Invalid JSON payload", err)
	}
	return nil
}

// GetQueryParam gets a query parameter as string
func (h *Handler) GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetQueryParamInt gets a query parameter as int
func (h *Handler) GetQueryParamInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetQueryParamBool gets a query parameter as bool
func (h *Handler) GetQueryParamBool(r *http.Request, key string, defaultValue bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// GetPathParam gets a path parameter
func (h *Handler) GetPathParam(r *http.Request, key string) string {
	// This is a placeholder. In a real application, you would use your router's
	// path parameter extraction method (e.g., gorilla/mux, chi, etc.)
	return ""
}

// ValidateRequest validates the request
func (h *Handler) ValidateRequest(r *http.Request) error {
	// Add common request validation logic here
	return nil
}

// SetResponseHeaders sets common response headers
func (h *Handler) SetResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
}

// LogRequest logs request details
func (h *Handler) LogRequest(r *http.Request) {
	h.logger.Info("Incoming request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
	)
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// GetPaginationParams extracts pagination parameters from request
func (h *Handler) GetPaginationParams(r *http.Request) PaginationParams {
	return PaginationParams{
		Page:  h.GetQueryParamInt(r, "page", 1),
		Limit: h.GetQueryParamInt(r, "limit", 10),
	}
}

// SortParams represents sorting parameters
type SortParams struct {
	Field string
	Order string
}

// GetSortParams extracts sorting parameters from request
func (h *Handler) GetSortParams(r *http.Request) SortParams {
	return SortParams{
		Field: h.GetQueryParam(r, "sort_field"),
		Order: h.GetQueryParam(r, "sort_order"),
	}
}

// FilterParams represents filtering parameters
type FilterParams struct {
	Field    string
	Operator string
	Value    string
}

// GetFilterParams extracts filtering parameters from request
func (h *Handler) GetFilterParams(r *http.Request) []FilterParams {
	// This is a placeholder. In a real application, you would implement
	// your own filter parameter parsing logic
	return nil
}
