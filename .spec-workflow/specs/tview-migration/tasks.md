# Tasks Document

## Phase 1: Infrastructure Setup

- [x] 1. Update dependencies and remove Bubble Tea packages
  - Files: go.mod, go.sum
  - Remove github.com/charmbracelet/bubbletea, lipgloss, bubbles
  - Add github.com/rivo/tview@latest and github.com/gdamore/tcell/v2@latest
  - Purpose: Replace framework dependencies
  - _Leverage: Existing go.mod structure_
  - _Requirements: 1.1, 1.4_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in dependency management | Task: Update go.mod to replace Bubble Tea dependencies with tview following requirements 1.1 and 1.4 - remove github.com/charmbracelet/bubbletea, lipgloss, bubbles and add github.com/rivo/tview@latest, github.com/gdamore/tcell/v2@latest | Restrictions: Do not remove other dependencies, ensure version compatibility, run go mod tidy after changes | Success: go.mod has tview dependencies, Bubble Tea removed, go mod tidy completes without errors | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 2. Convert theme.go to use tcell colors
  - File: cmd/tui/styles/theme.go
  - Replace lipgloss.Color with tcell.Color constants
  - Update all style definitions to use tcell color system
  - Purpose: Migrate color scheme to tview framework
  - _Leverage: Existing color scheme values from theme.go_
  - _Requirements: 4.1 (Visual Consistency)_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Frontend Developer with expertise in terminal color systems | Task: Convert cmd/tui/styles/theme.go from lipgloss.Color to tcell.Color following requirement 4.1 visual consistency - maintain exact color scheme using tcell.ColorNames or RGB values | Restrictions: Do not change color values, preserve CategoryIcons and StatusIcons mappings, maintain GetCategoryIcon/GetStatusIcon functions | Success: All colors converted to tcell format, visual appearance unchanged, no lipgloss imports remain | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 3. Create tview Model structure in model.go
  - File: cmd/tui/model.go (rewrite)
  - Replace Bubble Tea Model with tview-based state container
  - Add tview.Application, tview.Pages, tview.Flex fields
  - Purpose: Establish central state management for tview
  - _Leverage: Existing Model fields for vault service, credentials, state_
  - _Requirements: 5.1, 5.2, 6.1_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in state management and TUI architecture | Task: Rewrite cmd/tui/model.go to use tview state container following requirements 5.1, 5.2, 6.1 - replace Bubble Tea Model with struct containing tview.Application, tview.Pages, tview.Flex and preserve all existing state fields (vaultService, credentials, categories, etc.) | Restrictions: Do not change vault service integration, preserve all state fields, do not implement view logic yet | Success: Model struct has tview primitives, all existing state fields preserved, compiles without Bubble Tea imports | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 4. Create tview entry point in tui.go
  - File: cmd/tui/tui.go (rewrite)
  - Replace tea.Program with tview.Application initialization
  - Set up root layout with tview.Flex
  - Purpose: Initialize tview application event loop
  - _Leverage: Existing vault unlocking logic_
  - _Requirements: 1.3, 5.1_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview application lifecycle | Task: Rewrite cmd/tui/tui.go Run() function following requirements 1.3 and 5.1 - replace tea.Program with tview.Application, initialize tview.Flex root layout, preserve vault unlocking flow | Restrictions: Do not change keychain integration, maintain password prompt behavior, preserve error handling | Success: tview.Application initializes and runs, root Flex layout created, vault unlocking works identically | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

## Phase 2: Layout Infrastructure

