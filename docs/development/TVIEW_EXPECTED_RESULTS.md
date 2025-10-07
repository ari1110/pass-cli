# tview TUI Implementation - Expected Results Guide

**Spec**: tui-tview-implementation
**Purpose**: Detailed expected results for manual testing
**Audience**: QA testers, developers validating implementation

---

## Overview

This document provides detailed expected results for each test scenario in the tview TUI implementation. Use this as a reference when performing manual testing to validate that the implementation meets specifications.

---

## Table of Contents

1. [Vault Unlocking and Entry](#vault-unlocking-and-entry)
2. [Component Display and Layout](#component-display-and-layout)
3. [Keyboard Shortcuts](#keyboard-shortcuts)
4. [CRUD Operations](#crud-operations)
5. [Focus Management](#focus-management)
6. [Responsive Layout](#responsive-layout)
7. [Visual Styling](#visual-styling)
8. [Error Handling](#error-handling)

---

## Vault Unlocking and Entry

### Test: Launch TUI with Keychain Configured

**Command**: `pass-cli tui --vault test-vault-tview/vault.enc`

**Expected Behavior**:
1. Application starts immediately
2. No password prompt appears (keychain provides credentials)
3. TUI dashboard displays within 500ms
4. No error messages shown

**Visual Indicators**:
- ✅ Dashboard appears without interaction
- ✅ Status bar may show "🔓 Keychain" indicator (if implemented)

**If Keychain Fails**: Application should fall back to password prompt

---

### Test: Launch TUI without Keychain (Password Prompt)

**Command**: `pass-cli tui --vault test-vault-tview/vault.enc` (with keychain entry deleted)

**Expected Behavior**:
1. Terminal displays: `Enter master password: `
2. Cursor waits for input
3. Typing shows masked characters: `*` or hidden completely
4. Pressing Enter validates password

**Visual Indicators**:
- ✅ Prompt text is clear and descriptive
- ✅ Password input is masked (security)
- ✅ No password characters visible in terminal

**Correct Password** → Dashboard displays
**Incorrect Password** → Error message, retry (up to 3 attempts total)

---

### Test: Failed Unlock After 3 Attempts

**Scenario**: Enter wrong password 3 times

**Expected Behavior**:
1. **Attempt 1**: "Incorrect password. Try again. (2 attempts remaining)"
2. **Attempt 2**: "Incorrect password. Try again. (1 attempt remaining)"
3. **Attempt 3**: "Incorrect password. Maximum attempts exceeded."
4. Application exits
5. Terminal restored to normal state
6. Exit code 1 (verify with `echo %errorlevel%` on Windows or `echo $?` on Unix)

**Visual Indicators**:
- ✅ Clear error messages
- ✅ Attempt counter visible
- ✅ No terminal corruption after exit

---

### Test: Successful Dashboard Launch

**Expected Behavior**: After successful unlock, dashboard displays with:

**Layout** (assuming terminal width > 120 columns):
```
╭─ Sidebar ──────────╮  ╭─ Credentials ────────────────────╮  ╭─ Detail ──────╮
│                    │  │ Service    | Username | Category │  │               │
│ All Credentials    │  │ aws-prod   | admin    | Cloud    │  │ No credential │
│ ├─ Cloud (3)       │  │ github     | myuser   | Git      │  │ selected      │
│ ├─ Databases (3)   │  │ ...        | ...      | ...      │  │               │
│ ├─ APIs (1)        │  ╰──────────────────────────────────╯  ╰───────────────╯
│ └─ ...             │
╰────────────────────╯
[q] Quit  [n] New  [e] Edit  [d] Delete  [?] Help                    5 credentials
```

**Expected Components**:
1. **Sidebar** (left): Category tree with rounded borders
2. **Table** (center): Credential list with columns
3. **Detail** (right): Empty or showing selected credential
4. **Status Bar** (bottom): Keyboard shortcuts

**Visual Checks**:
- ✅ All panels have rounded borders (╭╮╰╯)
- ✅ Panel titles visible ("Sidebar", "Credentials", "Detail")
- ✅ Status bar shows context-aware shortcuts
- ✅ Active panel has highlighted border (different color)

---

## Component Display and Layout

### Sidebar Component (tview.TreeView)

**Expected Visual Structure**:
```
╭─ Sidebar ──────────╮
│ All Credentials    │  ← Root node (always visible)
│ ├─ Cloud (3)       │  ← Category with count
│ │  ├─ aws-prod     │  ← Credentials under category
│ │  ├─ aws-dev      │
│ │  └─ azure        │
│ ├─ Databases (3)   │
│ │  ├─ postgres     │
│ │  └─ ...          │
│ └─ Uncategorized   │
╰────────────────────╯
```

**Expected Behavior**:
- ✅ Tree structure with expand/collapse (▶ collapsed, ▼ expanded)
- ✅ Category counts shown in parentheses
- ✅ Credentials listed under their categories
- ✅ Arrow keys navigate tree
- ✅ Enter expands/collapses categories or selects credential
- ✅ Selected item highlighted

**Colors**:
- Border: Cyan (#8be9fd) when focused, Gray (#6272a4) when unfocused
- Selected item: Highlighted background

---

### Table Component (tview.Table)

**Expected Visual Structure**:
```
╭─ Credentials ──────────────────────────────────────────────╮
│ Service          | Username      | Category    | Last Used │  ← Header row
│ ────────────────────────────────────────────────────────── │
│ aws-production   | admin         | Cloud       | 2h ago    │  ← Data rows
│ github-personal  | myusername    | Git         | 1d ago    │
│ postgres-main    | dbadmin       | Databases   | 3h ago    │
│ ...              | ...           | ...         | ...       │
╰────────────────────────────────────────────────────────────╯
```

**Expected Behavior**:
- ✅ Fixed header row (doesn't scroll)
- ✅ Columns: Service, Username, Category, Last Used
- ✅ Arrow keys (↑↓) navigate rows
- ✅ Selected row highlighted
- ✅ Enter opens detail view
- ✅ Scrolls smoothly with long lists

**Filtering**:
- When category selected in sidebar → Only shows credentials from that category
- "All Credentials" → Shows all credentials

---

### Detail View Component (tview.TextView)

**Expected Display** (when credential selected):
```
╭─ Detail ───────────────────────╮
│ Service: aws-production        │
│ Username: admin                │
│ Password: ********             │  ← Masked by default
│ Category: Cloud                │
│ URL: https://aws.amazon.com    │
│ Notes: Production AWS account  │
│                                │
│ Created: 2024-01-15 10:30      │
│ Modified: 2024-01-20 14:22     │
│                                │
│ Last Used: 2 hours ago         │
│ Location: /home/user/project   │
╰────────────────────────────────╯
```

**Expected Behavior**:
- ✅ Shows all credential metadata
- ✅ Password masked as `********` by default
- ✅ Press 'p' to toggle password visibility
- ✅ Password shown in plain text when unmasked
- ✅ Scrollable if content exceeds panel height

**When No Credential Selected**:
```
╭─ Detail ───────────────────────╮
│                                │
│   No credential selected       │
│                                │
│   Select a credential to       │
│   view details                 │
│                                │
╰────────────────────────────────╯
```

---

### Status Bar Component

**Expected Display** (bottom of screen):
```
[q] Quit  [n] New  [e] Edit  [d] Delete  [c] Copy  [p] Toggle  [?] Help        15 credentials
```

**Context-Aware Shortcuts**:

**When Table Focused**:
```
[n] New  [e] Edit  [d] Delete  [/] Search  [?] Help  [q] Quit
```

**When Detail View Focused**:
```
[e] Edit  [d] Delete  [c] Copy  [p] Toggle  [?] Help  [q] Quit
```

**When Sidebar Focused**:
```
[Enter] Select  [↑↓] Navigate  [?] Help  [q] Quit
```

**When Modal Open**:
```
[Enter] Submit  [Esc] Cancel
```

**Expected Behavior**:
- ✅ Shortcuts update based on focused component
- ✅ Shortcut keys highlighted (e.g., different color or brackets)
- ✅ Credential count shown on right (e.g., "15 credentials")
- ✅ No border on status bar (flat bar)

---

### Modal Forms

#### Add Credential Form

**Expected Display** (centered on screen):
```
       ╭─ Add Credential ─────────────────────────╮
       │                                          │
       │  Service:    [                        ]  │
       │  Username:   [                        ]  │
       │  Password:   [                        ]  │  ← Masked input
       │  Category:   [                        ]  │
       │  URL:        [                        ]  │
       │  Notes:      [                        ]  │
       │                                          │
       │         [ Add ]      [ Cancel ]          │
       │                                          │
       ╰──────────────────────────────────────────╯
```

**Expected Behavior**:
- ✅ Modal centered on screen
- ✅ Background dimmed or modal clearly overlay
- ✅ Cursor in first field (Service)
- ✅ Tab moves between fields
- ✅ Password field shows masked input (*** as you type)
- ✅ Enter on last field OR clicking "Add" button submits
- ✅ Esc OR clicking "Cancel" closes without saving

**CRITICAL - Form Input Protection**:
- ✅ Typing 'e' in field → Letter 'e' appears in field
- ✅ Typing 'n' in field → Letter 'n' appears in field
- ✅ Typing 'd' in field → Letter 'd' appears in field
- ❌ Shortcuts MUST NOT trigger while typing in form

**Validation**:
- Empty service name → Error: "Service is required"
- Form stays open for correction

---

#### Edit Credential Form

**Expected Display** (pre-populated with existing data):
```
       ╭─ Edit Credential ────────────────────────╮
       │                                          │
       │  Service:    [aws-production          ]  │  ← Pre-filled
       │  Username:   [admin                   ]  │  ← Pre-filled
       │  Password:   [********************    ]  │  ← Pre-filled, masked
       │  Category:   [Cloud                   ]  │  ← Pre-filled
       │  URL:        [https://aws.amazon.com  ]  │  ← Pre-filled
       │  Notes:      [Production AWS account  ]  │  ← Pre-filled
       │                                          │
       │        [ Save ]      [ Cancel ]          │
       │                                          │
       ╰──────────────────────────────────────────╯
```

**Expected Behavior**:
- ✅ All fields pre-populated with existing values
- ✅ Can modify any field
- ✅ "Save" button (not "Add")
- ✅ Same input protection as Add form

---

#### Delete Confirmation Dialog

**Expected Display**:
```
       ╭─ Confirm Delete ─────────────────────────╮
       │                                          │
       │  Are you sure you want to delete this   │
       │  credential?                             │
       │                                          │
       │  Service: aws-production                 │
       │                                          │
       │  This action cannot be undone.           │
       │                                          │
       │         [ Yes ]      [ No ]              │
       │                                          │
       ╰──────────────────────────────────────────╯
```

**Expected Behavior**:
- ✅ Shows service name being deleted
- ✅ Clear warning message
- ✅ "Yes" button deletes credential
- ✅ "No" button cancels (credential preserved)
- ✅ Esc key same as "No"

---

## Keyboard Shortcuts

### Global Shortcuts (work outside forms)

| Key | Action | Expected Result |
|-----|--------|-----------------|
| `q` | Quit | Application exits gracefully, terminal restored |
| `Ctrl+C` | Force quit | Immediate exit, terminal restored |
| `n` | New credential | Add credential form modal appears |
| `e` | Edit credential | Edit form appears (only if credential selected) |
| `d` | Delete credential | Confirmation dialog appears (only if selected) |
| `c` | Copy password | Password copied to clipboard, status shows "Copied!" |
| `p` | Toggle password | Password visibility toggled in detail view |
| `/` | Search/filter | Search input appears (if implemented) |
| `?` | Help | Help screen modal appears with all shortcuts |
| `Tab` | Cycle focus | Focus moves to next panel (Sidebar → Table → Detail) |
| `Shift+Tab` | Reverse cycle | Focus moves to previous panel |
| `↑` | Navigate up | Selection moves up in focused component |
| `↓` | Navigate down | Selection moves down in focused component |
| `Enter` | Select/Submit | Selects item or submits form |
| `Esc` | Cancel/Close | Closes modal or returns to previous view |

**Critical Test - Form Input Protection**:

When form is open and cursor is in input field:
- ✅ Typing 'e' → Character 'e' appears in field (does NOT open edit modal)
- ✅ Typing 'n' → Character 'n' appears in field (does NOT open new credential modal)
- ✅ Typing 'd' → Character 'd' appears in field (does NOT open delete dialog)
- ✅ Typing '/' → Character '/' appears in field (does NOT activate search)
- ✅ Ctrl+C still works to quit application

**This is a critical bug fix from previous implementations**

---

## CRUD Operations

### Add Credential Operation

**Steps**:
1. Press 'n'
2. Fill form: Service=`test-service`, Username=`testuser`, Password=`testpass123`, Category=`Testing`
3. Press Enter or click "Add"

**Expected Results**:

**Immediate**:
- ✅ Modal closes automatically
- ✅ Status bar shows "Credential added successfully" (brief message)

**Table Update**:
- ✅ New row appears in table: `test-service | testuser | Testing | Just now`
- ✅ Row is automatically selected

**Sidebar Update**:
- ✅ "Testing" category appears (if new)
- ✅ Category count increments: "Testing (1)"
- ✅ Credential appears under "Testing" category in tree

**Persistence**:
- ✅ Credential saved to vault
- ✅ Visible in CLI: `pass-cli get test-service --vault test-vault-tview/vault.enc`

---

### Edit Credential Operation

**Steps**:
1. Select credential in table (e.g., `test-service`)
2. Press 'e'
3. Modify field (e.g., change Username to `modified-user`)
4. Press Enter or click "Save"

**Expected Results**:

**Immediate**:
- ✅ Modal closes
- ✅ Status bar shows "Credential updated successfully"

**Table Update**:
- ✅ Row shows new username: `test-service | modified-user | Testing | Just now`
- ✅ Last Used timestamp updates to "Just now"

**Detail View Update**:
- ✅ Shows modified username
- ✅ Modified timestamp updates

**Persistence**:
- ✅ Changes saved to vault
- ✅ CLI shows updated data: `pass-cli get test-service`

---

### Delete Credential Operation

**Steps**:
1. Select credential (e.g., `test-service`)
2. Press 'd'
3. Confirmation dialog appears
4. Click "Yes" or press Enter

**Expected Results**:

**Immediate**:
- ✅ Confirmation dialog closes
- ✅ Status bar shows "Credential deleted"

**Table Update**:
- ✅ Row removed from table
- ✅ Selection clears (no credential selected)
- ✅ If last credential in category, empty category behavior:
  - Category still visible with count (0)
  - OR category removed from tree (implementation choice)

**Sidebar Update**:
- ✅ Category count decrements: "Testing (0)" or category removed

**Detail View**:
- ✅ Clears to "No credential selected" state

**Persistence**:
- ✅ Credential removed from vault
- ✅ CLI shows not found: `pass-cli get test-service` → Error

**Cancel Delete**:
- ✅ Pressing "No" or Esc → Credential NOT deleted
- ✅ Everything remains unchanged

---

### Copy Password Operation

**Steps**:
1. Select credential with password
2. Press 'c'

**Expected Results**:

**Immediate**:
- ✅ Status bar shows "Password copied to clipboard!"
- ✅ Message displays for 2-3 seconds

**Clipboard**:
- ✅ Password is in clipboard (verify by pasting: Ctrl+V)
- ✅ Exact password copied (no extra whitespace or formatting)

**No Visual Change**:
- ✅ Password remains masked in detail view (if was masked)

**Verify**:
```bash
# After copying, paste in text editor
# Should see exact password: AWSProd@2024!SecureKey
```

---

### Toggle Password Visibility

**Steps**:
1. Select credential
2. Observe detail view shows: `Password: ********`
3. Press 'p'

**Expected Results**:

**After First Press**:
- ✅ Detail view shows: `Password: AWSProd@2024!SecureKey` (plain text)

**After Second Press**:
- ✅ Detail view shows: `Password: ********` (masked again)

**Toggle Behavior**:
- ✅ Repeatable (can toggle multiple times)
- ✅ State preserved while credential selected
- ✅ Resets to masked when different credential selected

---

## Focus Management

### Focus Cycling (Tab Key)

**Scenario**: Full layout (width > 120 cols) with Sidebar, Table, Detail visible

**Expected Sequence**:
1. **Initial**: Sidebar focused (border cyan/highlighted)
2. **Press Tab**: Table focused (sidebar border dims, table border highlights)
3. **Press Tab**: Detail focused (table border dims, detail border highlights)
4. **Press Tab**: Sidebar focused (wraps back to start)

**Visual Indicators**:
- ✅ Active panel: Border color #8be9fd (cyan/bright)
- ✅ Inactive panels: Border color #6272a4 (gray/dim)
- ✅ Change is immediate and obvious

**Reverse Cycling** (Shift+Tab):
- Same sequence in reverse order

---

### Focus Cycling in Medium Layout (80-120 cols)

**Layout**: Sidebar + Table (no Detail)

**Expected Sequence**:
1. Sidebar focused
2. Press Tab → Table focused
3. Press Tab → Sidebar focused (only 2 panels)

**Detail Panel**:
- ✅ Hidden (not part of focus cycle)

---

### Focus Cycling in Small Layout (<80 cols)

**Layout**: Table only (no Sidebar, no Detail)

**Expected Behavior**:
- ✅ Tab does nothing (only 1 panel)
- ✅ Or: Tab cycles but always returns to Table
- ✅ Table remains usable with arrow keys

---

### Modal Focus Lock

**Scenario**: Open any modal (Add form, Edit form, Delete confirm, Help)

**Expected Behavior**:
- ✅ Modal receives focus immediately
- ✅ Tab moves between fields within modal (for forms)
- ✅ Tab does NOT cycle to background panels
- ✅ Background panels still visible but not interactive
- ✅ Esc closes modal and returns focus to previous panel

---

## Responsive Layout

### Terminal Width: 150 columns (Large Layout)

**Expected Layout**:
```
┌─────────────┬───────────────────────────────────┬─────────────┐
│  Sidebar    │        Main Content Area          │   Detail    │
│   (20 cols) │            (90 cols)              │  (40 cols)  │
│             │                                   │             │
│  Category   │  Credential Table                 │  Metadata   │
│  Tree       │  Service | Username | Category   │  Panel      │
│             │                                   │             │
└─────────────┴───────────────────────────────────┴─────────────┘
Status Bar (full width)
```

**Visual Checks**:
- ✅ All 3 panels visible
- ✅ Sidebar ~20 columns (fixed width)
- ✅ Detail ~40 columns (fixed width)
- ✅ Table takes remaining space (flexible)
- ✅ All panels have borders
- ✅ No overlapping or gaps

---

### Terminal Width: 100 columns (Medium Layout)

**Expected Layout**:
```
┌─────────────┬───────────────────────────────────────────────┐
│  Sidebar    │        Main Content Area (Table)              │
│   (20 cols) │                (80 cols)                      │
│             │                                               │
│  Category   │  Credential Table                             │
│  Tree       │  Service | Username | Category | Last Used   │
│             │                                               │
└─────────────┴───────────────────────────────────────────────┘
Status Bar (full width)
```

**Visual Checks**:
- ✅ Sidebar visible (~20 columns)
- ✅ Table visible (remaining width)
- ❌ Detail panel hidden (not enough space)
- ✅ Status bar shows "Press [d] to view details" or similar hint

---

### Terminal Width: 70 columns (Small Layout)

**Expected Layout**:
```
┌──────────────────────────────────────────────────────────────┐
│           Main Content Area (Table Only)                     │
│              (Full width - 70 cols)                          │
│                                                              │
│  Credential Table                                            │
│  Service          | Username   | Category                   │
│                                                              │
└──────────────────────────────────────────────────────────────┘
Status Bar (full width)
```

**Visual Checks**:
- ❌ Sidebar hidden
- ✅ Table uses full width
- ❌ Detail panel hidden
- ✅ Table remains functional
- ✅ Can still navigate and select credentials

---

### Dynamic Resize Test

**Steps**:
1. Launch TUI at 150 columns
2. Gradually resize terminal smaller: 150 → 120 → 100 → 80 → 70

**Expected Behavior at Each Breakpoint**:

**150 → 120 columns**:
- ✅ Smooth transition
- ❌ Detail panel disappears
- ✅ Sidebar and Table remain

**100 → 80 columns**:
- ✅ Layout recalculates
- ✅ Sidebar and Table adjust widths

**80 → 70 columns**:
- ✅ Sidebar disappears
- ✅ Table takes full width

**Reverse Resize** (70 → 150):
- ✅ Panels reappear at appropriate breakpoints
- ✅ No visual glitches or corruption
- ✅ Selection and state preserved

**Performance**:
- ✅ Resize happens smoothly (no flicker)
- ✅ No lag or delay
- ✅ Application remains responsive

---

## Visual Styling

### Rounded Borders

**Expected Characters**:
```
╭─────────────╮  ← Top corners: ╭ ╮
│             │  ← Sides: │
│             │
╰─────────────╯  ← Bottom corners: ╰ ╯
```

**Terminals with Full Unicode Support**:
- ✅ Rounded corners display correctly
- ✅ Consistent across all panels

**Terminals without Unicode** (fallback):
```
+-------------+  ← ASCII fallback
|             |
|             |
+-------------+
```

**Visual Checks**:
- ✅ No broken characters (no `?` or `□`)
- ✅ Borders form complete rectangles
- ✅ Panel titles visible inside top border

---

### Color Palette (Dracula Theme)

**Expected Colors** (RGB values):

| Element | Color | RGB | Visual |
|---------|-------|-----|--------|
| Background | Background Dark | #282a36 | Dark purple-gray |
| Border Active | Cyan | #8be9fd | Bright cyan/blue |
| Border Inactive | Gray | #6272a4 | Muted blue-gray |
| Text Primary | White | #f8f8f2 | Off-white |
| Success | Green | #50fa7b | Bright green |
| Error | Red | #ff5555 | Bright red |
| Warning | Yellow | #f1fa8c | Pale yellow |

**Visual Checks**:
- ✅ Active panel border is noticeably brighter (cyan)
- ✅ Inactive borders are dimmer (gray)
- ✅ Success messages in green
- ✅ Error messages in red
- ✅ Overall theme is cohesive and professional

**Terminal Color Support**:
- **True Color (24-bit)**: Full Dracula colors
- **256-color**: Close approximations
- **16-color**: Basic fallback (still usable)

---

### Focus Visual Indication

**Active Panel**:
```
╭─ Sidebar ───────────╮  ← Bright cyan border #8be9fd
│ All Credentials     │
│ ├─ Cloud (3)        │
╰─────────────────────╯
```

**Inactive Panel**:
```
╭─ Credentials ───────╮  ← Dim gray border #6272a4
│ Service | Username  │
│ aws     | admin     │
╰──────────────────────╯
```

**Expected Visibility**:
- ✅ Immediately obvious which panel is focused
- ✅ Color difference is clear even in different lighting
- ✅ Works for colorblind users (brightness difference, not just hue)

---

### Status Bar Styling

**Expected Appearance**:
```
[n] New  [e] Edit  [d] Delete  [?] Help  [q] Quit                15 credentials
└─────────────────────────────────────────────────────────┘
```

**Styling Details**:
- ✅ No border (flat bar)
- ✅ Dark background (distinct from panels)
- ✅ Shortcut keys highlighted (brackets or color)
- ✅ Text readable against background
- ✅ Credential count right-aligned

---

## Error Handling

### Error: Empty Service Name in Add Form

**Steps**:
1. Press 'n' to open Add form
2. Leave Service field empty
3. Fill other fields
4. Press Enter to submit

**Expected Behavior**:
- ✅ Form does NOT close
- ✅ Error message appears: "Service name is required" (red text)
- ✅ Cursor remains in form
- ✅ Can correct error and resubmit
- ✅ Other field values preserved

---

### Error: Edit Non-Existent Credential

**Steps**:
1. No credential selected in table
2. Press 'e'

**Expected Behavior**:
- ✅ Nothing happens OR
- ✅ Status bar shows: "Please select a credential first"
- ✅ No crash or error modal

---

### Error: Delete Non-Existent Credential

**Steps**:
1. No credential selected
2. Press 'd'

**Expected Behavior**:
- ✅ Nothing happens OR
- ✅ Status bar shows: "Please select a credential first"
- ✅ No crash

---

### Error: Clipboard Copy Failure

**Scenario**: System clipboard unavailable (rare)

**Expected Behavior**:
- ✅ Status bar shows: "Failed to copy to clipboard"
- ✅ Error message in red
- ✅ Application remains stable
- ✅ Can retry copy

---

### Error: Empty Vault

**Steps**:
1. Initialize new empty vault
2. Launch TUI

**Expected Behavior**:
- ✅ TUI displays normally
- ✅ Sidebar shows "All Credentials" (empty)
- ✅ Table shows header but no rows OR message: "No credentials yet"
- ✅ Status bar shows: "0 credentials"
- ✅ Pressing 'n' works to add first credential

---

### Error: Terminal Too Small

**Steps**:
1. Resize terminal to < 60 columns or < 20 rows

**Expected Behavior** (if minimum size warning implemented):
```
╭─ Warning ───────────────────╮
│ Terminal too small          │
│                             │
│ Minimum size: 60x20         │
│ Current size: 50x15         │
│                             │
│ Please resize your terminal │
╰─────────────────────────────╯
```

**OR** (graceful degradation):
- ✅ Layout adapts to very small size
- ✅ Components remain usable (scrollable)
- ✅ No visual corruption

---

## Performance Expectations

### Startup Time

**Measurement**: Time from vault unlock to dashboard display

**Target**: < 500ms

**Acceptable**: < 1000ms (1 second)

**Method**:
```powershell
# PowerShell
Measure-Command { echo "test123456" | .\pass-cli.exe tui --vault test-vault-tview\vault.enc }

# Bash
time echo "test123456" | ./pass-cli tui --vault test-vault-tview/vault.enc
```

**Expected**: 100-500ms typical

---

### UI Responsiveness

**Keyboard Input**:
- ✅ Key press to action: < 50ms (feels instant)
- ✅ No input lag or buffering
- ✅ Smooth navigation

**Tab Focus Cycling**:
- ✅ Border color change: Immediate (< 16ms, single frame)

**Modal Open/Close**:
- ✅ Appears: < 100ms
- ✅ Closes: < 100ms

**Terminal Resize**:
- ✅ Layout recalculation: < 100ms
- ✅ Smooth redraw (no flicker)

---

### Large Data Set Performance

**Test Data**: 100+ credentials

**Table Scrolling**:
- ✅ Arrow key navigation: Smooth
- ✅ Jump to top/bottom: Fast
- ✅ No lag or stuttering

**Selection**:
- ✅ Detail view updates: < 50ms

**Search/Filter** (if implemented):
- ✅ Results update as you type
- ✅ No perceptible delay

---

## Terminal Compatibility

### Windows Terminal

**Expected**:
- ✅ Rounded borders display correctly
- ✅ True color support (24-bit)
- ✅ All keyboard shortcuts work
- ✅ Smooth resizing
- ✅ Alternate screen buffer works (terminal restores on quit)

---

### PowerShell / CMD

**Expected**:
- ✅ Basic functionality works
- ⚠️ Rounded borders may fall back to ASCII: `+-|`
- ✅ 256-color support (good enough for theme)
- ✅ Keyboard shortcuts work

---

### iTerm2 (macOS)

**Expected**:
- ✅ Excellent rendering (best terminal for macOS)
- ✅ Rounded borders display correctly
- ✅ True color support
- ✅ All shortcuts work
- ✅ Smooth performance

---

### Terminal.app (macOS)

**Expected**:
- ✅ Good rendering
- ✅ Rounded borders work
- ✅ 256-color support
- ✅ All shortcuts work

---

### gnome-terminal (Linux)

**Expected**:
- ✅ Good rendering
- ✅ Rounded borders work
- ✅ True color support
- ✅ Keyboard shortcuts work

---

### Alacritty (Cross-Platform)

**Expected**:
- ✅ Excellent performance (GPU-accelerated)
- ✅ Rounded borders work
- ✅ True color support
- ✅ Fastest rendering

---

## Conclusion

Use these expected results as your reference when testing. Any deviation from these expected behaviors should be documented as a bug or issue in your test report.

For questions or clarifications on expected behavior, refer to the tui-tview-implementation spec's requirements and design documents.

---

**Document Version**: 1.0
**Last Updated**: (To be filled during testing)
