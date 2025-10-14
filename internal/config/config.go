package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config represents the root configuration object containing all user settings
type Config struct {
	Terminal    TerminalConfig    `mapstructure:"terminal"`
	Keybindings map[string]string `mapstructure:"keybindings"`

	// LoadErrors populated during config loading (not in YAML)
	LoadErrors []string `mapstructure:"-"`

	// ParsedKeybindings stores parsed keybinding objects (populated during Validate)
	ParsedKeybindings map[string]*Keybinding `mapstructure:"-"`
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
	cfg := &Config{
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
			"toggle_detail":     "i",
			"toggle_sidebar":    "s",
			"help":              "?",
			"search":            "/",
			"confirm":           "enter",
			"cancel":            "esc",
		},
		LoadErrors: []string{},
	}

	// Validate to populate ParsedKeybindings
	// This ensures defaults are always ready for use
	cfg.Validate()

	return cfg
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

// GetEditor returns the editor to use for editing config files.
// Checks EDITOR environment variable first, then falls back to OS defaults.
func GetEditor() (string, error) {
	// Check EDITOR environment variable first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	// Platform-specific defaults
	switch runtime.GOOS {
	case "windows":
		return "notepad.exe", nil
	case "darwin", "linux":
		// Check for common editors in order of preference
		for _, ed := range []string{"nano", "vim", "vi"} {
			if _, err := exec.LookPath(ed); err == nil {
				return ed, nil
			}
		}
		return "", fmt.Errorf("no editor found. Please set EDITOR environment variable (e.g., export EDITOR=nano)")
	default:
		return "", fmt.Errorf("unsupported platform for editor detection")
	}
}

// OpenEditor opens the config file in the user's editor.
func OpenEditor(filePath string) error {
	editor, err := GetEditor()
	if err != nil {
		return err
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetDefaultConfigTemplate returns the default config file content with comments.
func GetDefaultConfigTemplate() string {
	return `# Pass-CLI Configuration File
# 
# This file allows you to customize terminal size warnings and keyboard shortcuts.
# All settings are optional - missing values will use defaults.

# Terminal size warning configuration
terminal:
  # Enable or disable terminal size warnings (default: true)
  warning_enabled: true
  
  # Minimum terminal width in columns before warning appears (default: 60)
  # Valid range: 1-10000
  min_width: 60
  
  # Minimum terminal height in rows before warning appears (default: 30)
  # Valid range: 1-1000
  min_height: 30

# Keyboard shortcuts
# Format: action: "key" or "modifier+key"
# Valid modifiers: ctrl, alt, shift
# Valid keys: letters, numbers, enter, esc, tab, space, f1-f12
#
# All shortcuts are lowercase (e.g., "ctrl+q" not "CTRL+Q")
keybindings:
  # Application control
  quit: "q"                    # Quit application (with confirmation)
  help: "?"                    # Show help modal with all shortcuts
  
  # Credential management
  add_credential: "a"          # Open form to add new credential
  edit_credential: "e"         # Edit selected credential
  delete_credential: "d"       # Delete selected credential (with confirmation)
  
  # View controls
  toggle_detail: "i"           # Toggle detail panel visibility
  toggle_sidebar: "s"          # Toggle sidebar visibility
  search: "/"                  # Activate search/filter mode
  
  # Form controls (used in modals/dialogs)
  confirm: "enter"             # Confirm action in forms
  cancel: "esc"                # Cancel action in forms

# Example custom keybindings:
#
# Vim-style bindings:
#   add_credential: "i"        # Insert mode
#   search: "/"                # Vim search
#   toggle_sidebar: "ctrl+w"   # Vim window command
#
# Emacs-style bindings:
#   quit: "ctrl+x"
#   search: "ctrl+s"
#   add_credential: "ctrl+n"
#
# Custom modifier keys:
#   quit: "ctrl+q"
#   add_credential: "n"
#   help: "f1"
`
}

// LoadFromPath loads configuration from a specific file path (useful for testing)

// T054: detectUnknownFields checks for unknown fields in the config YAML
func detectUnknownFields(v *viper.Viper) []ValidationWarning {
	var warnings []ValidationWarning

	// Get all keys from the config file
	allKeys := v.AllKeys()

	// Define known fields (all valid config keys)
	knownFields := map[string]bool{
		"terminal":                    true,
		"terminal.warning_enabled":    true,
		"terminal.min_width":          true,
		"terminal.min_height":         true,
		"keybindings":                 true,
		"keybindings.quit":            true,
		"keybindings.add_credential":  true,
		"keybindings.edit_credential": true,
		"keybindings.delete_credential": true,
		"keybindings.toggle_detail":   true,
		"keybindings.toggle_sidebar":  true,
		"keybindings.help":            true,
		"keybindings.search":          true,
		"keybindings.confirm":         true,
		"keybindings.cancel":          true,
	}

	// Check for unknown fields
	for _, key := range allKeys {
		if !knownFields[key] {
			warnings = append(warnings, ValidationWarning{
				Field:   key,
				Message: fmt.Sprintf("unknown field '%s' (will be ignored)", key),
			})
		}
	}

	return warnings
}

func LoadFromPath(configPath string) (*Config, *ValidationResult) {
	// T051: Log config load attempt
	fmt.Fprintf(os.Stderr, "[Config] Loading config from: %s\n", configPath)

	// Check if config file exists
	fileInfo, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		// No config file, use defaults (not an error)
		fmt.Fprintf(os.Stderr, "[Config] No config file found, using defaults\n")
		return GetDefaults(), &ValidationResult{Valid: true}
	}
	if err != nil {
		// T051: Log file access error
		fmt.Fprintf(os.Stderr, "[Config] Failed to access config file: %v\n", err)
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
		// T051: Log file size error
		fmt.Fprintf(os.Stderr, "[Config] Config file too large: %d KB (max: 100 KB)\n", fileInfo.Size()/1024)
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
		// T051: Log parse error
		fmt.Fprintf(os.Stderr, "[Config] Failed to parse YAML: %v\n", err)
		return GetDefaults(), &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Field: "config_file", Message: fmt.Sprintf("failed to parse YAML: %v", err)},
			},
		}
	}

	// T054: Detect unknown fields
	warnings := detectUnknownFields(v)

	// Unmarshal into Config struct (Viper will merge with defaults)
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		// T051: Log unmarshal error
		fmt.Fprintf(os.Stderr, "[Config] Failed to unmarshal config: %v\n", err)
		return GetDefaults(), &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Field: "config_file", Message: fmt.Sprintf("failed to unmarshal config: %v", err)},
			},
		}
	}

	// Validate the loaded config
	validationResult := cfg.Validate()

	// Add unknown field warnings to validation result
	validationResult.Warnings = append(validationResult.Warnings, warnings...)

	// T052: Log validation errors
	if !validationResult.Valid {
		fmt.Fprintf(os.Stderr, "[Config] Validation failed with %d error(s)\n", len(validationResult.Errors))
		for _, err := range validationResult.Errors {
			fmt.Fprintf(os.Stderr, "[Config]   - %s: %s\n", err.Field, err.Message)
		}
		return GetDefaults(), validationResult
	}

	// T051: Log successful load
	fmt.Fprintf(os.Stderr, "[Config] Successfully loaded config\n")

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

	// T032: Validate keybindings
	result = c.validateKeybindings(result)

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

