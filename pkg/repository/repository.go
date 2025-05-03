package repository

import (
	"context"
)

// Repository is a generic interface for database operations
type Repository[T any] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *T) error

	// FindByID finds an entity by ID
	FindByID(ctx context.Context, id string) (*T, error)

	// FindAll returns all entities
	FindAll(ctx context.Context, page, limit int) ([]T, error)

	// Update updates an entity
	Update(ctx context.Context, entity *T) error

	// Delete deletes an entity
	Delete(ctx context.Context, id string) error
}

// Options represents database connection options
type Options struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// Pagination represents pagination parameters
type Pagination struct {
	Page  int
	Limit int
	Total int64
}

// SortOrder represents sort order
type SortOrder string

const (
	// ASC represents ascending order
	ASC SortOrder = "ASC"
	// DESC represents descending order
	DESC SortOrder = "DESC"
)

// Sort represents sort parameters
type Sort struct {
	Field string
	Order SortOrder
}

// Filter represents filter parameters
type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

// Query represents a database query
type Query struct {
	Pagination *Pagination
	Sort       []Sort
	Filters    []Filter
}

// NewQuery creates a new Query instance
func NewQuery() *Query {
	return &Query{
		Pagination: &Pagination{
			Page:  1,
			Limit: 10,
		},
		Sort:    make([]Sort, 0),
		Filters: make([]Filter, 0),
	}
}

// WithPagination adds pagination to the query
func (q *Query) WithPagination(page, limit int) *Query {
	q.Pagination.Page = page
	q.Pagination.Limit = limit
	return q
}

// WithSort adds sort parameters to the query
func (q *Query) WithSort(field string, order SortOrder) *Query {
	q.Sort = append(q.Sort, Sort{
		Field: field,
		Order: order,
	})
	return q
}

// WithFilter adds a filter to the query
func (q *Query) WithFilter(field, operator string, value interface{}) *Query {
	q.Filters = append(q.Filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return q
}

// Transaction represents a database transaction
type Transaction interface {
	// Commit commits the transaction
	Commit() error
	// Rollback rolls back the transaction
	Rollback() error
}

// TransactionFunc is a function that executes within a transaction
type TransactionFunc func(Transaction) error

// TransactionManager manages database transactions
type TransactionManager interface {
	// Begin starts a new transaction
	Begin() (Transaction, error)
	// WithTransaction executes a function within a transaction
	WithTransaction(context.Context, TransactionFunc) error
}

// QueryBuilder builds database queries
type QueryBuilder interface {
	// Select builds a SELECT query
	Select(table string, columns []string) QueryBuilder
	// Where adds a WHERE clause
	Where(condition string, args ...interface{}) QueryBuilder
	// OrderBy adds an ORDER BY clause
	OrderBy(column string, order SortOrder) QueryBuilder
	// Limit adds a LIMIT clause
	Limit(limit int) QueryBuilder
	// Offset adds an OFFSET clause
	Offset(offset int) QueryBuilder
	// Build returns the final query string and arguments
	Build() (string, []interface{})
}

// DatabaseError represents a database error
type DatabaseError struct {
	Message string
	Err     error
}

// Error implements the error interface
func (e *DatabaseError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new DatabaseError
func NewDatabaseError(message string, err error) *DatabaseError {
	return &DatabaseError{
		Message: message,
		Err:     err,
	}
}
