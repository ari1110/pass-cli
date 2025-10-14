package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
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

// LoadFromPath loads configuration from a specific file path (useful for testing)
func LoadFromPath(configPath string) (*Config, *ValidationResult) {
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

	// T018+T020: Load YAML with Viper and merge with defaults
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults for merging
	defaults := GetDefaults()
	v.SetDefault("terminal.warning_enabled", defaults.Terminal.WarningEnabled)
	v.SetDefault("terminal.min_width", defaults.Terminal.MinWidth)
	v.SetDefault("terminal.min_height", defaults.Terminal.MinHeight)
	for action, key := range defaults.Keybindings {
		v.SetDefault(fmt.Sprintf("keybindings.%s", action), key)
	}

	// Read and parse YAML
	if err := v.ReadInConfig(); err != nil {
		return GetDefaults(), &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Field: "config_file", Message: fmt.Sprintf("failed to parse YAML: %v", err)},
			},
		}
	}

	// Unmarshal into Config struct (Viper will merge with defaults)
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return GetDefaults(), &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Field: "config_file", Message: fmt.Sprintf("failed to unmarshal config: %v", err)},
			},
		}
	}

	// Validate the loaded config
	validationResult := cfg.Validate()

	// If validation failed, return defaults instead
	if !validationResult.Valid {
		return GetDefaults(), validationResult
	}

	return &cfg, validationResult
}

// Load loads configuration from the default config path
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

	return LoadFromPath(configPath)
}

// Validate validates the configuration and returns a validation result
func (c *Config) Validate() *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// Validate terminal configuration
	result = c.validateTerminal(result)

	// TODO: Add keybinding validation in Phase 4

	// Set Valid flag based on error count
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

// validateTerminal validates terminal size configuration
func (c *Config) validateTerminal(result *ValidationResult) *ValidationResult {
	// T019: Validate min_width range (1-10000)
	if c.Terminal.MinWidth < 1 || c.Terminal.MinWidth > 10000 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "terminal.min_width",
			Message: fmt.Sprintf("must be between 1 and 10000 (got: %d)", c.Terminal.MinWidth),
		})
	}

	// T019: Validate min_height range (1-1000)
	if c.Terminal.MinHeight < 1 || c.Terminal.MinHeight > 1000 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "terminal.min_height",
			Message: fmt.Sprintf("must be between 1 and 1000 (got: %d)", c.Terminal.MinHeight),
		})
	}

	// T021: Warn if unusually large width (>500)
	if c.Terminal.MinWidth > 500 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "terminal.min_width",
			Message: fmt.Sprintf("unusually large value (%d) - most terminals are <300 columns", c.Terminal.MinWidth),
		})
	}

	// T021: Warn if unusually large height (>200)
	if c.Terminal.MinHeight > 200 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "terminal.min_height",
			Message: fmt.Sprintf("unusually large value (%d) - most terminals are <100 rows", c.Terminal.MinHeight),
		})
	}

	return result
}
