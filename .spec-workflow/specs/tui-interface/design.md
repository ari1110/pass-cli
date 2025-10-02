# Design Document: TUI Interface

## Overview

The TUI Interface feature adds a Terminal User Interface mode to Pass-CLI, providing an interactive, visual alternative to the existing command-line interface. The TUI will be built using the Bubble Tea framework (charmbracelet/bubbletea) following the Elm architecture pattern (Model-Update-View). The design ensures complete isolation from existing CLI commands while sharing the same underlying VaultService for all credential operations.

**Key Design Principles**:
- **Zero Impact on CLI**: TUI is an additive feature; removing the TUI package would not affect any CLI functionality
- **Shared Business Logic**: Both TUI and CLI use the same VaultService API, ensuring consistent behavior
- **Component-Based Architecture**: Each screen/view is an isolated, testable component
- **Keyboard-First UX**: All interactions via keyboard with vim-style shortcuts and intuitive navigation

## Steering Document Alignment

### Technical Standards (tech.md)

**Language and Tooling**:
- Implementation in Go 1.25+ (matching existing codebase)
- Uses Bubble Tea v0.25+ for TUI framework (proven, widely adopted)
- Follows existing dependency management with Go modules
- Maintains build system compatibility (Makefile, cross-compilation)

**Architecture Alignment**:
- **Layered Design**: TUI layer sits alongside CLI layer, both consuming VaultService
- **Service Layer Reuse**: Uses existing `internal/vault/VaultService`, `internal/keychain/KeychainService`
- **No Direct Storage Access**: TUI never touches `internal/storage`, `internal/crypto` directly

**Cryptography**:
- Zero changes to encryption/decryption logic
- TUI uses same AES-256-GCM encryption via VaultService
- Password handling uses existing `cmd/helpers.go` utilities

**Performance Requirements**:
- TUI startup: <100ms (REQ: <100ms for cached operations from tech.md)
- Memory overhead: <10MB additional (REQ: <50MB total from tech.md)
- Screen redraws: 60fps minimum for smooth user experience

### Project Structure (structure.md)

**Directory Organization**:
```
pass-cli/
├── main.go                     # Add TUI routing logic (minimal change)
├── cmd/
│   ├── root.go                 # Add TUI detection helper
│   ├── tui/                    # NEW: TUI package (self-contained)
│   │   ├── tui.go              # Entry point, program initialization
│   │   ├── model.go            # Application model/state
│   │   ├── update.go           # Message handling and state transitions
│   │   ├── view.go             # Main view assembly
│   │   ├── keys.go             # Keyboard bindings
│   │   ├── messages.go         # Custom message types
│   │   ├── views/              # Screen components
│   │   │   ├── list.go         # Credential list view
│   │   │   ├── detail.go       # Credential detail view
│   │   │   ├── form_add.go     # Add credential form
│   │   │   ├── form_edit.go    # Edit credential form
│   │   │   ├── confirm.go      # Confirmation dialogs
│   │   │   └── help.go         # Help overlay
│   │   ├── components/         # Reusable UI elements
│   │   │   ├── searchbar.go    # Search input component
│   │   │   ├── statusbar.go    # Status bar with indicators
│   │   │   └── notification.go # Notification bubbles
│   │   └── styles/             # Lipgloss styling
│   │       └── theme.go        # Color schemes, borders, styles
│   ├── [existing CLI commands]
│   └── helpers.go              # REUSED: password reading utilities
└── internal/
    └── vault/                  # UNCHANGED: shared by CLI and TUI
```

**Naming Conventions**:
- TUI files follow existing `snake_case.go` pattern
- TUI types use `PascalCase` (e.g., `ListView`, `DetailView`)
- TUI functions use `PascalCase` for public, `camelCase` for private
- TUI package imports follow: stdlib → external → internal order

