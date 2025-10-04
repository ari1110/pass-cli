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
│   ├── generate.go             # Generate password command
│   ├── version.go              # Version command
│   ├── helpers.go              # Helper functions (password reading)
│   └── tui/                    # TUI (Terminal User Interface) layer
│       ├── tui.go              # TUI entry point and initialization
│       ├── model.go            # Main tview model (application state)
│       ├── model_test.go       # Model unit tests
│       ├── commands.go         # tview commands (vault operations)
│       ├── messages.go         # tview message types
│       ├── events.go           # Global event handlers and keyboard shortcuts
│       ├── keys.go             # Keyboard bindings and key mappings
│       ├── helpers.go          # Layout rendering and helper functions
│       ├── helpers_test.go     # Helper function unit tests
│       ├── components/         # Reusable TUI components
│       │   ├── sidebar.go              # Category tree sidebar panel
│       │   ├── sidebar_test.go         # Sidebar unit tests
│       │   ├── category_tree.go        # Category tree logic
│       │   ├── category_tree_test.go   # Category tree unit tests
│       │   ├── metadata_panel.go       # Credential metadata panel
│       │   ├── process_panel.go        # Background process status panel
│       │   ├── command_bar.go          # Command palette / command bar
│       │   ├── command_bar_test.go     # Command bar unit tests
│       │   ├── statusbar.go            # Status bar with shortcuts
│       │   ├── statusbar_test.go       # Status bar unit tests
│       │   ├── breadcrumb.go           # Breadcrumb navigation
│       │   ├── layout_manager.go       # Responsive layout calculations
│       │   └── layout_manager_test.go  # Layout manager unit tests
│       ├── views/              # TUI view components (screens)
│       │   ├── list.go                 # Credential list view
│       │   ├── list_test.go            # List view unit tests
│       │   ├── detail.go               # Credential detail view
│       │   ├── detail_test.go          # Detail view unit tests
│       │   ├── form_add.go             # Add credential form view
│       │   ├── form_add_test.go        # Add form unit tests
│       │   ├── form_edit.go            # Edit credential form view
│       │   ├── help.go                 # Help screen view
│       │   └── confirm.go              # Confirmation dialog view
│       └── styles/             # TUI styling and theming
│           └── theme.go                # Color scheme and style definitions
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
│       ├── vault.go            # Credential management logic and types
│       └── vault_test.go       # Vault unit tests
├── test/                       # Integration and end-to-end tests
│   ├── integration_test.go            # Full workflow tests
│   ├── keychain_integration_test.go   # Keychain integration tests
│   ├── tui_integration_test.go        # TUI integration tests
│   ├── tui_dashboard_integration_test.go # TUI dashboard integration tests
│   ├── test-tui.bat                   # TUI manual testing script (Windows)
│   └── README.md                      # Test documentation
├── test-vault/                 # Test fixtures (encrypted vault for integration tests)
│   ├── vault.enc               # Test encrypted vault file
│   └── vault.enc.backup        # Backup of test vault file
├── docs/                       # Project documentation
│   ├── README.md               # Documentation index
│   ├── INSTALLATION.md         # Installation instructions
│   ├── USAGE.md                # Usage guide and examples
│   ├── SECURITY.md             # Security design and threat model
│   ├── DEVELOPMENT.md          # Development guide
│   ├── TROUBLESHOOTING.md      # Troubleshooting guide
│   ├── RELEASE.md              # Release process documentation
│   ├── CI-CD.md                # CI/CD pipeline documentation
│   ├── HOMEBREW.md             # Homebrew installation guide
│   ├── SCOOP.md                # Scoop installation guide
│   ├── development/            # Implementation tracking documents
│   │   ├── README.md                            # Development docs index
│   │   ├── DASHBOARD_IMPLEMENTATION_SUMMARY.md  # Dashboard implementation
│   │   ├── DASHBOARD_TESTING_CHECKLIST.md       # Testing checklist
│   │   ├── KEYBINDINGS_AUDIT.md                 # Keybindings audit
│   │   └── TVIEW_MIGRATION_CHECKLIST.md         # tview migration visual regression checklist
│   └── archive/                # Historical documentation
│       ├── README.md                    # Archive index
│       ├── RELEASE-DRY-RUN.md           # Pre-release validation (v0.0.1)
│       └── SECURITY-AUDIT.md            # Pre-release security audit (v0.0.1)
├── manifests/                  # Platform-agnostic package manager manifests
│   ├── winget/                 # Windows Package Manager manifest
│   │   ├── pass-cli.yaml       # WinGet package manifest
│   │   └── README.md           # WinGet manifest documentation
│   └── snap/                   # Snap package manifest
│       └── README.md           # Snap manifest documentation
├── homebrew/                   # Homebrew formula (platform-native, root location)
│   └── pass-cli.rb             # Homebrew formula file
├── scoop/                      # Scoop bucket (platform-native, root location)
│   └── pass-cli.json           # Scoop bucket manifest
├── .spec-workflow/             # Specification and workflow files
│   ├── steering/               # Project steering documents
│   │   ├── product.md          # Product vision and features
│   │   ├── tech.md             # Technology stack and architecture
│   │   └── structure.md        # Directory structure and conventions
│   ├── specs/                  # Feature specifications
│   │   ├── tview-migration/    # tview framework migration spec
│   │   └── tview-migration-remediation/ # tview migration remediation spec
│   ├── archive/                # Archived completed specs
│   │   └── specs/              # Completed specification archives
│   ├── templates/              # Spec document templates
│   ├── user-templates/         # User-customizable templates
│   ├── config.example.toml     # Example configuration file
│   └── session.json            # Active session state
├── .github/                    # GitHub-specific configuration
│   ├── workflows/              # GitHub Actions CI/CD workflows
│   │   ├── ci.yml              # Continuous integration workflow
│   │   └── release.yml         # Release automation workflow
│   └── dependabot.yml          # Dependabot configuration
├── .claude/                    # Claude Code configuration (git-ignored)
│   └── settings.local.json     # Local Claude settings
├── .serena/                    # Serena MCP server data (git-ignored)
│   ├── cache/                  # Cached analysis data
│   ├── memories/               # Project memory files
│   └── project.yml             # Serena project configuration
├── dist/                       # Build output directory (git-ignored)
│   └── [build artifacts]       # GoReleaser build artifacts
├── Makefile                    # Build targets and automation
├── .goreleaser.yml             # GoReleaser configuration
├── .mcp.json                   # Model Context Protocol server configuration
├── CLAUDE.md                   # Claude operational guide and standards
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── .gitignore                  # Git ignore patterns
├── README.md                   # Project overview and quick start
├── coverage.out                # Test coverage report (git-ignored)
└── pass-cli-test.exe           # Test binary (git-ignored)
```

## Git-Ignored Directories and Files

The following directories and files are excluded from version control (defined in `.gitignore`):

**Build Artifacts**:
- `dist/` - GoReleaser build output (binaries, archives, checksums)
- `*.exe` - Compiled executables (e.g., `pass-cli-test.exe`)
- `coverage.out` - Test coverage reports

**IDE and Tool Configuration**:
- `.claude/` - Claude Code local settings and configuration
- `.serena/` - Serena MCP server cache and memory files

**Test Data**:
- `test-vault/` - Test vault files (contains sensitive test data)

These directories exist in the working tree but are not tracked by git to keep the repository clean and prevent accidental commits of local configurations, build artifacts, or sensitive test data.

## Package Manager Organization

Pass-CLI uses two patterns for package manager files:

**Pattern A: Platform-Native (Root)**
- **homebrew/** - Homebrew formula files
- **scoop/** - Scoop bucket files
- **Rationale**: These tools have established conventions of root-level directories (e.g., homebrew-core uses Formula/, scoop uses bucket/)

**Pattern B: Platform-Agnostic (manifests/)**
- **manifests/winget/** - Windows Package Manager
- **manifests/snap/** - Snap packages
- **Rationale**: Cross-platform manifest systems consolidated under `manifests/` for clarity and organization

## Root-Level Configuration Files

**Build and Release**:
- `Makefile` - Build targets, test commands, and development automation
- `.goreleaser.yml` - GoReleaser configuration for multi-platform releases
- `go.mod` / `go.sum` - Go module dependencies and checksums

**Version Control**:
- `.gitignore` - Git exclusion patterns

**Documentation**:
- `README.md` - Project overview and quick start guide
- `CLAUDE.md` - Claude operational guide and development standards

**Tool Configuration**:
- `.mcp.json` - Model Context Protocol server configuration (Serena, spec-workflow)

## Naming Conventions

### Files
- **Commands**: `snake_case.go` (e.g., `add.go`, `generate.go`)
- **Packages**: `lowercase` single word (e.g., `crypto`, `storage`, `keychain`)
- **Internal modules**: `lowercase.go` or `snake_case.go` (e.g., `vault.go`, `category_tree.go`, `layout_manager.go`)
- **Tests**: `[filename]_test.go` (e.g., `crypto_test.go`, `vault_test.go`, `category_tree_test.go`)
- **TUI modules**: `lowercase.go` or `snake_case.go` (e.g., `model.go`, `events.go`, `form_add.go`, `command_bar.go`)
- **Test scripts**: Platform-specific extensions (e.g., `test-tui.bat` for Windows)

### Code
- **Types/Structs**: `PascalCase` (e.g., `Credential`, `EncryptedVault`)
- **Functions/Methods**: `PascalCase` for public, `camelCase` for private (e.g., `AddCredential`, `encryptData`)
- **Constants**: `PascalCase` (e.g., `DefaultVaultPath`, `MinPasswordLength`)
- **Variables**: `camelCase` (e.g., `masterPassword`, `vaultData`)

## Import Patterns

### Import Order
1. **Standard library**: `crypto/aes`, `encoding/json`, `os`, `path/filepath`
2. **External dependencies**: `github.com/spf13/cobra`, `github.com/zalando/go-keyring`, `github.com/rivo/tview`, `github.com/gdamore/tcell/v2`
3. **Internal packages**: `pass-cli/internal/crypto`, `pass-cli/internal/storage`, `pass-cli/cmd/tui/components`

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
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "github.com/spf13/cobra"
    "github.com/zalando/go-keyring"
)

// Internal packages
import (
    "pass-cli/cmd/tui/components"
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
CLI Commands (cmd/) & TUI Layer (cmd/tui/)
    ↓
Business Logic (internal/vault/)
    ↓
Service Layers (internal/crypto/, internal/storage/, internal/keychain/)
    ↓
Standard Library & External Dependencies
```

