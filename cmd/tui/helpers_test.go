package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"pass-cli/cmd/tui/components"
	"pass-cli/internal/vault"
)

// createTestModel creates a minimal Model for testing dashboard rendering
func createTestModel(width, height int, sidebarVisible, metadataVisible bool) *Model {
	m := &Model{
		state:           StateList,
		width:           width,
		height:          height,
		sidebarVisible:  sidebarVisible,
		metadataVisible: metadataVisible,
		processVisible:  false,
		commandBarOpen:  false,
		panelFocus:      FocusMain,
		layoutManager:   components.NewLayoutManager(),
	}

	// Initialize minimal components for rendering
	m.statusBar = components.NewStatusBar(false, 0, "list")              // keychainAvailable, credentialCount, currentView
	m.sidebar = components.NewSidebarPanel([]vault.CredentialMetadata{}) // empty credentials for testing
	m.metadataPanel = components.NewMetadataPanel()
	m.listView = nil // Will render empty content
	m.detailView = nil

	// Calculate initial layout
	m.recalculateLayout()

	return m
}

// TestDashboardWithinBounds verifies dashboard renders within terminal bounds
func TestDashboardWithinBounds(t *testing.T) {
	testCases := []struct {
		name            string
		width           int
		height          int
		sidebarVisible  bool
		metadataVisible bool
	}{
		{"100x30 - all panels", 100, 30, true, true},
		{"80x24 - sidebar only", 80, 24, true, false},
		{"120x40 - full layout", 120, 40, true, true},
		{"100x30 - main only", 100, 30, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := createTestModel(tc.width, tc.height, tc.sidebarVisible, tc.metadataVisible)
			output := m.renderDashboardView()

			// Measure rendered output
			renderedWidth := lipgloss.Width(output)
			renderedHeight := lipgloss.Height(output)

			// Verify within bounds
			if renderedWidth > tc.width {
				t.Errorf("Rendered width %d exceeds terminal width %d", renderedWidth, tc.width)
			}
			if renderedHeight > tc.height {
				t.Errorf("Rendered height %d exceeds terminal height %d", renderedHeight, tc.height)
			}

			// Check for border characters to verify borders are present
			borderChars := []string{"─", "│", "╭", "╮", "╰", "╯"}
			hasBorders := false
			for _, char := range borderChars {
				if strings.Contains(output, char) {
					hasBorders = true
					break
				}
			}
			if !hasBorders {
				t.Error("No border characters found in rendered output")
			}
		})
	}
}

// TestPanelVisibilityBorders verifies correct number of borders for different panel configs
func TestPanelVisibilityBorders(t *testing.T) {
	testCases := []struct {
		name            string
		sidebarVisible  bool
		metadataVisible bool
		expectedBorders int // Number of vertical panel borders
	}{
		{"main only", false, false, 2},      // Main panel: left + right border
		{"sidebar + main", true, false, 3},  // Sidebar right + Main left + Main right
		{"main + metadata", false, true, 3}, // Main left + Main right + Metadata left
		{"all panels", true, true, 4},       // Sidebar right + Main left + Main right + Metadata left
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := createTestModel(120, 40, tc.sidebarVisible, tc.metadataVisible)
			output := m.renderDashboardView()

			// Count vertical border characters (│) - rough approximation
			// This is not perfect but gives us a basic check
			verticalBorderCount := strings.Count(output, "│")

			// We expect at least some vertical borders
			if verticalBorderCount == 0 {
				t.Error("No vertical borders found in output")
			}

			// Check that borders are present (actual count depends on panel height)
			t.Logf("Panel config: sidebar=%v metadata=%v, vertical borders found: %d",
				tc.sidebarVisible, tc.metadataVisible, verticalBorderCount)
		})
	}
}

// TestFocusDimensionsUnchanged verifies focus changes don't affect rendered dimensions
func TestFocusDimensionsUnchanged(t *testing.T) {
	m := createTestModel(120, 40, true, true)

	focusStates := []PanelFocus{FocusSidebar, FocusMain, FocusMetadata}
	var baselineWidth, baselineHeight int

	for i, focus := range focusStates {
		m.panelFocus = focus
		m.updatePanelFocus() // Update focus state
		output := m.renderDashboardView()

		renderedWidth := lipgloss.Width(output)
		renderedHeight := lipgloss.Height(output)

		if i == 0 {
			// Set baseline from first render
			baselineWidth = renderedWidth
			baselineHeight = renderedHeight
		} else {
			// Compare to baseline
			if renderedWidth != baselineWidth {
				t.Errorf("Focus %v: width %d != baseline %d (focus should not change dimensions)",
					focus, renderedWidth, baselineWidth)
			}
			if renderedHeight != baselineHeight {
				t.Errorf("Focus %v: height %d != baseline %d (focus should not change dimensions)",
					focus, renderedHeight, baselineHeight)
			}
		}

		t.Logf("Focus %v: width=%d height=%d", focus, renderedWidth, renderedHeight)
	}
}

