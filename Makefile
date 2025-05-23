.PHONY: help build run test clean docker-build docker-run docker-stop lint fmt deps

# Default target
help:
	@echo "Available targets:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application locally"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Run with Docker Compose"
	@echo "  make docker-stop - Stop Docker containers"
	@echo "  make lint        - Run linter"
	@echo "  make fmt         - Format code"
	@echo "  make deps        - Download dependencies"

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application locally
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	docker build -t go-markdown-notes:latest .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Stop Docker containers
docker-stop:
	docker-compose down

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