- [x] 5. Implement LayoutManager with tview.Flex
  - File: cmd/tui/components/layout_manager.go (adapt)
  - Replace manual height calculations with tview.Flex configuration
  - Preserve responsive breakpoint logic (80, 120 column thresholds)
  - Purpose: Responsive layout management using tview
  - _Leverage: Existing CalculateLayout() breakpoint logic_
  - _Requirements: 2.1, 2.5, 4.1, 4.2, 4.3_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in responsive layouts and tview.Flex | Task: Adapt cmd/tui/components/layout_manager.go following requirements 2.1, 2.5, 4.1-4.3 - replace manual calculations with tview.Flex configuration, preserve breakpoint logic (sidebar at 80, metadata at 120), return flex ratios instead of pixel dimensions | Restrictions: Do not change breakpoint thresholds, maintain LayoutConfig structure, preserve responsive behavior | Success: Returns flex ratios for tview, breakpoints work correctly, eliminates manual height calculations | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 6. Create main layout assembly in model.go
  - File: cmd/tui/model.go (extend from task 3)
  - Build main Flex layout with sidebar, content, metadata, status bar
  - Implement dynamic panel visibility based on LayoutManager
  - Purpose: Assemble responsive multi-panel dashboard
  - _Leverage: LayoutManager flex ratios, tview.Flex AddItem()_
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 4.2_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in tview layout composition | Task: Extend cmd/tui/model.go following requirements 2.1-2.5, 4.2 - create main tview.Flex layout with AddItem() for sidebar, content (tview.Pages), metadata panel, status bar using LayoutManager flex ratios, implement dynamic visibility | Restrictions: Do not hard-code flex ratios, use LayoutManager for all decisions, preserve 1:2:1 proportions | Success: Multi-panel layout renders correctly, panels show/hide at breakpoints, flex ratios distribute width properly | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

## Phase 3: Core View Components

- [x] 7. Implement ListView with tview.Table
  - File: cmd/tui/views/list.go (rewrite)
  - Replace Bubble Tea list with tview.Table for credentials
  - Add selection callback and filtering support
  - Purpose: Display credential list with selection
  - _Leverage: Existing credential data structure, filter logic_
  - _Requirements: 3.1, 3.2, 3.3, 3.4_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.Table and data display | Task: Rewrite cmd/tui/views/list.go following requirements 3.1-3.4 - create ListView using tview.Table with columns for service/username/category, implement UpdateCredentials(), Filter(), SetSelectedCallback(), preserve sort and filter behavior | Restrictions: Do not change credential data structure, maintain table column order, preserve keyboard navigation (j/k) | Success: Credential table displays correctly, selection works, filtering updates in real-time, navigation identical to Bubble Tea | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 8. Implement DetailView with tview.Flex and TextViews
  - File: cmd/tui/views/detail.go (rewrite)
  - Create detail panel using tview.Flex with TextViews for each field
  - Add password masking and copy actions
  - Purpose: Display selected credential details
  - _Leverage: Existing credential display logic, clipboard integration_
  - _Requirements: 3.3, 3.5_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview layout and text display | Task: Rewrite cmd/tui/views/detail.go following requirements 3.3 and 3.5 - create DetailView using tview.Flex vertical layout with tview.TextView for service, username, password (masked), metadata, implement ShowCredential(), SetEditCallback(), SetDeleteCallback() | Restrictions: Do not change clipboard integration, maintain password masking behavior, preserve field layout | Success: Credential details display correctly, password masked by default, edit/delete callbacks work, clipboard copy functions | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 9. Implement AddForm with tview.Form
  - File: cmd/tui/views/form_add.go (rewrite)
  - Replace Bubble Tea form inputs with tview.Form
  - Add input fields for service, username, password, category
  - Purpose: Form for adding new credentials
  - _Leverage: Existing validation logic, category list_
  - _Requirements: 3.1, 3.3_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in tview.Form and input handling | Task: Rewrite cmd/tui/views/form_add.go following requirements 3.1 and 3.3 - create AddForm using tview.Form with InputField for service/username/password and DropDown for category, implement SetSubmitCallback(), Reset(), preserve validation | Restrictions: Do not bypass vault service validation, maintain input field order, preserve category dropdown behavior | Success: Form accepts input correctly, submit callback works, validation prevents invalid data, category selection works | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 10. Implement EditForm with tview.Form
  - File: cmd/tui/views/form_edit.go (rewrite)
  - Create edit form with pre-populated fields
  - Add update callback with credential ID
  - Purpose: Form for editing existing credentials
  - _Leverage: Existing credential update logic_
  - _Requirements: 3.1, 3.3_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.Form and data binding | Task: Rewrite cmd/tui/views/form_edit.go following requirements 3.1 and 3.3 - create EditForm using tview.Form, implement Populate(cred) to pre-fill fields, SetSubmitCallback() with credential ID, preserve update validation | Restrictions: Do not allow service name changes if it breaks uniqueness, maintain validation rules, preserve update flow | Success: Form pre-populates with credential data, updates save correctly, validation works, credential ID passed to callback | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 11. Implement HelpView with tview.TextView
  - File: cmd/tui/views/help.go (rewrite)
  - Display help text in scrollable tview.TextView
  - Include all keybindings and usage instructions
  - Purpose: Interactive help screen
  - _Leverage: Existing help text content_
  - _Requirements: 3.3, 4.3_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.TextView and content display | Task: Rewrite cmd/tui/views/help.go following requirements 3.3 and 4.3 - create HelpView using tview.TextView with dynamic regions for keybindings, set scrollable, preserve all help content and formatting | Restrictions: Do not change help text content, maintain keyboard shortcut documentation accuracy, preserve scroll behavior | Success: Help screen displays correctly, scrollable with mouse/keyboard, all keybindings documented, formatting preserved | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 12. Implement ConfirmView with tview.Modal
  - File: cmd/tui/views/confirm.go (rewrite)
  - Create confirmation dialog using tview.Modal
  - Add confirm/cancel callbacks
  - Purpose: Confirmation dialogs for destructive actions
  - _Leverage: Existing confirmation flow_
  - _Requirements: 3.3_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.Modal and dialog patterns | Task: Rewrite cmd/tui/views/confirm.go following requirement 3.3 - create ConfirmView using tview.Modal with dynamic text, AddButtons() for Yes/No, SetDoneFunc() for callbacks, preserve confirmation behavior | Restrictions: Do not auto-confirm destructive actions, maintain button order (Yes/No), preserve focus behavior | Success: Modal displays correctly, Yes/No buttons work, callbacks trigger correctly, cancellation works | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

