# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean test test-verbose test-coverage test-race deps run

all: deps test build

# Install dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GOTEST) ./...

# Run tests with verbose output
test-verbose:
	@echo "$(YELLOW)Running tests with verbose output...$(NC)"
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	$(GOTEST) -v -cover ./...
	@echo "$(YELLOW)Generating coverage report...$(NC)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Run tests with race detection
test-race:
	@echo "$(YELLOW)Running tests with race detection...$(NC)"
	$(GOTEST) -race ./...

# Run all tests (unit + integration + coverage)
test-all: test-verbose test-coverage test-race
	@echo "$(GREEN)All tests completed!$(NC)"

# Build the application
build:
	@echo "$(BLUE)Building application...$(NC)"
	$(GOBUILD) -o bin/ecommerce-server main.go

# Run the application
run:
	@echo "$(BLUE)Starting server...$(NC)"
	$(GOCMD) run main.go

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -f bin/ecommerce-server
	rm -f coverage.out coverage.html

# Run tests and show coverage
coverage: test-coverage
	@echo "$(GREEN)Coverage report available at coverage.html$(NC)"

# Run a specific test
test-specific:
	@echo "$(YELLOW)Running specific test...$(NC)"
	$(GOTEST) -v -run $(TEST) ./...

# Help
help:
	@echo "$(BLUE)Available commands:$(NC)"
	@echo "  make deps          - Install dependencies"
	@echo "  make test          - Run tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-race     - Run tests with race detection"
	@echo "  make test-all      - Run all tests (unit + integration + coverage)"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make coverage      - Generate and show coverage report"
	@echo "  make help          - Show this help message"
