# Tasks Document

- [x] 1. Add Bubble Tea dependencies to project
  - File: go.mod
  - Add Bubble Tea framework and companion libraries (bubbletea v0.25.0, bubbles v0.18.0, lipgloss v0.9.1) to enable TUI development
  - Purpose: Enable TUI development with industry-standard framework
  - _Leverage: Existing go.mod structure, Go module tooling_
  - _Requirements: REQ-1_
  - _Prompt: Role: DevOps Engineer with expertise in Go dependency management | Task: Add Bubble Tea framework dependencies to go.mod following REQ-1 by running go get for bubbletea v0.25.0, bubbles v0.18.0, and lipgloss v0.9.1, then running go mod tidy and go build to verify | Restrictions: Do not modify existing dependencies, do not use -u flag, maintain go.mod formatting | Success: Dependencies added with correct versions, go mod download completes without errors, project builds successfully_

- [x] 2. Create TUI entry point detection in main.go
  - File: main.go, cmd/tui/tui.go
  - Add logic to detect no-args invocation and route to TUI mode
  - Purpose: Enable TUI launch when user runs pass-cli without arguments
  - _Leverage: Existing cmd.Execute() pattern, Go os.Args for argument detection_
  - _Requirements: REQ-1, REQ-12_
  - _Prompt: Role: Go Developer specializing in CLI application architecture | Task: Modify main.go to check len(os.Args)==1 and call tui.Run() if true, otherwise call cmd.Execute() for CLI mode following REQ-1 and REQ-12, create cmd/tui/tui.go with stub Run() function | Restrictions: Must not break existing CLI behavior, preserve all existing functionality, TUI only when exactly zero arguments | Success: Running pass-cli without args launches TUI stub, running with any args executes CLI commands unchanged_

- [x] 3. Create main TUI model and initialization
  - File: cmd/tui/model.go, cmd/tui/tui.go, cmd/tui/messages.go, cmd/tui/keys.go
  - Implement core Bubble Tea model structure with state management and vault unlocking
  - Purpose: Establish foundation for all TUI views and state transitions
  - _Leverage: Bubble Tea Model-Update-View pattern, vault.VaultService, cmd.GetVaultPath(), cmd/helpers.go readPassword()_
  - _Requirements: REQ-1, REQ-2_
  - _Prompt: Role: Go TUI Developer with Bubble Tea expertise | Task: Create Model struct with state fields, implement Init/Update/View methods following Bubble Tea pattern per REQ-1 and REQ-2, implement vault unlock flow trying keychain first then password prompt | Restrictions: Follow Elm architecture strictly, use VaultService API only, return placeholder strings for View() for now, simple state transitions only | Success: Model contains necessary state fields, Init unlocks vault with keychain support, Update handles quit keys, TUI launches and unlocks vault, can quit with q_

- [x] 4. Implement list view with search
  - File: cmd/tui/views/list.go, cmd/tui/components/searchbar.go, cmd/tui/update.go, cmd/tui/view.go, cmd/tui/model.go
  - Create credential list view with real-time search filtering and keyboard navigation
  - Purpose: Display credentials in searchable table with vim-style navigation
  - _Leverage: bubbles/list for credential display, bubbles/textinput for search bar, formatRelativeTime() from cmd/list.go, VaultService.ListCredentialsWithMetadata()_
  - _Requirements: REQ-3_
  - _Prompt: Role: Frontend Developer specializing in TUI interfaces | Task: Implement searchable credential list using bubbles/list and textinput following REQ-3, implement case-insensitive filtering on service and username, support arrow/vim keys, / to focus search, Escape to clear search | Restrictions: Use bubbles components not custom list, maintain focus management between search and list, do not implement detail view navigation yet | Success: List displays all credentials, search filters real-time, / focuses search, Escape clears, arrow/j/k navigate, Tab switches focus, shows no results message when appropriate_

