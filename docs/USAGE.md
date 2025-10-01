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

#### Notes

- Master password must be at least 8 characters
- Strong passwords (20+ characters) are recommended
- Master password is stored in OS keychain for convenience
- Vault file is created with restricted permissions (0600)

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
| `--url` | | string | Service URL |
| `--notes` | | string | Additional notes |
| `--generate` | | bool | Generate a random password |
| `--length` | | int | Password length for generation (default: 20) |
| `--no-symbols` | | bool | Exclude symbols when generating |

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

# Generate a password during add
pass-cli add newservice --generate

# Generate with custom length and no symbols
pass-cli add newservice --generate --length 32 --no-symbols

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

#### Notes

- Service names must be unique
- Password input is hidden by default
- Passing password via `-p` flag is insecure (visible in shell history)
- Use `--generate` for strong random passwords
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
| `--copy` | `-c` | bool | Copy to clipboard only (no display) |
| `--no-clipboard` | | bool | Skip clipboard copy |
| `--masked` | | bool | Display password as asterisks |
| `--json` | | bool | Output as JSON |

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

# Copy to clipboard only (no display)
pass-cli get github --copy

# Display without clipboard
pass-cli get github --no-clipboard

# Display with masked password
pass-cli get github --masked

# JSON output
pass-cli get github --json
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

**JSON output:**
```json
{
  "service": "github",
  "username": "user@example.com",
  "password": "mySecretPassword123!",
  "url": "https://github.com",
  "notes": "Personal account",
  "created": "2025-01-15T10:30:00Z",
  "modified": "2025-01-15T10:30:00Z",
  "accessed": "2025-01-20T14:22:00Z",
  "used_in": ["/home/user/project-a", "/home/user/project-b"]
}
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

**JSON format:**
```json
[
  {
    "service": "github",
    "username": "user@example.com",
    "url": "https://github.com",
    "created": "2025-01-15T10:30:00Z",
    "modified": "2025-01-15T10:30:00Z",
    "accessed": "2025-01-20T14:22:00Z"
  },
  ...
]
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
| `--url` | | string | New URL |
| `--notes` | | string | New notes |
| `--generate` | | bool | Generate new random password |
| `--length` | | int | Password length for generation (default: 20) |
| `--no-symbols` | | bool | Exclude symbols when generating |

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

# Generate new password
pass-cli update github --generate

# Generate with custom options
pass-cli update github --generate --length 32 --no-symbols

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
| `--copy` | bool | Copy to clipboard only (no display) |
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

# Copy to clipboard without displaying
pass-cli generate --copy

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
pass-cli version 1.0.0
```

**Verbose:**
```
pass-cli version 1.0.0
  commit: abc123f
  built:  2025-01-20T10:30:00Z
  go:     go1.25.0
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

### JSON Mode

Structured data for parsing.

```bash
pass-cli get github --json | jq '.username'
# "user@example.com"

pass-cli list --format json | jq '.[].service'
# "github"
# "aws-prod"
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
import json

# Get credential
result = subprocess.run(
    ['pass-cli', 'get', 'github', '--json'],
    capture_output=True,
    text=True,
    check=True
)
cred = json.loads(result.stdout)

# Use credential
username = cred['username']
password = cred['password']

# Get password only
result = subprocess.run(
    ['pass-cli', 'get', 'github', '--quiet'],
    capture_output=True,
    text=True,
    check=True
)
password = result.stdout.strip()
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

Configuration file: `~/.pass-cli/config.yaml`

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
```

### Configuration Priority

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file
4. Built-in defaults (lowest priority)

## Usage Tracking

Pass-CLI automatically tracks where credentials are accessed based on your current working directory.

### How It Works

```bash
# Access from project directory
cd ~/projects/my-app
pass-cli get database

# Later, view tracking info
pass-cli get database --json | jq '.used_in'
# ["/home/user/projects/my-app"]
```

### Use Cases

- **Audit**: See which projects use which credentials
- **Cleanup**: Identify unused credentials
- **Documentation**: Auto-document credential dependencies

### Viewing Usage

```bash
# JSON output includes usage tracking
pass-cli get myservice --json

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
- Visit [GitHub Issues](https://github.com/yourusername/pass-cli/issues)
