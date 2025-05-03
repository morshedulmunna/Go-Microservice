package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Role      string    `json:"role" db:"role"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user
func NewUser(email, password, firstName, lastName, role string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UserRepository defines the user repository interface
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List(page, limit int) ([]User, error)
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=admin user"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty" validate:"omitempty,oneof=admin user"`
	Active    *bool  `json:"active,omitempty"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// Validate validates the create user request
func (r *CreateUserRequest) Validate() error {
	// TODO: Implement validation
	return nil
}

// Validate validates the update user request
func (r *UpdateUserRequest) Validate() error {
	// TODO: Implement validation
	return nil
}
