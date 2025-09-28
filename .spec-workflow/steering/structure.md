# Project Structure

## Directory Organization

```
pass-cli/
├── main.go                     # Application entry point
├── cmd/                        # CLI commands (Cobra)
│   ├── root.go                 # Root command and global flags
│   ├── init.go                 # Initialize vault command
│   ├── add.go                  # Add credential command
│   ├── get.go                  # Retrieve credential command
│   ├── list.go                 # List credentials command
│   ├── update.go               # Update credential command
│   ├── delete.go               # Delete credential command
│   └── generate.go             # Generate password command
├── internal/                   # Private application code
│   ├── crypto/                 # Encryption/decryption layer
│   │   ├── crypto.go           # AES-256-GCM implementation
│   │   └── crypto_test.go      # Crypto unit tests
│   ├── storage/                # File storage layer
│   │   ├── storage.go          # Vault file operations
│   │   └── storage_test.go     # Storage unit tests
│   ├── keychain/               # System keychain integration
│   │   ├── keychain.go         # Cross-platform keychain access
│   │   └── keychain_test.go    # Keychain unit tests
│   └── vault/                  # Business logic layer
│       ├── vault.go            # Credential management logic
│       ├── vault_test.go       # Vault unit tests
│       └── types.go            # Data structures and types
├── test/                       # Integration and end-to-end tests
│   ├── integration_test.go     # Full workflow tests
│   └── helpers.go              # Test utilities and fixtures
├── docs/                       # Project documentation
│   ├── installation.md         # Installation instructions
│   ├── usage.md                # Usage guide and examples
│   └── security.md             # Security design and threat model
├── scripts/                    # Build and utility scripts
│   ├── build.sh                # Cross-platform build script
│   └── test.sh                 # Comprehensive test runner
├── .spec-workflow/             # Specification and workflow files
│   ├── steering/               # Project steering documents
│   └── specs/                  # Feature specifications
├── Makefile                    # Build targets and automation
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── .gitignore                  # Git ignore patterns
└── README.md                   # Project overview and quick start
```

## Naming Conventions

### Files
- **Commands**: `snake_case.go` (e.g., `add.go`, `generate.go`)
- **Packages**: `lowercase` single word (e.g., `crypto`, `storage`, `keychain`)
- **Internal modules**: `lowercase.go` (e.g., `vault.go`, `types.go`)
- **Tests**: `[filename]_test.go` (e.g., `crypto_test.go`, `vault_test.go`)

### Code
- **Types/Structs**: `PascalCase` (e.g., `Credential`, `EncryptedVault`)
- **Functions/Methods**: `PascalCase` for public, `camelCase` for private (e.g., `AddCredential`, `encryptData`)
- **Constants**: `PascalCase` (e.g., `DefaultVaultPath`, `MinPasswordLength`)
- **Variables**: `camelCase` (e.g., `masterPassword`, `vaultData`)

## Import Patterns

### Import Order
1. **Standard library**: `crypto/aes`, `encoding/json`, `os`, `path/filepath`
2. **External dependencies**: `github.com/spf13/cobra`, `github.com/zalando/go-keyring`
3. **Internal packages**: `pass-cli/internal/crypto`, `pass-cli/internal/storage`

### Module/Package Organization
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
    "pass-cli/internal/vault"
)
```

## Code Structure Patterns

### Module/Package Organization
```go
// Package declaration and imports
package crypto

import (...)

// Constants and configuration
const (
    KeyLength = 32
    NonceLength = 12
)

// Type definitions
type CryptoService struct {
    // fields
}

// Public API functions
func NewCryptoService() *CryptoService { ... }
func (c *CryptoService) Encrypt(data []byte) ([]byte, error) { ... }

// Private helper functions
func deriveKey(password string, salt []byte) []byte { ... }
```

### Function/Method Organization
```go
func (v *Vault) AddCredential(service, username, value string) error {
    // Input validation
    if service == "" {
        return ErrEmptyService
    }

    // Core logic
    credential := Credential{
        Service:   service,
        Username:  username,
        Value:     value,
        CreatedAt: time.Now(),
    }

    // Error handling and persistence
    if err := v.storage.Save(credential); err != nil {
        return fmt.Errorf("failed to save credential: %w", err)
    }

    return nil
}
```

### File Organization Principles
- **One primary type per file**: Each file focuses on a single main struct or concept
- **Related functionality grouped**: Helper functions stay close to primary implementations
- **Public API first**: Exported functions and types at the top of files
- **Private implementation last**: Internal helpers and utilities at the bottom

## Code Organization Principles

1. **Single Responsibility**: Each package handles one domain (crypto, storage, keychain)
2. **Modularity**: Clear interfaces between layers enable testing and maintainability
3. **Testability**: Dependency injection and interfaces support comprehensive testing
4. **Consistency**: Follow established Go idioms and project patterns

## Module Boundaries

### Dependency Direction
```
CLI Commands (cmd/)
    ↓
Business Logic (internal/vault/)
    ↓
Service Layers (internal/crypto/, internal/storage/, internal/keychain/)
    ↓
Standard Library & External Dependencies
```

### Boundary Patterns
- **Public API vs Internal**: `cmd/` packages expose CLI interface, `internal/` packages are implementation details
- **Core vs Platform-specific**: Core crypto and vault logic is cross-platform, keychain integration handles OS differences
- **Stable vs Experimental**: Main packages are stable, future plugin system would be experimental
- **Dependencies direction**: Higher layers depend on lower layers, never vice versa

## Code Size Guidelines

**Suggested Guidelines:**
- **File size**: Maximum 500 lines per file (excluding tests)
- **Function/Method size**: Maximum 50 lines per function, prefer 10-20 lines
- **Struct complexity**: Maximum 10 fields per struct, consider composition for larger types
- **Nesting depth**: Maximum 4 levels of nesting, extract functions for complex logic

## Security Structure

### Sensitive Data Handling
```
Vault Layer: Plain text credentials (in memory only)
    ↓
Crypto Layer: Encryption/decryption operations
    ↓
Storage Layer: Encrypted data persistence
    ↓
File System: Encrypted files with secure permissions
```

### Security Boundaries
- **Memory management**: Clear sensitive data from memory after use
- **File permissions**: Vault files created with 600 permissions (user-only access)
- **Key derivation**: Master passwords never stored, only derived keys with salt
- **Error handling**: Security-sensitive errors don't leak implementation details

## Testing Structure

### Test Organization
```
Unit Tests: *_test.go files alongside source code
Integration Tests: test/ directory for cross-component testing
Test Utilities: test/helpers.go for shared test infrastructure
```

### Test Patterns
- **Table-driven tests**: Standard Go pattern for multiple test cases
- **Dependency injection**: Mock interfaces for external dependencies (keychain, filesystem)
- **Test fixtures**: Predefined data for consistent test scenarios
- **Cleanup**: Proper cleanup of test files and sensitive data

## Documentation Standards

- **Public APIs**: All exported functions and types have Go doc comments
- **Complex logic**: Inline comments for cryptographic operations and security decisions
- **Package documentation**: Each package has comprehensive package-level documentation
- **README files**: Usage examples and getting started guide in project root
- **Security documentation**: Dedicated security.md explaining threat model and design decisions