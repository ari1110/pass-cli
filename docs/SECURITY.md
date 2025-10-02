# Security Documentation

Comprehensive security architecture, cryptographic implementation, and best practices for Pass-CLI.

## Table of Contents

- [Security Overview](#security-overview)
- [Cryptographic Implementation](#cryptographic-implementation)
- [Master Password Management](#master-password-management)
- [Data Storage Security](#data-storage-security)
- [Threat Model](#threat-model)
- [Security Guarantees](#security-guarantees)
- [Limitations](#limitations)
- [Best Practices](#best-practices)
- [Security Checklist](#security-checklist)
- [Incident Response](#incident-response)
- [Security Audits](#security-audits)

## Security Overview

Pass-CLI is designed with security as the primary concern. All credentials are encrypted using industry-standard cryptography and stored locally on your machine with no cloud dependencies.

### Key Security Features

- **AES-256-GCM Encryption**: Military-grade authenticated encryption
- **PBKDF2 Key Derivation**: 100,000 iterations with SHA-256
- **System Keychain Integration**: Secure master password storage
- **Offline-First Design**: No network calls, no cloud dependencies
- **Secure Memory Handling**: Sensitive data cleared after use
- **File Permission Protection**: Vault files restricted to user-only access

## Cryptographic Implementation

### Encryption Algorithm

**AES-256-GCM (Galois/Counter Mode)**

- **Algorithm**: Advanced Encryption Standard
- **Key Size**: 256 bits (32 bytes)
- **Mode**: GCM (Galois/Counter Mode)
- **Authentication**: Built-in GMAC authentication tag
- **Implementation**: Go standard library `crypto/aes` and `crypto/cipher`

#### Why AES-256-GCM?

1. **NIST Approved**: Recommended by NIST for classified information
2. **Authenticated Encryption**: Prevents tampering and chosen-ciphertext attacks
3. **Parallelizable**: Fast performance on modern hardware
4. **Standard**: Widely used and well-audited implementation

### Key Derivation

**PBKDF2-SHA256**

- **Algorithm**: Password-Based Key Derivation Function 2
- **Hash Function**: SHA-256
- **Iterations**: 100,000
- **Salt Length**: 32 bytes (256 bits)
- **Output Length**: 32 bytes (256 bits)
- **Implementation**: `golang.org/x/crypto/pbkdf2`

#### Key Derivation Process

```
Master Key = PBKDF2(
    password = user's master password,
    salt = unique 32-byte random salt,
    iterations = 100,000,
    hash = SHA-256,
    key_length = 32 bytes
)
```

#### Why PBKDF2?

1. **Computationally Expensive**: 100,000 iterations slow down brute-force attacks
2. **Salted**: Unique salt prevents rainbow table attacks
3. **Standard**: NIST recommended for password-based cryptography
4. **Deterministic**: Same password + salt = same key

### Encryption Process

#### Encrypting Credentials

1. **Generate Salt** (first time only)
   ```
   salt = crypto/rand.Read(32 bytes)
   ```

2. **Derive Encryption Key**
   ```
   key = PBKDF2(master_password, salt, 100000, SHA256, 32)
   ```

3. **Generate Nonce**
   ```
   nonce = crypto/rand.Read(12 bytes)  // Per-encryption unique
   ```

4. **Encrypt Data**
   ```
   ciphertext = AES-256-GCM.Encrypt(
       plaintext = JSON(credentials),
       key = derived_key,
       nonce = nonce,
       additional_data = nil
   )
   ```

5. **Combine Components**
   ```
   vault_file = nonce || ciphertext || auth_tag
   ```

#### Decrypting Credentials

1. **Load Master Password** from system keychain
2. **Read Vault File** and extract salt, nonce, ciphertext
3. **Derive Key** using PBKDF2 with stored salt
4. **Decrypt and Verify**
   ```
   plaintext = AES-256-GCM.Decrypt(
       ciphertext,
       key,
       nonce
   )
   ```
5. **Parse JSON** to access credentials

### Random Number Generation

All random values use `crypto/rand`, which provides cryptographically secure random numbers from the operating system:

- **Windows**: `CryptGenRandom`
- **macOS/Linux**: `/dev/urandom`

Used for:
- Salt generation
- Nonce generation
- Password generation

## Master Password Management

### System Keychain Integration

Pass-CLI integrates with your operating system's secure credential storage to save your master password.

#### Windows - Credential Manager

- **Location**: Windows Credential Manager
- **Storage**: Encrypted by Windows using DPAPI
- **Access**: Protected by user's Windows login
- **Implementation**: `github.com/zalando/go-keyring`

**Viewing in Windows:**
1. Open Control Panel
2. User Accounts → Credential Manager
3. Windows Credentials
4. Look for "pass-cli" entry

#### macOS - Keychain

- **Location**: macOS Keychain (login keychain)
- **Storage**: Encrypted by macOS keychain services
- **Access**: Protected by user's macOS login password
- **Implementation**: `github.com/zalando/go-keyring`

**Viewing on macOS:**
1. Open Keychain Access app
2. Search for "pass-cli"
3. Double-click to view (requires password)

#### Linux - Secret Service

- **Backend**: GNOME Keyring, KWallet, or compatible
- **Protocol**: freedesktop.org Secret Service API
- **Storage**: Encrypted by keyring daemon
- **Access**: Protected by keyring password
- **Implementation**: `github.com/zalando/go-keyring`

**Viewing on Linux (GNOME):**
1. Open Seahorse (Passwords and Keys)
2. Login keyring
3. Search for "pass-cli"

### Master Password Requirements

- **Minimum Length**: 8 characters (enforced)
- **Recommended Length**: 20+ characters
- **Character Types**: All types allowed
- **Strength**: No maximum, use passphrase if preferred

### Master Password Security

**What Pass-CLI Does:**
- ✅ Stores master password in system keychain
- ✅ Clears password from memory after use
- ✅ Never writes password to disk in plaintext
- ✅ Never logs password

**What You Should Do:**
- ✅ Use a unique master password (not reused elsewhere)
- ✅ Make it strong (20+ characters or passphrase)
- ✅ Store backup securely (password manager, safe place)
- ❌ Don't share your master password
- ❌ Don't write it in plaintext files

## Data Storage Security

### Vault File Location

- **Windows**: `%USERPROFILE%\.pass-cli\vault.enc`
- **macOS/Linux**: `~/.pass-cli/vault.enc`

### File Permissions

Vault files are created with restricted permissions:

- **Unix (macOS/Linux)**: `0600` (owner read/write only)
- **Windows**: ACL restricting to current user

### Vault File Structure

```
+------------------+
| Salt (32 bytes)  |  ← PBKDF2 salt
+------------------+
| Nonce (12 bytes) |  ← AES-GCM nonce
+------------------+
| Ciphertext       |  ← Encrypted credentials (variable length)
+------------------+
| Auth Tag         |  ← GCM authentication tag (16 bytes)
+------------------+
```

### Atomic Writes

Vault updates use atomic write operations to prevent corruption:

1. Write to temporary file (`.vault.enc.tmp`)
2. Sync to disk (`fsync`)
3. Rename to actual vault file (atomic operation)
4. Delete temporary file on error

This ensures:
- No partial writes
- No corruption on crash
- Previous vault preserved on error

### Backup Strategy

Before each vault update:
1. Current vault backed up to `.vault.enc.backup`
2. New vault written atomically
3. Backup kept for disaster recovery

## Threat Model

### What Pass-CLI Protects Against

✅ **Offline Attacks**
- Vault file encryption protects against offline brute-force
- PBKDF2 slows down password cracking (100,000 iterations)
- No plaintext credentials stored anywhere

✅ **File System Compromise**
- Encrypted vault remains secure even if file is stolen
- File permissions prevent unauthorized local access

✅ **Process Memory Dumps**
- Sensitive data cleared from memory after use
- Master password not kept in memory permanently

✅ **Accidental Disclosure**
- No cloud storage = no cloud breach risk
- No network calls = no network interception

✅ **Unauthorized Local Access**
- System keychain protects master password
- File permissions restrict vault access

### What Pass-CLI Does NOT Protect Against

❌ **Malware on Your Machine**
- Keyloggers can capture master password when entered
- Memory scrapers can extract decrypted credentials
- Root/admin access bypasses file permissions

❌ **Physical Access Attacks**
- Attacker with physical access can copy vault file
- Vault encryption is only protection (strong password essential)

❌ **Side-Channel Attacks**
- Timing attacks, power analysis not mitigated
- Not designed for hostile multi-user systems

❌ **Weak Master Passwords**
- PBKDF2 slows attacks but doesn't prevent them
- Short/common passwords can be brute-forced

❌ **Social Engineering**
- Cannot protect against phishing for master password
- User education essential

## Security Guarantees

### What We Guarantee

1. **Confidentiality**: Credentials encrypted with AES-256-GCM
2. **Integrity**: Authentication tag prevents tampering
3. **Forward Secrecy**: Unique nonce per encryption
4. **Secure Defaults**: No insecure configuration options

### What We Cannot Guarantee

1. **Availability**: Forgot password = lost vault
2. **Recovery**: No backdoor or recovery mechanism
3. **Zero-Knowledge**: Master password accessible via keychain
4. **Perfect Security**: Subject to implementation bugs

## Limitations

### Known Limitations

1. **Master Password Recovery**: None available
   - If you forget master password, vault is unrecoverable
   - No "forgot password" mechanism
   - No backdoor or master key

2. **Keychain Dependency**
   - Master password security depends on OS keychain
   - Compromise of OS account = compromise of master password

3. **Single-User Design**
   - Not designed for multi-user systems
   - File permissions rely on OS access controls

4. **No Network Security**
   - Offline-only design
   - No secure sharing mechanism

5. **Memory Security**
   - Go garbage collector may leave memory traces
   - Sensitive data cleared but not guaranteed wiped

### Out of Scope

- ✗ Cloud synchronization
- ✗ Multi-user support
- ✗ Hardware security module (HSM) integration
- ✗ Biometric authentication
- ✗ Two-factor authentication for master password

## Best Practices

### Password Management

1. **Strong Master Password**
   ```
   ✅ Good: "correct-horse-battery-staple-29!" (33 chars)
   ✅ Good: "MyD0g!sN@med$potAnd1L0veH1m" (29 chars)
   ❌ Bad:  "password123" (11 chars, common)
   ❌ Bad:  "MyPassword1" (11 chars, predictable)
   ```

2. **Password Storage**
   - Write master password in password manager (ironic but practical)
   - Or write on paper, store in safe place
   - Don't store in plaintext file

3. **Password Rotation**
   - Change master password periodically
   - Rotate individual credentials regularly
   - Use `pass-cli generate` for new credentials

### Operational Security

1. **Vault Backups**
   ```bash
   # Regular backups
   cp ~/.pass-cli/vault.enc ~/backups/vault-$(date +%Y%m%d).enc

   # Store backups securely (encrypted drive, safe location)
   ```

2. **Clipboard Security**
   - Clipboard cleared automatically after 30 seconds
   - Avoid pasting into untrusted applications
   - Use `--no-clipboard` if concerned

3. **Script Security**
   ```bash
   # ✅ Good: Use quiet mode to avoid logging
   export API_KEY=$(pass-cli get service --quiet)

   # ❌ Bad: Full output might be logged
   export API_KEY=$(pass-cli get service)
   ```

4. **Audit Usage**
   ```bash
   # Review unused credentials monthly
   pass-cli list --unused --days 90

   # Delete obsolete credentials
   pass-cli delete old-service
   ```

### System Security

1. **Secure Your OS Account**
   - Use strong OS login password
   - Enable full-disk encryption
   - Keep system updated

2. **File System Security**
   - Don't commit vault to version control
   - Add to `.gitignore`:
     ```
     .pass-cli/
     *.enc
     ```

3. **Access Control**
   - Don't run Pass-CLI as root/admin
   - Use regular user account
   - Verify vault file permissions

### Development Security

1. **Testing**
   ```bash
   # Use separate vault for testing
   pass-cli --vault /tmp/test-vault.enc init

   # Clean up after testing
   rm -f /tmp/test-vault.enc
   ```

2. **Debugging**
   - Use `--verbose` flag, not hardcoded logging
   - Don't log credential values
   - Clear terminal after debugging

## Security Checklist

### Initial Setup
- [ ] Strong master password (20+ characters)
- [ ] Master password backed up securely
- [ ] Vault file permissions verified (0600)
- [ ] System keychain configured correctly

### Regular Maintenance
- [ ] Vault backed up monthly
- [ ] Unused credentials reviewed quarterly
- [ ] Master password rotated annually
- [ ] Pass-CLI updated to latest version

### Incident Response
- [ ] Master password changed if compromised
- [ ] Vault file restored from backup if corrupted
- [ ] All credentials rotated if vault possibly compromised
- [ ] System scan for malware if suspicious activity

## Incident Response

### Master Password Compromised

1. **Immediate Actions**
   - Change master password: `pass-cli init` (if you have access)
   - Or delete vault and start fresh
   - Rotate all credentials stored in vault

2. **Investigation**
   - Scan system for malware
   - Check keychain access logs (if available)
   - Review who had access to system

3. **Prevention**
   - Use stronger master password
   - Enable full-disk encryption
   - Review system security practices

### Vault File Corrupted

1. **Recovery**
   ```bash
   # Restore from backup
   cp ~/.pass-cli/vault.enc.backup ~/.pass-cli/vault.enc

   # Or from manual backup
   cp ~/backups/vault-20250120.enc ~/.pass-cli/vault.enc
   ```

2. **Verification**
   ```bash
   # Test vault access
   pass-cli list
   ```

3. **Prevention**
   - Regular backups
   - Atomic writes (built-in)
   - Don't manually edit vault file

### Credential Leaked

1. **Immediate**
   - Rotate credential immediately on actual service
   - Update in Pass-CLI: `pass-cli update service --generate`

2. **Investigation**
   - Identify leak source (logs, clipboard, screen share)
   - Review usage tracking: `pass-cli get service --json`

3. **Prevention**
   - Use `--quiet` mode in scripts
   - Clear shell history: `history -c`
   - Review script logging

## Security Audits

### Internal Audits

Run security checks regularly:

```bash
# Check vault permissions
ls -la ~/.pass-cli/

# Verify no plaintext secrets in code
grep -r "password.*=" .

# Run security scanner
make security-scan

# Check for vulnerable dependencies
go list -json -m all | nancy sleuth
```

### External Audits

Pass-CLI has not yet undergone external security audit. We welcome security researchers to review the code.

### Reporting Security Issues

**DO NOT** file public issues for security vulnerabilities.

Instead:
- Email: security@example.com
- Include: Detailed description, reproduction steps, impact assessment
- Response: Within 48 hours
- Disclosure: Coordinated disclosure after fix

### Security Updates

Security updates are released as:
- **Critical**: Immediate release, notification to users
- **High**: Release within 7 days
- **Medium**: Release in next version

Check for updates:
```bash
pass-cli version
# Compare with latest: https://github.com/ari1110/pass-cli/releases
```

## Cryptographic Algorithm Details

### AES-256-GCM Parameters

- **Block Size**: 128 bits
- **Key Size**: 256 bits
- **Nonce Size**: 96 bits (12 bytes) - NIST recommended
- **Tag Size**: 128 bits (16 bytes) - Full authentication
- **Additional Data**: None (not needed for our use case)

### PBKDF2 Parameters

- **Iteration Count**: 100,000
  - Provides ~0.1 second delay on modern hardware
  - Increases brute-force cost significantly
- **Salt Size**: 256 bits (32 bytes)
  - Unique per vault
  - Prevents rainbow table attacks
- **Hash Function**: SHA-256
  - NIST approved
  - 256-bit output matches key size

## Compliance and Standards

### Standards Compliance

- **NIST SP 800-38D**: AES-GCM mode
- **NIST SP 800-132**: PBKDF2 recommendations
- **NIST FIPS 197**: AES algorithm
- **RFC 5869**: PBKDF2 specification

### Best Practices Followed

- **OWASP**: Secure coding practices
- **CWE**: Common weakness mitigation
- **SANS**: Security implementation guidelines

## Further Reading

- [AES-GCM Specification (NIST SP 800-38D)](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf)
- [PBKDF2 Specification (RFC 2898)](https://www.rfc-editor.org/rfc/rfc2898)
- [Go Cryptography Documentation](https://pkg.go.dev/crypto)
- [OWASP Cryptographic Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cryptographic_Storage_Cheat_Sheet.html)

---

**Last Updated**: 2025-01-20
**Version**: v0.0.1
**Maintained By**: Pass-CLI Security Team
