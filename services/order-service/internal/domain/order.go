package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string  `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}

// Order represents an order entity
type Order struct {
	ID        string      `json:"id" db:"id"`
	UserID    string      `json:"user_id" db:"user_id"`
	Items     []OrderItem `json:"items" db:"items"`
	Total     float64     `json:"total" db:"total"`
	Status    OrderStatus `json:"status" db:"status"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// NewOrder creates a new order
func NewOrder(userID string, items []OrderItem) *Order {
	now := time.Now()
	total := calculateTotal(items)

	return &Order{
		ID:        uuid.New().String(),
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    OrderStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// calculateTotal calculates the total price of the order
func calculateTotal(items []OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// OrderRepository defines the order repository interface
type OrderRepository interface {
	Create(order *Order) error
	GetByID(id string) (*Order, error)
	GetByUserID(userID string, page, limit int) ([]Order, error)
	Update(order *Order) error
	Delete(id string) error
	List(page, limit int) ([]Order, error)
	UpdateStatus(id string, status OrderStatus) error
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	UserID string      `json:"user_id" validate:"required"`
	Items  []OrderItem `json:"items" validate:"required,min=1,dive"`
}

// UpdateOrderRequest represents a request to update an order
type UpdateOrderRequest struct {
	Status OrderStatus `json:"status" validate:"required,oneof=pending paid shipped delivered cancelled"`
}

// OrderResponse represents an order response
type OrderResponse struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ToResponse converts an Order to OrderResponse
func (o *Order) ToResponse() *OrderResponse {
	return &OrderResponse{
		ID:        o.ID,
		UserID:    o.UserID,
		Items:     o.Items,
		Total:     o.Total,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

// OrderListResponse represents a paginated list of orders
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

// Validate validates the create order request
func (r *CreateOrderRequest) Validate() error {
	// TODO: Implement validation
	return nil
}

// Validate validates the update order request
func (r *UpdateOrderRequest) Validate() error {
	// TODO: Implement validation
	return nil
}
