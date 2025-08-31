# Makefile for GoKit CLI

.PHONY: build test lint clean install uninstall help

# Build variables
BINARY_NAME=gokit
BUILD_DIR=bin
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Default target
all: build

# Build the CLI
build:
	@echo "Building GoKit CLI..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} main.go
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 main.go

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p ${BUILD_DIR}
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 main.go

build-windows:
	@echo "Building for Windows..."
	@mkdir -p ${BUILD_DIR}
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint the code
lint:
	@echo "Linting code..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ${BUILD_DIR}
	rm -f coverage.out coverage.html

# Install the CLI globally
install: build
	@echo "Installing GoKit CLI..."
	cp ${BUILD_DIR}/${BINARY_NAME} /usr/local/bin/
	@echo "GoKit CLI installed successfully"

# Uninstall the CLI
uninstall:
	@echo "Uninstalling GoKit CLI..."
	rm -f /usr/local/bin/${BINARY_NAME}
	@echo "GoKit CLI uninstalled successfully"

# Run the CLI
run: build
	@echo "Running GoKit CLI..."
	./${BUILD_DIR}/${BINARY_NAME}

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the CLI for current platform"
	@echo "  build-all     - Build for Linux, macOS, and Windows"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-darwin  - Build for macOS"
	@echo "  build-windows - Build for Windows"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Lint the code"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install CLI globally"
	@echo "  uninstall     - Uninstall CLI"
	@echo "  run           - Run the CLI"
	@echo "  help          - Show this help message"

# Development targets
dev: build
	@echo "Starting development mode..."
	./${BUILD_DIR}/${BINARY_NAME} --help

# Create release
release: clean build-all
	@echo "Creating release..."
	@mkdir -p release
	@cp ${BUILD_DIR}/* release/
	@echo "Release files created in release/ directory"

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod download

# Generate documentation
docs:
	@echo "Generating documentation..."
	@echo "Documentation is in README.md"

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	gosec ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...
