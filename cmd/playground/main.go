package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pass-cli/cmd/tui/components"
	"pass-cli/cmd/tui/styles"
	"pass-cli/internal/vault"
)

// Layer represents a rendering layer in the TUI
type Layer int

const (
	LayerNone Layer = iota
	LayerLayout
	LayerStatusBar
	LayerSidebar
	LayerMetadata
	LayerFull
)

func (l Layer) String() string {
	switch l {
	case LayerNone:
		return "None (Terminal Info)"
	case LayerLayout:
		return "Layout Foundation"
	case LayerStatusBar:
		return "Layout + Status Bar"
	case LayerSidebar:
		return "Layout + Status Bar + Sidebar"
	case LayerMetadata:
		return "Layout + Status Bar + Sidebar + Metadata"
	case LayerFull:
		return "Full Layout (All Components)"
	default:
		return "Unknown"
	}
}

type playgroundModel struct {
	currentLayer Layer
	width        int
	height       int
	layoutMgr    *components.LayoutManager
	layout       components.Layout
	statusBar    *components.StatusBar
	sidebar      *components.SidebarPanel
	metadata     *components.MetadataPanel

	// Test data
	testCredentials []vault.CredentialMetadata
}

func initialModel() playgroundModel {
	m := playgroundModel{
		currentLayer:    LayerNone,
		testCredentials: generateTestCredentials(),
	}
	return m
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
			// Progress to next layer
			if m.currentLayer < LayerFull {
				m.currentLayer++
				m.initializeLayer()
			}
		case "left", "h":
			// Go back to previous layer
			if m.currentLayer > LayerNone {
				m.currentLayer--
				m.initializeLayer()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initializeLayer()
	}

	return m, nil
}

func (m *playgroundModel) initializeLayer() {
	if m.width == 0 || m.height == 0 {
		return
	}

	// Always create layout manager
	if m.layoutMgr == nil {
		m.layoutMgr = components.NewLayoutManager()
	}

	// Define panel states based on current layer
	states := components.PanelStates{
		SidebarVisible:   false,
		MetadataVisible:  false,
		ProcessVisible:   false,
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}

	switch m.currentLayer {
	case LayerLayout:
		// Just layout calculations, no components

	case LayerStatusBar:
		m.statusBar = components.NewStatusBar(false, len(m.testCredentials), "playground")

	case LayerSidebar:
		states.SidebarVisible = true
		m.layout = m.layoutMgr.Calculate(m.width, m.height, states)
		m.statusBar = components.NewStatusBar(false, len(m.testCredentials), "list")
		m.sidebar = components.NewSidebarPanel(m.testCredentials)
		m.sidebar.SetSize(m.layout.Sidebar.ContentWidth, m.layout.Sidebar.ContentHeight)

	case LayerMetadata:
		states.SidebarVisible = true
		states.MetadataVisible = true
		m.layout = m.layoutMgr.Calculate(m.width, m.height, states)
		m.statusBar = components.NewStatusBar(false, len(m.testCredentials), "detail")
		m.sidebar = components.NewSidebarPanel(m.testCredentials)
		m.sidebar.SetSize(m.layout.Sidebar.ContentWidth, m.layout.Sidebar.ContentHeight)
		m.metadata = components.NewMetadataPanel()
		if len(m.testCredentials) > 0 {
			cred := &vault.Credential{
				Service:  m.testCredentials[0].Service,
				Username: m.testCredentials[0].Username,
				Password: "test-password-123",
				Notes:    m.testCredentials[0].Notes,
			}
			m.metadata.SetCredential(cred)
		}
		m.metadata.SetSize(m.layout.Metadata.ContentWidth, m.layout.Metadata.ContentHeight)

	case LayerFull:
		states.SidebarVisible = true
		states.MetadataVisible = true
		m.layout = m.layoutMgr.Calculate(m.width, m.height, states)
		m.statusBar = components.NewStatusBar(false, len(m.testCredentials), "detail")
		m.sidebar = components.NewSidebarPanel(m.testCredentials)
		m.sidebar.SetSize(m.layout.Sidebar.ContentWidth, m.layout.Sidebar.ContentHeight)
		m.metadata = components.NewMetadataPanel()
		if len(m.testCredentials) > 0 {
			cred := &vault.Credential{
				Service:  m.testCredentials[0].Service,
				Username: m.testCredentials[0].Username,
				Password: "test-password-123",
				Notes:    m.testCredentials[0].Notes,
			}
			m.metadata.SetCredential(cred)
		}
		m.metadata.SetSize(m.layout.Metadata.ContentWidth, m.layout.Metadata.ContentHeight)
	}
}