### TUI Layer Organization
The TUI follows tview's component-based architecture:
```
TUI Entry (tui.go)
    ↓
Model (model.go) - Application state and tview.Pages coordinator
    ↓
├── Events (events.go) - Global keyboard shortcuts and event handling
├── Components (components/) - Reusable UI elements
│   ├── Layout Manager - Responsive dimension calculations
│   ├── Sidebar - Category tree navigation (tview.TreeView)
│   ├── Status Bar - Context shortcuts and hints (tview.TextView)
│   ├── Breadcrumb - Navigation path display (tview.TextView)
│   ├── Metadata Panel - Credential metadata display (tview.Flex)
│   ├── Process Panel - Background process status
│   └── Command Bar - Command palette (tview.Modal + tview.InputField)
├── Views (views/) - Screen-level components
│   ├── List View - Credential browsing (tview.Table)
│   ├── Detail View - Credential details (tview.TextView)
│   ├── Forms - Add/Edit dialogs (tview.Form)
│   ├── Help - Help screen (tview.TextView)
│   └── Confirm - Confirmation dialogs (tview.Modal)
└── Styles (styles/) - Visual theming
    └── Theme - Colors, borders, typography (tcell colors)
```

### Boundary Patterns
- **Public API vs Internal**: `cmd/` packages expose CLI/TUI interface, `internal/` packages are implementation details
- **Core vs Platform-specific**: Core crypto and vault logic is cross-platform, keychain integration handles OS differences
- **CLI vs TUI**: Both layers use shared vault service, but have independent presentation logic
  - CLI: Script-friendly output, flags-based configuration
  - TUI: Interactive visual interface (tview), keyboard navigation, stateful UI
- **TUI Components vs Views**:
  - Components: Reusable, composable UI elements built on tview primitives (sidebar, panels, status bar)
  - Views: Screen-level components using tview primitives (list with tview.Table, detail with tview.TextView, forms with tview.Form)
- **TUI Framework**: Uses tview (rivo/tview) for terminal UI rendering, built on tcell for terminal control
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