**Code Organization Principles**:
- Each view file (list.go, detail.go, etc.) handles one screen's logic
- Components are reusable across views
- Styles are centralized in `styles/theme.go`
- No view directly accesses another view's state (message passing only)

## Code Reuse Analysis

### Existing Components to Leverage

**1. VaultService (internal/vault/vault.go)**:
```go
// TUI will use these existing methods:
- New(vaultPath string) (*VaultService, error)
- Initialize(masterPassword string, useKeychain bool) error
- Unlock(masterPassword string) error
- UnlockWithKeychain() error
- Lock()
- AddCredential(service, username, password, notes string) error
- GetCredential(service string) (*Credential, error)
- UpdateCredential(service, username, password, notes string) error
- DeleteCredential(service string) error
- ListCredentialsWithMetadata() ([]CredentialMetadata, error)
```

**2. KeychainService (internal/keychain/keychain.go)**:
```go
// TUI will use for status indicator:
- New() *KeychainService
- IsAvailable() bool
```

**3. Password Reading (cmd/helpers.go)**:
```go
// TUI will use for unlock prompt:
- readPassword() (string, error)
```

**4. Table Formatting (cmd/list.go)**:
```go
// TUI will adapt this logic for list view:
- formatRelativeTime(t time.Time) string
```

**5. Clipboard Integration (existing via atotto/clipboard)**:
```go
// TUI will use same clipboard library for password copy:
- clipboard.WriteAll(text string) error
```

### Integration Points

**main.go Entry Point**:
```go
// BEFORE (existing):
func main() {
    cmd.Execute()  // Always runs CLI
}

// AFTER (with TUI detection):
func main() {
    // If no args provided, launch TUI
    if len(os.Args) == 1 {
        tui.Run()
        return
    }

    // Otherwise run CLI as usual
    cmd.Execute()
}
```

**cmd/root.go Helper** (optional detection utility):
```go
// Add helper to detect if TUI should launch
func ShouldLaunchTUI() bool {
    return len(os.Args) == 1 && isatty.IsTerminal(os.Stdout.Fd())
}
```

**Vault Path Configuration**:
```go
// TUI will use same GetVaultPath() from cmd/root.go
vaultPath := cmd.GetVaultPath()
```

## Architecture

### Bubble Tea Framework (Elm Architecture)

Bubble Tea follows the **Model-Update-View** pattern:

```
┌─────────────────────────────────────────────────────┐
│                   Bubble Tea Loop                    │
│                                                      │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐  │
│  │  Model   │─Msg→│  Update  │─Model→│   View   │  │
│  │ (State)  │←────│ (Logic)  │←──────│ (Render) │  │
│  └──────────┘      └──────────┘      └──────────┘  │
│       │                  ↑                  │        │
│       │                  │                  │        │
│       └─── User Input ───┘                  │        │
│                                             │        │
│                                       Terminal       │
│                                         Output       │
└─────────────────────────────────────────────────────┘
```

**Model**: Holds application state (current view, credentials, selected index, etc.)
**Update**: Processes messages (key presses, API responses) and returns new state
**View**: Renders the current model state to a string for terminal display

### Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                              │
│  ┌────────────────┐                  ┌──────────────────┐   │
│  │ CLI Detection  │                  │ TUI Detection    │   │
│  │ (has args)     │                  │ (no args)        │   │
│  └────────┬───────┘                  └────────┬─────────┘   │
└───────────┼──────────────────────────────────┼──────────────┘
            │                                   │
            ▼                                   ▼
    ┌───────────────┐              ┌────────────────────────┐
    │  cmd/root.go  │              │  cmd/tui/tui.go        │
    │  Execute()    │              │  Run()                 │
    └───────┬───────┘              └────────┬───────────────┘
            │                                │
            │                                ▼
            │                    ┌───────────────────────────┐
            │                    │  cmd/tui/model.go         │
            │                    │  • Application State      │
            │                    │  • Current View           │
            │                    │  • Credentials Cache      │
            │                    └──────────┬────────────────┘
            │                               │
            │                               ▼
            │              ┌────────────────────────────────┐
            │              │  cmd/tui/update.go             │
            │              │  • Handle Key Presses          │
            │              │  • Process Messages            │
            │              │  • State Transitions           │
            │              └───────────┬────────────────────┘
            │                          │
            │                          ▼
            │              ┌───────────────────────────────┐
            │              │  cmd/tui/view.go              │
            │              │  • Assemble Screen            │
            │              │  • Delegate to Views          │
            │              └───────────┬───────────────────┘
            │                          │
            │         ┌────────────────┼────────────────┐
            │         ▼                ▼                ▼
            │   ┌──────────┐    ┌──────────┐    ┌──────────┐
            │   │ ListView │    │DetailView│    │ FormView │
            │   └──────────┘    └──────────┘    └──────────┘
            │
            └──────────────┬────────────────────────────┐
                           │  Shared Service Layer      │
                           ▼                            │
                  ┌──────────────────┐                 │
                  │ internal/vault/  │                 │
                  │  VaultService    │◄────────────────┘
                  └──────────────────┘
```

### State Management

**Application States** (AppState enum):
```go
type AppState int

const (
    StateUnlocking     AppState = iota  // Vault unlock screen
    StateList                            // Credential list view
    StateDetail                          // Credential detail view
    StateAdd                             // Add credential form
    StateEdit                            // Edit credential form
    StateConfirmDelete                   // Delete confirmation dialog
    StateHelp                            // Help overlay
    StateError                           // Error display
)
```

**State Transitions**:
```
┌──────────────┐
│ StateUnlocking│
└──────┬────────┘
       │ (vault unlocked)
       ▼
┌──────────────┐
│  StateList   │◄──────────────────┐
└──────┬────────┘                   │
       │                            │
       ├─ Enter ──► StateDetail ────┤
       ├─ 'a' ─────► StateAdd ───────┤
       ├─ 'e' ─────► StateEdit ──────┤
       ├─ 'd' ─────► StateConfirmDelete ─┤
       ├─ '?' ─────► StateHelp ──────────┤
       └─ 'q' ─────► Exit               │
                                         │
                (Escape / operation complete)
                                         │
                    ─────────────────────┘
```

## Components and Interfaces

### 1. Main Model (cmd/tui/model.go)

**Purpose**: Central application state container

```go
type Model struct {
    // Application state
    state          AppState
    prevState      AppState  // For returning from help/overlays

    // Vault integration
    vaultService   *vault.VaultService
    vaultPath      string
    keychainAvail  bool

    // Data cache
    credentials    []vault.CredentialMetadata
    selectedCred   *vault.Credential
    selectedIndex  int

    // UI state
    searchQuery    string
    notification   string
    notificationTime time.Time

    // View components (from bubbles library)
    listView       list.Model         // Credential list
    searchInput    textinput.Model    // Search bar
    formInputs     []textinput.Model  // Form fields
    viewport       viewport.Model     // Scrollable content

    // Window size
    width          int
    height         int

    // Error handling
    err            error
}
```

**Interfaces**:
```go
// Bubble Tea required methods
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

**Dependencies**:
- `internal/vault/VaultService`
- `internal/keychain/KeychainService`
- `github.com/charmbracelet/bubbles/list`
- `github.com/charmbracelet/bubbles/textinput`
- `github.com/charmbracelet/bubbles/viewport`

**Reuses**:
- `cmd.GetVaultPath()` for vault location
- `vault.New(path)` for service creation

### 2. Update Handler (cmd/tui/update.go)

**Purpose**: Process all messages and return new model state

