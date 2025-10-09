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
