# Product Overview

## Product Purpose
Pass-CLI is a secure, cross-platform command-line password and API key manager designed specifically for developers. It solves the critical problem of securely storing, managing, and accessing credentials in local development environments without relying on cloud services or exposing sensitive data in configuration files.

## Target Users
**Primary Users**: Software developers, DevOps engineers, and security-conscious technologists who:
- Need to manage multiple API keys, database passwords, and service credentials
- Work in local development environments with sensitive credentials
- Prefer command-line tools for workflow integration
- Value security, privacy, and offline-first solutions
- Want lightweight alternatives to full-featured password managers

**Pain Points Addressed**:
- Storing credentials in plaintext configuration files
- Sharing API keys through insecure channels
- Managing credentials across multiple projects and environments
- Remembering complex passwords and rotating API keys
- Lack of CLI-native credential management tools

## Key Features

1. **Local Encrypted Storage**: AES-256-GCM encryption with PBKDF2 key derivation (100,000 iterations) for maximum security
2. **System Keychain Integration**: Automatic master password storage in Windows Credential Manager, macOS Keychain, or Linux Secret Service (unique differentiator for CLI password managers)
3. **Interactive TUI Dashboard**: Full-featured terminal UI with hierarchical category tree, credential browsing, add/edit forms, and real-time updates via `pass-cli tui` command
4. **6-Field Credential Model**: Comprehensive credential storage with Service, Username, Password, Category, URL, and Notes fields for rich metadata
5. **Fast CLI Commands**: Script-friendly commands for automation and workflow integration with sub-100ms credential retrieval
6. **Script-Friendly Output**: Support for shell integration with `--quiet`, `--field`, and `--masked` flags for use in scripts like `$env:API_KEY=$(pass-cli get service -q)`
7. **Automatic Usage Tracking**: Intelligent tracking of where credentials are used based on working directory, with no manual flags required
8. **Cross-Platform Compatibility**: Single binary distribution for Windows, macOS (Intel/ARM), and Linux (amd64/arm64)
9. **Offline Operation**: No cloud dependencies, works completely offline
10. **Password Generation**: Cryptographically secure password generation with customizable length and character options
11. **Clipboard Integration**: Automatic credential copying with 30-second timeout and optional `--no-clipboard` flag to disable
12. **Masked Password Display**: Optional `--masked` flag to display passwords as asterisks for additional security
13. **Package Manager Ready**: Homebrew and Scoop manifests prepared for distribution (pending publication to official repositories)

## Business Objectives

- **Developer Productivity**: Reduce friction in credential management workflows
- **Security Enhancement**: Eliminate plaintext credential storage in development environments
- **Community Impact**: Provide open-source alternative to proprietary credential managers
- **Learning Platform**: Demonstrate Go security best practices and CLI design patterns

## Success Metrics

- **Adoption**: Target 1,000+ GitHub stars within 6 months
- **Distribution**: Homebrew and Scoop manifests ready for submission to official repositories
- **Security**: Zero known vulnerabilities in encryption implementation (AES-256-GCM with authentication)
- **Usability**: Sub-100ms response time for credential operations (verified in testing)
- **Quality**: Comprehensive test coverage with automated CI/CD testing
- **Feature Completeness**: ✅ TUI dashboard fully implemented with all core features (v0.1.0+)
- **User Experience**: Dual interface support (CLI for scripting, TUI for visual management) providing workflow flexibility

## Product Principles

1. **Security First**: Never compromise on cryptographic security or data protection
2. **Developer Experience**: Design for speed, simplicity, and CLI integration
3. **Privacy by Design**: Local-only storage with no telemetry or data collection
4. **Open Source**: Transparent, auditable code with community contributions
5. **Minimal Dependencies**: Lean binary with essential functionality only

## Monitoring & Visibility

- **CLI Output Modes**:
  - **Normal Mode**: Formatted table display with credential details (all 6 fields)
  - **Quiet Mode**: Password-only output for scripting (`--quiet`)
  - **Field Mode**: Extract specific fields like username, URL, category (`--field`)
  - **Masked Mode**: Display passwords as asterisks (`--masked`)
