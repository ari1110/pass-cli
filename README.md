# Pass-CLI

> A secure, cross-platform command-line password manager designed for developers

Pass-CLI is a fast, secure password and API key manager that stores credentials locally with AES-256-GCM encryption. Built for developers who need quick, script-friendly access to credentials without cloud dependencies.

## ✨ Key Features

- **🔒 Military-Grade Encryption**: AES-256-GCM with PBKDF2 key derivation (100,000 iterations)
- **🔐 System Keychain Integration**: Seamless integration with Windows Credential Manager, macOS Keychain, and Linux Secret Service
- **⚡ Lightning Fast**: Sub-100ms credential retrieval for smooth workflows
- **🖥️ Cross-Platform**: Single binary for Windows, macOS (Intel/ARM), and Linux (amd64/arm64)
- **📋 Clipboard Support**: Automatic credential copying with security timeouts
- **🔑 Password Generation**: Cryptographically secure random passwords
- **📊 Usage Tracking**: Automatic tracking of where credentials are used
- **🤖 Script-Friendly**: Clean output modes for shell integration (`--quiet`, `--field`, `--masked`, `--no-clipboard`)
- **🔌 Offline First**: No cloud dependencies, works completely offline
- **📦 Easy Installation**: Available via Homebrew and Scoop

## 🚀 Quick Start

### Installation

#### macOS / Linux (Homebrew)

```bash
# Add tap and install
brew tap ari1110/homebrew-tap
brew install pass-cli
```

#### Windows (Scoop)

```powershell
# Add bucket and install
scoop bucket add pass-cli https://github.com/ari1110/scoop-bucket
scoop install pass-cli
```

#### Manual Installation

