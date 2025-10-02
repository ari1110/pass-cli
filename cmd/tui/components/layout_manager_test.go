package components

import (
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
