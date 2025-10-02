# Tasks Document

- [x] 1. Create Layout Manager component
  - File: cmd/tui/components/layout_manager.go
  - Implement centralized dimension calculation system with responsive breakpoints
  - Purpose: Calculate panel dimensions based on terminal size and visibility states
  - _Leverage: Existing SetSize() pattern from views, WindowSizeMsg handling_
  - _Requirements: REQ-10, REQ-11_
  - _Prompt: Role: Go TUI Developer specializing in responsive layout systems | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create LayoutManager struct with Calculate() method following REQ-10 and REQ-11, implementing three responsive breakpoints (full: >=120 cols, medium: 80-119, small: <80), calculate panel dimensions (sidebar, main, metadata, process, command bar) respecting minimum constraints (sidebar 20 cols, main 40 cols, metadata 25 cols), return Layout struct with IsTooSmall flag when terminal < 60x20 | Restrictions: Pure calculation logic with no UI rendering, must handle panel visibility states, no external dependencies beyond standard library | _Leverage: None, new component | Success: LayoutManager.Calculate() returns correct dimensions for all breakpoints, respects minimum sizes, handles edge cases (very small terminals, many panels), unit tests pass for all scenarios_

- [x] 2. Create Category Tree utility functions
  - File: cmd/tui/components/category_tree.go
  - Implement pure functions for categorizing credentials by service patterns
  - Purpose: Organize credentials into categories (APIs, Cloud, Databases, Git, etc.)
  - _Leverage: vault.CredentialMetadata structure_
  - _Requirements: REQ-3_
  - _Prompt: Role: Go Developer specializing in data transformation and pattern matching | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create CategorizeCredentials() function following REQ-3 that takes []vault.CredentialMetadata and returns []Category, implement pattern matching for 8 categories (APIs & Services, Cloud Infrastructure, Databases, Version Control, Communication, Payment Processing, AI Services, Uncategorized), create GetCategoryIcon() and GetStatusIcon() functions returning icon strings, use case-insensitive substring matching on service names | Restrictions: Pure functions with no state, must handle empty inputs gracefully, all credentials must be categorized (use Uncategorized as fallback) | _Leverage: vault.CredentialMetadata type | Success: CategorizeCredentials() correctly categorizes test credentials, icons are defined for all categories, unit tests cover all category patterns and edge cases_

- [x] 3. Extend theme with panel styles and icons
  - File: cmd/tui/styles/theme.go
  - Add panel border styles and icon mappings for dashboard
  - Purpose: Provide consistent styling for focused/unfocused panels and category icons
  - _Leverage: Existing theme styles (PrimaryColor, SubtleColor, BorderStyle)_
  - _Requirements: REQ-12_
  - _Prompt: Role: UI Designer with terminal color scheme expertise | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Extend existing theme.go following REQ-12 by adding ActivePanelBorderStyle (cyan border), InactivePanelBorderStyle (gray border), CategoryIcons map (☁️ Cloud, 🔑 APIs, 💾 Databases, 📦 Git, 📧 Comm, 💰 Payment, 🤖 AI, 📁 Uncat), StatusIcons map (⏳ pending, ✓ success, ✗ failed, ▶ collapsed, ▼ expanded), ensure 4.5:1 contrast ratio, graceful ASCII fallback for unsupported terminals | Restrictions: Do not modify existing styles, only add new ones, must maintain backward compatibility | _Leverage: Existing theme constants and styles | Success: New panel styles render correctly with focus highlighting, icons display properly or fallback gracefully, all existing styles continue to work, theme tests pass_

