# TUI Keybinding Audit

## Global Keys (work in all states except when typing in inputs)

| Key | Action | Notes |
|-----|--------|-------|
| `?` or `F1` | Open help overlay | |
| `q` | Quit application | Not available when typing |
| `Ctrl+C` | Emergency quit | Works even when typing |

## Panel Controls (List & Detail states only)

| Key | Action | Notes |
|-----|--------|-------|
| `s` | Toggle sidebar panel | |
| `p` | Toggle process panel | Only shown if processes exist |
| `f` | Toggle all footer panels | Process + command bar |
| `Tab` | Next panel focus | Cycles through visible panels |
| `Shift+Tab` | Previous panel focus | |
| `:` | Open command bar | Vi-style command mode |

## State: List View

| Key | Action | Notes |
|-----|--------|-------|
| `/` | Activate search | |
| `a` | Add new credential | |
| `Enter` | View selected credential | Opens detail view |
| `↑/↓` | Navigate list | |
| `Esc` | Clear search | When search is active |

## State: Detail View

| Key | Action | Notes |
|-----|--------|-------|
| `m` | **CONFLICT** | See below |
| `c` | Copy password to clipboard | |
| `e` | Edit credential | |
| `d` | Delete credential | Opens confirmation |
| `Esc` | Back to list view | |

### 🔴 CONFLICT: 'm' key in Detail View

When focus is on **main panel**:
- `m` toggles **metadata panel visibility** (model.go:315)

When focus is on **other panels** OR metadata not focused:
- `m` toggles **password masking** (detail.go:61)

**This is confusing!** The same key does different things depending on which panel has focus.

## State: Add / Edit Forms

| Key | Action | Notes |
|-----|--------|-------|
| `Tab` | Next field | |
| `Shift+Tab` | Previous field | |
| `Ctrl+G` | Generate password | |
| `Ctrl+S` | Save credential | |
| `Esc` | Cancel (with confirmation) | Shows discard dialog if changes exist |

## State: Confirmation Dialogs

| Key | Action | Notes |
|-----|--------|-------|
| `y` | Confirm action | Simple yes/no |
| `n` | Cancel | |
| `Esc` | Cancel | |
| `Enter` | Confirm | For typed confirmations |

## State: Password Prompt

| Key | Action | Notes |
|-----|--------|-------|
| `Enter` | Unlock vault | |
| `Esc` | Quit application | |

## State: Command Bar (when open)

| Key | Action | Notes |
|-----|--------|-------|
| `Enter` | Execute command | |
| `Esc` | Close command bar | |
| `↑` | Previous command in history | |
| `↓` | Next command in history | |

## Recommendations

### Fix 'm' key conflict:
1. **Option A (Recommended)**: Use different key for metadata toggle
   - `i`: Toggle info/metadata panel
   - Keep `m`: Always toggle password mask in detail view

2. **Option B**: Make 'm' context-aware with clear visual feedback
   - Show in status bar which action 'm' will perform
   - Not recommended - confusing UX

### Standardize panel toggles:
- `s`: Sidebar
- `i`: Info/metadata panel (new)
- `p`: Process panel
- `f`: Footer (all bottom panels)

### Add missing keys:
- Consider adding `h` for help (in addition to `?` and `F1`)
- Consider adding `r` for refresh/reload credentials
