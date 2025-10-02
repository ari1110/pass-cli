package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pass-cli/cmd/tui/styles"
)

// Shortcut represents a keyboard shortcut with its description
type Shortcut struct {
	Key    string
	Action string
}

// HelpView displays keyboard shortcuts organized by category
type HelpView struct {
	viewport viewport.Model
	width    int
	height   int
}

// NewHelpView creates a new help overlay
func NewHelpView() *HelpView {
	vp := viewport.New(0, 0)

	return &HelpView{
		viewport: vp,
	}
}

// SetSize updates dimensions
func (h *HelpView) SetSize(width, height int) {
	h.width = width
	h.height = height

	// Size the viewport for modal display
	modalWidth := 80
	if modalWidth > width-10 {
		modalWidth = width - 10
	}

	modalHeight := height - 10
	if modalHeight < 10 {
		modalHeight = 10
	}

	h.viewport.Width = modalWidth - 4  // Account for padding
	h.viewport.Height = modalHeight - 4

	// Update content
	h.updateContent()
}

// Update handles messages
func (h *HelpView) Update(msg tea.Msg) (*HelpView, tea.Cmd) {
	var cmd tea.Cmd
	h.viewport, cmd = h.viewport.Update(msg)
	return h, cmd
}

// View renders the help overlay
func (h *HelpView) View() string {
	// Modal styling
	modalStyle := styles.ModalBorderStyle.
		Width(h.viewport.Width + 4).
		Align(lipgloss.Left)

	modal := modalStyle.Render(h.viewport.View())

	// Dim background by placing modal in center
	dimStyle := styles.SubtleStyle

	overlay := lipgloss.Place(
		h.width,
		h.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)

	return dimStyle.Render(overlay)
}

// updateContent generates the help content
func (h *HelpView) updateContent() {
	var content strings.Builder

	// Title
	titleStyle := styles.TitleStyle.Align(lipgloss.Center)

	content.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	content.WriteString("\n\n")

	// Global shortcuts
	content.WriteString(h.renderSection("Global", []Shortcut{
		{"q", "Quit application"},
		{"? or F1", "Show this help"},
		{"Esc", "Back/Cancel"},
		{"Ctrl+C", "Quit"},
	}))

	// Navigation
	content.WriteString(h.renderSection("Navigation", []Shortcut{
		{"↑/↓ or j/k", "Move up/down"},
		{"←/→ or h/l", "Move left/right"},
		{"Tab", "Next field/focus"},
		{"Shift+Tab", "Previous field/focus"},
		{"Enter", "Select/Confirm"},
	}))

	// List View
	content.WriteString(h.renderSection("List View", []Shortcut{
		{"/", "Focus search bar"},
		{"Esc (in search)", "Clear search"},
		{"a", "Add new credential"},
		{"Enter", "View credential details"},
	}))

	// Detail View
	content.WriteString(h.renderSection("Detail View", []Shortcut{
		{"m", "Toggle password visibility"},
		{"c", "Copy password to clipboard"},
		{"e", "Edit credential"},
		{"d", "Delete credential"},
		{"Esc", "Back to list"},
	}))

	// Forms (Add/Edit)
	content.WriteString(h.renderSection("Add/Edit Forms", []Shortcut{
		{"Tab/↓", "Next field"},
		{"Shift+Tab/↑", "Previous field"},
		{"g", "Generate secure password"},
		{"Ctrl+S", "Save credential"},
		{"Esc", "Cancel (with confirmation)"},
	}))

	// Confirmation Dialogs
	content.WriteString(h.renderSection("Confirmation Dialogs", []Shortcut{
		{"y", "Yes, confirm"},
		{"n", "No, cancel"},
		{"Enter", "Confirm (typed mode)"},
		{"Esc", "Cancel"},
	}))

	// Footer
	content.WriteString("\n")
	footerStyle := styles.SubtleStyle.
		Italic(true).
		Align(lipgloss.Center)

	content.WriteString(footerStyle.Render("Press any key to close"))

	h.viewport.SetContent(content.String())
}

// renderSection renders a category of shortcuts
func (h *HelpView) renderSection(title string, shortcuts []Shortcut) string {
	var section strings.Builder

	// Section title
	sectionStyle := styles.SubtitleStyle.
		MarginTop(1).
		MarginBottom(1)

	section.WriteString(sectionStyle.Render(title))
	section.WriteString("\n")

	// Shortcuts in two-column format
	keyStyle := styles.KeyStyle.Width(20)
	actionStyle := styles.ValueStyle

	for _, sc := range shortcuts {
		key := keyStyle.Render("  " + sc.Key)
		action := actionStyle.Render(sc.Action)
		section.WriteString(fmt.Sprintf("%s %s\n", key, action))
	}

	section.WriteString("\n")
	return section.String()
}
