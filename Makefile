.PHONY: all build clean test run docker-build docker-run docker-stop proto lint help

# Variables
BINARY_NAME=go-microservice
GO=go
DOCKER_COMPOSE=docker-compose
PROTOC=protoc

# Build targets
all: clean build test

build:
	@echo "Building services..."
	$(GO) build -o bin/user-service ./services/user-service/cmd/api
	$(GO) build -o bin/product-service ./services/product-service/cmd/api
	$(GO) build -o bin/order-service ./services/order-service/cmd/api

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf vendor/

test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

run:
	@echo "Running services..."
	./scripts/run.sh

# Docker targets
docker-build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-run:
	@echo "Starting services..."
	$(DOCKER_COMPOSE) up -d

docker-stop:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

# Protocol buffer targets
proto:
	@echo "Generating protobuf code..."
	$(PROTOC) --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/**/*.proto

# Development targets
lint:
	@echo "Running linters..."
	golangci-lint run

# Dependency management
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	cd services/user-service && $(GO) mod download
	cd services/product-service && $(GO) mod download
	cd services/order-service && $(GO) mod download

tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy
	cd services/user-service && $(GO) mod tidy
	cd services/product-service && $(GO) mod tidy
	cd services/order-service && $(GO) mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, build, and test"
	@echo "  build        - Build all services"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  run          - Run services locally"
	@echo "  docker-build - Build Docker images"
	@echo "  docker-run   - Start services with Docker"
	@echo "  docker-stop  - Stop Docker services"
	@echo "  proto        - Generate protobuf code"
	@echo "  lint         - Run linters"
	@echo "  deps         - Download dependencies"
	@echo "  tidy         - Tidy dependencies"
	@echo "  help         - Show this help message" 