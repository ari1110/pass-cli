package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"pass-cli/cmd/tui-tview/models"
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
	// Create root node
	root := tview.NewTreeNode("All Credentials").
		SetColor(tcell.NewRGBColor(139, 233, 253)). // Cyan accent color
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
	// Get categories from state (thread-safe read)
	categories := s.appState.GetCategories()

	// Clear existing children
	s.rootNode.ClearChildren()

	// Add category nodes
	for _, category := range categories {
		node := tview.NewTreeNode(category).
			SetSelectable(true).
			SetColor(tcell.NewRGBColor(248, 248, 242)) // White text
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
	s.SetBorder(true).
		SetTitle(" Categories ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.NewRGBColor(139, 233, 253)). // Cyan border
		SetBackgroundColor(tcell.NewRGBColor(40, 42, 54))  // Dark background
}
