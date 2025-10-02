package components

import (
	"pass-cli/cmd/tui/styles"
	"testing"
)

func TestLayoutManager_Calculate_FullLayout(t *testing.T) {
	lm := NewLayoutManager()

	states := PanelStates{
		SidebarVisible:   true,
		MetadataVisible:  true,
		ProcessVisible:   false,
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}

	layout := lm.Calculate(140, 40, states)

	if layout.IsTooSmall {
		t.Error("Layout should not be too small for 140x40")
	}

	// Check that all three main panels are allocated
	if layout.Sidebar.Width < MinSidebarWidth {
		t.Errorf("Sidebar width %d is less than minimum %d", layout.Sidebar.Width, MinSidebarWidth)
	}
	if layout.Main.Width < MinMainWidth {
		t.Errorf("Main width %d is less than minimum %d", layout.Main.Width, MinMainWidth)
	}
	if layout.Metadata.Width < MinMetadataWidth {
		t.Errorf("Metadata width %d is less than minimum %d", layout.Metadata.Width, MinMetadataWidth)
	}

	// Check status bar is at bottom
	if layout.StatusBar.Height != StatusBarHeight {
		t.Errorf("Status bar height %d != %d", layout.StatusBar.Height, StatusBarHeight)
	}
}

func TestLayoutManager_Calculate_MediumLayout(t *testing.T) {
	lm := NewLayoutManager()

	states := PanelStates{
		SidebarVisible:   true,
		MetadataVisible:  false, // Hidden on medium
		ProcessVisible:   false,
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}

	layout := lm.Calculate(100, 30, states)

	if layout.IsTooSmall {
		t.Error("Layout should not be too small for 100x30")
	}

	// Sidebar and main should be visible
	if layout.Sidebar.Width < MinSidebarWidth {
		t.Errorf("Sidebar width %d is less than minimum %d", layout.Sidebar.Width, MinSidebarWidth)
	}
	if layout.Main.Width < MinMainWidth {
		t.Errorf("Main width %d is less than minimum %d", layout.Main.Width, MinMainWidth)
	}

	// Metadata should have zero width (hidden)
	if layout.Metadata.Width != 0 {
		t.Errorf("Metadata should be hidden but width is %d", layout.Metadata.Width)
	}
}

func TestLayoutManager_Calculate_SmallLayout(t *testing.T) {
	lm := NewLayoutManager()

	states := PanelStates{
		SidebarVisible:   false, // Hidden on small
		MetadataVisible:  false, // Hidden on small
		ProcessVisible:   false,
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}

	layout := lm.Calculate(70, 25, states)

	if layout.IsTooSmall {
		t.Error("Layout should not be too small for 70x25")
	}

	// Only main should be visible
	if layout.Sidebar.Width != 0 {
		t.Errorf("Sidebar should be hidden but width is %d", layout.Sidebar.Width)
	}
	if layout.Metadata.Width != 0 {
		t.Errorf("Metadata should be hidden but width is %d", layout.Metadata.Width)
	}
	if layout.Main.Width < MinMainWidth {
		t.Errorf("Main width %d is less than minimum %d", layout.Main.Width, MinMainWidth)
	}
}

func TestLayoutManager_Calculate_TooSmall(t *testing.T) {
	lm := NewLayoutManager()

	states := PanelStates{
		SidebarVisible:   true,
		MetadataVisible:  true,
		ProcessVisible:   false,
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}

	layout := lm.Calculate(50, 15, states)

	if !layout.IsTooSmall {
		t.Error("Layout should be too small for 50x15")
	}

	if layout.MinWidth != lm.minWidth {
		t.Errorf("MinWidth %d != %d", layout.MinWidth, lm.minWidth)
	}
	if layout.MinHeight != lm.minHeight {
		t.Errorf("MinHeight %d != %d", layout.MinHeight, lm.minHeight)
	}
}

func TestLayoutManager_Calculate_WithProcessPanel(t *testing.T) {
	lm := NewLayoutManager()

	states := PanelStates{
		SidebarVisible:   true,
		MetadataVisible:  false,
		ProcessVisible:   true, // Process panel visible
		CommandBarOpen:   false,
		StatusBarVisible: true,
	}

	layout := lm.Calculate(100, 30, states)

	if layout.IsTooSmall {
		t.Error("Layout should not be too small for 100x30")
	}

	// Process panel should be allocated
	if layout.Process.Height < MinProcessHeight {
		t.Errorf("Process height %d is less than minimum %d", layout.Process.Height, MinProcessHeight)
	}
}

