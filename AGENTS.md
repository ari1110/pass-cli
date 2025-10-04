# Agent Guidelines for pass-cli

## Build/Lint/Test Commands

```bash
# Build
go build -ldflags "-X main.Version=dev -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X main.CommitHash=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")" -o pass-cli .
make build              # Build binary via Makefile
make build-dev          # Build with debug info

# Testing
go test ./...           # Run all unit tests
go test -v ./path/to/package -run TestName    # Run single test
make test-race          # Run with race detection
make test-integration   # Run integration tests
make test-coverage      # Generate coverage report

# Code Quality
go fmt ./...            # Format code
go vet ./...            # Static analysis
golangci-lint run       # Comprehensive linting
make check              # fmt + vet + lint + test
make pre-commit         # All pre-commit checks
goimports -w .          # Format and fix imports
```

## Code Style Guidelines

### Import Organization
```go
// Standard library imports
import (
    "crypto/aes"
    "encoding/json"
    "os"
)

// External dependencies
import (
    "github.com/spf13/cobra"
    "github.com/zalando/go-keyring"
)

// Internal packages
import (
    "pass-cli/internal/crypto"
    "pass-cli/internal/storage"
)
```

### File Organization
- **One primary type per file**: Each file focuses on a single main struct or concept
- **Public API first**: Exported functions and types at the top of files
- **Private implementation last**: Internal helpers and utilities at the bottom
- **Maximum file size**: 500 lines (excluding tests)
- **Maximum function size**: 50 lines, prefer 10-20 lines

### Naming Conventions
- **Files**: `snake_case.go` (commands), `lowercase.go` (internal), `*_test.go` (tests)
- **Types/Structs**: `PascalCase` (e.g., `Credential`, `EncryptedVault`)
- **Functions**: `PascalCase` for public, `camelCase` for private
- **Constants**: `PascalCase` (e.g., `DefaultVaultPath`, `MinPasswordLength`)
- **Variables**: `camelCase` (e.g., `masterPassword`, `vaultData`)

### Error Handling
- Always handle errors explicitly; never ignore with _
- Use `fmt.Errorf("context: %w", err)` for error wrapping
- Define package-level error variables
- Security-sensitive errors don't leak implementation details

### Security (CRITICAL)
- Never log sensitive data or expose credentials in error messages
- Use AES-256-GCM encryption with PBKDF2 key derivation (100k+ iterations)
- Clear sensitive memory with defer after use
- Set file permissions to 600 for vault files
- Validate all inputs before processing
- Follow defense-in-depth principles

### Testing
- Table-driven tests preferred for multiple test cases
- Unit tests in same package: `*_test.go`
- Integration tests in `test/` directory
- Target 90%+ code coverage
- Use dependency injection for mocking external dependencies
- Proper cleanup of test files and sensitive data

### Architecture
- **Module boundaries**: CLI → Business Logic → Service Layers → Dependencies
- **Dependency injection**: Use interfaces for testability and modularity
- **Package names**: Match directory names (lowercase single words)
- **Context**: Use `context.Context` for cancellation and timeouts
- **Composition**: Prefer composition over inheritance

## Spec-Driven Development

**ALWAYS use spec-workflow MCP tools** for feature work. Read steering docs first. Follow workflow: Requirements → Design → Tasks → Implementation. NEVER proceed without dashboard approval.