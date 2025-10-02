package tui

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/cmd/tui/components"
	"pass-cli/cmd/tui/styles"
	"pass-cli/cmd/tui/views"
)

// generatePassword generates a secure random password
func generatePassword(length int) (string, error) {
	// Use all character types for maximum security
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"

	password := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}

// nextPanelFocus returns the next panel focus in cycle
func (m *Model) nextPanelFocus() PanelFocus {
	visiblePanels := m.getVisiblePanels()
	if len(visiblePanels) == 0 {
		return FocusMain
	}

	// Find current index
	currentIdx := -1
	for i, panel := range visiblePanels {
		if panel == m.panelFocus {
			currentIdx = i
			break
		}
	}

	// Move to next
	nextIdx := (currentIdx + 1) % len(visiblePanels)
	return visiblePanels[nextIdx]
}

// previousPanelFocus returns the previous panel focus in cycle
func (m *Model) previousPanelFocus() PanelFocus {
	visiblePanels := m.getVisiblePanels()
	if len(visiblePanels) == 0 {
		return FocusMain
	}

	// Find current index
	currentIdx := -1
	for i, panel := range visiblePanels {
		if panel == m.panelFocus {
			currentIdx = i
			break
		}
	}

	// Move to previous
	prevIdx := currentIdx - 1
	if prevIdx < 0 {
		prevIdx = len(visiblePanels) - 1
	}
	return visiblePanels[prevIdx]
}

// getVisiblePanels returns a list of currently visible panels
func (m *Model) getVisiblePanels() []PanelFocus {
	panels := []PanelFocus{}

	if m.sidebarVisible {
		panels = append(panels, FocusSidebar)
	}

	// Main is always visible
	panels = append(panels, FocusMain)

	if m.metadataVisible {
		panels = append(panels, FocusMetadata)
	}

	return panels
}

// updatePanelFocus updates the focus state of all panels
func (m *Model) updatePanelFocus() {
	if m.sidebar != nil {
		m.sidebar.SetFocus(m.panelFocus == FocusSidebar)
	}
	if m.metadataPanel != nil {
		m.metadataPanel.SetFocus(m.panelFocus == FocusMetadata)
	}
}

// recalculateLayout calculates panel dimensions and propagates to components
func (m *Model) recalculateLayout() {
	if m.layoutManager == nil || m.width == 0 || m.height == 0 {
		return
	}

	// Get panel states
	states := components.PanelStates{
		SidebarVisible:   m.sidebarVisible,
		MetadataVisible:  m.metadataVisible,
		ProcessVisible:   m.processVisible,
		CommandBarOpen:   m.commandBarOpen,
		StatusBarVisible: true,
	}

	// Calculate layout
	layout := m.layoutManager.Calculate(m.width, m.height, states)

	// Check if too small
	if layout.IsTooSmall {
		// Terminal is too small, will show warning in View()
		return
	}

	// Update sidebar
	if m.sidebar != nil && m.sidebarVisible {
		m.sidebar.SetSize(layout.Sidebar.Width, layout.Sidebar.Height)
	}

	// Update main content (list or detail view)
	if m.listView != nil && m.state == StateList {
		m.listView.SetSize(layout.Main.Width, layout.Main.Height)
	}
	if m.detailView != nil && m.state == StateDetail {
		m.detailView.SetSize(layout.Main.Width, layout.Main.Height)
	}

	// Update metadata panel
	if m.metadataPanel != nil && m.metadataVisible {
		m.metadataPanel.SetSize(layout.Metadata.Width, layout.Metadata.Height)
	}

	// Update process panel
	if m.processPanel != nil && m.processVisible {
		m.processPanel.SetSize(layout.Process.Width, layout.Process.Height)
	}

	// Update command bar
	if m.commandBar != nil && m.commandBarOpen {
		m.commandBar.SetSize(layout.CommandBar.Width, layout.CommandBar.Height)
	}

	// Update breadcrumb
	if m.breadcrumb != nil {
		m.breadcrumb.SetSize(layout.Main.Width)
	}
}