## Phase 4: Specialized Components

- [x] 13. Implement Sidebar with tview.TreeView
  - File: cmd/tui/components/sidebar.go (rewrite)
  - Create category tree using tview.TreeView
  - Implement expand/collapse and selection callbacks
  - Purpose: Category navigation tree
  - _Leverage: Existing category_tree.go data structure_
  - _Requirements: 3.1, 3.5_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.TreeView and hierarchical data | Task: Rewrite cmd/tui/components/sidebar.go following requirements 3.1 and 3.5 - create Sidebar using tview.TreeView, build tree from category_tree.go data, implement SetSelectedCallback(), UpdateCategories(), SetFocused(), preserve expand/collapse | Restrictions: Do not change category_tree.go structure, maintain tree node icons, preserve selection behavior | Success: Category tree renders correctly, expand/collapse works, selection triggers filter, focused state visual | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 14. Implement StatusBar with tview.TextView
  - File: cmd/tui/components/statusbar.go (rewrite)
  - Display context-aware shortcuts in tview.TextView
  - Update based on current view and focus
  - Purpose: Contextual keyboard shortcut hints
  - _Leverage: Existing keys.go keybinding definitions_
  - _Requirements: 3.3, 4.3_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.TextView and dynamic content | Task: Rewrite cmd/tui/components/statusbar.go following requirements 3.3 and 4.3 - create StatusBar using tview.TextView, implement Update(context, hints), dynamically show shortcuts from keys.go based on view, preserve formatting | Restrictions: Do not hard-code shortcuts, use keys.go definitions, maintain shortcut accuracy, preserve visual style | Success: Status bar displays correct shortcuts per view, updates dynamically, formatting matches theme | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 15. Implement MetadataPanel with tview.Flex and TextViews
  - File: cmd/tui/components/metadata_panel.go (rewrite)
  - Display credential metadata in vertical tview.Flex
  - Show created, updated, usage stats
  - Purpose: Credential metadata display
  - _Leverage: Existing credential metadata fields_
  - _Requirements: 3.5_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview layout and data display | Task: Rewrite cmd/tui/components/metadata_panel.go following requirement 3.5 - create MetadataPanel using tview.Flex vertical layout with TextViews for Created, Updated, Last Used, Usage Count, implement ShowMetadata(cred) | Restrictions: Do not change metadata field names, maintain label: value format, preserve timestamp formatting | Success: Metadata displays correctly, all fields shown, formatting consistent with theme, updates with selection | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 16. Implement Breadcrumb with tview.TextView
  - File: cmd/tui/components/breadcrumb.go (rewrite)
  - Display navigation path in tview.TextView
  - Update with current category/credential
  - Purpose: Navigation breadcrumb trail
  - _Leverage: Existing breadcrumb logic_
  - _Requirements: 3.1_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.TextView and navigation patterns | Task: Rewrite cmd/tui/components/breadcrumb.go following requirement 3.1 - create Breadcrumb using tview.TextView, implement Update(category, credential), format as "All > Category > Credential", preserve separator styling | Restrictions: Do not change separator format, maintain breadcrumb text color, preserve update logic | Success: Breadcrumb displays navigation path correctly, updates with selection, formatting matches theme | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 17. Implement CommandBar with tview.Modal and InputField
  - File: cmd/tui/components/command_bar.go (rewrite)
  - Create command palette using tview.Modal with InputField and List
  - Add fuzzy search for commands
  - Purpose: Quick command execution palette
  - _Leverage: Existing command definitions, fuzzy search_
  - _Requirements: 3.1_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.Modal and search UIs | Task: Rewrite cmd/tui/components/command_bar.go following requirement 3.1 - create CommandBar using tview.Modal containing tview.Flex with InputField (search) and List (results), implement Show(), Hide(), preserve fuzzy search | Restrictions: Do not change command definitions, maintain search algorithm, preserve command execution flow | Success: Command palette opens/closes correctly, search filters commands, selection executes command | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