func (m playgroundModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var b strings.Builder

	// Header showing current layer
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.PrimaryColor).
		Width(m.width).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("🏗️  TUI Playground - Layer: %s", m.currentLayer))

	b.WriteString(header + "\n")

	switch m.currentLayer {
	case LayerNone:
		b.WriteString(m.renderTerminalInfo())

	case LayerLayout:
		b.WriteString(m.renderLayoutInfo())

	case LayerStatusBar:
		b.WriteString(m.renderStatusBarLayer())

	case LayerSidebar:
		b.WriteString(m.renderSidebarLayer())

	case LayerMetadata:
		b.WriteString(m.renderMetadataLayer())

	case LayerFull:
		b.WriteString(m.renderFullLayout())
	}

	// Footer with navigation hints
	footer := lipgloss.NewStyle().
		Foreground(styles.SubtleColor).
		Width(m.width).
		Align(lipgloss.Center).
		Render("← h/left: Previous Layer | l/right: Next Layer → | q: Quit")

	b.WriteString("\n" + footer)

	return b.String()
}

func (m playgroundModel) renderTerminalInfo() string {
	info := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.PrimaryColor).
		Padding(1, 2).
		Width(60).
		Render(fmt.Sprintf(
			"Terminal Dimensions\n\n"+
				"Width:  %d columns\n"+
				"Height: %d rows\n\n"+
				"Press → or 'l' to start building layers",
			m.width, m.height,
		))

	return lipgloss.Place(m.width, m.height-4, lipgloss.Center, lipgloss.Center, info)
}

func (m playgroundModel) renderLayoutInfo() string {
	states := components.PanelStates{
		SidebarVisible:   true,
		MetadataVisible:  true,
		ProcessVisible:   false,
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}

	layout := m.layoutMgr.Calculate(m.width, m.height, states)

	info := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.PrimaryColor).
		Padding(1, 2).
		Width(70).
		Render(fmt.Sprintf(
			"Layout Manager Calculations\n\n"+
				"Terminal:       %d×%d\n"+
				"Too Small:      %v\n\n"+
				"Sidebar:        %d×%d (Content: %d×%d)\n"+
				"Main:           %d×%d (Content: %d×%d)\n"+
				"Metadata:       %d×%d (Content: %d×%d)\n"+
				"Status Bar:     %d×%d\n\n"+
				"✓ Foundation ready",
			m.width, m.height,
			layout.IsTooSmall,
			layout.Sidebar.Width, layout.Sidebar.Height,
			layout.Sidebar.ContentWidth, layout.Sidebar.ContentHeight,
			layout.Main.Width, layout.Main.Height,
			layout.Main.ContentWidth, layout.Main.ContentHeight,
			layout.Metadata.Width, layout.Metadata.Height,
			layout.Metadata.ContentWidth, layout.Metadata.ContentHeight,
			layout.StatusBar.Width, layout.StatusBar.Height,
		))

	return lipgloss.Place(m.width, m.height-4, lipgloss.Center, lipgloss.Center, info)
}

func (m playgroundModel) renderStatusBarLayer() string {
	// Show layout info plus actual status bar
	layoutInfo := m.renderLayoutInfo()

	// Render actual status bar at bottom
	m.statusBar.SetSize(m.width)
	statusBarView := m.statusBar.Render()

	// Combine
	lines := strings.Split(layoutInfo, "\n")
	if len(lines) > 3 {
		lines = lines[:len(lines)-3] // Remove footer space
	}

	content := strings.Join(lines, "\n")
	content += "\n\n" + lipgloss.NewStyle().
		Foreground(styles.SuccessColor).
		Bold(true).
		Render("Status Bar Preview:") + "\n"
	content += statusBarView

	return content
}

