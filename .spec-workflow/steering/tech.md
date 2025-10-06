# Technology Stack

## Project Type
Command-line interface (CLI) tool for secure credential management. Single-binary application designed for local development environments with cross-platform compatibility.

## Core Technologies

### Primary Language(s)
- **Language**: Go 1.25.1 (modern Go with enhanced security features)
- **Runtime/Compiler**: Go compiler with CGO disabled for static linking
- **Language-specific tools**: Go modules for dependency management, built-in toolchain for cross-compilation

### Key Dependencies/Libraries
- **github.com/spf13/cobra v1.10.1**: CLI framework and command management
- **github.com/spf13/viper v1.21.0**: Configuration management and file handling
- **github.com/zalando/go-keyring v0.2.6**: Cross-platform system keychain integration
- **github.com/atotto/clipboard v0.1.4**: Cross-platform clipboard operations
- **github.com/howeyc/gopass v0.0.0-20210920133722-c8aef6fb66ef**: Secure password input with masking
- **github.com/olekukonko/tablewriter v1.1.0**: Formatted table output for credential display
- **github.com/charmbracelet/bubbletea v1.3.10**: TUI framework with Elm-inspired architecture (Model-Update-View pattern) - LEGACY implementation
- **github.com/rivo/tview v0.42.0**: Terminal UI framework (indirect dependency, actively used for migration)
- **github.com/gdamore/tcell/v2 v2.9.0**: Terminal cell-based view library (indirect, used by tview)
- **github.com/charmbracelet/lipgloss v1.1.0**: Terminal styling library for layout, colors, and borders
- **github.com/charmbracelet/bubbles v0.21.0**: Reusable TUI components (text input, viewport, list)
- **golang.org/x/crypto v0.42.0**: Extended cryptographic functions (PBKDF2)
- **golang.org/x/term v0.35.0**: Terminal detection and input handling
- **Standard library crypto packages**: AES-256-GCM implementation, secure random generation

### Application Architecture
Layered architecture with clear separation of concerns:
- **CLI Layer**: Command interface using Cobra framework with script-friendly output modes
- **TUI Layer**: Interactive dashboard currently in transition from Bubble Tea to tview
  - **Current State**: Dual implementation - Bubble Tea (legacy, default) and tview (migration in progress)
  - **Bubble Tea Components** (legacy): Model-Update-View pattern, component-based design with sidebar, panels, and status bar
  - **tview Components** (new): TreeView-based sidebar, Flex layouts, Modal components, InputField components
  - **Component-Based Design**: Reusable UI components (sidebar, metadata panel, command bar, breadcrumb, status bar)
  - **Both implementations coexist**: Toggle in cmd/tui/tui.go (useBubbleTea flag currently true)
  - **Responsive Layout System**: Multi-panel dashboard with breakpoint-based adaptation
    - Breakpoints: 80 columns (medium - shows sidebar + main), 120 columns (large - shows sidebar + main + metadata)
    - Layout calculation handled by LayoutManager component
    - Panels automatically show/hide based on terminal width
  - **Styles & Theming**: Lipgloss-based styling for Bubble Tea, tcell colors for tview
    - Theme defined in cmd/tui/styles/theme.go
    - Rounded borders for panels
    - Active/inactive panel color distinction
    - Context-aware keyboard shortcuts displayed in status bar
  - **Testing**: Multi-terminal emulator testing required (Windows Terminal, iTerm2, gnome-terminal, Alacritty)
- **Service Layer**: Business logic for credential management (Vault service) with automatic usage tracking
- **Storage Layer**: Encrypted file operations and persistence with atomic writes
- **Crypto Layer**: AES-256-GCM encryption and key derivation
- **Keychain Layer**: System integration for master password storage (Windows Credential Manager, macOS Keychain, Linux Secret Service)

### Data Storage (if applicable)
- **Primary storage**: Local encrypted JSON files in user home directory (`~/.pass-cli/vault.enc`)
- **Usage tracking**: Embedded within vault data structure, tracking location and timestamp per credential
- **Caching**: In-memory vault cache during active sessions
- **Data formats**: JSON for structured data, binary for encrypted blobs
- **Backup strategy**: Atomic writes with temporary files for corruption prevention, automatic backups before saves

