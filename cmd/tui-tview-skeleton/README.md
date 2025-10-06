# TUI-tview Architecture

## Overview

This directory contains the tview-based Terminal User Interface (TUI) implementation for pass-cli. This is a modern, widget-based TUI that uses the rivo/tview framework for layout and components.

## Architecture Philosophy

**Clean separation of concerns:**
- **Models**: Application state and business logic
- **Components**: Reusable UI widgets (sidebar, table, forms)
- **Layout**: Responsive layout management and view composition
- **Events**: Keyboard shortcuts and event handling
- **Styles**: Visual theming and colors

**Key principles:**
- No mutex deadlocks: State is managed cleanly without callback-while-locked patterns
- Component reuse: Single instances of components managed by application state
- Explicit dependencies: Each file documents what it uses and why
- tview primitives: Leverage Flex, Grid, Pages, TreeView, Table, Form

## Directory Structure

```
cmd/tui-tview/
├── README.md              # This file - architecture overview
├── main.go.md            # Entry point and initialization flow
├── app.go.md             # tview.Application lifecycle management
├── models/
│   ├── state.go.md       # Central application state
│   └── navigation.go.md  # Navigation state (selected items, current view)
├── components/
│   ├── sidebar.go.md     # Category tree sidebar (tview.TreeView)
│   ├── table.go.md       # Credential list table (tview.Table)
│   ├── detail.go.md      # Credential detail view (tview.TextView)
│   ├── statusbar.go.md   # Status bar with shortcuts (tview.TextView)
│   └── forms.go.md       # Add/Edit credential forms (tview.Form + Modal)
├── layout/
│   ├── manager.go.md     # Responsive layout with breakpoints (tview.Flex)
│   └── pages.go.md       # Page/modal management (tview.Pages)
├── events/
│   ├── handlers.go.md    # Global keyboard shortcuts
│   └── focus.go.md       # Focus management between components
└── styles/
    └── theme.go.md       # Color palette and styling (tcell.Color)
```

## Data Flow

```
User Input
    ↓
events/handlers.go (keyboard shortcuts)
    ↓
models/state.go (update state)
    ↓
components/* (refresh UI)
    ↓
layout/manager.go (compose layout)
    ↓
Screen
```

## Key Design Decisions

### 1. State Management
- **Single source of truth**: `models/state.go` holds all application state
- **No callbacks while locked**: Avoid mutex deadlocks by releasing locks before callbacks
- **Explicit updates**: Components pull from state, don't push updates

### 2. Component Architecture
- **Single instances**: One sidebar, one table, one detail view (no duplication)
- **Managed by state**: Components created once, stored in state, reused
- **tview primitives**: Use TreeView, Table, Form instead of building custom

### 3. Layout System
- **Responsive**: Three breakpoints (small <80, medium 80-120, large >120 columns)
- **Flex-based**: Use tview.Flex for automatic sizing
- **Pages for modals**: Use tview.Pages to show/hide modal dialogs

### 4. Event Handling
- **Global shortcuts**: SetInputCapture for app-wide keys (q, n, e, d)
- **Component-specific**: Let input components (Form, InputField) handle their own keys
- **Focus-aware**: Check focused component type before intercepting shortcuts

### 5. Styling
- **Modern aesthetic**: Rounded borders, gradient colors, padding
- **tcell colors**: Use tcell.NewRGBColor() for precise color control
- **Consistent palette**: Define colors once in theme.go, reuse everywhere

## Integration Points

### Vault Service
- `internal/vault/vault.go` - Business logic for credential management
- State wraps vault service, provides thread-safe access
- Components never call vault directly, always through state

### Existing CLI Commands
- TUI is a separate entry point from CLI commands
- Shares the same vault service and internal packages
- No code duplication, clean separation

## Reference Implementation

See `cmd/tview-playground/main.go` for working examples of:
- Flex layouts with multiple panels
- Rounded borders and modern styling
- Keyboard navigation
- Page switching

## Lessons Learned from Previous Attempts

**Avoid:**
- ❌ Calling callbacks while holding mutex locks
- ❌ Creating duplicate component instances (LayoutManager creating new sidebar)
- ❌ Intercepting keyboard input for Form/InputField components
- ❌ Forgetting to call refreshAllComponents() after state changes
- ❌ Mixing architectural concerns (UI logic in state, business logic in components)

**Do:**
- ✅ Release locks BEFORE calling callbacks
- ✅ Create components once, store in state, reuse
- ✅ Check focused component type before intercepting keyboard
- ✅ Refresh UI after every state mutation
- ✅ Keep clean separation: state = data, components = view, events = controller

## Development Workflow

1. **Design phase** (current): Document each file's purpose with .md files
2. **Review phase**: Ensure architecture makes sense, no deadlocks possible
3. **Implementation phase**: Convert .md → .go files one at a time
4. **Testing phase**: Verify each component works before moving to next

## Next Steps

Once skeleton is reviewed and approved:
1. Convert main.go.md → main.go (entry point)
2. Convert app.go.md → app.go (application setup)
3. Convert models/state.go.md → state.go (state management)
4. Convert components/*.md → *.go (UI components)
5. Convert layout/*.md → *.go (layout management)
6. Convert events/*.md → *.go (event handling)
7. Convert styles/theme.go.md → theme.go (styling)
8. Wire everything together and test
