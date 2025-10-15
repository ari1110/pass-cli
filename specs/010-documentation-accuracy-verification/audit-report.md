# Documentation Accuracy Audit Report

**Audit Date**: 2025-10-15
**Scope**: 7 primary documentation files + all docs/ subdirectory files
**Methodology**: Manual verification per [verification-procedures.md](./verification-procedures.md)
**Status**: 🚧 **IN PROGRESS** (Template - to be populated during implementation)

---

## Summary Statistics

**Total Discrepancies**: 9 (DISC-001 through DISC-009)

### By Category

| Category | Total | Critical | High | Medium | Low |
|----------|-------|----------|------|--------|-----|
| CLI Interface | 9 | 6 | 0 | 3 | 0 |
| Code Examples | - | - | - | - | - |
| File Paths | - | - | - | - | - |
| Configuration | - | - | - | - | - |
| Feature Claims | - | - | - | - | - |
| Architecture | - | - | - | - | - |
| Metadata | - | - | - | - | - |
| Output Examples | - | - | - | - | - |
| Cross-References | - | - | - | - | - |
| Behavioral Descriptions | - | - | - | - | - |

### By File

| File | Discrepancies | Status |
|------|---------------|--------|
| README.md | - | ❌ Not Started |
| docs/USAGE.md | 8 | 🔧 In Progress |
| docs/MIGRATION.md | - | ❌ Not Started |
| docs/SECURITY.md | - | ❌ Not Started |
| docs/TROUBLESHOOTING.md | - | ❌ Not Started |
| docs/KNOWN_LIMITATIONS.md | - | ❌ Not Started |
| CONTRIBUTING.md | - | ❌ Not Started |
| docs/INSTALLATION.md | - | ❌ Not Started |

---

## Discrepancy Details

### README.md

#### DISC-001 [CLI/Critical] Non-existent `--generate` flag documented for `add` command

- **Location**: README.md:158, 161
- **Category**: CLI Interface
- **Severity**: Critical
- **Documented**:
  ```bash
  pass-cli add newservice --generate
  pass-cli add newservice --generate --length 32 --no-symbols
  ```
- **Actual**: cmd/add.go does not define `--generate`, `--length`, or `--no-symbols` flags. These belong to `pass-cli generate` command.
- **Remediation**: Remove `--generate` flag references from `add` examples. Document `pass-cli generate` as separate command for password generation.
- **Status**: ❌ Open
- **Commit**: [TBD]

---

### docs/USAGE.md

#### DISC-002 [CLI/Critical] Non-existent `--generate` flag documented for `add` command

- **Location**: docs/USAGE.md:145-147, 165, 168
- **Category**: CLI Interface
- **Severity**: Critical
- **Documented**: Flag table shows `--generate`, `--length`, `--no-symbols` as flags for `add` command
- **Actual**: cmd/add.go does not define these flags
- **Remediation**: Remove from flag table, update examples to use separate `pass-cli generate` command
- **Status**: ❌ Open
- **Commit**: [TBD]

---

#### DISC-003 [CLI/Medium] Missing `--category` flag in USAGE.md flag table

- **Location**: docs/USAGE.md:139-147 (add command flags table)
- **Category**: CLI Interface
- **Severity**: Medium
- **Documented**: Flag table does not list `--category` / `-c` flag
- **Actual**: cmd/add.go:57 defines `addCmd.Flags().StringVarP(&addCategory, "category", "c", "", "category for organizing credentials")`
- **Remediation**: Add row to flag table:
  ```markdown
  | `--category` | `-c` | string | Category for organizing credentials |
  ```
- **Status**: ❌ Open
- **Commit**: [TBD]

---

### docs/MIGRATION.md

#### DISC-004 [CLI/Critical] Non-existent `--generate` flag in migration examples

- **Location**: docs/MIGRATION.md:141-142, 193, 259-260
- **Category**: CLI Interface / Code Examples
- **Severity**: Critical
- **Documented**: Migration examples show `pass-cli add service --generate`
- **Actual**: Flag does not exist
- **Remediation**: Update examples to use two-step process: `pass-cli generate` → copy password → `pass-cli add service -p [paste]` OR remove `--generate` and use interactive prompts
- **Status**: ❌ Open
- **Commit**: [TBD]

