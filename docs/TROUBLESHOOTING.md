# Troubleshooting Guide

Solutions to common issues and frequently asked questions for Pass-CLI.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Initialization Issues](#initialization-issues)
- [Keychain Access Issues](#keychain-access-issues)
- [Vault Access Issues](#vault-access-issues)
- [Command Issues](#command-issues)
- [Platform-Specific Issues](#platform-specific-issues)
- [Performance Issues](#performance-issues)
- [Vault Recovery](#vault-recovery)
- [Frequently Asked Questions](#frequently-asked-questions)

## Installation Issues

### Command Not Found After Installation

**Symptom**: `pass-cli: command not found` or `'pass-cli' is not recognized`

**Cause**: Binary not in system PATH

**Solutions**:

**macOS/Linux:**
```bash
# Check if binary exists
which pass-cli

# If not found, check installation location
ls -la /usr/local/bin/pass-cli
ls -la ~/.local/bin/pass-cli

# Add to PATH if needed (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:$HOME/.local/bin"
source ~/.bashrc

# Verify
pass-cli version
```

**Windows:**
```powershell
# Check if binary exists
where.exe pass-cli

# Add to PATH
$path = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = "$path;C:\path\to\pass-cli"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")

# Restart PowerShell
exit

# Verify
pass-cli version
```

---

### Permission Denied When Running

**Symptom**: `Permission denied` when executing pass-cli

**Cause**: Binary doesn't have execute permissions

**Solution (macOS/Linux)**:
```bash
# Add execute permission
chmod +x /path/to/pass-cli

# Or reinstall with correct permissions
sudo install -m 755 pass-cli /usr/local/bin/
```

---

### Homebrew Installation Fails

**Symptom**: `Error: No such file or directory` or tap not found

**Solutions**:
```bash
# Update Homebrew
brew update

# Check Homebrew status
brew doctor

# Remove and re-add tap
brew untap yourusername/pass-cli
brew tap yourusername/pass-cli

# Try verbose installation
brew install --verbose pass-cli

# Check logs if failed
brew gist-logs pass-cli
```

---

### Scoop Installation Fails

**Symptom**: Manifest not found or hash mismatch

**Solutions**:
```powershell
# Update Scoop
scoop update

# Check Scoop status
scoop status

# Remove and re-add bucket
scoop bucket rm pass-cli
scoop bucket add pass-cli https://github.com/yourusername/scoop-pass-cli

# Clear cache and retry
scoop cache rm pass-cli
scoop install pass-cli

# Check logs
scoop cat pass-cli
```

---

### macOS "Cannot Open" Security Warning

**Symptom**: "pass-cli cannot be opened because the developer cannot be verified"

**Cause**: macOS Gatekeeper blocks unsigned binaries

**Solutions**:

**Option 1: Remove quarantine attribute**
```bash
xattr -d com.apple.quarantine /usr/local/bin/pass-cli
```

**Option 2: Allow in System Preferences**
1. Try to run pass-cli
2. Open System Preferences → Security & Privacy
3. Click "Allow Anyway" next to pass-cli message
4. Run pass-cli again and confirm

**Option 3: Build from source** (trusted)
```bash
git clone https://github.com/yourusername/pass-cli
cd pass-cli
make build
sudo mv pass-cli /usr/local/bin/
```

---

## Initialization Issues

### "Vault Already Exists" Error

**Symptom**: `Error: vault already exists at ~/.pass-cli/vault.enc`

**Cause**: Trying to initialize when vault already exists

**Solutions**:

**Option 1: Use existing vault**
```bash
# Try to access existing vault
pass-cli list

# If you remember the password, continue using it
```

**Option 2: Backup and reinitialize**
```bash
# Backup existing vault
mv ~/.pass-cli/vault.enc ~/.pass-cli/vault.enc.old

# Initialize new vault
pass-cli init

# If needed, you can restore old vault later
# mv ~/.pass-cli/vault.enc.old ~/.pass-cli/vault.enc
```

**Option 3: Use different vault location**
```bash
pass-cli --vault /path/to/new/vault.enc init
```

---

### "Failed to Store Master Password" Error

**Symptom**: Error when saving master password to keychain

**Cause**: Keychain service not available or permission denied

**Solutions**:

**macOS:**
```bash
# Check keychain status
security list-keychains

# Unlock login keychain
security unlock-keychain ~/Library/Keychains/login.keychain-db

# Verify keychain access
security add-generic-password -a "$USER" -s "test" -w "test"
security delete-generic-password -a "$USER" -s "test"
```

**Linux:**
```bash
# Check if secret service is running
ps aux | grep -i "gnome-keyring\|kwallet"

# Start GNOME Keyring (if not running)
gnome-keyring-daemon --start

# Or install if missing
sudo apt install gnome-keyring  # Ubuntu/Debian
sudo dnf install gnome-keyring  # Fedora
```

**Windows:**
```powershell
# Run as administrator
# Check Credential Manager service
Get-Service -Name "VaultSvc"

# Start if stopped
Start-Service -Name "VaultSvc"
```

---

## Keychain Access Issues

### "Failed to Retrieve Master Password" Error

**Symptom**: Cannot get master password from keychain

**Cause**: Master password not stored or keychain locked

**Solutions**:

1. **Verify keychain entry exists**

   **macOS:**
   ```bash
   # Check Keychain Access app
   # Search for "pass-cli" entry
   ```

   **Linux:**
   ```bash
   # Check with Seahorse (GNOME) or KWalletManager (KDE)
   ```

   **Windows:**
   ```powershell
   # Check Credential Manager
   # Control Panel → User Accounts → Credential Manager
   # Windows Credentials → Look for "pass-cli"
   ```

2. **Reinitialize vault** (will prompt for password again)
   ```bash
   # Backup vault first
   cp ~/.pass-cli/vault.enc ~/vault-backup.enc

   # Reinitialize (keeps vault but updates keychain)
   pass-cli init
   ```

---

### Keychain Access Denied

**Symptom**: "Access denied" when accessing keychain

**Cause**: Keychain locked or permission issues

**Solutions**:

**macOS:**
```bash
# Unlock keychain
security unlock-keychain ~/Library/Keychains/login.keychain-db

# Grant access to pass-cli
# Will prompt when pass-cli runs - click "Always Allow"
```

**Linux (GNOME):**
```bash
# Unlock keyring
# Will prompt for keyring password when pass-cli runs

# If keyring password is different from login password
# Open Seahorse → Right-click Login → Change Password
```

**Windows:**
```powershell
# Ensure running as correct user
whoami

# Credential Manager uses Windows login credentials
# Ensure logged in as user who created vault
```

---

### "Secret Service Not Available" (Linux)

**Symptom**: Cannot access secret service on Linux

**Cause**: Secret service daemon not running

**Solutions**:

**GNOME:**
```bash
# Install GNOME Keyring
sudo apt install gnome-keyring  # Ubuntu/Debian
sudo dnf install gnome-keyring  # Fedora

# Start daemon
gnome-keyring-daemon --start --components=secrets

# Add to session startup
# Add to ~/.profile or ~/.bash_profile:
eval $(gnome-keyring-daemon --start --components=secrets)
```

**KDE:**
```bash
# Install KWallet
sudo apt install kwalletmanager  # Ubuntu/Debian

# Start KWallet
kwalletd5 &
```

**Alternative: File-based password** (less secure)
```bash
# Store password in encrypted file (not recommended)
# Use vault without keychain integration
# Enter password each time
```

---

## Vault Access Issues

### "Invalid Master Password" Error

**Symptom**: Password rejected when accessing vault

**Cause**: Incorrect password or vault corruption

**Solutions**:

1. **Verify password**
   - Check caps lock
   - Try typing slowly
   - Copy-paste if stored elsewhere

2. **Check keychain entry**
   ```bash
   # macOS: View in Keychain Access
   # Linux: View in Seahorse/KWallet
   # Windows: View in Credential Manager
   ```

3. **Restore from backup**
   ```bash
   # If vault is corrupted
   cp ~/.pass-cli/vault.enc.backup ~/.pass-cli/vault.enc
   ```

4. **Try manual backup**
   ```bash
   # If you have manual backup
   cp ~/backups/vault-20250120.enc ~/.pass-cli/vault.enc
   ```

---

### "Vault File Corrupted" Error

**Symptom**: Cannot decrypt vault, corruption detected

**Cause**: File corruption from crash or disk error

**Solutions**:

1. **Restore automatic backup**
   ```bash
   # Check if backup exists
   ls -la ~/.pass-cli/vault.enc.backup

   # Restore
   cp ~/.pass-cli/vault.enc.backup ~/.pass-cli/vault.enc

   # Test access
   pass-cli list
   ```

2. **Restore manual backup**
   ```bash
   # List available backups
   ls -la ~/backups/vault-*.enc

   # Restore most recent
   cp ~/backups/vault-20250120.enc ~/.pass-cli/vault.enc
   ```

3. **If no backup available**
   ```bash
   # Unfortunately, corrupted vault without backup is unrecoverable
   # Initialize new vault
   mv ~/.pass-cli/vault.enc ~/.pass-cli/vault.enc.corrupted
   pass-cli init
   # Re-add credentials manually
   ```

---

### "Permission Denied" Reading Vault

**Symptom**: Cannot read vault file

**Cause**: File permission issues

**Solutions**:

**macOS/Linux:**
```bash
# Check permissions
ls -la ~/.pass-cli/vault.enc

# Fix permissions (should be 0600)
chmod 600 ~/.pass-cli/vault.enc

# Fix ownership
sudo chown $USER:$USER ~/.pass-cli/vault.enc
```

**Windows:**
```powershell
# Check ACL
Get-Acl "$env:USERPROFILE\.pass-cli\vault.enc" | Format-List

# Reset permissions to current user
$acl = Get-Acl "$env:USERPROFILE\.pass-cli\vault.enc"
$acl.SetAccessRuleProtection($true, $false)
$rule = New-Object System.Security.AccessControl.FileSystemAccessRule(
    $env:USERNAME, "FullControl", "Allow"
)
$acl.AddAccessRule($rule)
Set-Acl "$env:USERPROFILE\.pass-cli\vault.enc" $acl
```

---

## Command Issues

### "Service Already Exists" Error

**Symptom**: Cannot add credential, service name already exists

**Cause**: Duplicate service name

**Solutions**:

```bash
# Check existing credentials
pass-cli list

# Update existing instead
pass-cli update <service>

# Or delete and re-add
pass-cli delete <service>
pass-cli add <service>

# Use different name
pass-cli add <service>-prod
```

---

### "Service Not Found" Error

**Symptom**: Cannot get/update/delete non-existent service

**Cause**: Service name doesn't exist or typo

**Solutions**:

```bash
# List all services to check spelling
pass-cli list

# Check for case sensitivity
pass-cli list | grep -i <service>

# If service was deleted, re-add
pass-cli add <service>
```

---

### Clipboard Copy Fails

**Symptom**: "Failed to copy to clipboard" error

**Cause**: Clipboard not available or permission denied

**Solutions**:

**macOS:**
```bash
# Usually works by default
# If fails, try without clipboard
pass-cli get <service> --no-clipboard
```

**Linux:**
```bash
# Install xclip or xsel
sudo apt install xclip  # Ubuntu/Debian

# Or use wl-clipboard for Wayland
sudo apt install wl-clipboard

# Verify
echo "test" | xclip -selection clipboard
```

**Windows:**
```powershell
# Usually works by default
# If fails, try without clipboard
pass-cli get <service> --no-clipboard
```

**Workaround:**
```bash
# Use quiet mode and manual copy
pass-cli get <service> --quiet
# Then copy output manually
```

---

## Platform-Specific Issues

### Windows

#### Antivirus Blocks Pass-CLI

**Symptom**: Antivirus quarantines or blocks pass-cli.exe

**Cause**: False positive from security software

**Solutions**:
1. Add exception in antivirus software
2. Whitelist `pass-cli.exe` and vault directory
3. Download from official source and verify checksum
4. Build from source if concerned

#### PowerShell Execution Policy

**Symptom**: Cannot run pass-cli in PowerShell

**Cause**: Execution policy restrictions

**Solution**:
```powershell
# Check policy
Get-ExecutionPolicy

# Set policy (if needed)
Set-ExecutionPolicy -Scope CurrentUser RemoteSigned

# Or run directly
.\pass-cli.exe version
```

---

### macOS

#### Keychain Prompt Every Time

**Symptom**: macOS asks for keychain password on every use

**Cause**: Pass-CLI not trusted for keychain access

**Solution**:
1. Open Keychain Access app
2. Search for "pass-cli"
3. Double-click the entry
4. Go to "Access Control" tab
5. Select "Allow all applications to access this item"
6. Save changes

---

### Linux

#### D-Bus Session Issues

**Symptom**: "D-Bus session not available" error

**Cause**: D-Bus session bus not running

**Solution**:
```bash
# Check D-Bus
echo $DBUS_SESSION_BUS_ADDRESS

# Start D-Bus if missing
eval $(dbus-launch --sh-syntax)

# Add to shell startup (~/.bashrc)
if [ -z "$DBUS_SESSION_BUS_ADDRESS" ]; then
    eval $(dbus-launch --sh-syntax)
fi
```

#### SELinux/AppArmor Blocking Access

**Symptom**: Permission denied despite correct file permissions

**Cause**: SELinux or AppArmor policy

**Solutions**:

**SELinux:**
```bash
# Check if SELinux is enforcing
getenforce

# Temporarily disable (for testing)
sudo setenforce 0

# Or create policy for pass-cli
sudo audit2allow -a -M pass-cli
sudo semodule -i pass-cli.pp
```

**AppArmor:**
```bash
# Check AppArmor status
sudo aa-status

# Disable for testing
sudo systemctl stop apparmor

# Or create profile
sudo aa-complain /usr/local/bin/pass-cli
```

---

## Performance Issues

### Slow Unlock Times

**Symptom**: First access takes several seconds

**Cause**: PBKDF2 key derivation (100,000 iterations)

**Expected Behavior**:
- First unlock: 300-500ms (key derivation)
- Cached access: <100ms

**Solutions**:

If slower than expected:
```bash
# Check system resources
top

# Ensure no CPU throttling
# Ensure SSD not slow (HDD will be slower)

# Run on faster machine if needed
```

---

### Large Vault Slow

**Symptom**: Operations slow with many credentials

**Cause**: Loading entire vault into memory

**Expected Behavior**:
- <100 credentials: No noticeable delay
- 100-1000 credentials: <1s
- >1000 credentials: May be slower

**Solutions**:
```bash
# Split into multiple vaults by purpose
pass-cli --vault ~/.pass-cli/work.enc init
pass-cli --vault ~/.pass-cli/personal.enc init

# Archive unused credentials
pass-cli list --unused --days 90
# Delete unused ones
```

---

## Vault Recovery

### Forgot Master Password

**Symptom**: Cannot remember master password

**Unfortunate Reality**: **There is no recovery mechanism**

**Options**:

1. **Try to remember**
   - Try common variations
   - Check password manager if stored there
   - Check secure notes

2. **Check keychain** (if still accessible)
   - macOS: Keychain Access → search "pass-cli"
   - Linux: Seahorse → search "pass-cli"
   - Windows: Credential Manager → search "pass-cli"

3. **If truly lost**
   ```bash
   # Vault is unrecoverable
   # Start fresh
   mv ~/.pass-cli/vault.enc ~/.pass-cli/vault.enc.lost
   pass-cli init
   # Re-add credentials from services
   ```

**Prevention**:
- Write master password in secure location
- Store in another password manager
- Keep backup of master password

---

### Vault File Deleted

**Symptom**: Vault file missing

**Solutions**:

1. **Check trash/recycle bin**

2. **Restore from backup**
   ```bash
   # Automatic backup
   cp ~/.pass-cli/vault.enc.backup ~/.pass-cli/vault.enc

   # Manual backup
   cp ~/backups/vault-*.enc ~/.pass-cli/vault.enc
   ```

3. **Restore from cloud sync** (if syncing vault)
   ```bash
   # From Dropbox, Google Drive, etc.
   cp ~/Dropbox/vault.enc ~/.pass-cli/vault.enc
   ```

4. **If no backup**
   ```bash
   # Unfortunately, must start over
   pass-cli init
   ```

---

### Corrupt Vault Recovery

**Symptom**: Vault fails to decrypt or shows corruption errors

**Solutions**:

1. **Try automatic backup**
   ```bash
   cp ~/.pass-cli/vault.enc.backup ~/.pass-cli/vault.enc
   pass-cli list
   ```

2. **Try older backups**
   ```bash
   # List backups by date
   ls -lt ~/backups/vault-*.enc

   # Try each, newest first
   cp ~/backups/vault-20250120.enc ~/.pass-cli/vault.enc
   pass-cli list
   ```

3. **Attempt partial recovery** (advanced)
   ```bash
   # Examine vault file
   hexdump -C ~/.pass-cli/vault.enc | head -n 20

   # Check file size
   ls -la ~/.pass-cli/vault.enc

   # If file is obviously truncated or wrong size
   # Recovery likely impossible, use backup
   ```

---

## Frequently Asked Questions

### General Questions

**Q: Where is my vault stored?**

A:
- Windows: `%USERPROFILE%\.pass-cli\vault.enc`
- macOS/Linux: `~/.pass-cli/vault.enc`

---

**Q: How do I backup my vault?**

A:
```bash
# Simple copy
cp ~/.pass-cli/vault.enc ~/backups/vault-$(date +%Y%m%d).enc

# Automated daily backup (cron)
0 0 * * * cp ~/.pass-cli/vault.enc ~/backups/vault-$(date +%Y%m%d).enc
```

---

**Q: Can I sync my vault across machines?**

A: Yes, but carefully:
```bash
# Using cloud storage (manual)
cp ~/.pass-cli/vault.enc ~/Dropbox/vault.enc

# On other machine
cp ~/Dropbox/vault.enc ~/.pass-cli/vault.enc

# Warning: Conflicts if editing on multiple machines simultaneously
```

---

**Q: How do I change my master password?**

A: Currently requires manual process:
```bash
# 1. Export all credentials to temporary file (manually)
pass-cli list --format json > /tmp/credentials.json

# 2. Initialize new vault with new password
mv ~/.pass-cli/vault.enc ~/.pass-cli/vault.enc.old
pass-cli init

# 3. Re-add credentials manually
# (No automated import yet)

# 4. Verify and delete old vault
pass-cli list
rm ~/.pass-cli/vault.enc.old
```

---

**Q: Is my data sent to the cloud?**

A: No. Pass-CLI:
- ✅ Works completely offline
- ✅ Never makes network calls
- ✅ Stores everything locally
- ✅ No telemetry or tracking

---

**Q: What happens if I lose my vault file?**

A:
- If you have backup: Restore from backup
- If no backup: All credentials lost, must start over
- Prevention: Regular backups essential

---

### Technical Questions

**Q: Can I use Pass-CLI in scripts?**

A: Yes, designed for it:
```bash
# Use quiet mode
export API_KEY=$(pass-cli get service --quiet)

# Extract specific field
export USERNAME=$(pass-cli get service --field username --quiet)

# JSON output for parsing
pass-cli list --format json | jq '.[] | .service'
```

---

**Q: How secure is Pass-CLI?**

A: See [SECURITY.md](SECURITY.md) for full details:
- AES-256-GCM encryption
- PBKDF2 key derivation (100,000 iterations)
- System keychain integration
- No cloud dependencies
- Limitations explained in security doc

---

**Q: Can multiple users share a vault?**

A: Not designed for this:
- Vault is single-user
- Master password would be shared (insecure)
- No access control mechanism
- Use separate vaults per user

---

**Q: What if I forget a specific credential password?**

A: Individual credentials cannot be recovered:
- Vault decrypts all-or-nothing
- If vault accessible, all credentials accessible
- If vault locked, all credentials inaccessible
- No per-credential password recovery

---

## Getting Further Help

### Before Asking for Help

1. **Check this troubleshooting guide**
2. **Read relevant documentation**:
   - [Installation Guide](INSTALLATION.md)
   - [Usage Guide](USAGE.md)
   - [Security Documentation](SECURITY.md)
3. **Search existing issues**: [GitHub Issues](https://github.com/yourusername/pass-cli/issues)

### Reporting Issues

When reporting issues, include:

```bash
# System information
pass-cli version --verbose
uname -a  # or `systeminfo` on Windows

# Error message (full output)
pass-cli <command> --verbose 2>&1

# Steps to reproduce
# 1. Run: pass-cli init
# 2. Run: pass-cli add test
# 3. Error occurs
```

### Getting Help

- **GitHub Issues**: [Report bugs](https://github.com/yourusername/pass-cli/issues/new)
- **GitHub Discussions**: [Ask questions](https://github.com/yourusername/pass-cli/discussions)
- **Security Issues**: Email security@example.com (don't file publicly)

---

**Last Updated**: 2025-01-20
**Version**: 1.0.0