- [x] 4. Create Sidebar Panel component
  - File: cmd/tui/components/sidebar.go
  - Implement sidebar with category tree, stats, and quick actions
  - Purpose: Display navigable category tree with credential counts and statistics
  - _Leverage: bubbles/viewport for scrolling, category_tree.go for categorization, styles/theme.go_
  - _Requirements: REQ-3, REQ-4, REQ-14_
  - _Prompt: Role: Go TUI Developer specializing in interactive navigation components | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create SidebarPanel struct following REQ-3, REQ-4, REQ-14 with fields (categories, selectedCategory, selectedCred, stats, viewport), implement NewSidebarPanel(), SetSize(), Update(), View(), GetSelectedCredential(), GetSelectedCategory() methods, support j/k/arrow navigation, Enter/l to expand/collapse categories, render category tree with expand/collapse indicators (▶/▼), display stats section (Total, Used, Recent) and quick actions ([a] Add, [:] Command, [?] Help), use viewport for scrolling | Restrictions: Must use bubbles/viewport, only render visible content, handle empty categories gracefully | _Leverage: viewport pattern from DetailView, theme.go styles | Success: Sidebar renders category tree correctly, navigation works smoothly, expansion/collapse state maintained, stats update when credentials change, scrolling works for long category lists_

- [x] 5. Create Metadata Panel component
  - File: cmd/tui/components/metadata_panel.go
  - Implement right-side panel showing credential details
  - Purpose: Display selected credential information with usage records
  - _Leverage: bubbles/viewport for scrolling, DetailView rendering logic_
  - _Requirements: REQ-6_
  - _Prompt: Role: Go TUI Developer specializing in information display components | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create MetadataPanel struct following REQ-6 with fields (credential, passwordMasked, viewport), implement NewMetadataPanel(), SetCredential(), SetSize(), Update(), View(), TogglePasswordMask() methods, render credential service (header), username, password (masked by default), created/updated timestamps (relative format), usage records table (path, access count, last accessed), support m key to toggle password mask, use viewport for scrolling long content | Restrictions: Password masked by default, must use viewport for overflow, display "(not set)" for missing fields | _Leverage: DetailView's rendering logic, formatTime() function, viewport pattern | Success: Metadata panel displays all credential fields correctly, password masking works, usage records formatted in readable table, scrolling works for long content, updates when credential selection changes_

- [x] 6. Create Process Panel component
  - File: cmd/tui/components/process_panel.go
  - Implement bottom panel showing async operation feedback
  - Purpose: Display operation status (generating password, saving, deleting)
  - _Leverage: Notification patterns from existing views_
  - _Requirements: REQ-8_
  - _Prompt: Role: Go TUI Developer specializing in status notification systems | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create ProcessPanel struct following REQ-8 with fields (processes []Process, maxDisplay 5), Process struct (ID, Description, Status, Error, Timestamp), implement NewProcessPanel(), AddProcess(), UpdateProcess(), SetSize(), View() methods, render most recent 5 processes with status icons (⏳ running, ✓ success, ✗ failed), use green for success, red for errors, display process descriptions and timestamps, auto-hide panel when all processes complete after 3 seconds | Restrictions: Maximum 5 visible processes, must trim old processes, simple list rendering | _Leverage: styles/theme.go status icons, notification patterns | Success: Process panel displays active operations, status updates reflect correctly, success/error colors applied, auto-hide works after completion, processes don't overflow allocated space_

- [x] 7. Create Command Bar component
  - File: cmd/tui/components/command_bar.go
  - Implement vim-style command input bar
  - Purpose: Execute commands like :add, :search, :category, :help, :quit
  - _Leverage: bubbles/textinput, existing form patterns_
  - _Requirements: REQ-7_
  - _Prompt: Role: Go Developer specializing in command parsing and input handling | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create CommandBar struct following REQ-7 with fields (input textinput.Model, history []string, historyIdx, error), Command struct (Name, Args []string), implement NewCommandBar(), Focus(), Blur(), SetSize(), Update(), View(), GetCommand(), SetError() methods, parse commands (:add [service], :search [query], :category [name], :help/:h, :quit/:q), support command history with up/down arrows, display : prompt, show errors in red below input, return parsed Command on Enter | Restrictions: Must use bubbles/textinput, commands start with :, Enter executes, Esc cancels | _Leverage: textinput patterns from forms, command parsing from main.go | Success: Command bar opens with : key, accepts text input, parses commands correctly, history navigation works, errors display properly, executes or cancels correctly_