// renderDashboardView renders the multi-panel dashboard layout
func (m *Model) renderDashboardView() string {
	if m.layoutManager == nil {
		// Fallback to simple view if layout manager not initialized
		var mainContent string
		if m.state == StateList && m.listView != nil {
			mainContent = m.listView.View()
		} else if m.state == StateDetail && m.detailView != nil {
			mainContent = m.detailView.View()
		}
		if m.statusBar != nil {
			return mainContent + "\n" + m.statusBar.Render()
		}
		return mainContent
	}

	// Calculate layout
	states := components.PanelStates{
		SidebarVisible:   m.sidebarVisible,
		MetadataVisible:  m.metadataVisible,
		ProcessVisible:   m.processVisible,
		CommandBarOpen:   m.commandBarOpen,
		StatusBarVisible: true,
	}
	layout := m.layoutManager.Calculate(m.width, m.height, states)

	// Check if terminal is too small
	if layout.IsTooSmall {
		warning := "Terminal too small\n\n"
		warning += fmt.Sprintf("Minimum size: %dx%d\n", layout.MinWidth, layout.MinHeight)
		warning += fmt.Sprintf("Current size: %dx%d\n", m.width, m.height)
		warning += "\nPlease resize your terminal"
		return warning
	}

	// Render main content area
	mainContent := ""
	if m.state == StateList && m.listView != nil {
		mainContent = m.listView.View()
	} else if m.state == StateDetail && m.detailView != nil {
		mainContent = m.detailView.View()
	}

	// Apply panel border styling
	mainPanelStyle := styles.InactivePanelBorderStyle
	if m.panelFocus == FocusMain {
		mainPanelStyle = styles.ActivePanelBorderStyle
	}
	mainPanel := mainPanelStyle.Width(layout.Main.Width).Height(layout.Main.Height).Render(mainContent)

	// Collect horizontal panels (sidebar, main, metadata)
	var horizontalPanels []string

	// Render sidebar if visible
	if m.sidebarVisible && m.sidebar != nil {
		sidebarContent := m.sidebar.View()
		sidebarStyle := styles.InactivePanelBorderStyle
		if m.panelFocus == FocusSidebar {
			sidebarStyle = styles.ActivePanelBorderStyle
		}
		sidebarPanel := sidebarStyle.Width(layout.Sidebar.Width).Height(layout.Sidebar.Height).Render(sidebarContent)
		horizontalPanels = append(horizontalPanels, sidebarPanel)
	}

	// Main panel is always visible
	horizontalPanels = append(horizontalPanels, mainPanel)

	// Render metadata panel if visible
	if m.metadataVisible && m.metadataPanel != nil {
		metadataContent := m.metadataPanel.View()
		metadataStyle := styles.InactivePanelBorderStyle
		if m.panelFocus == FocusMetadata {
			metadataStyle = styles.ActivePanelBorderStyle
		}
		metadataPanel := metadataStyle.Width(layout.Metadata.Width).Height(layout.Metadata.Height).Render(metadataContent)
		horizontalPanels = append(horizontalPanels, metadataPanel)
	}

	// Join horizontal panels
	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, horizontalPanels...)

	// Build vertical layout
	var verticalSections []string
	verticalSections = append(verticalSections, mainRow)

	// Add process panel if visible
	if m.processVisible && m.processPanel != nil {
		processContent := m.processPanel.View()
		if processContent != "" {
			verticalSections = append(verticalSections, processContent)
		}
	}

	// Add command bar if open
	if m.commandBarOpen && m.commandBar != nil {
		commandBarContent := m.commandBar.View()
		verticalSections = append(verticalSections, commandBarContent)
	}

	// Add status bar
	if m.statusBar != nil {
		verticalSections = append(verticalSections, m.statusBar.Render())
	}

	// Join all vertical sections
	return lipgloss.JoinVertical(lipgloss.Left, verticalSections...)
}

// executeCommand executes a parsed command from the command bar
func (m *Model) executeCommand(cmd *components.Command) (tea.Model, tea.Cmd) {
	switch cmd.Name {
	case "add", "a":
		// Open add form
		m.state = StateAdd
		m.addForm = views.NewAddFormView()
		// TODO: Pre-fill service if args provided: if len(cmd.Args) > 0 { m.addForm.SetService(cmd.Args[0]) }
		m.addForm.SetSize(m.width, m.height)
		return m, nil

	case "search":
		// Switch to list view (search is activated with '/')
		if m.state == StateDetail {
			m.state = StateList
			m.detailView = nil
		}
		// TODO: Activate search with query: if len(cmd.Args) > 0 { m.listView.ActivateSearch(cmd.Args[0]) }
		return m, nil

	case "category", "cat":
		// Navigate to specific category in sidebar
		if len(cmd.Args) > 0 && m.sidebar != nil {
			// TODO: Find category and expand it
			// For now, just switch to sidebar focus
			m.panelFocus = FocusSidebar
			m.sidebarVisible = true
			m.updatePanelFocus()
		}
		return m, nil

	case "help", "h":
		// Show help
		m.state = StateHelp
		if m.helpView == nil {
			m.helpView = views.NewHelpView()
		}
		m.helpView.SetSize(m.width, m.height)
		return m, nil

	case "quit", "q":
		// Quit the application
		return m, tea.Quit

	default:
		// Unknown command
		m.commandBar.SetError("Unknown command: " + cmd.Name)
		return m, nil
	}
}
