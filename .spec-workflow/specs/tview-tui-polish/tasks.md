# Tasks Document

## Phase 1: Data Model Extension (Foundation)

- [ ] 1. Add Category and URL fields to Credential struct
  - File: internal/vault/vault.go
  - Add `Category string` and `URL string` fields to Credential struct with JSON tags
  - Purpose: Extend data model to support credential organization and web interface links
  - _Leverage: Existing Credential struct (lines 37-45), JSON serialization patterns_
  - _Requirements: 1.1, 1.2_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go backend developer specializing in data models and JSON serialization | Task: Add Category and URL string fields to the Credential struct in internal/vault/vault.go (around lines 37-45), with proper json tags following the existing pattern. Category is for grouping credentials (empty string = "Uncategorized"), URL is for optional web interface links. Both fields are optional. | Restrictions: Do not modify existing fields, maintain JSON tag format consistency, do not change serialization behavior, ensure backward compatibility with existing vaults (new fields will deserialize as empty strings) | Success: Credential struct has Category and URL fields with correct json tags, existing vault files can load without errors, new fields serialize/deserialize correctly. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 2. Update VaultService.AddCredential method signature
  - File: internal/vault/vault.go
  - Extend AddCredential method to accept category, url, notes parameters and populate Credential struct
  - Purpose: Enable adding credentials with full metadata through the vault service
  - _Leverage: Existing AddCredential method implementation, Credential struct validation patterns_
  - _Requirements: 1.1, 1.3_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go backend developer with expertise in API design and service layers | Task: Update the AddCredential method signature in internal/vault/vault.go to accept (service, username, password, category, url, notes string) parameters. Update the method implementation to populate the new Category and URL fields when creating the Credential struct. Maintain existing validation for service/username/password (required), but allow category/url/notes to be empty strings. | Restrictions: Do not break existing callers (this will require updating CLI commands and TUI in subsequent tasks), maintain all existing validation logic, ensure CreatedAt/UpdatedAt timestamps are still set correctly, preserve usage tracking behavior | Success: AddCredential accepts 6 parameters, creates Credential with all fields populated, validation still works for required fields, method compiles successfully. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 3. Update VaultService.UpdateCredential method signature
  - File: internal/vault/vault.go
  - Extend UpdateCredential method to accept category, url, notes parameters and update Credential struct
  - Purpose: Enable updating credentials with full metadata through the vault service
  - _Leverage: Existing UpdateCredential method implementation, credential lookup patterns_
  - _Requirements: 1.2, 1.3_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go backend developer with expertise in CRUD operations and data updates | Task: Update the UpdateCredential method signature in internal/vault/vault.go to accept (service, username, password, category, url, notes string) parameters. Update the method implementation to update the Category and URL fields. If password is empty string, preserve the existing password (user didn't change it). Update UpdatedAt timestamp. | Restrictions: Do not break credential lookup logic, maintain service name as primary key, preserve existing password if new password is empty, ensure UpdatedAt is set correctly, do not modify CreatedAt or usage tracking | Success: UpdateCredential accepts 6 parameters, updates all modifiable fields correctly, empty password preserves existing password, UpdatedAt timestamp updates, method compiles successfully. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 4. Update CLI add command to use extended signature
  - File: cmd/add.go
  - Update add command to pass empty strings for category, url when calling VaultService.AddCredential
  - Purpose: Maintain CLI compatibility with extended vault service signature
  - _Leverage: Existing addCmd.RunE function, flag parsing patterns_
  - _Requirements: 1.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go CLI developer with expertise in Cobra command frameworks | Task: Update the runAdd function in cmd/add.go to call VaultService.AddCredential with 6 parameters: service, username, password, "" (empty category), "" (empty url), notes (existing notes variable). The CLI doesn't support category/url yet, so pass empty strings as placeholders. | Restrictions: Do not add new CLI flags (category/url support is future enhancement), maintain existing flag behavior, preserve all existing validation and error handling, ensure backward compatibility for CLI users | Success: CLI add command compiles and works correctly, passes empty strings for category/url, existing functionality unchanged, notes field still works. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 5. Update CLI update command to use extended signature
  - File: cmd/update.go
  - Update update command to pass empty strings for category, url when calling VaultService.UpdateCredential
  - Purpose: Maintain CLI compatibility with extended vault service signature
  - _Leverage: Existing updateCmd.RunE function, existing credential fetching pattern_
  - _Requirements: 1.2_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go CLI developer with expertise in Cobra frameworks and CRUD operations | Task: Update the runUpdate function in cmd/update.go to call VaultService.UpdateCredential with 6 parameters. Fetch the existing credential to preserve category/url (since CLI doesn't modify them), then call UpdateCredential(service, username, password, existingCred.Category, existingCred.URL, notes). | Restrictions: Do not add new CLI flags, preserve existing category/url values from the credential, maintain all existing validation and error handling, ensure password preservation logic still works | Success: CLI update command compiles and works correctly, preserves existing category/url, all existing functionality works, notes update still functions. After completion, update tasks.md by changing this task from [ ] to [x]._

## Phase 2: TUI Integration

- [ ] 6. Update AppState.AddCredential method signature
  - File: cmd/tui-tview/models/state.go
  - Extend AppState.AddCredential to accept category, url, notes and pass to VaultService
  - Purpose: Enable TUI forms to pass all credential fields through AppState
  - _Leverage: Existing AppState.AddCredential implementation (pass-through pattern to VaultService)_
  - _Requirements: 1.1, 2.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go developer with expertise in application state management and service integration | Task: Update the AddCredential method signature in cmd/tui-tview/models/state.go to accept (service, username, password, category, url, notes string) parameters. Pass all 6 parameters directly to as.vaultService.AddCredential(). Trigger credentials refresh callback after successful add. | Restrictions: Maintain existing error handling pattern, preserve onCredentialsChanged callback invocation, do not add business logic (AppState is pass-through layer), ensure thread-safe access to vaultService | Success: AppState.AddCredential accepts 6 parameters, passes them to VaultService, triggers UI refresh on success, compiles successfully. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 7. Update AppState.UpdateCredential method signature
  - File: cmd/tui-tview/models/state.go
  - Extend AppState.UpdateCredential to accept category, url, notes and pass to VaultService
  - Purpose: Enable TUI forms to update all credential fields through AppState
  - _Leverage: Existing AppState.UpdateCredential implementation (pass-through pattern)_
  - _Requirements: 1.2, 2.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go developer with expertise in application state management and CRUD operations | Task: Update the UpdateCredential method signature in cmd/tui-tview/models/state.go to accept (service, username, password, category, url, notes string) parameters. Pass all 6 parameters directly to as.vaultService.UpdateCredential(). Trigger credentials refresh and selection update callbacks after successful update. | Restrictions: Maintain existing error handling, preserve callback invocations (onCredentialsChanged, onSelectionChanged), do not add business logic, ensure proper refresh behavior | Success: AppState.UpdateCredential accepts 6 parameters, passes to VaultService, triggers UI updates on success, compiles successfully. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 8. Connect AddForm to pass all 6 fields to AppState
  - File: cmd/tui-tview/components/forms.go
  - Update AddForm.onAddPressed to extract category, url, notes from form and pass to AppState.AddCredential
  - Purpose: Complete add form integration with extended data model
  - _Leverage: Existing form field extraction pattern (GetFormItem), existing field indices (service=0, username=1, password=2, category=3, url=4, notes=5)_
  - _Requirements: 1.1, 2.1, 5.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go frontend developer with expertise in tview forms and type assertions | Task: Update onAddPressed method in AddForm (cmd/tui-tview/components/forms.go around line 72) to extract all 6 fields. Service/username/password already extracted. Add: category from dropdown GetFormItem(3).(*tview.DropDown).GetCurrentOption() (use second return value for selected text), url from GetFormItem(4).(*tview.InputField).GetText(), notes from GetFormItem(5).(*tview.TextArea).GetText(). Pass all 6 to AppState.AddCredential(). Remove TODO comment on line 84-85. | Restrictions: Maintain existing validation call, preserve error handling pattern, ensure form closes on success, do not change form field order or structure | Success: AddForm extracts and passes all 6 fields, category dropdown text extracted correctly, notes textarea extracted, TODO comment removed, form works end-to-end. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 9. Connect EditForm to pass all 6 fields to AppState
  - File: cmd/tui-tview/components/forms.go
  - Update EditForm.onSavePressed to extract category, url, notes and pass to AppState.UpdateCredential
  - Purpose: Complete edit form integration with extended data model
  - _Leverage: Existing field extraction in onSavePressed, existing password preservation logic (lines 226-234)_
  - _Requirements: 1.2, 2.1, 2.5_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go frontend developer with expertise in form handling and credential updates | Task: Update onSavePressed method in EditForm (cmd/tui-tview/components/forms.go around line 213) to extract all 6 fields. Service/username/password already extracted. Add: category from GetFormItem(3).(*tview.DropDown).GetCurrentOption() (second return value), url from GetFormItem(4).(*tview.InputField).GetText(), notes from GetFormItem(5).(*tview.TextArea).GetText(). Preserve existing password-empty-check logic (lines 226-234). Pass all 6 to AppState.UpdateCredential(). | Restrictions: Maintain password preservation logic, preserve error handling, ensure form closes on success, do not change validation behavior | Success: EditForm extracts and passes all 6 fields, password preservation still works, empty fields handled correctly, form works end-to-end. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 10. Implement EditForm pre-population for all fields
  - File: cmd/tui-tview/components/forms.go
  - Update buildFormFieldsWithValues to fetch full credential and pre-populate password, category, url, notes fields
  - Purpose: Show current values when editing credentials
  - _Leverage: GetFullCredential pattern, existing field initialization in buildFormFieldsWithValues (lines 188-209)_
  - _Requirements: 2.1, 2.2, 2.3, 2.4_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go frontend developer with expertise in form pre-population and data fetching | Task: Update buildFormFieldsWithValues in EditForm (cmd/tui-tview/components/forms.go around line 187) to pre-populate all fields. Call ef.appState.GetFullCredential(ef.credential.Service) at start. Use fullCred to: 1) Pre-populate password field (line 198): AddPasswordField("Password", fullCred.Password, 40, '*', nil) 2) Pre-populate url field (line 203): AddInputField("URL", fullCred.URL, 50, nil, nil) 3) Pre-populate notes field (line 204): AddTextArea("Notes", fullCred.Notes, 50, 5, 0, nil) 4) For category dropdown, find index of fullCred.Category in categories slice, use it as selectedOption parameter. Remove TODO comments on lines 201, 291. Handle GetFullCredential error by showing in status bar via ef.appState.onError callback. | Restrictions: Handle GetFullCredential errors gracefully, do not break form if credential fetch fails, maintain field order, preserve existing form structure, ensure dropdown index is valid (default to 0 if category not found) | Success: All fields pre-populate with current values, password shows masked, category dropdown selects correct category, url and notes display, TODO comments removed, errors handled gracefully. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 11. Update AppState.GetCategories to extract from credentials
  - File: cmd/tui-tview/models/state.go
  - Implement GetCategories to extract unique categories from credentials, including "Uncategorized" for empty categories
  - Purpose: Populate category dropdown with actual categories from vault
  - _Leverage: Existing GetCategories method stub (line 364 with TODO comment), credential iteration patterns in LoadCredentials_
  - _Requirements: 1.4, 5.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go developer with expertise in data aggregation and slice operations | Task: Implement the GetCategories method in cmd/tui-tview/models/state.go (line 364) to iterate through as.credentials, extract unique category values into a slice. If a credential has empty Category, it's implicitly "Uncategorized". Always include "Uncategorized" in the returned slice. Sort categories alphabetically with "Uncategorized" first, then sorted categories. Remove TODO comment. | Restrictions: Handle nil credentials slice safely, ensure unique categories only, maintain alphabetical sort (except Uncategorized first), do not modify as.credentials, ensure thread-safe access to credentials | Success: GetCategories returns unique category list, "Uncategorized" always present and first, remaining categories sorted alphabetically, TODO comment removed, method works correctly. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 12. Update Sidebar to group credentials by category
  - File: cmd/tui-tview/components/sidebar.go
  - Modify sidebar tree structure to group credentials by category, with "Uncategorized" for empty categories
  - Purpose: Organize credentials visually by category for better navigation
  - _Leverage: Existing sidebar treeView structure, GetCategories method from AppState, existing Refresh method_
  - _Requirements: 1.4, 1.5_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go UI developer with expertise in tview TreeView components and hierarchical data structures | Task: Update the Sidebar.Refresh method in cmd/tui-tview/components/sidebar.go to build a category-grouped tree. Get categories from sb.appState.GetCategories(). For each category, create a category tree node, then add credential nodes under it. Credentials with empty Category go under "Uncategorized". Maintain "ALL Credentials" root node, then add category nodes as children. Preserve existing node selection behavior and keyboard navigation. | Restrictions: Do not break existing tree navigation, maintain node selection state if possible, preserve keyboard shortcuts, ensure tree updates efficiently on Refresh, handle empty credentials gracefully | Success: Sidebar shows category groups, credentials organized under categories, "Uncategorized" shows credentials with empty category, tree navigation works, selection preserved on refresh. After completion, update tasks.md by changing this task from [ ] to [x]._