- [x] 8. Create Breadcrumb component
  - File: cmd/tui/components/breadcrumb.go
  - Implement navigation path display
  - Purpose: Show current location (Home > APIs > Cloud > aws-prod)
  - _Leverage: lipgloss for truncation, existing text rendering_
  - _Requirements: REQ-5_
  - _Prompt: Role: Go TUI Developer specializing in navigation UI components | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create Breadcrumb struct following REQ-5 with fields (path []string, width), implement NewBreadcrumb(), SetPath(), SetSize(), View() methods, render path segments joined with " > " separator, truncate middle segments with "..." when path exceeds panel width (e.g., "Home > ... > aws-prod"), use bold or colored styling to stand out from content, ensure minimum 3 segments visible (first, ..., last) when truncated | Restrictions: Must fit within allocated width, handle empty path gracefully, simple text rendering | _Leverage: lipgloss text truncation, styles/theme.go | Success: Breadcrumb displays path correctly, truncation works when width constrained, styling makes it visually distinct, updates when navigation changes_

- [x] 9. Extend Model with panel state management
  - File: cmd/tui/model.go
  - Add panel state fields and initialization to existing Model
  - Purpose: Manage panel visibility, focus, and instances
  - _Leverage: Existing Model structure, initialization patterns_
  - _Requirements: REQ-1, REQ-2, REQ-9, REQ-13_
  - _Prompt: Role: Go TUI Developer with Bubble Tea Model-Update-View expertise | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Extend existing Model struct following REQ-1, REQ-2, REQ-9, REQ-13 by adding panel fields (layoutManager *LayoutManager, sidebar *SidebarPanel, metadataPanel *MetadataPanel, processPanel *ProcessPanel, commandBar *CommandBar, breadcrumb *Breadcrumb), panel state fields (panelFocus PanelFocus, sidebarVisible bool, metadataVisible bool, processVisible bool, commandBarOpen bool), category fields (categories []Category, currentCategory string), create PanelFocus enum (FocusSidebar, FocusMain, FocusMetadata, FocusCommandBar), initialize all panels in NewModel(), categorize credentials when loaded | Restrictions: Do not modify existing Model fields, only add new ones, maintain backward compatibility | _Leverage: Existing Model structure, NewModel() pattern | Success: Model compiles with new fields, panels initialized correctly, existing TUI functionality unaffected, panel states managed properly_

- [x] 10. Implement panel toggle key handling
  - File: cmd/tui/model.go (continue), cmd/tui/keys.go (if exists)
  - Add keyboard shortcuts for panel visibility toggles
  - Purpose: Enable s/m/p/f keys to show/hide panels
  - _Leverage: Existing key handling in Update() method_
  - _Requirements: REQ-2, REQ-13_
  - _Prompt: Role: Go TUI Developer specializing in keyboard input handling | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Add panel toggle logic to Update() method following REQ-2, REQ-13, handle s key (toggle sidebar), m key (toggle metadata, but preserve existing m for password mask in detail view), p key (toggle process panel), f key (toggle all footer panels), Tab key (switch panel focus), Shift+Tab (previous panel focus), check hasInputFocus() before handling to avoid conflicts with forms/search, update panelFocus using nextPanelFocus() helper, trigger layout recalculation on visibility changes | Restrictions: Must not break existing key handlers, check input focus state, Tab only works when no input focused | _Leverage: Existing tea.KeyMsg handling, hasInputFocus() pattern | Success: Panel toggle keys work correctly, Tab switches focus between visible panels, existing keys (q, ?, /, a, e, d) still work, no conflicts with form input_

