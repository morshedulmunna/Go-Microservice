package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/morshedulmunna/go-microservice/pkg/errors"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/domain"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/service"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service *service.UserService
	logger  *zap.Logger
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// HandleUsers handles requests to /users
func (h *UserHandler) HandleUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleList(w, r)
		case http.MethodPost:
			h.handleCreate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// HandleUser handles requests to /users/{id}
func (h *UserHandler) HandleUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
		if id == "" {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.handleGet(w, r, id)
		case http.MethodPut:
			h.handleUpdate(w, r, id)
		case http.MethodDelete:
			h.handleDelete(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (h *UserHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := domain.NewUser(req.Email, req.Password, req.FirstName, req.LastName, req.Role)
	if err := h.service.Create(user); err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user.ToResponse())
}

func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	user, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

func (h *UserHandler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Active != nil {
		user.Active = *req.Active
	}

	if err := h.service.Update(user); err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.Delete(id); err != nil {
		if errors.IsNotFound(err) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to delete user", zap.Error(err))
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, err := h.service.List(page, limit)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	response := domain.UserListResponse{
		Users: make([]domain.UserResponse, len(users)),
		Page:  page,
		Limit: limit,
	}

	for i, user := range users {
		response.Users[i] = *user.ToResponse()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
