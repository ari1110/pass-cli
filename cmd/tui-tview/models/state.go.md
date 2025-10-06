# models/state.go

## Purpose
Central application state management with thread-safe access to vault data and UI state. This is the single source of truth for the entire TUI.

## Responsibilities

1. **Vault service wrapper**: Thread-safe access to vault operations
2. **Credential storage**: In-memory cache of loaded credentials
3. **UI state tracking**: Current view, selected items, navigation state
4. **Component management**: Store references to UI components (single instances)
5. **State change notifications**: Notify UI components when data changes

## CRITICAL: Avoid Mutex Deadlocks

**The Problem**: Previous implementation had deadlocks from calling callbacks while holding locks.

**The Solution**:
1. **Never call callbacks while holding a lock**
2. **Release lock BEFORE notifying**
3. **Use explicit notification helpers**

## Dependencies

### Internal Dependencies
- `pass-cli/internal/vault` - Vault service for credential operations

### External Dependencies
- `sync` - RWMutex for thread-safe access
- `github.com/rivo/tview` - Component types

### Component Dependencies
- `pass-cli/cmd/tui-tview/components` - For storing component instances

## Key Types

### `AppState`
**Purpose**: Holds all application state

**Fields**:
```go
type AppState struct {
    // Concurrency control
    mu sync.RWMutex  // Protects all fields below

    // Vault service
    vault *vault.Vault

    // Credential data
    credentials []vault.Credential
    categories  []string

    // Current selections
    selectedCategory   string
    selectedCredential *vault.Credential

    // UI components (single instances, created once)
    sidebar    *tview.TreeView
    table      *tview.Table
    detailView *tview.TextView
    statusBar  *tview.TextView

    // Notification callbacks
    onCredentialsChanged func()  // Called when credentials are loaded/modified
    onSelectionChanged   func()  // Called when selection changes
    onError              func(error)  // Called when errors occur
}
```

## Key Functions

### Constructor

#### `NewAppState(vaultService *vault.Vault) *AppState`
**Purpose**: Create new AppState with vault service

**Steps**:
1. Initialize AppState struct
2. Store vault service reference
3. Initialize empty credential slice
4. Return state (components added later)

### Data Access (Read Operations)

#### `GetCredentials() []vault.Credential`
**Purpose**: Thread-safe read of credentials

**Locking**: Read lock (RLock)

```go
func (s *AppState) GetCredentials() []vault.Credential {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.credentials
}
```

#### `GetSelectedCredential() *vault.Credential`
**Purpose**: Get currently selected credential

**Returns**: Copy of selected credential (or nil)

#### `GetCategories() []string`
**Purpose**: Get list of unique categories

**Returns**: Sorted category list

### Data Mutations (Write Operations)

#### `LoadCredentials() error`
**Purpose**: Load all credentials from vault

**CRITICAL - Deadlock Prevention**:
```go
func (s *AppState) LoadCredentials() error {
    s.mu.Lock()

    // Load credentials from vault
    creds, err := s.vault.GetCredentials()
    if err != nil {
        wrappedErr := fmt.Errorf("failed to load: %w", err)
        s.mu.Unlock()  // ✅ RELEASE LOCK FIRST
        s.notifyError(wrappedErr)  // ✅ THEN notify
        return wrappedErr
    }

    // Update state
    s.credentials = creds
    s.updateCategories()  // Internal helper, safe to call while locked

    s.mu.Unlock()  // ✅ RELEASE LOCK
    s.notifyCredentialsChanged()  // ✅ THEN notify

    return nil
}
```

#### `AddCredential(service, username, password string) error`
**Purpose**: Add new credential to vault

**Pattern**: Same as LoadCredentials - lock, mutate, unlock, notify

#### `UpdateCredential(id, service, username, password string) error`
**Purpose**: Update existing credential

**Pattern**: Lock → Update → Unlock → Notify

#### `DeleteCredential(id string) error`
**Purpose**: Delete credential from vault

**Pattern**: Lock → Delete → Unlock → Notify

### Selection Management

#### `SetSelectedCategory(category string)`
**Purpose**: Update selected category

**Effect**: Triggers selection changed notification

