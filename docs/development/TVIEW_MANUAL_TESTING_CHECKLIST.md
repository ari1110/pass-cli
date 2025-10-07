# TUI tview Implementation - Manual Testing Checklist

**Spec**: tui-tview-implementation (Task 17)
**Implementation**: cmd/tui-tview/
**Date**: _____________
**Tester**: _____________
**Platform**: _____________

## Overview

This checklist validates the tview-based TUI implementation against all requirements from the tui-tview-implementation spec. Complete each test and mark with ✅ (pass), ❌ (fail), or ⚠️ (partial/issues).

---

## Prerequisites Setup

- [ ] Build the latest binary: `go build -o pass-cli.exe .`
- [ ] Run test data setup script: `test/setup-tview-test-data.bat` (Windows) or `test/setup-tview-test-data.sh` (Unix)
- [ ] Verify test vault created at: `test-vault-tview/vault.enc`
- [ ] Test vault password: `test123456`
- [ ] Test vault contains at least 10 credentials across 5+ categories

---

## Test Environment Information

**Operating System**: _____________________ (Windows 10/11, macOS version, Linux distro)
**Terminal Emulator**: __________________ (Windows Terminal, PowerShell, iTerm2, gnome-terminal, Alacritty, etc.)
**Terminal Size**: ______________________ (columns x rows, e.g., 120x30)
**Color Support**: ______________________ (16-color, 256-color, true color)
**Unicode Support**: ____________________ (Full, Partial, None)

---

## Requirement 1: TUI Entry Point and Vault Unlocking

**User Story**: As a developer, I want to launch an interactive TUI dashboard, so that I can visually manage my credentials without using individual CLI commands.

### 1.1 - Keychain Unlock Attempt

**Test**: Launch TUI with keychain configured
```bash
# If keychain is set up from init:
pass-cli tui --vault test-vault-tview/vault.enc
```

- [ ] **Result**: Application attempts keychain unlock
- [ ] **Observed**: ________________________
- [ ] **Status**: ✅ ❌ ⚠️

### 1.2 - Password Prompt Fallback

**Test**: Launch TUI without keychain
```bash
# Delete keychain entry first if exists
pass-cli tui --vault test-vault-tview/vault.enc
```

- [ ] **Result**: Prompts for master password with masked input
- [ ] **Password Masked**: Characters shown as `*` or hidden
- [ ] **Status**: ✅ ❌ ⚠️

### 1.3 - Successful Vault Unlock and Dashboard Display

**Test**: Enter correct password (`test123456`)

- [ ] **Result**: TUI dashboard displays with all panels
- [ ] **Sidebar Visible**: Shows category tree
- [ ] **Table Visible**: Shows credential list
- [ ] **Status Bar Visible**: Shows shortcuts at bottom
- [ ] **No Crashes**: Application starts cleanly
- [ ] **Status**: ✅ ❌ ⚠️

### 1.4 - Failed Unlock After 3 Attempts

**Test**: Enter wrong password 3 times

- [ ] **Result**: Application exits with error code 1
- [ ] **Error Message**: Clear explanation shown
- [ ] **Terminal Restored**: No visual corruption
- [ ] **Exit Code**: Verify `echo %errorlevel%` (Windows) or `echo $?` (Unix) = 1
- [ ] **Status**: ✅ ❌ ⚠️

### 1.5 - Panic Recovery and Terminal Restoration

**Test**: Simulate panic (requires code modification or test build)

- [ ] **Result**: Terminal state restored to normal
- [ ] **Error Message**: Panic error shown to user
- [ ] **No Corruption**: Cursor and colors restored
- [ ] **Exit Code**: Non-zero exit code
- [ ] **Status**: ✅ ❌ ⚠️ (Skip if unable to simulate)

**Notes**: _______________________________________________________

---

## Requirement 2: Central State Management with Deadlock Prevention

**User Story**: As a developer implementing the TUI, I want thread-safe state management that prevents mutex deadlocks, so that the application remains responsive and never hangs.

### 2.1 - Lock Release Before Notifications

**Test**: Perform multiple credential operations rapidly (add, edit, delete)

- [ ] **Result**: No application hangs or freezes
- [ ] **Responsive**: UI updates immediately
- [ ] **No Deadlocks**: Application remains interactive
- [ ] **Status**: ✅ ❌ ⚠️

