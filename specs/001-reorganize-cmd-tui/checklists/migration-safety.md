# Migration Safety Checklist: Reorganize cmd/tui Directory Structure

**Purpose**: Validate requirements quality for safe, reversible directory reorganization with verification at each step
**Created**: 2025-10-09
**Feature**: [spec.md](../spec.md) | [plan.md](../plan.md) | [tasks.md](../tasks.md)
**Focus**: Migration safety requirements for rollback, verification checkpoints, and git history preservation
**Audience**: Peer reviewer validating requirements quality during PR review
**Scope**: Balanced coverage of UX requirements (TUI rendering) and technical migration mechanics

---

## Rollback & Recovery Requirements

- [ ] CHK001 - Are rollback procedures explicitly defined for each migration step? [Completeness, Gap]
- [ ] CHK002 - Is the rollback command syntax specified with concrete examples? [Clarity, Spec §Rollback]
- [ ] CHK003 - Are rollback requirements defined for partial failures during a step? [Coverage, Edge Case]
- [ ] CHK004 - Is the "clean working tree" validation requirement specified before rollback? [Completeness, Gap]
- [ ] CHK005 - Are recovery requirements defined for interrupted git operations? [Coverage, Exception Flow]
- [ ] CHK006 - Is the baseline branch requirement (pre-reorg-tui) explicitly documented? [Clarity, Spec §Assumptions]

## Verification Checkpoint Requirements

- [ ] CHK007 - Are verification criteria defined for EACH of the 4 user stories? [Completeness, Spec §US1-US4]
- [ ] CHK008 - Is "TUI renders correctly" quantified with specific observable criteria? [Measurability, Spec §SC-001]
- [ ] CHK009 - Are compilation verification requirements specified after each step? [Completeness, Spec §FR-008]
- [ ] CHK010 - Is "no black screen" defined with testable visual criteria? [Clarity, Spec §SC-001]
- [ ] CHK011 - Are manual testing procedures detailed enough for consistent execution? [Clarity, Plan §Verification]
- [ ] CHK012 - Is the verification sequence (compile → run → visual → interaction) explicitly ordered? [Completeness, Plan §Step 1-4]
- [ ] CHK013 - Are verification requirements consistent across all 4 migration steps? [Consistency, Spec §US1-US4]
- [ ] CHK014 - Is the test vault path and password documented for verification testing? [Completeness, Quickstart §Prerequisites]
- [ ] CHK015 - Are terminal capability requirements specified for TUI verification? [Gap, Spec §Assumptions]

## Git History Preservation Requirements

- [ ] CHK016 - Is the `git mv` command specified as the required directory move method? [Clarity, Spec §FR-004]
- [ ] CHK017 - Are git history verification commands explicitly documented? [Completeness, Tasks §T022]
- [ ] CHK018 - Is the verification criterion for "history preserved" measurable? [Measurability, Tasks §T022]
- [ ] CHK019 - Are git commit requirements defined after each user story completion? [Completeness, Spec §FR-006]
- [ ] CHK020 - Is the commit message format specified with examples? [Clarity, Plan §Decision 3]
- [ ] CHK021 - Are requirements for verifying git rename detection explicitly stated? [Completeness, Tasks §T021]

## Sequential Dependency Requirements

- [ ] CHK022 - Are the sequential dependencies between user stories explicitly documented? [Completeness, Tasks §Dependencies]
- [ ] CHK023 - Is the strict execution order (P1→P2→P3→P4) requirement clearly stated? [Clarity, Spec §User Stories]
- [ ] CHK024 - Are the consequences of out-of-order execution documented? [Coverage, Edge Case]
- [ ] CHK025 - Is the "MUST complete before next" requirement explicit for each phase? [Completeness, Tasks §Phase Dependencies]
- [ ] CHK026 - Are parallel execution opportunities clearly marked with [P] tags? [Clarity, Tasks §Format]
- [ ] CHK027 - Are the prerequisites for each user story explicitly listed? [Completeness, Spec §US1-US4]

## Package & Import Requirements

- [ ] CHK028 - Is the package declaration change (`main` → `tui`) specified for ALL files? [Completeness, Spec §FR-001]
- [ ] CHK029 - Is the function signature change (`main()` → `Run(vaultPath string) error`) precisely defined? [Clarity, Spec §FR-002]
- [ ] CHK030 - Are import path update requirements specified with exact find/replace patterns? [Clarity, Spec §FR-003]
- [ ] CHK031 - Is the verification requirement for "no missed occurrences" explicitly stated? [Completeness, Tasks §T017]
- [ ] CHK032 - Are requirements defined for handling vaultPath parameter (empty vs. provided)? [Completeness, Plan §Step 1]

## TUI Rendering & Visual Requirements

