# Tasks

- [x] 1. Update Model struct view field types to tview
  - File: cmd/tui/model.go
  - Update Model struct to use tview view types instead of Bubble Tea view types
  - Purpose: Enable integration of existing tview view implementations
  - _Leverage: Existing tview view implementations (ListViewTview, DetailViewTview, AddFormTview, EditFormTview, HelpViewTview, ConfirmViewTview)_
  - _Requirements: 1.1_
  - _Prompt: Role: Go developer completing a framework migration by integrating existing tview view implementations | Task: Update the Model struct in cmd/tui/model.go to use tview view field types: Change listView *views.ListView to listView *views.ListViewTview, Change detailView *views.DetailView to detailView *views.DetailViewTview, Change addForm *views.AddFormView to addForm *views.AddFormTview, Change editForm *views.EditFormView to editForm *views.EditFormTview, Change helpView *views.HelpView to helpView *views.HelpViewTview, Change confirmView *views.ConfirmView to confirmView *views.ConfirmViewTview | Restrictions: Only change field TYPE declarations do not change field names, Do not modify any methods yet, Do not touch component fields those come in Phase 2, The file must compile after changes compilation errors will be fixed in subsequent tasks | Success: Model struct view fields are typed as tview variants, Field names remain unchanged, No other code modified_

- [x] 2. Update NewModel() view constructors
  - File: cmd/tui/model.go
  - Update NewModel() function to instantiate tview view implementations
  - Purpose: Instantiate tview views instead of Bubble Tea views
  - _Leverage: Existing tview view constructors (NewListViewTview, NewDetailViewTview, NewAddFormTview, NewEditFormTview, NewHelpViewTview, NewConfirmViewTview)_
  - _Requirements: 1.2_
  - _Prompt: Role: Go developer completing a framework migration by integrating existing tview view implementations | Task: Update the NewModel function in cmd/tui/model.go to instantiate tview view constructors: Change views.NewListView to views.NewListViewTview, Change views.NewDetailView to views.NewDetailViewTview, Change views.NewAddFormView to views.NewAddFormTview, Change views.NewEditFormView to views.NewEditFormTview, Change views.NewHelpView to views.NewHelpViewTview, Change views.NewConfirmView to views.NewConfirmViewTview | Restrictions: Only change constructor calls for views not components yet, Do not modify constructor parameters, Do not change any other initialization logic, The file must compile after changes | Success: NewModel calls tview view constructors, All view fields properly initialized, No changes to component constructors yet_

- [x] 3. Fix view method call compilation errors
  - File: cmd/tui/model.go, cmd/tui/events.go
  - Fix compilation errors from view type changes by updating method calls
  - Purpose: Ensure code compiles with tview view types
  - _Leverage: Existing tview view method signatures (GetTable, GetForm, GetFlex methods)_
  - _Requirements: 1.4_
  - _Prompt: Role: Go developer completing a framework migration by integrating existing tview view implementations | Task: Fix all compilation errors resulting from view type changes by running go build to identify compilation errors, For each error related to view method calls identify the tview equivalent method such as GetTable instead of View and update the call site to use the correct tview method, Common changes include replacing Bubble Tea Update msg calls with tview event handlers, replacing View calls with GetTable GetForm GetFlex methods, updating any type assertions from ListView to ListViewTview | Restrictions: Only fix compilation errors related to view types, Do not change component-related code yet, Do not modify tview view implementations themselves, Maintain existing business logic and behavior | Success: go build completes successfully with no view-related compilation errors, All view method calls use correct tview methods, No runtime behavior changes_

- [x] 4. Verify view integration builds successfully
  - File: Project-wide
  - Verify view integration phase completed successfully and application compiles
  - Purpose: Ensure Phase 1 completion before proceeding to Phase 2
  - _Leverage: Go compiler for verification_
  - _Requirements: 1.5_
  - _Prompt: Role: Go developer completing a framework migration verification | Task: Verify view integration phase completion by running go build and ensuring it succeeds, running go test ./cmd/tui/views/... to check if view tests compile they may fail that is expected, verifying that old Bubble Tea view structs are NOT referenced in Model struct, documenting any remaining issues for next phase | Restrictions: Do not make code changes in this task, Only run verification commands, Do not proceed to component integration yet | Success: go build succeeds without errors, Model struct uses only tview view types, No references to old ListView DetailView AddFormView EditFormView HelpView ConfirmView in Model struct_

