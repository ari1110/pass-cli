# TUI Feature Testing Guide - Search, Delete, Copy

**Purpose**: Quick testing guide for three specific TUI features
**Scope**: Search UI, Delete Confirmation, Copy Password
**Prerequisites**: Built binary, test vault setup
**Estimated Time**: 15-20 minutes

---

## Prerequisites Setup

### Build the Binary

**Windows**:
```bash
go build -o pass-cli.exe .
```

**Unix/Linux/macOS**:
```bash
go build -o pass-cli .
```

### Setup Test Vault

**Windows**:
```bash
test\setup-tview-test-data.bat
```

**Unix/Linux/macOS**:
```bash
chmod +x test/setup-tview-test-data.sh
./test/setup-tview-test-data.sh
```

### Launch TUI

**Windows**:
```bash
pass-cli.exe tui --vault test-vault-tview\vault.enc
```

**Unix/Linux/macOS**:
```bash
./pass-cli tui --vault test-vault-tview/vault.enc
```

**Master Password**: `test123456`

---

## Test 1: Search Functionality (/) - Expected: Not Implemented

**Objective**: Verify that search functionality is not implemented in the current version.

**Background**: Code analysis of `cmd/tui-tview/events/handlers.go` shows no handler for the '/' key in the `handleGlobalKey()` function (lines 71-111). This feature is planned for a future release.

### Test Steps:

1. **Launch TUI** and unlock vault with password `test123456`

2. **Navigate to table view** (press Tab until credential table is focused)

3. **Press '/' key**
   - Observe what happens
   - Expected: Nothing happens, or key is ignored
   - NOT expected: Search input field appears

4. **Try from sidebar**
   - Press Tab to focus sidebar
   - Press '/' key
   - Observe behavior

5. **Try from detail view**
   - Select a credential (press Enter)
   - Press Tab to focus detail view
   - Press '/' key
   - Observe behavior

### Expected Results:
- ✅ Pressing '/' does nothing
- ✅ No search UI appears
- ✅ No error messages
- ✅ Application remains stable

### Document:
- What happened when you pressed '/': _______________
- Any unexpected behavior: _______________
- Conclusion: Feature not implemented (as expected) ✅ / Unexpected behavior ❌

**Code Reference**: See `cmd/tui-tview/events/handlers.go` lines 71-111 for keyboard handler implementation.

---

## Test 2: Delete Confirmation Modal (d)

**Objective**: Verify that deleting a credential shows a confirmation dialog and works correctly.

**Background**: Pressing 'd' triggers `handleDeleteCredential()` in `cmd/tui-tview/events/handlers.go` (lines 163-189), which displays a confirmation modal before deletion.

### Part A: Confirm Deletion (Yes)

1. **Launch TUI** and unlock vault

2. **Select a test credential**
   - Navigate to table (Tab key)
   - Use arrow keys to select "github.com" (or any test credential)
   - Note the service name: _______________

3. **Press 'd' key**
   - A modal should appear
   - Verify modal title: "Delete Credential"
   - Verify message: "Delete credential '<service>'?\nThis action cannot be undone."
   - Verify service name is shown: _______________
   - Verify Yes/No buttons are visible

4. **Navigate to "Yes" button**
   - Use Tab key to highlight "Yes" (if needed)
   - Press Enter

5. **Verify deletion**
   - Modal should close automatically
   - Credential should disappear from table
   - Check status bar for success message: _______________
   - Verify sidebar category count decreased (if applicable)

**Expected Results**:
- ✅ Modal appears with correct title and message
- ✅ Service name displayed in message
- ✅ Yes/No buttons visible and functional
- ✅ Credential deleted after confirmation
- ✅ Status message shown (e.g., "Credential deleted")
- ✅ Modal closes automatically

### Part B: Cancel Deletion (No)

1. **Select another credential**
   - Use arrow keys to select "gitlab.com" (or any test credential)
   - Note the service name: _______________

2. **Press 'd' key**
   - Modal appears

