# components/sidebar.go

## Purpose
Category tree sidebar using tview.TreeView. Displays hierarchical list of credential categories for navigation.

## Responsibilities

1. **Category display**: Show list of all unique categories
2. **Tree structure**: Hierarchical display with "All Credentials" root
3. **Selection handling**: Track selected category and notify state
4. **Refresh on data changes**: Update tree when credentials change
5. **Visual feedback**: Highlight selected node

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For AppState

### External Dependencies
- `github.com/rivo/tview` - TreeView primitive
- `github.com/gdamore/tcell/v2` - For colors

## Key Types

### `Sidebar`
**Purpose**: Wrapper around tview.TreeView with refresh logic

**Fields**:
```go
type Sidebar struct {
    *tview.TreeView  // Embedded TreeView

    appState *models.AppState  // Reference to state
    rootNode *tview.TreeNode   // Root of tree
}
```

## Key Functions

### Constructor

#### `NewSidebar(appState *models.AppState) *Sidebar`
**Purpose**: Create and configure sidebar

**Steps**:
1. Create tview.TreeView
2. Create root node: "All Credentials"
3. Build initial tree from state
4. Configure styles (borders, colors)
5. Setup selection handler
6. Return wrapped sidebar

### Tree Building

#### `Refresh()`
**Purpose**: Rebuild tree from current AppState

**Steps**:
1. Get categories from appState
2. Clear root node children
3. For each category:
   - Create tree node
   - Set node text as category name
   - Add as child of root
4. Expand root node
5. Restore selection if possible

**Called when**: Credentials are loaded or modified

#### `buildTree(categories []string)`
**Purpose**: Construct tree nodes from category list

**Structure**:
```
All Credentials (root)
├── AWS
├── GitHub
├── Databases
└── Uncategorized
```

### Selection Handling

#### `OnSelectionChanged(handler func(category string))`
**Purpose**: Register callback for selection changes

**When called**: User presses Enter or clicks on node

**Behavior**:
```go
treeView.SetSelectedFunc(func(node *tview.TreeNode) {
    if node != rootNode {
        category := node.GetText()
        appState.SetSelectedCategory(category)
    } else {
        appState.SetSelectedCategory("") // All credentials
    }
})
```

### Styling

#### `applyStyles()`
**Purpose**: Apply theme colors and borders

**Configuration**:
- Rounded borders
- Border color from theme
- Background color from theme
- Text color for selected/unselected nodes

**Example**:
```go
sidebar.SetBorder(true).SetTitle(" Categories ").SetTitleAlign(tview.AlignLeft)
sidebar.SetBorderColor(tcell.NewRGBColor(139, 233, 253)) // Cyan
sidebar.SetBackgroundColor(tcell.NewRGBColor(40, 42, 54))
```

## Example Structure

```go
package components

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/models"
    "pass-cli/cmd/tui-tview/styles"
)

type Sidebar struct {
    *tview.TreeView
    appState *models.AppState
    rootNode *tview.TreeNode
}

func NewSidebar(appState *models.AppState) *Sidebar {
    // Create root node
    root := tview.NewTreeNode("All Credentials").
        SetColor(styles.AccentColor).
        SetSelectable(true)

    // Create tree view
    tree := tview.NewTreeView().
        SetRoot(root).
        SetCurrentNode(root)

    sidebar := &Sidebar{
        TreeView: tree,
        appState: appState,
        rootNode: root,
    }

    // Apply styling
    sidebar.applyStyles()

    // Setup selection handler
    sidebar.SetSelectedFunc(sidebar.onSelect)

    // Initial tree build
    sidebar.Refresh()

    return sidebar
}

func (s *Sidebar) Refresh() {
    categories := s.appState.GetCategories()

    // Clear existing children
    s.rootNode.ClearChildren()

    // Add category nodes
    for _, category := range categories {
        node := tview.NewTreeNode(category).
            SetSelectable(true).
            SetColor(tcell.ColorWhite)
        s.rootNode.AddChild(node)
    }

    // Expand root
    s.rootNode.SetExpanded(true)
}

func (s *Sidebar) onSelect(node *tview.TreeNode) {
    if node == s.rootNode {
        // Root selected - show all
        s.appState.SetSelectedCategory("")
    } else {
        // Category selected
        category := node.GetText()
        s.appState.SetSelectedCategory(category)
    }
}

func (s *Sidebar) applyStyles() {
    s.SetBorder(true).
        SetTitle(" Categories ").
        SetTitleAlign(tview.AlignLeft).
        SetBorderColor(styles.BorderColor).
        SetBackgroundColor(styles.BackgroundColor)
}
```

## Interaction Flow

```
User presses Arrow Down
    ↓
tview handles navigation (built-in)
    ↓
User presses Enter on "AWS" node
    ↓
SetSelectedFunc callback fires
    ↓
onSelect() is called with "AWS" node
    ↓
appState.SetSelectedCategory("AWS")
    ↓
AppState notifies selection changed
    ↓
Table component refreshes to show AWS credentials
```

## Visual Design

### Width
- Fixed: 20 columns (or configurable based on breakpoint)
- Flexible: Use layout manager to adjust

### Height
- Full height minus status bar (calculated by layout manager)

### Borders
- Rounded corners
- Accent color (cyan) for border
- Title: " Categories " (with spaces for padding)

### Content
- Root node in accent color
- Category nodes in white
- Selected node highlighted (tview handles this)

## Keyboard Shortcuts

Handled by tview.TreeView automatically:
- **Arrow Up/Down**: Navigate
- **Enter**: Select category
- **Home**: Jump to first
- **End**: Jump to last

Custom shortcuts (handled in event handler):
- **Tab**: Move focus to table
- **n**: New credential (global)
- **?**: Help (global)

## State Integration

### Reads from AppState
- `GetCategories()` - For building tree

### Writes to AppState
- `SetSelectedCategory(category)` - When node selected

### Responds to AppState callbacks
- `onCredentialsChanged` - Triggers Refresh()

## Edge Cases

1. **No categories**: Show "Uncategorized" only
2. **Empty vault**: Show "No credentials" message
3. **Category deleted**: Refresh tree, select root
4. **Rapid selection changes**: Debounce if needed (unlikely with TUI)

## Testing Considerations

- **Test tree building**: Verify nodes created correctly
- **Test selection**: Verify state is updated
- **Test refresh**: Verify tree rebuilds on data change
- **Mock AppState**: Use mock for isolated testing

## Future Enhancements

- **Search in sidebar**: Filter categories by typing
- **Credential count per category**: Show "(5)" next to each category
- **Collapsible categories**: Multi-level hierarchy
- **Drag and drop**: Move credentials between categories (advanced)
- **Right-click menu**: Context menu for category operations