func (m playgroundModel) renderSidebarLayer() string {
	m.sidebar.SetFocus(true)
	sidebarView := m.sidebar.View()

	// Debug info for right side
	debugInfo := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SuccessColor).
		Padding(1, 2).
		Render(fmt.Sprintf(
			"Sidebar Component\n\n"+
				"Allocated Width:  %d\n"+
				"Content Width:    %d\n"+
				"Allocated Height: %d\n"+
				"Content Height:   %d\n\n"+
				"Credentials: %d\n\n"+
				"✓ Rendering with focus",
			m.layout.Sidebar.Width,
			m.layout.Sidebar.ContentWidth,
			m.layout.Sidebar.Height,
			m.layout.Sidebar.ContentHeight,
			len(m.testCredentials),
		))

	// Position sidebar at calculated X,Y
	mainWidth := m.width - m.layout.Sidebar.Width
	debugPlaced := lipgloss.Place(mainWidth, m.layout.Sidebar.Height, lipgloss.Center, lipgloss.Center, debugInfo)

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, debugPlaced)
	m.statusBar.SetSize(m.width)
	return content + "\n" + m.statusBar.Render()
}

func (m playgroundModel) renderMetadataLayer() string {
	m.sidebar.SetFocus(true)
	m.metadata.SetFocus(false)
	sidebarView := m.sidebar.View()
	metadataView := m.metadata.View()

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, metadataView)
	m.statusBar.SetSize(m.width)
	return content + "\n" + m.statusBar.Render()
}

func (m playgroundModel) renderFullLayout() string {
	// For now, same as metadata - we'll add main panel later
	m.sidebar.SetFocus(false)
	m.metadata.SetFocus(true)
	sidebarView := m.sidebar.View()
	metadataView := m.metadata.View()

	// Create a simple main panel placeholder
	mainPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SubtleColor).
		Width(m.layout.Main.Width).
		Height(m.layout.Main.Height).
		Render(lipgloss.Place(
			m.layout.Main.ContentWidth,
			m.layout.Main.ContentHeight,
			lipgloss.Center,
			lipgloss.Center,
			"[Main Content Area]\n\nList/Detail View Goes Here",
		))

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarView,
		mainPanel,
		metadataView,
	)

	m.statusBar.SetSize(m.width)
	return content + "\n" + m.statusBar.Render()
}

func generateTestCredentials() []vault.CredentialMetadata {
	now := time.Now()
	return []vault.CredentialMetadata{
		{
			Service:      "github.com",
			Username:     "developer@example.com",
			Notes:        "Personal GitHub account",
			CreatedAt:    now.AddDate(0, -6, 0),
			UpdatedAt:    now.AddDate(0, -1, 0),
			UsageCount:   15,
			LastAccessed: now.AddDate(0, 0, -2),
			Locations:    []string{"home", "work"},
		},
		{
			Service:      "aws-production",
			Username:     "admin",
			Notes:        "Production AWS credentials - DO NOT SHARE",
			CreatedAt:    now.AddDate(-1, 0, 0),
			UpdatedAt:    now.AddDate(0, -2, 0),
			UsageCount:   45,
			LastAccessed: now.AddDate(0, 0, -1),
			Locations:    []string{"work"},
		},
		{
			Service:      "stripe-api",
			Username:     "api-key",
			Notes:        "Test mode API key",
			CreatedAt:    now.AddDate(0, -3, 0),
			UpdatedAt:    now.AddDate(0, -3, 0),
			UsageCount:   8,
			LastAccessed: now.AddDate(0, 0, -7),
			Locations:    []string{"work"},
		},
		{
			Service:      "postgres-db",
			Username:     "dbadmin",
			Notes:        "Main application database",
			CreatedAt:    now.AddDate(0, -8, 0),
			UpdatedAt:    now.AddDate(0, -1, -5),
			UsageCount:   120,
			LastAccessed: now,
			Locations:    []string{"work", "staging"},
		},
		{
			Service:      "openai-api",
			Username:     "api-key",
			Notes:        "OpenAI API access",
			CreatedAt:    now.AddDate(0, -2, 0),
			UpdatedAt:    now.AddDate(0, -2, 0),
			UsageCount:   3,
			LastAccessed: now.AddDate(0, 0, -14),
			Locations:    []string{"home"},
		},
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running playground: %v\n", err)
		os.Exit(1)
	}
}