- [ ] 5. Implement detail view
  - File: cmd/tui/views/detail.go, cmd/tui/update.go, cmd/tui/view.go, cmd/tui/model.go
  - Create credential detail view showing full information with masked password and copy-to-clipboard
  - Purpose: Display complete credential details with password toggle and clipboard integration
  - _Leverage: bubbles/viewport for scrollable content, atotto/clipboard for copying, VaultService.GetCredential(), formatRelativeTime()_
  - _Requirements: REQ-4_
  - _Prompt: Role: UI/UX Developer specializing in secure information presentation | Task: Implement detail view showing service, username, masked password, notes, timestamps, usage records following REQ-4, support m to toggle password visibility, c to copy password, display usage records in table | Restrictions: Passwords masked by default with asterisks, use existing clipboard library, use bubbles/viewport for scrolling, do not implement edit/delete yet | Success: Enter on list item navigates to detail, password masked by default, m toggles visibility, c copies to clipboard with notification, usage records displayed in table, Escape returns to list_

- [ ] 6. Implement add credential form
  - File: cmd/tui/views/form_add.go, cmd/tui/update.go, cmd/tui/view.go, cmd/tui/model.go
  - Create interactive form for adding new credentials with validation and password generation
  - Purpose: Enable credential creation with Tab navigation and field validation
  - _Leverage: bubbles/textinput for text fields, bubbles/textarea for notes, password generation from cmd/generate.go, VaultService.AddCredential()_
  - _Requirements: REQ-5_
  - _Prompt: Role: Forms Developer with input validation expertise | Task: Create add form with service, username, password, notes fields following REQ-5, implement Tab/Shift+Tab navigation, validate service name required and unique, support g key for password generation, Ctrl+S to save | Restrictions: Use bubbles textinput and textarea, validate before saving, show inline errors, support password generation with g key | Success: a key opens form, Tab navigates fields, focused field highlighted, g generates password with notification, Ctrl+S saves with validation, Escape shows discard confirmation, successful save returns to list with new item_

- [ ] 7. Implement edit credential form
  - File: cmd/tui/views/form_edit.go, cmd/tui/update.go, cmd/tui/view.go, cmd/tui/model.go
  - Create form for editing existing credentials with usage warnings
  - Purpose: Enable credential updates with pre-filled values and change tracking
  - _Leverage: form_add.go patterns, bubbles textinput/textarea, VaultService.UpdateCredential(), password generation from cmd/generate.go_
  - _Requirements: REQ-6_
  - _Prompt: Role: Forms Developer with update workflow expertise | Task: Create edit form pre-filled with current values following REQ-6, make service name read-only, show usage warning if credential has usage records, track changes for discard confirmation, support password generation | Restrictions: Service name read-only, pre-fill all fields, show usage warning with locations, only confirm discard if changes made | Success: e opens edit form pre-filled, service name read-only, other fields editable, Tab navigation works, g generates password, Ctrl+S saves, Escape confirms only if changed, usage warning displays locations, UpdatedAt updated after save_

- [ ] 8. Implement delete confirmation dialog
  - File: cmd/tui/views/confirm.go, cmd/tui/update.go, cmd/tui/view.go, cmd/tui/model.go
  - Create confirmation dialog for credential deletion with usage warnings
  - Purpose: Prevent accidental deletions with typed confirmation for high-risk deletions
  - _Leverage: bubbles/textinput for typed confirmation, VaultService.DeleteCredential(), lipgloss for modal styling_
  - _Requirements: REQ-7_
  - _Prompt: Role: Safety Engineer specializing in destructive action prevention | Task: Create delete confirmation dialog following REQ-7 with simple y/n for credentials without usage, typed service name confirmation for credentials with usage, display usage locations in red/warning color | Restrictions: Two modes simple and typed, usage locations prominently displayed, typed name must match exactly, center dialog as modal overlay | Success: d shows confirmation, displays service and username, simple y/n for no usage, typed confirmation for credentials with usage, usage locations in red, y confirms simple, exact service name confirms typed, n or Escape cancels, wrong name shows error, successful deletion returns to list_

- [ ] 9. Create status bar with keychain indicator
  - File: cmd/tui/components/statusbar.go, cmd/tui/view.go, cmd/tui/model.go
  - Create persistent status bar showing keychain status, credential count, current view, shortcuts
  - Purpose: Display system status and context-relevant keyboard shortcuts
  - _Leverage: KeychainService.IsAvailable() for status, lipgloss for styling and layout_
  - _Requirements: REQ-8_
  - _Prompt: Role: UI Developer specializing in status displays | Task: Create status bar showing keychain indicator, credential count, current view, and shortcuts following REQ-8, left side shows keychain status and count, center shows view name, right side shows relevant shortcuts per state | Restrictions: Visible in all views at bottom, update real-time when data changes, use emoji icons, keep shortcuts relevant to current view, fill terminal width | Success: Status bar visible at bottom in all views, keychain available shows green lock, unavailable shows yellow, credential count updates on add/delete, view name displayed, shortcuts change per state, fills full width_