### 2.2 - Credential Loading and Cache Update

**Test**: Launch TUI, observe initial load

- [ ] **Result**: Credentials displayed in table
- [ ] **Category Tree Populated**: All categories shown
- [ ] **Component Refresh**: All panels show data
- [ ] **Status**: ✅ ❌ ⚠️

### 2.3 - Thread-Safe State Access

**Test**: Navigate rapidly between credentials while resizing terminal

- [ ] **Result**: No crashes or data corruption
- [ ] **Consistent Display**: Selection remains stable
- [ ] **No Race Warnings**: (if running with `-race` build)
- [ ] **Status**: ✅ ❌ ⚠️

### 2.4 - Selection State Updates and Callbacks

**Test**: Select different categories and credentials

- [ ] **Result**: Detail view updates to show selected credential
- [ ] **Sidebar Highlight**: Selected category highlighted
- [ ] **Table Highlight**: Selected credential highlighted
- [ ] **Status**: ✅ ❌ ⚠️

### 2.5 - Error Callbacks Without Lock Holding

**Test**: Trigger error (try deleting non-existent credential or invalid operation)

- [ ] **Result**: Error displayed in status bar or modal
- [ ] **No Hang**: Application remains responsive
- [ ] **Status**: ✅ ❌ ⚠️

### 2.6 - Single Component Instances

**Test**: Verify components are created once and reused

- [ ] **Result**: No duplicate panels or visual glitches
- [ ] **Memory Stable**: No memory leaks during extended use
- [ ] **Status**: ✅ ❌ ⚠️ (Visual inspection)

**Notes**: _______________________________________________________

---

## Requirement 3: Component-Based UI Architecture

**User Story**: As a user, I want a multi-panel dashboard with sidebar navigation, credential list, detail view, and status bar, so that I can efficiently navigate and manage my credentials.

### 3.1 - Sidebar with Category Tree (tview.TreeView)

**Test**: Launch TUI, observe sidebar

- [ ] **Result**: Sidebar panel displays on left
- [ ] **Tree Structure**: Categories shown as tree nodes
- [ ] **Root Node**: "All Credentials" or similar root
- [ ] **Expandable**: Categories can be expanded/collapsed
- [ ] **Border**: Rounded border with title
- [ ] **Status**: ✅ ❌ ⚠️

### 3.2 - Credential List Table (tview.Table)

**Test**: Observe main table panel

- [ ] **Result**: Table displays credentials in columns
- [ ] **Columns**: Service, Username, Category, Last Used (or similar)
- [ ] **Header Row**: Fixed header with column titles
- [ ] **Scrollable**: Can scroll with arrow keys
- [ ] **Border**: Rounded border with title
- [ ] **Status**: ✅ ❌ ⚠️

### 3.3 - Credential Detail View (tview.TextView)

**Test**: Press Enter on a credential in table

- [ ] **Result**: Detail view displays credential information
- [ ] **Service Name**: Displayed prominently
- [ ] **Username**: Visible
- [ ] **Password**: Masked by default (******)
- [ ] **Metadata**: Timestamps, URL, notes shown
- [ ] **Border**: Rounded border with title
- [ ] **Status**: ✅ ❌ ⚠️

### 3.4 - Status Bar (tview.TextView)

**Test**: Observe bottom of screen

- [ ] **Result**: Status bar always visible
- [ ] **Shortcuts Displayed**: Context-aware keyboard shortcuts
- [ ] **Updates**: Changes based on focused panel
- [ ] **No Border**: Flat bar at bottom
- [ ] **Status**: ✅ ❌ ⚠️

### 3.5 - Modal Forms (tview.Form + Modal)

**Test**: Press 'n' to add credential, press 'e' to edit

- [ ] **Result**: Modal form appears centered
- [ ] **Input Fields**: Service, Username, Password, Category, URL, Notes
- [ ] **Buttons**: Add/Save and Cancel buttons
- [ ] **Modal Overlay**: Background dimmed or modal centered
- [ ] **Status**: ✅ ❌ ⚠️

### 3.6 - Component Refresh from State

**Test**: Add credential, verify all components update

