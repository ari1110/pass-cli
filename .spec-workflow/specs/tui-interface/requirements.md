# Requirements Document: TUI Interface

## Introduction

Pass-CLI currently provides a robust command-line interface (CLI) for managing encrypted credentials. While the CLI excels for power users and scripting scenarios, it presents a learning curve for new users who must remember command syntax and options. This feature adds a Terminal User Interface (TUI) mode that provides a visual, interactive, keyboard-driven interface for credential management.

**Purpose**: Enhance user experience by providing an intuitive, discoverable interface while maintaining the existing CLI for power users and automation workflows.

**Value to Users**:
- **Discoverability**: Visual menus and help screens eliminate need to memorize commands
- **Safety**: Confirmation dialogs prevent accidental deletions
- **Productivity**: Real-time search and keyboard navigation for faster workflows
- **Accessibility**: Visual feedback and structured layouts improve usability
- **Complementary**: TUI and CLI coexist without conflict - users choose their preferred mode

## Alignment with Product Vision

This feature directly supports multiple principles outlined in product.md:

**Developer Experience (Principle #2)**: The TUI provides an alternate interface optimized for interactive sessions, reducing cognitive load for credential browsing and management tasks.

**Privacy by Design (Principle #3)**: The TUI maintains the same local-only storage model with no telemetry or network dependencies. All security properties remain unchanged.

**Minimal Dependencies (Principle #5)**: Using Bubble Tea framework adds minimal binary overhead (~1-2MB) while providing professional TUI capabilities without reinventing complex terminal handling.

**Product Positioning**: Enhances Pass-CLI's unique position as "the only pure CLI password manager with native OS keychain integration" by adding an exceptional user experience layer that competitors lack.

**Success Metrics Alignment**:
- **Usability**: Improved onboarding for new users reduces learning curve
- **Quality**: Component-based architecture maintains 90%+ test coverage target
- **Adoption**: Enhanced UX may contribute to GitHub star growth target

## Requirements

### Requirement 1: TUI Mode Activation and Navigation

**User Story**: As a developer, I want to launch TUI mode by running `pass-cli` without arguments, so that I can visually browse and manage my credentials without remembering command syntax.

#### Acceptance Criteria

1. WHEN user executes `pass-cli` without command arguments THEN the system SHALL launch interactive TUI mode
2. WHEN user executes `pass-cli` with any command arguments (e.g., `pass-cli list`) THEN the system SHALL execute CLI mode (existing behavior unchanged)
3. WHEN TUI mode is active THEN the system SHALL support keyboard navigation using:
   - Arrow keys (↑/↓/←/→) for navigation
   - Vim-style keys (h/j/k/l) as alternative navigation
   - Tab key for field/component switching
   - Enter key for selection/confirmation
   - Escape key for cancel/back actions
4. WHEN user presses `q` or `Ctrl+C` in TUI main view THEN the system SHALL exit gracefully and return to shell
5. WHEN terminal window is resized THEN the TUI SHALL redraw layout to fit new dimensions without crashing

### Requirement 2: Vault Unlock in TUI Mode

**User Story**: As a user, I want the TUI to handle vault unlocking with keychain support, so that I have a seamless experience consistent with CLI mode.

#### Acceptance Criteria

1. WHEN TUI launches and vault is not initialized THEN the system SHALL display a message with instructions to run `pass-cli init`
2. WHEN TUI launches and vault exists but is locked THEN the system SHALL:
   - First attempt to unlock using OS keychain (if configured)
   - IF keychain unlock succeeds THEN proceed to credential list view
   - IF keychain unavailable or fails THEN prompt user for master password
3. WHEN user enters master password THEN the system SHALL:
   - Validate password using existing vault service
   - IF password correct THEN unlock vault and display credential list
   - IF password incorrect THEN display error message and allow retry
4. WHEN vault unlock fails after 3 attempts THEN the system SHALL exit TUI mode with error message

### Requirement 3: Credential List View with Search

**User Story**: As a user, I want to see all my credentials in a searchable, sortable list, so that I can quickly find the credential I need without typing exact service names.

#### Acceptance Criteria

1. WHEN vault is unlocked THEN the system SHALL display a list view showing all credentials with columns:
   - Service name (primary identifier)
   - Username
   - Last updated timestamp
   - Usage indicator (if credential has been accessed)
2. WHEN credential list is displayed THEN the system SHALL provide a search bar at the top of the screen
3. WHEN user types in search bar THEN the system SHALL filter the credential list in real-time matching against:
   - Service name (case-insensitive substring match)
   - Username (case-insensitive substring match)
4. WHEN search query matches no credentials THEN the system SHALL display "No credentials found matching '[query]'" message
5. WHEN user presses `/` key from list view THEN the system SHALL focus the search bar
6. WHEN user presses `Escape` in search bar THEN the system SHALL clear search and return focus to credential list
7. WHEN credential list has more items than fit on screen THEN the system SHALL provide scrollable view with scroll indicators

### Requirement 4: Credential Detail View

**User Story**: As a user, I want to view full credential details with the ability to copy password to clipboard, so that I can use credentials without exposing them on screen.

#### Acceptance Criteria

1. WHEN user selects a credential from list (Enter key) THEN the system SHALL display detail view showing:
   - Service name
   - Username
   - Password (masked as asterisks by default)
   - Notes (if present)
   - Created timestamp
   - Updated timestamp
   - Usage records (locations and access counts)
2. WHEN detail view is active and user presses `c` THEN the system SHALL copy password to clipboard and display "Password copied to clipboard" confirmation
3. WHEN detail view is active and user presses `m` THEN the system SHALL toggle password visibility (show plaintext / mask)
4. WHEN detail view is active and user presses `Escape` or `Backspace` THEN the system SHALL return to credential list view
5. WHEN password is copied to clipboard THEN the system SHALL use existing clipboard utilities (consistent with CLI behavior)
6. WHEN usage records exist THEN the system SHALL display them in a formatted table showing:
   - Location (working directory path)
   - Last accessed timestamp
   - Access count
   - Git repository (if applicable)

### Requirement 5: Add Credential Interactive Form

**User Story**: As a user, I want to add new credentials through an interactive form with validation, so that I don't need to remember command syntax or flags.

#### Acceptance Criteria

1. WHEN user presses `a` key from list view THEN the system SHALL display "Add Credential" form with fields:
   - Service name (required)
   - Username (optional)
   - Password (optional - can be generated)
   - Notes (optional, multi-line)
2. WHEN user is in add form THEN the system SHALL:
   - Support Tab key to move between fields
   - Support Shift+Tab to move backward between fields
   - Highlight currently focused field
3. WHEN user is in password field and presses `g` THEN the system SHALL:
   - Generate a secure password using existing password generation logic
   - Display "Password generated (20 characters)" confirmation
   - Populate password field with generated value
4. WHEN user completes form and presses `Ctrl+S` or final Enter THEN the system SHALL:
   - Validate that service name is not empty
   - Validate that service name doesn't already exist
   - IF validation passes THEN save credential using VaultService and return to list view
   - IF validation fails THEN display inline error message next to invalid field
5. WHEN user presses `Escape` in add form THEN the system SHALL:
   - Display confirmation dialog "Discard new credential? (y/n)"
   - IF user confirms THEN discard and return to list view
   - IF user cancels THEN return to form editing
6. WHEN credential is successfully added THEN the system SHALL:
   - Display "Credential added successfully" notification
   - Refresh credential list to show new entry
   - Select the newly added credential in the list

### Requirement 6: Update Credential Interactive Form

**User Story**: As a user, I want to update existing credentials through an editable form, so that I can correct mistakes or rotate passwords easily.

#### Acceptance Criteria

1. WHEN user presses `e` key on selected credential in list or detail view THEN the system SHALL display "Edit Credential" form pre-filled with current values:
   - Service name (read-only, cannot be changed)
   - Username (editable)
   - Password (editable, masked)
   - Notes (editable, multi-line)
2. WHEN user modifies fields and presses `Ctrl+S` THEN the system SHALL:
   - Update credential using VaultService
   - Update the UpdatedAt timestamp
   - Display "Credential updated successfully" notification
   - Return to detail view showing updated values
3. WHEN user presses `Escape` in edit form THEN the system SHALL:
   - IF changes were made THEN display confirmation dialog "Discard changes? (y/n)"
   - IF no changes were made THEN immediately return to detail view
4. WHEN password field is focused and user presses `g` THEN the system SHALL generate new password (same as add form)
5. WHEN usage records exist for credential being edited THEN the system SHALL display warning:
   - "This credential is used in [N] locations. Update will affect existing usage."
   - Show usage locations in warning message

### Requirement 7: Delete Credential with Confirmation

**User Story**: As a user, I want to delete credentials with confirmation dialogs and usage warnings, so that I don't accidentally remove credentials still in use.

#### Acceptance Criteria

1. WHEN user presses `d` key on selected credential in list or detail view THEN the system SHALL display confirmation dialog:
   - "Delete '[service name]'?" message
   - Username displayed for verification
   - Warning if usage records exist: "Used in [N] locations: [list]"
   - Options: "Yes (y) / No (n)"
2. WHEN user confirms deletion (presses `y`) THEN the system SHALL:
   - Delete credential using VaultService
   - Display "Credential deleted successfully" notification
   - Return to credential list with next item selected
3. WHEN user cancels deletion (presses `n` or `Escape`) THEN the system SHALL:
   - Close confirmation dialog
   - Return to previous view (list or detail) with same credential selected
4. WHEN credential has usage records THEN the confirmation dialog SHALL:
   - Display usage locations prominently in red/warning color
   - Require explicit confirmation ("Type service name to confirm: ____")
   - IF typed service name doesn't match THEN prevent deletion and show error

### Requirement 8: Status Bar with System Indicators

**User Story**: As a user, I want to see system status information in a persistent status bar, so that I know keychain status, number of credentials, and available keyboard shortcuts.

#### Acceptance Criteria

1. WHEN TUI is active THEN the system SHALL display a status bar at bottom of screen showing:
   - Keychain status indicator (icon/text: "🔓 Keychain" if active, "🔒 Password" if not)
   - Total credential count (e.g., "15 credentials")
   - Current view name (e.g., "List", "Detail", "Add", "Edit")
   - Common keyboard shortcuts for current view
2. WHEN keychain is available and enabled THEN status bar SHALL display "🔓 Keychain" in green
3. WHEN keychain is not available or not enabled THEN status bar SHALL display "🔒 Password" in yellow
4. WHEN credential count changes (add/delete) THEN status bar SHALL update count in real-time
5. WHEN user switches views THEN status bar SHALL update to show context-relevant keyboard shortcuts

### Requirement 9: Help Overlay with Keyboard Shortcuts

**User Story**: As a user, I want to view all available keyboard shortcuts in a help overlay, so that I can discover features without reading documentation.

#### Acceptance Criteria

1. WHEN user presses `?` or `F1` key from any view THEN the system SHALL display help overlay showing:
   - Global shortcuts (q: quit, ?: help, Esc: back)
   - View-specific shortcuts for current screen
   - Navigation keys (arrows, vim keys, tab)
   - Action keys (enter, space, letters)
2. WHEN help overlay is displayed THEN the system SHALL:
   - Dim background content (semi-transparent overlay)
   - Display shortcuts in categorized groups
   - Show shortcuts in two columns: "Key" and "Action"
3. WHEN user presses any key while help is open THEN the system SHALL close help overlay and return to previous view
4. WHEN help overlay is too large for screen THEN the system SHALL make it scrollable with scroll indicators

### Requirement 10: Theme and Visual Styling

**User Story**: As a user, I want a visually consistent, professional interface with good contrast, so that I can use the TUI comfortably for extended periods.

#### Acceptance Criteria

1. WHEN TUI is active THEN the system SHALL use a cohesive color scheme with:
   - Distinct colors for selected/focused items
   - Subtle colors for borders and decorations
   - High contrast for text readability
   - Warning colors (red/orange) for destructive actions
   - Success colors (green) for confirmations
2. WHEN terminal supports 256 colors THEN the system SHALL use full color palette
3. WHEN terminal supports only basic colors THEN the system SHALL gracefully degrade to ANSI colors
4. WHEN text is too long for field width THEN the system SHALL truncate with ellipsis ("...")
5. WHEN displaying sensitive data (passwords) THEN the system SHALL:
   - Use asterisks (***) for masking by default
   - Use monospace font for password display
   - Maintain consistent field widths

### Requirement 11: Error Handling and User Feedback

**User Story**: As a user, I want clear error messages and feedback for all operations, so that I understand what went wrong and how to fix it.

#### Acceptance Criteria

1. WHEN an error occurs during vault operations THEN the system SHALL:
   - Display error message in a notification box
   - Provide actionable information (e.g., "Vault file not found. Run 'pass-cli init' to create one.")
   - Allow user to dismiss notification (Enter or Escape)
2. WHEN network or file system operation fails THEN the system SHALL:
   - Display error with technical details if verbose mode enabled
   - Display user-friendly message otherwise
   - Not crash or exit TUI
3. WHEN successful operation completes THEN the system SHALL:
   - Display brief success notification (e.g., "Credential saved")
   - Auto-dismiss notification after 2 seconds
   - Allow immediate dismissal with any key
4. WHEN clipboard operation fails THEN the system SHALL display "Clipboard unavailable" warning and continue operation

### Requirement 12: CLI and TUI Coexistence

**User Story**: As a power user, I want all existing CLI commands to work unchanged, so that my scripts and workflows are not disrupted by TUI addition.

#### Acceptance Criteria

1. WHEN user provides command arguments (e.g., `pass-cli list`, `pass-cli get github`) THEN the system SHALL execute CLI mode exactly as before TUI implementation
2. WHEN user runs `pass-cli --help` THEN the system SHALL display CLI help (not launch TUI)
3. WHEN user runs `pass-cli version` THEN the system SHALL display version info (not launch TUI)
4. WHEN scripts invoke `pass-cli` with arguments THEN behavior SHALL be identical to pre-TUI versions
5. WHEN TUI and CLI both access vault THEN they SHALL:
   - Use the same VaultService implementation
   - Use the same encryption/decryption logic
   - Use the same keychain integration
   - Maintain data format compatibility

## Non-Functional Requirements

### Code Architecture and Modularity

- **Single Responsibility Principle**: Each TUI component file (list view, detail view, forms) handles one specific screen or UI concern
- **Modular Design**: TUI package (`cmd/tui/`) is completely isolated from CLI commands; removing TUI package would not affect CLI functionality
- **Dependency Management**: TUI depends on `internal/vault/VaultService` but does not directly access storage, crypto, or keychain layers
- **Clear Interfaces**: TUI uses public VaultService API only; no access to internal implementation details

### Performance

- **Startup Time**: TUI mode SHALL launch within 100ms on modern hardware (comparable to CLI)
- **Search Responsiveness**: Credential list filtering SHALL update within 50ms of keystroke
- **Screen Redraw**: UI SHALL maintain 60fps during animations and transitions
- **Memory Usage**: TUI mode SHALL add no more than 10MB to process memory footprint beyond CLI mode
- **Large Vaults**: SHALL handle vaults with 1000+ credentials without performance degradation

### Security

- **Password Masking**: Passwords SHALL be masked (asterisks) by default in all views
- **Clipboard Security**: Clipboard operations SHALL use existing clipboard utility (same timeout behavior as CLI)
- **No Logging**: Sensitive data (passwords, master password) SHALL NOT be logged even in debug/verbose mode
- **Memory Clearing**: Password strings SHALL be cleared from memory when no longer needed (same as CLI)
- **Keychain Integration**: TUI SHALL use existing KeychainService without modifying security properties

### Reliability

- **Graceful Degradation**: TUI SHALL detect and handle:
  - Terminal size too small (display minimum size warning)
  - Missing terminal capabilities (fallback to basic colors)
  - Clipboard unavailable (show warning, continue operation)
- **Error Recovery**: Errors during operations (save, delete, update) SHALL NOT crash TUI; user returns to safe state
- **State Consistency**: Interrupt (Ctrl+C) SHALL cleanly exit without leaving vault in inconsistent state
- **Terminal Compatibility**: SHALL work on:
  - Windows Terminal, PowerShell, CMD
  - macOS Terminal, iTerm2
  - Linux GNOME Terminal, Konsole, Alacritty

### Usability

- **Keyboard-Only Operation**: 100% of TUI functionality SHALL be accessible via keyboard (no mouse required)
- **Discoverability**: Help overlay (`?` key) SHALL be accessible from any view
- **Consistent Navigation**: Same keys SHALL perform same actions across all views (Escape = back, Enter = select, etc.)
- **Visual Feedback**: All actions SHALL provide visual confirmation (highlights, notifications, animations)
- **Accessibility**: Text SHALL maintain 4.5:1 contrast ratio with background for readability

### Testing

- **Unit Test Coverage**: Each TUI component (views, components) SHALL have unit tests covering:
  - State transitions
  - Input handling
  - Rendering logic
- **Integration Tests**: End-to-end tests SHALL verify:
  - TUI launches and exits cleanly
  - Full CRUD workflow (add, view, update, delete)
  - Keychain integration in TUI mode
- **Manual Testing**: Cross-platform terminal testing on Windows, macOS, Linux before release