```go
// Main update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        return m.handleResize(msg)
    case vaultUnlockedMsg:
        return m.handleVaultUnlocked()
    case credentialSavedMsg:
        return m.handleCredentialSaved(msg)
    case errorMsg:
        return m.handleError(msg)
    }
    return m, nil
}

// Key press routing
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Global keys
    switch msg.String() {
    case "ctrl+c", "q":
        if m.state == StateList {
            return m, tea.Quit
        }
    case "?", "f1":
        return m.enterHelp(), nil
    case "esc":
        return m.handleEscape(), nil
    }

    // State-specific key handling
    switch m.state {
    case StateList:
        return m.updateList(msg)
    case StateDetail:
        return m.updateDetail(msg)
    case StateAdd, StateEdit:
        return m.updateForm(msg)
    // ... other states
    }
    return m, nil
}
```

**Dependencies**:
- `internal/vault/VaultService` for CRUD operations
- Custom message types for async operations

### 3. View Renderer (cmd/tui/view.go)

**Purpose**: Assemble final screen output

```go
func (m Model) View() string {
    // Handle different states
    switch m.state {
    case StateUnlocking:
        return m.renderUnlocking()
    case StateList:
        return m.renderList()
    case StateDetail:
        return m.renderDetail()
    case StateAdd, StateEdit:
        return m.renderForm()
    case StateConfirmDelete:
        return m.renderConfirmation()
    case StateHelp:
        return m.renderHelp()
    case StateError:
        return m.renderError()
    }
    return ""
}

// Example list view assembly
func (m Model) renderList() string {
    var b strings.Builder

    // Header
    b.WriteString(m.renderHeader())
    b.WriteString("\n\n")

    // Search bar
    b.WriteString(m.searchInput.View())
    b.WriteString("\n\n")

    // Credential list
    b.WriteString(m.listView.View())
    b.WriteString("\n")

    // Status bar
    b.WriteString(m.renderStatusBar())

    // Notification overlay (if active)
    if m.hasNotification() {
        b.WriteString(m.renderNotification())
    }

    return b.String()
}
```

**Dependencies**:
- `cmd/tui/views/*` for screen-specific rendering
- `cmd/tui/components/*` for reusable UI elements
- `cmd/tui/styles/theme.go` for styling

### 4. List View (cmd/tui/views/list.go)

**Purpose**: Display credential list with search

```go
type ListView struct {
    list        list.Model
    searchBar   textinput.Model
    credentials []vault.CredentialMetadata
    filtered    []vault.CredentialMetadata
}

// Render list item
func (l *ListView) renderItem(meta vault.CredentialMetadata) string {
    // Format: "[Service]  Username  (Last used: X days ago)"
    service := meta.Service
    username := meta.Username
    lastUsed := formatRelativeTime(meta.LastAccessed)

    return fmt.Sprintf("%-30s %-25s %s", service, username, lastUsed)
}

// Filter credentials by search query
func (l *ListView) filter(query string) {
    query = strings.ToLower(query)
    l.filtered = make([]vault.CredentialMetadata, 0)

    for _, cred := range l.credentials {
        if strings.Contains(strings.ToLower(cred.Service), query) ||
           strings.Contains(strings.ToLower(cred.Username), query) {
            l.filtered = append(l.filtered, cred)
        }
    }
}
```

**Dependencies**:
- `github.com/charmbracelet/bubbles/list` for list component
- `github.com/charmbracelet/bubbles/textinput` for search bar

**Reuses**:
- `formatRelativeTime()` from `cmd/list.go`

### 5. Detail View (cmd/tui/views/detail.go)

**Purpose**: Display full credential details

