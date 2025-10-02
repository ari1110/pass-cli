# Tasks Document

- [ ] 1. Create docs/development/ directory and README
  - File: docs/development/README.md (NEW)
  - Create new directory for implementation tracking documents
  - Write README explaining purpose and scope
  - Purpose: Establish location for development-focused documentation
  - _Leverage: Existing docs/archive/README.md as template_
  - _Requirements: 2.1-2.5 (Documentation Structure Alignment)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical writer and repository organizer | Task: Create docs/development/ directory and write docs/development/README.md following requirements 2.1-2.5. The README should explain: "This directory contains implementation tracking documents and development notes. These are NOT user-facing documentation. Contents: Dashboard implementation summaries, testing checklists, keybinding audits, and other development artifacts." Include a brief list of files that will be in this directory (dashboard summaries, testing checklists, keybinding audit). Use docs/archive/README.md as a reference for tone and structure | Restrictions: Keep README concise (< 100 lines), use clear markdown formatting, explain WHY files are here not just WHAT they are | Success: Directory created, README clearly explains purpose and distinguishes from user-facing docs, provides context for what belongs here | Instructions: After completing this task, update tasks.md to mark task 1 as complete [x] and mark task 2 as in-progress [-]_

- [ ] 2. Move implementation docs to docs/development/
  - Files: DASHBOARD_IMPLEMENTATION_SUMMARY.md, DASHBOARD_TESTING_CHECKLIST.md, KEYBINDINGS_AUDIT.md
  - Move from root to docs/development/ using git mv
  - Purpose: Clean root directory and organize implementation artifacts
  - _Leverage: git mv command to preserve file history_
  - _Requirements: 1.1-1.5 (Root Directory Cleanup)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Git operations specialist | Task: Move implementation documentation files from root to docs/development/ following requirements 1.1-1.5 using history-preserving git operations. Execute these commands: `git mv DASHBOARD_IMPLEMENTATION_SUMMARY.md docs/development/DASHBOARD_IMPLEMENTATION_SUMMARY.md`, `git mv DASHBOARD_TESTING_CHECKLIST.md docs/development/DASHBOARD_TESTING_CHECKLIST.md`, `git mv KEYBINDINGS_AUDIT.md docs/development/KEYBINDINGS_AUDIT.md`. Verify all files moved successfully and history preserved with `git log --follow docs/development/DASHBOARD_IMPLEMENTATION_SUMMARY.md` | Restrictions: MUST use `git mv` not manual move+add (preserves history), verify each move succeeded before next, do not modify file contents | Success: All three files moved to docs/development/, git history preserved (git log --follow shows original commits), root directory cleaner | Instructions: After completing this task, update tasks.md to mark task 2 as complete [x] and mark task 3 as in-progress [-]_

- [ ] 3. Move test utility to test/ directory
  - File: test-tui.bat
  - Move from root to test/ using git mv
  - Purpose: Organize test utilities with other test files
  - _Leverage: git mv command_
  - _Requirements: 1.1-1.5 (Root Directory Cleanup), 5.1-5.2 (Test Organization Review)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Test infrastructure organizer | Task: Move test-tui.bat from root to test/ directory following requirements 1.1-1.5 and 5.1-5.2. Execute: `git mv test-tui.bat test/test-tui.bat`. Verify the script still works from new location by checking if it references correct paths (should use relative paths or be run from project root) | Restrictions: Use git mv to preserve history, test that script still functions after move, do not modify script contents unless paths are broken | Success: test-tui.bat moved to test/ directory, git history preserved, script still functional from new location | Instructions: After completing this task, update tasks.md to mark task 3 as complete [x] and mark task 4 as in-progress [-]_

