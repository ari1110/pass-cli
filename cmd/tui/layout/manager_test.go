package layout

import (
	"testing"
)

// TestDetermineLayoutMode verifies breakpoint logic at critical boundaries.
func TestDetermineLayoutMode(t *testing.T) {
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
	}

	tests := []struct {
		name     string
		width    int
		expected LayoutMode
	}{
		// Small mode (< 80)
		{"Very narrow terminal", 40, LayoutSmall},
		{"Just below medium breakpoint", 79, LayoutSmall},

		// Medium mode (80-119)
		{"Exactly at medium breakpoint", 80, LayoutMedium},
		{"Middle of medium range", 100, LayoutMedium},
		{"Just below large breakpoint", 119, LayoutMedium},

		// Large mode (>= 120)
		{"Exactly at large breakpoint", 120, LayoutLarge},
		{"Wide terminal", 150, LayoutLarge},
		{"Very wide terminal", 200, LayoutLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lm.determineLayoutMode(tt.width)
			if result != tt.expected {
				t.Errorf("determineLayoutMode(%d) = %v, want %v", tt.width, result, tt.expected)
			}
		})
	}
}

// TestSetBreakpoints verifies custom breakpoint configuration.
func TestSetBreakpoints(t *testing.T) {
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
	}

	// Set custom breakpoints
	lm.SetBreakpoints(90, 140)

	if lm.mediumBreakpoint != 90 {
		t.Errorf("mediumBreakpoint = %d, want 90", lm.mediumBreakpoint)
	}
	if lm.largeBreakpoint != 140 {
		t.Errorf("largeBreakpoint = %d, want 140", lm.largeBreakpoint)
	}

	// Verify layout mode calculation uses new breakpoints
	if mode := lm.determineLayoutMode(85); mode != LayoutSmall {
		t.Errorf("With breakpoint 90, width 85 should be Small, got %v", mode)
	}
	if mode := lm.determineLayoutMode(95); mode != LayoutMedium {
		t.Errorf("With breakpoint 90-140, width 95 should be Medium, got %v", mode)
	}
	if mode := lm.determineLayoutMode(145); mode != LayoutLarge {
		t.Errorf("With breakpoint 140, width 145 should be Large, got %v", mode)
	}
}

// TestGetCurrentMode verifies mode tracking.
func TestGetCurrentMode(t *testing.T) {
	lm := &LayoutManager{
		currentMode: LayoutMedium,
	}

	if mode := lm.GetCurrentMode(); mode != LayoutMedium {
		t.Errorf("GetCurrentMode() = %v, want LayoutMedium", mode)
	}

	lm.currentMode = LayoutLarge
	if mode := lm.GetCurrentMode(); mode != LayoutLarge {
		t.Errorf("GetCurrentMode() = %v, want LayoutLarge", mode)
	}
}

// TestHandleResize verifies resize detection and mode changes.
func TestHandleResize(t *testing.T) {
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall,
	}

	// Note: contentRow is nil, but rebuildLayout() guards against this
	// In real usage, CreateMainLayout() initializes all components

	// Resize to medium mode
	lm.HandleResize(100, 40)

	if lm.width != 100 {
		t.Errorf("width = %d, want 100", lm.width)
	}
	if lm.height != 40 {
		t.Errorf("height = %d, want 40", lm.height)
	}
	if lm.currentMode != LayoutMedium {
		t.Errorf("currentMode = %v, want LayoutMedium", lm.currentMode)
	}

	// Resize to large mode
	lm.HandleResize(150, 50)

	if lm.width != 150 {
		t.Errorf("width = %d, want 150", lm.width)
	}
	if lm.height != 50 {
		t.Errorf("height = %d, want 50", lm.height)
	}
	if lm.currentMode != LayoutLarge {
		t.Errorf("currentMode = %v, want LayoutLarge", lm.currentMode)
	}

	// Resize within same mode (should not rebuild)
	previousMode := lm.currentMode
	lm.HandleResize(155, 50)

	if lm.currentMode != previousMode {
		t.Errorf("Mode should not change for resize within same range")
	}
}

// TestLayoutModeConstants verifies enum values are distinct.
func TestLayoutModeConstants(t *testing.T) {
	if LayoutSmall == LayoutMedium {
		t.Error("LayoutSmall and LayoutMedium should be distinct")
	}
	if LayoutMedium == LayoutLarge {
		t.Error("LayoutMedium and LayoutLarge should be distinct")
	}
	if LayoutSmall == LayoutLarge {
		t.Error("LayoutSmall and LayoutLarge should be distinct")
	}
}

