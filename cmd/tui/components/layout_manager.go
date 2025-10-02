package components

import (
	"pass-cli/cmd/tui/styles"
)

// LayoutManager handles responsive dimension calculations for the dashboard layout
type LayoutManager struct {
	minWidth  int
	minHeight int
}

// PanelDimensions represents the dimensions of a single panel
type PanelDimensions struct {
	X             int
	Y             int
	Width         int // Total allocated width (includes border overhead)
	Height        int // Total allocated height (includes border overhead)
	ContentWidth  int // Width available for content (Width minus horizontal frame size from border + padding)
	ContentHeight int // Height available for content (Height minus vertical frame size from border + padding)
}

// Layout represents the complete dashboard layout with all panel dimensions
type Layout struct {
	Sidebar    PanelDimensions
	Main       PanelDimensions
	Metadata   PanelDimensions
	Process    PanelDimensions
	CommandBar PanelDimensions
	StatusBar  PanelDimensions
	IsTooSmall bool
	MinWidth   int
	MinHeight  int
}

// PanelStates represents the visibility state of all panels
type PanelStates struct {
	SidebarVisible   bool
	MetadataVisible  bool
	ProcessVisible   bool
	CommandBarOpen   bool
	StatusBarVisible bool // Always true, but included for completeness
}

// LayoutBreakpoint represents the current responsive breakpoint
type LayoutBreakpoint int

const (
	BreakpointSmall  LayoutBreakpoint = iota // < 80 cols
	BreakpointMedium                         // 80-119 cols
	BreakpointFull                           // >= 120 cols
)

// Panel size constraints
const (
	MinSidebarWidth  = 20
	MinMainWidth     = 40
	MinMetadataWidth = 25
	MinProcessHeight = 3
	MinCommandHeight = 1
	StatusBarHeight  = 1
	BreadcrumbHeight = 1

	// Target percentages for full layout
	TargetSidebarPercent  = 0.22 // ~22% for sidebar
	TargetMetadataPercent = 0.23 // ~23% for metadata
	// Main gets the rest
)

// NewLayoutManager creates a new layout manager with default minimum sizes
func NewLayoutManager() *LayoutManager {
	return &LayoutManager{
		minWidth:  60,
		minHeight: 20,
	}
}

