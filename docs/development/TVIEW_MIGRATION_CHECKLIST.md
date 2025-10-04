# tview Migration Visual Regression Checklist

This checklist ensures visual parity and compatibility between the original Bubble Tea implementation and the new tview implementation.

## Overview

**Migration Goal**: Maintain identical visual appearance and functionality while migrating from Bubble Tea to tview framework.

**Testing Approach**: Manual visual comparison across multiple terminal emulators and operating systems.

---

## Terminal Compatibility Matrix

Test the TUI application on the following terminal emulators to ensure consistent rendering:

### Windows
- [ ] **Windows Terminal** (latest version)
  - [ ] Color rendering (256-color support)
  - [ ] Unicode characters (icons, borders)
  - [ ] Box drawing characters
  - [ ] Resize handling
  - [ ] Mouse support (if applicable)

- [ ] **ConEmu**
  - [ ] Color rendering
  - [ ] Unicode characters
  - [ ] Border styles

- [ ] **Command Prompt (cmd.exe)**
  - [ ] Fallback rendering for limited color support
  - [ ] Basic functionality

### macOS
- [ ] **iTerm2** (latest version)
  - [ ] True color support
  - [ ] Unicode rendering
  - [ ] Border styles (rounded borders)
  - [ ] Resize handling
  - [ ] Split pane rendering

- [ ] **macOS Terminal.app**
  - [ ] Color rendering
  - [ ] Unicode characters
  - [ ] Border styles
  - [ ] Resize handling

### Linux
- [ ] **gnome-terminal**
  - [ ] Color rendering
  - [ ] Unicode support
  - [ ] Border rendering
  - [ ] Resize handling

- [ ] **xterm**
  - [ ] 256-color support
  - [ ] Basic rendering

- [ ] **Konsole (KDE)**
  - [ ] Color rendering
  - [ ] Unicode support

- [ ] **Alacritty**
  - [ ] True color support
  - [ ] Performance with large credential lists

---

## Visual Components Checklist

Compare each component visually between Bubble Tea and tview implementations:

### Status Bar
- [ ] **Layout**
  - [ ] Status bar positioned at bottom
  - [ ] Full width spanning terminal
  - [ ] Correct height (1 line)

- [ ] **Content**
  - [ ] Keychain indicator (🔓/🔒) displays correctly
  - [ ] Credential count displays with correct pluralization
  - [ ] Current view name displays
  - [ ] Keyboard shortcuts display correctly
  - [ ] Text truncation works for narrow terminals

- [ ] **Styling**
  - [ ] Background color matches theme
  - [ ] Text color readable and matches theme
  - [ ] Separators between sections render correctly

### Sidebar
- [ ] **Layout**
  - [ ] Correct width (25 columns default)
  - [ ] Rounded border renders correctly
  - [ ] Border color changes when focused (active/inactive)
  - [ ] Proper height calculation

- [ ] **Content**
  - [ ] "Categories" header displays
  - [ ] Category icons render correctly (☁️ 📦 🗄️ 🔧 💬 💳 🤖)
  - [ ] Category names display
  - [ ] Credential counts per category accurate
  - [ ] Vault statistics section displays
  - [ ] Tree structure renders with proper indentation

- [ ] **Interactions**
  - [ ] Expand/collapse icons (▶/▼) work correctly
  - [ ] Selection highlighting visible
  - [ ] Scroll behavior with many categories
  - [ ] Focus border color change

### Main Panel (List View)
- [ ] **Layout**
  - [ ] Takes remaining horizontal space
  - [ ] Rounded border renders correctly
  - [ ] Border color changes when focused
  - [ ] Breadcrumb displays above panel

- [ ] **Content**
  - [ ] Table headers (Service, Username, Category) display
  - [ ] Credentials listed in rows
  - [ ] Category badges display correctly
  - [ ] Row selection highlighting visible
  - [ ] Empty state message displays when no credentials

- [ ] **Search**
  - [ ] Search input field displays at top
  - [ ] Search prompt visible
  - [ ] Real-time filtering works
  - [ ] Filtered results display correctly

- [ ] **Interactions**
  - [ ] Up/down navigation (j/k, arrow keys)
  - [ ] Selection follows cursor
  - [ ] Scroll behavior with many credentials
  - [ ] Enter key opens detail view

### Main Panel (Detail View)
- [ ] **Layout**
  - [ ] Rounded border renders correctly
  - [ ] Proper spacing between fields
  - [ ] Breadcrumb shows navigation path

- [ ] **Content**
  - [ ] Service name displays prominently
  - [ ] Username field displays
  - [ ] Password field displays (masked by default)
  - [ ] Password masking (••••••) works
  - [ ] Password reveal toggle works (m key)
  - [ ] Notes section displays
  - [ ] Timestamps display correctly
  - [ ] Field labels bold/emphasized