// =============================================================================
// User Story 1 Tests: Terminal Size Warning Display
// =============================================================================

// mockPageManager is a test double for PageManager to verify ShowSizeWarning calls.
type mockPageManager struct {
	showSizeWarningCalled bool
	showSizeWarningArgs   struct {
		currentWidth  int
		currentHeight int
		minWidth      int
		minHeight     int
	}
	hideSizeWarningCalled bool
	sizeWarningActive     bool
}

func (m *mockPageManager) ShowSizeWarning(currentWidth, currentHeight, minWidth, minHeight int) {
	m.showSizeWarningCalled = true
	m.showSizeWarningArgs.currentWidth = currentWidth
	m.showSizeWarningArgs.currentHeight = currentHeight
	m.showSizeWarningArgs.minWidth = minWidth
	m.showSizeWarningArgs.minHeight = minHeight
	m.sizeWarningActive = true
}

func (m *mockPageManager) HideSizeWarning() {
	m.hideSizeWarningCalled = true
	m.sizeWarningActive = false
}

func (m *mockPageManager) IsSizeWarningActive() bool {
	return m.sizeWarningActive
}

// TestHandleResize_BelowMinimum verifies HandleResize calls ShowSizeWarning
// when width < 60 OR height < 30.
func TestHandleResize_BelowMinimum(t *testing.T) {
	tests := []struct {
		name          string
		width         int
		height        int
		shouldTrigger bool
	}{
		{"Both dimensions below minimum", 50, 20, true},
		{"Width below minimum", 50, 40, true},
		{"Height below minimum", 80, 20, true},
		{"Both dimensions at minimum", 60, 30, false},
		{"Both dimensions above minimum", 80, 40, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPM := &mockPageManager{}
			lm := &LayoutManager{
				mediumBreakpoint: 80,
				largeBreakpoint:  120,
				currentMode:      LayoutSmall,
				pageManager:      mockPM,
			}

			lm.HandleResize(tt.width, tt.height)

			if tt.shouldTrigger && !mockPM.showSizeWarningCalled {
				t.Errorf("Expected ShowSizeWarning to be called for %dx%d", tt.width, tt.height)
			}
			if !tt.shouldTrigger && mockPM.showSizeWarningCalled {
				t.Errorf("Expected ShowSizeWarning NOT to be called for %dx%d", tt.width, tt.height)
			}

			// Verify correct arguments were passed
			if tt.shouldTrigger {
				if mockPM.showSizeWarningArgs.currentWidth != tt.width {
					t.Errorf("currentWidth = %d, want %d", mockPM.showSizeWarningArgs.currentWidth, tt.width)
				}
				if mockPM.showSizeWarningArgs.currentHeight != tt.height {
					t.Errorf("currentHeight = %d, want %d", mockPM.showSizeWarningArgs.currentHeight, tt.height)
				}
				if mockPM.showSizeWarningArgs.minWidth != MinTerminalWidth {
					t.Errorf("minWidth = %d, want %d", mockPM.showSizeWarningArgs.minWidth, MinTerminalWidth)
				}
				if mockPM.showSizeWarningArgs.minHeight != MinTerminalHeight {
					t.Errorf("minHeight = %d, want %d", mockPM.showSizeWarningArgs.minHeight, MinTerminalHeight)
				}
			}
		})
	}
}

// TestHandleResize_StartupCheck verifies startup size check triggers warning
// if terminal is already too small.
func TestHandleResize_StartupCheck(t *testing.T) {
	mockPM := &mockPageManager{}
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall,
		pageManager:      mockPM,
	}

	// Simulate startup with small terminal
	lm.HandleResize(50, 20)

	// Verify warning was shown
	if !mockPM.showSizeWarningCalled {
		t.Error("Expected ShowSizeWarning to be called on startup with small terminal")
	}

	// Verify correct dimensions passed
	if mockPM.showSizeWarningArgs.currentWidth != 50 {
		t.Errorf("currentWidth = %d, want 50", mockPM.showSizeWarningArgs.currentWidth)
	}
	if mockPM.showSizeWarningArgs.currentHeight != 20 {
		t.Errorf("currentHeight = %d, want 20", mockPM.showSizeWarningArgs.currentHeight)
	}
}

// =============================================================================
// User Story 2 Tests: Automatic Recovery
// =============================================================================