// Calculate computes panel dimensions based on terminal size and panel visibility
func (lm *LayoutManager) Calculate(width, height int, states PanelStates) Layout {
	layout := Layout{
		MinWidth:  lm.minWidth,
		MinHeight: lm.minHeight,
	}

	// Check if terminal is too small
	if width < lm.minWidth || height < lm.minHeight {
		layout.IsTooSmall = true
		return layout
	}

	// Determine breakpoint
	breakpoint := lm.getBreakpoint(width)

	// Calculate vertical space allocation
	availableHeight := height
	currentY := 0

	// Reserve space for status bar (always visible)
	statusBarY := height - StatusBarHeight
	layout.StatusBar = PanelDimensions{
		X:      0,
		Y:      statusBarY,
		Width:  width,
		Height: StatusBarHeight,
	}
	setContentDimensions(&layout.StatusBar, false) // StatusBar is non-bordered
	availableHeight -= StatusBarHeight

	// Reserve space for command bar if open
	if states.CommandBarOpen {
		commandBarHeight := MinCommandHeight
		commandBarY := statusBarY - commandBarHeight
		layout.CommandBar = PanelDimensions{
			X:      0,
			Y:      commandBarY,
			Width:  width,
			Height: commandBarHeight,
		}
		setContentDimensions(&layout.CommandBar, false) // CommandBar is non-bordered
		availableHeight -= commandBarHeight
	}

	// Reserve space for process panel if visible
	if states.ProcessVisible {
		processHeight := MinProcessHeight
		processY := (statusBarY - MinCommandHeight)
		if states.CommandBarOpen {
			processY = layout.CommandBar.Y - processHeight
		}
		layout.Process = PanelDimensions{
			X:      0,
			Y:      processY,
			Width:  width,
			Height: processHeight,
		}
		setContentDimensions(&layout.Process, false) // Process panel is non-bordered
		availableHeight -= processHeight
	}

	// Calculate main content area height (includes breadcrumb)
	mainContentHeight := availableHeight

	// Calculate horizontal layout based on breakpoint and visibility
	currentX := 0

	switch breakpoint {
	case BreakpointFull:
		// Full layout: sidebar | main | metadata (if all visible)
		if states.SidebarVisible && states.MetadataVisible {
			// All three panels visible
			sidebarWidth := max(MinSidebarWidth, int(float64(width)*TargetSidebarPercent))
			metadataWidth := max(MinMetadataWidth, int(float64(width)*TargetMetadataPercent))
			mainWidth := width - sidebarWidth - metadataWidth

			// Ensure main meets minimum
			if mainWidth < MinMainWidth {
				// Reduce sidebar and metadata proportionally
				excess := MinMainWidth - mainWidth
				sidebarReduction := excess / 2
				metadataReduction := excess - sidebarReduction

				sidebarWidth -= sidebarReduction
				metadataWidth -= metadataReduction
				mainWidth = MinMainWidth

				// Ensure minimums are still met
				if sidebarWidth < MinSidebarWidth {
					sidebarWidth = MinSidebarWidth
				}
				if metadataWidth < MinMetadataWidth {
					metadataWidth = MinMetadataWidth
				}
				mainWidth = width - sidebarWidth - metadataWidth
			}

			layout.Sidebar = PanelDimensions{X: currentX, Y: currentY, Width: sidebarWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Sidebar, true) // Sidebar is bordered
			currentX += sidebarWidth

			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered
			currentX += mainWidth

			layout.Metadata = PanelDimensions{X: currentX, Y: currentY, Width: metadataWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Metadata, true) // Metadata is bordered

		} else if states.SidebarVisible {
			// Sidebar + main only
			sidebarWidth := max(MinSidebarWidth, int(float64(width)*TargetSidebarPercent))
			mainWidth := width - sidebarWidth

			if mainWidth < MinMainWidth {
				sidebarWidth = width - MinMainWidth
				mainWidth = MinMainWidth
			}

			layout.Sidebar = PanelDimensions{X: currentX, Y: currentY, Width: sidebarWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Sidebar, true) // Sidebar is bordered
			currentX += sidebarWidth

			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered

		} else if states.MetadataVisible {
			// Main + metadata only
			metadataWidth := max(MinMetadataWidth, int(float64(width)*TargetMetadataPercent))
			mainWidth := width - metadataWidth

			if mainWidth < MinMainWidth {
				metadataWidth = width - MinMainWidth
				mainWidth = MinMainWidth
			}

			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered
			currentX += mainWidth

			layout.Metadata = PanelDimensions{X: currentX, Y: currentY, Width: metadataWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Metadata, true) // Metadata is bordered

		} else {
			// Main only
			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: width, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered
		}

	case BreakpointMedium:
		// Medium layout: prioritize sidebar + main, metadata only if specifically visible
		if states.SidebarVisible && states.MetadataVisible {
			// Try to fit all three, but with tighter constraints
			sidebarWidth := MinSidebarWidth
			metadataWidth := MinMetadataWidth
			mainWidth := width - sidebarWidth - metadataWidth

			if mainWidth < MinMainWidth {
				// Not enough space, hide metadata automatically
				states.MetadataVisible = false
				mainWidth = width - sidebarWidth
			} else {
				layout.Sidebar = PanelDimensions{X: currentX, Y: currentY, Width: sidebarWidth, Height: mainContentHeight}
				setContentDimensions(&layout.Sidebar, true) // Sidebar is bordered
				currentX += sidebarWidth

				layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
				setContentDimensions(&layout.Main, true) // Main panel is bordered
				currentX += mainWidth

				layout.Metadata = PanelDimensions{X: currentX, Y: currentY, Width: metadataWidth, Height: mainContentHeight}
				setContentDimensions(&layout.Metadata, true) // Metadata is bordered
			}
		}

		if states.SidebarVisible && !states.MetadataVisible {
			sidebarWidth := MinSidebarWidth
			mainWidth := width - sidebarWidth

			layout.Sidebar = PanelDimensions{X: currentX, Y: currentY, Width: sidebarWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Sidebar, true) // Sidebar is bordered
			currentX += sidebarWidth

			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered

		} else if !states.SidebarVisible && states.MetadataVisible {
			metadataWidth := MinMetadataWidth
			mainWidth := width - metadataWidth

			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered
			currentX += mainWidth

			layout.Metadata = PanelDimensions{X: currentX, Y: currentY, Width: metadataWidth, Height: mainContentHeight}
			setContentDimensions(&layout.Metadata, true) // Metadata is bordered

		} else {
			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: width, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered
		}

	case BreakpointSmall:
		// Small layout: main only by default, toggle other panels individually
		if states.SidebarVisible && !states.MetadataVisible {
			sidebarWidth := MinSidebarWidth
			mainWidth := width - sidebarWidth

			if mainWidth >= MinMainWidth {
				layout.Sidebar = PanelDimensions{X: currentX, Y: currentY, Width: sidebarWidth, Height: mainContentHeight}
				setContentDimensions(&layout.Sidebar, true) // Sidebar is bordered
				currentX += sidebarWidth
				layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
				setContentDimensions(&layout.Main, true) // Main panel is bordered
			} else {
				// Not enough space, main only
				layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: width, Height: mainContentHeight}
				setContentDimensions(&layout.Main, true) // Main panel is bordered
			}

		} else if !states.SidebarVisible && states.MetadataVisible {
			metadataWidth := MinMetadataWidth
			mainWidth := width - metadataWidth

			if mainWidth >= MinMainWidth {
				layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: mainWidth, Height: mainContentHeight}
				setContentDimensions(&layout.Main, true) // Main panel is bordered
				currentX += mainWidth
				layout.Metadata = PanelDimensions{X: currentX, Y: currentY, Width: metadataWidth, Height: mainContentHeight}
				setContentDimensions(&layout.Metadata, true) // Metadata is bordered
			} else {
				// Not enough space, main only
				layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: width, Height: mainContentHeight}
				setContentDimensions(&layout.Main, true) // Main panel is bordered
			}

		} else {
			// Main only (or both panels hidden due to space constraints)
			layout.Main = PanelDimensions{X: currentX, Y: currentY, Width: width, Height: mainContentHeight}
			setContentDimensions(&layout.Main, true) // Main panel is bordered
		}
	}

	return layout
}

