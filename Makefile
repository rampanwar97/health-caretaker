# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=health-caretaker
BINARY_UNIX=$(BINARY_NAME)_unix

# Build flags
LDFLAGS=-ldflags "-X 'main.Version=$(VERSION)' -X 'main.CommitSHA=$(GIT_COMMIT)' -X 'main.BuildDate=$(BUILD_TIME)'"
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

.PHONY: all build clean test deps run docker-build docker-run help

all: test build

# Build the application
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/server

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) ./cmd/server

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run: build
	./$(BINARY_NAME)

# Run with custom config
run-config: build
	./$(BINARY_NAME) -config=config.json

# Run with environment variables
run-env: build
	WEB_PORT=$(WEB_PORT) METRICS_PORT=$(METRICS_PORT) ./$(BINARY_NAME)

# Show version
version: build
	./$(BINARY_NAME) -version

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code
lint:
	golangci-lint run

# Docker build
docker-build:
	docker build -t health-caretaker:$(VERSION) .
	docker tag health-caretaker:$(VERSION) health-caretaker:latest

# Docker run
docker-run:
	docker run -p 8080:8080 -p 9091:9091 health-caretaker:latest

# Docker Compose up
compose-up:
	docker compose up --build

# Docker Compose up with custom ports
compose-up-env:
	WEB_PORT=$(WEB_PORT) METRICS_PORT=$(METRICS_PORT) docker compose up --build

# Docker Compose down
compose-down:
	docker compose down

# Install development tools
install-tools:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Download dependencies"
	@echo "  run           - Build and run the application"
	@echo "  run-config    - Run with custom config file"
	@echo "  run-env       - Run with environment variables (WEB_PORT, METRICS_PORT)"
	@echo "  version       - Show version information"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  compose-up    - Start with Docker Compose"
	@echo "  compose-up-env - Start with Docker Compose using env vars"
	@echo "  compose-down  - Stop Docker Compose"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help"