- [x] 5. Commit view integration changes
  - File: All files modified in Tasks 1-4
  - Create git commit capturing view integration phase
  - Purpose: Enable easy rollback if needed
  - _Leverage: Git version control_
  - _Requirements: Non-functional (Atomic Changes)_
  - _Prompt: Role: Go developer completing a framework migration with proper version control | Task: Commit the view integration changes by running git status to review changed files, running git diff to review exact changes, staging all modified files with git add, creating commit with message feat: Integrate tview view implementations into Model struct - Update Model view field types to tview variants - Update NewModel to use tview view constructors - Fix view method calls for tview compatibility Phase 1 of tview-migration-remediation spec. Generated with Claude Code Co-Authored-By: Claude noreply@anthropic.com | Restrictions: Do not push changes yet, Only commit view integration changes, Ensure commit message accurately reflects changes | Success: Changes committed to local repository, Commit message clearly describes view integration, Git history shows atomic commit for Phase 1_
  - Note: Remediation work was required after initial implementation. Initial work improperly stubbed functionality with TODOs. Remediation fixes: (1) Implemented SetupFormCallbacks() and added SetSubmitCallback at all form creation points to enable form submission, (2) Restored UpdateCredentials call instead of recreating list view, (3) Verified all tests pass. Commits: fcdfc3f (remediation), plus earlier commits.

- [x] 6. Update Model struct component field types to tview
  - File: cmd/tui/model.go
  - Update Model struct to use tview component types
  - Purpose: Enable integration of existing tview component implementations
  - _Leverage: Existing tview component implementations (SidebarTview, StatusBarTview, MetadataPanelTview, CommandBarTview)_
  - _Requirements: 2.1_
  - _Prompt: Role: Go developer completing a framework migration by integrating existing tview component implementations | Task: Update the Model struct in cmd/tui/model.go to use tview component field types: Change sidebar *components.SidebarPanel to sidebar *components.SidebarTview, Change statusBar *components.StatusBar to statusBar *components.StatusBarTview, Change metadataPanel *components.MetadataPanel to metadataPanel *components.MetadataPanelTview, Change commandBar *components.CommandBar to commandBar *components.CommandBarTview, Note that Breadcrumb already uses tview no change needed | Restrictions: Only change component field TYPE declarations, Do not change field names, Do not modify any methods yet, The file must compile after changes errors will be fixed in next task | Success: Model struct component fields are typed as tview variants, Field names remain unchanged, Breadcrumb field unchanged already tview-based_

- [x] 7. Update NewModel() component constructors and fix compilation errors
  - File: cmd/tui/model.go, cmd/tui/helpers.go, cmd/tui/helpers_test.go
  - Update component constructors and fix compilation errors
  - Purpose: Instantiate tview components and ensure code compiles
  - _Leverage: Existing tview component constructors (NewSidebarTview, NewStatusBarTview, NewMetadataPanelTview, NewCommandBarTview) and methods (GetTreeView, Render, GetFlex, GetModal)_
  - _Requirements: 2.2, 2.4_
  - _Prompt: Role: Go developer completing a framework migration by integrating existing tview component implementations | Task: Update component constructors and fix compilation errors by updating NewModel to call tview component constructors: Change components.NewSidebarPanel to components.NewSidebarTview, Change components.NewStatusBar to components.NewStatusBarTview, Change components.NewMetadataPanel to components.NewMetadataPanelTview, Change components.NewCommandBar to components.NewCommandBarTview, then running go build to identify compilation errors and fixing component method calls using tview component methods, Common changes include using GetTreeView for sidebar, using Render for status bar, using GetFlex for metadata panel, using GetModal for command bar | Restrictions: Only change component constructors and method calls, Do not modify tview component implementations, Maintain existing business logic | Success: NewModel calls tview component constructors, go build succeeds without component-related errors, All component method calls use correct tview methods_
  - Note: renderDashboardView() is Bubble Tea/lipgloss only and stubbed with placeholder text. Bubble Tea Update() method calls to component methods (View, Render, SetError, GetCommand, etc.) were stubbed or commented as they won't be used in tview mode. Updated method calls: SetFocus→SetFocused (sidebar), SetCredentialCount→SetCredentialCountTview (status bar), SetCurrentView→SetCurrentViewTview (status bar), SetShortcuts→SetShortcutsTview (status bar), SetCredential→ShowMetadata (metadata panel), Render→GetTextView().GetText() (status bar). Removed all SetSize() calls as tview handles sizing internally.

