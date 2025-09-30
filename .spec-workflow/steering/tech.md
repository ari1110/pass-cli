# Technology Stack

## Project Type
Command-line interface (CLI) tool for secure credential management. Single-binary application designed for local development environments with cross-platform compatibility.

## Core Technologies

### Primary Language(s)
- **Language**: Go 1.21+ (modern Go with enhanced security features)
- **Runtime/Compiler**: Go compiler with CGO disabled for static linking
- **Language-specific tools**: Go modules for dependency management, built-in toolchain for cross-compilation

### Key Dependencies/Libraries
- **github.com/spf13/cobra v1.10.1**: CLI framework and command management
- **github.com/spf13/viper v1.21.0**: Configuration management and file handling
- **github.com/zalando/go-keyring v0.2.6**: Cross-platform system keychain integration
- **golang.org/x/crypto v0.42.0**: Extended cryptographic functions (PBKDF2)
- **Standard library crypto packages**: AES-256-GCM implementation, secure random generation

### Application Architecture
Layered architecture with clear separation of concerns:
- **CLI Layer**: Command interface using Cobra framework with script-friendly output modes
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
- **Static Analysis**: golangci-lint v2.5.0 (comprehensive linter suite)
- **Formatting**: goimports (automatic import management and code formatting)
- **Testing Framework**: Go's built-in testing package with table-driven tests
- **Documentation**: Go doc comments and README.md with usage examples

### Version Control & Collaboration
- **VCS**: Git with conventional commit messages
- **Branching Strategy**: GitHub Flow with feature branches and pull requests
- **Code Review Process**: Required reviews for main branch, automated CI checks

## Deployment & Distribution
- **Target Platform(s)**: Windows 10+, macOS 10.15+, Linux distributions with glibc 2.17+
- **Distribution Method**:
  - GitHub Releases with automated binary builds
  - Homebrew formula for macOS/Linux
  - Scoop manifest for Windows
  - Direct binary download
- **Installation Requirements**: No runtime dependencies (static binary)
- **Update Mechanism**: Package manager updates, manual binary replacement

## Technical Requirements & Constraints

### Performance Requirements
- **Startup time**: <100ms for cached operations
- **Memory usage**: <50MB during normal operations
- **Response time**: <500ms for all credential operations
- **Binary size**: <20MB for cross-platform compatibility

### Compatibility Requirements
- **Platform Support**: Windows (amd64, arm64), macOS (amd64, arm64), Linux (amd64, arm64)
- **Go Version**: Minimum Go 1.21 for security and performance features
- **Standards Compliance**: NIST encryption standards, OWASP secure coding practices

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

## Known Limitations

- **Multi-user Support**: Single-user design, no concurrent access protection
- **Sync Capabilities**: No built-in synchronization across devices (future enhancement)
- **Audit Logging**: Basic operation logging only, no comprehensive audit trail
- **Plugin System**: Monolithic design, no extensibility via plugins (acceptable for v1.0)