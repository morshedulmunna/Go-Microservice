package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product entity
type Product struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Stock       int       `json:"stock" bson:"stock"`
	Category    string    `json:"category" bson:"category"`
	Active      bool      `json:"active" bson:"active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// NewProduct creates a new product
func NewProduct(name, description string, price float64, stock int, category string) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Category:    category,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ProductRepository defines the product repository interface
type ProductRepository interface {
	Create(product *Product) error
	GetByID(id string) (*Product, error)
	Update(product *Product) error
	Delete(id string) error
	List(page, limit int) ([]Product, error)
	FindByCategory(category string, page, limit int) ([]Product, error)
	UpdateStock(id string, quantity int) error
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	Category    string  `json:"category" validate:"required"`
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
	Category    string   `json:"category,omitempty"`
	Active      *bool    `json:"active,omitempty"`
}

// ProductResponse represents a product response
type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Product to ProductResponse
func (p *Product) ToResponse() *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		Active:      p.Active,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// Validate validates the create product request
func (r *CreateProductRequest) Validate() error {
	// TODO: Implement validation
	return nil
}

// Validate validates the update product request
func (r *UpdateProductRequest) Validate() error {
	// TODO: Implement validation
	return nil
}
