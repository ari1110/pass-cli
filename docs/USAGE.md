# Usage Guide

Complete reference for all Pass-CLI commands, flags, and features.

## Table of Contents

- [Global Options](#global-options)
- [Commands](#commands)
  - [init](#init---initialize-vault)
  - [add](#add---add-credential)
  - [get](#get---retrieve-credential)
  - [list](#list---list-credentials)
  - [update](#update---update-credential)
  - [delete](#delete---delete-credential)
  - [generate](#generate---generate-password)
  - [version](#version---show-version)
- [Output Modes](#output-modes)
- [Script Integration](#script-integration)
- [Environment Variables](#environment-variables)
- [Configuration](#configuration)
- [Usage Tracking](#usage-tracking)
- [Best Practices](#best-practices)

## Global Options

Available for all commands:

| Flag | Description | Example |
|------|-------------|---------|
| `--vault <path>` | Custom vault location | `--vault /custom/path/vault.enc` |
| `--verbose` | Enable verbose output | `--verbose` |
| `--help`, `-h` | Show help | `--help` |

### Global Flag Examples

```bash
# Use custom vault location
pass-cli --vault /secure/vault.enc list

# Enable verbose logging
pass-cli --verbose get github

# Get help for any command
pass-cli get --help
```

## Commands

### init - Initialize Vault

Create a new password vault.

#### Synopsis

```bash
pass-cli init
```

#### Description

Creates a new encrypted vault at `~/.pass-cli/vault.enc` and stores the master password in your system keychain. You will be prompted to create a master password.

#### Examples

```bash
# Initialize with default location
pass-cli init

# Initialize with custom location
pass-cli --vault /custom/path/vault.enc init
```

#### Flags

| Flag | Type | Description |
|------|------|-------------|
| `--enable-audit` | bool | Enable tamper-evident audit logging |
| `--use-keychain` | bool | Store master password in OS keychain (default: true) |

#### Password Policy (January 2025)

All master passwords must meet complexity requirements:
- **Minimum Length**: 12 characters
- **Uppercase**: At least one uppercase letter (A-Z)
- **Lowercase**: At least one lowercase letter (a-z)
- **Digit**: At least one digit (0-9)
- **Symbol**: At least one special symbol (!@#$%^&*()-_=+[]{}|;:,.<>?)

**Examples**:
- ✅ `MySecureP@ssw0rd2025!` (meets all requirements)
- ✅ `Correct-Horse-Battery-29!` (meets all requirements)
- ❌ `password123` (too short, no uppercase, no symbol)
- ❌ `MyPassword` (no digit, no symbol)

#### Audit Logging (Optional)

Enable audit logging to record vault operations with HMAC signatures:

```bash
# Initialize vault with audit logging
pass-cli init --enable-audit
```

**Audit features**:
- **Tamper-Evident**: HMAC-SHA256 signatures prevent log modification
- **Privacy**: Service names logged, passwords NEVER logged
- **Key Storage**: HMAC keys stored in OS keychain (separate from vault)
- **Auto-Rotation**: Logs rotate at 10MB with 7-day retention
- **Graceful Degradation**: Operations continue if logging fails

**Verification**:
```bash
# Verify audit log integrity
pass-cli verify-audit
```

#### Notes

- Master password must meet complexity requirements (12+ chars, uppercase, lowercase, digit, symbol)
- Strong passwords (20+ characters) recommended for master password
- Master password is stored in OS keychain for convenience
- Vault file is created with restricted permissions (0600)
- Audit logging is opt-in (disabled by default)

---

### add - Add Credential

Add a new credential to the vault.

#### Synopsis

```bash
pass-cli add <service> [flags]
```

#### Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--username` | `-u` | string | Username for the credential |
| `--password` | `-p` | string | Password (not recommended, use prompt) |
| `--category` | `-c` | string | Category for organizing credentials (e.g., 'Cloud', 'Databases') |
| `--url` | | string | Service URL |
| `--notes` | | string | Additional notes |

#### Examples

```bash
# Interactive mode (prompts for username/password)
pass-cli add github

# With username flag
pass-cli add github --username user@example.com

# With URL and notes
pass-cli add github \
  --username user@example.com \
  --url https://github.com \
  --notes "Personal account"

# With category
pass-cli add github -u user@example.com -c "Version Control"

# All flags (not recommended for password)
pass-cli add github \
  -u user@example.com \
  -p secret123 \
  --url https://github.com \
  --notes "Work account"
```

#### Interactive Prompts

When not using flags, you'll be prompted:

```
Enter username: user@example.com
Enter password: ******* (hidden input)
Enter URL (optional): https://github.com
Enter notes (optional): Personal account
```

#### Password Policy

Credential passwords must meet the same complexity requirements as master passwords:
- Minimum 12 characters with uppercase, lowercase, digit, and symbol
- TUI mode shows real-time strength indicator
- Generated passwords automatically meet policy requirements

#### Notes

- Service names must be unique
- Password input is hidden by default
- Passing password via `-p` flag is insecure (visible in shell history)
- Use `pass-cli generate` to create strong random passwords that meet policy requirements
- Usage tracking begins when credential is first accessed

---

### get - Retrieve Credential

Retrieve a credential from the vault.

#### Synopsis

```bash
pass-cli get <service> [flags]
```

#### Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--quiet` | `-q` | bool | Output password only (for scripts) |
| `--field` | `-f` | string | Extract specific field |
| `--no-clipboard` | | bool | Skip clipboard copy |
| `--masked` | | bool | Display password as asterisks |

#### Field Options

For `--field` flag:
- `username` - User's username
- `password` - User's password
- `url` - Service URL
- `notes` - Additional notes
- `service` - Service name
- `created` - Creation timestamp
- `modified` - Last modified timestamp
- `accessed` - Last accessed timestamp

#### Examples

```bash
# Default: Display credential and copy to clipboard
pass-cli get github

# Quiet mode (password only, for scripts)
pass-cli get github --quiet
pass-cli get github -q

# Get specific field
pass-cli get github --field username
pass-cli get github -f url

# Quiet mode with specific field
pass-cli get github --field username --quiet

# Display without clipboard
pass-cli get github --no-clipboard

# Display with masked password
pass-cli get github --masked
```

#### Output Examples

**Default output:**
```
Service:  github
Username: user@example.com
Password: mySecretPassword123!
URL:      https://github.com
Notes:    Personal account

✓ Password copied to clipboard (will clear in 30 seconds)
```

**Quiet mode:**
```bash
$ pass-cli get github --quiet
mySecretPassword123!
```

**Field extraction:**
```bash
$ pass-cli get github --field username --quiet
user@example.com
```

#### Notes

- Clipboard auto-clears after 30 seconds
- Usage tracking records current directory
- Accessing a credential updates the "last accessed" timestamp

---

### list - List Credentials

List all credentials in the vault.

#### Synopsis

```bash
pass-cli list [flags]
```

#### Flags

| Flag | Type | Description |
|------|------|-------------|
| `--format` | string | Output format: table, json, simple (default: table) |
| `--unused` | bool | Show only unused credentials |
| `--days` | int | Days threshold for unused (default: 30) |

#### Examples

```bash
# List all credentials (table format)
pass-cli list

# List as JSON
pass-cli list --format json

# Simple list (service names only)
pass-cli list --format simple

# Show unused credentials (not accessed in 30 days)
pass-cli list --unused

# Show credentials not used in 90 days
pass-cli list --unused --days 90
```

#### Output Examples

**Table format (default):**
```
+----------+----------------------+---------------------+
| SERVICE  | USERNAME             | LAST ACCESSED       |
+----------+----------------------+---------------------+
| github   | user@example.com     | 2025-01-20 14:22    |
| aws-prod | admin@company.com    | 2025-01-18 09:15    |
| database | dbuser               | 2025-01-15 16:30    |
+----------+----------------------+---------------------+
```

**Simple format:**
```
github
aws-prod
database
```

#### Notes

- Passwords are never shown in list output
- Table format is best for human viewing
- JSON format is best for parsing
- Simple format is best for shell scripts

---

### update - Update Credential

Update an existing credential.

#### Synopsis

```bash
pass-cli update <service> [flags]
```

#### Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--username` | `-u` | string | New username |
| `--password` | `-p` | string | New password (not recommended) |
| `--category` | | string | New category |
| `--url` | | string | New URL |
| `--notes` | | string | New notes |
| `--clear-category` | | bool | Clear category field to empty |
| `--clear-notes` | | bool | Clear notes field to empty |
| `--clear-url` | | bool | Clear URL field to empty |
| `--force` | `-f` | bool | Skip confirmation prompt |

#### Examples

```bash
# Update password (prompted)
pass-cli update github

# Update username
pass-cli update github --username newuser@example.com

# Update URL
pass-cli update github --url https://github.com/enterprise

# Update notes
pass-cli update github --notes "Updated account info"

# Update category
pass-cli update github --category "Work"

# Clear category field
pass-cli update github --clear-category

# Update multiple fields
pass-cli update github \
  --username newuser@example.com \
  --url https://github.com/enterprise \
  --notes "Corporate account"
```

#### Interactive Mode

If no flags provided, prompts for password:

```
Enter new password (leave blank to keep current): *******
Password updated successfully!
```

#### Notes

- At least one field must be updated
- Updating password clears usage history
- Original values preserved if not specified

---

### delete - Delete Credential

Delete a credential from the vault.

#### Synopsis

```bash
pass-cli delete <service> [flags]
```

#### Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--force` | `-f` | bool | Skip confirmation prompt |

#### Examples

```bash
# Delete with confirmation
pass-cli delete github

# Force delete (no confirmation)
pass-cli delete github --force
pass-cli delete github -f
```

#### Interactive Confirmation

Without `--force`:

```
Are you sure you want to delete 'github'? (yes/no): yes
Credential 'github' deleted successfully!
```

#### Notes

- Deletion is permanent (no undo)
- Confirmation required unless using `--force`
- Credential completely removed from vault

---

### generate - Generate Password

Generate a cryptographically secure password.

#### Synopsis

```bash
pass-cli generate [flags]
```

#### Aliases

`gen`, `pwd`

#### Flags

| Flag | Type | Description |
|------|------|-------------|
| `--length` | int | Password length (8-128, default: 20) |
| `--no-lower` | bool | Exclude lowercase letters |
| `--no-upper` | bool | Exclude uppercase letters |
| `--no-digits` | bool | Exclude digits |
| `--no-symbols` | bool | Exclude symbols |
| `--no-clipboard` | bool | Skip clipboard copy |

#### Examples

```bash
# Generate default password (20 chars, all character types)
pass-cli generate

# Custom length
pass-cli generate --length 32

# Alphanumeric only (no symbols)
pass-cli generate --no-symbols

# Digits and symbols only
pass-cli generate --no-lower --no-upper

# Letters only (no digits or symbols)
pass-cli generate --no-digits --no-symbols

# Display only (no clipboard)
pass-cli generate --no-clipboard
```

#### Character Sets

Default character sets:
- Lowercase: `a-z`
- Uppercase: `A-Z`
- Digits: `0-9`
- Symbols: `!@#$%^&*()_+-=[]{}|;:,.<>?`

#### Notes

- Uses `crypto/rand` for cryptographic randomness
- At least one character set must be enabled
- Minimum length: 8 characters
- Maximum length: 128 characters
- Clipboard auto-clears after 30 seconds

---

### version - Show Version

Display version information.

#### Synopsis

```bash
pass-cli version [flags]
```

#### Flags

| Flag | Type | Description |
|------|------|-------------|
| `--verbose` | bool | Show detailed version info |

#### Examples

```bash
# Show version
pass-cli version

# Verbose version info
pass-cli version --verbose
```

#### Output Examples

**Default:**
```
pass-cli version X.Y.Z
```

**Verbose:**
```
pass-cli version X.Y.Z
  commit: abc123f
  built:  2025-01-20T10:30:00Z
  go:     go1.25.1
```

---

## Output Modes

Pass-CLI supports multiple output modes for different use cases.

### Human-Readable (Default)

Formatted tables and colored output for terminal viewing.

```bash
pass-cli get github
# Service:  github
# Username: user@example.com
# Password: ****** (or full password)
```

### Quiet Mode

Single-line output, perfect for scripts.

```bash
pass-cli get github --quiet
# mySecretPassword123!

pass-cli get github --field username --quiet
# user@example.com
```

### Simple Mode (List Only)

Service names only, one per line.

```bash
pass-cli list --format simple
# github
# aws-prod
# database
```

## Script Integration

### Bash/Zsh Examples

**Export to environment variable:**

```bash
#!/bin/bash

# Export password
export DB_PASSWORD=$(pass-cli get database --quiet)

# Export specific field
export DB_USER=$(pass-cli get database --field username --quiet)

# Use in command
mysql -u "$(pass-cli get database -f username -q)" \
      -p"$(pass-cli get database -q)" \
      mydb
```

**Conditional execution:**

```bash
# Check if credential exists
if pass-cli get myservice --quiet &>/dev/null; then
    echo "Credential exists"
    export API_KEY=$(pass-cli get myservice --quiet)
else
    echo "Credential not found"
    exit 1
fi
```

**Loop through credentials:**

```bash
# Process all credentials
for service in $(pass-cli list --format simple); do
    echo "Processing $service..."
    username=$(pass-cli get "$service" --field username --quiet)
    echo "  Username: $username"
done
```

### PowerShell Examples

**Export to environment variable:**

```powershell
# Export password
$env:DB_PASSWORD = pass-cli get database --quiet

# Export specific field
$env:DB_USER = pass-cli get database --field username --quiet

# Use in command
$apiKey = pass-cli get openai --quiet
Invoke-RestMethod -Uri "https://api.openai.com" -Headers @{
    "Authorization" = "Bearer $apiKey"
}
```

**Conditional execution:**

```powershell
# Check if credential exists
try {
    $password = pass-cli get myservice --quiet 2>$null
    Write-Host "Credential exists"
    $env:API_KEY = $password
} catch {
    Write-Host "Credential not found"
    exit 1
}
```

### Python Examples

```python
import subprocess

# Get password only
result = subprocess.run(
    ['pass-cli', 'get', 'github', '--quiet'],
    capture_output=True,
    text=True,
    check=True
)
password = result.stdout.strip()

# Get specific field
result = subprocess.run(
    ['pass-cli', 'get', 'github', '--field', 'username', '--quiet'],
    capture_output=True,
    text=True,
    check=True
)
username = result.stdout.strip()
```

### Makefile Examples

```makefile
.PHONY: deploy
deploy:
	@export AWS_KEY=$$(pass-cli get aws --quiet --field username); \
	export AWS_SECRET=$$(pass-cli get aws --quiet); \
	./deploy.sh

.PHONY: test-db
test-db:
	@DB_URL="postgres://$$(pass-cli get testdb -f username -q):$$(pass-cli get testdb -q)@localhost/testdb" \
	go test ./...
```

## Environment Variables

### PASS_CLI_VAULT

Override default vault location.

```bash
# Bash
export PASS_CLI_VAULT=/custom/path/vault.enc
pass-cli list

# PowerShell
$env:PASS_CLI_VAULT = "C:\custom\path\vault.enc"
pass-cli list
```

### PASS_CLI_VERBOSE

Enable verbose logging.

```bash
export PASS_CLI_VERBOSE=1
pass-cli get github
```

## Configuration

**Configuration Location** (added January 2025):
- **Linux/macOS**: `~/.config/pass-cli/config.yml`
- **Windows**: `%APPDATA%\pass-cli\config.yml`

**Management Commands**:
```bash
# Initialize default config
pass-cli config init

# Edit config in default editor
pass-cli config edit

# Validate config syntax
pass-cli config validate

# Reset to defaults
pass-cli config reset
```

### Example Configuration

```yaml
# Default vault location
vault: ~/.pass-cli/vault.enc

# Enable verbose output
verbose: false

# Clipboard timeout (seconds)
clipboard_timeout: 30

# Default password generation length
password_length: 20

# Terminal display thresholds (TUI mode)
terminal:
  min_width: 60   # Minimum columns (default: 60)
  min_height: 30  # Minimum rows (default: 30)

# Custom keyboard shortcuts (TUI mode)
keybindings:
  quit: "q"                  # Quit application
  add_credential: "a"        # Add new credential
  edit_credential: "e"       # Edit credential
  delete_credential: "d"     # Delete credential
  toggle_detail: "i"         # Toggle detail panel
  toggle_sidebar: "s"        # Toggle sidebar
  help: "?"                  # Show help modal
  search: "/"                # Activate search

# Supported key formats for keybindings:
# - Single letters: a-z
# - Numbers: 0-9
# - Function keys: f1-f12
# - Modifiers: ctrl+, alt+, shift+
# Examples: ctrl+q, alt+a, shift+f1
```

### Keybinding Customization

**Configurable Actions**:
- `quit`, `add_credential`, `edit_credential`, `delete_credential`
- `toggle_detail`, `toggle_sidebar`, `help`, `search`

**Hardcoded Shortcuts** (cannot be changed):
- Navigation: Tab, Shift+Tab, ↑/↓, Enter, Esc
- Forms: Ctrl+H, Ctrl+S, Ctrl+C
- Detail view: p, c

**Validation**:
- Duplicate key assignments rejected (conflict detection)
- Unknown actions rejected
- Invalid config shows warning modal, app continues with defaults
- UI hints automatically update to reflect custom keybindings

### Configuration Priority

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file
4. Built-in defaults (lowest priority)

## TUI Mode

Pass-CLI includes an interactive Terminal User Interface (TUI) for visual credential management. The TUI provides an alternative to CLI commands with visual navigation, real-time search, and keyboard-driven workflows.

### Launching TUI Mode

```bash
# Launch TUI (no arguments)
pass-cli

# TUI opens automatically when no subcommand is provided
```

The TUI launches immediately and displays:
- **Left sidebar**: Category navigation (auto-hides on narrow terminals)
- **Center table**: Credential list with service name, username, last accessed time
- **Right panel**: Credential details with password, URL, notes, usage locations
- **Bottom status bar**: Context-aware keyboard shortcuts and status messages

### TUI vs CLI Mode

Pass-CLI operates in two modes:

| Mode | Activation | Use Case |
|------|------------|----------|
| **TUI Mode** | Run `pass-cli` with no arguments | Interactive browsing, visual credential management |
| **CLI Mode** | Run `pass-cli <command>` with explicit subcommand | Scripts, automation, quick single operations |

**Examples**:
```bash
# TUI Mode
pass-cli                        # Opens interactive interface

# CLI Mode
pass-cli list                   # Outputs credential table to stdout
pass-cli get github --quiet     # Outputs password only (script-friendly)
pass-cli add newcred            # Interactive prompts for credential data
```

Both modes access the same encrypted vault file (`~/.pass-cli/vault.enc`).

### TUI Keyboard Shortcuts

#### Navigation

| Shortcut | Action | Context |
|----------|--------|---------|
| `Tab` | Next component | Any view |
| `Shift+Tab` | Previous component | Any view |
| `↑` / `↓` | Navigate lists | Table, sidebar |
| `Enter` | Select credential / View details | Table |

#### Actions

| Shortcut | Action | Context |
|----------|--------|---------|
| `n` | New credential (opens add form) | Main view |
| `e` | Edit selected credential | Main view (credential selected) |
| `d` | Delete selected credential | Main view (credential selected) |
| `p` | Toggle password visibility | Detail panel |
| `c` | Copy password to clipboard | Detail panel |

#### View Controls

| Shortcut | Action | Context |
|----------|--------|---------|
| `i` | Toggle detail panel (Auto/Hide/Show) | Main view |
| `s` | Toggle sidebar (Auto/Hide/Show) | Main view |
| `/` | Activate search mode | Main view |

#### Forms (Add/Edit)

| Shortcut | Action | Context |
|----------|--------|---------|
| `Ctrl+S` | Save form | Add/edit forms |
| `Ctrl+H` | Toggle password visibility | Add/edit forms |
| `Tab` | Next field | Forms |
| `Shift+Tab` | Previous field | Forms |
| `Esc` | Cancel / Close form | Forms |

#### General

| Shortcut | Action | Context |
|----------|--------|---------|
| `?` | Show help modal | Any time |
| `q` | Quit application | Main view |
| `Esc` | Close modal / Cancel search | Modals, search mode |
| `Ctrl+C` | Quit application | Any time |

**Note**: Configurable shortcuts (a, e, d, i, s, ?, /, q) can be customized via `~/.config/pass-cli/config.yml`. See [Configuration](#configuration) section for keybinding customization details. Navigation shortcuts (Tab, arrows, Enter, Esc, Ctrl+H, Ctrl+S, Ctrl+C) are hardcoded and cannot be changed.

### Search & Filter

Press `/` to activate search mode. An input field appears at the top of the credential table.

**Search Behavior**:
- **Case-insensitive**: "git" matches "GitHub", "gitlab", "digit"
- **Substring matching**: Query can appear anywhere in field
- **Searchable fields**: Service name, username, URL, category (Notes field excluded)
- **Real-time filtering**: Results update as you type
- **Navigation**: Use `↑`/`↓` arrow keys to navigate filtered results

**Examples**:
```bash
# Search for GitHub credentials
/
github      # Type query → only GitHub credentials shown

# Search by category
/
dev         # Shows credentials in "Development" category

# Clear search
Esc         # Exits search mode, shows all credentials
```

**When searching**:
- Newly added credentials matching the query appear immediately in results
- Selection preserved if selected credential matches search
- Empty results show message: "No credentials match your search"

### Password Visibility Toggle

In add and edit forms, press `Ctrl+H` to toggle between masked and visible passwords.

**Use Cases**:
- Verify password spelling before saving
- Check for typos when editing existing passwords
- Confirm generated passwords meet requirements

**Behavior**:
- **Default state**: Password masked (asterisks: `******`)
- **After `Ctrl+H`**: Password visible (plaintext), label shows `[VISIBLE]`
- **After `Ctrl+H` again**: Password masked again
- **On form close**: Visibility resets to masked (secure default)
- **Cursor position**: Preserved when toggling (no text loss)

**Examples**:
```bash
# In add form
n                              # Open new credential form
Type: SecureP@ssw0rd!         # Password shows as ******
Ctrl+H                         # Password shows: SecureP@ssw0rd!
Ctrl+H                         # Password shows as ******
Ctrl+S                         # Save (password saved correctly)

# In edit form
e                              # Open edit form for selected credential
Focus password field           # Existing password loads (masked)
Ctrl+H                         # View current password
Type new password              # Update password
Ctrl+H                         # Mask again to verify asterisks
Ctrl+S                         # Save changes
```

**Security Note**: Password visibility is per-form. Switching between add and edit forms resets visibility to masked.

### Layout Controls

The TUI layout adapts to terminal size with manual override controls.

#### Detail Panel Toggle (`i` key)

Cycles through three states:
1. **Auto (responsive)**: Shows on wide terminals (>100 cols), hides on narrow
2. **Force Hide**: Always hidden regardless of terminal width
3. **Force Show**: Always visible regardless of terminal width

Status bar displays current state when toggling:
- "Detail Panel: Auto (responsive)"
- "Detail Panel: Hidden"
- "Detail Panel: Visible"

**Use Cases**:
- Hide detail panel to focus on credential list
- Force show on narrow terminal to view credential details
- Return to auto mode for responsive behavior

#### Sidebar Toggle (`s` key)

Cycles through three states:
1. **Auto (responsive)**: Shows on wide terminals (>80 cols), hides on narrow
2. **Force Hide**: Always hidden regardless of terminal width
3. **Force Show**: Always visible regardless of terminal width

Status bar displays current state when toggling:
- "Sidebar: Auto (responsive)"
- "Sidebar: Hidden"
- "Sidebar: Visible"

**Use Cases**:
- Hide sidebar to maximize table width
- Force show on narrow terminal to access category navigation
- Return to auto mode for responsive behavior

**Manual overrides persist** until user changes them or application restarts.

### Usage Location Display

The detail panel shows where each credential has been accessed.

**Information Displayed**:
- **File path**: Absolute path to working directory where `pass-cli get` was executed
- **Access count**: Number of times credential accessed from that location
- **Timestamp**: Hybrid format (relative for recent, absolute for old)
  - Recent (within 7 days): "2 hours ago", "3 days ago"
  - Older: "2025-09-15", "2024-12-01"
- **Git repository** (if available): Repository name extracted from working directory
- **Line number** (if available): File path with line number (e.g., `/path/file.go:42`)

**Display Format**:
```
Usage Locations:
  /home/user/projects/web-app
    Accessed: 12 times
    Last: 2 hours ago
    Repo: web-app

  /home/user/projects/api-server/src/config.go:156
    Accessed: 5 times
    Last: 2025-09-20
    Repo: api-server
```

**Empty State**: If credential has never been accessed, shows: "No usage recorded"

**Sorting**: Locations sorted by most recent access timestamp descending.

**Use Cases**:
- Audit which projects use which credentials
- Identify stale credentials not accessed recently
- Track credential usage patterns across repositories
- Understand credential dependencies for project cleanup

### Exiting TUI Mode

Press `q` or `Ctrl+C` at any time to quit the TUI and return to shell.

**Note**: If a modal is open (add form, edit form, help), pressing `q` or `Esc` closes the modal instead of quitting. Press `q` again from main view to quit application.

## TUI Best Practices

1. **Use `/` search for large vaults** - Faster than scrolling through 50+ credentials
2. **Press `?` to see all shortcuts** - Built-in help always available
3. **Toggle detail panel (`i`) on narrow terminals** - Maximize table visibility
4. **Use `Ctrl+H` to verify passwords** - Catch typos before saving
5. **Check usage locations before deleting** - Understand credential dependencies
6. **Press `c` to copy passwords** - Clipboard auto-clears after 30 seconds

## TUI Troubleshooting

**Problem**: TUI doesn't launch, shows "command not found"
**Solution**: Ensure you're running `pass-cli` with no arguments. If you pass any argument (even invalid), it attempts CLI mode.

**Problem**: Ctrl+H does nothing in forms
**Solution**: Ensure you're in add or edit form, not the main view. Password toggle only works in forms.

**Problem**: Search key `/` types "/" character instead of activating search
**Solution**: Ensure focus is on the main view (table/sidebar), not inside a form or modal. Press `Esc` to close any open modal first.

**Problem**: Sidebar doesn't appear
**Solution**: Press `s` to toggle sidebar. On narrow terminals (<80 cols), sidebar auto-hides in responsive mode. Press `s` twice to force show.

**Problem**: Usage locations not showing
**Solution**: Usage locations only appear after you've accessed credentials via `pass-cli get <service>` from different working directories. New credentials won't have usage data until first access.

## Usage Tracking

Pass-CLI automatically tracks where credentials are accessed based on your current working directory.

### How It Works

```bash
# Access from project directory
cd ~/projects/my-app
pass-cli get database

# Usage tracking is automatic based on current directory
```

### Use Cases

- **Audit**: See which projects use which credentials
- **Cleanup**: Identify unused credentials
- **Documentation**: Auto-document credential dependencies

### Viewing Usage

```bash
# List unused credentials
pass-cli list --unused --days 30
```

## Best Practices

### Security

1. **Never pass passwords via flags** - Use prompts or `--generate`
2. **Use quiet mode in scripts** - Prevents logging sensitive data
3. **Clear shell history** - When testing commands with passwords
4. **Use strong master passwords** - 20+ characters recommended

### Workflow

1. **Generate passwords** - Use `--generate` for new credentials
2. **Update regularly** - Rotate credentials periodically
3. **Track usage** - Review unused credentials monthly
4. **Backup vault** - Copy `~/.pass-cli/vault.enc` regularly

### Scripting

1. **Always use `--quiet`** - Clean output for variables
2. **Check exit codes** - Handle errors properly
3. **Use `--field`** - Extract exactly what you need
4. **Redirect stderr** - Control error output

### Examples

**Good:**
```bash
export API_KEY=$(pass-cli get service --quiet 2>/dev/null)
if [ -z "$API_KEY" ]; then
    echo "Failed to get credential" >&2
    exit 1
fi
```

**Bad:**
```bash
# Don't do this - exposes password in process list
pass-cli add service --password mySecretPassword
```

## Common Patterns

### CI/CD Pipeline

```bash
# Retrieve deployment credentials
export DEPLOY_KEY=$(pass-cli get production --quiet)
export DB_PASSWORD=$(pass-cli get prod-db --quiet)

# Run deployment
./deploy.sh
```

### Local Development

```bash
# Set up environment from credentials
export DB_HOST=$(pass-cli get dev-db --field url --quiet)
export DB_USER=$(pass-cli get dev-db --field username --quiet)
export DB_PASS=$(pass-cli get dev-db --quiet)

# Start development server
npm run dev
```

### Credential Rotation

```bash
# Generate new password
NEW_PWD=$(pass-cli generate --length 32 --quiet)

# Update service
pass-cli update myservice --password "$NEW_PWD"

# Use new password
echo "$NEW_PWD" | some-service-update-command
```

## Getting Help

- Run any command with `--help` flag
- See [README](../README.md) for overview
- Check [Troubleshooting Guide](TROUBLESHOOTING.md) for common issues
- Visit [GitHub Issues](https://github.com/ari1110/pass-cli/issues)

---

**Documentation Version**: v0.0.1
**Last Updated**: October 2025
**Status**: Production Ready