func TestLayoutManager_Calculate_WithCommandBar(t *testing.T) {
	lm := NewLayoutManager()

	states := PanelStates{
		SidebarVisible:   true,
		MetadataVisible:  false,
		ProcessVisible:   false,
		CommandBarOpen:   true, // Command bar open
		StatusBarVisible: true,
	}

	layout := lm.Calculate(100, 30, states)

	if layout.IsTooSmall {
		t.Error("Layout should not be too small for 100x30")
	}

	// Command bar should be allocated
	if layout.CommandBar.Height < MinCommandHeight {
		t.Errorf("Command bar height %d is less than minimum %d", layout.CommandBar.Height, MinCommandHeight)
	}
}

func TestLayoutManager_Calculate_MinimumConstraints(t *testing.T) {
	lm := NewLayoutManager()

	// Test that minimum constraints are enforced
	testCases := []struct {
		name   string
		width  int
		height int
		states PanelStates
	}{
		{
			name:   "Sidebar minimum",
			width:  80,
			height: 30,
			states: PanelStates{SidebarVisible: true, MetadataVisible: false, StatusBarVisible: true},
		},
		{
			name:   "Main minimum",
			width:  80,
			height: 30,
			states: PanelStates{SidebarVisible: false, MetadataVisible: false, StatusBarVisible: true},
		},
		{
			name:   "Metadata minimum",
			width:  100,
			height: 30,
			states: PanelStates{SidebarVisible: false, MetadataVisible: true, StatusBarVisible: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			layout := lm.Calculate(tc.width, tc.height, tc.states)

			if !layout.IsTooSmall {
				if tc.states.SidebarVisible && layout.Sidebar.Width > 0 && layout.Sidebar.Width < MinSidebarWidth {
					t.Errorf("Sidebar width %d violates minimum %d", layout.Sidebar.Width, MinSidebarWidth)
				}
				if layout.Main.Width > 0 && layout.Main.Width < MinMainWidth {
					t.Errorf("Main width %d violates minimum %d", layout.Main.Width, MinMainWidth)
				}
				if tc.states.MetadataVisible && layout.Metadata.Width > 0 && layout.Metadata.Width < MinMetadataWidth {
					t.Errorf("Metadata width %d violates minimum %d", layout.Metadata.Width, MinMetadataWidth)
				}
			}
		})
	}
}

func TestLayoutManager_GetBreakpoint(t *testing.T) {
	lm := NewLayoutManager()

	testCases := []struct {
		width    int
		expected LayoutBreakpoint
	}{
		{140, BreakpointFull},
		{120, BreakpointFull},
		{119, BreakpointMedium},
		{100, BreakpointMedium},
		{80, BreakpointMedium},
		{79, BreakpointSmall},
		{70, BreakpointSmall},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result := lm.getBreakpoint(tc.width)
			if result != tc.expected {
				t.Errorf("Width %d: expected breakpoint %d, got %d", tc.width, tc.expected, result)
			}
		})
	}
}

// TestGetFrameSizeValues verifies that GetFrameSize() returns correct values for border styles
func TestGetFrameSizeValues(t *testing.T) {
	// Get frame sizes from styles
	activeHorizontal := styles.ActivePanelBorderStyle.GetHorizontalFrameSize()
	activeVertical := styles.ActivePanelBorderStyle.GetVerticalFrameSize()
	inactiveHorizontal := styles.InactivePanelBorderStyle.GetHorizontalFrameSize()
	inactiveVertical := styles.InactivePanelBorderStyle.GetVerticalFrameSize()

	// RoundedBorder = 2 chars + Padding(0,1) = 2 chars horizontal = 4 total horizontal
	expectedHorizontal := 4
	// RoundedBorder = 2 chars + no vertical padding = 2 total vertical
	expectedVertical := 2

	if activeHorizontal != expectedHorizontal {
		t.Errorf("ActivePanelBorderStyle horizontal frame size: got %d, want %d", activeHorizontal, expectedHorizontal)
	}
	if activeVertical != expectedVertical {
		t.Errorf("ActivePanelBorderStyle vertical frame size: got %d, want %d", activeVertical, expectedVertical)
	}
	if inactiveHorizontal != expectedHorizontal {
		t.Errorf("InactivePanelBorderStyle horizontal frame size: got %d, want %d", inactiveHorizontal, expectedHorizontal)
	}
	if inactiveVertical != expectedVertical {
		t.Errorf("InactivePanelBorderStyle vertical frame size: got %d, want %d", inactiveVertical, expectedVertical)
	}

	// Both styles should have same frame size
	if activeHorizontal != inactiveHorizontal {
		t.Errorf("Active and Inactive styles have different horizontal frame sizes: %d vs %d", activeHorizontal, inactiveHorizontal)
	}
	if activeVertical != inactiveVertical {
		t.Errorf("Active and Inactive styles have different vertical frame sizes: %d vs %d", activeVertical, inactiveVertical)
	}
}

