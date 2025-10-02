package main

import (
	"fmt"
	"os"
	"pass-cli/cmd"
	"pass-cli/cmd/tui"
)

func main() {
	// Launch TUI mode if no command is specified (only flags or nothing)
	// Commands: init, add, get, list, update, delete, generate, version
	shouldUseTUI := true

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		// If argument doesn't start with - or --, it's a command
		if arg != "" && arg[0] != '-' {
			shouldUseTUI = false
			break
		}
		// Skip the value for flags like --vault <path>
		if (arg == "--vault" || arg == "-v") && i+1 < len(os.Args) {
			i++ // Skip next argument (the vault path)
		}
	}

	if shouldUseTUI {
		if err := tui.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return
	}

	// Otherwise, execute CLI mode
	cmd.Execute()
}