- [ ] **Interactions**
  - [ ] 'm' key toggles password visibility
  - [ ] 'c' key copies password
  - [ ] 'e' key enters edit mode
  - [ ] 'd' key shows delete confirmation
  - [ ] ESC returns to list view

### Metadata Panel
- [ ] **Layout**
  - [ ] Displays on right side (≥120 columns)
  - [ ] Correct width (25 columns)
  - [ ] Rounded border renders correctly
  - [ ] Border color changes when focused

- [ ] **Content**
  - [ ] "Metadata" header displays
  - [ ] Created timestamp displays
  - [ ] Updated timestamp displays
  - [ ] Last used timestamp displays
  - [ ] Usage count displays
  - [ ] Field labels and values aligned

- [ ] **Visibility**
  - [ ] Hidden when terminal < 120 columns
  - [ ] 'i' key toggles visibility
  - [ ] Shows/hides smoothly

### Forms (Add/Edit)
- [ ] **Layout**
  - [ ] Form centered on screen
  - [ ] Border renders correctly
  - [ ] Title displays ("Add Credential" / "Edit Credential")

- [ ] **Fields**
  - [ ] Service input field
  - [ ] Username input field
  - [ ] Password input field
  - [ ] Notes input field (multi-line)
  - [ ] Field labels display
  - [ ] Input focus indicator visible

- [ ] **Interactions**
  - [ ] Tab navigation between fields
  - [ ] Text input works in all fields
  - [ ] Enter submits form (when on button)
  - [ ] ESC shows discard confirmation (if changes)
  - [ ] Generate password button works

### Confirmation Dialogs
- [ ] **Layout**
  - [ ] Modal overlay/centered
  - [ ] Border renders correctly
  - [ ] Proper sizing for message

- [ ] **Content**
  - [ ] Confirmation message displays clearly
  - [ ] Service name shows in delete confirmations
  - [ ] Yes/No buttons visible
  - [ ] Keyboard shortcuts shown (y/n)

- [ ] **Interactions**
  - [ ] 'y' key confirms action
  - [ ] 'n' key cancels
  - [ ] ESC cancels
  - [ ] Button focus indicators work

### Help View
- [ ] **Layout**
  - [ ] Full screen overlay
  - [ ] Border renders correctly
  - [ ] Scrollable content area

- [ ] **Content**
  - [ ] Title "Help" displays
  - [ ] All keyboard shortcuts listed
  - [ ] Descriptions clear and readable
  - [ ] Section headers visible
  - [ ] Grouped by category (Navigation, Actions, etc.)

- [ ] **Interactions**
  - [ ] Scrolling works (j/k, arrows, page up/down)
  - [ ] '?' or F1 opens help
  - [ ] ESC or '?' closes help

### Command Bar
- [ ] **Layout**
  - [ ] Displays above status bar when activated
  - [ ] Input field with ':' prompt
  - [ ] Full width

- [ ] **Content**
  - [ ] ':' prompt visible
  - [ ] Placeholder text displays
  - [ ] Error messages display in red
  - [ ] Command history works (up/down arrows)

- [ ] **Interactions**
  - [ ] ':' key activates command bar
  - [ ] Text input works
  - [ ] Enter executes command
  - [ ] ESC closes command bar
  - [ ] Up/down navigate history

### Breadcrumb
- [ ] **Layout**
  - [ ] Displays above main panel (detail view only)
  - [ ] Single line height
  - [ ] Proper spacing

- [ ] **Content**
  - [ ] Navigation path displays (e.g., "All > Cloud > aws-prod")
  - [ ] Separator characters render correctly (>)
  - [ ] Text color matches theme

### Process Panel
- [ ] **Layout**
  - [ ] Displays above command bar when active
  - [ ] Correct height (3 lines)
  - [ ] No border (plain background)

- [ ] **Content**
  - [ ] Process status displays
  - [ ] Progress indicators render
  - [ ] Status icons display (⏳ ✓ ✗)
  - [ ] Process names and messages visible

---

## Responsive Layout Testing

Test layout behavior at different terminal sizes:

### Breakpoints
- [ ] **Small (< 80 columns)**
  - [ ] Sidebar hidden automatically
  - [ ] Metadata hidden automatically
  - [ ] Main panel takes full width
  - [ ] Status bar truncates gracefully

- [ ] **Medium (80-119 columns)**
  - [ ] Sidebar visible (if toggled on)
  - [ ] Metadata hidden
  - [ ] Main panel shares space with sidebar
  - [ ] Proper proportions (1:2 ratio)

- [ ] **Large (≥ 120 columns)**
  - [ ] All panels visible (if toggled on)
  - [ ] Sidebar, main, metadata proportions correct (1:2:1)
  - [ ] No layout overflow

