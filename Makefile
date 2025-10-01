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
	@echo "Coverage report generated: coverage.html"

# Run tests with coverage and show summary
.PHONY: test-coverage-report
test-coverage-report:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo ""
	@echo "HTML report: coverage.html"
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
.PHONY: test-race
test-race:
	go test -race -short ./...

# Run integration tests
.PHONY: test-integration
test-integration:
	go test -v -tags=integration -timeout 5m ./test

# Run integration tests in short mode (skip performance/stress tests)
.PHONY: test-integration-short
test-integration-short:
	go test -v -tags=integration -short -timeout 2m ./test

# Run all tests (unit + integration)
.PHONY: test-all
test-all: test test-integration

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

# Pre-commit checks (comprehensive)
.PHONY: pre-commit
pre-commit: fmt vet lint test-race security-scan
	@echo "All pre-commit checks passed!"

# Pre-release checks (full validation)
.PHONY: pre-release
pre-release: clean fmt vet lint test-all security-scan vuln-check release-check
	@echo "All pre-release checks passed!"

# Security scan with gosec
.PHONY: security-scan
security-scan:
	@echo "Running security scan with gosec..."
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	gosec -fmt=json -out=gosec-report.json ./...
	gosec ./...

# Vulnerability check for dependencies
.PHONY: vuln-check
vuln-check:
	@echo "Checking for vulnerable dependencies..."
	@command -v govulncheck >/dev/null 2>&1 || { echo "Installing govulncheck..."; go install golang.org/x/vuln/cmd/govulncheck@latest; }
	govulncheck ./...

# GoReleaser dry run (test release without publishing)
.PHONY: release-dry-run
release-dry-run:
	@echo "Running GoReleaser in snapshot mode (dry run)..."
	@command -v goreleaser >/dev/null 2>&1 || { echo "Error: goreleaser not installed"; exit 1; }
	goreleaser release --snapshot --clean --skip=publish

# GoReleaser build only (snapshot mode)
.PHONY: release-snapshot
release-snapshot:
	@echo "Building snapshot release..."
	@command -v goreleaser >/dev/null 2>&1 || { echo "Error: goreleaser not installed"; exit 1; }
	goreleaser build --snapshot --clean

# Validate GoReleaser configuration
.PHONY: release-check
release-check:
	@echo "Validating GoReleaser configuration..."
	@command -v goreleaser >/dev/null 2>&1 || { echo "Error: goreleaser not installed"; exit 1; }
	goreleaser check

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -f coverage.out coverage.html
	rm -f gosec-report.json
	rm -rf dist/

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

# Generate module dependency graph
.PHONY: deps-graph
deps-graph:
	@echo "Generating dependency graph..."
	go mod graph | grep pass-cli | head -20

# Tidy and verify dependencies
.PHONY: deps-tidy
deps-tidy:
	go mod tidy
	go mod verify

# Update dependencies
.PHONY: deps-update
deps-update:
	go get -u ./...
	go mod tidy

# Help target
.PHONY: help
help:
	@echo "Pass-CLI Makefile Commands"
	@echo ""
	@echo "Building:"
	@echo "  build                  Build the binary"
	@echo "  build-dev              Build with debug info"
	@echo "  build-all              Cross-compile for all platforms"
	@echo "  install                Install the binary to GOPATH"
	@echo ""
	@echo "Testing:"
	@echo "  test                   Run unit tests"
	@echo "  test-race              Run tests with race detection"
	@echo "  test-coverage          Run tests with HTML coverage report"
	@echo "  test-coverage-report   Run tests with coverage summary"
	@echo "  test-integration       Run integration tests (5min timeout)"
	@echo "  test-integration-short Run integration tests (skip perf/stress)"
	@echo "  test-all               Run all tests (unit + integration)"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt                    Format code with gofmt"
	@echo "  vet                    Run go vet"
	@echo "  lint                   Run golangci-lint"
	@echo "  check                  Run fmt + vet + lint + test"
	@echo "  pre-commit             Run all pre-commit checks"
	@echo "  pre-release            Run comprehensive pre-release validation"
	@echo ""
	@echo "Security:"
	@echo "  security-scan          Run gosec security scanner"
	@echo "  vuln-check             Check for vulnerable dependencies"
	@echo ""
	@echo "Release:"
	@echo "  release-check          Validate GoReleaser configuration"
	@echo "  release-dry-run        Test full release process (no publish)"
	@echo "  release-snapshot       Build snapshot release locally"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps-tidy              Tidy and verify go.mod"
	@echo "  deps-update            Update all dependencies"
	@echo "  deps-graph             Show dependency graph"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean                  Remove build artifacts and reports"
	@echo ""
	@echo "Help:"
	@echo "  help                   Show this help message"