- [ ] **Result**: Table shows new credential
- [ ] **Sidebar**: Category count increments
- [ ] **No Duplicate State**: Changes reflected everywhere
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Requirement 4: Responsive Layout Management

**User Story**: As a user working in different terminal sizes, I want the layout to adapt to my terminal width, so that I can use the TUI in various environments.

### 4.1 - Small Layout (<80 columns)

**Test**: Resize terminal to 70 columns

- [ ] **Result**: Only table visible (sidebar hidden, no metadata panel)
- [ ] **Table**: Full width of terminal
- [ ] **Status Bar**: Still visible
- [ ] **Usable**: Can still navigate and use TUI
- [ ] **Status**: ✅ ❌ ⚠️

### 4.2 - Medium Layout (80-120 columns)

**Test**: Resize terminal to 100 columns

- [ ] **Result**: Sidebar and table visible, no metadata panel
- [ ] **Sidebar**: ~20 columns on left
- [ ] **Table**: Remaining width
- [ ] **Status Bar**: Visible
- [ ] **Status**: ✅ ❌ ⚠️

### 4.3 - Large Layout (>120 columns)

**Test**: Resize terminal to 140 columns (if in detail view)

- [ ] **Result**: Sidebar, table, and metadata panel visible
- [ ] **Sidebar**: ~20 columns
- [ ] **Table**: Flexible middle
- [ ] **Metadata Panel**: ~40 columns on right
- [ ] **Status Bar**: Visible
- [ ] **Status**: ✅ ❌ ⚠️

### 4.4 - Dynamic Resize Handling

**Test**: Resize terminal while TUI is running

- [ ] **Result**: Layout recalculates and redraws
- [ ] **Smooth Transition**: No flicker or corruption
- [ ] **Components Adapt**: Panels show/hide correctly
- [ ] **No Crashes**: Application stable during resize
- [ ] **Status**: ✅ ❌ ⚠️

### 4.5 - tview.Flex Layout

**Test**: Visual inspection of layout structure

- [ ] **Result**: Panels arranged with Flex layout
- [ ] **Proportional Sizing**: Components sized appropriately
- [ ] **Modal Management**: Modals overlay main layout
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Requirement 5: Global Event Handling and Keyboard Shortcuts

**User Story**: As a user, I want intuitive keyboard shortcuts for all actions, so that I can navigate and manage credentials efficiently without a mouse.

### 5.1 - Quit Application (q / Ctrl+C)

**Test**: Press 'q', then test Ctrl+C

- [ ] **'q' Result**: Application exits gracefully
- [ ] **Ctrl+C Result**: Application exits gracefully
- [ ] **Terminal Restored**: No visual corruption
- [ ] **Exit Code**: 0 (normal exit)
- [ ] **Status**: ✅ ❌ ⚠️

### 5.2 - Add Credential (n)

**Test**: Press 'n' from table view

- [ ] **Result**: Add credential form modal appears
- [ ] **Form Focused**: Cursor in first input field
- [ ] **Status**: ✅ ❌ ⚠️

### 5.3 - Edit Credential (e)

**Test**: Select credential, press 'e'

- [ ] **Result**: Edit form modal appears
- [ ] **Pre-filled**: Form shows existing credential data
- [ ] **Status**: ✅ ❌ ⚠️

### 5.4 - Delete Credential (d)

**Test**: Select credential, press 'd'

- [ ] **Result**: Confirmation dialog appears
- [ ] **Confirmation Text**: Asks "Are you sure?"
- [ ] **Yes/No Options**: Can confirm or cancel
- [ ] **Status**: ✅ ❌ ⚠️

### 5.5 - Focus Cycling (Tab)

**Test**: Press Tab multiple times

- [ ] **Result**: Focus cycles through panels
- [ ] **Order**: Sidebar → Table → Detail (if visible)
- [ ] **Border Highlight**: Focused panel has distinct border
- [ ] **Status**: ✅ ❌ ⚠️

### 5.6 - Form Input Protection (CRITICAL)

**Test**: Open add/edit form, type 'e', 'n', 'd' in text fields