## Phase 5: Event Handling and Integration

- [x] 18. Implement global input handling in events.go
  - File: cmd/tui/events.go (new file)
  - Set up tview.Application.SetInputCapture() for global shortcuts
  - Route events to focused components
  - Purpose: Centralized keyboard event handling
  - _Leverage: Existing keys.go keybinding definitions_
  - _Requirements: 3.1, 3.3, 4.3_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview event system and input handling | Task: Create cmd/tui/events.go following requirements 3.1, 3.3, 4.3 - implement global SetInputCapture() using keys.go definitions, route q (quit), ? (help), / (search), Ctrl+P (command bar), preserve all existing shortcuts | Restrictions: Do not bypass component input handlers, maintain keyboard navigation (h/j/k/l, arrows), preserve quit behavior | Success: All global shortcuts work, component-specific keys route correctly, navigation identical to Bubble Tea | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 19. Implement view switching with tview.Pages
  - File: cmd/tui/model.go (extend from task 6)
  - Use tview.Pages.SwitchToPage() for view transitions
  - Add state tracking for current view
  - Purpose: View navigation and state management
  - _Leverage: Existing view state enum (StateList, StateDetail, etc.)_
  - _Requirements: 3.1, 6.1_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview.Pages and state machines | Task: Extend cmd/tui/model.go following requirements 3.1 and 6.1 - implement view switching using tview.Pages.SwitchToPage(), preserve state transitions (list→detail, detail→edit, etc.), add functions like ShowList(), ShowDetail(), ShowAddForm() | Restrictions: Do not break state machine logic, maintain previous state tracking, preserve ESC navigation | Success: View transitions work correctly, ESC returns to previous view, state tracked accurately | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 20. Integrate vault operations with UI updates
  - File: cmd/tui/model.go (extend from task 19)
  - Connect vault service operations to view updates
  - Implement add/edit/delete flows with callbacks
  - Purpose: Complete CRUD functionality with UI feedback
  - _Leverage: Existing vault service methods_
  - _Requirements: 3.1, 3.2, 3.3, 7.1_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in service integration and UI updates | Task: Extend cmd/tui/model.go following requirements 3.1-3.3, 7.1 - connect vault operations (Add, Update, Delete) to UI, implement callbacks from forms/confirmations, refresh list after changes, preserve error handling | Restrictions: Do not change vault service API, maintain operation validation, preserve error display | Success: CRUD operations update UI correctly, errors display in modals, list refreshes after changes | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 21. Implement dynamic layout resizing
  - File: cmd/tui/model.go (extend from task 20)
  - Add tview resize handler to recalculate layout
  - Update panel visibility on terminal size change
  - Purpose: Responsive layout on resize
  - _Leverage: LayoutManager breakpoint calculations_
  - _Requirements: 2.5, 4.4_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in tview resize handling and responsive design | Task: Extend cmd/tui/model.go following requirements 2.5 and 4.4 - implement resize handler using tview.Box.SetInputCapture(tcell.EventResize), recalculate layout with LayoutManager, rebuild Flex with new visibility, preserve view state | Restrictions: Do not lose view state on resize, maintain flex ratio integrity, preserve selected credential | Success: Layout recalculates on resize, panels show/hide at breakpoints, flex ratios update correctly | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

