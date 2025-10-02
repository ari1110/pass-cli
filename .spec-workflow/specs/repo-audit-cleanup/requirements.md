# Requirements Document

## Introduction

The pass-cli repository has grown organically with various implementation documents, test files, and package manager configurations added at different times. This has resulted in inconsistent file organization, root-level clutter, and misalignment with the documented structure in `structure.md`.

This spec addresses systematic repository organization to ensure:
- All files follow documented structure conventions
- Root directory is clean with only essential project files
- Documentation is properly organized and discoverable
- Package manager configurations follow consistent patterns
- The actual structure matches what's documented in steering docs

## Alignment with Product Vision

From `product.md`:
- **Product Principle #2**: "Developer Experience: Design for speed, simplicity, and CLI integration" - A well-organized repo improves developer onboarding and contribution experience
- **Product Principle #4**: "Open Source: Transparent, auditable code with community contributions" - Clear organization makes the codebase more accessible to contributors

From `tech.md`:
- **Go Best Practices**: "Follow established Go idioms and project patterns" - Proper file organization is a Go ecosystem standard
- **Code Quality Tools**: Organized structure makes it easier to run linters, tests, and maintain the codebase

From `structure.md`:
- **Project Structure**: Documents expected organization - actual structure should match documented structure
- **File Organization Principles**: "One primary type per file", "Related functionality grouped" - These principles should extend to the entire repo structure

This audit ensures the repository reflects the professionalism and quality standards expected of a developer tool.

## Requirements

### Requirement 1: Root Directory Cleanup

**User Story:** As a developer browsing the repository, I want the root directory to contain only essential project files, so that I can quickly understand the project structure without digging through implementation details.

#### Acceptance Criteria

1. WHEN viewing the root directory THEN it SHALL contain only: README.md, LICENSE, go.mod, go.sum, main.go, Makefile, .gitignore, and essential config files (.goreleaser.yml, .github/)
2. WHEN implementation/development documentation exists at root THEN it SHALL be moved to appropriate docs/ subdirectories
3. IF a file is a temporary implementation artifact (DASHBOARD_IMPLEMENTATION_SUMMARY.md, DASHBOARD_TESTING_CHECKLIST.md, KEYBINDINGS_AUDIT.md) THEN it SHALL be moved to docs/development/ or docs/archive/
4. WHEN test scripts exist at root (test-tui.bat) THEN they SHALL be moved to test/ or a scripts/ directory
5. WHEN the root directory is listed THEN developers SHALL see a clean, standard Go project structure

### Requirement 2: Documentation Structure Alignment

**User Story:** As a developer contributing to pass-cli, I want documentation organized by purpose and lifecycle, so that I can find relevant docs quickly and understand their current status.

#### Acceptance Criteria

1. WHEN documentation is categorized THEN it SHALL follow structure: docs/ for current, docs/archive/ for historical, docs/development/ for implementation notes
2. WHEN a document is implementation-specific (dashboard summaries, testing checklists) THEN it SHALL be in docs/development/
3. IF a document is outdated or superseded THEN it SHALL be moved to docs/archive/ with a README explaining why
4. WHEN docs/development/ exists THEN it SHALL contain a README explaining: "Development notes and implementation tracking documents. Not user-facing."
5. WHEN docs are organized THEN each subdirectory SHALL have a README listing contents and purpose

### Requirement 3: Package Manager Configuration Consistency

**User Story:** As a maintainer managing package releases, I want all package manager configurations in a consistent location, so that I can update them together during releases.

#### Acceptance Criteria

1. WHEN package manager configs are organized THEN they SHALL follow ONE of these patterns:
   - Pattern A: All in root (homebrew/, scoop/, snap/, winget/) - current mix
   - Pattern B: All in manifests/ (manifests/homebrew/, manifests/scoop/, manifests/snap/, manifests/winget/)
   - Pattern C: Keep platform-native at root (homebrew/, scoop/), platform-agnostic in manifests/ (snap/, winget/)
2. WHEN the chosen pattern is applied THEN ALL package configs SHALL follow it consistently
3. IF configs are moved THEN .goreleaser.yml, CI workflows, and docs SHALL be updated to reflect new paths
4. WHEN reviewing package configs THEN the organization SHALL be immediately obvious and documented in structure.md
5. WHEN adding a new package manager THEN the pattern SHALL be clear from existing structure

### Requirement 4: Structure.md Synchronization

**User Story:** As a developer reading the steering docs, I want structure.md to accurately reflect the actual repository structure, so that I can navigate the codebase confidently.

#### Acceptance Criteria

1. WHEN structure.md documents a directory THEN that directory SHALL exist in the actual repository
2. WHEN the actual repository has directories not in structure.md THEN structure.md SHALL be updated to document them
3. IF a directory serves a specific purpose (manifests/, homebrew/, scripts/) THEN structure.md SHALL explain its purpose and organization
4. WHEN structure.md is updated THEN it SHALL include file counts and key files for major directories
5. WHEN the audit is complete THEN a developer SHALL be able to navigate the repo using only structure.md as a guide

### Requirement 5: Test Organization Review

**User Story:** As a developer running tests, I want test files organized by type and scope, so that I can run appropriate test suites efficiently.

#### Acceptance Criteria

1. WHEN test files are organized THEN they SHALL be clearly categorized: unit tests (*_test.go adjacent to source), integration tests (test/)
2. WHEN test utilities exist (test-tui.bat) THEN they SHALL be in test/ directory with clear naming
3. IF test data directories exist (test-vault/) THEN they SHALL be documented in test/README.md
4. WHEN test organization is reviewed THEN test/README.md SHALL document: test types, how to run them, test data setup
5. WHEN developers run tests THEN they SHALL know exactly which test suite they're executing (unit vs integration)

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Each directory has one clear purpose (docs for documentation, test for tests, etc.)
- **Modular Design**: Package manager configs isolated from core code, documentation separated from implementation
- **Dependency Management**: File moves must not break import paths or build processes
- **Clear Interfaces**: Directory purposes documented in structure.md and local README files

### Performance
- No performance impact from reorganization (Go imports use module paths, not file locations for most files)
- CI/CD pipeline updates must maintain or improve build times

### Reliability
- All file moves must be verified with successful builds and test runs
- Package manager configuration paths must be updated atomically (update code, configs, and CI together)
- No broken links in documentation after reorganization

### Usability
- Root directory clean and navigable (< 15 items at root level excluding hidden files/dirs)
- Documentation discoverable with consistent naming and organization
- New contributors can find relevant docs within 30 seconds of landing in the repo

### Maintainability
- structure.md becomes single source of truth for repository organization
- Adding new documentation or configs follows established patterns
- File organization supports future growth without requiring restructuring
