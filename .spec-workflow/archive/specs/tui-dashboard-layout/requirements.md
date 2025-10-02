# Requirements Document: TUI Dashboard Layout

## Introduction

Pass-CLI's current TUI provides a functional, single-column interface for credential management. While effective, it doesn't maximize available screen real estate or provide the visual hierarchy and discoverability that modern terminal applications offer. This feature transforms the TUI into a **dashboard-style, panel-based interface** inspired by modern terminal file managers like Superfile, providing a file-explorer-like experience with better organization, navigation, and window sizing.

**Purpose**: Enhance the TUI with a professional, panel-based layout that improves information density, discoverability, and usability while maintaining backward compatibility and ensuring proper sizing across all terminal dimensions.

**Value to Users**:
- **Better Information Density**: See credential details, categories, and statistics simultaneously
- **Improved Navigation**: Hierarchical organization with breadcrumbs and tree views
- **Flexible Workspace**: Toggle panels on/off to optimize screen usage
- **Professional Aesthetics**: Modern, icon-enhanced design matching tools like Superfile
- **Proper Sizing**: Dynamic layout that adapts to terminal size without cutoff or overflow
- **Power User Features**: Command bar for vim-style commands and multi-panel workflows
- **Enhanced Productivity**: Quick access to stats, categories, and recent credentials

## Alignment with Product Vision

This feature directly supports principles outlined in product.md:

