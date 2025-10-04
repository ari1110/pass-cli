# Requirements Document

## Introduction

This specification defines the requirements for completing the incomplete tview migration that was started under the `tview-migration` spec. Through detailed code analysis, we discovered that while tview view and component implementations were created (e.g., `ListViewTview`, `SidebarTview`), they were never integrated into the active codebase. The application still runs on Bubble Tea (`useBubbleTea := true`) and the Model struct still instantiates old Bubble Tea views instead of the new tview ones.

**Purpose**: Complete the tview migration by integrating existing tview implementations, removing obsolete Bubble Tea code, and ensuring the application actually runs on tview instead of Bubble Tea.

**Value**: Achieves the original migration goals (reliable layout, no height overflow bugs) by actually activating the tview code that was written but never integrated. Eliminates duplicate codebases and technical debt from the incomplete migration.

## Alignment with Product Vision

This remediation directly supports:

- **Technical Excellence**: Completes a migration instead of leaving dual implementations that confuse maintainability
- **Developer Experience**: Removes confusion from having two parallel implementations of the same features
- **Code Quality**: Eliminates dead code and reduces codebase complexity
- **Original Migration Goals**: Actually achieves the layout reliability benefits that justified the tview migration

## Requirements

### Requirement 1: View Integration

**User Story:** As a developer, I want the Model struct to use tview view implementations, so that the application actually runs tview views instead of Bubble Tea views.

#### Acceptance Criteria

1. WHEN Model struct is examined THEN view fields SHALL be typed as `*views.ListViewTview`, `*views.DetailViewTview`, etc. (not old `*views.ListView` types)
2. WHEN NewModel() is called THEN the system SHALL instantiate `NewListViewTview()`, `NewDetailViewTview()`, etc. (not old Bubble Tea constructors)
3. WHEN the application runs THEN tview views SHALL render credentials, forms, and detail screens
4. WHEN view methods are called THEN the system SHALL use tview view methods (e.g., `GetTable()`, `GetForm()`)
5. IF the application compiles successfully THEN the system SHALL NOT reference old Bubble Tea view structs (`ListView`, `DetailView`, `AddFormView`, `EditFormView`, `HelpView`, `ConfirmView`)

### Requirement 2: Component Integration

**User Story:** As a developer, I want the Model struct to use tview component implementations, so that dashboard panels render with tview primitives.

#### Acceptance Criteria

1. WHEN Model struct is examined THEN component fields SHALL be typed as `*components.SidebarTview`, `*components.StatusBarTview`, etc.
2. WHEN NewModel() is called THEN the system SHALL instantiate `NewSidebarTview()`, `NewStatusBarTview()`, etc.
3. WHEN the dashboard renders THEN tview components SHALL display sidebar, status bar, metadata panel, breadcrumb, and command bar
4. WHEN component methods are called THEN the system SHALL use tview component methods
5. IF the application compiles successfully THEN the system SHALL NOT reference old Bubble Tea component structs (`SidebarPanel`, `StatusBar`, `MetadataPanel`, `CommandBar`)

### Requirement 3: Application Activation

**User Story:** As a developer, I want the application to run on tview by default, so that the tview migration is actually active instead of dormant.

#### Acceptance Criteria

1. WHEN tui.go Run() function is examined THEN `useBubbleTea` SHALL be set to `false`
2. WHEN the TUI launches THEN the system SHALL execute the `runTview()` code path (not `tea.NewProgram()`)
3. WHEN the TUI runs THEN the system SHALL use `tview.Application` as the event loop
4. WHEN keyboard events are received THEN the system SHALL use `tview.Application.SetInputCapture()` for event handling
5. IF the application runs successfully THEN Bubble Tea code paths SHALL NOT be executed

### Requirement 4: Old Code Deletion

**User Story:** As a developer, I want obsolete Bubble Tea view and component code removed, so that the codebase doesn't maintain duplicate implementations.

#### Acceptance Criteria

