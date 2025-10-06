# layout/manager.go

## Purpose
Responsive layout manager using tview.Flex. Handles terminal size detection and adaptive panel visibility based on breakpoints.

## Responsibilities

1. **Responsive layouts**: Adjust panel visibility based on terminal width
2. **Breakpoint management**: Three layouts (small, medium, large)
3. **Component composition**: Combine sidebar, table, detail, statusbar
4. **Resize handling**: Detect terminal size changes and reflow layout
5. **Dimension calculation**: Calculate optimal component sizes

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For AppState
- `pass-cli/cmd/tui-tview/components` - For UI components

### External Dependencies
- `github.com/rivo/tview` - Flex layout
- `github.com/gdamore/tcell/v2` - For screen size detection

## Key Types

### `LayoutManager`
**Purpose**: Manages responsive layout composition

**Fields**:
```go
type LayoutManager struct {
    app      *tview.Application
    appState *models.AppState

    // Current dimensions
    width  int
    height int

    // Layout primitives
    mainLayout *tview.Flex  // Root layout
    contentRow *tview.Flex  // Main content area (sidebar + main + metadata)
    mainColumn *tview.Flex  // Main column (table or table+detail)

    // Component references
    sidebar    *tview.TreeView
    table      *tview.Table
    detailView *tview.TextView
    statusBar  *tview.TextView

    // Breakpoints
    mediumBreakpoint int  // 80 columns
    largeBreakpoint  int  // 120 columns
}
```

### `LayoutMode`
**Purpose**: Enum for layout configurations

```go
type LayoutMode int

const (
    LayoutSmall  LayoutMode = iota  // < 80 cols: Table only (no sidebar)
    LayoutMedium                     // 80-120 cols: Sidebar + Table
    LayoutLarge                      // > 120 cols: Sidebar + Table + Detail
)
```

## Key Functions

### Constructor

#### `NewLayoutManager(app *tview.Application, appState *models.AppState) *LayoutManager`
**Purpose**: Create layout manager with default breakpoints

**Initial state**: Create layout structure but don't populate yet

### Layout Building

#### `CreateMainLayout() *tview.Flex`
**Purpose**: Build complete layout structure

**Structure**:
```
┌───────────────────────────────┐
│       mainLayout (FlexRow)    │
├───────────────────────────────┤
│   contentRow (FlexColumn)     │  ← Main content area
│   ┌──────┬───────┬──────────┐ │
│   │      │       │          │ │
│   │ Side │ Table │  Detail  │ │
│   │ bar  │       │          │ │
│   │      │       │          │ │
│   └──────┴───────┴──────────┘ │
├───────────────────────────────┤
│       statusBar (1 row)       │  ← Fixed height
└───────────────────────────────┘
```

**Steps**:
1. Get components from appState
2. Create content row (Flex horizontal)
3. Create main layout (Flex vertical)
4. Add content row and status bar
5. Setup resize detection
6. Return main layout

#### `rebuildLayout()`
**Purpose**: Reconstruct layout based on current dimensions

**Steps**:
1. Determine layout mode from width
2. Clear contentRow
3. Add components based on mode:
   - Small: Table only (full width)
   - Medium: Sidebar (20) + Table (flex)
   - Large: Sidebar (20) + Table (flex) + Detail (40)
4. Update component sizes

### Resize Handling

#### `HandleResize(width, height int)`
**Purpose**: Respond to terminal size changes

**Steps**:
1. Store new dimensions
2. Determine if layout mode changed
3. If mode changed, rebuild layout
4. Otherwise, just update sizes
5. Force redraw

**Called when**: Terminal resize detected

#### `detectTerminalSize() (int, int)`
**Purpose**: Get current terminal dimensions

**Implementation**: Use tcell screen.Size()

**Setup**:
```go
// In CreateMainLayout, use SetDrawFunc for initial size detection
mainLayout.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
    if lm.width == 0 && lm.height == 0 {
        termWidth, termHeight := screen.Size()
        lm.HandleResize(termWidth, termHeight)
    }
    return x, y, width, height
})
```

### Layout Mode Logic

#### `determineLayoutMode(width int) LayoutMode`
**Purpose**: Calculate layout mode from terminal width

**Logic**:
```go
func (lm *LayoutManager) determineLayoutMode(width int) LayoutMode {
    if width < lm.mediumBreakpoint {
        return LayoutSmall
    }
    if width < lm.largeBreakpoint {
        return LayoutMedium
    }
    return LayoutLarge
}
```

### Component Sizing

#### `calculateSidebarWidth() int`
**Purpose**: Get sidebar width for current layout

**Returns**:
- Small mode: 0 (hidden)
- Medium/Large: 20 columns (fixed)

#### `calculateDetailWidth() int`
**Purpose**: Get detail view width for current layout

**Returns**:
- Small/Medium: 0 (hidden)
- Large: 40 columns (fixed or flex)

#### `calculateTableWidth() int`
**Purpose**: Get table width for current layout

**Returns**: Flex (fills remaining space)

## Example Structure

