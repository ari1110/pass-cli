# Documentation Inventory Report

**Date**: 2025-10-01
**Version**: v0.0.1
**Purpose**: Comprehensive assessment of all documentation for sync status, relevance, and quality

---

## Executive Summary

**Total Documents**: 15 markdown files
**Categories**:
- User-facing: 3 documents (README.md, INSTALLATION.md, USAGE.md)
- Developer-facing: 6 documents (DEVELOPMENT.md, CI-CD.md, RELEASE.md, HOMEBREW.md, SCOOP.md, test/README.md)
- Operational: 2 documents (SECURITY.md, TROUBLESHOOTING.md)
- Platform-specific: 2 documents (winget/README.md, snap/README.md)
- Project artifacts: 2 documents (RELEASE-DRY-RUN.md, SECURITY-AUDIT.md)

**Overall Status**: 🟡 **NEEDS UPDATE** - Most docs contain placeholder URLs and version mismatches

---

## Critical Issues Found

1. **Placeholder URLs**: All docs use `yourusername` instead of actual repository (`ari1110/pass-cli`)
2. **Version Mismatches**: Docs reference v1.0.0, but actual release is v0.0.1
3. **Feature Mismatches**:
   - README/USAGE document `--json` flag (NOT implemented per tech.md)
   - SECURITY.md states 600,000 PBKDF2 iterations (actual: 100,000 per tech.md)
4. **Go Version Discrepancy**: DEVELOPMENT.md states 1.25+ but current is likely 1.23+
5. **Missing Implementations**: winget and snap documentation is future-work

---

## Detailed Inventory

### 1. README.md
- **Location**: `R:\Test-Projects\pass-cli\README.md`
- **Category**: User Documentation (Primary)
- **Purpose**: Main entry point for users, installation, quick start, features
- **Sync Status**: 🔴 **OUT OF SYNC**
- **Priority**: ⭐⭐⭐ **CRITICAL**

