package config

import (
	"github.com/gdamore/tcell/v2"
)

// Keybinding represents a parsed keybinding for runtime matching
type Keybinding struct {
	Action    string // Action name (e.g., "add_credential")
	KeyString string // Original string from config (e.g., "ctrl+a")

	// Parsed tcell representation
	Key       tcell.Key
	Rune      rune
	Modifiers tcell.ModMask
}