- [ ] 4. Update test/README.md with comprehensive test documentation
  - File: test/README.md
  - Enhance existing README to document test types, utilities, and test data
  - Purpose: Provide clear guidance on test organization and execution
  - _Leverage: Existing test/README.md, structure.md testing section_
  - _Requirements: 5.1-5.5 (Test Organization Review)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Test documentation specialist | Task: Update test/README.md following requirements 5.1-5.5 to comprehensively document test organization. Add or update sections: 1) Test Types: "Unit tests (*_test.go files adjacent to source code), Integration tests (test/*.go files)", 2) Running Tests: "Run all: `go test ./...`, Run integration only: `go test ./test`, Run unit only: `go test ./cmd/... ./internal/...`", 3) Test Utilities: "test-tui.bat - Manual TUI testing script (run from project root)", 4) Test Data: "test-vault/ - Integration test fixture directory with encrypted vault for testing". Read existing test/README.md first and enhance it, don't completely rewrite | Restrictions: Keep README under 150 lines, use clear examples, maintain existing content that's still relevant, ensure instructions are accurate | Success: test/README.md clearly documents all test types, provides runnable commands, explains test utilities and test data, developers can understand test structure from README alone | Instructions: After completing this task, update tasks.md to mark task 4 as complete [x] and mark task 5 as in-progress [-]_

- [ ] 5. Create or update docs/README.md as documentation index
  - File: docs/README.md
  - Create comprehensive index of all documentation with descriptions
  - Purpose: Help users and developers find relevant documentation quickly
  - _Leverage: Existing docs/ files, documentation structure from design.md_
  - _Requirements: 2.1-2.5 (Documentation Structure Alignment)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Documentation architect | Task: Create or enhance docs/README.md following requirements 2.1-2.5 as a comprehensive documentation index. Structure: 1) Introduction: "Pass-CLI Documentation Index", 2) User Documentation section listing: INSTALLATION.md, USAGE.md, SECURITY.md, TROUBLESHOOTING.md with brief descriptions, 3) Contributor Documentation section: DEVELOPMENT.md, 4) Maintainer Documentation section: RELEASE.md, CI-CD.md, HOMEBREW.md, SCOOP.md, 5) Development Notes section: "See development/ for implementation tracking docs", 6) Archive section: "See archive/ for historical documentation". Check if docs/README.md exists first with Read tool | Restrictions: Keep descriptions concise (one line per doc), organize by audience (user/contributor/maintainer), link to actual files, mention subdirectories (development/, archive/) | Success: docs/README.md serves as complete index, all documentation categorized by audience, easy to navigate, new contributors can find what they need within 30 seconds | Instructions: After completing this task, update tasks.md to mark task 5 as complete [x] and mark task 6 as in-progress [-]_

- [ ] 6. Update structure.md comprehensively
  - File: .spec-workflow/steering/structure.md
  - Add missing directories and update to reflect actual structure
  - Purpose: Make structure.md the single source of truth for repository navigation
  - _Leverage: Current structure.md, design.md directory purpose matrix_
  - _Requirements: 4.1-4.5 (Structure.md Synchronization)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Repository architect and technical writer | Task: Update .spec-workflow/steering/structure.md following requirements 4.1-4.5 to comprehensively document repository structure. Changes: 1) Add to Directory Organization tree: `docs/development/` (implementation tracking), `homebrew/` (Homebrew formula), `scoop/` (Scoop manifest), `manifests/` (platform-agnostic package managers), `test-vault/` (test fixtures), 2) Add new section after "Directory Organization" called "Package Manager Organization" explaining Pattern C: "Platform-Native (Root): homebrew/, scoop/ - tools with established root conventions. Platform-Agnostic (manifests/): snap/, winget/ - cross-platform manifest systems", 3) Update documentation section to mention docs/development/ and docs/archive/ with purposes, 4) Add file counts for major directories. Read current structure.md first to preserve existing content | Restrictions: Follow existing structure.md format and style, preserve all current content, only ADD missing information, ensure tree structure is accurate (verify with actual repo), keep explanations concise | Success: structure.md documents ALL directories in repo, package manager pattern clearly explained, documentation structure comprehensive, developers can navigate repo using structure.md as guide | Instructions: After completing this task, update tasks.md to mark task 6 as complete [x] and mark task 7 as in-progress [-]_

