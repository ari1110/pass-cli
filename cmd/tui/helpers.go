package tui

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

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

	// Update sidebar (using ContentWidth/ContentHeight for bordered panels)
	if m.sidebar != nil && m.sidebarVisible {
		m.sidebar.SetSize(layout.Sidebar.ContentWidth, layout.Sidebar.ContentHeight)
	}

	// Update main content (list or detail view) (using ContentWidth/ContentHeight for bordered panels)
	if m.listView != nil && m.state == StateList {
		m.listView.SetSize(layout.Main.ContentWidth, layout.Main.ContentHeight)
	}
	if m.detailView != nil && m.state == StateDetail {
		m.detailView.SetSize(layout.Main.ContentWidth, layout.Main.ContentHeight)
	}

	// Update metadata panel (using ContentWidth/ContentHeight for bordered panels)
	if m.metadataPanel != nil && m.metadataVisible {
		m.metadataPanel.SetSize(layout.Metadata.ContentWidth, layout.Metadata.ContentHeight)
	}

	// Update process panel (non-bordered, ContentWidth == Width)
	if m.processPanel != nil && m.processVisible {
		m.processPanel.SetSize(layout.Process.ContentWidth, layout.Process.ContentHeight)
	}

	// Update command bar (non-bordered, ContentWidth == Width)
	if m.commandBar != nil && m.commandBarOpen {
		m.commandBar.SetSize(layout.CommandBar.ContentWidth, layout.CommandBar.ContentHeight)
	}

	// Update breadcrumb (uses main panel's content width)
	if m.breadcrumb != nil {
		m.breadcrumb.SetSize(layout.Main.ContentWidth)
	}
}

// renderDashboardView renders the multi-panel dashboard.
// Panels use Lipgloss GetFrameSize() to calculate border/padding overhead.
//
// IMPORTANT: Lipgloss .Width(n) sets content width and adds frame ON TOP,
// so we subtract GetHorizontalFrameSize() before applying .Width() to ensure
// total rendered width equals allocated layout width. Same applies for height.
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
	// Lipgloss .Width() adds border but not padding, so subtract padding only
	// GetFrameSize()=4 (border=2 + padding=2), but .Width() already accounts for border
	// So we subtract padding: Width - 2
	mainPanelWidth := layout.Main.Width - 2
	mainPanelHeight := layout.Main.Height - 2
	if mainPanelWidth < 10 {
		mainPanelWidth = 10
	}
	if mainPanelHeight < 5 {
		mainPanelHeight = 5
	}
	mainPanel := mainPanelStyle.Width(mainPanelWidth).Height(mainPanelHeight).Render(mainContent)

	// Collect horizontal panels (sidebar, main, metadata)
	var horizontalPanels []string

	// Render sidebar if visible
	if m.sidebarVisible && m.sidebar != nil {
		sidebarContent := m.sidebar.View()
		sidebarStyle := styles.InactivePanelBorderStyle
		if m.panelFocus == FocusSidebar {
			sidebarStyle = styles.ActivePanelBorderStyle
		}
		// Lipgloss .Width() adds border but not padding, so subtract padding only (Width - 2)
		sidebarWidth := layout.Sidebar.Width - 2
		sidebarHeight := layout.Sidebar.Height - 2
		if sidebarWidth < 10 {
			sidebarWidth = 10
		}
		if sidebarHeight < 5 {
			sidebarHeight = 5
		}
		sidebarPanel := sidebarStyle.Width(sidebarWidth).Height(sidebarHeight).Render(sidebarContent)
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
		// Lipgloss .Width() adds border but not padding, so subtract padding only (Width - 2)
		metadataWidth := layout.Metadata.Width - 2
		metadataHeight := layout.Metadata.Height - 2
		if metadataWidth < 10 {
			metadataWidth = 10
		}
		if metadataHeight < 5 {
			metadataHeight = 5
		}
		metadataPanel := metadataStyle.Width(metadataWidth).Height(metadataHeight).Render(metadataContent)
		horizontalPanels = append(horizontalPanels, metadataPanel)
	}

	// Join horizontal panels
	mainRow := lipgloss.JoinHorizontal(lipgloss.Top, horizontalPanels...)

	// Build vertical layout
	var verticalSections []string

	// Add breadcrumb at the top (if in detail view and we have a path)
	if m.state == StateDetail && m.breadcrumb != nil {
		breadcrumbView := m.breadcrumb.View()
		if breadcrumbView != "" {
			verticalSections = append(verticalSections, breadcrumbView)
		}
	}

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

	// Join all vertical sections and constrain to terminal dimensions
	output := lipgloss.JoinVertical(lipgloss.Left, verticalSections...)

	// Use a style to constrain the entire output to terminal size
	return lipgloss.NewStyle().
		MaxWidth(m.width).
		MaxHeight(m.height).
		Render(output)
}

// getPanelShortcuts returns the panel toggle shortcuts for status bar
func (m *Model) getPanelShortcuts() string {
	shortcuts := []string{}

	// Always show sidebar toggle
	if m.sidebarVisible {
		shortcuts = append(shortcuts, "s: hide sidebar")
	} else {
		shortcuts = append(shortcuts, "s: show sidebar")
	}

	// Show metadata/info toggle in Detail state
	if m.state == StateDetail {
		if m.metadataVisible {
			shortcuts = append(shortcuts, "i: hide info")
		} else {
			shortcuts = append(shortcuts, "i: show info")
		}
	}

	// Show process toggle if there are active processes
	if m.processPanel != nil && m.processPanel.HasActiveProcesses() {
		if m.processVisible {
			shortcuts = append(shortcuts, "p: hide processes")
		} else {
			shortcuts = append(shortcuts, "p: show processes")
		}
	}

	// Show tab navigation if multiple panels visible
	visiblePanels := m.getVisiblePanels()
	if len(visiblePanels) > 1 {
		shortcuts = append(shortcuts, "tab: switch panel")
	}

	if len(shortcuts) == 0 {
		return ""
	}

	return strings.Join(shortcuts, " | ")
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
