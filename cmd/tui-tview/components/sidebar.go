package components

import (
	"github.com/rivo/tview"
	"pass-cli/cmd/tui-tview/models"
	"pass-cli/cmd/tui-tview/styles"
)

// Sidebar wraps tview.TreeView to display credential categories.
// Provides category navigation with "All Credentials" root and category children.
type Sidebar struct {
	*tview.TreeView

	appState *models.AppState
	rootNode *tview.TreeNode
}

// NewSidebar creates and configures a new Sidebar component.
// Creates TreeView with root "All Credentials" node and builds initial tree.
func NewSidebar(appState *models.AppState) *Sidebar {
	theme := styles.GetCurrentTheme()

	// Create root node
	root := tview.NewTreeNode("All Credentials").
		SetColor(theme.BorderColor). // Cyan accent color
		SetSelectable(true).
		SetExpanded(true)

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

// Refresh rebuilds the category tree from current AppState.
// Clears existing children and adds category nodes from appState.GetCategories().
func (s *Sidebar) Refresh() {
	theme := styles.GetCurrentTheme()

	// Get categories from state (thread-safe read)
	categories := s.appState.GetCategories()

	// Clear existing children
	s.rootNode.ClearChildren()

	// Add category nodes
	for _, category := range categories {
		node := tview.NewTreeNode(category).
			SetSelectable(true).
			SetColor(theme.TextPrimary) // White text
		s.rootNode.AddChild(node)
	}

	// Ensure root is expanded
	s.rootNode.SetExpanded(true)
}

// onSelect handles node selection by updating AppState.
// Root node selection shows all credentials, category nodes filter by category.
func (s *Sidebar) onSelect(node *tview.TreeNode) {
	if node == s.rootNode {
		// Root selected - show all credentials
		s.appState.SetSelectedCategory("")
	} else {
		// Category selected - filter by category
		category := node.GetText()
		s.appState.SetSelectedCategory(category)
	}
}

// applyStyles applies borders, colors, and title to the sidebar.
// Uses rounded borders with cyan accent color and dark background.
func (s *Sidebar) applyStyles() {
	styles.ApplyBorderedStyle(s.TreeView, "Categories", true)
}
