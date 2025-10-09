# TUI Black Screen Issue - Final Diagnosis
## Branch: `tui-tview-skeleton` vs `pre-reorg-tui`

**Date**: 2025-10-08
**Issue**: Black screen when running `./pass-cli.exe` with no arguments
**Expected**: TUI dashboard should render with sidebar, table, and status bar

---

## Architecture Discovery

### How TUI Launches (Both Branches)

**Entry Point**: `main.go` (root of project)

```go
func main() {
    // Detect if user provided a subcommand
    shouldUseTUI := true
    for arg in os.Args {
        if arg doesn't start with '-', it's a command
        shouldUseTUI = false
    }

    if shouldUseTUI {
        tui.Run(vaultPath)  // ← Launch TUI
    } else {
        cmd.Execute()  // ← Launch CLI
    }
}
```

**Behavior**:
- `./pass-cli.exe` → **TUI mode**
- `./pass-cli.exe --vault path` → **TUI mode** (flags only)
- `./pass-cli.exe add` → **CLI mode** (subcommand present)

---

## Branch Comparison

### pre-reorg-tui (WORKING)

**TUI Implementation**: Bubble Tea (old framework)
**Location**: `cmd/tui/` (Bubble Tea) + `cmd/tui-tview/` (tview prototype)

```
cmd/
├── tui/
│   ├── tui.go              # func Run() error - Bubble Tea entry point
│   ├── model.go            # Bubble Tea model
│   ├── commands.go         # Bubble Tea commands
│   └── views/              # Bubble Tea views
└── tui-tview/              # Separate tview prototype (not used by main.go)
    └── main.go             # package main, func main()
```

**Import in main.go**: `import "pass-cli/cmd/tui"`
**Function call**: `tui.Run()` (no parameters)
**Flow**:
1. main.go calls `tui.Run()`
2. `cmd/tui/tui.go` launches Bubble Tea TUI
3. Bubble Tea model handles all rendering
4. Works correctly ✅

---

### tui-tview-skeleton (BROKEN - BLACK SCREEN)

**TUI Implementation**: tview (new framework)
**Location**: `cmd/tui/` (tview only, Bubble Tea removed)

```
cmd/
└── tui/
    ├── main.go             # func Run(vaultPath string) error - tview entry
    ├── app.go              # tview.Application helpers
    ├── components/         # tview components (sidebar, table, detail, etc.)
    ├── events/             # Event handlers
    ├── layout/             # Layout manager
    ├── models/             # AppState, NavigationState
    └── styles/             # Theme and styling
```

**Import in main.go**: `import "pass-cli/cmd/tui"`
**Function call**: `tui.Run(vaultPath)` (vault path parameter)
**Flow**:
1. main.go calls `tui.Run(vaultPath)`
2. `cmd/tui/main.go` attempts to unlock vault and launch tview TUI
3. **BLACK SCREEN occurs here** ❌

---

## Function Signature Mismatch (FIXED)

### Initial Suspicion (RESOLVED)

| Branch | main.go calls | cmd/tui signature |
|--------|---------------|-------------------|
| pre-reorg-tui | `tui.Run()` | `func Run() error` ✅ Match |
| tui-tview-skeleton | `tui.Run(vaultPath)` | `func Run(vaultPath string) error` ✅ Match |

**Status**: ✅ Signatures are compatible - this is NOT the issue

---

## Initialization Sequence Analysis

### tui-tview-skeleton Initialization Flow

**File**: `cmd/tui/main.go`

```go
func Run(vaultPath string) error {
    // 1. Get vault path (use provided path or default)
    if vaultPath == "" {
        vaultPath = getDefaultVaultPath()
    }

    // 2. Initialize vault service
    vaultService, err := vault.New(vaultPath)
    if err != nil {
        return fmt.Errorf("failed to initialize vault service: %w", err)
    }

    // 3. Try keychain unlock first
    unlocked := false
    if err := vaultService.UnlockWithKeychain(); err == nil {
        unlocked = true
    } else {
        fmt.Println("Keychain unlock unavailable, prompting for password...")

        // Interactive prompt with limited attempts
        for attempt := 1; attempt <= 3; attempt++ {
            password, err := promptForPassword()  // ← POTENTIAL HANG POINT
            if err != nil {
                return fmt.Errorf("failed to read password: %w", err)
            }

            if err := vaultService.Unlock(password); err == nil {
                unlocked = true
                break
            }
        }
    }

    // 5. Launch TUI
    if err := LaunchTUI(vaultService); err != nil {  // ← POTENTIAL HANG POINT
        return fmt.Errorf("TUI error: %w", err)
    }
    return nil
}
```