```go
type DetailView struct {
    credential      *vault.Credential
    passwordVisible bool
    viewport        viewport.Model
}

func (d *DetailView) Render() string {
    var b strings.Builder

    // Service header
    b.WriteString(styles.TitleStyle.Render(d.credential.Service))
    b.WriteString("\n\n")

    // Fields
    b.WriteString(d.renderField("Username", d.credential.Username))
    b.WriteString(d.renderField("Password", d.renderPassword()))
    if d.credential.Notes != "" {
        b.WriteString(d.renderField("Notes", d.credential.Notes))
    }
    b.WriteString(d.renderField("Created", d.credential.CreatedAt.Format("2006-01-02 15:04")))
    b.WriteString(d.renderField("Updated", d.credential.UpdatedAt.Format("2006-01-02 15:04")))

    // Usage records
    if len(d.credential.UsageRecord) > 0 {
        b.WriteString("\n")
        b.WriteString(d.renderUsageRecords())
    }

    return b.String()
}

func (d *DetailView) renderPassword() string {
    if d.passwordVisible {
        return d.credential.Password
    }
    return strings.Repeat("*", len(d.credential.Password))
}
```

**Dependencies**:
- `github.com/charmbracelet/bubbles/viewport` for scrollable content

### 6. Add/Edit Form Views (cmd/tui/views/form_add.go, form_edit.go)

**Purpose**: Interactive forms for credential input

```go
type FormView struct {
    mode        FormMode  // Add or Edit
    service     textinput.Model
    username    textinput.Model
    password    textinput.Model
    notes       textarea.Model
    focusedField int
    originalCred *vault.Credential  // For edit mode
}

// Form navigation
func (f *FormView) NextField() {
    f.focusedField = (f.focusedField + 1) % 4
    f.updateFocus()
}

func (f *FormView) PrevField() {
    f.focusedField = (f.focusedField - 1 + 4) % 4
    f.updateFocus()
}

// Validation
func (f *FormView) Validate() error {
    if strings.TrimSpace(f.service.Value()) == "" {
        return errors.New("service name is required")
    }
    return nil
}

// Get credential from form
func (f *FormView) GetCredential() vault.Credential {
    return vault.Credential{
        Service:  strings.TrimSpace(f.service.Value()),
        Username: strings.TrimSpace(f.username.Value()),
        Password: f.password.Value(),
        Notes:    f.notes.Value(),
    }
}
```

**Dependencies**:
- `github.com/charmbracelet/bubbles/textinput`
- `github.com/charmbracelet/bubbles/textarea`

**Reuses**:
- Password generation logic from existing `cmd/generate.go`

### 7. Status Bar Component (cmd/tui/components/statusbar.go)

**Purpose**: Display system status and shortcuts

```go
type StatusBar struct {
    keychainAvailable bool
    credentialCount   int
    currentView       string
    width             int
}

func (s *StatusBar) Render() string {
    // Left side: keychain status and credential count
    var left strings.Builder
    if s.keychainAvailable {
        left.WriteString(styles.SuccessStyle.Render("🔓 Keychain"))
    } else {
        left.WriteString(styles.WarningStyle.Render("🔒 Password"))
    }
    left.WriteString(fmt.Sprintf("  %d credentials", s.credentialCount))

    // Right side: current view and shortcuts
    right := fmt.Sprintf("%s  |  ?: Help  q: Quit", s.currentView)

    // Center spacing
    leftStr := left.String()
    spacing := strings.Repeat(" ", s.width-len(leftStr)-len(right))

    return styles.StatusBarStyle.Render(leftStr + spacing + right)
}
```

**Dependencies**:
- `cmd/tui/styles/theme.go`

**Reuses**:
- Keychain availability check from `internal/keychain`

### 8. Help Overlay (cmd/tui/views/help.go)

**Purpose**: Display keyboard shortcuts

```go
type HelpView struct {
    viewport viewport.Model
}

func (h *HelpView) Render(state AppState) string {
    var b strings.Builder

    // Title
    b.WriteString(styles.TitleStyle.Render("Keyboard Shortcuts"))
    b.WriteString("\n\n")

    // Global shortcuts
    b.WriteString(renderShortcutSection("Global", globalShortcuts))

    // State-specific shortcuts
    switch state {
    case StateList:
        b.WriteString(renderShortcutSection("List View", listShortcuts))
    case StateDetail:
        b.WriteString(renderShortcutSection("Detail View", detailShortcuts))
    // ... other states
    }

    b.WriteString("\n")
    b.WriteString(styles.SubtleStyle.Render("Press any key to close"))

    return h.viewport.View()
}

var globalShortcuts = []Shortcut{
    {"q", "Quit"},
    {"?/F1", "Help"},
    {"Esc", "Back/Cancel"},
}

var listShortcuts = []Shortcut{
    {"↑/↓ or j/k", "Navigate"},
    {"Enter", "View details"},
    {"/", "Search"},
    {"a", "Add credential"},
    {"e", "Edit credential"},
    {"d", "Delete credential"},
}
```