**Developer Experience (Principle #2)**: The dashboard interface reduces cognitive load by organizing credentials visually, showing context at a glance, and providing multiple navigation methods (sidebar, breadcrumbs, search, commands).

**Privacy by Design (Principle #3)**: All features remain local-only with no network dependencies. Metadata panel shows usage tracking without compromising security.

**Minimal Dependencies (Principle #5)**: Built entirely with existing Bubble Tea and Lipgloss libraries. No additional dependencies required.

**Product Positioning**: Enhances Pass-CLI's unique position by matching the quality and user experience of modern terminal applications like Superfile, setting a new standard for CLI password managers.

**Success Metrics Alignment**:
- **Usability**: Dashboard interface reduces learning curve and improves workflow efficiency
- **Quality**: Component-based architecture maintains 90%+ test coverage target
- **Adoption**: Professional, modern UI may contribute to GitHub star growth and community adoption

## Requirements

### Requirement 1: Multi-Panel Layout System

**User Story**: As a user, I want a panel-based interface with sidebar, main content, and metadata areas, so that I can see multiple types of information simultaneously without switching views.

#### Acceptance Criteria

1. WHEN TUI launches in dashboard mode THEN the system SHALL display three primary panel areas:
   - Left sidebar (20-25% width): Categories, stats, and quick actions
   - Center main content (50-60% width): Credential list, detail view, or forms
   - Right metadata panel (20-25% width): Selected credential details and usage info
2. WHEN terminal width is ≥120 columns THEN the system SHALL display all three panels simultaneously
3. WHEN terminal width is 80-119 columns THEN the system SHALL display sidebar + main content (metadata available on toggle)
4. WHEN terminal width is <80 columns THEN the system SHALL display main content only (sidebar and metadata available on toggle)
5. WHEN panels are displayed THEN each SHALL have visual borders distinguishing panel boundaries
6. WHEN a panel has focus THEN the system SHALL highlight its border with primary color (cyan)
7. WHEN a panel is inactive THEN the system SHALL display its border in subtle gray

### Requirement 2: Panel Toggle System

**User Story**: As a user, I want to show or hide individual panels using keyboard shortcuts, so that I can maximize space for the content I'm currently viewing.

#### Acceptance Criteria

1. WHEN user presses `s` key THEN the system SHALL toggle sidebar panel visibility
2. WHEN user presses `m` key THEN the system SHALL toggle metadata panel visibility
3. WHEN user presses `p` key THEN the system SHALL toggle process panel visibility
4. WHEN user presses `f` key THEN the system SHALL toggle all footer panels (processes, clipboard if implemented)
5. WHEN a panel is toggled off THEN the system SHALL redistribute its space to remaining visible panels
6. WHEN a panel is toggled on THEN the system SHALL recalculate layout and restore panel to appropriate size
7. WHEN panel visibility changes THEN the system SHALL smoothly redraw the interface without flicker
8. WHEN all panels are hidden THEN the system SHALL display main content in full-screen mode
9. WHEN user presses panel toggle key again THEN the system SHALL restore panel to previous state

### Requirement 3: Credential Categorization and Sidebar Navigation

**User Story**: As a user, I want credentials automatically organized into categories with a navigable sidebar, so that I can quickly browse credentials by type without searching.

#### Acceptance Criteria

1. WHEN TUI displays credentials THEN the system SHALL automatically categorize them based on service name patterns:
   - **APIs & Services**: Services containing "api", "key", or common API providers
   - **Cloud Infrastructure**: AWS, Azure, GCP, DigitalOcean, Heroku, etc.
   - **Databases**: PostgreSQL, MySQL, MongoDB, Redis, etc.
   - **Version Control**: GitHub, GitLab, Bitbucket, etc.
   - **Communication**: SendGrid, Twilio, Mailgun, Slack, etc.
   - **Payment Processing**: Stripe, Square, PayPal, etc.
   - **AI Services**: OpenAI, Anthropic, Cohere, etc.
   - **Uncategorized**: Credentials not matching patterns
2. WHEN sidebar is visible THEN the system SHALL display category tree with:
   - Category name with icon
   - Credential count in parentheses (e.g., "☁️ Cloud (3)")
   - Expand/collapse indicator (▶/▼)
3. WHEN user navigates category with arrow keys or vim keys (j/k) THEN the system SHALL move selection up or down
4. WHEN user presses Enter or `l` on collapsed category THEN the system SHALL expand category to show credentials
5. WHEN user presses Enter or `l` on expanded category THEN the system SHALL collapse category
6. WHEN user selects a credential in sidebar THEN the system SHALL display it in main content area
7. WHEN user presses `h` or backspace on expanded category THEN the system SHALL collapse it
8. WHEN sidebar displays "All Credentials" option THEN selecting it SHALL show all credentials regardless of category

### Requirement 4: Statistics and Quick Actions Panel

**User Story**: As a user, I want to see credential statistics and quick action buttons in the sidebar, so that I can understand vault status and perform common tasks without memorizing commands.

#### Acceptance Criteria

1. WHEN sidebar is visible THEN the system SHALL display statistics section showing:
   - Total credential count
   - Number of credentials used (with usage records)
   - Number of credentials modified recently (last 7 days)
2. WHEN sidebar is visible THEN the system SHALL display quick actions section with:
   - `[a]` Add new credential
   - `[:]` Open command bar
   - `[?]` Show help
3. WHEN statistics update (add/delete/use credential) THEN the system SHALL update counts in real-time
4. WHEN user presses hotkey shown in quick actions THEN the system SHALL execute corresponding action

### Requirement 5: Breadcrumb Navigation

**User Story**: As a user, I want to see my current location in a breadcrumb path, so that I always know my navigation context.

#### Acceptance Criteria

1. WHEN viewing credentials in main content THEN the system SHALL display breadcrumb path at top showing:
   - Navigation hierarchy (e.g., "Home > APIs > Cloud > aws-prod")
   - Current category or view name
2. WHEN viewing all credentials THEN breadcrumb SHALL show "Home > All Credentials"
3. WHEN viewing specific category THEN breadcrumb SHALL show "Home > [Category Name]"
4. WHEN viewing credential detail THEN breadcrumb SHALL show "Home > [Category] > [Service Name]"
5. WHEN breadcrumb path is too long for panel width THEN the system SHALL truncate middle segments with "..." (e.g., "Home > ... > aws-prod")
6. WHEN breadcrumb is displayed THEN it SHALL use distinctive styling (bold or colored) to stand out from content

### Requirement 6: Metadata Panel for Credential Details

**User Story**: As a user, I want to see credential details in a dedicated metadata panel, so that I can view information while keeping the main list visible.

#### Acceptance Criteria

1. WHEN credential is selected in main content THEN the system SHALL display in metadata panel:
   - Service name (header)
   - Username
   - Password (masked by default)
   - Created timestamp (relative format)
   - Updated timestamp (relative format)
   - Usage records with locations and counts
2. WHEN metadata panel is visible and credential has no username THEN the system SHALL display "(not set)" or similar placeholder
3. WHEN metadata panel shows usage records THEN it SHALL display:
   - Working directory path
   - Access count
   - Last accessed timestamp (relative)
4. WHEN password field is focused in metadata panel THEN user SHALL be able to press `m` to toggle mask
5. WHEN password is visible in metadata panel THEN user SHALL be able to press `c` to copy to clipboard
6. WHEN metadata panel displays long text (paths, notes) THEN it SHALL wrap text appropriately or provide scrolling
7. WHEN no credential is selected THEN metadata panel SHALL display helpful message (e.g., "Select a credential to view details")

### Requirement 7: Command Bar (Vim-Style Commands)

**User Story**: As a power user, I want a vim-style command bar for executing actions by typing commands, so that I can perform operations quickly without navigating menus.

#### Acceptance Criteria

1. WHEN user presses `:` key THEN the system SHALL open command bar at bottom of screen with `:` prompt
2. WHEN command bar is open THEN user SHALL be able to type commands and see text in input field
3. WHEN user presses Enter in command bar THEN the system SHALL execute the typed command
4. WHEN user presses Esc or Ctrl+C in command bar THEN the system SHALL close command bar without executing
5. WHEN user types `:add [service]` THEN the system SHALL open add form pre-filled with service name
6. WHEN user types `:search [query]` THEN the system SHALL filter credentials matching query
7. WHEN user types `:category [name]` THEN the system SHALL navigate to specified category
8. WHEN user types `:help` or `:h` THEN the system SHALL open help overlay
9. WHEN user types `:quit` or `:q` THEN the system SHALL exit TUI
10. WHEN user types invalid command THEN the system SHALL display error message "Unknown command: [command]"
11. WHEN command bar is open THEN it SHALL display command history with up/down arrow keys

### Requirement 8: Process Panel for Async Operations

**User Story**: As a user, I want to see feedback for background operations in a dedicated panel, so that I know when tasks are in progress or completed.

#### Acceptance Criteria

1. WHEN async operation starts (password generation, save, delete) THEN the system SHALL display process panel
2. WHEN process panel is visible THEN it SHALL show:
   - Operation name/description
   - Status indicator (⏳ in progress, ✓ completed, ✗ failed)
   - Timestamp
3. WHEN operation completes successfully THEN process SHALL show green checkmark and success message
4. WHEN operation fails THEN process SHALL show red X and error message
5. WHEN multiple operations are in progress THEN process panel SHALL show most recent 3-5 processes
6. WHEN all operations complete THEN process panel SHALL auto-hide after 3 seconds
7. WHEN user presses `p` key THEN the system SHALL toggle process panel visibility
8. WHEN process panel is displayed THEN it SHALL appear at bottom of screen above status bar

### Requirement 9: Multiple Credential Panels (Multi-Pane View)

**User Story**: As a user, I want to open multiple credential panels side-by-side, so that I can compare credentials or view different categories simultaneously.

#### Acceptance Criteria

1. WHEN user presses `n` key in main content area THEN the system SHALL create new credential panel
2. WHEN multiple credential panels exist THEN the system SHALL split main content area horizontally
3. WHEN 2 panels exist THEN each SHALL receive ~50% of main content width
4. WHEN 3 panels exist THEN each SHALL receive ~33% of main content width
5. WHEN user presses Tab or `L` (shift+l) THEN the system SHALL move focus to next credential panel
6. WHEN user presses Shift+Tab or `H` (shift+h) THEN the system SHALL move focus to previous credential panel
7. WHEN user presses `w` on focused credential panel THEN the system SHALL close that panel and redistribute space
8. WHEN only one credential panel remains and user presses `w` THEN the system SHALL not close panel (minimum one required)
9. WHEN multiple panels are visible THEN each SHALL have independent navigation and selection state
10. WHEN panel has focus THEN its border SHALL be highlighted with primary color

### Requirement 10: Responsive Layout System

**User Story**: As a user, I want the dashboard to adapt gracefully to different terminal sizes, so that I can use it on both large and small screens without content being cut off.

#### Acceptance Criteria

1. WHEN terminal width is ≥120 columns THEN the system SHALL display full three-panel layout (sidebar + main + metadata)
2. WHEN terminal width is 80-119 columns THEN the system SHALL display two-panel layout (sidebar + main) with metadata toggle-able
3. WHEN terminal width is <80 columns THEN the system SHALL display single-panel layout (main only) with sidebar/metadata toggle-able
4. WHEN terminal height is <20 rows THEN the system SHALL display warning overlay: "Terminal too small. Resize to at least 80x20"
5. WHEN terminal is resized THEN the system SHALL recalculate panel dimensions and redraw layout smoothly
6. WHEN panel content is too wide for allocated space THEN the system SHALL truncate with ellipsis ("...") or wrap text appropriately
7. WHEN calculating panel sizes THEN the system SHALL respect minimum widths:
   - Sidebar: minimum 20 columns
   - Main content: minimum 40 columns
   - Metadata: minimum 25 columns
8. WHEN available space is insufficient for all panels THEN the system SHALL hide lower priority panels (metadata first, then sidebar)

### Requirement 11: Layout Manager and Dimension Calculations

**User Story**: As a developer, I want a centralized layout manager that calculates panel dimensions, so that sizing logic is consistent and maintainable.

#### Acceptance Criteria

1. WHEN TUI initializes THEN the system SHALL create LayoutManager component to manage all dimension calculations
2. WHEN terminal size changes THEN LayoutManager SHALL calculate new dimensions for all visible panels based on:
   - Terminal width and height
   - Panel visibility states (which panels are shown/hidden)
   - Number of credential panels in main area
   - Minimum size constraints
3. WHEN LayoutManager calculates dimensions THEN it SHALL return structured layout with:
   - Sidebar dimensions (x, y, width, height)
   - Main content area dimensions (x, y, width, height)
   - Metadata panel dimensions (x, y, width, height)
   - Process panel dimensions (if visible)
4. WHEN dimension calculations complete THEN LayoutManager SHALL propagate sizes to all components via SetSize() methods
5. WHEN calculation results in invalid layout (e.g., negative width) THEN LayoutManager SHALL fallback to single-panel mode
6. WHEN LayoutManager detects terminal below minimum size THEN it SHALL return error state triggering warning overlay

### Requirement 12: Visual Enhancements and Icon Integration

**User Story**: As a user, I want the dashboard to use icons and modern styling, so that it has a professional, contemporary appearance matching tools like Superfile.

#### Acceptance Criteria

1. WHEN dashboard displays categories THEN each SHALL have appropriate icon:
   - ☁️ Cloud Infrastructure
   - 🔑 APIs & Services
   - 💾 Databases
   - 📦 Version Control
   - 📧 Communication
   - 💰 Payment Processing
   - 🤖 AI Services
   - 📁 Uncategorized
2. WHEN dashboard displays status indicators THEN it SHALL use:
   - 🔓 Keychain available
   - 🔒 Password mode
   - ⏳ Operation in progress
   - ✓ Operation succeeded
   - ✗ Operation failed
   - ▶ Collapsed category
   - ▼ Expanded category
3. WHEN panel has focus THEN its border SHALL use bold primary color (cyan)
4. WHEN panel is inactive THEN its border SHALL use subtle gray
5. WHEN displaying panel headers THEN the system SHALL use bold text with icons
6. WHEN displaying statistics THEN the system SHALL use appropriate icons:
   - 📊 Total count
   - ⚡ Recently used
   - 📝 Recently updated
7. WHEN icons are not supported by terminal font THEN the system SHALL gracefully fallback to text symbols (>, v, *, +, -)

### Requirement 13: Keyboard Navigation Between Panels

**User Story**: As a user, I want intuitive keyboard shortcuts to switch focus between panels, so that I can navigate the dashboard efficiently without using a mouse.

#### Acceptance Criteria

1. WHEN user presses Tab THEN the system SHALL move focus to next panel in order: sidebar → main → metadata → sidebar
2. WHEN user presses Shift+Tab THEN the system SHALL move focus to previous panel in reverse order
3. WHEN panel receives focus THEN its border SHALL change to primary color (cyan)
4. WHEN panel loses focus THEN its border SHALL change to subtle color (gray)
5. WHEN sidebar has focus THEN arrow keys/vim keys SHALL navigate categories and credentials in sidebar
6. WHEN main content has focus THEN arrow keys/vim keys SHALL navigate credentials in list or scroll detail view
7. WHEN metadata panel has focus THEN arrow keys SHALL scroll content if it exceeds panel height
8. WHEN command bar is open THEN Tab SHALL not switch panels (used for command completion instead)
9. WHEN help overlay is open THEN Tab SHALL not switch panels (overlay has full focus)

### Requirement 14: Category Tree Expand/Collapse

**User Story**: As a user, I want to expand and collapse categories in the sidebar tree, so that I can focus on relevant categories and reduce visual clutter.

#### Acceptance Criteria

1. WHEN sidebar displays categories THEN each category SHALL show expand/collapse indicator (▶ collapsed, ▼ expanded)
2. WHEN user presses Enter or `l` on collapsed category THEN the system SHALL expand category showing its credentials
3. WHEN user presses Enter or `l` on expanded category THEN the system SHALL collapse category hiding its credentials
4. WHEN category is expanded THEN credentials SHALL be indented beneath category name
5. WHEN user presses `h` or backspace on expanded category THEN the system SHALL collapse it
6. WHEN user presses `h` on credential item THEN the system SHALL collapse parent category
7. WHEN sidebar scrolls and category is collapsed THEN only category header SHALL be visible
8. WHEN all categories are collapsed THEN sidebar SHALL show compact list of category names with counts
9. WHEN TUI launches THEN all categories SHALL be collapsed by default (user expands as needed)

### Requirement 15: Integration with Existing TUI Features

**User Story**: As a user, I want all existing TUI functionality (add, edit, delete, search, help) to work seamlessly within the new dashboard layout, so that I don't lose any features.

#### Acceptance Criteria

1. WHEN user presses `a` to add credential THEN add form SHALL display in main content area (panels remain visible)
2. WHEN user presses `e` to edit credential THEN edit form SHALL display in main content area
3. WHEN user presses `/` to search THEN search bar SHALL appear in main content area
4. WHEN user presses `?` or F1 for help THEN help overlay SHALL appear over entire dashboard
5. WHEN user selects credential from sidebar OR main list THEN detail view SHALL display in main content area AND metadata panel (if visible)
6. WHEN user performs any CRUD operation THEN sidebar statistics SHALL update in real-time
7. WHEN credential is added/deleted/updated THEN category counts SHALL update in sidebar
8. WHEN user copies password (press `c`) THEN process panel SHALL show "Password copied to clipboard" notification
9. WHEN user generates password (Ctrl+G in form) THEN process panel SHALL show "Password generated" notification
10. WHEN all existing keyboard shortcuts are used THEN they SHALL behave identically to pre-dashboard TUI

## Non-Functional Requirements

### Code Architecture and Modularity

- **Single Responsibility Principle**: Each panel component (sidebar, metadata, process panel) handles only its own rendering and interaction logic
- **Modular Design**: LayoutManager is isolated from panel implementations; panels know nothing about layout calculations
- **Dependency Management**: Panels depend only on shared theme and standard Bubble Tea interfaces; no inter-panel dependencies
- **Clear Interfaces**: All panels implement standard Sizeable interface with SetSize(), Update(), View() methods
- **Component Reusability**: Existing views (ListView, DetailView, Forms) are reused as main content within panels
- **Backward Compatibility**: Dashboard features are additive; existing TUI can function if dashboard components are removed

### Performance

- **Startup Time**: Dashboard TUI SHALL launch within 150ms on modern hardware (50ms increase from current TUI acceptable for enhanced features)
- **Panel Switching**: Focus changes SHALL complete within 16ms (60fps) for smooth transitions
- **Layout Recalculation**: Terminal resize SHALL recalculate and redraw within 50ms
- **Category Rendering**: Sidebar with 100+ credentials across 8 categories SHALL render within 30ms
- **Memory Usage**: Dashboard mode SHALL add no more than 5MB to process memory beyond current TUI
- **Large Vaults**: SHALL handle vaults with 1000+ credentials without performance degradation in any panel

### Security

- **Password Masking**: Passwords SHALL remain masked by default in ALL panels (main, metadata, forms)
- **No Logging**: Panel state, credential content, and command bar history SHALL NOT be logged
- **Memory Clearing**: Sensitive data in panel buffers SHALL be cleared when panels are hidden or destroyed
- **Clipboard Security**: Clipboard operations SHALL use existing secure clipboard utility with timeout
- **Command Bar Security**: Command history SHALL NOT persist commands containing sensitive data
- **Process Panel Security**: Process notifications SHALL NOT display password values in operation descriptions

### Reliability

- **Graceful Degradation**: Dashboard SHALL detect and handle:
  - Terminal size too small (show warning, don't crash)
  - Missing nerd font support (fallback to ASCII symbols)
  - Panel layout calculation failures (fallback to single-panel mode)
- **Error Recovery**: Layout calculation errors SHALL NOT crash TUI; system SHALL fallback to safe minimal layout
- **State Consistency**: Panel visibility changes SHALL maintain correct focus state
- **Interrupt Handling**: Ctrl+C SHALL cleanly exit from any panel or command bar state
- **Terminal Compatibility**: Dashboard SHALL work correctly on:
  - Windows: Windows Terminal, PowerShell, CMD
  - macOS: Terminal.app, iTerm2, Alacritty
  - Linux: GNOME Terminal, Konsole, Kitty, Alacritty, st

### Usability

- **Keyboard-Only Operation**: 100% of dashboard functionality SHALL be accessible via keyboard
- **Discoverability**: Panel toggle keys SHALL be shown in status bar and help overlay
- **Consistent Navigation**: Same keys SHALL perform same actions across panels (j/k for up/down, h/l for collapse/expand)
- **Visual Feedback**: Panel focus changes SHALL be immediately visible via border color change
- **Accessibility**: All text SHALL maintain 4.5:1 contrast ratio for readability
- **Progressive Disclosure**: Complex features (command bar, multi-panel) are opt-in; basic workflow works without them
- **Responsive Behavior**: Layout changes SHALL be smooth and predictable when resizing terminal

### Testing

- **Unit Test Coverage**: Each panel component SHALL have unit tests covering:
  - Rendering with various data states (empty, single item, many items)
  - Keyboard input handling
  - Focus state transitions
  - Size calculation edge cases
- **Integration Tests**: End-to-end tests SHALL verify:
  - Panel layout across different terminal sizes
  - Panel toggle functionality
  - Multi-panel workflows
  - Command bar execution
  - Category navigation and expansion
- **Manual Testing**: Cross-platform testing on Windows, macOS, Linux before release
- **Regression Testing**: All existing TUI integration tests SHALL pass unchanged
