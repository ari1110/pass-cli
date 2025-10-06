# Requirements Document

## Introduction

The tui-tview-implementation feature converts the comprehensive architectural documentation in `cmd/tui-tview-skeleton/` into a fully functional tview-based Terminal User Interface (TUI) for pass-cli. This implementation will replace the existing Bubble Tea TUI, providing a modern, widget-based interface for managing credentials with improved component reusability, cleaner state management, and prevention of mutex deadlocks that plagued previous attempts.

The skeleton contains 14 detailed .md files documenting every aspect of the architecture: entry points, state management, UI components, layout system, event handling, and styling. This spec implements that complete design into production-ready Go code.

## Alignment with Product Vision

This feature directly supports pass-cli's product vision as outlined in product.md:

1. **Developer Experience**: Provides visual, interactive credential management with keyboard-driven navigation, complementing the script-friendly CLI commands
2. **Security First**: Implements secure vault unlocking with keychain integration and masked password input
3. **Developer Productivity**: Enables quick visual browsing, searching, and editing of credentials without leaving the terminal
4. **Features in Development**: Completes the "Interactive TUI Dashboard" feature explicitly listed as "In Progress" in product.md

The TUI dashboard enhances pass-cli from a pure CLI tool to a comprehensive credential management system with both command-line and interactive modes.

## Requirements

### Requirement 1: TUI Entry Point and Vault Unlocking

**User Story:** As a developer, I want to launch an interactive TUI dashboard, so that I can visually manage my credentials without using individual CLI commands.

#### Acceptance Criteria

1.1. WHEN the user runs `pass-cli tui` THEN the system SHALL attempt to unlock the vault using system keychain

1.2. IF keychain unlock fails THEN the system SHALL prompt for master password with masked input

1.3. WHEN the vault is successfully unlocked THEN the system SHALL initialize the tview application and display the main dashboard

1.4. IF vault unlocking fails after 3 attempts THEN the system SHALL exit with error code 1

1.5. WHEN the application encounters a panic THEN the system SHALL restore terminal state and display error message

### Requirement 2: Central State Management with Deadlock Prevention

**User Story:** As a developer implementing the TUI, I want thread-safe state management that prevents mutex deadlocks, so that the application remains responsive and never hangs.

#### Acceptance Criteria

2.1. WHEN any state mutation occurs THEN the system SHALL release all locks BEFORE invoking notification callbacks

2.2. WHEN credentials are loaded from the vault THEN the system SHALL update the in-memory cache and notify registered components

2.3. WHEN a component requests current state THEN the system SHALL provide thread-safe read access using RWMutex

2.4. WHEN the user selects a category or credential THEN the system SHALL update selection state and trigger selection change callbacks

2.5. IF an error occurs during vault operations THEN the system SHALL invoke error callbacks without holding locks

2.6. WHEN components are created THEN the system SHALL store single instances for reuse (no duplicate components)

### Requirement 3: Component-Based UI Architecture

**User Story:** As a user, I want a multi-panel dashboard with sidebar navigation, credential list, detail view, and status bar, so that I can efficiently navigate and manage my credentials.

#### Acceptance Criteria

3.1. WHEN the TUI launches THEN the system SHALL display a sidebar with category tree navigation using tview.TreeView

3.2. WHEN the TUI launches THEN the system SHALL display a credential list table using tview.Table

3.3. WHEN a credential is selected THEN the system SHALL display credential details in a dedicated panel using tview.TextView

3.4. WHEN the TUI is displayed THEN the system SHALL show a status bar with context-aware keyboard shortcuts using tview.TextView

3.5. WHEN the user requests to add/edit credentials THEN the system SHALL display modal forms using tview.Form and tview.Modal

3.6. WHEN components need to refresh THEN the system SHALL update UI by reading from centralized AppState (not maintaining duplicate state)

### Requirement 4: Responsive Layout Management

**User Story:** As a user working in different terminal sizes, I want the layout to adapt to my terminal width, so that I can use the TUI in various environments.

#### Acceptance Criteria

4.1. WHEN terminal width is less than 80 columns THEN the system SHALL hide the sidebar and metadata panel (list view only)

