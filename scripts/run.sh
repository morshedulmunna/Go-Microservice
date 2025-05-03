#!/bin/bash

# Stop all running containers
docker-compose down

# Build and start all services
docker-compose up --build -d

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Run database migrations
echo "Running database migrations..."

# User Service migrations
docker-compose exec -T user-service sh -c "cd migrations && ls *.up.sql | sort | xargs -I {} psql -h postgres -U postgres -d user_service -f {}"

# Order Service migrations
docker-compose exec -T order-service sh -c "cd migrations && ls *.up.sql | sort | xargs -I {} psql -h postgres -U postgres -d order_service -f {}"

# Product Service migrations (MongoDB)
docker-compose exec -T mongodb mongosh --username mongodb --password mongodb --authenticationDatabase admin product_service /app/migrations/*.js

# Show running containers
docker-compose ps

echo "All services are up and running!"
echo "Services endpoints:"
echo "- User Service: http://localhost:8081"
echo "- Product Service: http://localhost:8082"
echo "- Order Service: http://localhost:8083"
echo "- Consul UI: http://localhost:8500"
echo "- RabbitMQ UI: http://localhost:15672"
echo "- Prometheus: http://localhost:9090"
echo "- Grafana: http://localhost:3000"
echo "- Jaeger UI: http://localhost:16686"
echo "- Vault UI: http://localhost:8200" 