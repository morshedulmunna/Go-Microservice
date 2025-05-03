# Go Microservices Architecture

A production-grade microservices architecture using Go 1.24.2, following clean architecture principles and best practices.

## Services

- **User Service**: Handles user management and authentication
- **Product Service**: Manages product catalog and inventory
- **Order Service**: Processes and manages orders

## Features

- Clean Architecture
- Domain-Driven Design
- CQRS Pattern
- Event-Driven Architecture
- Service Discovery with Consul
- Message Queue with RabbitMQ
- Distributed Tracing with Jaeger
- Metrics with Prometheus
- Monitoring with Grafana
- Secrets Management with Vault
- Rate Limiting
- Circuit Breaking
- Logging with Zap
- Configuration with Viper
- HTTP and gRPC APIs
- OpenAPI Documentation
- Docker Containerization
- Kubernetes Ready

## Prerequisites

- Go 1.24.2 or later
- Docker and Docker Compose
- Make
- Protocol Buffers Compiler
- PostgreSQL (for User and Order services)
- MongoDB (for Product service)

## Project Structure

```
.
├── cmd/                    # Application entry points
├── config/                 # Configuration files
├── deployments/            # Deployment configurations
│   ├── docker/            # Docker configurations
│   ├── k8s/               # Kubernetes manifests
│   └── prometheus/        # Prometheus configurations
├── docs/                   # Documentation
├── pkg/                    # Shared packages
│   ├── config/            # Configuration package
│   ├── errors/            # Error handling
│   ├── logger/            # Logging package
│   ├── middleware/        # HTTP middleware
│   └── messaging/         # Message queue
├── proto/                  # Protocol buffer definitions
├── scripts/               # Build and deployment scripts
└── services/              # Microservices
    ├── user-service/      # User service
    ├── product-service/   # Product service
    └── order-service/     # Order service
```

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/morshedulmunna/go-microservice.git
   cd go-microservice
   ```

2. Install dependencies:

   ```bash
   make deps
   ```

3. Build the services:

   ```bash
   make build
   ```

4. Start the services:
   ```bash
   make run
   ```

## Development

- **Build**: `make build`
- **Test**: `make test`
- **Run**: `make run`
- **Clean**: `make clean`
- **Lint**: `make lint`
- **Generate Protocol Buffers**: `make proto`

## Docker

- **Build Images**: `make docker-build`
- **Start Services**: `make docker-run`
- **Stop Services**: `make docker-stop`

## Service Endpoints

- User Service: http://localhost:8081
- Product Service: http://localhost:8082
- Order Service: http://localhost:8083
- Consul UI: http://localhost:8500
- RabbitMQ UI: http://localhost:15672
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- Jaeger UI: http://localhost:16686
- Vault UI: http://localhost:8200

## API Documentation

- User Service: http://localhost:8081/swagger/index.html
- Product Service: http://localhost:8082/swagger/index.html
- Order Service: http://localhost:8083/swagger/index.html

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
