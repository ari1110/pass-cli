# Requirements Document

## Introduction

This specification addresses polish and completion items for the tview-based TUI dashboard implementation. The current implementation has core functionality working, but several features are incomplete, data model fields are missing, and UI elements need refinement for production readiness. This work will complete the transition from Bubble Tea to tview by resolving all TODOs, adding missing data model support, and improving visual consistency and usability.

The tview TUI is an interactive terminal dashboard for pass-cli that provides visual credential management with keyboard shortcuts, multi-panel layout, and real-time updates. Currently, the forms display Category and URL fields that don't exist in the data model, edit forms don't pre-populate values, and several UI elements lack visual polish.

## Alignment with Product Vision

**From product.md:**
- **Developer Experience**: "Design for speed, simplicity, and CLI integration" - The TUI must be polished and intuitive for daily use
- **Interactive TUI Dashboard** (In Progress): "Multi-panel terminal interface with sidebar navigation, metadata panel, and visual credential management"
- **Key Features**: The TUI provides an interactive alternative to CLI commands for developers who prefer visual interfaces

This specification directly supports:
1. **Completing the TUI feature** listed in "Features in Development"
2. **Developer productivity** by ensuring the TUI is production-ready
3. **Usability** through improved visual feedback and complete functionality

## Requirements

### Requirement 1: Data Model Completion

**User Story:** As a developer, I want to categorize my credentials and associate URLs with them, so that I can organize credentials by project/service type and quickly access related web interfaces.

#### Acceptance Criteria

1.1. WHEN I add a credential THEN the system SHALL accept optional `category`, `url`, and `notes` fields in addition to service, username, and password

1.2. WHEN I edit a credential THEN the system SHALL preserve existing `category`, `url`, and `notes` values if not modified

1.3. WHEN I view credential details THEN the system SHALL display the `category` and `url` if present

1.4. WHEN I list credentials in the TUI sidebar THEN the system SHALL group credentials by category

1.5. WHEN credentials have no category THEN the system SHALL display them under "Uncategorized"

### Requirement 2: Form Pre-Population

**User Story:** As a user, I want edit forms to show current credential values, so that I can see what I'm changing and avoid accidentally clearing fields.

#### Acceptance Criteria

2.1. WHEN I open the edit form for a credential THEN the system SHALL pre-populate the password field with the current password (masked by default)

2.2. WHEN I open the edit form for a credential THEN the system SHALL pre-populate the notes field with existing notes

2.3. WHEN I open the edit form for a credential THEN the system SHALL pre-populate the category dropdown with the current category

2.4. WHEN I open the edit form for a credential THEN the system SHALL pre-populate the URL field with the existing URL

2.5. WHEN I leave a field unchanged in the edit form THEN the system SHALL preserve the original value

### Requirement 3: Keyboard Shortcut Clarity

**User Story:** As a user, I want to see available keyboard shortcuts clearly displayed, so that I can discover and use features without memorizing commands.

#### Acceptance Criteria

3.1. WHEN I view the status bar THEN the system SHALL display keyboard shortcuts with clear key-action pairing (e.g., "[n] New" instead of just "New")

3.2. WHEN I view the Details panel THEN the system SHALL display action shortcuts prominently (not in small gray text)

3.3. WHEN I press "?" THEN the system SHALL display a well-formatted help modal with organized sections

3.4. WHEN the help modal is displayed THEN the system SHALL use proper alignment and spacing for readability

### Requirement 4: Visual Polish and Consistency

**User Story:** As a user, I want a visually consistent and polished interface, so that the TUI feels professional and is easy to read.

#### Acceptance Criteria

4.1. WHEN I view form input fields THEN the system SHALL ensure background colors fill the entire field area without gaps

4.2. WHEN I view the URL input field THEN the system SHALL prevent background color from extending beyond the field boundaries

4.3. WHEN I open add/edit forms THEN the system SHALL ensure Save and Cancel buttons are fully visible

4.4. WHEN I select a row in the credentials table THEN the system SHALL highlight it with clear visual indication

4.5. WHEN I select a credential in the table or sidebar THEN the Details panel SHALL immediately update to show the selected credential

### Requirement 5: Complete TODO Items

**User Story:** As a developer, I want all code TODO comments resolved, so that the codebase is complete and maintainable.

#### Acceptance Criteria

5.1. WHEN reviewing `cmd/tui-tview/components/forms.go:85` THEN the code SHALL pass Category/URL/Notes fields to AppState.AddCredential()

5.2. WHEN reviewing `cmd/tui-tview/components/forms.go:201` THEN the code SHALL pre-populate Category, URL, and Notes from the credential

5.3. WHEN reviewing `cmd/tui-tview/components/forms.go:291` THEN the code SHALL match the credential's category field in the dropdown

5.4. WHEN reviewing `cmd/tui-tview/models/state.go:364` THEN the code SHALL implement proper category extraction

### Requirement 6: Feature Verification

**User Story:** As a user, I want all advertised features to work correctly, so that I can rely on the TUI for credential management.

#### Acceptance Criteria

6.1. WHEN I press "/" THEN the system SHALL display a search/filter UI

**Note:** Requirements 6.2-6.4 are already implemented and verified working:
- ✅ 6.2: Delete confirmation modal displays when pressing 'd' (verified working)
- ✅ 6.3: Copy password shows status bar feedback when pressing 'c' (verified working)
- ✅ 6.4: Clipboard paste confirms password copied successfully (verified working)

These features require **verification testing only** (no implementation needed).

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Each component (forms, detail view, status bar) has a clear purpose
- **Modular Design**: Data model changes in vault.go should not require changes to multiple TUI components
- **Dependency Management**: TUI components depend on AppState, which manages vault interactions
- **Clear Interfaces**: AppState provides clean methods for credential operations

### Performance
- Form pre-population SHALL load credential data in <100ms
- Category grouping SHALL not degrade sidebar rendering performance
- Detail panel updates SHALL be instantaneous (<50ms)

### Security
- Passwords in edit forms SHALL remain masked by default
- Password visibility toggle SHALL be required to reveal passwords
- Category and URL fields SHALL accept any user input without validation (they are metadata, not security-critical)

### Reliability
- Edit forms SHALL never lose data if a field is left empty (preserve existing values)
- Missing Category/URL SHALL default to empty string, not cause errors
- All TODO items SHALL be resolved before marking spec complete

### Usability
- Keyboard shortcuts SHALL be discoverable through the help modal (?)
- Visual feedback SHALL confirm all user actions (copy, save, delete)
- Form layouts SHALL be consistent between add and edit modes
- Help modal SHALL be readable and well-organized
