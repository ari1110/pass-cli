# Implementation Notes - TUI Layout System

## Potential Issues to Validate (Pre-Task 6)

These concerns were raised before completing integration tests. They should be validated through testing and real-world usage.

### 1. Double Subtraction Bug Concern

**Location:** `cmd/tui/helpers.go` - `recalculateLayout()` and `renderDashboardView()`

**Concern:** Potential double subtraction of frame overhead:
- In `recalculateLayout()`: Components receive `layout.Panel.ContentWidth` (already `Width - frameSize`)
- In `renderDashboardView()`: We calculate `contentWidth = layout.Panel.Width - horizontalFrame` again

**Expected behavior:**
- Component gets `ContentWidth` → renders content
- Apply border style with `.Width(contentWidth).Height(contentHeight)` where `contentWidth = Width - frame`
- Lipgloss adds frame back on top → total rendered = Width

**Test to verify:**
- Measure actual rendered panel width using `lipgloss.Width()`
- Should equal `layout.Panel.Width` exactly
- If smaller → double subtraction bug
- If larger → not subtracting enough

**Code paths to check:**
```go
// recalculateLayout() - line 134
m.sidebar.SetSize(layout.Sidebar.ContentWidth, layout.Sidebar.ContentHeight)

// renderDashboardView() - line 241
horizontalFrame := sidebarStyle.GetHorizontalFrameSize()
sidebarWidth := layout.Sidebar.Width - horizontalFrame  // ← Is this correct?
```

### 2. Focus State Frame Size Mismatch Concern

**Location:** `cmd/tui/components/layout_manager.go` vs `cmd/tui/helpers.go`

**Concern:** Frame size calculated once in layout manager but used with different styles in rendering:
- Layout manager uses `styles.ActivePanelBorderStyle.GetFrameSize()` for all panels
- Rendering uses actual panel style (active OR inactive based on focus)

**Assumption made:** Both `ActivePanelBorderStyle` and `InactivePanelBorderStyle` have identical frame sizes (both RoundedBorder + Padding(0,1))

**What if assumption is wrong:**
- Panels would have inconsistent sizing when focus changes
- Components would be sized for one frame size but rendered with another

**Test to verify:**
- Verify `TestGetFrameSizeValues` passes (already done ✓)
- Integration test: Change focus, measure rendered dimensions
- Should remain constant regardless of focus state

### 3. Component Expectations Concern

**Location:** All component `SetSize()` implementations

**Concern:** Components might expect to receive total allocated space and handle borders themselves, not pre-calculated content space.

**What changed:**
- **Before:** Components received `layout.Panel.Width` and `layout.Panel.Height` (total space including borders)
- **After:** Components receive `layout.Panel.ContentWidth` and `layout.Panel.ContentHeight` (space excluding borders)

**Potential issues:**
- If components have internal border handling, they now get less space than expected
- If components expect to render borders themselves, borders might be missing

**Components affected:**
- `SidebarPanel.SetSize()`
- `ListView.SetSize()` / `DetailView.SetSize()`
- `MetadataPanel.SetSize()`
- `ProcessPanel.SetSize()`
- `CommandBar.SetSize()`
- `Breadcrumb.SetSize()`

**Test to verify:**
- Check if components render correctly with new dimensions
- Verify no missing content or unexpected truncation
- Check that components don't try to add their own borders

## Validation Plan (Task 6)

Task 6 integration tests should specifically verify:

1. **Dimension Correctness:**
   - `lipgloss.Width(renderedPanel) == layout.Panel.Width` (exactly)
   - `lipgloss.Height(renderedPanel) == layout.Panel.Height` (exactly)
   - No overflow beyond terminal bounds

2. **Focus Invariance:**
   - Changing `panelFocus` should NOT change rendered dimensions
   - Only colors should change, not sizes

3. **Border Presence:**
   - All four borders visible (top, bottom, left, right)
   - Border characters present in rendered output (─, │, ╭, ╮, ╰, ╯)

4. **Component Rendering:**
   - Components render content correctly with ContentWidth/ContentHeight
   - No unexpected truncation or missing content

## Real-World Testing Results

_To be filled in after manual testing with actual TUI application_

**Test Date:** [To be filled]

**Issues Found:**
- [List any issues discovered during manual testing]

**Verification:**
- [ ] All borders visible at various terminal sizes (80x24, 100x30, 120x40, 140x50)
- [ ] No border overflow or cut-off
- [ ] Panels resize smoothly on terminal resize
- [ ] Focus changes don't affect dimensions
- [ ] Components render correctly with new sizing

**Resolution:**
- [Document how issues were resolved, if any]