## Phase 3: UI Polish

- [ ] 13. Add Category and URL display to DetailView
  - File: cmd/tui-tview/components/detail.go
  - Add Category and URL fields to credential detail display in formatCredential method
  - Purpose: Show all credential metadata in detail panel
  - _Leverage: Existing formatCredential method (lines 60-109), existing field formatting pattern_
  - _Requirements: 1.3_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go UI developer with expertise in text formatting and tview color tags | Task: Update the formatCredential method in cmd/tui-tview/components/detail.go (around line 60) to display Category and URL fields. After the Username field (around line 70), add: if cred.Category != "" { b.WriteString(fmt.Sprintf("[gray]Category:[-]  [white]%s[-]\n", cred.Category)) }. After Category, add: if cred.URL != "" { b.WriteString(fmt.Sprintf("[gray]URL:[-]       [white]%s[-]\n", cred.URL)) }. Maintain field alignment with existing Service/Username/Password fields. | Restrictions: Only display if fields are non-empty, maintain color tag consistency ([gray] for labels, [white] for values), preserve existing field spacing and alignment, do not modify password field behavior | Success: Category displays if present, URL displays if present, fields align with existing layout, empty fields don't show, formatting matches existing pattern. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 14. Improve status bar keyboard shortcut formatting
  - File: cmd/tui-tview/components/statusbar.go
  - Update getShortcutsForContext to use white color for keys, gray for action descriptions
  - Purpose: Improve shortcut scannability and visual hierarchy
  - _Leverage: Existing getShortcutsForContext method (lines 96-113), tview color tag system_
  - _Requirements: 3.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go UI developer with expertise in terminal colors and readability | Task: Update getShortcutsForContext method in cmd/tui-tview/components/statusbar.go (lines 96-113) to improve key-action formatting. Change from "[gray][Tab] Next..." to "[white]Tab[-][gray] Next  [white]↑↓[-][gray] Navigate..." pattern. Apply to all four context cases (FocusSidebar, FocusTable, FocusDetail, FocusModal). Use [white] for keys, [-] to reset, [gray] for action text, ensure consistent spacing. | Restrictions: Maintain existing context logic, do not change which shortcuts appear per context, preserve spacing between shortcut groups, ensure readability on dark terminals | Success: Keys display in white, actions in gray, all four contexts updated consistently, shortcuts remain readable and well-spaced. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 15. Improve help modal formatting and alignment
  - File: cmd/tui-tview/events/handlers.go
  - Reformat help modal text with header separator, aligned columns, and consistent color scheme
  - Purpose: Improve help modal readability and professional appearance
  - _Leverage: Existing showHelpModal method (lines 216-246), tview modal component, color tag system_
  - _Requirements: 3.3, 3.4_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go UI developer with expertise in modal formatting and text alignment | Task: Update the helpText string in showHelpModal method (cmd/tui-tview/events/handlers.go lines 216-236) to improve formatting. Add header separator: "[yellow]═══════════════════════════════════════════════[-]\n[yellow]           Keyboard Shortcuts[-]\n[yellow]═══════════════════════════════════════════════[-]". For each shortcut, use format: "  [white]Key[-]              Description". Left-align keys (15 char width), descriptions start at column 20. Use [white] for keys, [gray] for section headers ([cyan] for section titles like "Navigation"). Add spacing between sections. Ensure all shortcuts from existing help are included. | Restrictions: Preserve all existing shortcuts, do not change modal width/height parameters (60x25), maintain tview color tag syntax, ensure readability, do not exceed modal width | Success: Help modal displays header separator, shortcuts aligned in columns, keys in white, section headers in cyan/gray, all shortcuts present, visually balanced and readable. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 16. Fix form input field background sizing issues
  - File: cmd/tui-tview/styles/theme.go or cmd/tui-tview/styles/forms.go
  - Investigate and fix input field background rendering issues (top 3 fields don't fill space, URL extends)
  - Purpose: Ensure consistent form field appearance
  - _Leverage: Existing ApplyFormStyle function, tview.Form field styling API, theme color system_
  - _Requirements: 4.1, 4.2_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go UI developer with expertise in tview form styling and terminal rendering | Task: Investigate form field background issues in ApplyFormStyle (cmd/tui-tview/styles/theme.go or forms.go). Issues: 1) InputField backgrounds (service, username, password) don't fill allocated height 2) URL InputField background extends past boundaries. Solutions: Ensure SetFieldBackgroundColor is set correctly, investigate if field width needs explicit setting, check if tview version has known background rendering issues. May need to set individual field styles instead of global form style. Test with actual form instances (AddForm, EditForm). | Restrictions: Maintain existing theme color scheme, do not break existing form functionality, ensure changes work across Windows Terminal, iTerm2, gnome-terminal, preserve form layout and spacing | Success: All input fields have consistent background rendering, backgrounds fill allocated space, no overflow on URL field, forms look polished and professional. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 17. Ensure Save and Cancel buttons visible in forms
  - File: cmd/tui-tview/components/forms.go or event handlers that show modals
  - Increase modal height or adjust form layout to ensure Save/Cancel buttons fully visible
  - Purpose: Ensure users can see and access form action buttons
  - _Leverage: Existing modal display code in events/handlers.go, form button setup in forms.go (AddButton calls)_
  - _Requirements: 4.3_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go UI developer with expertise in modal layouts and tview flex containers | Task: Find where AddForm and EditForm modals are displayed (likely in cmd/tui-tview/events/handlers.go or forms.go). Check modal height parameters. Forms have 6 fields + 2 buttons = minimum ~20 lines needed. Increase modal height if needed (e.g., from 25 to 30). Alternatively, adjust form field heights if textarea is too tall. Test that Save and Cancel buttons are fully visible and clickable in both add and edit modals. | Restrictions: Do not reduce field sizes excessively, maintain form usability, ensure buttons remain at bottom, preserve existing button alignment (AlignRight), test on standard terminal sizes (80x24, 120x40) | Success: Save and Cancel buttons fully visible in add form, fully visible in edit form, buttons accessible with Tab/Enter navigation, forms fit in modal without scrolling. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 18. Improve table row selection visibility
  - File: cmd/tui-tview/components/table.go or cmd/tui-tview/styles/theme.go
  - Enhance selected row highlighting with clearer visual indication (background color, border, etc.)
  - Purpose: Make it obvious which credential is currently selected
  - _Leverage: Existing CredentialTable component, tview.Table selection styling, theme color system_
  - _Requirements: 4.4_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go UI developer with expertise in tview tables and selection styling | Task: Investigate table row selection styling in cmd/tui-tview/components/table.go or styles/theme.go. Current issue: selected row not clearly visible. Solutions: 1) Use table.SetSelectedStyle() to set bold, background color, or different text color 2) Consider using theme.SelectedBg or theme.ActiveBorderColor for selected row 3) Test with arrow key navigation and mouse selection. Ensure selected row stands out clearly from unselected rows. | Restrictions: Maintain existing table structure and data, preserve keyboard navigation behavior, ensure colors work on dark and light terminal themes, do not break mouse selection | Success: Selected table row has clear visual indication (background color or bold text), easy to see which row is selected, selection visible during arrow key navigation, works on multiple terminal emulators. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 19. Ensure DetailView syncs with table/sidebar selection
  - File: cmd/tui-tview/components/detail.go, cmd/tui-tview/components/table.go, cmd/tui-tview/components/sidebar.go
  - Verify and fix selection change callbacks to trigger DetailView.Refresh immediately
  - Purpose: Keep detail panel in sync with selected credential
  - _Leverage: Existing AppState.onSelectionChanged callback, DetailView.Refresh method, table/sidebar selection handlers_
  - _Requirements: 4.5_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go developer with expertise in event-driven UI and callback patterns | Task: Investigate DetailView sync issue. Check: 1) Table selection handler - does it call AppState.SetSelectedCredential()? 2) Sidebar selection handler - same check 3) AppState.SetSelectedCredential - does it trigger onSelectionChanged callback? 4) DetailView registered for onSelectionChanged in main.go or tui.go? Ensure selection change → SetSelectedCredential → callback → DetailView.Refresh chain works. Add debug logging if needed to trace flow. | Restrictions: Do not add unnecessary refreshes (performance), maintain callback pattern architecture, preserve existing selection logic, ensure both keyboard and mouse selection trigger updates | Success: DetailView updates immediately when table row selected, DetailView updates when sidebar node selected, selection changes via keyboard update detail, selection via mouse updates detail. After completion, update tasks.md by changing this task from [ ] to [x]._

