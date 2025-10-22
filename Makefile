# Makefile for bonsai - Git Branch Cleanup CLI Tool

# Binary name
BINARY_NAME=bonsai

# Build directory
BUILD_DIR=bin

# Installation directory (user's local bin)
INSTALL_DIR=$(HOME)/.local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-s -w"

# Main package path
MAIN_PATH=./cmd/$(BINARY_NAME)

.PHONY: all build install clean test coverage lint fmt help deps

# Default target
all: clean build

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  make build           - Build the binary"
	@echo "  make install         - Build and install the binary to $(INSTALL_DIR)"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make test            - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-all        - Run all tests (unit + integration)"
	@echo "  make coverage        - Run unit tests with coverage report"
	@echo "  make coverage-all    - Run all tests with coverage report"
	@echo "  make lint            - Run linter (requires golangci-lint)"
	@echo "  make fmt             - Format Go source code"
	@echo "  make deps            - Download and tidy dependencies"
	@echo "  make run             - Build and run the application"
	@echo "  make all             - Clean and build"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

## install: Build and install the binary locally
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)
	@echo "Clean complete"

## test: Run unit tests only
test:
	@echo "Running unit tests..."
	$(GOTEST) -v ./...

## test-integration: Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./internal/git

## test-all: Run all tests (unit + integration)
test-all:
	@echo "Running all tests..."
	$(GOTEST) -v -tags=integration ./...

## coverage: Run tests with coverage report
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## coverage-all: Run all tests with coverage report
coverage-all:
	@echo "Running all tests with coverage..."
	$(GOTEST) -v -tags=integration -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(INSTALL_DIR)"; \
	fi

## fmt: Format Go source code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## deps: Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

## run: Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## uninstall: Remove installed binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled"
