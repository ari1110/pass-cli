package main

import (
	"os"
	"pass-cli/cmd"
	"pass-cli/cmd/tui"
)

func main() {
	// If no arguments provided, launch TUI mode
	if len(os.Args) == 1 {
		if err := tui.Run(); err != nil {
			os.Exit(1)
		}
		return
	}

	// Otherwise, execute CLI mode
	cmd.Execute()
}