4.2. WHEN terminal width is between 80-120 columns THEN the system SHALL show sidebar and main content (no metadata panel)

4.3. WHEN terminal width exceeds 120 columns THEN the system SHALL show sidebar, main content, and metadata panel

4.4. WHEN terminal size changes THEN the system SHALL recalculate layout and redraw components

4.5. WHEN creating the layout THEN the system SHALL use tview.Flex for flexible sizing and tview.Pages for modal management

### Requirement 5: Global Event Handling and Keyboard Shortcuts

**User Story:** As a user, I want intuitive keyboard shortcuts for all actions, so that I can navigate and manage credentials efficiently without a mouse.

#### Acceptance Criteria

5.1. WHEN the user presses 'q' or Ctrl+C THEN the system SHALL gracefully exit the application

5.2. WHEN the user presses 'n' THEN the system SHALL open the add credential form modal

5.3. WHEN the user presses 'e' AND a credential is selected THEN the system SHALL open the edit credential form modal

5.4. WHEN the user presses 'd' AND a credential is selected THEN the system SHALL show delete confirmation dialog

5.5. WHEN the user presses 'Tab' THEN the system SHALL cycle focus between panels (sidebar → table → detail)

5.6. WHEN the user is in a form input field THEN the system SHALL NOT intercept keyboard shortcuts (let tview.Form handle input)

5.7. WHEN the user presses '/' THEN the system SHALL open a search/filter input

5.8. WHEN the user presses '?' THEN the system SHALL display the help screen with all keyboard shortcuts

### Requirement 6: Focus Management Between Components

**User Story:** As a user, I want clear visual indication of which panel is focused, so that I know where my keyboard input will be directed.

#### Acceptance Criteria

6.1. WHEN a panel gains focus THEN the system SHALL highlight its border with the active color

6.2. WHEN a panel loses focus THEN the system SHALL dim its border to the inactive color

6.3. WHEN the user presses Tab THEN the system SHALL move focus to the next panel in sequence

6.4. WHEN the user presses Shift+Tab THEN the system SHALL move focus to the previous panel in sequence

6.5. WHEN a modal is displayed THEN the system SHALL focus the modal and prevent focus changes to background panels

### Requirement 7: Credential CRUD Operations

**User Story:** As a user, I want to create, view, update, and delete credentials through the TUI, so that I can manage my vault without using CLI commands.

#### Acceptance Criteria

7.1. WHEN the user submits the add credential form THEN the system SHALL validate inputs, add to vault, refresh UI, and close modal

7.2. WHEN the user submits the edit credential form THEN the system SHALL update the credential in vault, refresh UI, and close modal

7.3. WHEN the user confirms deletion THEN the system SHALL remove credential from vault, refresh UI, and clear selection

7.4. WHEN a credential operation fails THEN the system SHALL display error message in a modal without crashing

7.5. WHEN credentials are modified THEN the system SHALL trigger onCredentialsChanged callback to refresh all affected components

### Requirement 8: Modern Styling and Theming

**User Story:** As a user, I want a visually appealing interface with modern styling, so that the TUI feels polished and professional.

#### Acceptance Criteria

8.1. WHEN panels are rendered THEN the system SHALL use rounded borders with consistent color palette

8.2. WHEN defining colors THEN the system SHALL use tcell.NewRGBColor() for precise color control

8.3. WHEN styling components THEN the system SHALL use colors defined in styles/theme.go for consistency

8.4. WHEN displaying the status bar THEN the system SHALL show keyboard shortcuts with visual highlighting

8.5. WHEN rendering modals THEN the system SHALL center them on screen with semi-transparent background

## Non-Functional Requirements

### Code Architecture and Modularity

- **Single Responsibility Principle**: Each file SHALL have one primary purpose (main.go for entry, state.go for state management, sidebar.go for sidebar component, etc.)
- **Modular Design**: Components SHALL be isolated and reusable (sidebar, table, detail view, forms as separate files)
- **Dependency Management**: State SHALL depend on vault service, components SHALL depend on state, layout SHALL depend on components
- **Clear Interfaces**: AppState SHALL provide explicit getters/setters for all state access, components SHALL provide Refresh() methods
- **No Circular Dependencies**: Import flow SHALL be main → models → components → styles (never circular)