## Phase 4: Verification and Cleanup

- [ ] 20. Verify search UI exists and functions
  - File: cmd/tui-tview/events/handlers.go or cmd/tui-tview/components/
  - Test that pressing "/" displays search UI and filters credentials
  - Purpose: Confirm search feature is implemented and working
  - _Leverage: Global keyboard handler in events/handlers.go, any existing search components_
  - _Requirements: 6.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA engineer with expertise in feature verification and TUI testing | Task: Verify search UI implementation. 1) Launch TUI and press "/" key 2) Check if search modal/input appears 3) Type search text and verify credentials filter 4) Test Esc to close search 5) If search doesn't exist: search for "/" key handler in cmd/tui-tview/events/handlers.go - may be unimplemented or commented out. If unimplemented, create GitHub issue or note in spec. If implemented, document how it works. | Restrictions: This is verification only - do not implement search if missing (separate spec needed), document findings clearly, test with multiple search terms | Success: Search UI verified working OR documented as not implemented, if working: search filters credentials correctly, Esc closes search, search results update dynamically. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 21. Test delete confirmation modal (already verified working)
  - File: Manual testing only
  - Confirm delete confirmation displays and deletion works
  - Purpose: Verify delete feature confirmed working by user
  - _Leverage: N/A - manual testing_
  - _Requirements: 6.2 (verified)_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA engineer performing regression testing | Task: Perform manual test of delete functionality. 1) Launch TUI 2) Select a credential 3) Press "d" key 4) Verify modal asks for confirmation 5) Test both confirm and cancel paths 6) Verify credential deleted from vault when confirmed 7) Verify credential remains when cancelled. Document test results. User has already confirmed this works - this task is for thoroughness. | Restrictions: This is verification only, do not modify code, test with test vault (test-vault-tview/vault.enc), do not delete real credentials | Success: Delete confirmation modal displays correctly, deletion works when confirmed, cancel preserves credential, test documented as passing. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 22. Test copy password feedback (already verified working)
  - File: Manual testing only
  - Confirm copy to clipboard works and status bar shows feedback
  - Purpose: Verify copy feature confirmed working by user
  - _Leverage: N/A - manual testing_
  - _Requirements: 6.3, 6.4 (verified)_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA engineer performing integration testing | Task: Perform manual test of copy password functionality. 1) Launch TUI 2) Select a credential 3) Press "c" key 4) Verify status bar shows "Password copied" or similar message 5) Paste into text editor to verify clipboard contains correct password 6) Test with multiple credentials to ensure consistency. Document test results. User has already confirmed this works. | Restrictions: This is verification only, do not modify code, test with test vault credentials, ensure clipboard is cleared after test | Success: Copy shows status bar feedback, clipboard contains correct password, paste works in external application, test documented as passing. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 23. Resolve TODO comment in forms.go line 85
  - File: cmd/tui-tview/components/forms.go
  - Remove or update TODO comment about passing Category/URL/Notes fields
  - Purpose: Clean up code by removing obsolete TODO comments
  - _Leverage: Task 8 implementation (Category/URL/Notes now passed)_
  - _Requirements: 5.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Developer performing code cleanup | Task: Review TODO comment at line 85 in cmd/tui-tview/components/forms.go. After task 8 is complete, Category/URL/Notes are now passed to AppState.AddCredential(). Remove the TODO comment entirely since the feature is now implemented. Verify surrounding code is clean and well-commented. | Restrictions: Do not modify functionality, only remove TODO comment, ensure code is still readable without the comment | Success: TODO comment removed, code compiles, functionality unchanged, no obsolete comments remain. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 24. Resolve TODO comment in forms.go line 201
  - File: cmd/tui-tview/components/forms.go
  - Remove or update TODO comment about pre-populating Category/URL/Notes
  - Purpose: Clean up code by removing obsolete TODO comments
  - _Leverage: Task 10 implementation (pre-population now implemented)_
  - _Requirements: 5.2_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Developer performing code cleanup | Task: Review TODO comment at line 201 in cmd/tui-tview/components/forms.go. After task 10 is complete, Category/URL/Notes are now pre-populated from the credential. Remove the TODO comment entirely. Verify EditForm pre-population is working correctly. | Restrictions: Do not modify functionality, only remove TODO comment, ensure code remains clear | Success: TODO comment removed, pre-population works, code compiles, no obsolete comments. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 25. Resolve TODO comment in forms.go line 291
  - File: cmd/tui-tview/components/forms.go
  - Remove or update TODO comment about matching credential category field
  - Purpose: Clean up code by removing obsolete TODO comments
  - _Leverage: Task 10 implementation (category matching now implemented)_
  - _Requirements: 5.3_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Developer performing code cleanup | Task: Review TODO comment at line 291 in cmd/tui-tview/components/forms.go. After task 10 is complete, category dropdown now correctly selects the credential's category. Remove the TODO comment entirely. Verify category selection works correctly in edit form. | Restrictions: Do not modify functionality, only remove TODO comment, ensure category matching logic is clear | Success: TODO comment removed, category selection works, code compiles, no obsolete comments. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 26. Resolve TODO comment in state.go line 364
  - File: cmd/tui-tview/models/state.go
  - Remove or update TODO comment about category extraction
  - Purpose: Clean up code by removing obsolete TODO comments
  - _Leverage: Task 11 implementation (GetCategories now implemented)_
  - _Requirements: 5.4_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Developer performing code cleanup | Task: Review TODO comment at line 364 in cmd/tui-tview/models/state.go. After task 11 is complete, GetCategories method is now fully implemented with proper category extraction. Remove the TODO comment entirely. Verify GetCategories returns correct unique categories. | Restrictions: Do not modify GetCategories implementation, only remove TODO comment, ensure method documentation is clear | Success: TODO comment removed, GetCategories works correctly, code compiles, no obsolete comments remain. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 27. Run full integration test suite
  - File: Manual testing across all components
  - Test complete workflow: add credential with category/url, edit it, view details, delete
  - Purpose: Verify all components work together correctly
  - _Leverage: Test vault (test-vault-tview), all implemented features_
  - _Requirements: All requirements_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA engineer performing comprehensive integration testing | Task: Execute full end-to-end test: 1) Launch TUI with test vault 2) Press "n" to add credential with all fields (service, username, password, category, url, notes) 3) Verify credential appears in sidebar under category 4) Select credential, verify all fields display in detail view 5) Press "e" to edit, verify all fields pre-populated 6) Change category, save, verify sidebar updates 7) Test password visibility toggle 8) Test copy password 9) Test delete with confirmation 10) Verify help modal formatting 11) Test keyboard shortcuts across all contexts. Document any issues found. | Restrictions: Use test vault only, document all test steps and results, note any bugs or UX issues, test on primary development terminal emulator | Success: All features work in integration, no critical bugs found, UX is polished and intuitive, spec requirements met, test results documented. After completion, update tasks.md by changing this task from [ ] to [x]._

- [ ] 28. Update test data setup script for new fields
  - File: test/setup-tview-test-data.bat
  - Update script to add credentials with category and url using --notes flag workaround or direct vault manipulation
  - Purpose: Provide realistic test data with categories for testing
  - _Leverage: Existing test setup script, CLI add command (now supports extended signature)_
  - _Requirements: 1.1_
  - _Prompt: Implement the task for spec tview-tui-polish, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps engineer with expertise in test data generation and bash scripting | Task: Update test/setup-tview-test-data.bat to create credentials with categories. Since CLI doesn't expose category/url flags yet, credentials added via CLI will have empty category. For testing, either: 1) Wait for CLI to add category flags (separate enhancement) OR 2) Directly edit test vault JSON after creation to add category/url fields OR 3) Document that test data script creates credentials without categories for now. Choose option 3 for minimal change - add comment noting categories can be added via TUI after vault is created. | Restrictions: Do not break existing test data script, maintain 15 test credentials, ensure script still runs successfully, preserve backward compatibility | Success: Test script runs without errors, creates 15 credentials, script documented regarding category support, credentials can be edited in TUI to add categories. After completion, update tasks.md by changing this task from [ ] to [x]._