- [ ] 7. Update documentation references and links
  - Files: README.md, docs/*.md (any that reference moved files)
  - Check and update any links to files that were moved
  - Purpose: Prevent broken documentation links
  - _Leverage: grep to find references, Edit tool to update links_
  - _Requirements: All (ensure no broken references after moves)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Documentation maintainer and link validator | Task: Find and update all documentation references to moved files following all requirements. Steps: 1) Use Grep tool to search for references to moved files: `DASHBOARD_IMPLEMENTATION_SUMMARY.md`, `DASHBOARD_TESTING_CHECKLIST.md`, `KEYBINDINGS_AUDIT.md`, `test-tui.bat` in README.md and docs/*.md files, 2) For each reference found, update path to new location (e.g., DASHBOARD_IMPLEMENTATION_SUMMARY.md becomes docs/development/DASHBOARD_IMPLEMENTATION_SUMMARY.md), 3) Check README.md specifically for any links to docs/ files, 4) Verify no broken links remain. Use Grep with output_mode:"content" to see actual references | Restrictions: Only update paths that reference moved files, preserve link text and context, do not modify other content, verify each update is correct new path | Success: No references to old file paths remain, all links point to new locations, documentation remains coherent and navigable, no broken links | Instructions: After completing this task, update tasks.md to mark task 7 as complete [x] and mark task 8 as in-progress [-]_

- [ ] 8. Verify build and test suite execution
  - Files: Entire codebase
  - Run build and tests to ensure file moves didn't break anything
  - Purpose: Confirm reorganization doesn't affect functionality
  - _Leverage: go build, go test commands_
  - _Requirements: All (verification step)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Quality assurance engineer | Task: Verify build and test execution after file reorganization following all requirements as validation step. Execute using Bash tool: 1) `go build` - verify binary builds successfully, 2) `go test ./...` - verify all tests pass (unit and integration), 3) `go test ./test` - specifically verify integration tests pass, 4) Navigate to test/ and run `./test-tui.bat` (or `.\\test-tui.bat` on Windows) - verify test utility works from new location. Capture output of each command | Restrictions: Do not modify any code, only verify functionality, if any failures occur document them clearly for fixing, test execution from project root directory | Success: `go build` succeeds with no errors, all tests pass (`go test ./...` shows PASS), integration tests pass, test-tui.bat executes without path errors, file reorganization confirmed to have no functional impact | Instructions: After completing this task, update tasks.md to mark task 8 as complete [x] and mark task 9 as in-progress [-]_

- [ ] 9. Create git commit with all changes
  - Files: All modified and moved files
  - Create single atomic commit with reorganization changes
  - Purpose: Ensure all changes are bundled together for easy rollback if needed
  - _Leverage: git add, git commit commands_
  - _Requirements: All (final step to complete reorganization)_
  - _Prompt: Implement the task for spec repo-audit-cleanup, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Git operations specialist | Task: Create comprehensive git commit for all repository reorganization changes following all requirements. Steps: 1) Review changed files with `git status`, 2) Add all changes: `git add .` (files were already moved with git mv so they're staged, this catches any edits to READMEs and structure.md), 3) Create commit with message following this template using heredoc for multiline: "docs: reorganize repository structure for clarity\n\n- Move implementation docs to docs/development/\n- Move test utilities to test/\n- Create docs/development/README.md\n- Update test/README.md with comprehensive test docs\n- Create docs/README.md as documentation index\n- Update structure.md with missing directories\n- Document package manager organization pattern\n\nFixes root directory clutter and aligns structure with steering docs.\n\n🤖 Generated with [Claude Code](https://claude.com/claude-code)\n\nCo-Authored-By: Claude <noreply@anthropic.com>". Use git commit with heredoc pattern for multi-line message | Restrictions: Ensure all files are staged before commit (git status should show staged moves and modifications), use conventional commit format (docs: prefix), include descriptive body explaining changes, do NOT push (that's user's decision) | Success: Single commit created containing all reorganization changes, commit message is clear and descriptive, git log shows atomic changeset, repo ready to push or PR | Instructions: After completing this task, update tasks.md to mark task 9 as complete [x]_
