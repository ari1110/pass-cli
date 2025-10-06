# components/table.go

## Purpose
Credential list table using tview.Table. Displays filtered credentials based on selected category with selection support.

## Responsibilities

1. **Credential display**: Show list of credentials in table format
2. **Column layout**: Service, Username, Last Used columns
3. **Filtering**: Display only credentials matching selected category
4. **Selection handling**: Track selected credential and enable navigation
5. **Sorting**: Sort by service name (or other columns)
6. **Refresh on changes**: Update table when credentials or selection changes

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For AppState
- `pass-cli/internal/vault` - For Credential type

### External Dependencies
- `github.com/rivo/tview` - Table primitive
- `github.com/gdamore/tcell/v2` - For colors and alignment

## Key Types

### `CredentialTable`
**Purpose**: Wrapper around tview.Table with filtering and refresh logic

**Fields**:
```go
type CredentialTable struct {
    *tview.Table  // Embedded Table

    appState         *models.AppState
    filteredCreds    []vault.Credential  // Currently displayed credentials
    selectedIndex    int                 // Selected row index
}
```

## Key Functions

### Constructor

#### `NewCredentialTable(appState *models.AppState) *CredentialTable`
**Purpose**: Create and configure table

**Steps**:
1. Create tview.Table
2. Configure fixed header row
3. Setup column layout
4. Configure styles (borders, colors)
5. Setup selection handler
6. Enable selection and highlighting
7. Return wrapped table

### Table Building

#### `Refresh()`
**Purpose**: Rebuild table from filtered credentials

**Steps**:
1. Get all credentials from appState
2. Get selected category from appState
3. Filter credentials by category
4. Clear table (keep header)
5. Populate rows with filtered credentials
6. Restore selection if possible

**Called when**: Credentials change OR category selection changes

#### `buildHeader()`
**Purpose**: Create fixed header row

**Columns**:
| Service | Username | Last Used |

**Styling**:
- Header in accent color (bold)
- Fixed row (not selectable)
- Column widths: flexible based on content

#### `populateRows(credentials []vault.Credential)`
**Purpose**: Add credential rows to table

**For each credential**:
- Row with service, username, last used timestamp
- Store credential reference in cell metadata
- Apply row styling (alternating colors optional)

### Filtering

#### `filterByCategory(credentials []vault.Credential, category string) []vault.Credential`
**Purpose**: Filter credentials by selected category

**Logic**:
```go
if category == "" {
    return credentials  // Show all
}

filtered := []vault.Credential{}
for _, cred := range credentials {
    if cred.Category == category {
        filtered = append(filtered, cred)
    }
}
return filtered
```

### Selection Handling

#### `OnEnterPressed(handler func(credential *vault.Credential))`
**Purpose**: Register callback for Enter key on selected row

**Behavior**:
```go
table.SetSelectedFunc(func(row, column int) {
    if row == 0 {
        return  // Header row, ignore
    }

    cred := getCredentialFromRow(row)
    appState.SetSelectedCredential(cred)
    // Trigger detail view
})
```

#### `getCredentialFromRow(row int) *vault.Credential`
**Purpose**: Retrieve credential from table cell metadata

**Implementation**: Store credential in cell reference during population

### Styling

#### `applyStyles()`
**Purpose**: Apply theme colors and borders

**Configuration**:
- Rounded borders
- Border color from theme
- Background color from theme
- Selected row highlight color
- Fixed column widths or auto-sizing

**Example**:
```go
table.SetBorder(true).
    SetTitle(" Credentials ").
    SetTitleAlign(tview.AlignLeft).
    SetBorderColor(styles.BorderColor).
    SetBackgroundColor(styles.BackgroundColor)

table.SetSelectable(true, false)  // Select rows, not columns
table.SetFixed(1, 0)  // Fix header row
```

## Example Structure

