package components

import (
	"fmt"
	"strings"

	"pass-cli/cmd/tui/styles"
	"pass-cli/internal/vault"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SidebarPanel represents the left sidebar with category navigation
type SidebarPanel struct {
	categories       []Category
	selectedCategory int
	selectedCred     int // Index within the selected category's credentials
	stats            Stats
	viewport         viewport.Model
	width            int
	height           int
	focused          bool
}

// Stats represents vault statistics
type Stats struct {
	Total  int
	Used   int
	Recent int
}

// NewSidebarPanel creates a new sidebar panel
func NewSidebarPanel(credentials []vault.CredentialMetadata) *SidebarPanel {
	vp := viewport.New(20, 10)
	categories := CategorizeCredentials(credentials)

	// Auto-expand first category if it has items
	if len(categories) > 0 && categories[0].Count > 0 {
		categories[0].Expanded = true
	}

	stats := Stats{
		Total:  len(credentials),
		Used:   countUsedCredentials(credentials),
		Recent: countRecentCredentials(credentials),
	}

	return &SidebarPanel{
		categories:       categories,
		selectedCategory: 0,
		selectedCred:     -1, // -1 means category header is selected
		stats:            stats,
		viewport:         vp,
		focused:          false,
	}
}

// SetSize updates the panel dimensions
func (s *SidebarPanel) SetSize(width, height int) {
	s.width = width
	s.height = height

	// Reserve space for title (1), stats (4), quick actions (4), borders/padding (2)
	contentHeight := height - 11
	if contentHeight < 5 {
		contentHeight = 5
	}

	s.viewport.Width = width - 2 // Account for padding
	s.viewport.Height = contentHeight
}

// SetFocus sets the focus state of the panel
func (s *SidebarPanel) SetFocus(focused bool) {
	s.focused = focused
}

// Update handles tea messages
func (s *SidebarPanel) Update(msg tea.Msg) (*SidebarPanel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !s.focused {
			return s, nil
		}

		switch msg.String() {
		case "j", "down":
			s.moveDown()
		case "k", "up":
			s.moveUp()
		case "enter", "l", "right":
			// Toggle expansion or select credential
			if s.selectedCred == -1 {
				// On category header, toggle expansion
				if s.selectedCategory >= 0 && s.selectedCategory < len(s.categories) {
					s.categories[s.selectedCategory].Expanded = !s.categories[s.selectedCategory].Expanded
				}
			}
		case "h", "left":
			// Collapse current category
			if s.selectedCategory >= 0 && s.selectedCategory < len(s.categories) {
				s.categories[s.selectedCategory].Expanded = false
				s.selectedCred = -1 // Go back to header
			}
		}
	}

	s.updateViewportContent()
	return s, nil
}

// moveDown moves selection down
func (s *SidebarPanel) moveDown() {
	if len(s.categories) == 0 {
		return
	}

	currentCat := &s.categories[s.selectedCategory]

	if s.selectedCred == -1 {
		// On category header
		if currentCat.Expanded && currentCat.Count > 0 {
			// Move to first credential
			s.selectedCred = 0
		} else {
			// Move to next category
			if s.selectedCategory < len(s.categories)-1 {
				s.selectedCategory++
				s.selectedCred = -1
			}
		}
	} else {
		// On credential
		if s.selectedCred < currentCat.Count-1 {
			// Move to next credential
			s.selectedCred++
		} else {
			// Move to next category
			if s.selectedCategory < len(s.categories)-1 {
				s.selectedCategory++
				s.selectedCred = -1
			}
		}
	}
}

// moveUp moves selection up
func (s *SidebarPanel) moveUp() {
	if len(s.categories) == 0 {
		return
	}

	if s.selectedCred == -1 {
		// On category header, move to previous category
		if s.selectedCategory > 0 {
			s.selectedCategory--
			// If previous category is expanded, select last credential
			prevCat := &s.categories[s.selectedCategory]
			if prevCat.Expanded && prevCat.Count > 0 {
				s.selectedCred = prevCat.Count - 1
			}
		}
	} else {
		// On credential
		if s.selectedCred > 0 {
			// Move to previous credential
			s.selectedCred--
		} else {
			// Move to category header
			s.selectedCred = -1
		}
	}
}

// GetSelectedCredential returns the currently selected credential, or nil if on header
func (s *SidebarPanel) GetSelectedCredential() *vault.CredentialMetadata {
	if s.selectedCred == -1 || s.selectedCategory >= len(s.categories) {
		return nil
	}

	cat := s.categories[s.selectedCategory]
	if s.selectedCred >= 0 && s.selectedCred < len(cat.Credentials) {
		return &cat.Credentials[s.selectedCred]
	}

	return nil
}