3. **Navigate to "No" button**
   - Use Tab key to highlight "No"
   - Press Enter

4. **Verify cancellation**
   - Modal should close
   - Credential should remain in table
   - No changes to vault
   - No deletion message

**Expected Results**:
- ✅ Modal closes without deleting
- ✅ Credential remains in table
- ✅ No status message (or "Cancelled" message)

### Part C: Cancel with Esc

1. **Select a credential and press 'd'**

2. **Press Esc key**

3. **Verify**
   - Modal closes
   - Credential not deleted

**Expected Results**:
- ✅ Esc closes modal without deleting

### Part D: Edge Cases

1. **No selection test**
   - Deselect all (if possible) or navigate to empty area
   - Press 'd' key
   - Expected: Error message "no credential selected" in status bar
   - Observed: _______________

2. **Rapid press test**
   - Select credential
   - Press 'd' multiple times rapidly
   - Expected: Only one modal appears, no crashes
   - Observed: _______________

### Document:
- Modal appearance: ✅ Correct / ❌ Issues: _______________
- Yes button works: ✅ / ❌
- No button works: ✅ / ❌
- Esc cancels: ✅ / ❌
- Error handling: ✅ / ❌
- Any issues: _______________

**Code Reference**: See `cmd/tui-tview/events/handlers.go` lines 163-189 for delete handler implementation.

---

## Test 3: Copy Password to Clipboard (c)

**Objective**: Verify that pressing 'c' copies the password to clipboard and shows feedback.

**Background**: Pressing 'c' triggers `handleCopyPassword()` in `cmd/tui-tview/events/handlers.go` (lines 200-212), which uses the `github.com/atotto/clipboard` library to copy the password and displays a success message via the status bar.

### Part A: Basic Copy Operation

1. **Launch TUI** and unlock vault

2. **Select a credential**
   - Navigate to table
   - Select "github.com" (or any credential with known password)
   - Note the service name: _______________
   - Note the expected password: _______________ (check with 'p' to reveal)

3. **Press 'c' key**
   - Observe status bar at bottom of screen
   - Expected: Green message "Password copied to clipboard!"
   - Observed message: _______________
   - Observed color: _______________
   - Message should display for ~3 seconds then revert to shortcuts

4. **Verify clipboard content**
   - Open a text editor (Notepad, TextEdit, nano, etc.)
   - Paste clipboard content (Ctrl+V on Windows/Linux, Cmd+V on macOS)
   - Verify pasted text matches expected password
   - Password correct: ✅ / ❌ (actual: _______________)

**Expected Results**:
- ✅ 'c' key responds immediately
- ✅ Green success message appears in status bar
- ✅ Message text: "Password copied to clipboard!"
- ✅ Message displays for ~3 seconds
- ✅ Clipboard contains correct password
- ✅ No corruption or truncation

### Part B: Error Handling (No Selection)

1. **Deselect credential** (if possible, or start fresh TUI session)

2. **Press 'c' key without selecting a credential**
   - Expected: Red error message in status bar
   - Expected text: "no credential selected" or similar
   - Observed message: _______________
   - Observed color: _______________
   - Message should display for ~5 seconds

**Expected Results**:
- ✅ Red error message appears
- ✅ Message indicates no selection
- ✅ No crash or hang

### Part C: Unicode Password Test

1. **Add credential with Unicode password**
   - Press 'n' to open add form
   - Service: "unicode-test"
   - Username: "testuser"
   - Password: "パスワード123" (or any Unicode string)
   - Category: "Testing"
   - Submit form

2. **Select the new credential**

3. **Press 'c' to copy**

4. **Paste and verify**
   - Paste in text editor
   - Verify Unicode characters preserved: _______________
   - Characters correct: ✅ / ❌

**Expected Results**:
- ✅ Unicode password copied correctly
- ✅ No encoding issues
- ✅ Characters display correctly when pasted

### Part D: Context Testing