---

### docs/SECURITY.md

#### DISC-005 [CLI/Critical] Non-existent `--generate` flag in security best practices

- **Location**: docs/SECURITY.md:608
- **Category**: CLI Interface
- **Severity**: Critical
- **Documented**: Security best practices recommend `pass-cli update service --generate`
- **Actual**: cmd/update.go does not define `--generate` flag
- **Remediation**: Update recommendation to use separate `pass-cli generate` command, then `pass-cli update service --password [generated]`
- **Status**: ❌ Open
- **Commit**: [TBD]

---

#### DISC-006 [CLI/Critical] Non-existent `--copy`/`-c` flag documented for `get` command

- **Location**: docs/USAGE.md:217
- **Category**: CLI Interface
- **Severity**: Critical
- **Documented**: Flag table shows `--copy | -c | bool | Copy to clipboard only (no display)`
- **Actual**: cmd/get.go defines only `--quiet/-q`, `--field/-f`, `--masked`, `--no-clipboard` flags. NO `--copy` flag exists.
- **Remediation**: Remove `--copy` row from flag table, remove example at line 251 (`pass-cli get github --copy`)
- **Status**: ❌ Open
- **Commit**: [TBD]

---

#### DISC-007 [CLI/Critical] Non-existent `--generate` flag documented for `update` command

- **Location**: docs/USAGE.md:254-256, 274, 277
- **Category**: CLI Interface
- **Severity**: Critical
- **Documented**: Flag table shows `--generate`, `--length`, `--no-symbols` as flags for `update` command
- **Actual**: cmd/update.go does not define these flags. Only defines: `--username/-u`, `--password/-p`, `--category`, `--url`, `--notes`, `--clear-category`, `--clear-notes`, `--clear-url`, `--force`
- **Remediation**: Remove `--generate`, `--length`, `--no-symbols` from flag table and examples (lines 274, 277)
- **Status**: ❌ Open
- **Commit**: [TBD]

---

#### DISC-008 [CLI/Medium] Missing flags in `update` command documentation

- **Location**: docs/USAGE.md:246-253
- **Category**: CLI Interface
- **Severity**: Medium
- **Documented**: Flag table incomplete
- **Actual**: cmd/update.go defines `--category`, `--clear-category`, `--clear-notes`, `--clear-url`, `--force` flags not documented in table
- **Remediation**: Add missing flags to table:
  ```markdown
  | `--category` | | string | New category |
  | `--clear-category` | | bool | Clear category field to empty |
  | `--clear-notes` | | bool | Clear notes field to empty |
  | `--clear-url` | | bool | Clear URL field to empty |
  | `--force` | `-f` | bool | Skip confirmation prompt |
  ```
- **Status**: ❌ Open
- **Commit**: [TBD]

---

#### DISC-009 [CLI/Medium] Non-existent `--copy` flag documented for `generate` command

- **Location**: docs/USAGE.md:370, 392
- **Category**: CLI Interface
- **Severity**: Medium
- **Documented**: Flag table shows `--copy | bool | Copy to clipboard only (no display)` and example at line 392
- **Actual**: cmd/generate.go defines only `--length/-l`, `--no-clipboard`, `--no-digits`, `--no-lower`, `--no-symbols`, `--no-upper`. NO `--copy` flag exists.
- **Remediation**: Remove `--copy` from flag table (line 370) and example (line 392: `pass-cli generate --copy`)
- **Status**: ❌ Open
- **Commit**: [TBD]

---

## Known Issues (Pre-Audit Findings)

The following discrepancies were identified during initial USAGE.md spot check (conversation leading to this spec):

1. **README.md:158, 161** - `--generate` flag for `add` command (DISC-001)
2. **docs/USAGE.md:145-147** - `--generate` flag table entry (DISC-002)
3. **docs/USAGE.md:139-147** - Missing `--category` flag (DISC-003)
4. **docs/MIGRATION.md** - Multiple `--generate` examples (DISC-004)
5. **docs/SECURITY.md:608** - `--generate` recommendation (DISC-005)

**Estimated Total**: 50-100 additional discrepancies anticipated across all 10 categories and remaining files.

