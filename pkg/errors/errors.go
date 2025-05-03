package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error types
var (
	ErrNotFound           = NewAppError(http.StatusNotFound, "Resource not found", nil)
	ErrBadRequest         = NewAppError(http.StatusBadRequest, "Bad request", nil)
	ErrUnauthorized       = NewAppError(http.StatusUnauthorized, "Unauthorized", nil)
	ErrForbidden          = NewAppError(http.StatusForbidden, "Forbidden", nil)
	ErrInternalServer     = NewAppError(http.StatusInternalServerError, "Internal server error", nil)
	ErrValidation         = NewAppError(http.StatusUnprocessableEntity, "Validation error", nil)
	ErrDuplicateEntry     = NewAppError(http.StatusConflict, "Duplicate entry", nil)
	ErrServiceUnavailable = NewAppError(http.StatusServiceUnavailable, "Service unavailable", nil)
)

// Error response structure
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   *AppError   `json:"error"`
	Data    interface{} `json:"data,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) *ErrorResponse {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		appErr = ErrInternalServer
		appErr.Err = err
	}

	return &ErrorResponse{
		Success: false,
		Error:   appErr,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprintf("validation failed on field: %s", v[0].Field)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusNotFound
}

// IsBadRequest checks if the error is a bad request error
func IsBadRequest(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusBadRequest
}

// IsUnauthorized checks if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusUnauthorized
}

// IsForbidden checks if the error is a forbidden error
func IsForbidden(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusForbidden
}

// IsInternalServer checks if the error is an internal server error
func IsInternalServer(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusInternalServerError
}

// IsValidation checks if the error is a validation error
func IsValidation(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusUnprocessableEntity
}

// IsDuplicateEntry checks if the error is a duplicate entry error
func IsDuplicateEntry(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusConflict
}

// IsServiceUnavailable checks if the error is a service unavailable error
func IsServiceUnavailable(err error) bool {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	}
	return appErr != nil && appErr.Code == http.StatusServiceUnavailable
}