- [x] 11. Implement command bar integration
  - File: cmd/tui/model.go (continue)
  - Add command bar opening, command execution logic
  - Purpose: Handle : key to open command bar and execute parsed commands
  - _Leverage: CommandBar.GetCommand(), existing navigation logic_
  - _Requirements: REQ-7_
  - _Prompt: Role: Go Developer specializing in command execution and state management | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Integrate CommandBar following REQ-7, handle : key to open command bar (set commandBarOpen = true, commandBar.Focus()), delegate Update() to commandBar when open, on Enter get Command from commandBar.GetCommand(), execute commands (:add opens form with pre-filled service, :search activates search with query, :category navigates to category, :help opens help, :quit exits), close command bar on Esc or after execution, display errors using commandBar.SetError() | Restrictions: Must not interfere with existing command handling, validate command before execution, close bar after successful command | _Leverage: Existing state transitions (StateAdd, StateList), search activation pattern | Success: : key opens command bar, commands parse correctly, :add/:search/:category/:help/:quit execute as specified, errors display in command bar, Esc cancels without executing_

- [x] 12. Implement Layout Manager integration in Update
  - File: cmd/tui/model.go (continue)
  - Add layout recalculation on WindowSizeMsg and panel visibility changes
  - Purpose: Dynamically calculate panel dimensions and propagate to all components
  - _Leverage: LayoutManager.Calculate(), existing SetSize() calls_
  - _Requirements: REQ-10, REQ-11_
  - _Prompt: Role: Go TUI Developer specializing in responsive layout systems | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Integrate LayoutManager following REQ-10, REQ-11 in Update() method, on WindowSizeMsg or panel visibility change call layoutManager.Calculate() with current width, height, and getPanelStates() (returns PanelStates struct with visibility flags), check layout.IsTooSmall and return early if true, propagate dimensions to all panels using layout.Sidebar/Main/Metadata dimensions with each panel's SetSize() method, update existing views (listView, detailView, forms) with layout.Main dimensions | Restrictions: Must call SetSize() for all visible panels, handle IsTooSmall case, don't resize hidden panels | _Leverage: Existing WindowSizeMsg handling, SetSize() pattern | Success: Layout recalculates correctly on resize, panel dimensions update smoothly, minimum size detection works, all panels receive correct dimensions, existing views fit within main panel area_

- [x] 13. Implement multi-panel View rendering
  - File: cmd/tui/view.go or cmd/tui/model.go View() method
  - Render multi-panel layout with lipgloss horizontal joins
  - Purpose: Display sidebar, main content, and metadata panels side-by-side
  - _Leverage: lipgloss.JoinHorizontal(), lipgloss.JoinVertical(), existing View() logic_
  - _Requirements: REQ-1, REQ-8_
  - _Prompt: Role: Go TUI Developer specializing in layout rendering and Lipgloss | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Refactor View() method following REQ-1, REQ-8, check layout.IsTooSmall first and render warning if true, collect visible panels in slice (sidebar if visible, main content always, metadata if visible), render each with getPanelBorderStyle() using focused/unfocused styles, apply Width() and Height() from layout dimensions, use lipgloss.JoinHorizontal(lipgloss.Top, panels...) for main row, append process panel vertically if visible, append command bar if open, append status bar at bottom, maintain existing overlay logic (help, confirmations render over panels) | Restrictions: Must use layout dimensions from LayoutManager, respect panel visibility flags, preserve overlay behavior | _Leverage: Existing View() logic, lipgloss joining, getPanelBorderStyle() helper | Success: Multi-panel layout renders correctly, panels displayed side-by-side with proper borders, focused panel highlighted, overlays work over panels, status bar always at bottom, layout adapts to panel visibility_

- [x] 14. Update status bar with panel indicators
  - File: cmd/tui/components/statusbar.go
  - Add panel toggle indicators to shortcuts display
  - Purpose: Show available panel toggle keys in status bar
  - _Leverage: Existing StatusBar.SetShortcuts() method_
  - _Requirements: REQ-2_
  - _Prompt: Role: Go TUI Developer specializing in status indicators | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Update StatusBar following REQ-2, modify updateStatusBar() in model.go to include panel shortcuts in status bar, add "s: sidebar | m: metadata | p: processes" to shortcut strings for relevant states, update shortcuts when panel visibility changes, maintain existing view-specific shortcuts (List: "/: search | a: add", Detail: "m: toggle | c: copy | e: edit") | Restrictions: Must not remove existing shortcuts, keep shortcut strings concise, update dynamically based on state | _Leverage: Existing updateStatusBar() function, StatusBar.SetShortcuts() | Success: Status bar displays panel toggle shortcuts, shortcuts update correctly per view, all existing shortcuts preserved, panel indicators show current state_

