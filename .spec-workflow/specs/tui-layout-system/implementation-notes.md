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

## Integration Test Results

**Test Date:** 2025-10-02

**Issues Found:**
1. ✅ **CONFIRMED: Concern #1 was valid** - Width subtraction bug existed
   - Rendered output was smaller than terminal dimensions
   - Pattern: losing 2 pixels per panel (Full=6, Medium=4, Small=2)

**Root Cause Discovered:**
Lipgloss `.Width(n)` behavior is **neither** "content width" nor "total width":
- **Actual behavior**: `.Width(n)` sets width to `n`, then adds BORDER but NOT PADDING
- With `RoundedBorder() + Padding(0,1)`:
  - Frame size = 4 (border=2 + padding=2)
  - `.Width(16)` produces total width = **18** (not 20!)
  - Adds border (2 chars) automatically, but NOT padding (2 chars)

**Correct Formula:**
```go
// WRONG (was using):
panelWidth = layout.Panel.Width - GetHorizontalFrameSize()  // Width - 4
// → Renders too small (subtracts border+padding, but Lipgloss adds back only border)

// CORRECT (now using):
panelWidth = layout.Panel.Width - paddingSize  // Width - 2
// → Renders exact (subtracts padding, Lipgloss adds border)
```

**Resolution:**
- Changed from `Width - 4` to `Width - 2` for all bordered panels
- Test `TestLipglossWidthBehavior` validates Lipgloss behavior empirically
- Test `TestExactDimensionMatch` ensures rendered dimensions match terminal exactly

**Verification:**
- [x] All borders visible at various terminal sizes
- [x] Rendered dimensions exactly match terminal dimensions
- [x] No pixel loss (TestDiagnosePanelWidths shows "Pixels lost: 0")
- [x] Focus changes don't affect dimensions (TestFocusDimensionsUnchanged passes)
- [x] All integration tests pass

**Concerns Status:**
1. ✅ **RESOLVED**: Double subtraction bug - fixed by using `Width - 2` instead of `Width - 4`
2. ✅ **VALIDATED**: Focus state frame sizes are identical - test confirms no dimension change on focus
3. ⏸️ **PENDING**: Component expectations - need real-world testing with actual vault data