- **TUI Dashboard**:
  - **Hierarchical View**: Visual category tree with credential counts
  - **Real-Time Detail View**: Credential details update on navigation
  - **Search & Filter**: Quick credential location by service name or category
  - **Visual Indicators**: Password masking toggle, clipboard copy confirmations
- **Usage Tracking**:
  - Automatic tracking of credential access by working directory
  - Last used timestamps for all credentials
  - Usage count statistics
  - Unused credential detection (`--unused --days N`)
- **Key Metrics**:
  - Total credential count
  - Credentials organized by category (hierarchical tree navigation in TUI)
  - Last accessed timestamps
  - Creation and modification dates
  - Credential metadata (Category, URL, Notes fields)

## Competitive Landscape

### Existing Solutions Analysis

**gopass**
- Free, open source, uses GPG + Git
- No OS keychain integration
- Asymmetric encryption (GPG hybrid approach)
- Strong sync capabilities via Git

**Bitwarden CLI**
- Free tier + premium, requires server (cloud or self-hosted)
- No native CLI keychain integration
- Excellent web/mobile integration
- Team collaboration features

**KeePassXC**
- Free, single .kdbx file
- GUI can use keychain but CLI cannot
- Strong security track record
- Cross-platform compatibility

**1Password CLI**
- Paid only, has keychain integration
- Requires desktop app for full functionality
- Excellent UX and team features
- Enterprise-grade security

**Chezmoi**
- Dotfiles manager with keychain integration
- Not a pure password manager
- Uses go-keyring library

**Pass-CLI Positioning**: A pure CLI password manager with **dual interfaces** (CLI + TUI), native OS keychain integration, symmetric encryption (AES-256-GCM), and automatic usage tracking. Unique features include:
- **Interactive TUI dashboard** with hierarchical category navigation and responsive layouts (competitors lack visual TUI interfaces)
- **6-field credential model** (Service, Username, Password, Category, URL, Notes) for rich metadata management
- Sub-100ms credential retrieval optimized for developer workflows
- Automatic directory-based usage tracking without manual flags
- Multiple output modes (quiet, field-specific, masked, clipboard) optimized for shell scripting and CI/CD pipelines
- **Single binary distribution** with no runtime dependencies (vs. gopass requiring GPG, Bitwarden CLI requiring server)

## Future Vision

Pass-CLI aims to become the standard CLI credential manager for developers, setting the benchmark for security, usability, and integration capabilities in the developer tools ecosystem.

### Completed Features (v0.1.0+)
- **Interactive TUI Dashboard** ✅ COMPLETE: Full-featured terminal UI accessible via `pass-cli tui` command
  - **Hierarchical Category Tree**: Navigate credentials by category with expandable/collapsible tree structure
  - **Credential Table**: Browse and search credentials with real-time filtering
  - **Detail View**: View credential details with password masking toggle and clipboard copy support
  - **Add/Edit Forms**: Complete 6-field forms with validation (Service, Username, Password, Category, URL, Notes)
  - **Visual Theme**: Dracula color palette with rounded borders and active/inactive panel distinction
  - **Responsive Layout**: Adaptive multi-panel layout with breakpoints (80/120 columns)
  - **Keyboard Navigation**: Comprehensive shortcuts with context-aware status bar hints
  - **Focus Management**: Layout-aware tab cycling across panels
  - **Modal System**: Delete confirmations and command palette for quick actions
  - **Comprehensive Testing**: Unit tests and manual testing suite covering all TUI features
- **6-Field Credential Model** ✅ COMPLETE: Extended credential storage with Category, URL, and Notes fields
  - CLI commands updated: `add --category --url`, `update --category --url`, `get` displays all fields
  - Full TUI integration with forms and detail view

### Planned Enhancements
- **JSON Output Format**: Add `--json` flag for structured output to enable parsing with jq and other tools
- **Export/Import**: Credential backup and migration functionality for vault portability
- **Browser Extensions**: Integration with web browsers for auto-fill capabilities
- **Remote Sync**: Secure vault synchronization across multiple development machines
- **Team Collaboration**: Shared credential vaults with granular access controls
- **Two-Factor Authentication**: TOTP support for credentials requiring 2FA
- **Audit Logging**: Comprehensive access logs for security compliance
- **Advanced Shell Completion**: Enhanced completion for bash, zsh, fish, and PowerShell