- [ ] **CRITICAL Result**: Characters appear in input field
- [ ] **NOT Triggered**: Shortcuts do NOT fire
- [ ] **Can Type 'e'**: Letter 'e' appears in field, does NOT open edit dialog
- [ ] **Can Type 'n'**: Letter 'n' appears in field, does NOT open new credential form
- [ ] **Can Type 'd'**: Letter 'd' appears in field, does NOT open delete dialog
- [ ] **Ctrl+C Works**: Ctrl+C still quits from form
- [ ] **Status**: ✅ ❌ ⚠️

**NOTE**: This is a CRITICAL bug fix. Previous implementations incorrectly intercepted form input.

### 5.7 - Search/Filter (/)

**Test**: Press '/' from table view

- [ ] **Result**: Search input appears
- [ ] **Can Type**: Search query accepted
- [ ] **Filters Results**: Credential list filters as you type
- [ ] **Status**: ✅ ❌ ⚠️ (Skip if search not implemented yet)

### 5.8 - Help Screen (?)

**Test**: Press '?' from any view

- [ ] **Result**: Help screen modal appears
- [ ] **Content**: Lists all keyboard shortcuts
- [ ] **Dismissible**: Press Esc or 'q' to close
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Requirement 6: Focus Management Between Components

**User Story**: As a user, I want clear visual indication of which panel is focused, so that I know where my keyboard input will be directed.

### 6.1 - Active Panel Border Highlight

**Test**: Tab through panels, observe borders

- [ ] **Result**: Focused panel has distinct border color
- [ ] **Color Change**: Noticeable difference (e.g., cyan vs gray)
- [ ] **Clear Indication**: Easy to see which panel is active
- [ ] **Status**: ✅ ❌ ⚠️

### 6.2 - Inactive Panel Border Dim

**Test**: Observe unfocused panels

- [ ] **Result**: Unfocused panels have dimmed borders
- [ ] **Consistent**: All unfocused panels use same dim color
- [ ] **Status**: ✅ ❌ ⚠️

### 6.3 - Forward Focus Cycling (Tab)

**Test**: Press Tab repeatedly

- [ ] **Result**: Focus moves forward through panels
- [ ] **Wraps**: After last panel, returns to first
- [ ] **Status**: ✅ ❌ ⚠️

### 6.4 - Reverse Focus Cycling (Shift+Tab)

**Test**: Press Shift+Tab repeatedly

- [ ] **Result**: Focus moves backward through panels
- [ ] **Wraps**: Before first panel, goes to last
- [ ] **Status**: ✅ ❌ ⚠️

### 6.5 - Modal Focus Lock

**Test**: Open add/edit form modal

- [ ] **Result**: Modal receives focus
- [ ] **Background Locked**: Tab does NOT cycle to background panels
- [ ] **Modal Only**: Focus stays within modal
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Requirement 7: Credential CRUD Operations

**User Story**: As a user, I want to create, view, update, and delete credentials through the TUI, so that I can manage my vault without using CLI commands.

### 7.1 - Add Credential

**Test**: Press 'n', fill form, submit

1. Press 'n'
2. Fill: Service=`test-service`, Username=`testuser`, Password=`testpass123`, Category=`Testing`
3. Press Enter or click Add button

- [ ] **Result**: Credential added to vault
- [ ] **Table Updated**: New credential appears in table
- [ ] **Sidebar Updated**: Category count increments
- [ ] **Modal Closed**: Form dismisses after success
- [ ] **Status Message**: Success message in status bar (optional)
- [ ] **Status**: ✅ ❌ ⚠️

### 7.2 - Edit Credential

**Test**: Select credential, press 'e', modify, submit

1. Select credential in table
2. Press 'e'
3. Change Username to `modified-user`
4. Submit form

- [ ] **Result**: Credential updated in vault
- [ ] **Table Shows Changes**: Modified username visible
- [ ] **Modal Closed**: Form dismisses after success
- [ ] **Status**: ✅ ❌ ⚠️

### 7.3 - Delete Credential

**Test**: Select credential, press 'd', confirm

1. Select credential
2. Press 'd'
3. Confirm deletion (Yes)

- [ ] **Result**: Credential removed from vault
- [ ] **Table Updated**: Credential no longer in list
- [ ] **Sidebar Updated**: Category count decrements
- [ ] **Selection Cleared**: No credential selected
- [ ] **Status**: ✅ ❌ ⚠️

**Test Cancel**: Repeat but cancel deletion (No)