1. **Test from sidebar**
   - Focus sidebar (Tab key)
   - Select a category node
   - Press 'c'
   - Observe behavior: _______________

2. **Test from table**
   - Focus table (Tab key)
   - Select credential
   - Press 'c'
   - Verify copy works: ✅ / ❌

3. **Test from detail view**
   - Focus detail view (Tab key)
   - Press 'c'
   - Verify copy works: ✅ / ❌

**Expected Results**:
- ✅ Copy works from all contexts where a credential is selected
- ✅ Consistent behavior across contexts

### Part E: Rapid Copy Operations

1. **Select credential A** (e.g., "github.com")
   - Press 'c'
   - Note password: _______________

2. **Immediately select credential B** (e.g., "gitlab.com")
   - Press 'c' immediately
   - Note password: _______________

3. **Paste clipboard**
   - Paste in text editor
   - Verify clipboard contains credential B password (most recent)
   - Correct: ✅ / ❌

**Expected Results**:
- ✅ Clipboard updated to most recent copy
- ✅ No crashes or hangs
- ✅ Status messages display correctly for each operation

### Part F: Platform-Specific Testing

**Windows**:
- [ ] Test in Windows Terminal
- [ ] Test in PowerShell
- [ ] Test in CMD
- [ ] Verify Ctrl+V paste works
- [ ] Verify right-click paste works
- [ ] Any issues: _______________

**macOS** (if available):
- [ ] Test in iTerm2
- [ ] Test in Terminal.app
- [ ] Verify Cmd+V paste works
- [ ] Any issues: _______________

**Linux** (if available):
- [ ] Desktop environment: _______________
- [ ] Test in terminal (Ctrl+Shift+V)
- [ ] Test in GUI app (Ctrl+V)
- [ ] Clipboard mechanism: X11 / Wayland / Other
- [ ] Any issues: _______________

### Document:
- Copy works: ✅ / ❌
- Success message correct: ✅ / ❌
- Error handling works: ✅ / ❌
- Unicode preserved: ✅ / ❌ / N/A
- Platform-specific issues: _______________
- Overall status: ✅ Pass / ❌ Fail / ⚠️ Issues

**Code Reference**:
- Handler: `cmd/tui-tview/events/handlers.go` lines 200-212
- Copy implementation: `cmd/tui-tview/components/detail.go` lines 167-188
- Status bar feedback: `cmd/tui-tview/components/statusbar.go` lines 67-77

---

## Test Completion Checklist

- [ ] Test 1: Search (/) - Verified not implemented
- [ ] Test 2: Delete (d) - All parts completed (A, B, C, D)
- [ ] Test 3: Copy (c) - All parts completed (A, B, C, D, E, F)
- [ ] Results documented in test report
- [ ] Screenshots captured (if needed)
- [ ] Issues logged (if any)

---

## Next Steps

1. **Fill out test report**: Use `TVIEW_TEST_REPORT_TEMPLATE.md` to document detailed results

2. **Update checklist**: Mark completed sections in `TVIEW_MANUAL_TESTING_CHECKLIST.md`

3. **Report issues**: If bugs found, document in the Bug Report Template section of the checklist

4. **Sign off**: Complete the test summary and sign-off sections

---

## Quick Reference

### Keyboard Shortcuts:
- `q` - Quit
- `n` - New credential
- `e` - Edit credential
- `d` - Delete credential (shows confirmation)
- `c` - Copy password to clipboard
- `p` - Toggle password visibility
- `?` - Help
- `Tab` - Cycle focus
- `Esc` - Cancel/close modal

### Test Vault:
- Path: `test-vault-tview/vault.enc`
- Password: `test123456`
- Contains: 10+ test credentials across multiple categories

### Expected Behavior Summary:
- **Search (/)**: Not implemented (no action)
- **Delete (d)**: Shows confirmation modal, deletes on Yes, cancels on No/Esc
- **Copy (c)**: Copies password to clipboard, shows green success message for 3s, red error for 5s if no selection

---

**End of Testing Guide**
