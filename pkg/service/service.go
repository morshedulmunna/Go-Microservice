package service

import (
	"context"

	"github.com/morshedulmunna/go-microservice/pkg/repository"
)

// Service is a generic interface for business logic operations
type Service[T any] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *T) error

	// Get retrieves an entity by ID
	Get(ctx context.Context, id string) (*T, error)

	// List retrieves entities based on query parameters
	List(ctx context.Context, query *repository.Query) ([]T, error)

	// Update updates an entity
	Update(ctx context.Context, entity *T) error

	// Delete deletes an entity
	Delete(ctx context.Context, id string) error
}

// BaseService provides common functionality for services
type BaseService[T any] struct {
	repo repository.Repository[T]
}

// NewBaseService creates a new BaseService instance
func NewBaseService[T any](repo repository.Repository[T]) *BaseService[T] {
	return &BaseService[T]{
		repo: repo,
	}
}

// Create implements Service.Create
func (s *BaseService[T]) Create(ctx context.Context, entity *T) error {
	return s.repo.Create(ctx, entity)
}

// Get implements Service.Get
func (s *BaseService[T]) Get(ctx context.Context, id string) (*T, error) {
	return s.repo.FindByID(ctx, id)
}

// List implements Service.List
func (s *BaseService[T]) List(ctx context.Context, query *repository.Query) ([]T, error) {
	return s.repo.FindAll(ctx, query.Pagination.Page, query.Pagination.Limit)
}

// Update implements Service.Update
func (s *BaseService[T]) Update(ctx context.Context, entity *T) error {
	return s.repo.Update(ctx, entity)
}

// Delete implements Service.Delete
func (s *BaseService[T]) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// ServiceError represents a service-level error
type ServiceError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new ServiceError
func NewServiceError(code int, message string, err error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ServiceOption represents a service configuration option
type ServiceOption func(interface{}) error

// WithRepository sets the repository for a service
func WithRepository[T any](repo repository.Repository[T]) ServiceOption {
	return func(s interface{}) error {
		if service, ok := s.(Service[T]); ok {
			if baseService, ok := service.(*BaseService[T]); ok {
				baseService.repo = repo
				return nil
			}
		}
		return NewServiceError(500, "invalid service type", nil)
	}
}

// ServiceMetrics represents service metrics
type ServiceMetrics interface {
	// RecordRequest records a request metric
	RecordRequest(method string, duration float64)
	// RecordError records an error metric
	RecordError(method string)
	// RecordSuccess records a success metric
	RecordSuccess(method string)
}

// ServiceHealth represents service health information
type ServiceHealth struct {
	Status    string                 `json:"status"`
	Version   string                 `json:"version"`
	BuildTime string                 `json:"buildTime"`
	Uptime    string                 `json:"uptime"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthChecker checks service health
type HealthChecker interface {
	// Check performs a health check
	Check(ctx context.Context) (*ServiceHealth, error)
}

// ServiceRegistry represents service registration functionality
type ServiceRegistry interface {
	// Register registers the service
	Register(ctx context.Context) error
	// Deregister deregisters the service
	Deregister(ctx context.Context) error
	// GetService gets service information
	GetService(ctx context.Context, name string) ([]string, error)
}

// CircuitBreaker represents circuit breaker functionality
type CircuitBreaker interface {
	// Execute executes a function with circuit breaker protection
	Execute(func() error) error
	// State returns the current state of the circuit breaker
	State() string
}

// RateLimiter represents rate limiting functionality
type RateLimiter interface {
	// Allow checks if a request is allowed
	Allow() bool
	// Reset resets the rate limiter
	Reset()
}