- [ ] **Result**: Credential NOT deleted
- [ ] **Still in List**: Credential remains
- [ ] **Status**: ✅ ❌ ⚠️

### 7.4 - Operation Error Handling

**Test**: Try to add credential with empty service name

- [ ] **Result**: Error displayed in modal or status bar
- [ ] **No Crash**: Application remains stable
- [ ] **Form Remains Open**: Can correct and retry
- [ ] **Status**: ✅ ❌ ⚠️

### 7.5 - Callback-Driven Refresh

**Test**: After any CRUD operation, check all components

- [ ] **Result**: All components refresh automatically
- [ ] **No Manual Refresh Needed**: Updates immediate
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Requirement 8: Modern Styling and Theming

**User Story**: As a user, I want a visually appealing interface with modern styling, so that the TUI feels polished and professional.

### 8.1 - Rounded Borders

**Test**: Visual inspection of all panels

- [ ] **Result**: Panels use rounded border characters (╭╮╰╯)
- [ ] **Consistent**: All panels have rounded borders
- [ ] **Rendering**: Characters display correctly (no boxes or fallback)
- [ ] **Status**: ✅ ❌ ⚠️ (Note if ASCII fallback occurs)

**Terminal-Specific Rendering**: ___________________________

### 8.2 - tcell.NewRGBColor() Precision

**Test**: Code inspection (styles/theme.go)