// getBreakpoint determines the current responsive breakpoint
func (lm *LayoutManager) getBreakpoint(width int) LayoutBreakpoint {
	if width >= 120 {
		return BreakpointFull
	} else if width >= 80 {
		return BreakpointMedium
	}
	return BreakpointSmall
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// setContentDimensions calculates and sets ContentWidth/ContentHeight for a panel.
// For bordered panels, subtracts frame size using Lipgloss GetFrameSize() to account for border/padding overhead.
// For non-bordered panels, content dimensions equal total dimensions.
//
// IMPORTANT: This uses GetFrameSize() instead of hardcoded constants because:
// - GetFrameSize() automatically accounts for border style (e.g., RoundedBorder) and padding configuration
// - If border/padding changes in theme.go, calculations adapt automatically
// - Both ActivePanelBorderStyle and InactivePanelBorderStyle have identical frame sizes (same border + padding config)
func setContentDimensions(panel *PanelDimensions, bordered bool) {
	if bordered {
		// Calculate content dimensions using Lipgloss GetFrameSize() to account for border/padding overhead
		horizontalFrame := styles.ActivePanelBorderStyle.GetHorizontalFrameSize()
		verticalFrame := styles.ActivePanelBorderStyle.GetVerticalFrameSize()
		panel.ContentWidth = panel.Width - horizontalFrame
		panel.ContentHeight = panel.Height - verticalFrame
	} else {
		// Non-bordered panels: content dimensions equal total dimensions
		panel.ContentWidth = panel.Width
		panel.ContentHeight = panel.Height
	}
}
