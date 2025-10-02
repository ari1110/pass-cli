# Requirements Document - CI Quality Fixes

## Introduction

The GitHub Actions CI/CD workflows are currently failing due to linting errors, test failures, and security scan issues. These failures are blocking automated releases and preventing the Homebrew/Scoop automation from functioning. This spec addresses all CI/workflow failures to restore full automation and ensure code quality standards.

## Alignment with Product Vision

Reliable CI/CD is critical for maintaining the production-ready quality standard established in the pass-cli-release spec. Automated testing and releases enable rapid iteration while maintaining security and code quality. This directly supports the goal of delivering a professional, enterprise-grade password manager.

## Requirements

### FR-1: Fix Linting Errors

**User Story:** As a developer, I want golangci-lint to pass without errors, so that code quality standards are enforced automatically

**Current Issue:** `Error: can't load config: the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.25.1)`

The project uses Go 1.25.1, and the workflow was using golangci-lint with `version: latest`, which pulled a version built with Go 1.24.

**Root Cause:** golangci-lint v2.4.0+ supports Go 1.25 when built with Go 1.25. Official binaries since v2.4.0 are built with Go 1.25. The issue is that `version: latest` in the GitHub Action was not pulling the correct version.

**Solution:** Pin golangci-lint to v2.5 (or later) explicitly in workflow files, which ensures a version built with Go 1.25+ is used.

#### Acceptance Criteria

1. WHEN golangci-lint runs THEN it SHALL complete with exit code 0 (no errors)
2. WHEN code is pushed to main THEN the CI lint job SHALL pass successfully
3. WHEN a release tag is created THEN the release workflow lint step SHALL pass successfully
4. IF linting errors exist THEN they SHALL be fixed or explicitly exempted with justification
5. WHEN Go version is updated THEN golangci-lint version SHALL be pinned to a compatible version

### FR-2: Fix Windows Test Failures

**User Story:** As a developer, I want all tests to pass on Windows, so that cross-platform compatibility is verified

#### Acceptance Criteria

1. WHEN unit tests run on windows-latest THEN all tests SHALL pass
2. WHEN integration tests run on windows-latest THEN all tests SHALL pass
3. IF platform-specific behavior exists THEN it SHALL be properly isolated and tested
4. WHEN tests run with race detector on Windows THEN no race conditions SHALL be detected

### FR-3: Resolve Security Scan Issues

**User Story:** As a security-conscious developer, I want gosec to run without blocking errors, so that security issues are caught early

**Current Issues (4 findings):**
1. **G407 (HIGH)** - `internal/crypto/crypto.go:80`: Use of hardcoded IV/nonce (false positive - nonce is randomly generated)
2. **G304 (MEDIUM)** - `internal/storage/storage.go:348`: Potential file inclusion via variable (user-controlled path)
3. **G304 (MEDIUM)** - `internal/storage/storage.go:323`: Potential file inclusion via variable (user-controlled path)
4. **G304 (MEDIUM)** - `internal/storage/storage.go:271`: Potential file inclusion via variable (user-controlled path)

#### Acceptance Criteria

1. WHEN gosec runs THEN it SHALL complete without critical findings
2. IF gosec produces warnings THEN they SHALL be reviewed and either fixed or documented as accepted risks
3. WHEN SARIF upload occurs THEN it SHALL succeed or be gracefully skipped
4. IF SARIF format issues exist THEN the security scan configuration SHALL be updated to produce valid output
5. WHEN false positives are identified THEN they SHALL be suppressed with `#nosec` comments and documented justification

### FR-4: Enable Successful Automated Releases

**User Story:** As a release manager, I want GitHub Actions to automatically release new versions, so that Homebrew/Scoop are updated without manual intervention

#### Acceptance Criteria

1. WHEN a version tag is pushed THEN the release workflow SHALL complete successfully
2. WHEN the release workflow completes THEN homebrew-tap SHALL be updated automatically
3. WHEN the release workflow completes THEN scoop-bucket SHALL be updated automatically
4. WHEN the release workflow completes THEN GitHub Releases SHALL contain all platform binaries and checksums
5. IF workflow steps fail THEN clear error messages SHALL indicate the problem

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Lint fixes should not alter functional behavior
- **Modular Design**: Test fixes should be isolated to specific test files
- **Dependency Management**: No new dependencies unless absolutely necessary
- **Clear Interfaces**: Workflow configurations should be well-documented

### Performance
- CI workflows should complete within 5 minutes for standard commits
- Release workflows should complete within 10 minutes
- Test suites should run efficiently with proper parallelization

### Security
- All security findings must be addressed or explicitly accepted
- gosec must run without critical vulnerabilities
- No secrets or sensitive data in workflow logs

### Reliability
- CI must pass consistently across all platforms (Windows, macOS, Linux)
- Flaky tests must be identified and fixed
- Workflow failures must provide actionable error messages

### Maintainability
- Workflow configurations should use pinned action versions
- Lint rules should be documented and justifiable
- Test code should follow the same quality standards as production code
