# Smart Fit Girl - Development Makefile

.PHONY: help build up down logs clean test

# Default target
help:
	@echo "Smart Fit Girl Development Commands:"
	@echo "  make build    - Build all Docker images"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make logs     - View logs from all services"
	@echo "  make clean    - Remove all containers and volumes"
	@echo "  make test     - Run tests for all services"
	@echo "  make db-init  - Initialize database with migrations"

# Build all services
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Clean up everything
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Run tests (will be implemented later)
test:
	@echo "Running tests..."
	@echo "TODO: Implement test commands for each service"

# Initialize database
db-init:
	@echo "Initializing database..."
	docker-compose up -d postgres
	@echo "Database initialization complete"

# Generate Go code from proto files
proto-gen:
	@echo "Generating Go code from proto files..."
	protoc --go_out=. --go-grpc_out=. proto/*.proto
	@echo "Proto generation complete"

# Development helpers
dev-web-api:
	cd services/web-api && go run cmd/main.go

dev-client:
	cd web-client && npm start