### External Integrations (if applicable)
- **System Keychains**: Windows Credential Manager, macOS Keychain, Linux Secret Service
- **Clipboard**: Cross-platform clipboard integration for credential copying
- **File System**: Secure file permissions (600) for vault storage

## Development Environment

### Build & Development Tools
- **Build System**: Go toolchain with custom Makefile for cross-compilation
- **Package Management**: Go modules with dependency pinning for reproducible builds
- **Development workflow**: Live reload via `go run`, integrated testing with `go test`
- **Cross-compilation**: Native Go support for Windows, macOS, Linux (amd64, arm64)

### Code Quality Tools
- **Static Analysis**: golangci-lint v2.5 (comprehensive linter suite, configured in CI)
- **Security Scanning**: gosec v2 (Go security checker, runs in CI/CD pipeline)
- **Vulnerability Checking**: govulncheck (Go vulnerability database scanner)
- **Formatting**: goimports and gofmt (automatic import management and code formatting)
- **Testing Framework**: Go's built-in testing package with table-driven tests
- **Code Coverage**: go test with coverage analysis and Codecov integration
- **Documentation**: Go doc comments and README.md with usage examples

### Version Control & Collaboration
- **VCS**: Git with conventional commit messages
- **Branching Strategy**: GitHub Flow with feature branches and pull requests
- **Code Review Process**: Required reviews for main branch, automated CI checks
- **CI/CD Platform**: GitHub Actions
  - **CI Workflow**: Unit tests, integration tests, linting, security scans on all PRs and main branch pushes
  - **Test Matrix**: Cross-platform testing on Ubuntu, macOS, and Windows with Go 1.25
  - **Release Workflow**: Automated releases with GoReleaser on version tags
  - **Actions Used**: actions/checkout@v5, actions/setup-go@v6, golangci/golangci-lint-action@v8, goreleaser/goreleaser-action@v6, codecov/codecov-action@v5

## Deployment & Distribution
- **Target Platform(s)**: Windows 10+, macOS 10.15+, Linux distributions with glibc 2.17+
- **Build & Release Tool**: GoReleaser v2 (latest version, automated via GitHub Actions)
  - **Build Configuration**: .goreleaser.yml with CGO_ENABLED=0, -trimpath, -mod=readonly flags
  - **Version Injection**: Build-time ldflags for version, commit hash, and build date
  - **Universal Binaries**: Automatic creation of macOS universal binaries (amd64 + arm64)
- **Distribution Method**:
  - GitHub Releases with automated binary builds (Windows, macOS, Linux for amd64 and arm64)
  - Homebrew tap (ari1110/homebrew-tap) with automated formula updates
  - Scoop bucket (ari1110/scoop-bucket) with automated manifest updates
  - Direct binary download with SHA256 checksums
- **Archive Formats**: .tar.gz for Unix-like systems, .zip for Windows
- **Installation Requirements**: No runtime dependencies (static binary, CGO disabled)
- **Update Mechanism**: Package manager updates (Homebrew/Scoop), manual binary replacement
- **Release Verification**: SHA256 checksums provided in checksums.txt for all releases

## Technical Requirements & Constraints

### Performance Requirements
- **Startup time**: <100ms for cached operations
- **Memory usage**: <50MB during normal operations
- **Response time**: <500ms for all credential operations
- **Binary size**: <20MB for cross-platform compatibility

### Compatibility Requirements
- **Platform Support**: Windows (amd64, arm64), macOS (amd64, arm64 + universal binary), Linux (amd64, arm64)
- **Go Version**: Go 1.25.1 (specified in go.mod, tested in CI with 1.25)
- **Standards Compliance**: NIST encryption standards, OWASP secure coding practices
- **Build Tags**: netgo (for static linking without CGO)

### Security & Compliance
- **Encryption**: AES-256-GCM with cryptographically secure key derivation (PBKDF2, 100k iterations)
- **Key Management**: Never store master passwords in plaintext, secure memory clearing
- **File Permissions**: Vault files created with 600 permissions on Unix (macOS/Linux), Windows ACLs on Windows
- **Platform-Specific Security**:
  - macOS: App-level keychain isolation
  - Windows: User-level credential isolation
  - Linux: D-Bus Secret Service integration