- [ ] **Result**: Colors defined using tcell.NewRGBColor()
- [ ] **RGB Values**: Precise color specifications (e.g., #8be9fd)
- [ ] **Status**: ✅ ❌ ⚠️ (Code review)

### 8.3 - Consistent Color Palette (styles/theme.go)

**Test**: Visual inspection across all components

- [ ] **Result**: Components use consistent theme colors
- [ ] **No Hardcoded Colors**: All colors from theme.go
- [ ] **Cohesive Look**: Professional color scheme
- [ ] **Status**: ✅ ❌ ⚠️

### 8.4 - Status Bar Shortcuts with Highlighting

**Test**: Observe status bar

- [ ] **Result**: Keyboard shortcuts displayed
- [ ] **Highlighting**: Shortcut keys highlighted (e.g., color or brackets)
- [ ] **Readable**: Clear and easy to scan
- [ ] **Status**: ✅ ❌ ⚠️

### 8.5 - Centered Modals with Background

**Test**: Open any modal (add, edit, delete, help)

- [ ] **Result**: Modal centered on screen
- [ ] **Background**: Background visible but dimmed OR modal clearly centered
- [ ] **Proper Sizing**: Modal not too large or too small
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Cross-Platform Testing

### Windows Testing

**Terminals to Test**:
- [ ] **Windows Terminal** (Recommended)
- [ ] **PowerShell**
- [ ] **CMD** (Command Prompt)

**Windows-Specific Checks**:
- [ ] **Rounded Borders**: Render correctly or fallback cleanly
- [ ] **Colors**: Display correctly (true color or 256-color)
- [ ] **Keyboard**: All shortcuts work (no OS conflicts)
- [ ] **Keychain**: Windows Credential Manager integration works
- [ ] **No CRLF Issues**: No visible line ending problems

**Notes**: _______________________________________________________

### macOS Testing (if available)

**Terminals to Test**:
- [ ] **iTerm2** (Recommended)
- [ ] **Terminal.app**

**macOS-Specific Checks**:
- [ ] **Rounded Borders**: Render correctly
- [ ] **Colors**: Display correctly (true color)
- [ ] **Keyboard**: All shortcuts work (no Cmd conflicts)
- [ ] **Keychain**: macOS Keychain integration works
- [ ] **Font Rendering**: Unicode characters display properly

**Notes**: _______________________________________________________

### Linux Testing (if available)

**Terminals to Test**:
- [ ] **gnome-terminal**
- [ ] **Alacritty** (Recommended for performance)
- [ ] **Konsole** (KDE)

**Linux-Specific Checks**:
- [ ] **Rounded Borders**: Render correctly
- [ ] **Colors**: Display correctly
- [ ] **Keyboard**: All shortcuts work
- [ ] **Keychain**: Secret Service (D-Bus) integration works
- [ ] **Font Rendering**: Unicode characters display properly

**Notes**: _______________________________________________________

---

## Terminal Compatibility Testing

### Color Support Tests

**Test Terminal Color Capabilities**:

- [ ] **True Color (24-bit)**: Run `echo $COLORTERM` (should show "truecolor" or "24bit")
- [ ] **256-color**: TUI displays with rich colors
- [ ] **16-color Fallback**: TUI still usable with basic colors
- [ ] **Status**: ✅ ❌ ⚠️

### Border Character Tests

**Test Unicode Box Drawing**:

- [ ] **Rounded Borders**: ╭╮╰╯ characters display correctly
- [ ] **Straight Lines**: ─│ characters display correctly
- [ ] **No Broken Boxes**: No `?` or `□` characters instead of borders
- [ ] **ASCII Fallback**: If Unicode fails, clean ASCII fallback (+-|)
- [ ] **Status**: ✅ ❌ ⚠️

### Mouse Support Tests

**Test Mouse Interaction** (if enabled):

- [ ] **Click to Focus**: Clicking panel focuses it
- [ ] **Click to Select**: Clicking table row selects credential
- [ ] **Scroll Wheel**: Scrolling works in table/detail view
- [ ] **Status**: ✅ ❌ ⚠️ (Skip if mouse disabled)

### Alternate Screen Tests

**Test Terminal State Restoration**:

1. Note current terminal content
2. Launch TUI
3. Quit TUI
4. Check terminal

- [ ] **Result**: Previous terminal content restored
- [ ] **Alternate Screen**: TUI used alternate screen buffer
- [ ] **Cursor Position**: Cursor restored to correct position
- [ ] **No Artifacts**: No leftover TUI elements
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Performance Testing

### Startup Time

**Test**: Measure TUI launch time after vault unlock

```bash
# Windows PowerShell
Measure-Command { echo "test123456" | .\pass-cli.exe tui --vault test-vault-tview\vault.enc }

# Unix
time echo "test123456" | ./pass-cli tui --vault test-vault-tview/vault.enc
```

- [ ] **Target**: < 500ms after vault unlock
- [ ] **Actual**: _______ ms
- [ ] **Status**: ✅ ❌ ⚠️

### UI Responsiveness

**Test**: Navigate rapidly through credentials

- [ ] **Arrow Keys**: Immediate response to navigation
- [ ] **Tab Focus**: Immediate panel focus changes
- [ ] **Modal Open/Close**: < 100ms
- [ ] **No Lag**: Smooth, responsive feel
- [ ] **Status**: ✅ ❌ ⚠️

### Large Credential Set

**Test**: Create vault with 100+ credentials (use helper script)

- [ ] **Table Scrolling**: Smooth scrolling through long list
- [ ] **Sidebar Rendering**: Categories load quickly
- [ ] **Selection**: Credential selection responsive
- [ ] **Memory Usage**: No excessive memory consumption
- [ ] **Status**: ✅ ❌ ⚠️ (Skip if helper script unavailable)

**Notes**: _______________________________________________________

---

## Edge Cases and Error Handling

### Empty Vault

**Test**: Initialize empty vault, launch TUI

- [ ] **Result**: TUI displays gracefully
- [ ] **Empty State**: Shows "No credentials" message
- [ ] **Can Add**: Pressing 'n' opens add form
- [ ] **No Crashes**: Application stable
- [ ] **Status**: ✅ ❌ ⚠️

### Single Credential

**Test**: Vault with only 1 credential

- [ ] **Result**: Layout works with minimal data
- [ ] **Selection**: Can select single credential
- [ ] **Status**: ✅ ❌ ⚠️

### Long Service Names

**Test**: Add credential with very long service name (50+ characters)

- [ ] **Result**: Name truncated with ellipsis (...)
- [ ] **No Overflow**: Text doesn't overflow panel
- [ ] **Full Name Visible**: Full name in detail view
- [ ] **Status**: ✅ ❌ ⚠️

### Special Characters

**Test**: Add credential with special characters in service/username (e.g., `@`, `#`, `$`, `&`)

- [ ] **Result**: Special characters display correctly
- [ ] **No Escaping Issues**: Characters render as-is
- [ ] **Status**: ✅ ❌ ⚠️

### Unicode Passwords

**Test**: Add credential with Unicode password (e.g., `পাসওয়ার্ড`, `パスワード`, `密码`)

- [ ] **Result**: Unicode characters stored correctly
- [ ] **Display**: Characters visible when password unmasked
- [ ] **Copy Works**: Clipboard copy preserves Unicode
- [ ] **Status**: ✅ ❌ ⚠️

### Rapid Key Presses

**Test**: Press keys rapidly (mash keyboard)

- [ ] **Result**: Application doesn't crash
- [ ] **No Lag**: Handles input buffer gracefully
- [ ] **Recovers**: Returns to stable state
- [ ] **Status**: ✅ ❌ ⚠️

### Resize During Operation

**Test**: Resize terminal while form is open or during delete confirmation

- [ ] **Result**: Modal remains centered
- [ ] **Layout Adapts**: Background layout recalculates
- [ ] **No Corruption**: Display remains correct
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Regression Testing

**Verify No Breaking Changes to Core Features**:

### CLI Commands Still Work

**Test CLI commands while TUI exists in binary**:

- [ ] `pass-cli init` - Still works
- [ ] `pass-cli add` - Still works
- [ ] `pass-cli get` - Still works
- [ ] `pass-cli list` - Still works
- [ ] `pass-cli update` - Still works
- [ ] `pass-cli delete` - Still works
- [ ] `pass-cli generate` - Still works
- [ ] **Status**: ✅ ❌ ⚠️

### Vault Compatibility

**Test**: Create credential in TUI, retrieve via CLI (and vice versa)

1. Add credential via TUI
2. Retrieve with `pass-cli get <service>`
3. Add credential via CLI
4. View in TUI

- [ ] **Result**: Credentials compatible between CLI and TUI
- [ ] **Data Integrity**: No corruption or data loss
- [ ] **Status**: ✅ ❌ ⚠️

**Notes**: _______________________________________________________

---

## Bug Report Template

**If bugs are found, document them here**:

### Bug #1

- **Summary**: _______________________________________
- **Severity**: Critical / High / Medium / Low
- **Steps to Reproduce**:
  1. _______________________________________
  2. _______________________________________
  3. _______________________________________
- **Expected Behavior**: _______________________________________
- **Actual Behavior**: _______________________________________
- **Screenshot/Log**: _______________________________________
- **Environment**: OS: ______, Terminal: ______, Size: ______
- **Workaround**: _______________________________________

### Bug #2

_(Add more as needed)_

---

## Test Summary

### Overall Results

- **Total Tests**: ______
- **Passed**: ______
- **Failed**: ______
- **Partial/Warnings**: ______
- **Skipped**: ______

### Critical Issues

List any critical blockers found:
1. _______________________________________
2. _______________________________________

### Recommended Actions

- [ ] **Ready for Release**: All critical tests passed
- [ ] **Needs Fixes**: Address critical bugs before release
- [ ] **Needs Further Testing**: Additional testing required

### Sign-Off

- **Tester Name**: _______________________
- **Date**: _______________________
- **Platform**: _______________________
- **Build Version**: _______________________
- **Approval**: ✅ Approved for Release | ❌ Needs Work | ⚠️ Conditional Approval

**Additional Notes**:
________________________________________________________
________________________________________________________
________________________________________________________

---

## Appendix: Quick Reference

### Keyboard Shortcuts (Expected)

- `q` / `Ctrl+C` - Quit
- `n` - New credential
- `e` - Edit credential
- `d` - Delete credential
- `c` - Copy password
- `p` - Toggle password visibility
- `/` - Search/filter
- `?` - Help
- `Tab` - Cycle focus forward
- `Shift+Tab` - Cycle focus backward
- `Enter` - Select/submit
- `Esc` - Cancel/close modal

### Layout Breakpoints (Expected)

- **< 80 columns**: Table only
- **80-120 columns**: Sidebar + Table
- **> 120 columns**: Sidebar + Table + Detail/Metadata

### Theme Colors (Expected - Dracula-inspired)

- Background: #282a36
- Border Active: #8be9fd (cyan)
- Border Inactive: #6272a4 (gray)
- Text: #f8f8f2 (white)
- Success: #50fa7b (green)
- Error: #ff5555 (red)
- Warning: #f1fa8c (yellow)

---

**End of Manual Testing Checklist**