// GetSelectedCategory returns the currently selected category name
func (s *SidebarPanel) GetSelectedCategory() string {
	if s.selectedCategory >= 0 && s.selectedCategory < len(s.categories) {
		return s.categories[s.selectedCategory].Name
	}
	return ""
}

// UpdateCredentials updates the sidebar with new credentials
func (s *SidebarPanel) UpdateCredentials(credentials []vault.CredentialMetadata) {
	s.categories = CategorizeCredentials(credentials)

	// Auto-expand first category if it has items
	if len(s.categories) > 0 && s.categories[0].Count > 0 {
		s.categories[0].Expanded = true
	}

	s.stats = Stats{
		Total:  len(credentials),
		Used:   countUsedCredentials(credentials),
		Recent: countRecentCredentials(credentials),
	}

	// Reset selection if out of bounds
	if s.selectedCategory >= len(s.categories) {
		s.selectedCategory = 0
		s.selectedCred = -1
	}

	s.updateViewportContent()
}

// updateViewportContent refreshes the viewport with current category tree
func (s *SidebarPanel) updateViewportContent() {
	content := s.renderCategoryTree()
	s.viewport.SetContent(content)
}

// renderCategoryTree renders the category tree as a string
func (s *SidebarPanel) renderCategoryTree() string {
	if len(s.categories) == 0 {
		return styles.SubtleStyle.Render("No credentials yet")
	}

	var lines []string

	for catIdx, cat := range s.categories {
		// Render category header
		expandIcon := styles.GetStatusIcon("collapsed")
		if cat.Expanded {
			expandIcon = styles.GetStatusIcon("expanded")
		}

		categoryLine := fmt.Sprintf("%s %s %s (%d)",
			expandIcon,
			cat.Icon,
			cat.Name,
			cat.Count,
		)

		// Highlight if selected
		if catIdx == s.selectedCategory && s.selectedCred == -1 {
			if s.focused {
				categoryLine = styles.SelectedStyle.Render(categoryLine)
			} else {
				categoryLine = styles.FocusedLabelStyle.Render(categoryLine)
			}
		} else {
			categoryLine = styles.ValueStyle.Render(categoryLine)
		}

		lines = append(lines, categoryLine)

		// Render credentials if expanded
		if cat.Expanded {
			for credIdx, cred := range cat.Credentials {
				credLine := fmt.Sprintf("  • %s", cred.Service)
				if len(cred.Username) > 0 {
					credLine += fmt.Sprintf(" (%s)", truncate(cred.Username, 15))
				}

				// Highlight if selected
				if catIdx == s.selectedCategory && credIdx == s.selectedCred {
					if s.focused {
						credLine = styles.SelectedStyle.Render(credLine)
					} else {
						credLine = styles.FocusedLabelStyle.Render(credLine)
					}
				} else {
					credLine = styles.SubtleStyle.Render(credLine)
				}

				lines = append(lines, credLine)
			}
		}
	}

	return strings.Join(lines, "\n")
}

// View renders the sidebar panel
func (s *SidebarPanel) View() string {
	titleStyle := styles.InactivePanelTitleStyle
	if s.focused {
		titleStyle = styles.PanelTitleStyle
	}

	title := titleStyle.Render("📋 Categories")

	// Stats section
	statsLines := []string{
		"",
		styles.LabelStyle.Render("Vault Statistics:"),
		fmt.Sprintf("  Total: %d", s.stats.Total),
		fmt.Sprintf("  Used: %d", s.stats.Used),
		fmt.Sprintf("  Recent: %d", s.stats.Recent),
	}

	// Quick actions
	actionLines := []string{
		"",
		styles.LabelStyle.Render("Quick Actions:"),
		styles.SubtleStyle.Render("  [a] Add"),
		styles.SubtleStyle.Render("  [:] Command"),
		styles.SubtleStyle.Render("  [?] Help"),
	}

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		s.viewport.View(),
		strings.Join(statsLines, "\n"),
		strings.Join(actionLines, "\n"),
	)

	return content
}

// Helper functions

func countUsedCredentials(credentials []vault.CredentialMetadata) int {
	count := 0
	for _, cred := range credentials {
		if cred.UsageCount > 0 {
			count++
		}
	}
	return count
}

func countRecentCredentials(credentials []vault.CredentialMetadata) int {
	// Consider credentials accessed in last 7 days as recent
	// For now, return count of credentials with non-zero usage
	// This can be enhanced with actual time-based filtering
	return countUsedCredentials(credentials)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