---

## Appendix: Verification Test Log

### Category 1: CLI Interface Verification

**Test Date**: [TBD]
**Methodology**: Execute `pass-cli [command] --help`, compare against USAGE.md flag tables

| Command | Documented Flags | Actual Flags (from --help) | Discrepancies | Status |
|---------|------------------|---------------------------|---------------|--------|
| init | --use-keychain, --enable-audit | [TBD] | [TBD] | ❌ Not Tested |
| add | --username/-u, --password/-p, --category/-c, --url, --notes, --generate, --length, --no-symbols | [TBD] | DISC-002, DISC-003 | ❌ Not Tested |
| get | [TBD] | [TBD] | [TBD] | ❌ Not Tested |
| list | [TBD] | [TBD] | [TBD] | ❌ Not Tested |
| update | [TBD] | [TBD] | [TBD] | ❌ Not Tested |
| delete | [TBD] | [TBD] | [TBD] | ❌ Not Tested |
| generate | [TBD] | [TBD] | [TBD] | ❌ Not Tested |
| config | [TBD] | [TBD] | [TBD] | ❌ Not Tested |
| verify-audit | [TBD] | [TBD] | [TBD] | ❌ Not Tested |
| tui | [TBD] | [TBD] | [TBD] | ❌ Not Tested |

**Total Discrepancies Found**: [TBD]

---

### Category 2: Code Examples Verification

**Test Date**: [TBD]
**Methodology**: Extract bash/PowerShell blocks, execute in test vault

| File | Line | Code Block | Exit Code | Output Match | Discrepancies | Status |
|------|------|-----------|-----------|--------------|---------------|--------|
| README.md | 158 | `pass-cli add newservice --generate` | [TBD] | [TBD] | DISC-001 | ❌ Not Tested |
| [TBD] | [TBD] | [TBD] | [TBD] | [TBD] | [TBD] | ❌ Not Tested |

**Total Discrepancies Found**: [TBD]

---

### Categories 3-10: [To Be Populated During Implementation]

---

## Remediation Progress Tracker

### Phase 1: Critical/High Priority Fixes (Immediate User Impact)

- [ ] DISC-001: README.md `--generate` flag (Critical)
- [ ] DISC-002: USAGE.md `--generate` flag table (Critical)
- [ ] DISC-004: MIGRATION.md `--generate` examples (Critical)
- [ ] DISC-005: SECURITY.md `--generate` recommendation (Critical)
- [ ] [Additional Critical/High findings TBD]

**Target**: Fix all Critical/High within first implementation phase

---

### Phase 2: Medium Priority Fixes (Incomplete Documentation)

- [ ] DISC-003: USAGE.md missing `--category` flag (Medium)
- [ ] [Additional Medium findings TBD]

**Target**: Batch fix by file, commit after each file complete

---

### Phase 3: Low Priority Fixes (Cosmetic/Metadata)

- [ ] [Low priority findings TBD - dates, links, formatting]

**Target**: Final cleanup phase

---

## Final Validation Checklist

**Success Criteria Verification** (per spec.md):

- [ ] **SC-001**: 100% of documented CLI commands, flags, and aliases match actual implementation
- [ ] **SC-002**: 100% of code examples execute successfully without errors
- [ ] **SC-003**: 100% of file path references resolve to actual locations
- [ ] **SC-004**: 100% of configuration YAML examples pass validation
- [ ] **SC-005**: 100% of documented features verified to exist and function as described
- [ ] **SC-006**: All architecture descriptions match actual internal/ package structure
- [ ] **SC-007**: All version numbers and dates current as of remediation completion date
- [ ] **SC-008**: 100% of command output examples match actual CLI output format
- [ ] **SC-009**: 100% of internal markdown links resolve correctly
- [ ] **SC-010**: Audit report documents all discrepancies with file path, line number, issue description, and remediation action ✅ (this document)
- [ ] **SC-011**: All identified discrepancies remediated with git commits documenting rationale
- [ ] **SC-012**: User trust restored - documentation can be followed without encountering "command not found" or "unknown flag" errors

---

**Report Status**: 🚧 **IN PROGRESS** - 5 discrepancies documented (from pre-audit findings), full audit pending