// TestHideSizeWarning verifies HideSizeWarning removes the warning page
// and clears the state flag.
func TestHideSizeWarning(t *testing.T) {
	mockPM := &mockPageManager{sizeWarningActive: true}
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall,
		pageManager:      mockPM,
	}

	// Resize to adequate size
	lm.HandleResize(80, 40)

	// Verify HideSizeWarning was called
	if !mockPM.hideSizeWarningCalled {
		t.Error("Expected HideSizeWarning to be called when resizing to adequate size")
	}

	// Verify state was cleared
	if mockPM.sizeWarningActive {
		t.Error("Expected sizeWarningActive to be false after HideSizeWarning")
	}
}

// TestHideSizeWarning_WhenNotActive verifies safe no-op when warning not showing.
func TestHideSizeWarning_WhenNotActive(t *testing.T) {
	mockPM := &mockPageManager{sizeWarningActive: false}
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall,
		pageManager:      mockPM,
	}

	// Resize to adequate size when warning already hidden
	lm.HandleResize(80, 40)

	// Should still call HideSizeWarning (method handles idempotency)
	if !mockPM.hideSizeWarningCalled {
		t.Error("Expected HideSizeWarning to be called even when not active")
	}

	// State should remain false
	if mockPM.sizeWarningActive {
		t.Error("Expected sizeWarningActive to remain false")
	}
}

// TestHandleResize_ExactlyAtMinimum verifies 60×30 does NOT trigger warning
// (inclusive boundary).
func TestHandleResize_ExactlyAtMinimum(t *testing.T) {
	mockPM := &mockPageManager{}
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall,
		pageManager:      mockPM,
	}

	// Resize to exactly minimum dimensions
	lm.HandleResize(60, 30)

	// Should NOT trigger warning
	if mockPM.showSizeWarningCalled {
		t.Error("Expected NO warning at exactly 60×30 (inclusive boundary)")
	}

	// Should call HideSizeWarning (to clear any existing warning)
	if !mockPM.hideSizeWarningCalled {
		t.Error("Expected HideSizeWarning to be called at adequate size")
	}
}

// TestHandleResize_PartialFailure verifies 70×25 triggers warning
// (height < 30, OR logic).
func TestHandleResize_PartialFailure(t *testing.T) {
	mockPM := &mockPageManager{}
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall,
		pageManager:      mockPM,
	}

	// Width OK (70 >= 60), but height too small (25 < 30)
	lm.HandleResize(70, 25)

	// Should trigger warning (OR logic: width OK but height fails)
	if !mockPM.showSizeWarningCalled {
		t.Error("Expected warning to be shown for 70×25 (height < 30)")
	}

	// Verify correct dimensions passed
	if mockPM.showSizeWarningArgs.currentWidth != 70 {
		t.Errorf("currentWidth = %d, want 70", mockPM.showSizeWarningArgs.currentWidth)
	}
	if mockPM.showSizeWarningArgs.currentHeight != 25 {
		t.Errorf("currentHeight = %d, want 25", mockPM.showSizeWarningArgs.currentHeight)
	}
}

// TestFullResizeFlow_ShowAndHide verifies end-to-end resize flow:
// start 50×20 (warning shows), resize 80×40 (warning hides), interface functional.
func TestFullResizeFlow_ShowAndHide(t *testing.T) {
	mockPM := &mockPageManager{}
	lm := &LayoutManager{
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall,
		pageManager:      mockPM,
	}

	// Step 1: Startup with small terminal
	lm.HandleResize(50, 20)

	// Verify warning shown
	if !mockPM.showSizeWarningCalled {
		t.Error("Expected warning to be shown at 50×20")
	}
	if !mockPM.sizeWarningActive {
		t.Error("Expected sizeWarningActive=true after showing warning")
	}

	// Step 2: Resize to adequate size
	mockPM.hideSizeWarningCalled = false // Reset flag
	lm.HandleResize(80, 40)

	// Verify warning hidden
	if !mockPM.hideSizeWarningCalled {
		t.Error("Expected HideSizeWarning to be called at 80×40")
	}
	if mockPM.sizeWarningActive {
		t.Error("Expected sizeWarningActive=false after hiding warning")
	}

	// Step 3: Verify layout mode updated correctly (recovery functional)
	if lm.currentMode != LayoutMedium {
		t.Errorf("Expected LayoutMedium at width 80, got %v", lm.currentMode)
	}
	if lm.width != 80 {
		t.Errorf("Expected width=80, got %d", lm.width)
	}
	if lm.height != 40 {
		t.Errorf("Expected height=40, got %d", lm.height)
	}
}
