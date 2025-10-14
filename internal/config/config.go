package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the root configuration object containing all user settings
type Config struct {
	Terminal    TerminalConfig    `mapstructure:"terminal"`
	Keybindings map[string]string `mapstructure:"keybindings"`

	// LoadErrors populated during config loading (not in YAML)
	LoadErrors []string `mapstructure:"-"`
}

// TerminalConfig represents terminal size warning configuration
type TerminalConfig struct {
	WarningEnabled bool `mapstructure:"warning_enabled"`
	MinWidth       int  `mapstructure:"min_width"`
	MinHeight      int  `mapstructure:"min_height"`
}

// ValidationResult represents the outcome of checking configuration correctness
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationWarning
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string // e.g., "keybindings.add_credential"
	Message string // e.g., "conflicts with keybindings.delete_credential (both use 'd')"
	Line    int    // Line number in YAML (if available)
}

// ValidationWarning represents a non-fatal validation warning
type ValidationWarning struct {
	Field   string
	Message string
}

// GetDefaults returns the default configuration with hardcoded terminal and keybinding values
func GetDefaults() *Config {
	return &Config{
		Terminal: TerminalConfig{
			WarningEnabled: true,
			MinWidth:       60,
			MinHeight:      30,
		},
		Keybindings: map[string]string{
			"quit":              "q",
			"add_credential":    "a",
			"edit_credential":   "e",
			"delete_credential": "d",
			"toggle_detail":     "tab",
			"toggle_sidebar":    "s",
			"help":              "?",
			"search":            "/",
			"confirm":           "enter",
			"cancel":            "esc",
		},
		LoadErrors: []string{},
	}
}

// GetConfigPath returns the OS-appropriate config file path using os.UserConfigDir()
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory if UserConfigDir fails
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine config directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".pass-cli")
	} else {
		configDir = filepath.Join(configDir, "pass-cli")
	}

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.yml"), nil
}

// Load loads configuration from file and returns defaults for now (skeleton implementation)
func Load() (*Config, *ValidationResult) {
	configPath, err := GetConfigPath()
	if err != nil {
		// Cannot determine config path, use defaults
		return GetDefaults(), &ValidationResult{
			Valid: true,
			Warnings: []ValidationWarning{
				{Field: "config_path", Message: fmt.Sprintf("cannot determine config path: %v", err)},
			},
		}
	}

	// Check if config file exists
	fileInfo, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		// No config file, use defaults (not an error)
		return GetDefaults(), &ValidationResult{Valid: true}
	}
	if err != nil {
		// File stat error, use defaults
		return GetDefaults(), &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Field: "config_file", Message: fmt.Sprintf("cannot access config file: %v", err)},
			},
		}
	}

	// Check file size limit (100 KB)
	const maxFileSize = 100 * 1024 // 100 KB
	if fileInfo.Size() > maxFileSize {
		return GetDefaults(), &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "config_file",
					Message: fmt.Sprintf("config file too large (size: %d KB, max: 100 KB)", fileInfo.Size()/1024),
				},
			},
		}
	}

	// TODO: Implement YAML parsing and validation in later tasks
	return GetDefaults(), &ValidationResult{Valid: true}
}

// Validate validates the configuration and returns a validation result (skeleton implementation)
func (c *Config) Validate() *ValidationResult {
	// TODO: Implement actual validation logic in later tasks
	return &ValidationResult{Valid: true}
}