- [x] 8. Verify component integration and commit
  - File: Project-wide then all modified files
  - Verify component integration phase completion and commit changes
  - Purpose: Ensure Phase 2 completion
  - _Leverage: Go compiler for verification, Git for version control_
  - _Requirements: 2.5_
  - _Prompt: Role: Go developer completing a framework migration verification and version control | Task: Verify component integration and commit by running go build and ensuring it succeeds, verifying Model struct uses only tview component types, running git status and git diff to review changes, committing changes with message feat: Integrate tview component implementations into Model struct - Update Model component field types to tview variants - Update NewModel to use tview component constructors - Fix component method calls for tview compatibility Phase 2 of tview-migration-remediation spec. Generated with Claude Code Co-Authored-By: Claude noreply@anthropic.com | Restrictions: Do not proceed to activation yet, Only commit component integration changes | Success: go build succeeds, Model struct uses only tview types no SidebarPanel StatusBar MetadataPanel CommandBar, Changes committed to local repository_

- [ ] 9. Activate tview by setting useBubbleTea to false
  - File: cmd/tui/tui.go
  - Activate tview code path by changing framework toggle
  - Purpose: Make tview the active implementation
  - _Leverage: Existing runTview function, All integrated tview views and components from Phases 1-2_
  - _Requirements: 3.1, 3.2, 3.3_
  - _Prompt: Role: Go developer activating the tview framework after successful integration | Task: Activate tview in cmd/tui/tui.go by locating the line useBubbleTea := true and changing it to useBubbleTea := false, verifying the runTview code path is now active | Restrictions: Only change the boolean value, Do not modify runTview function, Do not delete Bubble Tea code path yet deletion comes in Phase 4 | Success: useBubbleTea set to false, Application will execute runTview instead of tea.NewProgram, File compiles successfully_

- [ ] 10. Manual testing and runtime verification
  - File: N/A (runtime testing)
  - Verify application runs correctly on tview and basic functionality works
  - Purpose: Ensure tview implementation works at runtime
  - _Leverage: Visual regression checklist (docs/development/TVIEW_MIGRATION_CHECKLIST.md)_
  - _Requirements: 3.4, 3.5, 7 (Functional Preservation)_
  - _Prompt: Role: Go developer verifying a framework migration through manual testing | Task: Test the tview application manually by building the application with go build, initializing test vault with ./pass-cli.exe init --vault test-vault/vault.enc, launching TUI with ./pass-cli.exe --vault test-vault/vault.enc, verifying basic functionality: Application launches without crashes, Credential list displays, Navigation works with j/k keys, Press a to test add form displays, Press ESC to return to list, Press question mark to verify help view, documenting any runtime errors or visual issues, If critical issues found set useBubbleTea back to true and investigate | Restrictions: Do not modify code based on minor visual differences yet, Only verify application does not crash and basic views render, Full visual regression testing comes after deletion phase | Success: Application launches successfully, Main views render list detail forms help, No critical crashes or errors, Basic keyboard navigation functional_

