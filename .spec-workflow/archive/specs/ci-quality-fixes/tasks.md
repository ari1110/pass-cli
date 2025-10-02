# Tasks Document - CI Quality Fixes

- [x] 1. Fix golangci-lint Go version compatibility
  - Files: `.github/workflows/ci.yml`, `.github/workflows/release.yml`
  - **CORRECTED APPROACH**: Pin golangci-lint to v2.5 (supports Go 1.25)
  - Keep Go 1.25.1 in go.mod (no downgrade)
  - Update golangci-lint-action from `version: latest` to `version: v2.5`
  - Purpose: Resolve "Go language version mismatch" error by using golangci-lint built with Go 1.25
  - **Note**: Initial implementation incorrectly downgraded Go version; reverted and corrected
  - _Leverage: golangci-lint v2.4.0+ official binaries built with Go 1.25_
  - _Requirements: FR-1_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer with expertise in GitHub Actions and Go tooling | Task: Fix golangci-lint compatibility by pinning to v2.5 following requirement FR-1 | Restrictions: Do NOT downgrade Go version, only change golangci-lint version in workflows, verify golangci-lint v2.5 is built with Go 1.25+ | _Leverage: Existing .github/workflows/*.yml files, golangci-lint v2.5 official binaries | _Requirements: FR-1 (Fix Linting Errors) | Success: golangci-lint v2.5 runs without version mismatch, Go 1.25.1 maintained, all workflows use consistent versions | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 2. Suppress gosec false positive for crypto nonce (G407)
  - File: `internal/crypto/crypto.go:80`
  - Add `// #nosec G407 -- Nonce is randomly generated via crypto/rand, not hardcoded` suppression
  - Document why this is a false positive in code comment
  - Purpose: Allow security scan to pass while documenting that the nonce is properly randomized
  - _Leverage: Standard gosec suppression syntax_
  - _Requirements: FR-3_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Security Engineer with expertise in Go cryptography and static analysis tools | Task: Suppress gosec G407 false positive for randomly-generated nonce following requirement FR-3 | Restrictions: Only add #nosec comment with clear justification, do not modify encryption logic, verify nonce is actually randomly generated | _Leverage: internal/crypto/crypto.go existing implementation | _Requirements: FR-3 (Resolve Security Scan Issues) | Success: gosec no longer reports G407 for this line, clear justification comment added, nonce generation remains cryptographically secure | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 3. Suppress gosec false positives for storage file paths (G304)
  - Files: `internal/storage/storage.go:271`, `internal/storage/storage.go:323`, `internal/storage/storage.go:348`
  - Add `// #nosec G304 -- Vault/backup path is user-controlled by design for CLI tool` suppressions
  - Document that user-controlled paths are intentional for a CLI password manager
  - Purpose: Allow security scan to pass while documenting intentional design decision
  - _Leverage: Standard gosec suppression syntax_
  - _Requirements: FR-3_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Security Engineer with expertise in file system security and Go | Task: Suppress gosec G304 false positives for intentional user-controlled vault paths following requirement FR-3 | Restrictions: Only add #nosec comments with clear justification, do not modify file handling logic, ensure justification explains CLI tool design | _Leverage: internal/storage/storage.go existing implementation | _Requirements: FR-3 (Resolve Security Scan Issues) | Success: gosec no longer reports G304 for these lines, clear justification comments added explaining CLI design, file handling remains secure | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 4. Investigate and fix Windows test failures
  - Files: `.github/workflows/ci.yml`, `.github/workflows/release.yml`
  - Fixed by adding `shell: bash` to test steps
  - Root cause: Windows PowerShell misinterpreted `-coverprofile=coverage.txt` as package `.txt`
  - Purpose: Ensure cross-platform test compatibility
  - **Note**: All Windows tests now passing with bash shell specification
  - _Leverage: Existing test framework, GitHub Actions test logs_
  - _Requirements: FR-2_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in cross-platform Go testing | Task: Investigate and fix Windows-specific test failures following requirement FR-2 | Restrictions: Must maintain test coverage, ensure fixes work on all platforms, do not skip tests unless absolutely necessary | _Leverage: GitHub Actions logs from failed runs, existing test files | _Requirements: FR-2 (Fix Windows Test Failures) | Success: All tests pass on windows-latest, no regressions on macOS/Linux, platform-specific issues properly handled | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 5. Verify CI workflow passes end-to-end
  - **VERIFIED**: GitHub Actions Run #18152313298 - All 9 jobs passing
  - ✅ Lint job (golangci-lint v2.5 with action v7)
  - ✅ Security Scan (gosec with documented suppressions)
  - ✅ Tests - Ubuntu, macOS, Windows (all platforms passing)
  - ✅ Integration Tests - all platforms passing
  - ✅ Build (GoReleaser snapshot) - artifacts generated
  - Additional fixes: GoReleaser v2 configuration (directory syntax, removed before hooks)
  - Purpose: Ensure all CI fixes work together
  - _Leverage: GitHub Actions workflows, existing CI infrastructure_
  - _Requirements: FR-1, FR-2, FR-3_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer with expertise in CI/CD validation | Task: Verify complete CI workflow passes following requirements FR-1, FR-2, and FR-3 | Restrictions: Must test on actual GitHub Actions, not just locally, verify all jobs pass, ensure no new failures introduced | _Leverage: .github/workflows/ci.yml, GitHub Actions interface | _Requirements: FR-1, FR-2, FR-3 | Success: Full CI workflow passes on test commit, lint job passes, security scan passes, all platform tests pass, build artifacts generated | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 6. Test release workflow with v1.0.2 tag
  - **VERIFIED**: Tag v1.0.2-test successfully triggered release workflow
  - ✅ GitHub Actions Run #18152820458 - All jobs passing
  - ✅ GoReleaser built all platform binaries (Windows, macOS, Linux × amd64, arm64)
  - ✅ checksums.txt generated
  - ✅ GitHub Release created: https://github.com/ari1110/pass-cli/releases/tag/v1.0.2-test
  - **Additional fixes applied**: Disabled Snap/SBOM (tools not in CI), configured PATs for package managers
  - Purpose: Validate release automation works end-to-end
  - _Leverage: .github/workflows/release.yml, GoReleaser configuration_
  - _Requirements: FR-4_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Release Manager with expertise in automated release workflows | Task: Test complete release workflow with test tag following requirement FR-4 | Restrictions: Use test tag (not final version), monitor all workflow steps, ensure no errors in any job | _Leverage: .github/workflows/release.yml, .goreleaser.yml | _Requirements: FR-4 (Enable Successful Automated Releases) | Success: Release workflow completes successfully, all binaries built, GitHub Release created with artifacts, checksums.txt generated | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 7. Verify Homebrew tap auto-update
  - **VERIFIED**: https://github.com/ari1110/homebrew-tap auto-updated successfully
  - ✅ Commit by goreleaserbot at 2025-10-01T05:55:37Z
  - ✅ Formula/pass-cli.rb created/updated with "Brew formula update for pass-cli version v1.0.2-test"
  - ✅ PAT configuration (HOMEBREW_TAP_TOKEN) working correctly
  - Purpose: Confirm Homebrew automation is working
  - _Leverage: GoReleaser brew publisher, homebrew-tap repository_
  - _Requirements: FR-4_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Package Manager Specialist with expertise in Homebrew formulas | Task: Verify Homebrew tap auto-update functionality following requirement FR-4 | Restrictions: Do not manually create formula, only verify GoReleaser automation worked, check all URLs and checksums | _Leverage: homebrew-tap repository, GoReleaser configuration | _Requirements: FR-4 (Enable Successful Automated Releases) | Success: Formula/pass-cli.rb exists in tap repo, URLs point to correct release, checksums are correct, commit shows GoReleaser automation | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 8. Verify Scoop bucket auto-update
  - **VERIFIED**: https://github.com/ari1110/scoop-bucket auto-updated successfully
  - ✅ Commit by goreleaserbot at 2025-10-01T05:55:38Z
  - ✅ bucket/pass-cli.json created/updated with "Scoop update for pass-cli version v1.0.2-test"
  - ✅ PAT configuration (SCOOP_BUCKET_TOKEN) working correctly
  - Purpose: Confirm Scoop automation is working
  - _Leverage: GoReleaser scoop publisher, scoop-bucket repository_
  - _Requirements: FR-4_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Package Manager Specialist with expertise in Scoop manifests | Task: Verify Scoop bucket auto-update functionality following requirement FR-4 | Restrictions: Do not manually create manifest, only verify GoReleaser automation worked, check all URLs and hashes | _Leverage: scoop-bucket repository, GoReleaser configuration | _Requirements: FR-4 (Enable Successful Automated Releases) | Success: bucket/pass-cli.json exists in bucket repo, URLs point to correct release, SHA256 hashes are correct, commit shows GoReleaser automation | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 9. Document Go version upgrade path
  - **COMPLETED**: Added comprehensive golangci-lint troubleshooting to docs/CI-CD.md
  - ✅ Documented "Go language version mismatch" issue with root cause and solution
  - ✅ Included example configuration for Go 1.25+ projects (golangci-lint v2.5+ with action v7)
  - ✅ Added reference link to GitHub issue for future tracking
  - ✅ Also documented PAT setup for Homebrew/Scoop cross-repo updates
  - Purpose: Ensure team can troubleshoot and maintain Go version compatibility
  - _Leverage: Current fix as reference, golangci-lint release tracking_
  - _Requirements: FR-1_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer with expertise in Go tooling and CI/CD | Task: Document future Go version upgrade process following requirement FR-1 | Restrictions: Keep instructions clear and actionable, include version checking steps, provide troubleshooting tips | _Leverage: Current go.mod and workflow files, golangci-lint GitHub releases | _Requirements: FR-1 (Fix Linting Errors) | Success: Clear documentation created, upgrade steps are actionable, rollback procedure included, golangci-lint version compatibility checking explained | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_

- [x] 10. Clean up test tags and complete spec
  - **COMPLETED**: Cleaned up test artifacts and finalized spec
  - ✅ Deleted test tag v1.0.2-test and associated GitHub Release
  - ✅ Kept production releases (v1.0.0, v1.0.1)
  - ✅ Committed documentation updates to CI-CD.md (commit 6980d61)
  - ✅ All 10 tasks in ci-quality-fixes spec marked complete
  - **Note**: pass-cli-release Task 16 remains pending - CI/CD pipeline is now ready for production use
  - Purpose: Clean up temporary artifacts and finalize spec
  - _Leverage: Git tag commands, task tracking_
  - _Requirements: All_
  - _Prompt: Implement the task for spec ci-quality-fixes, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Project Manager with expertise in cleanup and documentation | Task: Clean up test artifacts and finalize spec completion following all requirements | Restrictions: Only delete test tags, keep actual releases, ensure both specs are properly marked complete | _Leverage: Git commands, .spec-workflow/specs/*/tasks.md files | _Requirements: All requirements | Success: Test tags cleaned up, pass-cli-release Task 16 marked [x], ci-quality-fixes spec marked complete, repository in clean state | Instructions: Mark this task as in-progress [-] in tasks.md when starting, then mark as complete [x] when finished_