Download the latest release for your platform from [GitHub Releases](https://github.com/ari1110/pass-cli/releases):

```bash
# Extract and move to PATH
tar -xzf pass-cli_*_<os>_<arch>.tar.gz
sudo mv pass-cli /usr/local/bin/  # macOS/Linux
# Or move pass-cli.exe to a directory in PATH (Windows)
```

### First Steps

```bash
# Initialize your vault
pass-cli init

# Add your first credential
pass-cli add github
# Enter username and password when prompted

# Retrieve it
pass-cli get github

# Copy password to clipboard
pass-cli get github --copy

# Use in scripts (quiet mode)
export API_KEY=$(pass-cli get myservice --quiet --field password)
```

## 📖 Usage

### Initialize Vault

```bash
# Create a new vault
pass-cli init

# The vault is stored at ~/.pass-cli/
```

### Add Credentials

```bash
# Interactive mode (prompts for username/password)
pass-cli add myservice

# With URL and notes
pass-cli add github --url https://github.com --notes "Personal account"

# Generate a strong password
pass-cli add newservice --generate

# Custom password generation
pass-cli add newservice --generate --length 32 --no-symbols
```

### Retrieve Credentials

```bash
# Display credential (formatted)
pass-cli get myservice

# Copy password to clipboard
pass-cli get myservice --copy

# Quiet mode for scripts (password only)
pass-cli get myservice --quiet

# Get specific field
pass-cli get myservice --field username

# Display with masked password
pass-cli get myservice --masked
```

### List Credentials

```bash
# List all credentials
pass-cli list

# Show unused credentials
pass-cli list --unused
```

### Update Credentials

```bash
# Update password (prompted)
pass-cli update myservice

# Update specific fields
pass-cli update myservice --username newuser@example.com
pass-cli update myservice --url https://new-url.com
pass-cli update myservice --notes "Updated notes"

# Regenerate password
pass-cli update myservice --generate
```

### Delete Credentials

```bash
# Delete a credential (with confirmation)
pass-cli delete myservice

# Force delete (no confirmation)
pass-cli delete myservice --force
```

### Generate Passwords

```bash
# Generate a password (default: 20 chars, alphanumeric + symbols)
pass-cli generate

# Custom length
pass-cli generate --length 32

# Alphanumeric only (no symbols)
pass-cli generate --no-symbols

# Copy to clipboard
pass-cli generate --copy
```

### Version Information

```bash
# Check version
pass-cli version

# Verbose version info
pass-cli version --verbose
```

## 🔐 Security

### Encryption

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Derivation**: PBKDF2-SHA256 with 100,000 iterations
- **Salt**: Unique 32-byte random salt per vault
- **Authentication**: Built-in authentication tag (GCM) prevents tampering
- **IV**: Unique initialization vector per credential

### Master Password Storage

Pass-CLI integrates with your operating system's secure credential storage:

- **Windows**: Windows Credential Manager
- **macOS**: Keychain
- **Linux**: Secret Service (GNOME Keyring, KWallet)

Your master password is stored securely and unlocked automatically when needed.

### Clipboard Security

When using `--copy`, the clipboard is:
1. Cleared after 30 seconds automatically
2. Only contains the password (no metadata)
3. Can be cleared immediately with Ctrl+C

### Vault Location

The encrypted vault is stored at:
- **Windows**: `%USERPROFILE%\.pass-cli\vault.enc`
- **macOS/Linux**: `~/.pass-cli/vault.enc`

### Best Practices

- ✅ Use a strong, unique master password (20+ characters)
- ✅ Keep your vault backed up (it's just a file!)
- ✅ Use `--generate` for new passwords
- ✅ Regularly update credentials
- ✅ Use `--quiet` mode in scripts to avoid logging sensitive data
- ❌ Don't commit vault files to version control
- ❌ Don't share your master password

## 🤖 Script Integration

### Shell Integration Examples

**Bash/Zsh**:

```bash
#!/bin/bash
# Export API key for use in script
export API_KEY=$(pass-cli get openai --quiet --field password)

# Use in curl
curl -H "Authorization: Bearer $(pass-cli get github --quiet)" \
     https://api.github.com/user

# Conditional on success
if pass-cli get myservice --quiet > /dev/null 2>&1; then
    echo "Credential exists"
fi
```

**PowerShell**:

```powershell
# Store credential in variable
$apiKey = pass-cli get myservice --quiet --field password

# Use in web request
$headers = @{
    "Authorization" = "Bearer $apiKey"
}
Invoke-RestMethod -Uri "https://api.example.com" -Headers $headers

# Use with environment variable
$env:DATABASE_PASSWORD = pass-cli get postgres --quiet
```

**CI/CD Integration**:

```yaml
# GitHub Actions example
steps:
  - name: Retrieve credentials
    run: |
      export DB_PASSWORD=$(pass-cli get database --quiet)
      ./deploy.sh
```

### Output Modes

| Flag | Output | Use Case |
|------|--------|----------|
| (default) | Formatted table | Human-readable display |
| `--quiet` | Password only | Scripts, export to variables |
| `--field <name>` | Specific field | Extract username, URL, etc. |
| `--masked` | Masked password | Display password as asterisks |
| `--copy` | Clipboard | Quick copy without display |

## 📊 Usage Tracking

Pass-CLI automatically tracks where credentials are accessed based on your current working directory:

```bash
# Access from different directories
cd ~/project-a
pass-cli get database

cd ~/project-b
pass-cli get database

# View usage information with list command
pass-cli list --unused --days 30
```

This helps you:
- Identify which projects use which credentials
- Track credential usage patterns
- Audit credential access
- Find unused credentials for cleanup

## 🛠️ Advanced Usage

### Custom Vault Location

```bash
# Use a custom vault location
pass-cli --vault /path/to/custom/vault.enc list

# Or set via environment variable
export PASS_CLI_VAULT=/path/to/custom/vault.enc
pass-cli list
```

### Verbose Logging

```bash
# Enable verbose output for debugging
pass-cli --verbose get myservice
```

### Configuration File

Create `~/.pass-cli/config.yaml`:

```yaml
# Default configuration
verbose: false
vault: ~/.pass-cli/vault.enc
```

## 🏗️ Building from Source

### Prerequisites

- Go 1.25 or later
- Make (optional, for convenience)

### Build

```bash
# Clone the repository
git clone https://github.com/ari1110/pass-cli.git
cd pass-cli

# Build binary
go build -o pass-cli .

# Or use Make
make build

# Run tests
make test

# Run with coverage
make test-coverage
```

## 📝 Development

### Running Tests

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./test/

# All tests
make test-all
```

### Code Quality

```bash
# Run linter
make lint

# Security scan
make security-scan

# Format code
make fmt
```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) first.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📋 Requirements

- **Operating System**: Windows 10+, macOS 10.15+, or Linux (any modern distribution)
- **Architecture**: amd64 or arm64
- **Dependencies**: None (static binary)

## 🗺️ Roadmap

- [x] Core credential management (add, get, list, update, delete)
- [x] AES-256-GCM encryption
- [x] System keychain integration
- [x] Password generation
- [x] Clipboard support
- [x] Usage tracking
- [ ] Import from other password managers
- [ ] Export functionality
- [ ] Credential sharing (encrypted)
- [ ] Two-factor authentication support
- [ ] Browser extensions

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Documentation**: [docs/](docs/)
- **Releases**: [GitHub Releases](https://github.com/ari1110/pass-cli/releases)
- **Issues**: [GitHub Issues](https://github.com/ari1110/pass-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ari1110/pass-cli/discussions)

## ❓ FAQ

### How is this different from `pass` (the standard Unix password manager)?

Pass-CLI offers:
- System keychain integration (no GPG required)
- Built-in clipboard support
- Usage tracking
- Cross-platform Windows support
- Script-friendly output modes (--quiet, --field, --masked)
- Single binary distribution

### Is my data stored in the cloud?

No. Pass-CLI stores everything locally on your machine. There are no cloud dependencies or network calls.

### Can I sync my vault across machines?

The vault is a single encrypted file (`~/.pass-cli/vault.enc`). You can sync it using any file sync service (Dropbox, Google Drive, etc.), but be aware of potential conflicts if editing from multiple machines simultaneously.

### What happens if I forget my master password?

Unfortunately, there's no way to recover your vault without the master password. The encryption is designed to be unbreakable. Keep your master password safe and consider backing it up securely.

### How do I backup my vault?

Simply copy the vault file:
```bash
cp ~/.pass-cli/vault.enc ~/backup/vault-$(date +%Y%m%d).enc
```

### Can I use Pass-CLI in my company?

Yes! Pass-CLI is MIT licensed and free for commercial use. It's designed for professional developer workflows.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [go-keyring](https://github.com/zalando/go-keyring) for system keychain integration
- Inspired by the Unix `pass` password manager

## 📞 Support

- **Bug Reports**: [GitHub Issues](https://github.com/ari1110/pass-cli/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/ari1110/pass-cli/discussions)
- **Security Issues**: Email security@example.com (please don't file public issues)

---

Made with ❤️ by developers, for developers. Star ⭐ this repo if you find it useful!
