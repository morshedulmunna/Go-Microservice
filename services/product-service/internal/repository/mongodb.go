package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/morshedulmunna/go-microservice/pkg/config"
	"github.com/morshedulmunna/go-microservice/services/product-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepository implements domain.ProductRepository
type ProductRepository struct {
	collection *mongo.Collection
}

// NewProductRepository creates a new ProductRepository instance
func NewProductRepository(cfg *config.Config) (*ProductRepository, error) {
	// Create MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Get collection
	collection := client.Database(cfg.MongoDB.Database).Collection("products")

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %v", err)
	}

	return &ProductRepository{
		collection: collection,
	}, nil
}

// Create implements domain.ProductRepository
func (r *ProductRepository) Create(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %v", err)
	}

	return nil
}

// GetByID implements domain.ProductRepository
func (r *ProductRepository) GetByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	return &product, nil
}

// Update implements domain.ProductRepository
func (r *ProductRepository) Update(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product.UpdatedAt = time.Now()
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": product.ID}, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Delete implements domain.ProductRepository
func (r *ProductRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// List implements domain.ProductRepository
func (r *ProductRepository) List(page, limit int) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %v", err)
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %v", err)
	}

	return products, nil
}

// FindByCategory implements domain.ProductRepository
func (r *ProductRepository) FindByCategory(category string, page, limit int) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"category": category}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find products by category: %v", err)
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %v", err)
	}

	return products, nil
}

// UpdateStock implements domain.ProductRepository
func (r *ProductRepository) UpdateStock(id string, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"stock": quantity},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Close closes the MongoDB connection
func (r *ProductRepository) Close(ctx context.Context) error {
	return r.collection.Database().Client().Disconnect(ctx)
}
