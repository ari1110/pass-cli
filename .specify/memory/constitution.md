# Pass-CLI Constitution

<!--
Sync Impact Report:
Version Change: Initial constitution (no previous version) → 1.0.0
Modified Principles: N/A (initial creation)
Added Sections: All sections (Security First, Spec-Driven Development, Testing Discipline, Layered Architecture, Code Quality Standards, Development Workflow, Governance)
Removed Sections: N/A
Templates Requiring Updates:
  ✅ plan-template.md - Constitution Check section aligns with principles
  ✅ spec-template.md - Requirements structure supports security and testing principles
  ✅ tasks-template.md - Task organization supports layered architecture and test-first approach
Follow-up TODOs: None - all placeholders resolved
-->

## Core Principles

### I. Security First (NON-NEGOTIABLE)

Pass-CLI is a **credential management tool**. Security MUST never be compromised for convenience, performance, or feature velocity.

**Mandatory Security Rules**:
- MUST use AES-256-GCM encryption with authenticated encryption for all credential storage
- MUST use PBKDF2 key derivation with minimum 100,000 iterations
- MUST use cryptographically secure random generation (`crypto/rand`) for all random data (IVs, salts, passwords)
- MUST clear sensitive data from memory after use
- MUST set vault file permissions to 600 (Unix) or user-only ACLs (Windows)
- MUST validate all inputs before cryptographic operations
- MUST use `golang.org/x/crypto` for extended cryptographic functions

**Prohibited Practices**:
- NEVER store passwords or master passwords in plaintext (logs, errors, debug output, memory dumps)
- NEVER use weak encryption algorithms (DES, RC4, MD5, SHA1 for passwords)
- NEVER skip input validation on credential data
- NEVER expose credentials in error messages or stack traces
- NEVER log sensitive data (passwords, keys, decrypted credentials)

**Rationale**: Security vulnerabilities in a password manager are catastrophic. Users trust this tool with their most sensitive data. Any security compromise undermines the entire product value proposition and user trust.

---

### II. Spec-Driven Development (NON-NEGOTIABLE)

All feature work MUST follow the spec-workflow process: Requirements → Design → Tasks → Implementation.

**Workflow Requirements**:
- MUST use spec-workflow MCP server tools for all features
- MUST read steering documents (product.md, tech.md, structure.md) before creating specs
- MUST request, poll, and delete approvals via dashboard (verbal approval NOT accepted)
- MUST follow specs exactly as written with NO deviations, shortcuts, or reinterpretations
- MUST mark tasks accurately in tasks.md ([ ] pending, [-] in-progress, [x] completed)
- MUST work on ONE spec at a time until completion

**Transparency Requirements**:
- MUST report accurate state of work (no aspirational claims)
- MUST stop immediately if discovering incomplete work and document gaps
- MUST surface spec errors or ambiguities before implementation
- MUST test thoroughly before marking tasks complete

**Rationale**: Spec-driven development ensures deliberate planning, reduces rework, maintains architectural consistency, and provides clear audit trails. Accuracy and transparency prevent compounding errors and enable effective collaboration.

---

### III. Testing Discipline

All production code MUST have corresponding tests. Tests MUST be written before implementation (Test-Driven Development) when feasible.

**Testing Requirements**:
- MUST write unit tests for all new code (`*_test.go` files)
- MUST use table-driven tests for Go (idiomatic pattern)
- MUST write integration tests for cross-layer interactions (in `test/` directory)
- MUST achieve minimum 90% code coverage for new features
- MUST test security-critical paths (encryption, decryption, key derivation)
- MUST test error paths and edge cases

**Test Quality Standards**:
- Tests MUST be deterministic (no flaky tests)
- Tests MUST be isolated (no shared state between tests)
- Tests MUST clean up resources (temp files, test vaults)
- Tests MUST not expose sensitive data (use fixtures, not real credentials)

**Pre-Commit Testing**:
- `go fmt ./...` - Format code
- `go vet ./...` - Static analysis
- `golangci-lint run` - Comprehensive linting
- `gosec ./...` - Security scanning
- `go test ./...` - Run all tests