```go
package layout

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/models"
)

type LayoutMode int

const (
    LayoutSmall LayoutMode = iota
    LayoutMedium
    LayoutLarge
)

type LayoutManager struct {
    app      *tview.Application
    appState *models.AppState

    width  int
    height int

    mainLayout *tview.Flex
    contentRow *tview.Flex

    sidebar    *tview.TreeView
    table      *tview.Table
    detailView *tview.TextView
    statusBar  *tview.TextView

    mediumBreakpoint int
    largeBreakpoint  int

    currentMode LayoutMode
}

func NewLayoutManager(app *tview.Application, appState *models.AppState) *LayoutManager {
    return &LayoutManager{
        app:              app,
        appState:         appState,
        mediumBreakpoint: 80,
        largeBreakpoint:  120,
        currentMode:      LayoutSmall,
    }
}

func (lm *LayoutManager) CreateMainLayout() *tview.Flex {
    // Get component references from appState
    lm.sidebar = lm.appState.GetSidebar()
    lm.table = lm.appState.GetTable()
    lm.detailView = lm.appState.GetDetailView()
    lm.statusBar = lm.appState.GetStatusBar()

    // Create content row (horizontal layout)
    lm.contentRow = tview.NewFlex().SetDirection(tview.FlexColumn)

    // Create main layout (vertical)
    lm.mainLayout = tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(lm.contentRow, 0, 1, true).   // Content area (flex)
        AddItem(lm.statusBar, 1, 0, false)    // Status bar (fixed 1 row)

    // Setup resize detection
    lm.mainLayout.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
        if lm.width == 0 && lm.height == 0 {
            termWidth, termHeight := screen.Size()
            lm.HandleResize(termWidth, termHeight)
        }
        return x, y, width, height
    })

    return lm.mainLayout
}

func (lm *LayoutManager) HandleResize(width, height int) {
    lm.width = width
    lm.height = height

    newMode := lm.determineLayoutMode(width)
    if newMode != lm.currentMode {
        lm.currentMode = newMode
        lm.rebuildLayout()
    }
}

func (lm *LayoutManager) rebuildLayout() {
    lm.contentRow.Clear()

    switch lm.currentMode {
    case LayoutSmall:
        // Table only (full width)
        lm.contentRow.AddItem(lm.table, 0, 1, true)

    case LayoutMedium:
        // Sidebar + Table
        lm.contentRow.
            AddItem(lm.sidebar, 20, 0, false).
            AddItem(lm.table, 0, 1, true)

    case LayoutLarge:
        // Sidebar + Table + Detail
        lm.contentRow.
            AddItem(lm.sidebar, 20, 0, false).
            AddItem(lm.table, 0, 1, true).
            AddItem(lm.detailView, 40, 0, false)
    }
}

func (lm *LayoutManager) determineLayoutMode(width int) LayoutMode {
    if width < lm.mediumBreakpoint {
        return LayoutSmall
    }
    if width < lm.largeBreakpoint {
        return LayoutMedium
    }
    return LayoutLarge
}

func (lm *LayoutManager) GetCurrentMode() LayoutMode {
    return lm.currentMode
}

func (lm *LayoutManager) SetBreakpoints(medium, large int) {
    lm.mediumBreakpoint = medium
    lm.largeBreakpoint = large
}
```

## Interaction Flow

```
Terminal resized to 100 columns
    ↓
SetDrawFunc detects size change
    ↓
HandleResize(100, 40) called
    ↓
determineLayoutMode(100) → LayoutMedium
    ↓
Mode changed from Small → Medium
    ↓
rebuildLayout()
    ↓
contentRow cleared
    ↓
Add sidebar (20 cols) + table (flex)
    ↓
UI redraws with new layout
```

## Breakpoint Strategy

### Small (< 80 columns)
**Why**: Not enough space for sidebar
**Layout**: Table only
**Use case**: Narrow terminals, vertical split screens

### Medium (80-120 columns)
**Why**: Enough for sidebar + table
**Layout**: Sidebar + Table
**Use case**: Standard terminal windows

### Large (> 120 columns)
**Why**: Room for all panels
**Layout**: Sidebar + Table + Detail
**Use case**: Wide terminals, full screen

## Component Widths

| Component | Small | Medium | Large |
|-----------|-------|--------|-------|
| Sidebar   | 0     | 20     | 20    |
| Table     | Flex  | Flex   | Flex  |
| Detail    | 0     | 0      | 40    |
| Status    | Full  | Full   | Full  |

**Flex**: Component takes remaining available space

## Edge Cases

1. **Very narrow terminal** (< 40 cols): Show warning, minimal layout
2. **Rapid resizing**: Debounce or just rebuild (tview is fast enough)
3. **Detail hidden in medium**: Show detail as modal instead
4. **Component sizes**: Ensure minimum widths are respected

## Performance Considerations

- **Rebuild is fast**: tview handles layout efficiently
- **No debouncing needed**: Immediate response feels better
- **Component reuse**: Don't recreate, just rearrange

## Integration with Navigation

When layout mode changes, update navigation:
```go
// In HandleResize:
if newMode != lm.currentMode {
    if newMode == LayoutSmall {
        // Sidebar hidden - move focus to table
        nav.SetFocus(models.FocusTable)
    }
}
```

## Testing Considerations

- **Test each breakpoint**: Verify layout at 79, 80, 119, 120 columns
- **Test resize transitions**: Verify smooth mode changes
- **Test component visibility**: Ensure hidden components don't receive input
- **Mock screen size**: Use test harness with controllable dimensions

## Future Enhancements

- **Custom breakpoints**: User-configurable break points
- **Saved layouts**: Remember user's preferred layout
- **Manual override**: Allow user to force layout mode
- **Animations**: Smooth transitions between modes (challenging with tview)
- **Vertical layouts**: Support vertical stacking for tall terminals
- **Grid layout**: Use tview.Grid for more complex arrangements
