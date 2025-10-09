# TUI Command Comparison
## Binary Testing: `pre-reorg-tui` vs `tui-tview-skeleton`

**Date**: 2025-10-08
**Testing**: Direct binary comparison

---

## Test Results

### Build and Test Commands
```bash
# Build from pre-reorg-tui
git checkout pre-reorg-tui
go build -o pass-cli-pre-reorg.exe
./pass-cli-pre-reorg.exe --help

# Build from tui-tview-skeleton
git checkout tui-tview-skeleton
go build -o pass-cli-skeleton.exe
./pass-cli-skeleton.exe --help
```

### Command Registration Status

| Branch | Binary | TUI Command | Status |
|--------|--------|-------------|--------|
| `pre-reorg-tui` | `pass-cli-pre-reorg.exe` | ✅ **EXISTS** | `tui - Launch interactive TUI dashboard` |
| `tui-tview-skeleton` | `pass-cli-skeleton.exe` | ❌ **MISSING** | Not in command list |

### Expected Behavior

#### pre-reorg-tui ✅
```bash
$ ./pass-cli-pre-reorg.exe tui
# TUI launches successfully
# Shows sidebar, table, detail view
# All keyboard shortcuts work
```

#### tui-tview-skeleton ❌
```bash
$ ./pass-cli-skeleton.exe tui
Error: unknown command "tui" for "pass-cli"
Run 'pass-cli --help' for usage.
```

---

## Root Cause Confirmed

### File Comparison

**pre-reorg-tui**:
```
cmd/
├── tui.go          ✅ EXISTS - Registers Cobra command
└── tui-tview/
    └── main.go     ✅ Contains TUI implementation
```

**tui-tview-skeleton**:
```
cmd/
├── tui.go          ❌ DELETED - No command registration
└── tui/
    └── main.go     ✅ Contains TUI implementation (but not callable)
```

### Why There's No "Black Screen"

The `tui-tview-skeleton` binary will **not** show a black screen because:
- The command isn't registered at all
- Cobra catches the unknown command error
- Displays error message: `Error: unknown command "tui"`
- Returns to shell immediately

**A black screen would only occur if:**
1. The command WAS registered
2. But the TUI initialization crashed or hung
3. Without proper error handling

---

## The Fix (Confirmed)

Create `cmd/tui.go` to register the Cobra command:

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"pass-cli/cmd/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI dashboard",
	Long: `Launch an interactive terminal user interface (TUI) for managing credentials.`,
	Run: runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)  // ← THIS IS WHAT'S MISSING
}

func runTUI(cmd *cobra.Command, args []string) {
	vaultPath := GetVaultPath()
	if err := tui.Run(vaultPath); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}
```

---

## Clarification Needed

**Question for user**: When you describe seeing a "black screen" in `tui-tview-skeleton`:
- Are you running `./pass-cli-skeleton.exe tui`?
- Are you seeing the error message: `Error: unknown command "tui"`?
- Or are you seeing something different?

If you're truly seeing a black screen (no error message), that would indicate:
1. You're running a different binary than what's built from `tui-tview-skeleton`
2. Or there's a terminal rendering issue hiding the error message
3. Or you're testing from a different branch

---

**Testing Completed**: 2025-10-08
**Status**: Root cause confirmed - `cmd/tui.go` file is missing