- [x] 15. Implement sidebar category selection integration
  - File: cmd/tui/model.go (continue)
  - Connect sidebar credential selection to main content display
  - Purpose: When credential selected in sidebar, show in main area and metadata panel
  - _Leverage: Sidebar.GetSelectedCredential(), existing credential loading_
  - _Requirements: REQ-3, REQ-6_
  - _Prompt: Role: Go TUI Developer with component integration expertise | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Integrate sidebar selection following REQ-3, REQ-6, when sidebar has focus and user presses Enter on credential, get selected credential using sidebar.GetSelectedCredential(), load full credential details using loadCredentialDetailsCmd(), display in detail view in main content area, update metadata panel with SetCredential(), update breadcrumb with category path, trigger layout update if needed | Restrictions: Only handle selection when sidebar focused, must load full credential (not just metadata), update all relevant panels | _Leverage: Existing loadCredentialDetailsCmd, state transitions to StateDetail | Success: Selecting credential in sidebar loads details in main area, metadata panel updates simultaneously, breadcrumb shows path, focus can switch to main area, existing credential selection from list still works_

- [x] 16. Add unit tests for Layout Manager
  - File: cmd/tui/components/layout_manager_test.go
  - Write comprehensive tests for dimension calculations
  - Purpose: Ensure layout calculations work correctly across all scenarios
  - _Leverage: Go testing package, table-driven test pattern_
  - _Requirements: REQ-10, REQ-11_
  - _Prompt: Role: QA Engineer with Go testing expertise | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create unit tests for LayoutManager following REQ-10, REQ-11, test full layout (140x40, all panels visible), medium layout (100x30, metadata hidden), small layout (70x25, sidebar and metadata hidden), too small (50x15, warning shown), test minimum constraints enforcement (sidebar min 20, main min 40, metadata min 25), test multi-panel main area division, test panel visibility combinations, use table-driven tests for breakpoint scenarios | Restrictions: Test calculations only, no UI rendering, cover all edge cases | _Leverage: Existing test patterns from model_test.go | Success: All layout calculation scenarios tested, edge cases covered (negative dimensions, zero panels, maximum panels), breakpoint logic verified, minimum constraints enforced correctly, tests pass reliably_

- [x] 17. Add unit tests for Category Tree
  - File: cmd/tui/components/category_tree_test.go
  - Write tests for credential categorization logic
  - Purpose: Verify category matching patterns work correctly
  - _Leverage: Go testing package, test credential fixtures_
  - _Requirements: REQ-3_
  - _Prompt: Role: QA Engineer specializing in data transformation testing | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create unit tests for category_tree.go following REQ-3, test CategorizeCredentials() with various service names (aws-prod → Cloud, github → Git, postgres → Databases, stripe-api → Payment, openai → AI, random-service → Uncategorized), test case-insensitivity, test empty input, test all credentials categorized (none lost), test icon mappings for all categories, verify category counts match actual credentials | Restrictions: Pure function testing, no state or side effects, comprehensive pattern coverage | _Leverage: vault.CredentialMetadata, test helpers | Success: All categorization patterns tested and working, edge cases handled (empty, unknown services), icon mappings verified, count calculations correct, all tests pass_

