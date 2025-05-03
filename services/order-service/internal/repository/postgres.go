package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/morshedulmunna/go-microservice/pkg/config"
	"github.com/morshedulmunna/go-microservice/services/order-service/internal/domain"
)

// OrderRepository implements domain.OrderRepository
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new OrderRepository instance
func NewOrderRepository(cfg *config.Config) (*OrderRepository, error) {
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

	return &OrderRepository{db: db}, nil
}

// Create implements domain.OrderRepository
func (r *OrderRepository) Create(order *domain.Order) error {
	// Convert items to JSON
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %v", err)
	}

	query := `
		INSERT INTO orders (id, user_id, items, total, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(query,
		order.ID,
		order.UserID,
		itemsJSON,
		order.Total,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}

	return nil
}

// GetByID implements domain.OrderRepository
func (r *OrderRepository) GetByID(id string) (*domain.Order, error) {
	query := `
		SELECT id, user_id, items, total, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order domain.Order
	var itemsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&itemsJSON,
		&order.Total,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	// Unmarshal items
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %v", err)
	}

	return &order, nil
}

// GetByUserID implements domain.OrderRepository
func (r *OrderRepository) GetByUserID(userID string, page, limit int) ([]domain.Order, error) {
	query := `
		SELECT id, user_id, items, total, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var itemsJSON []byte

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&itemsJSON,
			&order.Total,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		// Unmarshal items
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %v", err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %v", err)
	}

	return orders, nil
}

// Update implements domain.OrderRepository
func (r *OrderRepository) Update(order *domain.Order) error {
	// Convert items to JSON
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %v", err)
	}

	query := `
		UPDATE orders
		SET user_id = $1, items = $2, total = $3, status = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(query,
		order.UserID,
		itemsJSON,
		order.Total,
		order.Status,
		time.Now(),
		order.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Delete implements domain.OrderRepository
func (r *OrderRepository) Delete(id string) error {
	query := `
		DELETE FROM orders
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// List implements domain.OrderRepository
func (r *OrderRepository) List(page, limit int) ([]domain.Order, error) {
	query := `
		SELECT id, user_id, items, total, status, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var itemsJSON []byte

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&itemsJSON,
			&order.Total,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		// Unmarshal items
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %v", err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %v", err)
	}

	return orders, nil
}

// UpdateStatus implements domain.OrderRepository
func (r *OrderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Close closes the database connection
func (r *OrderRepository) Close() error {
	return r.db.Close()
}