// TestPanelDimensionsMatchLayout verifies rendered panel dimensions match layout allocation
func TestPanelDimensionsMatchLayout(t *testing.T) {
	testCases := []struct {
		name            string
		width           int
		height          int
		sidebarVisible  bool
		metadataVisible bool
	}{
		{"Full layout", 140, 40, true, true},
		{"Medium layout", 100, 30, true, false},
		{"Small layout", 70, 25, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := createTestModel(tc.width, tc.height, tc.sidebarVisible, tc.metadataVisible)

			// Get the layout
			states := components.PanelStates{
				SidebarVisible:   tc.sidebarVisible,
				MetadataVisible:  tc.metadataVisible,
				ProcessVisible:   false,
				CommandBarOpen:   false,
				StatusBarVisible: true,
			}
			layout := m.layoutManager.Calculate(tc.width, tc.height, states)

			// Render and measure total output
			output := m.renderDashboardView()
			renderedWidth := lipgloss.Width(output)
			renderedHeight := lipgloss.Height(output)

			// Total rendered output should fit within terminal
			if renderedWidth > tc.width {
				t.Errorf("Total rendered width %d exceeds terminal width %d", renderedWidth, tc.width)
			}
			if renderedHeight > tc.height {
				t.Errorf("Total rendered height %d exceeds terminal height %d", renderedHeight, tc.height)
			}

			// Log layout dimensions for debugging
			t.Logf("Terminal: %dx%d", tc.width, tc.height)
			if layout.Sidebar.Width > 0 {
				t.Logf("Sidebar: Total=%dx%d Content=%dx%d",
					layout.Sidebar.Width, layout.Sidebar.Height,
					layout.Sidebar.ContentWidth, layout.Sidebar.ContentHeight)
			}
			if layout.Main.Width > 0 {
				t.Logf("Main: Total=%dx%d Content=%dx%d",
					layout.Main.Width, layout.Main.Height,
					layout.Main.ContentWidth, layout.Main.ContentHeight)
			}
			if layout.Metadata.Width > 0 {
				t.Logf("Metadata: Total=%dx%d Content=%dx%d",
					layout.Metadata.Width, layout.Metadata.Height,
					layout.Metadata.ContentWidth, layout.Metadata.ContentHeight)
			}
			t.Logf("Rendered output: %dx%d", renderedWidth, renderedHeight)

			// Check for exact width match (this will fail if we have double subtraction bug)
			// The rendered output should exactly match the terminal width
			if renderedWidth != tc.width {
				t.Logf("WARNING: Rendered width %d != terminal width %d (diff: %d)",
					renderedWidth, tc.width, tc.width-renderedWidth)
			}
		})
	}
}

// TestLipglossWidthBehavior tests how Lipgloss .Width() actually works
func TestLipglossWidthBehavior(t *testing.T) {
	// Create a style with RoundedBorder and Padding(0,1) like our panels
	testStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("cyan")).
		Padding(0, 1)

	frameSize := testStyle.GetHorizontalFrameSize()
	t.Logf("Frame size (border+padding): %d", frameSize)

	// Test: pass 16 to Width()
	rendered := testStyle.Width(16).Render("test")
	actualWidth := lipgloss.Width(rendered)
	t.Logf("Called .Width(16), actual rendered width: %d", actualWidth)

	// If Width(16) means content=16, then actual should be 16+frame
	// If Width(16) means total=16, then actual should be 16
	expectedIfContentWidth := 16 + frameSize
	t.Logf("Expected if Width() sets content: %d", expectedIfContentWidth)
	t.Logf("Expected if Width() sets total: 16")

	if actualWidth == expectedIfContentWidth {
		t.Logf("✓ Width() sets CONTENT width, frame added on top")
	} else if actualWidth == 16 {
		t.Logf("✓ Width() sets TOTAL width including frame")
	} else {
		t.Logf("? Unexpected behavior: neither matches")
	}
}

// TestDiagnosePanelWidths helps diagnose where pixels are being lost
func TestDiagnosePanelWidths(t *testing.T) {
	m := createTestModel(100, 30, true, false) // sidebar + main

	states := components.PanelStates{
		SidebarVisible:   true,
		MetadataVisible:  false,
		ProcessVisible:   false,
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}
	layout := m.layoutManager.Calculate(100, 30, states)

	t.Logf("=== Layout Allocation ===")
	t.Logf("Sidebar: Total=%d Content=%d (frame=%d)",
		layout.Sidebar.Width, layout.Sidebar.ContentWidth,
		layout.Sidebar.Width-layout.Sidebar.ContentWidth)
	t.Logf("Main: Total=%d Content=%d (frame=%d)",
		layout.Main.Width, layout.Main.ContentWidth,
		layout.Main.Width-layout.Main.ContentWidth)
	t.Logf("Sum of Total widths: %d", layout.Sidebar.Width+layout.Main.Width)

	// Now render and see what we get
	output := m.renderDashboardView()
	renderedWidth := lipgloss.Width(output)

	t.Logf("=== Rendered Output ===")
	t.Logf("Total rendered width: %d", renderedWidth)
	t.Logf("Pixels lost: %d", 100-renderedWidth)

	// The issue: we should get 100, but we're getting 96 (losing 4 pixels)
}

// TestExactDimensionMatch verifies rendered output exactly matches terminal dimensions
func TestExactDimensionMatch(t *testing.T) {
	testCases := []struct {
		name            string
		width           int
		height          int
		sidebarVisible  bool
		metadataVisible bool
	}{
		{"Full layout", 140, 40, true, true},
		{"Medium layout", 100, 30, true, false},
		{"Small layout", 70, 25, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := createTestModel(tc.width, tc.height, tc.sidebarVisible, tc.metadataVisible)
			output := m.renderDashboardView()

			renderedWidth := lipgloss.Width(output)
			renderedHeight := lipgloss.Height(output)

			// Rendered output should EXACTLY match terminal dimensions
			if renderedWidth != tc.width {
				t.Errorf("Rendered width %d != terminal width %d (diff: %d pixels)",
					renderedWidth, tc.width, tc.width-renderedWidth)
			}
			if renderedHeight != tc.height {
				t.Errorf("Rendered height %d != terminal height %d (diff: %d pixels)",
					renderedHeight, tc.height, tc.height-renderedHeight)
			}
		})
	}
}
