package config

import (
	"testing"
)

// Placeholder for config package unit tests
// Tests will be added in later phases following TDD approach

func TestGetDefaults(t *testing.T) {
	cfg := GetDefaults()
	if cfg == nil {
		t.Fatal("GetDefaults() returned nil")
	}

	// Verify terminal defaults
	if cfg.Terminal.WarningEnabled != true {
		t.Errorf("expected WarningEnabled=true, got %v", cfg.Terminal.WarningEnabled)
	}
	if cfg.Terminal.MinWidth != 60 {
		t.Errorf("expected MinWidth=60, got %d", cfg.Terminal.MinWidth)
	}
	if cfg.Terminal.MinHeight != 30 {
		t.Errorf("expected MinHeight=30, got %d", cfg.Terminal.MinHeight)
	}

	// Verify keybindings defaults exist
	if len(cfg.Keybindings) == 0 {
		t.Error("expected default keybindings, got empty map")
	}
	if cfg.Keybindings["quit"] != "q" {
		t.Errorf("expected quit='q', got '%s'", cfg.Keybindings["quit"])
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}
	if path == "" {
		t.Error("GetConfigPath() returned empty string")
	}
}