#### `SetSelectedCredential(credential *vault.Credential)`
**Purpose**: Update selected credential

**Effect**: Triggers selection changed notification

### Component Management

#### `SetSidebar(sidebar *tview.TreeView)`
**Purpose**: Store sidebar component reference

**Why**: Components created once, stored here, reused

#### `GetSidebar() *tview.TreeView`
**Purpose**: Retrieve sidebar component

**Why**: Layout manager uses this instead of creating new instance

#### `SetTable(table *tview.Table)` / `GetTable()`
#### `SetDetailView(view *tview.TextView)` / `GetDetailView()`
#### `SetStatusBar(bar *tview.TextView)` / `GetStatusBar()`

**Pattern**: Same as sidebar - store and retrieve single instances

### Notification Management

#### `SetOnCredentialsChanged(callback func())`
**Purpose**: Register callback for credential changes

**When called**: After LoadCredentials(), AddCredential(), etc.

#### `SetOnSelectionChanged(callback func())`
**Purpose**: Register callback for selection changes

**When called**: After SetSelectedCategory(), SetSelectedCredential()

#### `SetOnError(callback func(error))`
**Purpose**: Register callback for errors

**When called**: When operations fail

#### Private Notification Helpers

```go
// These are ALWAYS called AFTER releasing locks

func (s *AppState) notifyCredentialsChanged() {
    if s.onCredentialsChanged != nil {
        s.onCredentialsChanged()
    }
}

func (s *AppState) notifySelectionChanged() {
    if s.onSelectionChanged != nil {
        s.onSelectionChanged()
    }
}

func (s *AppState) notifyError(err error) {
    if s.onError != nil {
        s.onError(err)
    }
}
```

## Example Usage

### Initialization
```go
// In main.go:
vaultService, _ := vault.NewVault()
appState := models.NewAppState(vaultService)

// Load credentials
if err := appState.LoadCredentials(); err != nil {
    return err
}

// Create and store components
sidebar := components.NewSidebar(appState)
appState.SetSidebar(sidebar)

table := components.NewTable(appState)
appState.SetTable(table)

// Register callbacks
appState.SetOnCredentialsChanged(func() {
    // Refresh all UI components
    sidebar.Refresh()
    table.Refresh()
})

appState.SetOnSelectionChanged(func() {
    // Update detail view
    detailView.Refresh()
})
```

### Adding a Credential
```go
// In form submit handler:
err := appState.AddCredential(service, username, password)
if err != nil {
    // Error callback will be invoked automatically
    return
}

// onCredentialsChanged callback will be invoked automatically
// UI will refresh via the callback
```

## Deadlock Prevention Checklist

✅ **Do**:
- Lock for the shortest time possible
- Release lock BEFORE calling callbacks
- Use defer for read locks only (simple getters)
- Keep locked sections simple (no function calls that might block)

❌ **Don't**:
- Call any callback while holding a lock
- Call external functions while locked (they might try to acquire lock)
- Use defer for write locks if you need to call callbacks before return
- Hold lock across I/O operations

## Pattern Examples

### ✅ Correct Pattern
```go
func (s *AppState) UpdateData() error {
    s.mu.Lock()

    // Do work
    s.data = newData

    s.mu.Unlock()  // Release FIRST
    s.notifyChanged()  // Notify AFTER

    return nil
}
```

### ❌ Wrong Pattern (Deadlock Risk)
```go
func (s *AppState) UpdateData() error {
    s.mu.Lock()
    defer s.mu.Unlock()  // Lock held until return

    s.data = newData
    s.notifyChanged()  // ❌ Called while locked!

    return nil
}
```

## Testing Considerations

- **Mock vault service**: Inject mock for testing state operations
- **Test callbacks**: Verify callbacks are invoked after operations
- **Test thread safety**: Use go routines to test concurrent access
- **Test deadlock prevention**: Verify locks are released before callbacks

## Future Enhancements

- **Undo/Redo**: Track state changes for undo functionality
- **Search state**: Store search query and results
- **Filter state**: Store active filters
- **Sort state**: Store current sort order
- **Session persistence**: Save state on exit, restore on launch
