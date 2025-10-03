package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Dashbrew-inspired playground model
type playgroundModel struct {
	width       int
	height      int
	currentPage int // 0 = background, 1 = two panels (1:2), 2 = three panels (1:2:1)
}

func initialModel() playgroundModel {
	return playgroundModel{}
}

func (m playgroundModel) Init() tea.Cmd {
	return nil
}

func (m playgroundModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "right", "l":
			if m.currentPage < 2 {
				m.currentPage++
			}
		case "left", "h":
			if m.currentPage > 0 {
				m.currentPage--
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m playgroundModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	// Simple three-section layout:
	// 1. Header (1 line)
	// 2. Content (height - 2 lines)
	// 3. Footer (1 line)

	// Header - purple background
	var headerText string
	switch m.currentPage {
	case 0:
		headerText = "🧪 Dashbrew Playground - Page 0: Background Only"
	case 1:
		headerText = "🧪 Dashbrew Playground - Page 1: Two Panels (1:2 Ratio)"
	case 2:
		headerText = "🧪 Dashbrew Playground - Page 2: Three Panels (1:2:1 Ratio)"
	}

	header := lipgloss.NewStyle().
		Background(lipgloss.Color("53")).
		Foreground(lipgloss.Color("15")).
		Width(m.width).
		Bold(true).
		Align(lipgloss.Center).
		Render(headerText)

	// Content - dark gray background
	contentHeight := m.height - 2
	content := m.renderContent(contentHeight)

	// Footer - medium gray background
	footer := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("15")).
		Width(m.width).
		Align(lipgloss.Center).
		Render("h/left: Previous | l/right: Next | q: Quit")

	// Combine all sections
	return header + "\n" + content + "\n" + footer
}

func (m playgroundModel) renderContent(contentHeight int) string {
	if m.currentPage == 0 {
		// Page 0: Just background
		return lipgloss.NewStyle().
			Background(lipgloss.Color("234")).
			Width(m.width).
			Height(contentHeight).
			Render("")
	}

	if m.currentPage == 1 {
		// Page 1: Two-panel flex layout (1:2 ratio)
		return m.renderTwoPanelLayout(contentHeight)
	}

	// Page 2: Three-panel flex layout (1:2:1 ratio)
	return m.renderThreePanelLayout(contentHeight)
}

// renderTwoPanelLayout demonstrates two-panel flex layout (1:2 ratio)
func (m playgroundModel) renderTwoPanelLayout(contentHeight int) string {
	// Flex values
	leftFlex := 1
	rightFlex := 2
	totalFlex := leftFlex + rightFlex

	// Calculate widths based on flex ratio
	leftWidth := m.width * leftFlex / totalFlex
	rightWidth := m.width - leftWidth // Remaining space to avoid rounding issues

	// Create a sample style to get frame size
	sampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)
	horizontalFrame := sampleStyle.GetHorizontalFrameSize()
	verticalFrame := sampleStyle.GetVerticalFrameSize()

	// Create left panel - Width/Height are CONTENT dimensions
	leftContentWidth := leftWidth - horizontalFrame
	leftContentHeight := contentHeight - verticalFrame

	leftPanel := lipgloss.NewStyle().
		Width(leftContentWidth).
		Height(leftContentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		BorderBackground(lipgloss.Color("234")).
		Background(lipgloss.Color("237")).
		Padding(0, 1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(fmt.Sprintf("Left Panel\n\nFlex: %d\nAllocated: %d\nContent: %d", leftFlex, leftWidth, leftContentWidth))

	// Create right panel - Width/Height are CONTENT dimensions
	rightContentWidth := rightWidth - horizontalFrame
	rightContentHeight := contentHeight - verticalFrame

	rightPanel := lipgloss.NewStyle().
		Width(rightContentWidth).
		Height(rightContentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		BorderBackground(lipgloss.Color("234")).
		Background(lipgloss.Color("237")).
		Padding(0, 1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(fmt.Sprintf("Right Panel\n\nFlex: %d\nAllocated: %d\nContent: %d", rightFlex, rightWidth, rightContentWidth))

	// Join panels horizontally
	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Wrap with background to fill any remaining space
	// Note: JoinHorizontal doesn't fill background, so we need to do it manually
	return lipgloss.NewStyle().
		Background(lipgloss.Color("234")).
		Width(m.width).
		Height(contentHeight).
		Render(joined)
}

// renderThreePanelLayout demonstrates three-panel flex layout (1:2:1 ratio)
func (m playgroundModel) renderThreePanelLayout(contentHeight int) string {
	// Flex values
	leftFlex := 1
	centerFlex := 2
	rightFlex := 1
	totalFlex := leftFlex + centerFlex + rightFlex

	// Calculate widths based on flex ratio
	leftWidth := m.width * leftFlex / totalFlex
	centerWidth := m.width * centerFlex / totalFlex
	rightWidth := m.width - leftWidth - centerWidth // Remaining space to avoid rounding issues

	// Create a sample style to get frame size
	sampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)
	horizontalFrame := sampleStyle.GetHorizontalFrameSize()
	verticalFrame := sampleStyle.GetVerticalFrameSize()

	// Create left panel - Width/Height are CONTENT dimensions
	leftContentWidth := leftWidth - horizontalFrame
	leftContentHeight := contentHeight - verticalFrame

	leftPanel := lipgloss.NewStyle().
		Width(leftContentWidth).
		Height(leftContentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		BorderBackground(lipgloss.Color("234")).
		Background(lipgloss.Color("237")).
		Padding(0, 1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(fmt.Sprintf("Left\n\nFlex: %d\nWidth: %d", leftFlex, leftWidth))

	// Create center panel
	centerContentWidth := centerWidth - horizontalFrame
	centerContentHeight := contentHeight - verticalFrame

	centerPanel := lipgloss.NewStyle().
		Width(centerContentWidth).
		Height(centerContentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		BorderBackground(lipgloss.Color("234")).
		Background(lipgloss.Color("237")).
		Padding(0, 1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(fmt.Sprintf("Center\n\nFlex: %d\nWidth: %d", centerFlex, centerWidth))

	// Create right panel
	rightContentWidth := rightWidth - horizontalFrame
	rightContentHeight := contentHeight - verticalFrame

	rightPanel := lipgloss.NewStyle().
		Width(rightContentWidth).
		Height(rightContentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		BorderBackground(lipgloss.Color("234")).
		Background(lipgloss.Color("237")).
		Padding(0, 1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(fmt.Sprintf("Right\n\nFlex: %d\nWidth: %d", rightFlex, rightWidth))

	// Join panels horizontally
	joined := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, centerPanel, rightPanel)

	// Wrap with background to fill any remaining space
	return lipgloss.NewStyle().
		Background(lipgloss.Color("234")).
		Width(m.width).
		Height(contentHeight).
		Render(joined)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running dashbrew playground: %v\n", err)
		os.Exit(1)
	}
}
