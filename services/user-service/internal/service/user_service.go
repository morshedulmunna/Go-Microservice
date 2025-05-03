package service

import (
	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/domain"
	"go.uber.org/zap"
)

// UserService implements the user service
type UserService struct {
	repo   domain.UserRepository
	logger *zap.Logger
}

// NewUserService creates a new UserService instance
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger.GetLogger(),
	}
}

// Create creates a new user
func (s *UserService) Create(user *domain.User) error {
	s.logger.Info("Creating user", zap.String("email", user.Email))
	return s.repo.Create(user)
}

// GetByID gets a user by ID
func (s *UserService) GetByID(id string) (*domain.User, error) {
	s.logger.Info("Getting user by ID", zap.String("id", id))
	return s.repo.GetByID(id)
}

// GetByEmail gets a user by email
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	s.logger.Info("Getting user by email", zap.String("email", email))
	return s.repo.GetByEmail(email)
}

// Update updates a user
func (s *UserService) Update(user *domain.User) error {
	s.logger.Info("Updating user", zap.String("id", user.ID))
	return s.repo.Update(user)
}

// Delete deletes a user
func (s *UserService) Delete(id string) error {
	s.logger.Info("Deleting user", zap.String("id", id))
	return s.repo.Delete(id)
}

// List lists users
func (s *UserService) List(page, limit int) ([]domain.User, error) {
	s.logger.Info("Listing users", zap.Int("page", page), zap.Int("limit", limit))
	return s.repo.List(page, limit)
}

// Close closes the service resources
func (s *UserService) Close() error {
	return nil
}