**Potential Issues**:
1. **Password prompt hang**: `promptForPassword()` uses `gopass.GetPasswdMasked()` which may hang or corrupt terminal state
2. **LaunchTUI hang**: tview initialization may have an infinite loop or blocking call

---

### LaunchTUI Initialization Steps

```go
func LaunchTUI(vaultService *vault.VaultService) error {
    defer RestoreTerminal()  // Panic recovery

    styles.SetRoundedBorders()

    // 1. Create tview.Application
    app := NewApp()

    // 2. Initialize AppState with vault service
    appState := models.NewAppState(vaultService)

    // 3. Load credentials
    if err := appState.LoadCredentials(); err != nil {  // ← POTENTIAL HANG
        return fmt.Errorf("failed to load credentials: %w", err)
    }

    // 4. Create UI components
    sidebar := components.NewSidebar(appState)          // ← POTENTIAL HANG
    table := components.NewCredentialTable(appState)    // ← POTENTIAL HANG
    detailView := components.NewDetailView(appState)    // ← POTENTIAL HANG
    statusBar := components.NewStatusBar(app, appState) // ← POTENTIAL HANG

    // 5. Store components in AppState
    appState.SetSidebar(sidebar.TreeView)
    appState.SetTable(table.Table)
    appState.SetDetailView(detailView.TextView)
    appState.SetStatusBar(statusBar.TextView)

    // 6. Register callbacks (with QueueUpdateDraw wrappers)
    appState.SetOnCredentialsChanged(func() {
        app.QueueUpdateDraw(func() {
            sidebar.Refresh()
            table.Refresh()
            detailView.Refresh()
        })
    })

    // ... more callback registration ...

    // 7. Create NavigationState
    nav := models.NewNavigationState(app, appState)

    // 8. Create LayoutManager and build layout
    layoutMgr := layout.NewLayoutManager(app, appState)
    mainLayout := layoutMgr.CreateMainLayout()          // ← POTENTIAL HANG

    // 9. Create PageManager
    pageManager := layout.NewPageManager(app)

    // 10. Create EventHandler and setup shortcuts
    eventHandler := events.NewEventHandler(...)
    eventHandler.SetupGlobalShortcuts()                 // ← POTENTIAL HANG

    // 11. Set root primitive
    pageManager.ShowPage("main", mainLayout)
    app.SetRoot(pageManager.Pages, true)

    // 12. Set initial focus
    nav.SetFocus(models.FocusSidebar)

    // 13. Run application (BLOCKING)
    return app.Run()                                     // ← SHOULD BLOCK HERE (not hang)
}
```