```go
package components

import (
    "fmt"
    "time"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/models"
    "pass-cli/cmd/tui-tview/styles"
    "pass-cli/internal/vault"
)

type CredentialTable struct {
    *tview.Table
    appState      *models.AppState
    filteredCreds []vault.Credential
}

func NewCredentialTable(appState *models.AppState) *CredentialTable {
    table := tview.NewTable()

    ct := &CredentialTable{
        Table:    table,
        appState: appState,
    }

    // Configure table
    ct.SetSelectable(true, false)
    ct.SetFixed(1, 0)  // Fix header row
    ct.applyStyles()

    // Build header
    ct.buildHeader()

    // Setup selection handler
    ct.SetSelectedFunc(ct.onSelect)

    // Initial population
    ct.Refresh()

    return ct
}

func (ct *CredentialTable) buildHeader() {
    headers := []string{"Service", "Username", "Last Used"}
    for col, header := range headers {
        cell := tview.NewTableCell(header).
            SetTextColor(styles.AccentColor).
            SetAlign(tview.AlignLeft).
            SetSelectable(false).
            SetExpansion(1)
        ct.SetCell(0, col, cell)
    }
}

func (ct *CredentialTable) Refresh() {
    // Get credentials and filter
    allCreds := ct.appState.GetCredentials()
    category := ct.appState.GetSelectedCategory()
    ct.filteredCreds = ct.filterByCategory(allCreds, category)

    // Clear table (keep header)
    for row := ct.GetRowCount() - 1; row > 0; row-- {
        ct.RemoveRow(row)
    }

    // Populate rows
    for i, cred := range ct.filteredCreds {
        ct.addCredentialRow(i+1, cred)
    }

    // Update title with count
    ct.SetTitle(fmt.Sprintf(" Credentials (%d) ", len(ct.filteredCreds)))
}

func (ct *CredentialTable) addCredentialRow(row int, cred vault.Credential) {
    // Service column
    serviceCell := tview.NewTableCell(cred.Service).
        SetTextColor(tcell.ColorWhite).
        SetAlign(tview.AlignLeft).
        SetReference(cred)  // Store credential in cell

    // Username column
    usernameCell := tview.NewTableCell(cred.Username).
        SetTextColor(tcell.ColorGray).
        SetAlign(tview.AlignLeft)

    // Last used column
    lastUsed := "Never"
    if !cred.LastUsed.IsZero() {
        lastUsed = formatRelativeTime(cred.LastUsed)
    }
    lastUsedCell := tview.NewTableCell(lastUsed).
        SetTextColor(tcell.ColorGray).
        SetAlign(tview.AlignLeft)

    ct.SetCell(row, 0, serviceCell)
    ct.SetCell(row, 1, usernameCell)
    ct.SetCell(row, 2, lastUsedCell)
}

func (ct *CredentialTable) onSelect(row, column int) {
    if row == 0 {
        return  // Header
    }

    // Get credential from cell reference
    cell := ct.GetCell(row, 0)
    if cell != nil {
        if cred, ok := cell.GetReference().(vault.Credential); ok {
            ct.appState.SetSelectedCredential(&cred)
        }
    }
}

func (ct *CredentialTable) filterByCategory(creds []vault.Credential, category string) []vault.Credential {
    if category == "" {
        return creds
    }

    filtered := []vault.Credential{}
    for _, cred := range creds {
        if cred.Category == category {
            filtered = append(filtered, cred)
        }
    }
    return filtered
}

func (ct *CredentialTable) applyStyles() {
    ct.SetBorder(true).
        SetTitle(" Credentials ").
        SetTitleAlign(tview.AlignLeft).
        SetBorderColor(styles.BorderColor).
        SetBackgroundColor(styles.BackgroundColor)
}

func formatRelativeTime(t time.Time) string {
    duration := time.Since(t)
    switch {
    case duration < time.Hour:
        return fmt.Sprintf("%dm ago", int(duration.Minutes()))
    case duration < 24*time.Hour:
        return fmt.Sprintf("%dh ago", int(duration.Hours()))
    default:
        return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
    }
}
```

## Interaction Flow

```
Category "AWS" selected in sidebar
    ↓
appState.SetSelectedCategory("AWS")
    ↓
AppState notifies selection changed
    ↓
Table.Refresh() called
    ↓
Filter credentials by category
    ↓
Rebuild table rows
    ↓
Display AWS credentials only
    ↓
User presses Enter on row
    ↓
onSelect() fires
    ↓
appState.SetSelectedCredential(cred)
    ↓
Detail view shows credential details
```

## Visual Design

### Column Layout
- **Service**: 30% width or auto-expand
- **Username**: 30% width or auto-expand
- **Last Used**: 20% width or fixed

### Row Styling
- **Header**: Accent color, not selectable
- **Selected row**: Highlighted (tview handles this)
- **Normal rows**: White text
- **Empty table**: Show "No credentials" message

### Borders
- Rounded corners
- Title shows count: " Credentials (5) "
- Border color from theme

## Keyboard Shortcuts

Handled by tview.Table automatically:
- **Arrow Up/Down**: Navigate rows
- **Home/End**: Jump to first/last
- **Page Up/Down**: Scroll by page
- **Enter**: Select credential

Custom shortcuts (handled in event handler):
- **Tab**: Move focus to detail view
- **n**: New credential
- **e**: Edit selected credential
- **d**: Delete selected credential

## State Integration

### Reads from AppState
- `GetCredentials()` - All credentials
- `GetSelectedCategory()` - For filtering

### Writes to AppState
- `SetSelectedCredential(cred)` - When row selected

### Responds to AppState callbacks
- `onCredentialsChanged` - Triggers Refresh()
- `onSelectionChanged` - Triggers Refresh() if category changed

## Edge Cases

1. **No credentials**: Show "No credentials" message row
2. **Empty category**: Show "No credentials in this category"
3. **Long service names**: Truncate with ellipsis
4. **Selection after delete**: Select next row, or previous if last
5. **Sorting stability**: Maintain order when refreshing

## Performance Considerations

- **Large credential lists**: tview.Table handles thousands of rows efficiently
- **Filter caching**: Cache filtered results if category unchanged
- **Incremental updates**: For now, full refresh is fine (fast enough)

## Testing Considerations

- **Test filtering**: Verify correct credentials displayed per category
- **Test selection**: Verify state updates on Enter
- **Test refresh**: Verify table rebuilds correctly
- **Test empty states**: Verify empty messages display
- **Mock AppState**: Use mock for isolated testing

## Future Enhancements

- **Column sorting**: Click column header to sort
- **Search/filter**: Type to filter by service or username
- **Multi-select**: Bulk operations on multiple credentials
- **Copy to clipboard**: Quick copy password without opening detail
- **Color coding**: Color rows by strength, age, or usage
- **Custom columns**: User-configurable column visibility
