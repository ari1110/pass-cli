package components

import (
	"sort"

	"github.com/rivo/tview"
	"pass-cli/cmd/tui-tview/models"
	"pass-cli/cmd/tui-tview/styles"
	"pass-cli/internal/vault"
)

// NodeReference identifies the type and value of a tree node.
// Used to distinguish categories from credentials without relying on tree position.
type NodeReference struct {
	Kind  string // "category" or "credential"
	Value string // Category name or service name
}

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
// Clears existing children and builds category-grouped tree with credential nodes.
func (s *Sidebar) Refresh() {
	theme := styles.GetCurrentTheme()

	// Get credentials from state (single snapshot for consistency)
	credentials := s.appState.GetCredentials()

	// Clear existing children
	s.rootNode.ClearChildren()

	// Pre-group credentials by category for O(N+C) performance instead of O(C×N)
	groups := make(map[string][]vault.CredentialMetadata)
	for _, cred := range credentials {
		category := cred.Category
		if category == "" {
			category = "Uncategorized"
		}
		groups[category] = append(groups[category], cred)
	}

	// Build category list from groups (avoids snapshot mismatch with credentials)
	categories := make([]string, 0, len(groups))
	for category := range groups {
		categories = append(categories, category)
	}
	sort.Strings(categories) // Sort categories alphabetically

	// Build category nodes with credential children
	for _, category := range categories {
		// Create category node with NodeReference
		categoryNode := tview.NewTreeNode(category).
			SetSelectable(true).
			SetColor(theme.TextPrimary). // White text
			SetReference(NodeReference{Kind: "category", Value: category}).
			SetExpanded(false) // Collapsed by default

		// Sort credentials within category for deterministic ordering
		credList := groups[category]
		sort.Slice(credList, func(i, j int) bool {
			return credList[i].Service < credList[j].Service
		})

		// Add credential nodes from sorted list
		for _, cred := range credList {
			// Create credential node with NodeReference
			credNode := tview.NewTreeNode(cred.Service).
				SetSelectable(true).
				SetColor(theme.TextSecondary). // Gray text to distinguish from category
				SetReference(NodeReference{Kind: "credential", Value: cred.Service})

			categoryNode.AddChild(credNode)
		}

		// Add category node to root
		s.rootNode.AddChild(categoryNode)
	}

	// Ensure root is expanded
	s.rootNode.SetExpanded(true)
}

// onSelect handles node selection by updating AppState.
// Root node shows all, category nodes filter by category, credential nodes select specific credential.
func (s *Sidebar) onSelect(node *tview.TreeNode) {
	if node == s.rootNode {
		// Root selected - show all credentials and clear detail view
		s.appState.SetSelectedCategory("")
		s.appState.SetSelectedCredential(nil)
		return
	}

	// Get node reference to determine type
	ref := node.GetReference()
	if ref == nil {
		// Safety fallback - treat as root
		s.appState.SetSelectedCategory("")
		s.appState.SetSelectedCredential(nil)
		return
	}

	// Type assert to NodeReference and switch on Kind
	if nodeRef, ok := ref.(NodeReference); ok {
		switch nodeRef.Kind {
		case "category":
			// Category node - filter by category and clear credential selection
			s.appState.SetSelectedCategory(nodeRef.Value)
			s.appState.SetSelectedCredential(nil)

		case "credential":
			// Credential node - lookup credential by service and select it
			if credMeta, found := s.appState.FindCredentialByService(nodeRef.Value); found {
				// Select specific credential (fresh lookup avoids stale pointers)
				s.appState.SetSelectedCredential(credMeta)
			}

		default:
			// Unknown kind - treat as root
			s.appState.SetSelectedCategory("")
			s.appState.SetSelectedCredential(nil)
		}
	}
}

// applyStyles applies borders, colors, and title to the sidebar.
// Uses rounded borders with cyan accent color and dark background.
func (s *Sidebar) applyStyles() {
	styles.ApplyBorderedStyle(s.TreeView, "Categories", true)
}