- [ ] CHK033 - Are all visual elements required for "TUI renders completely" enumerated? [Completeness, Spec §SC-001]
- [ ] CHK034 - Is "sidebar shows categories" verifiable with specific observable criteria? [Measurability, Tasks §T033]
- [ ] CHK035 - Are requirements for "credentials listed in table" defined with clear success criteria? [Clarity, Tasks §T033]
- [ ] CHK036 - Is "detail view shows credential details" specified with required fields? [Completeness, Tasks §T033]
- [ ] CHK037 - Are navigation testing requirements (arrow keys, Tab, Enter) explicitly documented? [Coverage, Tasks §T033]
- [ ] CHK038 - Are form interaction requirements (Ctrl+A) specified for verification? [Completeness, Tasks §T033]
- [ ] CHK039 - Is the password masking toggle requirement included in verification criteria? [Completeness, Tasks §T033]
- [ ] CHK040 - Are requirements for "no visual corruption" defined with testable criteria? [Measurability, Spec §SC-001]

## CLI Compatibility Requirements

- [ ] CHK041 - Are CLI command preservation requirements explicitly stated? [Completeness, Spec §FR-010]
- [ ] CHK042 - Is the argument parsing logic for TUI vs. CLI routing specified? [Clarity, Plan §Step 4]
- [ ] CHK043 - Are requirements defined for all CLI subcommands (list, get, add, etc.)? [Coverage, Tasks §T034]
- [ ] CHK044 - Are help flag requirements (--help, -h) explicitly documented? [Completeness, Tasks §T030]
- [ ] CHK045 - Are version flag requirements (--version, -v) specified? [Completeness, Tasks §T031]
- [ ] CHK046 - Is the edge case of `--vault` + `--help` combination addressed? [Edge Case, Tasks §T035]

## Error Handling & Exception Requirements

- [ ] CHK047 - Are error handling requirements defined for compilation failures at each step? [Coverage, Quickstart §Troubleshooting]
- [ ] CHK048 - Are requirements specified for detecting "black screen" issues during verification? [Completeness, Quickstart §Issue: Black screen]
- [ ] CHK049 - Are import error detection requirements explicitly stated? [Coverage, Quickstart §Issue: Import errors]
- [ ] CHK050 - Is the fallback behavior defined when CLI commands incorrectly launch TUI? [Edge Case, Quickstart §Issue: CLI commands launch TUI]
- [ ] CHK051 - Are requirements for handling partially completed steps documented? [Gap, Recovery Flow]

## Time & Performance Requirements

- [ ] CHK052 - Is the 2-hour completion time requirement measurable? [Measurability, Spec §SC-004]
- [ ] CHK053 - Is the <3 second TUI launch requirement explicitly stated? [Clarity, Spec §SC-006]
- [ ] CHK054 - Are performance degradation criteria defined if launch time exceeds threshold? [Gap, Non-Functional]

## Completeness & Coverage Validation

- [ ] CHK055 - Are requirements defined for all 10 functional requirements (FR-001 to FR-010)? [Traceability, Spec §Requirements]
- [ ] CHK056 - Are acceptance criteria defined for all 6 success criteria (SC-001 to SC-006)? [Traceability, Spec §Success Criteria]
- [ ] CHK057 - Do all 4 user stories have measurable acceptance scenarios? [Completeness, Spec §User Stories]
- [ ] CHK058 - Are edge cases documented in the spec addressed in verification tasks? [Coverage, Spec §Edge Cases]
- [ ] CHK059 - Is the out-of-scope boundary clearly defined? [Completeness, Spec §Out of Scope]

## Ambiguities & Conflicts

- [ ] CHK060 - Is "features function identically" quantified with specific test criteria? [Ambiguity, Spec §SC-002]
- [ ] CHK061 - Is "zero new compiler errors" verification method specified? [Clarity, Spec §SC-005]
- [ ] CHK062 - Are there any conflicting requirements between user stories? [Conflict, Spec §US1-US4]
- [ ] CHK063 - Is the distinction between "package tui" and "tui-tview" directory naming clear? [Clarity, Spec §FR-001, FR-004]

---

## Summary

**Total Items**: 63 checklist items
**Traceability**: 100% of items include spec/plan/tasks references or gap markers
**Focus Areas**:
- ✅ Migration Safety (rollback, git history, sequential execution)
- ✅ Verification Rigor (checkpoints, manual testing, visual criteria)
- ✅ Balanced Coverage (TUI rendering + technical migration mechanics)

**Depth Level**: Peer review gate - Requirements quality validation before implementation

**Critical Migration Safety Gaps to Address**:
- Partial failure recovery procedures
- Terminal capability requirements for TUI testing
- Concurrent interaction scenarios during migration

**Next Steps**: Address incomplete items before beginning implementation to ensure safe, reversible migration with comprehensive verification.