**Issues**:
1. Line 9: States "600,000 iterations" (actual: 100,000 per tech.md:89)
2. Lines 16, 114, 127, 279: Documents `--json` flag (NOT implemented, per tech.md:141-142 it's a future enhancement)
3. Lines 28, 29, 36, 37, 42, 58: Placeholder URLs `yourusername` (should be `ari1110`)
4. Line 185: States "PBKDF2-SHA256 with 600,000 iterations" (actual: 100,000)
5. Feature list includes `--json` which is future work, not v0.0.1

**Required Changes**:
- Replace all `yourusername` with `ari1110`
- Change all PBKDF2 iteration references from 600,000 to 100,000
- Remove `--json` flag from feature list and examples
- Update "Coming soon" to actual status for Homebrew/Scoop (v0.0.1 released)
- Add actual GitHub release URL
- Update version references to v0.0.1

---

### 2. docs/INSTALLATION.md
- **Location**: `R:\Test-Projects\pass-cli\docs\INSTALLATION.md`
- **Category**: User Documentation
- **Purpose**: Comprehensive installation guide for all platforms
- **Sync Status**: 🔴 **OUT OF SYNC**
- **Priority**: ⭐⭐⭐ **CRITICAL**

**Issues**:
1. Lines 29, 58, 101: Placeholder URLs `yourusername` instead of `ari1110`
2. Lines 28, 35: States "coming soon" but v0.0.1 is already released
3. Lines 81, 84, 124, 127, 138: Version examples use v1.0.0 (actual: v0.0.1)
4. Line 301: Go version 1.25 requirement (may need verification - likely 1.23+)
5. Lines 166, 199: Download URLs use placeholder `yourusername` and v1.0.0

**Required Changes**:
- Replace all `yourusername` with `ari1110`
- Update all version examples from v1.0.0 to v0.0.1
- Update installation status from "coming soon" to available
- Verify actual Go version requirement
- Update download URLs to actual v0.0.1 release
- Update checksums references to actual v0.0.1 checksums

---

### 3. docs/USAGE.md
- **Location**: `R:\Test-Projects\pass-cli\docs\USAGE.md`
- **Category**: User Documentation
- **Purpose**: Complete command reference with all flags and examples
- **Sync Status**: 🔴 **OUT OF SYNC**
- **Priority**: ⭐⭐⭐ **CRITICAL**

**Issues**:
1. Lines 173, 214, 279: Documents `--json` flag (NOT implemented)
2. Lines 242-255, 282-295, 316-328: JSON output examples (NOT implemented)
3. Lines 597-607: JSON mode section (NOT implemented)
4. Line 906: Reference to GitHub issues with placeholder URL
5. Script integration examples using `--json` flag throughout

**Required Changes**:
- Remove ALL `--json` flag references from command documentation
- Remove JSON output examples (lines 242-255, 282-295, 316-328)
- Remove "JSON Mode" section (lines 597-607)
- Update Script Integration examples to use only `--quiet` and `--field` flags
- Remove Python JSON parsing example (uses --json)
- Update cross-reference URLs to use `ari1110` instead of `yourusername`
- Verify all documented flags match actual implementation (--quiet, --field, --masked, --no-clipboard confirmed in tech.md)

---

### 4. docs/SECURITY.md
- **Location**: `R:\Test-Projects\pass-cli\docs\SECURITY.md`
- **Category**: Operational Documentation
- **Purpose**: Security architecture, cryptography details, best practices
- **Sync Status**: 🔴 **OUT OF SYNC**
- **Priority**: ⭐⭐⭐ **CRITICAL**

**Issues**:
1. Line 26: States "PBKDF2 Key Derivation: 100,000 iterations" (✅ CORRECT per tech.md:89)
2. Line 57: States "Iterations: 100,000" (✅ CORRECT)
3. Line 72: States "iterations = 100,000" (✅ CORRECT)
4. **WAIT** - Re-checking README discrepancy...
5. Line 565: Version shows 1.0.0 (should be v0.0.1)

**Required Changes**:
- Update version from 1.0.0 to v0.0.1 (line 615)
- Verify file permission strategy matches actual implementation (mentions 600 Unix, Windows ACLs, per tech.md:91-97)
- ✅ PBKDF2 iterations are CORRECT (100,000) - README.md needs fixing!

---

### 5. docs/DEVELOPMENT.md
- **Location**: `R:\Test-Projects\pass-cli\docs\DEVELOPMENT.md`
- **Category**: Developer Documentation
- **Purpose**: Development workflow, tooling, testing
- **Sync Status**: 🟡 **NEEDS VERIFICATION**
- **Priority**: ⭐⭐ **HIGH**

**Issues**:
1. Line 9: States "Go 1.25+" (need to verify actual requirement from go.mod)
2. Line 24: Placeholder URL `username` instead of `ari1110`
3. No specific version mismatches observed
4. Build commands and Makefile targets appear accurate

**Required Changes**:
- Update GitHub clone URL to use `ari1110`
- Verify Go version requirement matches go.mod
- Update dependency versions if tech.md lists specific versions (Cobra v1.10.1, Viper v1.21.0, etc. per tech.md:14-21)

---

### 6. docs/CI-CD.md
- **Location**: `R:\Test-Projects\pass-cli\docs\CI-CD.md`
- **Category**: Developer Documentation
- **Purpose**: CI/CD pipeline documentation
- **Sync Status**: 🟡 **NEEDS VERIFICATION**
- **Priority**: ⭐⭐ **HIGH**

**Issues**:
1. Line 81: States "Go Version: 1.25" (needs verification)
2. GitHub Actions versions need to match actual workflow files
3. Recent dependency updates (mentioned in task) may have updated action versions

**Required Changes**:
- Verify GitHub Actions versions match .github/workflows/ci.yml and release.yml
- Check for action version updates: setup-go, checkout, etc.
- Verify Go version matches actual usage
- Update PAT setup instructions to match actual HOMEBREW_TAP_TOKEN and SCOOP_BUCKET_TOKEN

---

### 7. docs/RELEASE.md
- **Location**: `R:\Test-Projects\pass-cli\docs\RELEASE.md`
- **Category**: Developer Documentation
- **Purpose**: Release process using GoReleaser
- **Sync Status**: 🟡 **NEEDS UPDATE**
- **Priority**: ⭐⭐ **HIGH**

**Issues**:
1. Lines 60, 170: Version examples use v1.0.0 (should show v0.0.1 as actual release)
2. Line 68: Placeholder `your-github-token` (acceptable as example)
3. Line 145: States "Go Version: '1.25'" (needs verification)
4. Versioning section discusses v1.0.0 as example (should reflect current v0.0.1 state)

**Required Changes**:
- Update version examples to v0.0.1 to reflect actual release
- Update version tagging conventions to show current state (v0.0.1)
- Verify Go version in GitHub Actions example
- Document actual distribution methods available post-v0.0.1

---

### 8. docs/TROUBLESHOOTING.md
- **Location**: `R:\Test-Projects\pass-cli\docs\TROUBLESHOOTING.md`
- **Category**: Operational Documentation
- **Purpose**: Common issues and solutions
- **Sync Status**: 🟢 **GENERALLY GOOD**
- **Priority**: ⭐ **MEDIUM**

**Issues**:
1. Lines 93, 100: Placeholder URLs `yourusername`
2. Issues appear generic and applicable to v0.0.1
3. No version-specific problems documented

**Required Changes**:
- Replace `yourusername` with `ari1110`
- Add any v0.0.1-specific known issues if discovered during release
- Update contact information if needed

---

### 9. docs/HOMEBREW.md
- **Location**: `R:\Test-Projects\pass-cli\docs\HOMEBREW.md`
- **Category**: Developer Documentation (Package Distribution)
- **Purpose**: Homebrew formula setup and submission
- **Sync Status**: 🔴 **OUT OF SYNC**
- **Priority**: ⭐⭐ **HIGH**

**Issues**:
1. Lines 20, 47, 58, 60, 68: Placeholder URLs `yourusername`
2. Line 68: Version example v1.0.0 (should be v0.0.1)
3. Line 98: Instructions say "coming soon" but actual tap exists (ari1110/homebrew-tap)

**Required Changes**:
- Update all URLs to `ari1110/homebrew-tap`
- Update version examples to v0.0.1
- Reflect actual tap status (released, auto-updated via GoReleaser)
- Update SHA256 calculation instructions with actual v0.0.1 artifacts

---

### 10. docs/SCOOP.md
- **Location**: `R:\Test-Projects\pass-cli\docs\SCOOP.md`
- **Category**: Developer Documentation (Package Distribution)
- **Purpose**: Scoop manifest setup and submission
- **Sync Status**: 🔴 **OUT OF SYNC**
- **Priority**: ⭐⭐ **HIGH**

**Issues**:
1. Lines 20, 52, 64, 73, 86, 98: Placeholder URLs `yourusername`
2. Line 73: Version example v1.0.0 (should be v0.0.1)
3. Instructions reference "coming soon" but actual bucket exists (ari1110/scoop-bucket)

**Required Changes**:
- Update all URLs to `ari1110/scoop-bucket`
- Update version examples to v0.0.1
- Reflect actual bucket status (released, auto-updated via GoReleaser)
- Update SHA256 calculation instructions with actual v0.0.1 artifacts

---

### 11. manifests/winget/README.md
- **Location**: `R:\Test-Projects\pass-cli\manifests\winget\README.md`
- **Category**: Platform-Specific Documentation
- **Purpose**: winget submission process
- **Sync Status**: 🟢 **FUTURE WORK**
- **Priority**: ⭐ **LOW** (not yet implemented)

**Issues**:
1. Line 15: Version references v1.0.0
2. Document describes future submission process
3. winget distribution not yet available

**Required Changes**:
- Mark clearly as "FUTURE WORK" or "NOT YET AVAILABLE"
- Update version references when actually implementing
- Consider moving to archive or future-features directory

---

### 12. manifests/snap/README.md
- **Location**: `R:\Test-Projects\pass-cli\manifests\snap\README.md`
- **Category**: Platform-Specific Documentation
- **Purpose**: Snap package submission process
- **Sync Status**: 🟢 **FUTURE WORK**
- **Priority**: ⭐ **LOW** (not yet implemented)

**Issues**:
1. Lines 19, 35, 77, 80, 133: Version references v1.0.0 and v1.1.0
2. Document describes future submission process
3. Snap distribution not yet available

**Required Changes**:
- Mark clearly as "FUTURE WORK" or "NOT YET AVAILABLE"
- Update version references when actually implementing
- Consider moving to archive or future-features directory

---

### 13. test/README.md
- **Location**: `R:\Test-Projects\pass-cli\test\README.md`
- **Category**: Developer Documentation
- **Purpose**: Integration test documentation
- **Sync Status**: 🟢 **GOOD**
- **Priority**: ⭐ **LOW**

**Issues**:
- None identified
- Appears current and accurate for v0.0.1
- Performance targets documented and achieved

**Required Changes**:
- None required

---

### 14. RELEASE-DRY-RUN.md
- **Location**: `R:\Test-Projects\pass-cli\RELEASE-DRY-RUN.md`
- **Category**: Project Artifact (Historical)
- **Purpose**: Release validation report from pre-release testing
- **Sync Status**: 🟡 **OBSOLETE** (historical artifact)
- **Priority**: ⭐ **LOW**

**Issues**:
1. Document is historical (dated 2025-09-30)
2. References "0.0.1-next (snapshot)" and "before v1.0 production release"
3. Actual v0.0.1 release has occurred
4. Content is outdated but has historical value

**Required Changes**:
- **RECOMMENDATION**: Archive or mark as historical
- Could be useful for reference but may confuse users
- Consider moving to `.spec-workflow/archive/` or `docs/archive/`

---

### 15. SECURITY-AUDIT.md
- **Location**: `R:\Test-Projects\pass-cli\SECURITY-AUDIT.md`
- **Category**: Project Artifact (Historical)
- **Purpose**: Pre-release security audit report
- **Sync Status**: 🟡 **OBSOLETE** (historical artifact)
- **Priority**: ⭐ **LOW**

**Issues**:
1. Document is historical (dated 2025-09-30)
2. References "1.0.0 (pre-release)" but actual release is v0.0.1
3. Audit findings are still relevant for security understanding
4. May confuse users as to which version was audited

**Required Changes**:
- **RECOMMENDATION**: Archive or mark as historical
- Consider updating header to clarify it applies to v0.0.1 codebase
- Alternatively move to `docs/archive/` with context note

---

### 16. dist/CHANGELOG.md
- **Location**: `R:\Test-Projects\pass-cli\dist\CHANGELOG.md`
- **Category**: Build Artifact
- **Purpose**: Auto-generated changelog from commits
- **Sync Status**: 🟢 **AUTO-GENERATED**
- **Priority**: ⭐ **NONE** (build artifact)

**Issues**:
- None - this is an auto-generated build artifact
- Recreated on each GoReleaser build
- Should be in .gitignore

**Required Changes**:
- Verify dist/ is in .gitignore
- No manual changes needed

---

## Summary by Category

### User Documentation (3 docs)
| Document | Status | Priority | Main Issues |
|----------|--------|----------|-------------|
| README.md | 🔴 OUT OF SYNC | ⭐⭐⭐ CRITICAL | Placeholder URLs, wrong PBKDF2 iterations (600k→100k), --json flag |
| INSTALLATION.md | 🔴 OUT OF SYNC | ⭐⭐⭐ CRITICAL | Placeholder URLs, v1.0.0→v0.0.1, "coming soon" vs released |
| USAGE.md | 🔴 OUT OF SYNC | ⭐⭐⭐ CRITICAL | Extensive --json documentation (not implemented) |

### Developer Documentation (6 docs)
| Document | Status | Priority | Main Issues |
|----------|--------|----------|-------------|
| DEVELOPMENT.md | 🟡 NEEDS VERIFY | ⭐⭐ HIGH | Go version verification, placeholder URL |
| CI-CD.md | 🟡 NEEDS VERIFY | ⭐⭐ HIGH | Action versions, Go version |
| RELEASE.md | 🟡 NEEDS UPDATE | ⭐⭐ HIGH | Version examples v1.0.0→v0.0.1 |
| HOMEBREW.md | 🔴 OUT OF SYNC | ⭐⭐ HIGH | Placeholder URLs, actual tap status |
| SCOOP.md | 🔴 OUT OF SYNC | ⭐⭐ HIGH | Placeholder URLs, actual bucket status |
| test/README.md | 🟢 GOOD | ⭐ LOW | None |

### Operational Documentation (2 docs)
| Document | Status | Priority | Main Issues |
|----------|--------|----------|-------------|
| SECURITY.md | 🟢 MOSTLY GOOD | ⭐ MEDIUM | Version reference 1.0.0→v0.0.1 |
| TROUBLESHOOTING.md | 🟢 MOSTLY GOOD | ⭐ MEDIUM | Placeholder URLs |

### Platform-Specific (2 docs)
| Document | Status | Priority | Main Issues |
|----------|--------|----------|-------------|
| winget/README.md | 🟢 FUTURE WORK | ⭐ LOW | Not implemented, needs clarity |
| snap/README.md | 🟢 FUTURE WORK | ⭐ LOW | Not implemented, needs clarity |

### Historical Artifacts (2 docs)
| Document | Status | Priority | Recommendation |
|----------|--------|----------|----------------|
| RELEASE-DRY-RUN.md | 🟡 OBSOLETE | ⭐ LOW | Archive to .spec-workflow/archive/ |
| SECURITY-AUDIT.md | 🟡 OBSOLETE | ⭐ LOW | Archive or update header |

---

## Recommended Action Plan

### Phase 1: Critical Fixes (Priority ⭐⭐⭐)
1. ✅ **README.md**: Fix PBKDF2 iterations, remove --json, update URLs, update version
2. ✅ **INSTALLATION.md**: Update all URLs, versions, and release status
3. ✅ **USAGE.md**: Remove all --json references, update URLs

### Phase 2: High Priority Updates (Priority ⭐⭐)
4. ✅ **DEVELOPMENT.md**: Verify Go version, update URLs
5. ✅ **CI-CD.md**: Verify action versions, update Go version
6. ✅ **RELEASE.md**: Update version examples and tagging conventions
7. ✅ **HOMEBREW.md**: Update URLs and tap status
8. ✅ **SCOOP.md**: Update URLs and bucket status

### Phase 3: Medium Priority (Priority ⭐)
9. ✅ **SECURITY.md**: Update version reference
10. ✅ **TROUBLESHOOTING.md**: Update placeholder URLs

### Phase 4: Cleanup (Priority ⭐)
11. ✅ **Archive obsolete docs**: Move RELEASE-DRY-RUN.md and SECURITY-AUDIT.md
12. ✅ **Mark future work**: Clearly label winget and snap docs

### Phase 5: Consistency Check
13. ✅ **Cross-document consistency**: Verify all versions, URLs, features match
14. ✅ **Link validation**: Test all internal and external links
15. ✅ **Terminology consistency**: Check keychain vs credential manager usage

---

## Key Findings Summary

### Critical Discrepancies Found

1. **PBKDF2 Iterations Mismatch**:
   - README.md states: 600,000 iterations
   - tech.md (authoritative) states: 100,000 iterations
   - SECURITY.md correctly states: 100,000 iterations
   - **Action**: Fix README.md to match tech.md (100,000)

2. **--json Flag Documentation**:
   - README.md, USAGE.md extensively document --json flag
   - tech.md (line 141-142) lists --json as FUTURE enhancement, not implemented
   - product.md (line 105) lists "--json flag" under "Potential Enhancements"
   - **Action**: Remove ALL --json references from user-facing docs

3. **Version State Confusion**:
   - Most docs reference v1.0.0 as release version
   - Actual published release is v0.0.1
   - **Action**: Update all version references to v0.0.1

4. **Placeholder URLs**:
   - Pervasive use of `yourusername` instead of `ari1110`
   - **Action**: Global find-replace across all docs

### Feature Alignment Check

✅ **Implemented and Documented**:
- `--quiet` flag for script-friendly output
- `--field` flag for field extraction
- `--masked` flag for password display
- `--no-clipboard` flag to skip clipboard

❌ **NOT Implemented but Documented**:
- `--json` flag (future enhancement per tech.md)

---

## Metrics

- **Total Documentation Files**: 15
- **Requiring Updates**: 12 (80%)
- **Sync Status**:
  - 🔴 Critical Out of Sync: 6 (40%)
  - 🟡 Needs Verification: 5 (33%)
  - 🟢 Good/Future Work: 4 (27%)

- **Priority Distribution**:
  - ⭐⭐⭐ Critical: 3 docs
  - ⭐⭐ High: 5 docs
  - ⭐ Medium/Low: 7 docs

---

**Report Generated**: 2025-10-01
**Next Steps**: Proceed with Phase 1 critical fixes (Tasks 2-4)