### Dynamic Resizing
- [ ] **Resize while running**
  - [ ] Layout recalculates smoothly
  - [ ] No visual glitches
  - [ ] Panels show/hide at breakpoints
  - [ ] Selected credential preserved
  - [ ] View state maintained

- [ ] **Too small warning**
  - [ ] Warning message displays when terminal < minimum
  - [ ] Message shows current and minimum size
  - [ ] Readable instructions

---

## Color and Styling Testing

### Theme Colors
- [ ] **Primary colors**
  - [ ] Primary color (blue) renders correctly
  - [ ] Secondary color (purple) renders correctly
  - [ ] Accent color (cyan) renders correctly

- [ ] **Semantic colors**
  - [ ] Success green displays correctly
  - [ ] Warning yellow displays correctly
  - [ ] Error red displays correctly
  - [ ] Muted gray for secondary text

- [ ] **Panel borders**
  - [ ] Active border color (bright blue)
  - [ ] Inactive border color (dim gray)
  - [ ] Border color transitions on focus change

### Text Styles
- [ ] **Emphasis**
  - [ ] Bold text renders correctly
  - [ ] Italic text renders correctly (if used)
  - [ ] Dimmed/muted text visible but subdued

- [ ] **Syntax highlighting**
  - [ ] Field labels styled consistently
  - [ ] Values styled consistently
  - [ ] Timestamps styled with muted color

---

## Keyboard Navigation Testing

### Global Shortcuts
- [ ] 'q' quits from list/detail/help views
- [ ] '?' or F1 opens help
- [ ] '/' activates search (list view)
- [ ] ':' opens command bar
- [ ] 's' toggles sidebar
- [ ] 'i' toggles metadata panel (detail view)
- [ ] 'p' toggles process panel
- [ ] 'f' toggles footer (status bar)
- [ ] Tab switches panel focus
- [ ] Shift+Tab switches panel focus backwards
- [ ] ESC navigates back / closes dialogs

### View-Specific Shortcuts
- [ ] **List View**
  - [ ] 'a' opens add form
  - [ ] 'j' / Down arrow moves selection down
  - [ ] 'k' / Up arrow moves selection up
  - [ ] Enter opens detail view

- [ ] **Detail View**
  - [ ] 'e' enters edit mode
  - [ ] 'd' shows delete confirmation
  - [ ] 'm' toggles password mask
  - [ ] 'c' copies password to clipboard

- [ ] **Forms**
  - [ ] Tab moves to next field
  - [ ] Shift+Tab moves to previous field
  - [ ] Enter submits (from button)
  - [ ] ESC shows discard confirmation (if changes)

---

## Performance Testing

- [ ] **Startup time**
  - [ ] Application starts in < 200ms
  - [ ] Initial render completes quickly

- [ ] **Rendering performance**
  - [ ] No visible lag when switching views
  - [ ] Smooth scrolling with 100+ credentials
  - [ ] Resize recalculation < 50ms

- [ ] **Memory usage**
  - [ ] Memory usage ≤ 50MB during normal operation
  - [ ] No memory leaks during extended use

---

## Comparison Screenshots

For visual regression testing, capture screenshots of key views in both implementations:

### Required Screenshots
1. [ ] List view with sidebar and status bar (140x40 terminal)
2. [ ] Detail view with all panels visible (140x40 terminal)
3. [ ] Add form centered
4. [ ] Delete confirmation dialog
5. [ ] Help screen
6. [ ] List view at medium breakpoint (100x30)
7. [ ] List view at small breakpoint (70x25)
8. [ ] Search active with filtered results
9. [ ] Command bar active
10. [ ] Too small warning message

### Screenshot Comparison Checklist
For each screenshot pair (Bubble Tea vs tview):
- [ ] Layout dimensions identical
- [ ] Border styles match
- [ ] Colors match (or are improved)
- [ ] Text positioning identical
- [ ] Icons and symbols render identically
- [ ] Spacing and padding consistent

---

## Known Differences (Acceptable)

Document any intentional differences from the original Bubble Tea implementation:

1. **tview-specific improvements**:
   - Potentially improved performance with tview's rendering
   - Better terminal compatibility with tview's tcell backend
   - Enhanced Unicode support

2. **Minor acceptable differences**:
   - Slight timing differences in animations (if any)
   - Backend-specific rendering optimizations
   - Terminal-specific rendering differences (expected)

---

## Final Validation

- [ ] All checklist items completed
- [ ] All terminal emulators tested
- [ ] All responsive breakpoints verified
- [ ] All keyboard shortcuts functional
- [ ] Visual parity confirmed with screenshots
- [ ] No regressions in functionality
- [ ] Performance meets or exceeds requirements
- [ ] Documentation updated to reflect any changes

---

## Sign-off

**Tester**: ________________
**Date**: ________________
**Migration Version**: ________________

**Notes**:
