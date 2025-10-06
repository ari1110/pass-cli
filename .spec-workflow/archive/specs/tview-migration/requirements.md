# Requirements Document

## Introduction

This specification defines the requirements for migrating the pass-cli TUI framework from Bubble Tea + Lipgloss to tview. Through systematic playground testing, we discovered fundamental layout management issues in Lipgloss that cause header/footer components to disappear at certain terminal heights due to manual height calculation errors. tview provides native Flex containers that eliminate these issues while simplifying the codebase.

**Purpose**: Replace the current Bubble Tea + Lipgloss TUI implementation with tview to achieve reliable layout rendering, eliminate off-by-one height errors, and reduce layout management complexity.

**Value**: Ensures the TUI dashboard remains functional and visually correct at all terminal sizes, improving user experience and reducing future maintenance burden.

## Alignment with Product Vision

This migration directly supports several product principles and objectives from product.md:

- **Developer Experience** (Product Principle #2): Eliminates frustrating layout bugs where header/footer disappear, improving daily usability
- **Monitoring & Visibility**: Preserves the multi-panel dashboard design while ensuring it works reliably at all terminal sizes
- **Technical Excellence**: Addresses technical debt by replacing manual layout calculations with a purpose-built layout engine
- **Offline Operation**: Maintains zero cloud dependencies while improving UI reliability

The migration aligns with the dashboard features outlined in product.md:
- Multi-Panel Layout (sidebar, main content, metadata, status bar)
- Responsive Design with terminal size adaptation
- Keyboard Navigation with context-aware shortcuts

## Requirements

### Requirement 1: Framework Migration

**User Story:** As a developer, I want the TUI to use tview instead of Bubble Tea + Lipgloss, so that layout rendering is reliable and maintainable.

#### Acceptance Criteria

1. WHEN the application is built THEN the system SHALL use github.com/rivo/tview and github.com/gdamore/tcell/v2 as TUI dependencies
2. WHEN the application is built THEN the system SHALL NOT include github.com/charmbracelet/bubbletea, github.com/charmbracelet/lipgloss, or github.com/charmbracelet/bubbles in the dependency tree
3. WHEN the TUI is initialized THEN the system SHALL use tview.Application as the main event loop
4. WHEN layout components are created THEN the system SHALL use tview.Flex containers for all multi-panel layouts

### Requirement 2: Layout Reliability

**User Story:** As a user, I want the header and footer to remain visible at all terminal sizes, so that I can always see navigation controls and context.

#### Acceptance Criteria

1. WHEN the terminal height is reduced to minimum size (20 lines) THEN the header SHALL remain visible
2. WHEN the terminal height is reduced to minimum size (20 lines) THEN the footer SHALL remain visible
3. WHEN vertical panel layouts are rendered THEN the system SHALL NOT produce height overflow errors
4. WHEN horizontal panel layouts are rendered THEN the system SHALL correctly distribute width according to flex ratios
5. IF the terminal size is below minimum (80x20) THEN the system SHALL display a "terminal too small" message

### Requirement 3: Functional Preservation

**User Story:** As a user, I want all existing TUI features to work exactly as before, so that my workflow is not disrupted.

#### Acceptance Criteria

1. WHEN the TUI is launched THEN the system SHALL display the sidebar, main content area, metadata panel, and status bar
2. WHEN I navigate with keyboard shortcuts THEN the system SHALL respond to all existing keybindings (h/j/k/l, arrows, enter, etc.)
3. WHEN I add/edit/delete credentials THEN the system SHALL perform operations identically to the Bubble Tea version
4. WHEN I search for credentials THEN the system SHALL filter results in real-time as before
5. WHEN I expand/collapse categories THEN the category tree SHALL function identically to the current implementation
6. WHEN vault operations complete THEN the system SHALL display success/error messages in the same manner

### Requirement 4: Responsive Layout

**User Story:** As a user, I want the dashboard to adapt to my terminal size, so that I can use it on different screen configurations.

#### Acceptance Criteria

1. WHEN terminal width is >= 120 columns THEN the system SHALL display all four panels (sidebar, main, metadata, status)
2. WHEN terminal width is 80-119 columns THEN the system SHALL hide the metadata panel and show three panels
3. WHEN terminal width is < 80 columns THEN the system SHALL display a minimum size warning
4. WHEN the terminal is resized THEN the system SHALL dynamically adjust panel widths according to flex ratios
5. WHEN the terminal is resized THEN the system SHALL maintain correct flex proportions (e.g., 1:2:1 for three panels)

### Requirement 5: Component Architecture

**User Story:** As a developer, I want TUI components to be modular and testable, so that I can maintain and extend the codebase easily.

#### Acceptance Criteria

1. WHEN implementing sidebar THEN the system SHALL create a dedicated sidebar component with tview primitives
2. WHEN implementing views (list, detail, forms) THEN the system SHALL use tview components (TextView, Form, List, etc.)
3. WHEN implementing layout management THEN the system SHALL centralize responsive logic in LayoutManager
4. WHEN implementing status bar THEN the system SHALL use tview.TextView with dynamic content updates
5. IF a component needs custom rendering THEN the system SHALL use tview's SetDrawFunc for custom draw logic

### Requirement 6: State Management

**User Story:** As a developer, I want a clear state management pattern, so that UI updates are predictable and debuggable.

#### Acceptance Criteria

1. WHEN the TUI initializes THEN the system SHALL maintain a central Model struct containing all application state
2. WHEN user input is received THEN the system SHALL update state through event handlers
3. WHEN state changes occur THEN the system SHALL trigger re-renders of affected components
4. WHEN vault operations are in progress THEN the system SHALL reflect loading states in the UI
5. WHEN navigation occurs THEN the system SHALL update currentView state and re-render appropriate panels

### Requirement 7: Migration Compatibility

**User Story:** As a developer, I want the migration to be completed without breaking existing functionality, so that users experience zero disruption.

#### Acceptance Criteria

1. WHEN the migration is complete THEN all existing CLI commands SHALL continue to work unchanged
2. WHEN the TUI is launched THEN the entry point in main.go SHALL correctly route to tview implementation
3. WHEN the migration is complete THEN internal/vault service SHALL remain unchanged
4. WHEN the migration is complete THEN all tests in cmd/tui/ SHALL pass with tview implementation
5. IF any existing tests fail THEN they SHALL be updated to reflect tview's event model

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Each TUI component (sidebar, panels, views) shall have a single, well-defined purpose
- **Modular Design**: tview components shall be isolated and composable, following tview's primitive-based architecture
- **Dependency Management**: TUI layer shall depend only on tview/tcell and internal/vault service
- **Clear Interfaces**: Define clean contracts between view components and state management

### Performance
- **Startup Time**: TUI initialization shall complete in <200ms (same or better than current implementation)
- **Render Performance**: Layout recalculation on resize shall complete in <50ms
- **Memory Usage**: tview implementation shall use <=50MB memory during normal operation
- **Input Latency**: Keyboard input shall be processed in <16ms for 60fps responsiveness

### Security
- **Sensitive Data Display**: Password masking and secure display shall function identically to current implementation
- **Clipboard Operations**: Credential copying shall use the same security measures as current implementation
- **No Data Leakage**: tview implementation shall not log or expose sensitive credential data

### Reliability
- **Error Handling**: All tview operations shall handle errors gracefully without crashing
- **Terminal Compatibility**: TUI shall work on Windows Terminal, iTerm2, macOS Terminal, and major Linux terminals
- **Recovery**: Layout errors shall not crash the application; show error state and allow retry
- **Regression Prevention**: All existing TUI test cases shall pass or be updated to reflect tview patterns

### Usability
- **Learning Curve**: Existing users shall not need to learn new navigation patterns
- **Visual Consistency**: Color scheme, borders, and styling shall match current dashboard design
- **Keyboard Shortcuts**: All existing shortcuts shall work identically (or be improved if tview enables better UX)
- **Help Documentation**: Help screen shall be updated to reflect any tview-specific improvements