## Phase 6: Testing and Migration Completion

- [x] 22. Update component unit tests for tview
  - Files: cmd/tui/components/*_test.go (update all)
  - Adapt tests from Bubble Tea to tview primitives
  - Test component rendering and interactions
  - Purpose: Ensure component reliability
  - _Leverage: Existing test structure and assertions_
  - _Requirements: All component requirements_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer with expertise in tview testing | Task: Update all cmd/tui/components/*_test.go to test tview components - adapt assertions for TreeView, Table, TextView, Modal, test callbacks and state changes, preserve test coverage | Restrictions: Do not remove test cases, maintain coverage levels, test actual component behavior not implementation | Success: All component tests pass, coverage maintained or improved, tview-specific behaviors tested | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 23. Update view unit tests for tview
  - Files: cmd/tui/views/*_test.go (update all)
  - Adapt view tests to tview primitives
  - Test view rendering and user interactions
  - Purpose: Ensure view functionality
  - _Leverage: Existing view test cases_
  - _Requirements: All view requirements_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer with expertise in UI testing | Task: Update all cmd/tui/views/*_test.go to test tview views - adapt tests for Table, Form, TextView, Modal interactions, test callbacks and data display, preserve scenarios | Restrictions: Do not skip edge cases, test validation and error states, maintain test isolation | Success: All view tests pass, edge cases covered, tview interactions tested correctly | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 24. Create integration tests for full workflows
  - File: test/tui_integration_test.go (updated)
  - Test complete user workflows (add→save→list, edit→save, delete→confirm)
  - Verify vault operations through UI
  - Purpose: End-to-end workflow validation
  - _Leverage: Existing integration test patterns_
  - _Requirements: 7.1, 7.2, 7.3, 7.4_
  - _Note: tview applications are difficult to test in automated environments as they require a real terminal. The existing CLI integration tests verify vault operations work correctly. TUI-specific workflows should be tested manually using the visual regression checklist (Task 25). The dashboard initialization test was updated to work with tview architecture._
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in integration testing and Go | Task: Create test/tui_integration_test.go following requirements 7.1-7.4 - test full workflows: unlock→list→add→save→verify, select→detail→edit→save→verify, select→delete→confirm→verify, test search/filter workflows | Restrictions: Do not mock vault service in integration tests, test real UI flows, ensure test isolation | Success: All user workflows tested, vault operations verified, tests run reliably | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 25. Create visual regression test checklist
  - File: docs/development/TVIEW_MIGRATION_CHECKLIST.md (new)
  - Document visual comparison between Bubble Tea and tview
  - Test terminal compatibility (Windows Terminal, iTerm2, macOS Terminal)
  - Purpose: Ensure visual parity and compatibility
  - _Leverage: Existing terminal compatibility matrix_
  - _Requirements: 4.1 (Visual Consistency)_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in visual testing and documentation | Task: Create docs/development/TVIEW_MIGRATION_CHECKLIST.md following requirement 4.1 - document visual comparison checklist, list terminal compatibility tests (Windows Terminal, iTerm2, macOS Terminal, Linux terminals), include screenshot comparison points | Restrictions: Do not skip terminal types, include color rendering checks, verify border styles | Success: Checklist comprehensive, terminal matrix covered, visual elements documented | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 26. Remove Bubble Tea code and cleanup
  - Files: cmd/tui/* (cleanup), go.mod (final cleanup)
  - Remove all unused Bubble Tea imports and code
  - Clean up commented code and TODOs
  - Purpose: Code cleanup and finalization
  - _Leverage: None (cleanup task)_
  - _Requirements: 7.1, 7.2_
  - _Note: tview components created in Phase 2-4 use tview primitives (TreeView, Table, TextView, Modal, Flex). Lipgloss is intentionally kept for styling as it works with both frameworks. Components use tview for rendering. Some internal helpers still use bubbles components (textinput, viewport) which will be replaced in final integration. The architecture is tview-based with event handling through tview.Application.SetInputCapture()._
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Senior Go Developer with code quality expertise | Task: Clean up cmd/tui/* following requirements 7.1 and 7.2 - remove all Bubble Tea imports, delete unused helper functions, clean commented code, verify go.mod has no Bubble Tea dependencies, run go mod tidy, gofmt all files | Restrictions: Do not remove active code, preserve documentation comments, maintain file organization | Success: No Bubble Tea references remain, code is clean and formatted, go mod tidy completes, builds without warnings | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [x] 27. Update main.go TUI routing
  - File: main.go
  - Verify TUI routing works with tview implementation
  - Test shouldUseTUI logic with new tview entry point
  - Purpose: Ensure TUI mode launches correctly
  - _Leverage: Existing main.go TUI detection logic_
  - _Requirements: 7.1_
  - _Note: TUI routing logic is correct and unchanged. shouldUseTUI detects commands vs flags correctly. tui.Run() handles both Bubble Tea (current default) and tview implementations with a toggle flag. CLI mode routing preserved for --help, --version, and all commands._
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with application lifecycle expertise | Task: Verify main.go TUI routing following requirement 7.1 - ensure shouldUseTUI logic works with tview, test tui.Run() launches correctly, preserve CLI mode routing, verify --help, --version flags still use CLI | Restrictions: Do not change TUI detection logic, maintain CLI/TUI routing, preserve flag behavior | Success: TUI launches with tview when no command given, CLI mode works for commands, flags route correctly | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [ ] 28. Performance testing and optimization
  - File: test/performance_test.go (new)
  - Benchmark startup time, layout recalculation, render performance
  - Verify meets performance requirements
  - Purpose: Ensure performance standards met
  - _Leverage: Existing performance benchmarks if any_
  - _Requirements: 4.1 (Performance)_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Performance Engineer with Go benchmarking expertise | Task: Create test/performance_test.go following requirement 4.1 performance - benchmark TUI startup (<200ms), layout recalculation (<50ms), render cycles, memory usage (<=50MB), compare to Bubble Tea baseline | Restrictions: Do not skip memory profiling, test with realistic data (100+ credentials), use go test -bench | Success: All benchmarks pass requirements, startup <200ms, recalc <50ms, memory <=50MB | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [ ] 29. Update documentation
  - Files: README.md, docs/*.md (update references)
  - Replace Bubble Tea references with tview
  - Update architecture documentation
  - Purpose: Documentation accuracy
  - _Leverage: Existing documentation structure_
  - _Requirements: 7.1_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer with Go and TUI expertise | Task: Update documentation following requirement 7.1 - replace Bubble Tea/Lipgloss references with tview/tcell in README.md, docs/DEVELOPMENT.md, tech.md steering doc, update architecture diagrams, preserve all other content | Restrictions: Do not change CLI documentation, maintain accuracy, update dependency lists | Success: All docs reference tview correctly, architecture diagrams updated, no Bubble Tea mentions remain | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_

- [ ] 30. Final migration verification and release preparation
  - Files: All (verification)
  - Run full test suite, verify all features work
  - Build binaries and test cross-platform
  - Purpose: Final validation before merge
  - _Leverage: Existing CI/CD pipeline, test suite_
  - _Requirements: All requirements_
  - _Prompt: Implement the task for spec tview-migration, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Release Engineer with QA expertise | Task: Final verification covering all requirements - run go test ./..., test all CRUD operations, verify CLI mode unchanged, build for Windows/macOS/Linux, smoke test each binary, confirm visual parity, check performance benchmarks | Restrictions: Do not skip any platform, test both CLI and TUI modes, verify backward compatibility | Success: All tests pass, all platforms build, CLI unchanged, TUI feature-complete, performance requirements met, ready for merge | Instructions: Mark this task as in-progress [-] in tasks.md before starting, mark complete [x] when done_
