# GoKit CLI Makefile

# Variables
BINARY_NAME=gokit
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X github.com/kumarabd/gokit/cli/commands.Version=${VERSION} -X github.com/kumarabd/gokit/cli/commands.BuildTime=${BUILD_TIME} -X github.com/kumarabd/gokit/cli/commands.GitCommit=${GIT_COMMIT}"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	cd cli && go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	cd cli && GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64 .

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	cd cli && GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-amd64 .
	cd cli && GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-arm64 .

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	cd cli && GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-windows-amd64.exe .

# Install the binary
.PHONY: install
install: build
	@echo "Installing ${BINARY_NAME}..."
	cp cli/${BINARY_NAME} /usr/local/bin/

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

# Run linting
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Run all checks
.PHONY: check
check: lint test

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f cli/${BINARY_NAME}
	rm -f cli/${BINARY_NAME}-*
	rm -f coverage.txt coverage.html

# Show version
.PHONY: version
version:
	@echo "Version: ${VERSION}"
	@echo "Build Time: ${BUILD_TIME}"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for all platforms"
	@echo "  install       - Install the binary"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  check         - Run lint and tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  version       - Show version information"
	@echo "  help          - Show this help"