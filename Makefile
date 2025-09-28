# Variables
BINARY_NAME=pass-cli
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Build for development with debug info
.PHONY: build-dev
build-dev:
	go build -gcflags="all=-N -l" -o $(BINARY_NAME) .

# Install the binary
.PHONY: install
install:
	go install $(LDFLAGS) .

# Run tests
.PHONY: test
test:
	go test ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run integration tests
.PHONY: test-integration
test-integration:
	go test -tags=integration ./test/...

# Lint the code
.PHONY: lint
lint:
	golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet the code
.PHONY: vet
vet:
	go vet ./...

# Check code quality
.PHONY: check
check: fmt vet lint test

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Cross-compile for all platforms
.PHONY: build-all
build-all:
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-arm64.exe .
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .

# Create release directory
.PHONY: prepare-dist
prepare-dist:
	mkdir -p dist

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           Build the binary"
	@echo "  build-dev       Build with debug info"
	@echo "  build-all       Cross-compile for all platforms"
	@echo "  install         Install the binary"
	@echo "  test            Run unit tests"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  test-integration Run integration tests"
	@echo "  lint            Run linter"
	@echo "  fmt             Format code"
	@echo "  vet             Run go vet"
	@echo "  check           Run all code quality checks"
	@echo "  clean           Clean build artifacts"
	@echo "  help            Show this help message"