### Performance

- **Startup Time**: TUI SHALL launch in under 500ms after vault is unlocked
- **Memory Usage**: Application SHALL use less than 50MB during normal operation
- **Responsiveness**: UI updates SHALL occur within 100ms of state changes
- **Smooth Navigation**: Keyboard input SHALL be responsive with no perceptible lag

### Security

- **Password Masking**: Master password input SHALL use masked input (howeyc/gopass or similar)
- **Credential Display**: Passwords in detail view MAY be masked by default with toggle to reveal
- **Secure Memory**: Credentials SHALL be cleared from memory when application exits
- **Keychain Integration**: SHALL use existing keychain service for automatic unlock (no new security code needed)

### Reliability

- **Error Handling**: ALL vault operations SHALL have error handling with user-friendly messages
- **Panic Recovery**: Application SHALL restore terminal state on panic
- **Graceful Shutdown**: Quitting SHALL clean up resources and restore terminal
- **State Consistency**: State mutations SHALL be atomic (complete success or complete rollback)

### Usability

- **Intuitive Navigation**: Keyboard shortcuts SHALL follow common TUI conventions (q=quit, n=new, e=edit, d=delete)
- **Visual Feedback**: Actions SHALL provide immediate visual feedback (modals, status updates)
- **Help Accessibility**: Help screen SHALL be accessible via '?' key
- **Consistent Theming**: All panels SHALL use consistent styling from theme.go

### Maintainability

- **Comprehensive Documentation**: Each file SHALL include purpose, responsibilities, dependencies in comments
- **Clear Component Boundaries**: No component SHALL modify another component's internal state directly
- **Deadlock Prevention**: ALL state mutations SHALL follow the lock → mutate → unlock → notify pattern
- **Testing**: Each component SHALL be testable in isolation with mock state

## Technical Constraints

1. **Framework**: SHALL use github.com/rivo/tview v0.42.0 for UI components
2. **Terminal Library**: SHALL use github.com/gdamore/tcell/v2 v2.9.0 for terminal control
3. **Vault Service**: SHALL reuse existing internal/vault package (no duplication)
4. **Go Version**: SHALL compile with Go 1.25.1
5. **Directory Structure**: SHALL implement in new cmd/tui-tview/ directory (skeleton location)
6. **Keychain Integration**: SHALL reuse existing internal/keychain package
7. **Password Input**: SHALL use existing howeyc/gopass or similar for masked input
8. **Cross-Platform**: SHALL work on Windows, macOS, and Linux terminals

## Implementation Scope

**In Scope**:
- Convert all 14 .md skeleton files to .go implementation files
- Implement main.go (entry point)
- Implement app.go (application lifecycle)
- Implement models/state.go (state management with deadlock prevention)
- Implement models/navigation.go (navigation state)
- Implement components/ (sidebar, table, detail, statusbar, forms)
- Implement layout/ (manager.go for responsive layout, pages.go for modal management)
- Implement events/ (handlers.go for shortcuts, focus.go for focus management)
- Implement styles/theme.go (color palette and styling)
- Register new `tui` command in cmd/root.go to launch tview TUI

**Out of Scope**:
- Removing existing Bubble Tea implementation (separate migration task)
- Advanced features (search, filter, sort - future enhancements)
- Multi-vault support (future enhancement)
- Theme customization (future enhancement)
- Exporting the TUI as a standalone binary (current architecture is fine)

## Success Metrics

1. **Functional Completeness**: All 14 skeleton .md files converted to working .go code
2. **Zero Deadlocks**: No mutex deadlocks during normal operation (verified by testing)
3. **All Requirements Met**: All acceptance criteria pass manual testing
4. **Cross-Platform**: TUI works on Windows Terminal, iTerm2, and gnome-terminal
5. **Clean Code**: Passes go fmt, go vet, golangci-lint with no errors
6. **Documentation**: All files have clear purpose and responsibility comments