1. WHEN old view structs are deleted THEN files SHALL NOT contain `type ListView struct`, `type DetailView struct`, `type AddFormView struct`, `type EditFormView struct`, `type HelpView struct`, or `type ConfirmView struct`
2. WHEN old component structs are deleted THEN files SHALL NOT contain `type SidebarPanel struct`, `type StatusBar struct`, `type MetadataPanel struct`, or `type CommandBar struct`
3. WHEN struct deletion is complete THEN only `*Tview` suffix structs SHALL remain (e.g., `ListViewTview`, `SidebarTview`)
4. WHEN deletion is complete THEN old constructors (`NewListView()`, `NewSidebarPanel()`) SHALL NOT exist
5. IF files are cleaned up THEN the system SHALL compile without referencing deleted structs

### Requirement 5: Dependency Cleanup

**User Story:** As a developer, I want Bubble Tea dependencies removed from go.mod, so that the migration meets the original Requirement 1.2 from tview-migration spec.

#### Acceptance Criteria

1. WHEN go.mod is examined THEN the system SHALL NOT list `github.com/charmbracelet/bubbletea` as a dependency
2. WHEN go.mod is examined THEN the system SHALL NOT list `github.com/charmbracelet/bubbles` as a dependency
3. WHEN `go mod tidy` is run THEN the system SHALL NOT re-add Bubble Tea dependencies
4. WHEN the application is built THEN the system SHALL compile without Bubble Tea imports
5. IF lipgloss is used only for styling with tview THEN `github.com/charmbracelet/lipgloss` MAY remain (acceptable for tview compatibility)

### Requirement 6: Test Updates

**User Story:** As a developer, I want unit tests to verify tview implementations, so that test coverage validates the active code path.

#### Acceptance Criteria

1. WHEN view tests are examined THEN tests SHALL instantiate tview view constructors (`NewListViewTview()`, etc.)
2. WHEN component tests are examined THEN tests SHALL verify tview component behavior
3. WHEN tests are run THEN all unit tests SHALL pass
4. WHEN tests execute THEN the system SHALL test actual tview rendering and event handling
5. IF old Bubble Tea test code exists THEN it SHALL be updated or removed to test tview implementations

### Requirement 7: Functional Preservation

**User Story:** As a user, I want all existing TUI functionality to work identically after remediation, so that my workflow is not disrupted.

#### Acceptance Criteria

1. WHEN the TUI launches THEN the system SHALL display the credential list identically to before
2. WHEN I navigate with keyboard shortcuts THEN the system SHALL respond to all existing keybindings (h/j/k/l, arrows, /, :, etc.)
3. WHEN I add/edit/delete credentials THEN vault operations SHALL complete successfully
4. WHEN I search for credentials THEN filtering SHALL work in real-time
5. WHEN the dashboard renders THEN all panels (sidebar, main, metadata, status) SHALL display correctly
6. WHEN the terminal is resized THEN responsive breakpoints SHALL trigger at 80 and 120 columns as designed

## Non-Functional Requirements

### Code Architecture and Modularity
- **Clean Integration**: View and component integration must preserve existing file structure
- **No Duplication**: Old and new implementations must not coexist after integration
- **Compilation Safety**: Each integration step must compile successfully before proceeding
- **Atomic Changes**: Each file deletion should be committed separately to enable easy rollback if needed

### Performance
- **Startup Time**: Application startup must remain under 200ms as per original tview-migration spec
- **Render Performance**: Layout recalculation must complete under 50ms
- **Memory Usage**: Application memory usage must stay within 50MB during normal operation

### Reliability
- **No Regressions**: All existing functionality must continue to work
- **Error Handling**: tview error paths must be equivalent to Bubble Tea error paths
- **State Preservation**: Application state must be maintained identically during migration

### Testability
- **Test Coverage**: Unit test coverage must remain at current levels or improve
- **Integration Tests**: Existing CLI integration tests must continue to pass
- **Manual Testing**: Visual regression checklist must be executed to verify appearance
