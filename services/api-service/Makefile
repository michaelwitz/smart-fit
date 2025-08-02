# Smart Fit Girl API Service Makefile

# Variables
BINARY_NAME=api-service
SWAG_BIN=/Users/michael/go/bin/swag

# Build the application
build:
	go build -o $(BINARY_NAME) .

# Run the application
run:
	go run .

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	go mod tidy
	go mod download

# Generate Swagger documentation
docs:
	$(SWAG_BIN) init

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test -v ./...

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) .

# Docker build
docker-build:
	docker build -t smart-fit-girl/api-service .

# Help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  docs        - Generate Swagger documentation"
	@echo "  fmt         - Format code"
	@echo "  test        - Run tests"
	@echo "  dev         - Run with live reload"
	@echo "  build-prod  - Build for production"
	@echo "  docker-build- Build Docker image"
	@echo "  help        - Show this help"

.PHONY: build run clean deps docs fmt test dev build-prod docker-build help
