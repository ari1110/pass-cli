# models/navigation.go

## Purpose
Manages navigation state including focus tracking, breadcrumbs, and view history. Separate from main AppState to keep concerns focused.

## Responsibilities

1. **Focus tracking**: Which component currently has focus
2. **View history**: Stack of previous views for back navigation
3. **Breadcrumb management**: Current navigation path display
4. **Focus cycling**: Move focus between components

## Dependencies

### Internal Dependencies
- None - This is a standalone utility

### External Dependencies
- `github.com/rivo/tview` - For Primitive type

## Key Types

### `FocusableComponent`
**Purpose**: Enum for components that can receive focus

```go
type FocusableComponent int

const (
    FocusSidebar FocusableComponent = iota
    FocusTable
    FocusDetail
    FocusForm
    FocusStatusBar
)
```

### `NavigationState`
**Purpose**: Tracks navigation and focus state

**Fields**:
```go
type NavigationState struct {
    // Current focus
    currentFocus FocusableComponent

    // Component references (for focus switching)
    components map[FocusableComponent]tview.Primitive

    // View history (for back navigation)
    viewHistory []string  // Stack of view identifiers

    // Breadcrumb path
    breadcrumbs []string  // ["Home", "AWS", "Credentials"]

    // Callback when focus changes
    onFocusChanged func(FocusableComponent)
}
```

## Key Functions

### Constructor

#### `NewNavigationState() *NavigationState`
**Purpose**: Create new navigation state

**Initial state**: Focus on sidebar, empty history

### Component Registration

#### `RegisterComponent(focus FocusableComponent, primitive tview.Primitive)`
**Purpose**: Register a component for focus management

**Usage**:
```go
nav.RegisterComponent(FocusSidebar, sidebar)
nav.RegisterComponent(FocusTable, table)
nav.RegisterComponent(FocusDetail, detailView)
```

### Focus Management

#### `GetCurrentFocus() FocusableComponent`
**Purpose**: Get currently focused component

**Returns**: FocusableComponent enum value

#### `SetFocus(focus FocusableComponent) error`
**Purpose**: Change focus to specified component

**Steps**:
1. Validate component is registered
2. Update currentFocus
3. Call app.SetFocus() on the primitive
4. Invoke onFocusChanged callback

**Returns**: Error if component not registered

#### `CycleFocus()`
**Purpose**: Move focus to next component in tab order

**Tab order**: Sidebar → Table → Detail → (back to Sidebar)

**Skips**: Components not currently visible

#### `CycleFocusReverse()`
**Purpose**: Move focus to previous component (Shift+Tab)

**Tab order**: Sidebar ← Table ← Detail ← (back to Detail)

### View History

#### `PushView(viewID string)`
**Purpose**: Add current view to history stack

**Usage**:
```go
nav.PushView("credential-list")
// User navigates to detail view
nav.PushView("credential-detail")
```

#### `PopView() (string, error)`
**Purpose**: Go back to previous view

**Returns**: Previous view ID, or error if history is empty

**Usage**:
```go
viewID, err := nav.PopView()
if err == nil {
    // Navigate to viewID
}
```

#### `GetViewHistory() []string`
**Purpose**: Get full view history (for debugging)

**Returns**: Copy of history stack

### Breadcrumb Management

#### `SetBreadcrumbs(path ...string)`
**Purpose**: Update breadcrumb path

**Usage**:
```go
nav.SetBreadcrumbs("Home")
nav.SetBreadcrumbs("Home", "AWS")
nav.SetBreadcrumbs("Home", "AWS", "Production Database")
```

#### `GetBreadcrumbs() []string`
**Purpose**: Get current breadcrumb path

**Returns**: Copy of breadcrumb slice

**For display**:
```go
path := nav.GetBreadcrumbs()
display := strings.Join(path, " > ")
// "Home > AWS > Production Database"
```

### Callback Management

#### `SetOnFocusChanged(callback func(FocusableComponent))`
**Purpose**: Register callback for focus changes

**When called**: After SetFocus() or CycleFocus()

**Usage**:
```go
nav.SetOnFocusChanged(func(focus FocusableComponent) {
    // Update status bar to show context-appropriate shortcuts
    statusBar.UpdateForFocus(focus)
})
```

## Example Usage

### Initialization
```go
// Create navigation state
nav := models.NewNavigationState()

// Register components
nav.RegisterComponent(models.FocusSidebar, sidebar)
nav.RegisterComponent(models.FocusTable, table)
nav.RegisterComponent(models.FocusDetail, detailView)

// Setup focus change callback
nav.SetOnFocusChanged(func(focus models.FocusableComponent) {
    switch focus {
    case models.FocusSidebar:
        statusBar.SetText("Arrow keys: navigate | Enter: select")
    case models.FocusTable:
        statusBar.SetText("Arrow keys: navigate | Enter: view | n: new")
    case models.FocusDetail:
        statusBar.SetText("Esc: back | e: edit | d: delete")
    }
})

// Set initial focus
nav.SetFocus(models.FocusSidebar)
```

### Keyboard Navigation
```go
// In event handler:
app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    switch event.Key() {
    case tcell.KeyTab:
        nav.CycleFocus()
        return nil
    case tcell.KeyBacktab:  // Shift+Tab
        nav.CycleFocusReverse()
        return nil
    }
    return event
})
```

### View Navigation
```go
// When opening detail view:
nav.PushView("list")
nav.SetBreadcrumbs("Home", category, credentialName)
// Show detail view

// When pressing Escape (back):
if prevView, err := nav.PopView(); err == nil {
    // Navigate back to previous view
}
```

## Integration with AppState

**Separate but related**:
- **AppState**: Holds *what* is selected (data)
- **NavigationState**: Holds *where* focus is (UI)

```go
// AppState tracks data:
appState.SetSelectedCredential(cred)

// NavigationState tracks focus:
nav.SetFocus(FocusDetail)
```

## Tab Order Strategy

### Normal Flow (Tab)
1. **Sidebar** → Select category
2. **Table** → Select credential
3. **Detail** → View/edit details
4. **Back to Sidebar**

### When Modals Open
- Modal forms capture focus
- Tab within form only
- Escape returns focus to previous component

## Testing Considerations

- **Test focus changes**: Verify callbacks are invoked
- **Test tab order**: Verify correct cycling
- **Test view history**: Push/pop operations
- **Mock components**: Use test primitives for registration

## Future Enhancements

- **Focus memory**: Remember last focused element per view
- **Focus hints**: Visual indicators for focused component
- **Keyboard shortcuts overlay**: Show shortcuts for current focus
- **Search mode**: Special focus state for search input
- **Modal stack**: Track nested modals for proper focus restoration