// TestLayoutContentDimensions verifies ContentWidth/ContentHeight calculations for all breakpoints and panel configs
func TestLayoutContentDimensions(t *testing.T) {
	lm := NewLayoutManager()

	// Get expected frame sizes from actual styles
	expectedHorizontalFrame := styles.ActivePanelBorderStyle.GetHorizontalFrameSize()
	expectedVerticalFrame := styles.ActivePanelBorderStyle.GetVerticalFrameSize()

	testCases := []struct {
		name   string
		width  int
		height int
		states PanelStates
	}{
		// Small breakpoint (< 80)
		{"Small - main only", 70, 25, PanelStates{SidebarVisible: false, MetadataVisible: false, StatusBarVisible: true}},
		{"Small - sidebar + main", 70, 25, PanelStates{SidebarVisible: true, MetadataVisible: false, StatusBarVisible: true}},
		{"Small - main + metadata", 70, 25, PanelStates{SidebarVisible: false, MetadataVisible: true, StatusBarVisible: true}},

		// Medium breakpoint (80-119)
		{"Medium - main only", 100, 30, PanelStates{SidebarVisible: false, MetadataVisible: false, StatusBarVisible: true}},
		{"Medium - sidebar + main", 100, 30, PanelStates{SidebarVisible: true, MetadataVisible: false, StatusBarVisible: true}},
		{"Medium - main + metadata", 100, 30, PanelStates{SidebarVisible: false, MetadataVisible: true, StatusBarVisible: true}},
		{"Medium - all panels", 100, 30, PanelStates{SidebarVisible: true, MetadataVisible: true, StatusBarVisible: true}},

		// Full breakpoint (>= 120)
		{"Full - main only", 140, 40, PanelStates{SidebarVisible: false, MetadataVisible: false, StatusBarVisible: true}},
		{"Full - sidebar + main", 140, 40, PanelStates{SidebarVisible: true, MetadataVisible: false, StatusBarVisible: true}},
		{"Full - main + metadata", 140, 40, PanelStates{SidebarVisible: false, MetadataVisible: true, StatusBarVisible: true}},
		{"Full - all panels", 140, 40, PanelStates{SidebarVisible: true, MetadataVisible: true, StatusBarVisible: true}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			layout := lm.Calculate(tc.width, tc.height, tc.states)

			// Check bordered panels (Sidebar, Main, Metadata)
			if layout.Sidebar.Width > 0 {
				expectedContent := layout.Sidebar.Width - expectedHorizontalFrame
				if layout.Sidebar.ContentWidth != expectedContent {
					t.Errorf("Sidebar ContentWidth: got %d, want %d (Width=%d - frame=%d)",
						layout.Sidebar.ContentWidth, expectedContent, layout.Sidebar.Width, expectedHorizontalFrame)
				}
				expectedContentH := layout.Sidebar.Height - expectedVerticalFrame
				if layout.Sidebar.ContentHeight != expectedContentH {
					t.Errorf("Sidebar ContentHeight: got %d, want %d (Height=%d - frame=%d)",
						layout.Sidebar.ContentHeight, expectedContentH, layout.Sidebar.Height, expectedVerticalFrame)
				}
			}

			if layout.Main.Width > 0 {
				expectedContent := layout.Main.Width - expectedHorizontalFrame
				if layout.Main.ContentWidth != expectedContent {
					t.Errorf("Main ContentWidth: got %d, want %d (Width=%d - frame=%d)",
						layout.Main.ContentWidth, expectedContent, layout.Main.Width, expectedHorizontalFrame)
				}
				expectedContentH := layout.Main.Height - expectedVerticalFrame
				if layout.Main.ContentHeight != expectedContentH {
					t.Errorf("Main ContentHeight: got %d, want %d (Height=%d - frame=%d)",
						layout.Main.ContentHeight, expectedContentH, layout.Main.Height, expectedVerticalFrame)
				}
			}

			if layout.Metadata.Width > 0 {
				expectedContent := layout.Metadata.Width - expectedHorizontalFrame
				if layout.Metadata.ContentWidth != expectedContent {
					t.Errorf("Metadata ContentWidth: got %d, want %d (Width=%d - frame=%d)",
						layout.Metadata.ContentWidth, expectedContent, layout.Metadata.Width, expectedHorizontalFrame)
				}
				expectedContentH := layout.Metadata.Height - expectedVerticalFrame
				if layout.Metadata.ContentHeight != expectedContentH {
					t.Errorf("Metadata ContentHeight: got %d, want %d (Height=%d - frame=%d)",
						layout.Metadata.ContentHeight, expectedContentH, layout.Metadata.Height, expectedVerticalFrame)
				}
			}

			// Check non-bordered panels (StatusBar, Process, CommandBar)
			// For non-bordered panels, ContentWidth == Width and ContentHeight == Height
			if layout.StatusBar.ContentWidth != layout.StatusBar.Width {
				t.Errorf("StatusBar ContentWidth: got %d, want %d (should equal Width for non-bordered)",
					layout.StatusBar.ContentWidth, layout.StatusBar.Width)
			}
			if layout.StatusBar.ContentHeight != layout.StatusBar.Height {
				t.Errorf("StatusBar ContentHeight: got %d, want %d (should equal Height for non-bordered)",
					layout.StatusBar.ContentHeight, layout.StatusBar.Height)
			}
		})
	}
}