**Dependencies**:
- `github.com/charmbracelet/bubbles/viewport`

### 9. Theme and Styles (cmd/tui/styles/theme.go)

**Purpose**: Centralized styling and colors

```go
package styles

import "github.com/charmbracelet/lipgloss"

var (
    // Colors
    PrimaryColor   = lipgloss.Color("#00ADD8")  // Go cyan
    SecondaryColor = lipgloss.Color("#007D9C")
    SuccessColor   = lipgloss.Color("#00C853")  // Green
    WarningColor   = lipgloss.Color("#FF6F00")  // Orange
    ErrorColor     = lipgloss.Color("#DD2C00")  // Red
    SubtleColor    = lipgloss.Color("#666666")  // Gray

    // Text styles
    TitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(PrimaryColor).
        MarginBottom(1)

    SelectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Background(PrimaryColor).
        Bold(true)

    SuccessStyle = lipgloss.NewStyle().
        Foreground(SuccessColor).
        Bold(true)

    ErrorStyle = lipgloss.NewStyle().
        Foreground(ErrorColor).
        Bold(true)

    SubtleStyle = lipgloss.NewStyle().
        Foreground(SubtleColor)

    // Borders
    BorderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(SubtleColor).
        Padding(1, 2)

    // Status bar
    StatusBarStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("#1A1A1A")).
        Foreground(lipgloss.Color("#FFFFFF")).
        Padding(0, 1)
)
```

**Dependencies**:
- `github.com/charmbracelet/lipgloss`

## Data Models

### 1. Application Model (already defined above in Components)

### 2. Custom Message Types (cmd/tui/messages.go)

```go
// Async operation messages
type vaultUnlockedMsg struct {
    credentials []vault.CredentialMetadata
}

type credentialSavedMsg struct {
    service string
}

type credentialDeletedMsg struct {
    service string
}

type credentialCopiedMsg struct{}

type errorMsg struct {
    err error
}

// Notification message
type notificationMsg struct {
    message string
    level   NotificationLevel
}

type NotificationLevel int
const (
    NotificationInfo NotificationLevel = iota
    NotificationSuccess
    NotificationWarning
    NotificationError
)
```

### 3. Reused Models from internal/vault

```go
// These are REUSED, not redefined:
- vault.Credential
- vault.CredentialMetadata
- vault.VaultData
```

## Error Handling

### Error Scenarios

1. **Vault Not Found**:
   - **Handling**: Display friendly error message with instructions
   - **User Impact**: "Vault not found. Run 'pass-cli init' to create one."
   - **Recovery**: Exit TUI, user runs init command

2. **Unlock Failure (Wrong Password)**:
   - **Handling**: Display error, allow retry (max 3 attempts)
   - **User Impact**: "Incorrect password. Try again (2 attempts remaining)"
   - **Recovery**: User re-enters password or exits with Ctrl+C

3. **Keychain Unavailable**:
   - **Handling**: Gracefully fall back to password prompt
   - **User Impact**: Status bar shows "🔒 Password" instead of "🔓 Keychain"
   - **Recovery**: User enters password manually

4. **Duplicate Credential (Add)**:
   - **Handling**: Show validation error in form
   - **User Impact**: "Credential 'github' already exists"
   - **Recovery**: User changes service name or cancels

