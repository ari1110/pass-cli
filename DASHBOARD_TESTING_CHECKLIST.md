# TUI Dashboard Manual Testing Checklist

This document provides a comprehensive checklist for manually testing the TUI dashboard implementation (Task 21 of tui-dashboard-layout spec).

## Prerequisites

- Build the application: `go build -o pass-cli.exe .`
- Initialize a test vault with credentials from different categories:
  - Cloud: aws-production, azure-storage, gcp-project
  - Databases: postgres-main, mysql-dev, mongodb-cluster
  - Version Control: github-personal, gitlab-work
  - APIs: stripe-api, custom-api
  - AI Services: openai-api, anthropic-key
  - Payment: stripe-payments, paypal-business
  - Uncategorized: random-service

## Test Environment

- [ ] **Windows Terminal** (Primary)
- [ ] **PowerShell** (Secondary)
- [ ] **CMD** (Optional)
- [ ] **macOS Terminal.app** (If available)
- [ ] **Linux GNOME Terminal** (If available)

## Dashboard Layout Tests

### Responsive Layout Breakpoints

Test at different terminal sizes to verify responsive behavior:

- [ ] **Full Layout (140x40+)**
  - Sidebar visible on left
  - Main content in center
  - Metadata panel on right (when in detail view)
  - All panels render with proper borders
  - Status bar at bottom shows all information

- [ ] **Medium Layout (100x30 to 119x30)**
  - Sidebar visible on left
  - Main content in center
  - Metadata panel hidden
  - Status bar visible and functional

- [ ] **Small Layout (70x25 to 99x29)**
  - Sidebar hidden
  - Only main content visible
  - Status bar visible
  - Status bar shows "s: show sidebar" shortcut

- [ ] **Minimum Size Warning (<60x20)**
  - Application shows minimum size warning
  - Warning displays current size (e.g., "Current size: 50x15")
  - Warning displays required size (e.g., "Minimum size: 60x20")
  - Warning message: "Please resize your terminal"

### Panel Visibility Tests

- [ ] **Sidebar Panel**
  - Press `s` to toggle sidebar on/off
  - Sidebar displays all credential categories
  - Categories show credential counts
  - Categories have proper icons (☁️ 🔑 💾 📦 📧 💰 🤖 📁)
  - Stats section shows total/used/recent counts
  - Categories can be expanded/collapsed

- [ ] **Metadata Panel (Detail View Only)**
  - Press `m` to toggle metadata panel on/off
  - Panel shows selected credential details
  - Service name displayed as header
  - Username visible
  - Password masked by default with asterisks
  - Press `m` again to toggle password visibility
  - Created/updated timestamps shown in relative format
  - Notes section displayed if present
  - Usage records table shows access history

- [ ] **Process Panel**
  - Shows when async operations are running
  - Displays operation descriptions
  - Shows status icons (⏳ ✓ ✗)
  - Uses green for success, red for errors
  - Auto-hides after operations complete (3 second delay)
  - Maximum 5 processes visible at once

- [ ] **Command Bar**
  - Press `:` to open command bar
  - Input shows `:` prompt
  - Can type commands
  - Press `Esc` to cancel
  - Press `Enter` to execute
  - Errors displayed in red below input

- [ ] **Status Bar**
  - Always visible at bottom
  - Shows keychain indicator (🔓/🔒)
  - Shows credential count (e.g., "5 credentials")
  - Shows current view (List/Detail/Help)
  - Shows contextual shortcuts
  - Panel shortcuts update based on visibility

## Keyboard Shortcut Tests

### Panel Navigation

- [ ] **Tab** - Cycle through visible panels (Sidebar → Main → Metadata)
- [ ] **Shift+Tab** - Cycle through panels in reverse
- [ ] **s** - Toggle sidebar visibility
- [ ] **m** - Toggle metadata panel (Detail view) or password mask
- [ ] **p** - Toggle process panel (when processes active)
- [ ] **f** - Toggle all footer panels