// TestLayoutTotalDimensionsPreserved verifies that sum of panel widths equals terminal width
func TestLayoutTotalDimensionsPreserved(t *testing.T) {
	lm := NewLayoutManager()

	testCases := []struct {
		name   string
		width  int
		height int
		states PanelStates
	}{
		{"Full layout", 140, 40, PanelStates{SidebarVisible: true, MetadataVisible: true, StatusBarVisible: true}},
		{"Medium layout", 100, 30, PanelStates{SidebarVisible: true, MetadataVisible: false, StatusBarVisible: true}},
		{"Small layout", 70, 25, PanelStates{SidebarVisible: false, MetadataVisible: false, StatusBarVisible: true}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			layout := lm.Calculate(tc.width, tc.height, tc.states)

			// Calculate sum of horizontal panel widths
			totalWidth := 0
			if layout.Sidebar.Width > 0 {
				totalWidth += layout.Sidebar.Width
			}
			if layout.Main.Width > 0 {
				totalWidth += layout.Main.Width
			}
			if layout.Metadata.Width > 0 {
				totalWidth += layout.Metadata.Width
			}

			// For non-too-small layouts, total should equal terminal width
			if !layout.IsTooSmall && totalWidth != tc.width {
				t.Errorf("Total panel width %d != terminal width %d (Sidebar=%d, Main=%d, Metadata=%d)",
					totalWidth, tc.width, layout.Sidebar.Width, layout.Main.Width, layout.Metadata.Width)
			}

			// StatusBar should span full width
			if layout.StatusBar.Width != tc.width {
				t.Errorf("StatusBar width %d != terminal width %d", layout.StatusBar.Width, tc.width)
			}
		})
	}
}

// TestBreakpointContentDimensions verifies content dimensions are correct at breakpoint transitions
func TestBreakpointContentDimensions(t *testing.T) {
	lm := NewLayoutManager()

	// Get expected frame sizes
	expectedHorizontalFrame := styles.ActivePanelBorderStyle.GetHorizontalFrameSize()
	expectedVerticalFrame := styles.ActivePanelBorderStyle.GetVerticalFrameSize()

	testCases := []struct {
		name   string
		width  int
		states PanelStates
	}{
		// Test 79→80 transition (Small to Medium)
		{"Small breakpoint (79)", 79, PanelStates{SidebarVisible: true, MetadataVisible: false, StatusBarVisible: true}},
		{"Medium breakpoint (80)", 80, PanelStates{SidebarVisible: true, MetadataVisible: false, StatusBarVisible: true}},

		// Test 119→120 transition (Medium to Full)
		{"Medium breakpoint (119)", 119, PanelStates{SidebarVisible: true, MetadataVisible: true, StatusBarVisible: true}},
		{"Full breakpoint (120)", 120, PanelStates{SidebarVisible: true, MetadataVisible: true, StatusBarVisible: true}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			layout := lm.Calculate(tc.width, 30, tc.states)

			// Verify all visible bordered panels have correct content dimensions
			verifyPanel := func(name string, dims PanelDimensions) {
				if dims.Width > 0 {
					expectedContent := dims.Width - expectedHorizontalFrame
					if dims.ContentWidth != expectedContent {
						t.Errorf("%s ContentWidth at width %d: got %d, want %d",
							name, tc.width, dims.ContentWidth, expectedContent)
					}
					expectedContentH := dims.Height - expectedVerticalFrame
					if dims.ContentHeight != expectedContentH {
						t.Errorf("%s ContentHeight at width %d: got %d, want %d",
							name, tc.width, dims.ContentHeight, expectedContentH)
					}
				}
			}

			verifyPanel("Sidebar", layout.Sidebar)
			verifyPanel("Main", layout.Main)
			verifyPanel("Metadata", layout.Metadata)
		})
	}
}