- [ ] 11. Delete old Bubble Tea view structs
  - File: cmd/tui/views/list.go, cmd/tui/views/detail.go, cmd/tui/views/form_add.go, cmd/tui/views/form_edit.go, cmd/tui/views/help.go, cmd/tui/views/confirm.go
  - Remove obsolete Bubble Tea view implementations
  - Purpose: Eliminate duplicate code now that tview is active
  - _Leverage: Confirmation that tview is working from Task 10_
  - _Requirements: 4.1, 4.3, 4.4_
  - _Prompt: Role: Go developer removing obsolete code after a successful framework migration | Task: Delete old Bubble Tea view structs and their methods by in cmd/tui/views/list.go deleting type ListView struct and all its methods keeping type ListViewTview struct and all its methods, in cmd/tui/views/detail.go deleting type DetailView struct and all its methods keeping type DetailViewTview struct and all its methods, in cmd/tui/views/form_add.go deleting type AddFormView struct and all its methods keeping type AddFormTview struct and all its methods, in cmd/tui/views/form_edit.go deleting type EditFormView struct and all its methods keeping type EditFormTview struct and all its methods, in cmd/tui/views/help.go deleting type HelpView struct and all its methods keeping type HelpViewTview struct and all its methods, in cmd/tui/views/confirm.go deleting type ConfirmView struct and all its methods keeping type ConfirmViewTview struct and all its methods, deleting old constructors NewListView NewDetailView NewAddFormView NewEditFormView NewHelpView NewConfirmView, running go build to verify no references remain | Restrictions: Only delete old Bubble Tea view code, Do not delete tview implementations, Do not modify any tview code, File must compile after deletion | Success: Old view structs ListView DetailView AddFormView EditFormView HelpView ConfirmView deleted, Only tview suffix structs remain, go build succeeds, No compilation errors from missing references_

- [ ] 12. Delete old Bubble Tea component structs and commit
  - File: cmd/tui/components/sidebar.go, cmd/tui/components/statusbar.go, cmd/tui/components/metadata_panel.go, cmd/tui/components/command_bar.go
  - Remove obsolete Bubble Tea component implementations and commit deletion phase
  - Purpose: Complete code cleanup and commit Phase 4
  - _Leverage: Confirmation that tview components work from Task 10_
  - _Requirements: 4.2, 4.3, 4.4, 4.5_
  - _Prompt: Role: Go developer removing obsolete code after a successful framework migration | Task: Delete old Bubble Tea component structs and commit by in cmd/tui/components/sidebar.go deleting type SidebarPanel struct and all its methods keeping type SidebarTview struct and all its methods, in cmd/tui/components/statusbar.go deleting type StatusBar struct and all its methods keeping type StatusBarTview struct and all its methods, in cmd/tui/components/metadata_panel.go deleting type MetadataPanel struct and all its methods keeping type MetadataPanelTview struct and all its methods, in cmd/tui/components/command_bar.go deleting type CommandBar struct and all its methods keeping type CommandBarTview struct and all its methods, deleting old constructors NewSidebarPanel NewStatusBar NewMetadataPanel NewCommandBar, running go build to verify compilation, committing with message refactor: Delete obsolete Bubble Tea view and component implementations - Remove old ListView DetailView AddFormView EditFormView HelpView ConfirmView - Remove old SidebarPanel StatusBar MetadataPanel CommandBar - Only tview implementations remain Phase 4 of tview-migration-remediation spec. Generated with Claude Code Co-Authored-By: Claude noreply@anthropic.com | Restrictions: Only delete old component code, Do not modify tview implementations, File must compile after deletion | Success: Old component structs deleted, Only tview suffix structs remain, go build succeeds, Changes committed_

