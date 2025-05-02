.PHONY: build run clean test docker-build docker-run help

# Application settings
APP_NAME = gps_server
MAIN_PATH = cmd/tcp_app/main/main.go

# Default environment variables
export TCP_SERVER_HOST ?= localhost
export TCP_SERVER_PORT ?= 8181

help:
	@echo "TCP GPS Server - Make targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help message"

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) $(MAIN_PATH)
	@echo "Build complete"

run:
	@echo "Running $(APP_NAME) on $(TCP_SERVER_HOST):$(TCP_SERVER_PORT)..."
	@go run $(MAIN_PATH)

clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)
	@echo "Clean complete"

test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests complete"

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .
	@echo "Docker build complete"

docker-run:
	@echo "Running Docker container..."
	@docker run -p $(TCP_SERVER_PORT):$(TCP_SERVER_PORT) -e TCP_SERVER_HOST=0.0.0.0 $(APP_NAME)

default: help 