- [ ] 10. Create help overlay
  - File: cmd/tui/views/help.go, cmd/tui/update.go, cmd/tui/view.go, cmd/tui/model.go
  - Create help screen showing all keyboard shortcuts organized by category
  - Purpose: Make features discoverable without reading documentation
  - _Leverage: bubbles/viewport for scrollable content, lipgloss for overlay styling and centering_
  - _Requirements: REQ-9_
  - _Prompt: Role: Documentation Engineer specializing in interactive help systems | Task: Create help overlay displaying shortcuts organized by category following REQ-9, show global shortcuts plus current view shortcuts, two-column layout with Key and Action, scrollable if needed | Restrictions: Accessible from any view with ? or F1, display as modal overlay dimming background, organize in logical groups, any key closes overlay | Success: ? or F1 shows help from any view, background dimmed, title displayed, global shortcuts always shown, view-specific shortcuts shown based on state, two-column format, scrollable with indicators if too large, any key closes and returns to previous view_

- [ ] 11. Create theme and styling
  - File: cmd/tui/styles/theme.go, cmd/tui/views/*.go, cmd/tui/components/*.go
  - Centralize all styling in theme file and apply throughout TUI
  - Purpose: Ensure visual consistency and professional appearance with 4.5:1 contrast
  - _Leverage: lipgloss for all styling, design.md color palette_
  - _Requirements: REQ-10_
  - _Prompt: Role: UI Designer with terminal color scheme expertise | Task: Create centralized theme file with all styling following REQ-10, use Go cyan as primary brand color, maintain 4.5:1 contrast ratio, apply consistently throughout all views | Restrictions: All colors in one place, graceful degradation for basic terminals, rounded borders for modals, subtle colors for decorative elements, bold colors only for important info | Success: theme.go created with color constants and style definitions, all views import and use theme styles, selected items use bold primary background, errors in red, success in green, keychain indicator uses success/warning colors, borders use subtle gray, consistent padding, professional cohesive look_

- [ ] 12. Add TUI unit tests
  - File: cmd/tui/model_test.go, cmd/tui/update_test.go, cmd/tui/views/list_test.go, cmd/tui/views/detail_test.go, cmd/tui/views/form_add_test.go, cmd/tui/components/statusbar_test.go
  - Create unit tests for TUI components and logic
  - Purpose: Ensure reliability and catch regressions with 80%+ coverage
  - _Leverage: Go testing package with table-driven tests, mock VaultService patterns from existing tests_
  - _Requirements: All_
  - _Prompt: Role: QA Engineer with Go testing expertise | Task: Create comprehensive unit tests for TUI covering all requirements, test state transitions, key handling, search filtering, form validation, password toggling | Restrictions: Mock VaultService for all tests, test components in isolation, use table-driven tests, aim for 80%+ coverage | Success: State transition tests, key handling tests for all states, search filter tests, form validation tests, form navigation tests, password visibility toggle tests, status bar rendering tests, confirmation dialog tests, all tests pass with go test, coverage ≥80%_

- [ ] 13. Add TUI integration tests
  - File: test/tui_integration_test.go
  - Create end-to-end integration tests for TUI workflows with real vault and keychain
  - Purpose: Verify TUI works correctly with real VaultService and system keychain
  - _Leverage: Existing integration test patterns from test/integration_test.go, build tag integration, existing TestMain for binary building_
  - _Requirements: All_
  - _Prompt: Role: Integration Test Engineer with end-to-end testing expertise | Task: Create integration tests for TUI covering full workflows following all requirements, test TUI launch, vault unlock, keychain integration, arg detection | Restrictions: Use build tag integration, use real VaultService with temp vaults, test with real keychain and skip if unavailable, clean up temp vaults and keychain entries | Success: TestIntegration_TUILaunch verifies TUI starts, TestIntegration_TUINoArgs verifies no args launches TUI, TestIntegration_TUIWithArgs verifies args launch CLI, TestIntegration_TUIKeychainIndicator verifies keychain detection, tests create temp vaults, cleanup with defer, tests skip gracefully if keychain unavailable, all tests pass with go test -v -tags=integration_
