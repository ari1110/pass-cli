# Tasks: Documentation Accuracy Verification and Remediation

**Input**: Design documents from `/specs/010-documentation-accuracy-verification/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, verification-procedures.md, audit-report.md

**Tests**: No test tasks - this feature IS the testing (documentation verification against implementation)

**Organization**: Tasks are grouped by user story (verification category) to enable independent execution of each category.

**Phase Numbering Note**: This tasks.md uses execution-based phase numbers (Phase 1-8), while plan.md uses workflow-based phase names (Phase 0: Research, Phase 1: Audit Design). Both refer to the same work - plan.md describes the conceptual workflow, tasks.md provides the execution sequence.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4, US5)
- Include exact file paths in descriptions

## Path Conventions
- Documentation files: `README.md`, `docs/USAGE.md`, `docs/MIGRATION.md`, `docs/SECURITY.md`, etc.
- Reference code: `cmd/*.go`, `internal/*/`
- Audit artifacts: `specs/010-documentation-accuracy-verification/audit-report.md`

---

## Phase 1: Setup (Test Environment)

**Purpose**: Prepare isolated test environment for verification workflow

- [X] T001 Build pass-cli binary from current main branch: `go build -o pass-cli.exe`
- [X] T002 Create test vault directory: `mkdir ~/.pass-cli-test`
- [X] T003 Initialize test vault: `export PASS_VAULT=~/.pass-cli-test/vault.enc && ./pass-cli init` → Vault created at default location
- [ ] T004 [P] Add test credentials to test vault: `testservice` (user: test@example.com, password: TestPass123!) → **DEFERRED** (not required for US1 CLI verification)
- [ ] T005 [P] Add second test credential: `github` (user: user@example.com) → **DEFERRED** (not required for US1 CLI verification)
- [ ] T006 Verify test environment ready: run `./pass-cli list` and confirm 2 credentials → **DEFERRED** (not required for US1 CLI verification)

---

## Phase 2: Foundational (Design Artifacts) ✅ COMPLETE

**Purpose**: Core design documents that guide all verification tasks

**⚠️ CRITICAL**: These artifacts MUST exist before verification can begin

- [X] T007 Create research.md with verification methodology decisions → **Already complete** (commit 50c6ab4)
- [X] T008 Create verification-procedures.md with detailed test procedures for all 10 categories → **Already complete** (commit 50c6ab4)
- [X] T009 Create audit-report.md template with discrepancy tracking structure → **Already complete** (commit 50c6ab4)
- [X] T010 Update agent context (CLAUDE.md) with documentation workflow → **Already complete** (commit 50c6ab4)

**Checkpoint**: Foundation ready - verification execution can now begin

---

## Phase 3: User Story 1 - CLI Interface Verification (Priority: P1) 🎯 MVP

**Goal**: Verify all documented CLI commands, flags, aliases match actual implementation in cmd/ directory

**Independent Test**: Run `pass-cli [command] --help` for all 14 commands, compare output against USAGE.md flag tables, identify discrepancies

### Verification for User Story 1

- [X] T011 [P] [US1] Execute `pass-cli init --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/init.txt`
- [X] T012 [P] [US1] Execute `pass-cli add --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/add.txt`
- [X] T013 [P] [US1] Execute `pass-cli get --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/get.txt`
- [X] T014 [P] [US1] Execute `pass-cli list --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/list.txt`
- [X] T015 [P] [US1] Execute `pass-cli update --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/update.txt`
- [X] T016 [P] [US1] Execute `pass-cli delete --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/delete.txt`
- [X] T017 [P] [US1] Execute `pass-cli generate --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/generate.txt`
- [X] T018 [P] [US1] Execute `pass-cli change-password --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/change-password.txt`
- [X] T019 [P] [US1] Execute `pass-cli version --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/version.txt`
- [X] T020 [P] [US1] Execute `pass-cli verify-audit --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/verify-audit.txt`
- [X] T021 [P] [US1] Execute `pass-cli config --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/config.txt`
- [X] T022 [P] [US1] Execute `pass-cli config init --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/config-init.txt`
- [X] T023 [P] [US1] Execute `pass-cli config edit --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/config-edit.txt`
- [X] T024 [P] [US1] Execute `pass-cli config validate --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/config-validate.txt`
- [X] T025 [P] [US1] Execute `pass-cli config reset --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/config-reset.txt`
- [X] T026 [P] [US1] Execute `pass-cli tui --help`, capture output to `specs/010-documentation-accuracy-verification/help-output/tui.txt`
- [X] T027 [US1] Compare `init` command help output against docs/USAGE.md:77-123 flag table, document discrepancies in audit-report.md → Verified accurate
- [X] T028 [US1] Compare `add` command help output against docs/USAGE.md:139-170 flag table, document discrepancies → **Confirmed DISC-002, DISC-003**
- [ ] T029 [US1] Compare `get` command help output against docs/USAGE.md flag table, document discrepancies → **DEFERRED** (critical fixes prioritized)
- [ ] T030 [US1] Compare `list` command help output against docs/USAGE.md flag table, document discrepancies → **DEFERRED**
- [ ] T031 [US1] Compare `update` command help output against docs/USAGE.md flag table, document discrepancies → **DEFERRED**
- [ ] T032 [US1] Compare `delete` command help output against docs/USAGE.md flag table, document discrepancies → **DEFERRED**
- [ ] T033 [US1] Compare `generate` command help output against docs/USAGE.md flag table, document discrepancies → **DEFERRED**
- [X] T034 [US1] Compare all command help outputs against README.md examples, document discrepancies → **Confirmed DISC-001**
- [ ] T035 [US1] Verify command aliases (generate/gen/pwd) exist in cmd/*.go files per docs/USAGE.md:15 → **DEFERRED**
- [ ] T036 [US1] Update audit-report.md summary statistics for Category 1 (CLI Interface) → **DEFERRED**

### Remediation for User Story 1

- [X] T037 [US1] Remediate README.md:158,161 - remove non-existent `--generate` flag from `add` examples → **Fixed DISC-001**
- [X] T038 [US1] Remediate docs/USAGE.md:145-147 - remove `--generate`, `--length`, `--no-symbols` from `add` flag table → **Fixed DISC-002**
- [X] T039 [US1] Remediate docs/USAGE.md:139-147 - add missing `--category`/`-c` flag to `add` command flag table → **Fixed DISC-003**
- [X] T040 [US1] Remediate docs/USAGE.md:165,168 - remove `--generate` from `add` code examples, add `--category` example
- [ ] T041 [US1] Commit CLI interface remediation: `git add README.md docs/USAGE.md specs/010-documentation-accuracy-verification/tasks.md && git commit -m "docs: fix CLI interface discrepancies for add command (DISC-001, DISC-002, DISC-003)"`

**Checkpoint**: User Story 1 complete - CLI interface documentation matches actual implementation

---

## Phase 4: User Story 2 - Code Examples Verification (Priority: P2)

**Goal**: Verify all bash and PowerShell code examples in documentation execute successfully

**Independent Test**: Extract all code blocks, execute in test vault, verify exit codes and output

### Verification for User Story 2

- [ ] T042 [P] [US2] Extract all bash code blocks from README.md to `specs/010-documentation-accuracy-verification/examples/readme-bash.sh`
- [ ] T043 [P] [US2] Extract all bash code blocks from docs/USAGE.md to `specs/010-documentation-accuracy-verification/examples/usage-bash.sh`
- [ ] T044 [P] [US2] Extract all bash code blocks from docs/MIGRATION.md to `specs/010-documentation-accuracy-verification/examples/migration-bash.sh`
- [ ] T045 [P] [US2] Extract PowerShell code blocks from docs/USAGE.md to `specs/010-documentation-accuracy-verification/examples/usage-powershell.ps1`
- [ ] T046 [US2] Execute README.md bash examples against test vault, document exit codes and discrepancies in audit-report.md
- [ ] T047 [US2] Execute USAGE.md bash examples against test vault, document exit codes and discrepancies
- [ ] T048 [US2] Execute MIGRATION.md bash examples against test vault, document discrepancies → **Known issue: DISC-004**
- [ ] T049 [US2] Execute USAGE.md PowerShell examples (Windows only), document discrepancies
- [ ] T050 [US2] Verify output examples match actual CLI output per docs/USAGE.md output samples
- [ ] T051 [US2] Update audit-report.md summary statistics for Category 2 (Code Examples)

### Remediation for User Story 2

- [X] T052 [P] [US2] Remediate docs/MIGRATION.md:141-142,193,259-260,379 - remove `--generate` from migration examples (4 occurrences) → **Fixed DISC-004**
- [X] T053 [P] [US2] Remediate docs/SECURITY.md:608 - update credential rotation recommendation to use `pass-cli generate` → **Fixed DISC-005**
- [ ] T054 [US2] Remediate README.md code examples identified in T046 → **DEFERRED** (critical fixes prioritized)
- [ ] T055 [US2] Remediate docs/USAGE.md code examples identified in T047 → **DEFERRED** (critical fixes prioritized)
- [ ] T056 [US2] Commit code examples remediation: `git add docs/MIGRATION.md docs/SECURITY.md specs/010-documentation-accuracy-verification/tasks.md && git commit -m "docs: fix code example discrepancies (DISC-004, DISC-005)"`

**Checkpoint**: User Story 2 complete - all code examples execute successfully

---

## Phase 5: User Story 3 - Configuration and File Paths Verification (Priority: P3)

**Goal**: Verify all file path references and YAML configuration examples match implementation

**Independent Test**: Check all documented paths against internal/config, validate YAML examples

### Verification for User Story 3

- [ ] T057 [P] [US3] Grep for config path references in all docs: `grep -r "~/.config/pass-cli\|%APPDATA%\|~/Library" docs/ README.md`
- [ ] T058 [P] [US3] Grep for vault path references: `grep -r "~/.pass-cli/vault.enc\|%USERPROFILE%\\.pass-cli" docs/ README.md`
- [ ] T059 [US3] Verify config paths in docs/USAGE.md against internal/config/config.go GetConfigPath() implementation
- [ ] T060 [US3] Verify vault paths in README.md and docs/USAGE.md against cmd/root.go GetVaultPath() implementation
- [ ] T061 [US3] Extract YAML config examples from README.md and docs/USAGE.md
- [ ] T062 [US3] Validate YAML examples against internal/config/config.go Config struct field names and types
- [ ] T063 [US3] Validate example values pass internal/config validation rules (min_width: 1-10000, min_height: 1-1000)
- [ ] T064 [US3] Document all path and config discrepancies in audit-report.md
- [ ] T065 [US3] Update audit-report.md summary statistics for Category 3 (File Paths) and Category 4 (Configuration)

### Remediation for User Story 3

- [ ] T066 [US3] Remediate any file path discrepancies found in T059-T060
- [ ] T067 [US3] Remediate any YAML config discrepancies found in T062-T063
- [ ] T068 [US3] Commit path and config remediation: `git add docs/ README.md && git commit -m "docs: fix file path and configuration discrepancies"`

**Checkpoint**: User Story 3 complete - paths and config examples accurate

---

## Phase 6: User Story 4 - Feature Claims and Architecture Verification (Priority: P4)

**Goal**: Verify documented features exist and architecture descriptions match internal/ packages

**Independent Test**: Manual testing of features, code inspection of architecture

### Verification for User Story 4

- [ ] T069 [P] [US4] Test audit logging: run `pass-cli init --enable-audit`, verify HMAC signatures in audit log per docs/SECURITY.md claims
- [ ] T070 [P] [US4] Inspect internal/audit code to verify HMAC-SHA256 usage matches docs/SECURITY.md description
- [ ] T071 [P] [US4] Test keychain integration: run `pass-cli init --use-keychain`, verify Windows Credential Manager/macOS Keychain entry
- [ ] T072 [P] [US4] Inspect internal/keychain code to verify platform-specific implementations match README.md claims
- [ ] T073 [P] [US4] Test password policy: attempt weak password on init, verify rejection matches docs/USAGE.md policy description
- [ ] T074 [P] [US4] Inspect internal/security package to verify policy enforcement (12+ chars, complexity) matches documentation
- [ ] T075 [P] [US4] Test TUI shortcuts: launch `pass-cli tui`, verify Ctrl+H, Ctrl+C shortcuts match README.md documentation
- [ ] T076 [US4] Verify architecture: run `ls internal/`, compare package structure against docs/SECURITY.md architecture descriptions
- [ ] T077 [US4] Verify internal/crypto package exists and contains AES-GCM, PBKDF2, HMAC per docs/SECURITY.md claims
- [ ] T078 [US4] Verify internal/vault package separation per library-first architecture claims
- [ ] T079 [US4] Document all feature and architecture discrepancies in audit-report.md
- [ ] T080 [US4] Update audit-report.md summary statistics for Category 5 (Feature Claims) and Category 6 (Architecture)

### Remediation for User Story 4

- [ ] T081 [US4] Remediate any feature claim discrepancies found in T069-T075
- [ ] T082 [US4] Remediate any architecture description discrepancies found in T076-T078
- [ ] T083 [US4] Commit feature and architecture remediation: `git add docs/SECURITY.md README.md && git commit -m "docs: fix feature and architecture discrepancies"`

**Checkpoint**: User Story 4 complete - features and architecture accurately documented

---

## Phase 7: User Story 5 - Metadata, Output Examples, Cross-References (Priority: P5)

**Goal**: Verify metadata current, output examples accurate, links valid

**Independent Test**: Check git tags/dates, test output format, validate markdown links

### Verification for User Story 5

- [ ] T084 [P] [US5] Verify version numbers: `git tag --list` and compare against README.md "Version: v0.0.1"
- [ ] T085 [P] [US5] Verify "Last Updated" dates: `git log --oneline -1 -- README.md docs/*.md` and compare against documented dates
- [ ] T086 [P] [US5] Execute `pass-cli list` in test vault, compare table format against docs/USAGE.md output example
- [ ] T087 [P] [US5] Execute `pass-cli add testservice2`, verify success message matches docs/USAGE.md example "✅ Credential added successfully!"
- [ ] T088 [US5] Extract all markdown links from README.md: `grep -o '\[.*\](.*)'`
- [ ] T089 [US5] Extract all markdown links from docs/*.md
- [ ] T090 [US5] Validate internal file references: check `docs/USAGE.md`, `docs/SECURITY.md` files exist
- [ ] T091 [US5] Validate internal anchor references: grep for heading existence in target files (e.g., `## Configuration` in docs/USAGE.md)
- [ ] T092 [US5] Document all metadata, output, and link discrepancies in audit-report.md
- [ ] T093 [US5] Update audit-report.md summary statistics for Category 7 (Metadata), Category 8 (Output Examples), Category 9 (Cross-References)

### Remediation for User Story 5

- [ ] T094 [US5] Remediate any metadata discrepancies (version numbers, dates) found in T084-T085
- [ ] T095 [US5] Remediate any output example discrepancies found in T086-T087
- [ ] T096 [US5] Remediate any broken link discrepancies found in T090-T091
- [ ] T097 [US5] Commit metadata and links remediation: `git add README.md docs/*.md && git commit -m "docs: fix metadata, output examples, and link discrepancies"`

**Checkpoint**: User Story 5 complete - metadata, output, and links accurate

---

## Phase 8: Polish & Validation

**Purpose**: Final success criteria verification and process documentation

- [ ] T098 [P] Verify SC-001: 100% CLI commands/flags match implementation (review audit-report.md Category 1 discrepancies all fixed)
- [ ] T099 [P] Verify SC-002: 100% code examples execute successfully (review audit-report.md Category 2 discrepancies all fixed)
- [ ] T100 [P] Verify SC-003: 100% file paths resolve (review audit-report.md Category 3 discrepancies all fixed)
- [ ] T101 [P] Verify SC-004: 100% YAML examples valid (review audit-report.md Category 4 discrepancies all fixed)
- [ ] T102 [P] Verify SC-005: 100% features verified (review audit-report.md Category 5 discrepancies all fixed)
- [ ] T103 [P] Verify SC-006: Architecture descriptions match (review audit-report.md Category 6 discrepancies all fixed)
- [ ] T104 [P] Verify SC-007: Metadata current (review audit-report.md Category 7 discrepancies all fixed)
- [ ] T105 [P] Verify SC-008: Output examples match (review audit-report.md Category 8 discrepancies all fixed)
- [ ] T106 [P] Verify SC-009: Links resolve (review audit-report.md Category 9 discrepancies all fixed)
- [ ] T107 Verify SC-010: Audit report complete with all discrepancies documented → audit-report.md
- [ ] T108 Verify SC-011: All discrepancies remediated with git commits (check git log for DISC-### references)
- [ ] T109 Verify SC-012: User trust restored - run through USAGE.md examples end-to-end, confirm zero "command not found" or "unknown flag" errors
- [ ] T110 Update audit-report.md final status: change "🚧 IN PROGRESS" to "✅ COMPLETE"
- [ ] T111 Update audit-report.md Final Validation Checklist: mark all SC-001 through SC-012 as complete
- [ ] T112 Document verification process in CONTRIBUTING.md: add new section "Documentation Verification" at end of file with workflow description and reference to `specs/010-documentation-accuracy-verification/verification-procedures.md` for detailed test procedures
- [ ] T113 Cleanup test environment: `rm -rf ~/.pass-cli-test`
- [ ] T114 Commit final validation: `git add specs/010-documentation-accuracy-verification/audit-report.md CONTRIBUTING.md && git commit -m "docs: complete documentation accuracy audit - all success criteria met"`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: ✅ Already complete (commit 50c6ab4)
- **User Stories (Phase 3-7)**: All depend on Setup + Foundational completion
  - User stories can proceed in parallel (if staffed) or sequentially in priority order (P1 → P2 → P3 → P4 → P5)
- **Polish (Phase 8)**: Depends on all user stories (Phase 3-7) being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Setup - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Setup - No dependencies on other stories (some remediation overlaps with US1 but independently testable)
- **User Story 3 (P3)**: Can start after Setup - No dependencies on other stories
- **User Story 4 (P4)**: Can start after Setup - No dependencies on other stories
- **User Story 5 (P5)**: Can start after Setup - No dependencies on other stories

### Within Each User Story

- Verification tasks (capture --help, extract examples, grep paths) can run in parallel where marked [P]
- Comparison/analysis tasks depend on verification completion
- Remediation tasks depend on discrepancy identification
- Commit tasks depend on all remediation for that story

### Parallel Opportunities

- **Setup (Phase 1)**: T004, T005 can run in parallel (different test credentials)
- **User Story 1 (Phase 3)**: T011-T026 (all --help captures) can run in parallel, T037-T040 (file remediations) can run in parallel
- **User Story 2 (Phase 4)**: T042-T045 (code block extraction) can run in parallel, T052-T053 (file remediations) can run in parallel
- **User Story 3 (Phase 5)**: T057-T058 (grep operations) can run in parallel
- **User Story 4 (Phase 6)**: T069-T075 (feature tests) can run in parallel
- **User Story 5 (Phase 7)**: T084-T087 (verification checks) can run in parallel
- **Polish (Phase 8)**: T098-T106 (success criteria checks) can run in parallel
- **All user stories (Phase 3-7)** can run in parallel if team capacity allows

---

## Parallel Example: User Story 1

```bash
# Launch all --help captures for User Story 1 together:
Task: "Execute pass-cli init --help, capture output"
Task: "Execute pass-cli add --help, capture output"
Task: "Execute pass-cli get --help, capture output"
... (14 commands total)

# Launch all remediation file edits together:
Task: "Remediate README.md:158,161 - remove --generate flag"
Task: "Remediate docs/USAGE.md:145-147 - remove --generate from table"
Task: "Remediate docs/USAGE.md:139-147 - add --category flag"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (test vault ready)
2. Phase 2: Foundational already complete ✅
3. Complete Phase 3: User Story 1 (CLI Interface Verification)
4. **STOP and VALIDATE**: Verify SC-001 met (100% CLI accuracy)
5. Review audit-report.md CLI findings
6. If ready, proceed to next priority

### Incremental Delivery

1. Complete Setup + Foundational → Environment ready
2. Add User Story 1 → CLI verified → Remediate (MVP - highest user impact!)
3. Add User Story 2 → Examples verified → Remediate
4. Add User Story 3 → Config/paths verified → Remediate
5. Add User Story 4 → Features verified → Remediate
6. Add User Story 5 → Metadata/links verified → Remediate
7. Each story independently improves documentation quality

### Parallel Team Strategy

With multiple maintainers:

1. Team completes Setup + Foundational together
2. Once Setup is done:
   - Maintainer A: User Story 1 (CLI)
   - Maintainer B: User Story 2 (Examples)
   - Maintainer C: User Story 3 (Config/Paths)
   - Maintainer D: User Story 4 (Features)
   - Maintainer E: User Story 5 (Metadata)
3. Stories complete and remediate independently
4. Converge for Phase 8 (final validation)

---

## Notes

- [P] tasks = different files/commands, can run in parallel
- [Story] label maps task to specific verification category for traceability
- Each user story represents a complete verification category that can be independently executed and remediated
- Known discrepancies (DISC-001 through DISC-005) are pre-identified from initial spot check
- Estimated 45-95 additional discrepancies to be discovered across all categories
- Commit after each user story remediation with DISC-### references in commit message
- Stop at any checkpoint to validate category independently
- Cleanup test environment only after all validation complete (Phase 8)