5. **Credential Not Found (Get/Update/Delete)**:
   - **Handling**: Display error notification
   - **User Impact**: "Credential not found"
   - **Recovery**: Return to list view, refresh credentials

6. **Clipboard Failure**:
   - **Handling**: Show warning, continue operation
   - **User Impact**: "Clipboard unavailable. Password not copied."
   - **Recovery**: User can view password with 'm' toggle instead

7. **Terminal Too Small**:
   - **Handling**: Display minimum size message
   - **User Impact**: "Terminal too small (min 80x24). Please resize."
   - **Recovery**: User resizes terminal, TUI auto-adjusts

8. **Save/Update Failure**:
   - **Handling**: Show error, preserve form data
   - **User Impact**: "Failed to save credential: [error details]"
   - **Recovery**: User can retry or cancel, form data not lost

## Testing Strategy

### Unit Testing

**Components to Test**:
1. **Model State Transitions** (`model_test.go`, `update_test.go`):
   - Test state changes (List → Detail → Edit → List)
   - Test key press handling for each state
   - Test search filtering logic
   - Test form validation

2. **View Rendering** (`view_test.go`, `views/*_test.go`):
   - Test view output contains expected elements
   - Test password masking/unmasking
   - Test status bar updates
   - Test notification display

3. **Components** (`components/*_test.go`):
   - Test search bar filtering
   - Test status bar formatting
   - Test notification timeout logic

**Testing Approach**:
```go
func TestModelStateTransition(t *testing.T) {
    m := Model{state: StateList}

    // Simulate pressing Enter (view details)
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    if m.state != StateDetail {
        t.Errorf("Expected StateDetail, got %v", m.state)
    }
}

func TestSearchFilter(t *testing.T) {
    lv := ListView{
        credentials: []vault.CredentialMetadata{
            {Service: "github", Username: "user1"},
            {Service: "gitlab", Username: "user2"},
        },
    }

    lv.filter("git")

    if len(lv.filtered) != 2 {
        t.Errorf("Expected 2 filtered, got %d", len(lv.filtered))
    }

    lv.filter("hub")

    if len(lv.filtered) != 1 || lv.filtered[0].Service != "github" {
        t.Error("Filter did not match 'github'")
    }
}
```

**Mock Strategy**:
- Mock VaultService for all tests (no real file I/O)
- Use test credentials with known data
- Test components in isolation

### Integration Testing

**Test File**: `test/tui_integration_test.go` (with `//go:build integration` tag)

**Tests to Implement**:

1. **TUI Launch and Exit**:
```go
func TestIntegration_TUILaunchExit(t *testing.T) {
    // Start TUI in background
    // Send 'q' key
    // Verify clean exit
}
```

2. **Full CRUD Workflow**:
```go
func TestIntegration_TUICRUDWorkflow(t *testing.T) {
    // Init vault with keychain
    // Launch TUI
    // Navigate to add form
    // Add credential
    // View detail
    // Edit credential
    // Delete credential
    // Verify all operations
}
```

3. **Keychain Integration**:
```go
func TestIntegration_TUIKeychainUnlock(t *testing.T) {
    // Init vault with keychain enabled
    // Launch TUI
    // Verify no password prompt (auto-unlock)
    // Verify keychain indicator in status bar
}
```

4. **Search Functionality**:
```go
func TestIntegration_TUISearch(t *testing.T) {
    // Create vault with multiple credentials
    // Launch TUI
    // Type search query
    // Verify filtering works
}
```

### Manual Testing

**Cross-Platform Terminal Testing**:
- **Windows**: Windows Terminal, PowerShell, CMD
- **macOS**: Terminal.app, iTerm2
- **Linux**: GNOME Terminal, Konsole, Alacritty

**Test Scenarios**:
1. Color rendering (256-color vs basic)
2. Unicode character support (borders, icons)
3. Resize handling
4. Keyboard input (especially special keys)
5. Clipboard operations

## Dependencies

### New Dependencies to Add

