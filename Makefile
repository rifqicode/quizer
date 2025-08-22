# Makefile for Quizer Backend

.PHONY: run build test clean install dev migrate

# Variables
BINARY_NAME=quizer
MAIN_PATH=cmd/main.go

# Default target
all: build

# Install dependencies
install:
	go mod download
	go mod tidy

# Run the application in development mode
dev:
	go run $(MAIN_PATH)

# Build the application
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the built binary
run: build
	./bin/$(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	go clean
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME)

migrate-create: 
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations $$name

migrate-up:
	migrate -database "$$DATABASE_URL" -path migrations -verbose up

migrate-down:
	migrate -database "$$DATABASE_URL" -path migrations -verbose down

migrate-rollback:
	migrate -database "$$DATABASE_URL" -path migrations -verbose down

# Load environment variables from .env file
include .env
export
