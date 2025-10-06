# events/focus.go

## Purpose
Focus management utilities and helpers. Provides convenient functions for setting focus with proper state updates.

## Responsibilities

1. **Focus helpers**: Convenient functions for common focus operations
2. **Focus validation**: Ensure focus target is valid and visible
3. **Focus callbacks**: Trigger appropriate callbacks on focus changes
4. **Focus restoration**: Remember and restore previous focus

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For NavigationState, AppState

### External Dependencies
- `github.com/rivo/tview` - Application

## Key Functions

### Focus Setting

#### `SetFocusToComponent(app *tview.Application, nav *models.NavigationState, target models.FocusableComponent) error`
**Purpose**: Safely set focus to a component with validation

**Steps**:
1. Validate target component is registered
2. Check if component is currently visible
3. Update navigation state
4. Call app.SetFocus() on primitive
5. Trigger focus changed callback

**Returns**: Error if component not registered or not visible

#### `RestorePreviousFocus(app *tview.Application, nav *models.NavigationState) error`
**Purpose**: Return focus to previous component

**Use case**: After closing modal, restore focus

**Implementation**: Uses navigation state's focus history

### Focus Validation

#### `IsComponentVisible(component models.FocusableComponent, layoutMode layout.LayoutMode) bool`
**Purpose**: Check if component is visible in current layout

**Logic**:
```go
func IsComponentVisible(component models.FocusableComponent, mode layout.LayoutMode) bool {
    switch component {
    case models.FocusSidebar:
        return mode != layout.LayoutSmall  // Hidden in small mode

    case models.FocusDetail:
        return mode == layout.LayoutLarge  // Only in large mode

    case models.FocusTable:
        return true  // Always visible

    case models.FocusStatusBar:
        return true  // Always visible

    default:
        return false
    }
}
```

#### `GetNextVisibleComponent(current models.FocusableComponent, mode layout.LayoutMode) models.FocusableComponent`
**Purpose**: Find next visible component for Tab navigation

**Logic**: Skip hidden components in tab order

### Focus Callbacks

#### `OnFocusChanged(focus models.FocusableComponent, statusBar *components.StatusBar)`
**Purpose**: Standard callback for focus changes

**Actions**:
- Update status bar with context-appropriate shortcuts
- Update component highlights
- Log focus change (if debug mode)

### Focus Helpers for Common Operations

#### `FocusOnNewCredential(app *tview.Application, appState *models.AppState, nav *models.NavigationState, credentialID string)`
**Purpose**: Focus on newly added credential in table

**Steps**:
1. Find credential in table
2. Set table selection to that row
3. Set focus to table
4. Update navigation state

#### `FocusOnFirstCredential(app *tview.Application, nav *models.NavigationState)`
**Purpose**: Focus on first credential in table

**Use case**: After loading vault, start at first item

## Example Structure