- [ ] 13. Remove Bubble Tea imports and useBubbleTea code path
  - File: cmd/tui/tui.go, all view files, all component files
  - Remove all Bubble Tea imports and obsolete code path
  - Purpose: Clean up dependencies before running go mod tidy
  - _Leverage: Confirmation that tview is fully active_
  - _Requirements: 5.3, 5.4_
  - _Prompt: Role: Go developer completing dependency cleanup after framework migration | Task: Remove Bubble Tea imports and code paths by searching for Bubble Tea imports with grep or search for github.com/charmbracelet/bubbletea and github.com/charmbracelet/bubbles, removing all Bubble Tea imports from files, in cmd/tui/tui.go deleting the entire useBubbleTea conditional block and replacing with direct call to runTview model and removing tea import, removing any remaining Bubble Tea type references, running go build to verify no compilation errors | Restrictions: Do not remove lipgloss imports may be used with tview, Only remove bubbletea and bubbles packages, Do not modify tview code | Success: No imports of github.com/charmbracelet/bubbletea, No imports of github.com/charmbracelet/bubbles, useBubbleTea conditional removed from tui.go, go build succeeds without Bubble Tea dependencies_

- [ ] 14. Run go mod tidy and verify dependency removal
  - File: go.mod, go.sum
  - Clean up Go module dependencies to remove Bubble Tea packages
  - Purpose: Achieve final dependency cleanup goal
  - _Leverage: Removed imports from Task 13_
  - _Requirements: 5.1, 5.2, 5.3_
  - _Prompt: Role: Go developer completing dependency cleanup after framework migration | Task: Clean up Go module dependencies by running go mod tidy to remove unused dependencies, verifying Bubble Tea removed by checking go.mod does NOT contain github.com/charmbracelet/bubbletea, checking go.mod does NOT contain github.com/charmbracelet/bubbles, noting that Lipgloss MAY remain acceptable for tview styling, running go build to confirm application builds without Bubble Tea, committing changes with message chore: Remove Bubble Tea dependencies from go.mod - Remove bubbletea and bubbles imports - Run go mod tidy to clean dependencies - Application now exclusively uses tview Phase 5 of tview-migration-remediation spec. Generated with Claude Code Co-Authored-By: Claude noreply@anthropic.com | Restrictions: Do not manually edit go.mod, Only run go mod tidy, Do not remove tview or tcell dependencies | Success: go.mod does not list bubbletea or bubbles, go build succeeds, Application runs correctly, Changes committed_

- [ ] 15. Update tests and final validation
  - File: cmd/tui/views/*_test.go, cmd/tui/components/*_test.go, docs/development/TVIEW_MIGRATION_CHECKLIST.md
  - Update unit tests to verify tview implementations and perform final validation
  - Purpose: Complete migration with full test coverage and validation
  - _Leverage: Visual regression checklist, Existing test structure, Updated tview implementations_
  - _Requirements: 6 (Test Updates), 7 (Functional Preservation), Non-functional (Performance, Testability)_
  - _Prompt: Role: Go developer completing test updates and final validation after framework migration | Task: Update tests and perform final validation. First, update view tests by changing constructors from NewListView() to NewListViewTview() and updating type assertions from *ListView to *ListViewTview, applying these changes to all view test files. Second, update component tests (if not already done in Phase 6) by changing constructors to tview variants and updating type assertions. Third, run all tests with `go test ./...` and fix any test failures. Fourth, perform manual testing using TVIEW_MIGRATION_CHECKLIST.md: test on primary terminal emulator, verify all views render correctly, test keyboard shortcuts, test responsive breakpoints by resizing terminal, verify performance with startup < 200ms. Fifth, document any remaining issues or acceptable differences. Finally, create final commit with message: "test: Update tests for tview implementations and verify migration\n\n- Update view and component tests to use tview constructors\n- Verify all unit tests pass\n- Complete manual testing with visual regression checklist\n- Migration from Bubble Tea to tview complete\n\nPhase 5 of tview-migration-remediation spec complete.\n\n🤖 Generated with Claude Code\n\nCo-Authored-By: Claude <noreply@anthropic.com>" | Restrictions: All tests must pass. Manual testing must show no critical regressions. Document acceptable differences from Bubble Tea version. | Success: All unit tests pass (go test ./... succeeds), Application launches and runs correctly, Visual regression checklist shows acceptable parity, No Bubble Tea code remains in codebase, Migration complete and verified_