// T032: validateKeybindings validates and parses keybinding configuration
func (c *Config) validateKeybindings(result *ValidationResult) *ValidationResult {
	// If keybindings is empty, merge with defaults
	if len(c.Keybindings) == 0 {
		c.Keybindings = GetDefaults().Keybindings
	}

	// Step 1: Check for unknown actions
	actionErrors := ValidateActions(c.Keybindings)
	for _, errMsg := range actionErrors {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "keybindings",
			Message: errMsg,
		})
	}

	// Step 2: Check for conflicts (duplicate key assignments)
	conflicts := DetectKeybindingConflicts(c.Keybindings)
	for _, conflict := range conflicts {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "keybindings",
			Message: conflict,
		})
	}

	// Step 3: Parse each keybinding and store parsed versions
	c.ParsedKeybindings = make(map[string]*Keybinding)
	for action, keyStr := range c.Keybindings {
		key, r, mods, err := ParseKeybinding(keyStr)
		if err != nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("keybindings.%s", action),
				Message: fmt.Sprintf("invalid key format '%s': %v", keyStr, err),
			})
			continue
		}

		// Store parsed keybinding
		c.ParsedKeybindings[action] = &Keybinding{
			Action:    action,
			KeyString: keyStr,
			Key:       key,
			Rune:      r,
			Modifiers: mods,
		}
	}

	return result
}

// GetParsedKeybindings returns the parsed keybindings map
// Must call Validate() first to populate ParsedKeybindings
func (c *Config) GetParsedKeybindings() map[string]*Keybinding {
	if c.ParsedKeybindings == nil {
		return make(map[string]*Keybinding)
	}
	return c.ParsedKeybindings
}