```go
package events

import (
    "fmt"

    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/components"
    "pass-cli/cmd/tui-tview/layout"
    "pass-cli/cmd/tui-tview/models"
)

// SetFocusToComponent safely sets focus to a component
func SetFocusToComponent(app *tview.Application, nav *models.NavigationState, target models.FocusableComponent) error {
    // Get current layout mode (would need to be passed in or accessed from layout manager)
    // For now, assume we have it
    if !IsComponentVisible(target, getCurrentLayoutMode()) {
        return fmt.Errorf("component %d is not visible", target)
    }

    return nav.SetFocus(target)
}

// IsComponentVisible checks if component is visible in current layout
func IsComponentVisible(component models.FocusableComponent, mode layout.LayoutMode) bool {
    switch component {
    case models.FocusSidebar:
        return mode != layout.LayoutSmall

    case models.FocusDetail:
        return mode == layout.LayoutLarge

    case models.FocusTable, models.FocusStatusBar:
        return true

    default:
        return false
    }
}

// GetNextVisibleComponent finds next visible component in tab order
func GetNextVisibleComponent(current models.FocusableComponent, mode layout.LayoutMode) models.FocusableComponent {
    order := []models.FocusableComponent{
        models.FocusSidebar,
        models.FocusTable,
        models.FocusDetail,
    }

    // Find current position
    currentIndex := -1
    for i, comp := range order {
        if comp == current {
            currentIndex = i
            break
        }
    }

    // Cycle through order, skipping hidden components
    for i := 1; i <= len(order); i++ {
        nextIndex := (currentIndex + i) % len(order)
        next := order[nextIndex]

        if IsComponentVisible(next, mode) {
            return next
        }
    }

    // Fallback: return current
    return current
}

// GetPreviousVisibleComponent finds previous visible component
func GetPreviousVisibleComponent(current models.FocusableComponent, mode layout.LayoutMode) models.FocusableComponent {
    order := []models.FocusableComponent{
        models.FocusSidebar,
        models.FocusTable,
        models.FocusDetail,
    }

    currentIndex := -1
    for i, comp := range order {
        if comp == current {
            currentIndex = i
            break
        }
    }

    for i := 1; i <= len(order); i++ {
        prevIndex := (currentIndex - i + len(order)) % len(order)
        prev := order[prevIndex]

        if IsComponentVisible(prev, mode) {
            return prev
        }
    }

    return current
}

// OnFocusChanged standard callback for focus changes
func OnFocusChanged(focus models.FocusableComponent, statusBar *components.StatusBar) {
    statusBar.UpdateForFocus(focus)
}

// FocusOnNewCredential focuses on newly added credential
func FocusOnNewCredential(app *tview.Application, appState *models.AppState, nav *models.NavigationState, credentialID string) {
    // Find credential in table
    table := appState.GetTable()
    if table == nil {
        return
    }

    // Set table selection (would need table method to find row by ID)
    // table.SelectCredential(credentialID)

    // Set focus
    nav.SetFocus(models.FocusTable)
}

// FocusOnFirstCredential focuses on first credential
func FocusOnFirstCredential(app *tview.Application, nav *models.NavigationState) {
    nav.SetFocus(models.FocusTable)
}

// RestoreFocusAfterModal restores focus after closing modal
func RestoreFocusAfterModal(app *tview.Application, nav *models.NavigationState) {
    // Get previous focus from navigation state
    // For now, default to table
    nav.SetFocus(models.FocusTable)
}
```

## Usage Examples

### Setting Focus with Validation
```go
// In event handler:
err := focus.SetFocusToComponent(app, nav, models.FocusSidebar)
if err != nil {
    // Component not visible, try next
    focus.SetFocusToComponent(app, nav, models.FocusTable)
}
```

### Tab Navigation with Skip Hidden
```go
// In Tab key handler:
current := nav.GetCurrentFocus()
layoutMode := layoutManager.GetCurrentMode()
next := focus.GetNextVisibleComponent(current, layoutMode)
nav.SetFocus(next)
```

### Focus on New Credential
```go
// After adding credential:
err := appState.AddCredential(service, username, password)
if err != nil {
    return err
}

// Focus on the new credential
credID := appState.GetLastAddedCredentialID()
focus.FocusOnNewCredential(app, appState, nav, credID)
```

## Focus Order Strategy

### Default Tab Order
1. **Sidebar** (if visible)
2. **Table**
3. **Detail** (if visible)

### Skip Hidden Components
- Small layout: Skip sidebar
- Medium layout: Skip detail
- Large layout: All visible

### Reverse Order (Shift+Tab)
Same order, reversed

## Focus States

### Focused Component
- **Visual indicator**: Highlighted border or background
- **Keyboard input**: Receives key events
- **Status bar**: Shows relevant shortcuts

### Unfocused Component
- **Visual indicator**: Dimmed border or background
- **Keyboard input**: Ignored (unless global shortcut)
- **Status bar**: N/A

## Integration with Layout Manager

Focus management must be aware of layout mode:

```go
// When layout mode changes:
layoutManager.OnModeChanged(func(newMode layout.LayoutMode) {
    current := nav.GetCurrentFocus()

    // If current focus is now hidden, move to visible component
    if !focus.IsComponentVisible(current, newMode) {
        next := focus.GetNextVisibleComponent(current, newMode)
        nav.SetFocus(next)
    }
})
```

## Testing Considerations

- **Test visibility**: Verify hidden components are skipped
- **Test cycling**: Verify Tab/Shift+Tab work correctly
- **Test restoration**: Verify focus restores after modal
- **Test validation**: Verify errors on invalid focus targets
- **Mock layout mode**: Use test values for layout mode

## Future Enhancements

- **Focus hints**: Visual indicators for tab order
- **Focus history**: Remember last N focus positions
- **Focus search**: Find component by name
- **Focus logging**: Debug mode for focus tracking
- **Auto-focus**: Focus on most relevant component automatically
- **Focus traps**: Prevent focus from leaving modal