**Rationale**: Comprehensive testing prevents regressions, validates security properties, documents expected behavior, and enables confident refactoring. Test-driven development catches design issues early.

---

### IV. Layered Architecture

Pass-CLI follows strict layered architecture. Layers MUST only depend on layers below them, never above or across.

**Architecture Layers** (top to bottom):
1. **CLI Layer** (`cmd/`) - Command interface, flags, help text, output formatting
2. **TUI Layer** (`cmd/tui/`) - Interactive terminal UI (tview-based dashboard)
3. **Service Layer** (`internal/vault/`) - Business logic, credential management, usage tracking
4. **Storage Layer** (`internal/storage/`) - Encrypted file operations, atomic writes, backups
5. **Crypto Layer** (`internal/crypto/`) - AES-256-GCM encryption, key derivation, random generation
6. **Keychain Layer** (`internal/keychain/`) - System keychain integration (Windows/macOS/Linux)
7. **Standard Library & External Dependencies**

**Layer Dependency Rules**:
- MUST depend only on layers below (e.g., Service → Storage → Crypto)
- MUST NOT depend on layers above (e.g., Crypto MUST NOT depend on Service)
- MUST NOT bypass layers (e.g., CLI MUST NOT call Crypto directly, use Service)
- MUST use interfaces for cross-layer communication (enables testing and mocking)

**Boundary Enforcement**:
- CLI and TUI share Vault service but have independent presentation logic
- Internal packages (`internal/`) are implementation details, NOT public API
- Core crypto/vault logic is cross-platform, keychain layer handles OS differences

**Rationale**: Layered architecture enforces separation of concerns, enables independent testing, prevents circular dependencies, and allows layer-specific optimizations without affecting other layers.

---

### V. Code Quality Standards

All code MUST meet quality standards before merging. Quality gates are enforced in CI/CD and pre-commit.

**Code Style**:
- MUST follow Go idioms and conventions (effective Go principles)
- MUST use `goimports` and `gofmt` for consistent formatting
- MUST write clear, descriptive variable and function names
- MUST add Go doc comments for all exported types, functions, and methods
- MUST keep functions focused (single responsibility, max 50 lines preferred)
- MUST limit file size (max 500 lines per file, excluding tests)

**Error Handling**:
- MUST handle all errors explicitly (no ignored errors)
- MUST wrap errors with context (`fmt.Errorf("context: %w", err)`)
- MUST avoid panic except for unrecoverable errors
- MUST sanitize errors before displaying to users (no sensitive data leaks)

**Dependencies**:
- MUST minimize external dependencies (lean binary principle)
- MUST use only well-maintained, security-audited libraries
- MUST pin dependency versions in go.mod for reproducible builds
- MUST run `govulncheck` to scan for vulnerable dependencies

**Documentation**:
- MUST document complex algorithms (especially cryptographic operations)
- MUST explain security decisions in comments
- MUST update README.md when adding user-facing features
- MUST update steering docs (tech.md, structure.md) when changing architecture

**Rationale**: Code quality standards ensure maintainability, readability, and long-term project health. Consistent style reduces cognitive load. Comprehensive documentation enables collaboration and auditing.

---

### VI. Cross-Platform Compatibility

Pass-CLI MUST work identically across Windows, macOS, and Linux with no platform-specific degradation.

**Platform Support**:
- Windows 10+ (amd64, arm64)
- macOS 10.15+ (amd64, arm64, universal binary)
- Linux (amd64, arm64) with glibc 2.17+

**Compatibility Requirements**:
- MUST build with CGO_ENABLED=0 (static binary, no C dependencies)
- MUST use cross-platform libraries (go-keyring, clipboard, etc.)
- MUST test on all target platforms in CI/CD (GitHub Actions matrix)
- MUST handle platform-specific features gracefully (e.g., keychain fallback)
- MUST set appropriate file permissions per platform (600 Unix, ACLs Windows)

**Performance Targets** (across all platforms):
- Startup time: <100ms for cached operations
- Credential retrieval: <500ms
- Memory usage: <50MB during normal operations
- Binary size: <20MB

**Rationale**: Users expect consistent behavior across platforms. Static binaries simplify distribution. Platform-specific optimizations (keychain integration) enhance UX without breaking cross-platform guarantees.

---

### VII. Offline-First & Privacy

Pass-CLI MUST operate completely offline with no cloud dependencies or network calls.

**Privacy Requirements**:
- MUST store all data locally (no cloud sync, no telemetry, no analytics)
- MUST never transmit credentials over network
- MUST not phone home for updates, license checks, or feature flags
- MUST not embed tracking pixels or analytics in documentation

**User Control**:
- Users control vault location and backup strategy
- Users control master password (no account creation, no SSO required)
- Users can audit encrypted vault files (documented format)
- Users can export/import credentials (planned feature)

**Rationale**: Privacy is a core product differentiator. Developers need tools that work offline (planes, secure networks, air-gapped environments). No cloud dependencies eliminates attack surface and vendor lock-in.

---

## Development Workflow

### Workflow Stages

**1. Spec Creation**:
- Read steering docs first (product.md, tech.md, structure.md)
- Create requirements.md → request approval → poll status → delete approval
- Create design.md → request approval → poll status → delete approval
- Create tasks.md → request approval → poll status → delete approval

**2. Implementation**:
- Execute tasks systematically in order
- Update task checkboxes: [ ] → [-] → [x]
- Commit frequently (after each task, after each phase, before context switches)
- Test before marking complete

**3. Quality Assurance**:
- Run all pre-commit checks (fmt, vet, lint, security scan, tests)
- Verify security properties (encryption, file permissions, memory clearing)
- Test cross-platform behavior (CI matrix runs automatically)

**4. Documentation**:
- Update README.md for user-facing changes
- Update steering docs (tech.md, structure.md) for architecture changes
- Add inline documentation for complex logic
- Update CHANGELOG.md

---

## Commit Standards

**Commit Frequency**:
- Commit after completing each task
- Commit after completing each spec phase
- Commit before switching tasks or contexts
- Commit when updating steering docs

**Commit Message Format**:
```
<type>: <description>

<body explaining changes>

<optional phase reference>

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Commit Types**:
- `feat:` - New feature or enhancement
- `fix:` - Bug fix
- `refactor:` - Code restructuring without behavior change
- `test:` - Adding or updating tests
- `docs:` - Documentation updates
- `chore:` - Maintenance tasks (dependencies, build config)
- `security:` - Security improvements or fixes

**Rationale**: Frequent commits create audit trails, enable rollback to working states, and demonstrate systematic progress. Conventional commit format enables automated changelog generation and semantic versioning.

---

## Governance

### Constitution Authority

This constitution supersedes all other development practices. When conflicts arise between this document and other guidance, this constitution takes precedence.

### Amendment Process

Constitutional amendments require:
1. Documentation of proposed change with rationale
2. Impact analysis on existing specs and code
3. Dashboard approval via spec-workflow
4. Migration plan for existing code if applicable
5. Version bump per semantic versioning rules

### Version Semantics

Constitution versioning follows semantic versioning:
- **MAJOR**: Backward-incompatible governance changes (principle removal/redefinition)
- **MINOR**: New principles added or materially expanded guidance
- **PATCH**: Clarifications, wording improvements, typo fixes

### Compliance Review

**All code changes MUST verify compliance**:
- Pull request checklist includes constitution compliance
- CI/CD enforces security scanning, testing, and quality gates
- Code reviews verify architectural layer boundaries
- Spec approvals verify adherence to spec-driven workflow

**Complexity Justification**:
- Any complexity added MUST be justified in spec documentation
- Violations of simplicity principles MUST be documented in plan.md Complexity Tracking
- Alternative simpler approaches MUST be documented and explained why rejected

### Runtime Development Guidance

For day-to-day development guidance, consult `CLAUDE.md` (Claude operational guide). That file provides detailed instructions for using spec-workflow tools, handling errors, and following coding standards.

---

**Version**: 1.0.0 | **Ratified**: 2025-10-07 | **Last Amended**: 2025-10-07