```go
// go.mod additions:
require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.9.1
)
```

**Dependency Justification**:
- **bubbletea**: Industry-standard TUI framework, used by GitHub CLI and many popular tools
- **bubbles**: Pre-built components (saves development time, battle-tested)
- **lipgloss**: Styling library (consistent with Bubble Tea ecosystem)

**Binary Size Impact**: +1-2MB (acceptable per tech.md requirements)

### Existing Dependencies (Reused)

```go
// Already in go.mod:
- github.com/atotto/clipboard v0.1.4
- github.com/spf13/cobra v1.10.1
- golang.org/x/term v0.35.0
```

## Implementation Phases

### Phase 1: Foundation (Tasks 1-3)
- Add Bubble Tea dependencies
- Create TUI entry point detection in main.go
- Implement basic Model, Update, View structure
- Get "Hello TUI" working

### Phase 2: Core Views (Tasks 4-5)
- Implement list view with search
- Implement detail view
- Test navigation between views

### Phase 3: Forms (Tasks 6-7)
- Implement add credential form
- Implement edit credential form
- Test CRUD operations

### Phase 4: Polish (Tasks 8-11)
- Add delete confirmation dialog
- Create status bar with keychain indicator
- Add help overlay
- Apply theme and styling

### Phase 5: Testing (Tasks 12-13)
- Write unit tests for all components
- Write integration tests
- Cross-platform manual testing

## Performance Considerations

**Optimization Strategies**:
1. **Lazy Loading**: Only render visible list items
2. **Credential Caching**: Load credentials once on unlock, cache in model
3. **Efficient Redraws**: Only redraw changed components
4. **Debounce Search**: Delay filtering until user stops typing

**Performance Targets** (from tech.md):
- TUI startup: <100ms ✓
- Search filter: <50ms ✓
- Screen redraw: 60fps (16ms frame time) ✓
- Memory overhead: <10MB ✓

## Security Considerations

**Security Properties** (unchanged from CLI):
1. **Password Masking**: Passwords masked by default in all views
2. **Clipboard Timeout**: Same behavior as CLI (existing clipboard utility)
3. **No Logging**: Sensitive data never logged (enforced in TUI code reviews)
4. **Memory Clearing**: Password strings cleared when model state changes
5. **Keychain Integration**: Uses existing KeychainService (no new security surface)

**TUI-Specific Security**:
- Screen content can be viewed by others looking at terminal (same risk as CLI)
- No additional attack surface introduced by TUI
- TUI does not weaken existing security properties

## Accessibility Considerations

**Keyboard-Only Operation**:
- 100% functionality via keyboard (no mouse required)
- Consistent key bindings across views
- Help overlay accessible from any screen

**Visual Accessibility**:
- High contrast colors (4.5:1 ratio minimum)
- Graceful degradation to basic colors if terminal doesn't support 256
- Clear focus indicators for selected items

**Screen Reader Compatibility**:
- TUI renders plain text (screen reader compatible)
- Descriptive labels for all UI elements
- Logical reading order

## Future Enhancements (Out of Scope for v1)

1. **Mouse Support**: Enable mouse clicks for selection
2. **Custom Themes**: User-configurable color schemes
3. **Animations**: Smooth transitions between views
4. **Multi-select**: Bulk operations on credentials
5. **Export View**: Visual export interface (currently CLI only)
6. **Import View**: Visual import interface (currently CLI only)
7. **Password Strength Indicator**: Visual bar showing password strength
8. **Recent Credentials**: Quick access to recently used credentials
9. **Favorites**: Pin frequently used credentials to top

## Migration and Rollback

**Migration**: None required - TUI is additive feature

**Rollback Strategy**:
- Remove `cmd/tui/` directory
- Revert changes to `main.go` (TUI detection logic)
- Remove Bubble Tea dependencies from `go.mod`
- CLI functionality completely unaffected

**Backward Compatibility**: 100% - all CLI commands work identically