**Analysis**:
- Steps 1-12 should execute quickly (< 100ms)
- Step 13 (`app.Run()`) is expected to block (it's the event loop)
- **If black screen occurs BEFORE UI renders**, hang is in steps 1-12
- **If black screen occurs AFTER UI should render**, issue is with rendering logic

---

## Likely Root Causes

### 1. Terminal State Corruption (HIGH PROBABILITY)

**Symptom**: Password prompt leaves terminal in raw/corrupted state, causing subsequent tview rendering to fail

**Evidence**:
- Password prompt uses `gopass.GetPasswdMasked()` which manipulates terminal flags
- If terminal isn't properly restored after password input, tview can't render

**Test**:
```bash
# If vault is unlocked via keychain (no password prompt), does TUI work?
./pass-cli.exe --vault test-vault/vault.enc
```

If TUI works when keychain unlocks but fails when password prompt is needed, this is the issue.

**Fix**: Add terminal state restoration after password prompt

---

### 2. Component Initialization Infinite Loop/Hang (MEDIUM PROBABILITY)

**Symptom**: One of the component constructors has an infinite loop or blocking call

**Suspects**:
- `components.NewSidebar(appState)` - Builds category tree
- `components.NewCredentialTable(appState)` - Builds credential table
- `layoutMgr.CreateMainLayout()` - Calls `rebuildLayout()` immediately

**Evidence**:
- `rebuildLayout()` was added in tui-tview-skeleton (line 106 of layout/manager.go)
- This could trigger premature layout calculations before tview app is running

**Test**: Add debug logging before each initialization step

---

### 3. Nil Pointer / Missing Component (LOW PROBABILITY)

**Symptom**: A component is nil and causes a panic, which is caught by `defer RestoreTerminal()` but leaves terminal in black screen state

**Evidence**:
- `defer RestoreTerminal()` exists, so panics are caught
- But terminal may not be properly restored on panic

**Test**: Check if error messages are being swallowed

---

## Recommended Diagnostic Steps

### Step 1: Test with Pre-Unlocked Vault (Bypass Password Prompt)

Set up keychain to avoid password prompt:

```bash
# On Windows (if working):
./pass-cli.exe --vault test-vault/vault.enc
```

**Expected**:
- ✅ If TUI renders: Password prompt is the issue
- ❌ If still black screen: Issue is in LaunchTUI initialization

---

### Step 2: Add Debug Logging

Edit `cmd/tui/main.go` and add logging:

```go
func LaunchTUI(vaultService *vault.VaultService) error {
    fmt.Fprintln(os.Stderr, "[DEBUG] LaunchTUI: Starting...")
    defer RestoreTerminal()

    fmt.Fprintln(os.Stderr, "[DEBUG] Setting rounded borders...")
    styles.SetRoundedBorders()

    fmt.Fprintln(os.Stderr, "[DEBUG] Creating app...")
    app := NewApp()

    fmt.Fprintln(os.Stderr, "[DEBUG] Creating AppState...")
    appState := models.NewAppState(vaultService)

    fmt.Fprintln(os.Stderr, "[DEBUG] Loading credentials...")
    if err := appState.LoadCredentials(); err != nil {
        return fmt.Errorf("failed to load credentials: %w", err)
    }

    fmt.Fprintln(os.Stderr, "[DEBUG] Creating sidebar...")
    sidebar := components.NewSidebar(appState)

    fmt.Fprintln(os.Stderr, "[DEBUG] Creating table...")
    table := components.NewCredentialTable(appState)

    // ... continue for all steps ...

    fmt.Fprintln(os.Stderr, "[DEBUG] Running app.Run() - should block here...")
    return app.Run()
}
```

**Run and observe**:
```bash
./pass-cli.exe --vault test-vault/vault.enc 2>&1 | tee debug.log
```

**Expected**: Debug output will show exactly where the hang occurs

---

### Step 3: Compare with pre-reorg-tui tview Standalone

The `pre-reorg-tui` branch has a standalone tview implementation in `cmd/tui-tview/main.go`. Try building and running that:

```bash
git checkout pre-reorg-tui
cd cmd/tui-tview
go build -o ../../tui-tview-standalone.exe
cd ../..
./tui-tview-standalone.exe
```

**Expected**:
- ✅ If this works: Confirms tview implementation is sound, issue is in integration
- ❌ If this also has black screen: tview implementation itself has bugs

---

## Key Differences Between Branches

| Aspect | pre-reorg-tui | tui-tview-skeleton |
|--------|---------------|-------------------|
| **TUI Framework** | Bubble Tea | tview |
| **Entry Function** | `func Run() error` | `func Run(vaultPath string) error` |
| **Password Handling** | Bubble Tea handles terminal | gopass.GetPasswdMasked() |
| **Component Creation** | Bubble Tea model init | Explicit component constructors |
| **Layout System** | Bubble Tea auto-layout | Custom LayoutManager with breakpoints |
| **Focus Management** | Bubble Tea built-in | Custom NavigationState |
| **Event Loop** | Bubble Tea (tea.Program) | tview (app.Run()) |

---

## Immediate Next Steps

1. **Test**: Run with vault path and enter password manually
   ```bash
   ./pass-cli.exe --vault test-vault/vault.enc
   # Enter: test1234
   ```

2. **Observe**: Note exactly when black screen appears:
   - During password prompt?
   - After entering password?
   - After "Launching TUI..." message (if any)?

3. **Add Debug Logging**: Follow Step 2 above to pinpoint hang location

4. **Report Findings**: Share debug output for further analysis

---

## Hypothesis Summary

**Most Likely**: Terminal state corruption from password prompt
**Secondary**: Component initialization hang (sidebar/table/layout)
**Least Likely**: Nil pointer causing silent panic

---

**Report Status**: Diagnostic steps ready for execution
**Next**: User to run tests and provide debug output
