# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=health-monitor
BINARY_UNIX=$(BINARY_NAME)_unix

# Build flags
LDFLAGS=-ldflags "-X health-monitoring/pkg/version.Version=$(VERSION) -X health-monitoring/pkg/version.BuildTime=$(BUILD_TIME) -X health-monitoring/pkg/version.GitCommit=$(GIT_COMMIT)"
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
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
	docker build -t health-monitor:$(VERSION) .
	docker tag health-monitor:$(VERSION) health-monitor:latest

# Docker run
docker-run:
	docker run -p 8080:8080 -p 9091:9091 health-monitor:latest

# Docker Compose up
compose-up:
	docker compose up --build

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
	@echo "  version       - Show version information"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  compose-up    - Start with Docker Compose"
	@echo "  compose-down  - Stop Docker Compose"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help"