### List View Shortcuts

- [ ] **/** - Activate search
- [ ] **a** - Add new credential
- [ ] **?** - Show help
- [ ] **q** - Quit application
- [ ] **:** - Open command bar
- [ ] **Up/Down** - Navigate credential list
- [ ] **Enter** - View credential details

### Detail View Shortcuts

- [ ] **Esc** - Return to list view
- [ ] **e** - Edit credential
- [ ] **d** - Delete credential (with confirmation)
- [ ] **c** - Copy password to clipboard
- [ ] **m** - Toggle password mask
- [ ] **Up/Down** - Scroll metadata panel (when focused)

### Command Bar Commands

Test each command:

- [ ] `:help` or `:h` - Opens help view
- [ ] `:quit` or `:q` - Quits application
- [ ] `:add <service>` - Opens add form (optionally pre-fills service)
- [ ] `:search <query>` - Activates search with query
- [ ] `:category <name>` - Navigates to category in sidebar
- [ ] Invalid command - Shows error message
- [ ] Empty command - Shows "Empty command" error
- [ ] Command without `:` - Shows "Commands must start with :" error

### Sidebar Interactions (When Focused)

- [ ] **Up/Down** or **j/k** - Navigate categories/credentials
- [ ] **Enter** - Select credential or expand/collapse category
- [ ] **l** - Expand category
- [ ] **h** - Collapse category
- [ ] Selecting credential loads details in main panel
- [ ] Breadcrumb updates with category path

## Visual Polish Tests

### Icons and Symbols

- [ ] **Category Icons** display correctly or fallback to ASCII
  - ☁️ Cloud / [CLD]
  - 🔑 APIs / [API]
  - 💾 Databases / [DB]
  - 📦 Version Control / [GIT]
  - 📧 Communication / [MSG]
  - 💰 Payment / [PAY]
  - 🤖 AI Services / [AI]
  - 📁 Uncategorized / [???]

- [ ] **Status Icons** display correctly
  - ⏳ Pending/Running / [.]
  - ✓ Success / [+]
  - ✗ Failed / [X]
  - ▶ Collapsed / [>]
  - ▼ Expanded / [v]

- [ ] **Keychain Indicators**
  - 🔓 Available / Keychain
  - 🔒 Password mode / Password

### Border Styling

- [ ] Focused panel has highlighted border (different color/style)
- [ ] Unfocused panels have subtle borders
- [ ] Borders render correctly (no broken box drawing characters)
- [ ] Panel borders don't overlap
- [ ] No visual artifacts or rendering glitches

### Text and Formatting

- [ ] Breadcrumb path displays with " > " separator
- [ ] Breadcrumb truncates middle segments with "..." when too long
- [ ] Long service names truncate with "..." in sidebar
- [ ] Timestamps show relative format ("2 hours ago", "3 days ago")
- [ ] Password masking uses asterisks consistently
- [ ] Text alignment is correct in all panels
- [ ] No text overflow beyond panel boundaries

### Colors and Theming

- [ ] Success messages in green
- [ ] Error messages in red
- [ ] Focused panel border distinct from unfocused
- [ ] Status bar uses subtle color
- [ ] Help text uses appropriate styling
- [ ] Color scheme is consistent throughout

## Performance Tests

- [ ] **Startup Time** - Application starts in < 150ms
- [ ] **View Transitions** - Switching views < 16ms (no perceived lag)
- [ ] **Panel Toggles** - Panel visibility changes are instant
- [ ] **Resize Handling** - Layout recalculates smoothly on window resize
- [ ] **Large Credential Lists** - Test with 50+ credentials
  - List view scrolls smoothly
  - Sidebar categories render quickly
  - Search remains responsive

## Categorization Tests

Verify credentials are automatically categorized correctly:

- [ ] **Cloud Services** - AWS, Azure, GCP, DigitalOcean, Heroku, etc.
- [ ] **Databases** - PostgreSQL, MySQL, MongoDB, Redis, etc.
- [ ] **Version Control** - GitHub, GitLab, Bitbucket
- [ ] **APIs & Services** - REST APIs, GraphQL, webhooks
- [ ] **AI Services** - OpenAI, Anthropic, Claude, GPT
- [ ] **Payment** - Stripe, PayPal, Square
- [ ] **Communication** - Slack, Discord, email services
- [ ] **Uncategorized** - Unknown services fall here
- [ ] **Case Insensitive** - "AWS-PROD" and "aws-prod" both categorize correctly
- [ ] **Stats Accurate** - Used/Recent counts match actual data

## Integration Tests

- [ ] **Existing Features Still Work**
  - Add credential workflow
  - Edit credential workflow
  - Delete credential with confirmation
  - Search functionality
  - Generate password
  - Copy to clipboard
  - Export functionality

- [ ] **No Regressions**
  - All pre-dashboard features functional
  - No crashes or panics
  - Error messages display correctly
  - Help system works

## Cross-Platform Considerations

### Windows-Specific

- [ ] Icons display correctly or ASCII fallback works
- [ ] Box drawing characters render properly
- [ ] Colors display correctly
- [ ] Keyboard shortcuts work (no conflicts)
- [ ] No CRLF/LF issues visible

### macOS-Specific (If Available)

- [ ] Terminal.app renders correctly
- [ ] iTerm2 renders correctly
- [ ] Keyboard shortcuts work (no conflicts with system shortcuts)
- [ ] Icons display correctly

### Linux-Specific (If Available)

- [ ] GNOME Terminal renders correctly
- [ ] Konsole renders correctly
- [ ] Alacritty renders correctly
- [ ] Keyboard shortcuts work

## Edge Cases

- [ ] **Empty Vault** - Dashboard handles zero credentials gracefully
- [ ] **Single Credential** - Layout works with minimal data
- [ ] **Very Long Service Names** - Truncation works correctly
- [ ] **Special Characters** - Service names with symbols render correctly
- [ ] **Unicode in Passwords** - Non-ASCII characters handled properly
- [ ] **Rapid Key Presses** - Application doesn't crash or lag
- [ ] **Window Resize During Operation** - Layout updates smoothly
- [ ] **Focus Lost/Regained** - Application state maintained

## Known Limitations

Document any issues found during testing:

1. **Issue**: _Description_
   - **Severity**: Critical/High/Medium/Low
   - **Steps to Reproduce**: _Steps_
   - **Expected**: _Expected behavior_
   - **Actual**: _Actual behavior_
   - **Workaround**: _If available_

## Test Coverage Summary

### Automated Test Coverage

**Unit Tests (56 tests)**:
- ✅ Layout Manager: 8 tests (all breakpoints, constraints, calculations)
- ✅ Category Tree: 13 tests (categorization, icons, edge cases)
- ✅ Sidebar Panel: 18 tests (navigation, selection, stats)
- ✅ Command Bar: 17 tests (parsing, history, errors)

**Integration Tests (8 tests)**:
- ✅ Dashboard initialization with real vault
- ✅ Credential categorization
- ✅ Responsive layout calculations
- ✅ Sidebar navigation
- ✅ Command bar parsing
- ✅ Status bar display
- ✅ Minimum size detection
- ✅ Category icon mappings

**Total**: 64 automated tests covering core functionality

### Manual Testing Required

- Panel visibility toggles and visual feedback
- Keyboard shortcut conflicts
- Cross-platform rendering
- Visual polish and UX feel
- Performance under real usage
- Icon display vs ASCII fallback

## Sign-Off

- [ ] All critical tests passing
- [ ] No critical bugs found
- [ ] Performance acceptable
- [ ] Visual polish complete
- [ ] Documentation complete

**Tested By**: _Name_
**Date**: _Date_
**Platform**: _OS and Terminal_
**Notes**: _Additional observations_