- [x] 18. Add unit tests for Sidebar Panel
  - File: cmd/tui/components/sidebar_test.go
  - Write tests for sidebar navigation and rendering
  - Purpose: Verify sidebar interaction and state management
  - _Leverage: Go testing package, mock categories_
  - _Requirements: REQ-3, REQ-4, REQ-14_
  - _Prompt: Role: QA Engineer with UI component testing expertise | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create unit tests for SidebarPanel following REQ-3, REQ-4, REQ-14, test category navigation (j/k/arrows move selection), test expansion/collapse (Enter/l toggles, h collapses), test credential selection within categories, test stats display updates, test viewport scrolling for long lists, test SetSize() updates dimensions correctly, test GetSelectedCredential() returns correct item, test empty categories gracefully handled | Restrictions: Mock viewport, test state management not rendering, cover all navigation paths | _Leverage: Existing view test patterns | Success: Navigation tests pass (up/down/expand/collapse), selection state maintained correctly, stats calculations verified, scrolling behavior tested, all interaction scenarios covered_

- [x] 19. Add unit tests for Command Bar
  - File: cmd/tui/components/command_bar_test.go
  - Write tests for command parsing and execution
  - Purpose: Verify command bar correctly parses and handles commands
  - _Leverage: Go testing package, table-driven tests_
  - _Requirements: REQ-7_
  - _Prompt: Role: QA Engineer specializing in parser testing | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create unit tests for CommandBar following REQ-7, test command parsing (:add github → Command{Name: "add", Args: ["github"]}, :search query → Command{Name: "search", Args: ["query"]}, :quit → Command{Name: "quit", Args: []}), test invalid commands, test empty input, test command history (up/down navigation), test error display, test Focus()/Blur() state changes, use table-driven tests for command parsing scenarios | Restrictions: Test parsing logic, mock textinput for interaction tests, cover all command types | _Leverage: parseCommand() function, test patterns | Success: All command parsing scenarios tested, history navigation verified, error handling tested, Focus/Blur state changes verified, invalid commands handled gracefully, all tests pass_

- [x] 20. Add integration tests for panel system
  - File: test/tui_dashboard_integration_test.go
  - Write end-to-end tests for dashboard functionality
  - Purpose: Verify complete panel system works with real vault and interactions
  - _Leverage: Existing integration test patterns, build tag integration_
  - _Requirements: All dashboard requirements_
  - _Prompt: Role: Integration Test Engineer with Go and Bubble Tea expertise | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Create integration tests covering all dashboard requirements, test panel toggle functionality (s/m/p/f keys toggle visibility correctly), test panel focus switching (Tab cycles through visible panels), test responsive layout (140x40 shows all, 100x30 hides metadata, 70x25 shows main only), test command bar execution (: opens, commands execute, Esc cancels), test sidebar navigation (category expansion, credential selection), test minimum size warning (<60x20 shows warning), test metadata panel updates (displays selected credential details), use real VaultService with temp vaults | Restrictions: Use build tag integration, create/cleanup temp vaults, test with real components (no mocks except keychain if unavailable) | _Leverage: Existing integration test setup, createTestVault() helper | Success: Panel toggles work correctly, focus switching functions, responsive breakpoints trigger correctly, command bar executes commands, sidebar integration works, all panels interact properly, tests pass reliably_

- [ ] 21. Manual testing and polish
  - No specific file
  - Perform cross-platform testing and visual polish
  - Purpose: Ensure dashboard works and looks good on all supported terminals
  - _Requirements: All_
  - _Prompt: Role: QA Engineer with cross-platform terminal testing expertise | Task: Implement spec tui-dashboard-layout, first run spec-workflow-guide to get the workflow guide then implement the task: Perform comprehensive manual testing covering all requirements, test on Windows (Windows Terminal, PowerShell, CMD), macOS (Terminal.app, iTerm2), Linux (GNOME Terminal, Konsole, Alacritty), verify all panel interactions work smoothly, verify responsive layouts at different terminal sizes, verify icons display or fallback gracefully, verify border rendering and colors, verify all keyboard shortcuts work without conflicts, verify performance (smooth rendering, no lag), verify existing TUI features still work (add/edit/delete/search), document any terminal-specific issues or quirks | Restrictions: Test on actual terminals (not emulators), verify full user workflows, ensure professional appearance | Success: Dashboard works correctly on all target platforms, icons display or fallback appropriately, performance is acceptable (<150ms startup, <16ms transitions), all existing features functional, visual polish complete, no critical bugs found_
