package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Simple playground model
type playgroundModel struct {
	width       int
	height      int
	currentPage int // 0 = no box, 1 = empty box, 2 = bordered box, 3 = bordered box with text
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
			if m.currentPage < 3 {
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
		headerText = "🏗️  Page 0: No Box (Just Background)"
	case 1:
		headerText = "🏗️  Page 1: Empty Box (No Border, No Text)"
	case 2:
		headerText = "🏗️  Page 2: Bordered Box (No Text)"
	case 3:
		headerText = "🏗️  Page 3: Bordered Box With Text"
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
		// Page 0: No box, just solid background
		return lipgloss.NewStyle().
			Background(lipgloss.Color("234")).
			Width(m.width).
			Height(contentHeight).
			Render("")
	}

	if m.currentPage == 1 {
		// Page 1: Empty box with no border, no text
		emptyBox := lipgloss.NewStyle().
			Width(40).
			Height(10).
			Background(lipgloss.Color("237")). // Slightly lighter gray than background
			Render("")

		// Center the empty box with background whitespace styling
		centered := lipgloss.Place(
			m.width,
			contentHeight,
			lipgloss.Center,
			lipgloss.Center,
			emptyBox,
			lipgloss.WithWhitespaceBackground(lipgloss.Color("234")),
		)

		return centered
	}

	if m.currentPage == 2 {
		// Page 2: Bordered box with no text
		borderedBox := lipgloss.NewStyle().
			Width(40).
			Height(10).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("6")).   // Cyan border
			BorderBackground(lipgloss.Color("234")). // Dark gray background for border cells
			Background(lipgloss.Color("237")).       // Light gray background for content
			Render("")

		// Center the bordered box with background whitespace styling
		centered := lipgloss.Place(
			m.width,
			contentHeight,
			lipgloss.Center,
			lipgloss.Center,
			borderedBox,
			lipgloss.WithWhitespaceBackground(lipgloss.Color("234")),
		)

		return centered
	}

	// Page 3: Bordered box with text
	boxContent := "Hello, TUI!\n\nThis is a test box\nwith multiple lines\nof text content."

	borderedBox := lipgloss.NewStyle().
		Width(40).
		Height(10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).   // Cyan border
		BorderBackground(lipgloss.Color("234")). // Dark gray background for border cells
		Background(lipgloss.Color("237")).       // Light gray background for content
		Align(lipgloss.Center, lipgloss.Center). // Center text within the box
		Render(boxContent)

	// Center the bordered box with background whitespace styling
	centered := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		borderedBox,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("234")),
	)

	return centered
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running playground: %v\n", err)
		os.Exit(1)
	}
}