- **Threat Model**: Protection against local file access, memory dumps, and weak passwords
- **Defense in Depth**: Encryption is primary security layer, file permissions provide additional protection

### Scalability & Reliability
- **Expected Load**: Single-user, local operations with hundreds of credentials
- **Availability**: Offline-first design, no network dependencies
- **Vault Size**: Efficient handling of vaults up to 10MB (thousands of credentials)

## Technical Decisions & Rationale

### Decision Log
1. **Go Language Selection**:
   - Strong cryptography standard library
   - Excellent cross-compilation support
   - Single binary distribution simplicity
   - Memory safety and performance

2. **AES-256-GCM Encryption**:
   - NIST recommended authenticated encryption
   - Built-in integrity verification
   - Resistance to padding oracle attacks
   - Standard library implementation

3. **Cobra CLI Framework**:
   - Industry standard (used by kubectl, docker, gh)
   - Excellent help system and command organization
   - Automatic completion and validation
   - Consistent with Go ecosystem tools

4. **Local File Storage**:
   - Offline-first approach for privacy
   - No cloud dependencies or attack surface
   - User control over data location
   - Simple backup and migration

5. **System Keychain Integration**:
   - Native OS credential storage when available
   - Graceful fallback to password prompts
   - Enhanced user experience for daily workflows
   - OS-level security protections
   - Unique differentiator among pure CLI password managers

6. **Script-Friendly Output Design**:
   - `--quiet` flag for clean output (no prompts or formatting)
   - `--field` flag for extracting specific credential fields
   - `--json` flag for structured output
   - Enables shell integration: `$env:API_KEY=$(pass-cli get service -q)`
   - `--no-clipboard` flag to prevent automatic clipboard copying

7. **Automatic Usage Tracking**:
   - Track credential usage based on $PWD (current working directory)
   - No manual flags required (fully automatic)
   - Store location, timestamps, access count, git repo info
   - Enables intelligent warnings on deletion/updates
   - Supports usage analysis and credential rotation insights

8. **Platform-Specific File Security**:
   - Unix (macOS/Linux): 0600 permissions for vault file
   - Windows: User-level access via Windows ACLs
   - Encryption as primary security (defense in depth approach)
   - Accept Windows can't enforce Unix-style permissions

9. **TUI Framework Migration (Bubble Tea → tview)**:
   - **Current State**: Bubble Tea (default) with tview migration in progress
   - **Reason for Migration**: Better component ecosystem and native tree/list views
   - **Bubble Tea (Legacy)**: Elm-inspired Model-Update-View architecture
   - **tview (New)**: Widget-based framework with TreeView, Flex, Modal, Pages components
   - **Migration Strategy**: Dual implementation with feature flag toggle, complete component parity before switchover
   - **Code Structure**: Each component has both Bubble Tea and tview implementations in same files
   - **Recent Progress**: Implemented tview versions of Sidebar (TreeView), StatusBar (TextView), MetadataPanel (Flex), Breadcrumb (TextView), CommandBar (Modal + InputField)

10. **Lipgloss and tcell Styling**:
   - **Lipgloss**: Declarative styling with CSS-like API for Bubble Tea components
   - **tcell**: Direct terminal cell manipulation and color management for tview
   - **Styling Approach**: Consistent color palette defined in styles/theme.go, converted between Lipgloss and tcell color formats
   - **Automatic width/height calculation** with border and padding support
   - **Composable layout primitives**: JoinHorizontal, JoinVertical for Lipgloss; Flex for tview
   - **Terminal color profile detection** for broad compatibility

## Known Limitations

- **Multi-user Support**: Single-user design, no concurrent access protection
- **Sync Capabilities**: No built-in synchronization across devices (future enhancement)
- **Audit Logging**: Basic operation logging only, no comprehensive audit trail
- **Plugin System**: Monolithic design, no extensibility via plugins (acceptable for v1.0)
- **TUI Framework Migration**: Active migration from Bubble Tea to tview in progress
  - Current default: Bubble Tea (useBubbleTea=true in tui.go)
  - Some components have dual implementations (increases code size temporarily)
  - Full migration expected to reduce dependencies and improve component reusability
  - No user-facing impact until migration is complete and toggle is switched