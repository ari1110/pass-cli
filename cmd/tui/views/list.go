package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pass-cli/internal/vault"
)

// credentialItem implements list.Item for the bubbles list component
type credentialItem struct {
	metadata vault.CredentialMetadata
}

func (i credentialItem) FilterValue() string {
	return i.metadata.Service + " " + i.metadata.Username
}

func (i credentialItem) Title() string {
	return i.metadata.Service
}

func (i credentialItem) Description() string {
	username := i.metadata.Username
	if username == "" {
		username = "(no username)"
	}
	return username
}

// ListView manages the credential list view
type ListView struct {
	list        list.Model
	searchInput textinput.Model
	searchFocused bool
	allItems    []list.Item
	width       int
	height      int
}

// NewListView creates a new list view
func NewListView(credentials []vault.CredentialMetadata) *ListView {
	// Convert credentials to list items
	items := make([]list.Item, len(credentials))
	for i, cred := range credentials {
		items[i] = credentialItem{metadata: cred}
	}

	// Create list with custom delegate for better spacing
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0) // Reduce spacing between items

	l := list.New(items, delegate, 0, 0)
	l.Title = "Credentials"
	l.SetShowStatusBar(false) // Hide status bar to save space
	l.SetFilteringEnabled(false) // We have our own search
	l.SetShowHelp(false) // We show our own help
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")) // Cyan

	// Create search input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100

	return &ListView{
		list:        l,
		searchInput: ti,
		allItems:    items,
	}
}

// SetSize updates the dimensions
func (v *ListView) SetSize(width, height int) {
	v.width = width
	v.height = height

	// Reserve space for title (1), search bar (3 with border), help (1), status bar (1)
	listHeight := height - 6
	if listHeight < 5 {
		listHeight = 5
	}

	// Set list width (full width available)
	v.list.SetSize(width, listHeight)

	// Set search input width accounting for border (2) + padding (2)
	searchWidth := width - 6
	if searchWidth < 20 {
		searchWidth = 20
	}
	v.searchInput.Width = searchWidth
}

// Update handles messages
func (v *ListView) Update(msg tea.Msg) (*ListView, tea.Cmd) {
	var cmd tea.Cmd

	if v.searchFocused {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				v.searchFocused = false
				v.searchInput.SetValue("")
				v.list.SetItems(v.allItems)
				return v, nil
			case "enter", "tab":
				v.searchFocused = false
				return v, nil
			}
		}

		v.searchInput, cmd = v.searchInput.Update(msg)

		// Filter list based on search query
		query := strings.ToLower(v.searchInput.Value())
		if query == "" {
			v.list.SetItems(v.allItems)
		} else {
			filtered := []list.Item{}
			for _, item := range v.allItems {
				filterVal := strings.ToLower(item.FilterValue())
				if strings.Contains(filterVal, query) {
					filtered = append(filtered, item)
				}
			}
			v.list.SetItems(filtered)
		}

		return v, cmd
	}

	// Handle list navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/", "tab":
			v.searchFocused = true
			v.searchInput.Focus()
			return v, textinput.Blink
		}
	}

	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

// View renders the list view
func (v *ListView) View() string {
	searchBar := v.renderSearchBar()

	var listView string
	if len(v.list.Items()) == 0 && v.searchInput.Value() != "" {
		// Show "no results" message when search has no matches
		noResultsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Padding(2, 0)
		listView = noResultsStyle.Render("No credentials found matching '" + v.searchInput.Value() + "'")
	} else {
		listView = v.list.View()
	}

	help := v.renderHelp()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		searchBar,
		listView,
		help,
	)
}

// renderSearchBar renders the search input
func (v *ListView) renderSearchBar() string {
	if v.searchFocused {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("6")). // Cyan when focused
			Padding(0, 1).
			Render(v.searchInput.View())
	}

	placeholder := "Press / to search"
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")). // Gray when not focused
		Padding(0, 1).
		Render(placeholder)
}

// renderHelp renders the help line
func (v *ListView) renderHelp() string {
	if v.searchFocused {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("esc: clear search | tab/enter: back to list")
	}

	help := "a: add | tab//: search | ↑↓/jk: navigate | enter: view | q: quit"
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(help)
}

// SelectedCredential returns the currently selected credential
func (v *ListView) SelectedCredential() *vault.CredentialMetadata {
	item := v.list.SelectedItem()
	if item == nil {
		return nil
	}

	credItem, ok := item.(credentialItem)
	if !ok {
		return nil
	}

	return &credItem.metadata
}

// IsSearchFocused returns whether the search input is currently focused
func (v *ListView) IsSearchFocused() bool {
	return v.searchFocused
}
