package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/morshedulmunna/go-microservice/pkg/config"
	"github.com/morshedulmunna/go-microservice/services/user-service/internal/domain"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(cfg *config.Config) (*UserRepository, error) {
	// Create connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.Database.MaxConnections)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetimeConnections) * time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &UserRepository{db: db}, nil
}

// Create implements domain.UserRepository
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// GetByID implements domain.UserRepository
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetByEmail implements domain.UserRepository
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// Update implements domain.UserRepository
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, first_name = $3, last_name = $4, role = $5, active = $6, updated_at = $7
		WHERE id = $8
	`

	result, err := r.db.Exec(query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Active,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete implements domain.UserRepository
func (r *UserRepository) Delete(id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List implements domain.UserRepository
func (r *UserRepository) List(page, limit int) ([]domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %v", err)
	}

	return users, nil
}

// Close closes the database connection
func (r *UserRepository) Close() error {
	return r.db.Close()
}
