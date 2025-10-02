# Development Guide

This guide covers the development workflow for Pass-CLI contributors.

## Prerequisites

### Required Tools

- **Go 1.25+**: [Download](https://go.dev/dl/)
- **Git**: For version control
- **golangci-lint**: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **GoReleaser**: `go install github.com/goreleaser/goreleaser/v2@latest`

### Optional Tools

- **Make**: For convenient command shortcuts (Windows: via Git Bash, WSL, or install GNU Make)
- **gosec**: Security scanner (`go install github.com/securego/gosec/v2/cmd/gosec@latest`)
- **govulncheck**: Vulnerability checker (`go install golang.org/x/vuln/cmd/govulncheck@latest`)

## Quick Start

```bash
# Clone the repository
git clone https://github.com/ari1110/pass-cli.git
cd pass-cli

# Install dependencies
go mod download

# Build the binary
go build -o pass-cli .

# Run tests
go test ./...

# Run the binary
./pass-cli --help
```

## Makefile Commands

The project includes a comprehensive Makefile with many convenient targets.

### Building

```bash
make build                  # Build the binary
make build-dev              # Build with debug info
make build-all              # Cross-compile for all platforms
make install                # Install to GOPATH
```

### Testing

```bash
make test                   # Run unit tests
make test-race              # Run tests with race detection
make test-coverage          # Generate HTML coverage report
make test-coverage-report   # Show coverage summary + HTML
make test-integration       # Run integration tests (5min)
make test-integration-short # Quick integration tests
make test-all               # All tests (unit + integration)
```

### Code Quality

```bash
make fmt                    # Format code
make vet                    # Run go vet
make lint                   # Run golangci-lint
make check                  # Run fmt + vet + lint + test
make pre-commit             # Comprehensive pre-commit checks
make pre-release            # Full pre-release validation
```

### Security

```bash
make security-scan          # Run gosec security scanner
make vuln-check             # Check for vulnerable dependencies
```

### Release

```bash
make release-check          # Validate GoReleaser config
make release-dry-run        # Test full release (no publish)
make release-snapshot       # Build snapshot release locally
```

### Dependencies

```bash
make deps-tidy              # Tidy and verify go.mod
make deps-update            # Update all dependencies
make deps-graph             # Show dependency graph
```

### Cleanup

```bash
make clean                  # Remove build artifacts
```

## Development Workflow

### Making Changes

1. **Create a branch**:
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes**:
   - Write code following Go best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   make test-all              # Run all tests
   make lint                  # Check code quality
   make security-scan         # Security check
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **Push and create PR**:
   ```bash
   git push origin feature/my-feature
   # Then create a pull request on GitHub
   ```

### Before Committing

Run the pre-commit checks:

```bash
make pre-commit
```

This runs:
- Code formatting
- `go vet`
- golangci-lint
- Tests with race detection
- Security scanning

### Before Creating a PR

Ensure all tests pass:

```bash
make test-all              # Unit + integration tests
make lint                  # Linting
make security-scan         # Security check
```

### Before Releasing

Run full pre-release validation:

```bash
make pre-release
```

This runs:
- All code quality checks
- All tests (unit + integration)
- Security scanning
- Vulnerability checking
- GoReleaser validation

## Testing

### Unit Tests

```bash
# Run all unit tests
go test ./...

# With verbose output
go test -v ./...

# With race detection
go test -race ./...

# With coverage
go test -cover ./...
```

### Integration Tests

Integration tests are marked with build tags:

```bash
# Run integration tests
go test -v -tags=integration ./test

# Skip slow tests
go test -v -tags=integration -short ./test
```

### Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

## Code Style

### Formatting

Code is formatted with `gofmt`:

```bash
go fmt ./...
# or
make fmt
```

### Linting

We use golangci-lint with strict configuration:

```bash
golangci-lint run
# or
make lint
```

### Documentation

- All exported functions must have Go doc comments
- Comments should explain "why" not "what"
- Keep comments concise and clear

Example:

```go
// GeneratePassword creates a cryptographically secure random password
// with the specified length and character requirements.
func GeneratePassword(length int, opts PasswordOptions) (string, error) {
    // Implementation...
}
```

## Project Structure

```
pass-cli/
├── cmd/                   # Cobra command definitions
│   ├── root.go           # Root command
│   ├── init.go           # Init command
│   ├── add.go            # Add command
│   └── ...
├── internal/             # Private application code
│   ├── crypto/           # Encryption/decryption
│   ├── storage/          # Vault file operations
│   ├── keychain/         # OS keychain integration
│   ├── vault/            # Vault service (business logic)
│   └── models/           # Data models
├── test/                 # Integration tests
│   └── integration_test.go
├── docs/                 # Documentation
├── .github/              # GitHub Actions workflows
├── main.go               # Application entry point
├── Makefile              # Build commands
└── .goreleaser.yml       # Release configuration
```

## Security

### Secure Coding Practices

1. **Never log sensitive data**: Passwords, encryption keys, etc.
2. **Use constant-time operations**: For cryptographic comparisons
3. **Clear sensitive data from memory**: After use
4. **Validate all inputs**: Especially user-provided data
5. **Use secure defaults**: Safe configuration out of the box

### Security Scanning

Run gosec regularly:

```bash
make security-scan
```

Check for vulnerable dependencies:

```bash
make vuln-check
```

## Debugging

### Enable Verbose Logging

```bash
./pass-cli --verbose <command>
```

### Build with Debug Info

```bash
make build-dev
# or
go build -gcflags="all=-N -l" -o pass-cli .
```

### Use Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug . -- init
```

## Performance

### Benchmarking

```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.

# Profile memory
go test -memprofile=mem.prof -bench=.

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Performance Targets

- First vault unlock: <500ms
- Cached operations: <100ms
- Support 1000+ credentials efficiently

## Release Process

See [RELEASE.md](RELEASE.md) for detailed release instructions.

Quick reference:

```bash
# 1. Test everything
make pre-release

# 2. Create tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 3. Push tag (triggers CI/CD)
git push origin v1.0.0

# 4. Monitor GitHub Actions
# 5. Verify release artifacts
```

## Troubleshooting

### Common Issues

**"golangci-lint not found"**:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**"make: command not found" (Windows)**:
- Use Git Bash, WSL, or install GNU Make
- Or run commands directly (see Makefile for exact commands)

**Tests failing with keychain errors**:
- Some keychain tests may fail without proper OS configuration
- Integration tests handle this gracefully

**Build fails with module errors**:
```bash
go mod tidy
go mod verify
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/ari1110/pass-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ari1110/pass-cli/discussions)
- **Documentation**: [docs/](.)

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [GoReleaser](https://goreleaser.com/)
- [golangci-lint](https://golangci-lint.run/)