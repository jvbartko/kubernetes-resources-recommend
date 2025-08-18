# Makefile for kubernetes-resources-recommend

# Variables
BINARY_NAME=kubernetes-resources-recommend
MAIN_PATH=cmd/kubernetes-resources-recommend/main.go
BUILD_DIR=bin

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

# Make Directory to store built binary
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Build the binary
build: $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(GOBIN)/$(BINARY_NAME)"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(GOBIN)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

# Run the application (requires parameters)
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH) $(ARGS)

# Build for multiple platforms
build-all: $(BUILD_DIR)
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(GOBIN)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(GOBIN)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Cross-platform build completed"

# Show help
help:
	@echo "Available commands:"
	@echo "  build        - Build the binary"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  vet          - Vet code"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install the binary"
	@echo "  run          - Run the application (use ARGS= for parameters)"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  help         - Show this help"

.PHONY: build deps test test-coverage fmt lint vet clean install run build-all help
