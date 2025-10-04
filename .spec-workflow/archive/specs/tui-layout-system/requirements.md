# Requirements Document

## Introduction

The TUI dashboard currently has persistent layout issues where panel borders are cut off or missing regardless of terminal window size. Users cannot see the top borders of components, right-side borders are truncated, and the layout feels inconsistent and unprofessional. This spec addresses the root cause: **lack of systematic layout calculation that properly accounts for Lipgloss border and padding overhead across all panels**.

Instead of applying ad-hoc fixes to individual panels, we need a **unified, reusable layout calculation system** that ensures all panels (sidebar, main, metadata, process, command bar, status bar) render perfectly within terminal bounds with all borders visible.

## Alignment with Product Vision

From `product.md`:
- **Key Feature #3 - Dual Interface Design**: "Interactive TUI Dashboard: Full-featured terminal UI with multi-panel layout, category tree navigation, real-time search, and visual credential management" - Cut-off borders undermine this visual management promise
- **Dashboard Features**: "Responsive Design: Adaptive layout based on terminal size (breakpoints at 80, 120 columns)" - Current layout fails to adapt properly, with borders extending beyond bounds
- **Product Principles**: "Developer Experience: Design for speed, simplicity, and CLI integration" - Users should focus on credentials, not fighting with cut-off borders

From `tech.md`:
- **TUI Layer Architecture**: "Responsive Layout System: Multi-panel dashboard with breakpoint-based adaptation" and "Lipgloss-based styling with consistent color scheme and borders" - Current implementation has inconsistent border handling
- **Decision #10 - Lipgloss Layout System**: "Automatic width/height calculation with border and padding support" and "Composable layout primitives" - We're not properly using Lipgloss's layout primitives, leading to manual overhead calculations
- **Go Best Practices**: Named constants, clear interfaces, encapsulation - Current layout code violates these with magic numbers and duplicated logic

From `structure.md`:
- **TUI Components**: Layout Manager documented as "Responsive dimension calculations" - Current implementation doesn't properly account for border overhead in those calculations

This feature ensures the TUI is **professional, consistent, and properly bounded** within terminal dimensions, reflecting the quality standards expected of a developer tool.

## Requirements

### Requirement 1: Systematic Border Overhead Calculation

**User Story:** As a developer using pass-cli TUI, I want all panel borders to be fully visible regardless of terminal size, so that the interface looks professional and complete.

#### Acceptance Criteria

1. WHEN the TUI renders THEN all panels (sidebar, main, metadata) SHALL display complete borders on all four sides (top, bottom, left, right)
2. WHEN the user resizes the terminal THEN all panel borders SHALL remain fully visible within the new dimensions
3. WHEN calculating panel dimensions THEN the system SHALL use a centralized function that accounts for Lipgloss border and padding overhead
4. IF a panel has border style with padding THEN the system SHALL subtract overhead (border: 2 chars, padding: 2 chars per side = 4 total width, 2 total height) from the allocated dimension before applying the Lipgloss style
5. WHEN a panel style is applied with `.Width(n)` THEN Lipgloss SHALL add decorations on top, resulting in total rendered width of `n + overhead`

### Requirement 2: Use Lipgloss Built-in Frame Size Methods

**User Story:** As a developer maintaining the TUI codebase, I want to use Lipgloss's built-in frame size methods to calculate border overhead, so that I leverage proven library functionality instead of custom calculations.

#### Acceptance Criteria

1. WHEN calculating border overhead THEN the system SHALL use `style.GetHorizontalFrameSize()` and `style.GetVerticalFrameSize()` methods
2. WHEN rendering any bordered panel THEN the code SHALL call frame size methods on the specific style being applied (ActivePanelBorderStyle or InactivePanelBorderStyle)
3. WHEN calculating panel content dimensions THEN the function SHALL subtract overhead using: `contentWidth = totalWidth - style.GetHorizontalFrameSize()`
4. IF a new panel is added with different border/padding configuration THEN the frame size methods SHALL automatically return correct overhead without code changes
5. WHEN reviewing the code THEN there SHALL be zero hardcoded border overhead constants (no `const borderWidth = 4`)

### Requirement 3: Layout Manager Enhancement

**User Story:** As a developer using pass-cli TUI, I want the layout manager to provide complete panel specifications including border overhead, so that rendering is consistent and predictable.

#### Acceptance Criteria

1. WHEN layout manager calculates dimensions THEN it SHALL return both `TotalWidth/Height` (allocated space) and `ContentWidth/Height` (space available for content after border overhead)
2. WHEN a panel is rendered THEN it SHALL use `ContentWidth/Height` for internal sizing and `TotalWidth/Height` for validation
3. WHEN border styles change (e.g., from thick to normal border) THEN the overhead calculation SHALL adapt automatically
4. IF terminal is too small THEN the layout manager SHALL provide clear minimum dimensions that account for all border overhead
5. WHEN debugging layout issues THEN developers SHALL be able to inspect both allocated and content dimensions for each panel

### Requirement 4: Comprehensive Layout Testing

**User Story:** As a developer maintaining the TUI, I want automated tests that verify layout calculations, so that future changes don't reintroduce border overflow bugs.

#### Acceptance Criteria

1. WHEN running unit tests THEN there SHALL be tests verifying border overhead calculation for all panel types
2. WHEN testing `renderBorderedPanel()` THEN tests SHALL verify that rendered output matches allocated dimensions exactly
3. WHEN testing layout manager THEN tests SHALL verify content dimensions equal total dimensions minus overhead
4. IF a panel style is modified THEN tests SHALL catch any mismatch between allocated and rendered sizes
5. WHEN tests run THEN they SHALL cover edge cases: minimum terminal size, single-panel layout, all-panels-visible layout

### Requirement 5: Clear Documentation and Constants

**User Story:** As a developer maintaining the TUI, I want clear documentation explaining Lipgloss width behavior and border overhead, so that I don't make the same mistakes again.

#### Acceptance Criteria

1. WHEN reviewing layout code THEN there SHALL be inline comments explaining: "Lipgloss .Width(n) sets CONTENT width, then ADDS borders/padding on top"
2. WHEN border overhead is calculated THEN it SHALL use named constants (e.g., `BorderWidth`, `PaddingWidth`) instead of magic numbers
3. WHEN a developer reads `renderBorderedPanel()` THEN the function SHALL have documentation explaining the sizing model
4. IF border or padding configuration changes THEN constants SHALL be updated in one central location
5. WHEN onboarding a new developer THEN layout documentation SHALL be available in code comments and/or a dedicated layout design doc

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Border overhead calculation logic shall be in one reusable function, not duplicated across helpers.go
- **Modular Design**: Panel rendering shall be abstracted into a reusable utility that works for all panel types
- **Dependency Management**: Layout manager shall not depend on specific panel implementations (sidebar, metadata, etc.)
- **Clear Interfaces**: Panel rendering function shall have a clear contract: accepts allocated dimensions and content, returns properly sized output

### Performance
- Layout calculations must remain under 5ms for typical terminal sizes (80x24 to 200x60)
- Border overhead calculation must be O(1) complexity
- Panel rendering must not allocate excessive memory (avoid unnecessary string copies)

### Reliability
- Layout must never cause borders to extend beyond terminal bounds
- Layout must handle edge cases gracefully: minimum size, maximum size, asymmetric panel configurations
- Border rendering must be idempotent: same inputs produce same outputs

### Usability
- All panel borders must be visible and complete in typical terminal sizes (100x30 minimum)
- Layout must adapt smoothly to terminal resize events
- Users should never see truncated borders, cut-off content, or layout overflow

### Maintainability
- Future panel additions must require <10 lines of code to integrate with layout system
- Border overhead constants must be centralized and easily adjustable
- Layout logic must be testable without requiring actual terminal rendering
