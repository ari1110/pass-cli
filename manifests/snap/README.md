# Snap Package for Pass-CLI

GoReleaser is configured to build Snap packages automatically. This directory documents the Snap submission process.

## Snap Configuration

The Snap configuration is in `.goreleaser.yml` under the `snapcrafts` section.

**Key settings:**
- **Grade:** `stable` (production-ready)
- **Confinement:** `strict` (sandboxed for security)
- **Base:** `core22` (Ubuntu 22.04 LTS base)
- **Plugs:** `home`, `network`, `password-manager-service`

## Building Snap Package

The Snap package is built automatically by GoReleaser:

```bash
# Build with GoReleaser (includes Snap)
goreleaser release --clean

# Or build just the Snap in snapshot mode
goreleaser build --snapshot --clean --id pass-cli-snap
```

The `.snap` file will be in the `dist/` directory.

## Testing Locally

Before submitting to the Snap Store, test locally:

```bash
# Install from local file (requires --dangerous for unsigned snap)
sudo snap install ./dist/pass-cli_1.0.0_amd64.snap --dangerous

# Test the application
pass-cli version
pass-cli --help

# View snap info
snap info --verbose pass-cli

# Uninstall
sudo snap remove pass-cli
```

## Submission Process

### 1. Create Snapcraft Account

Visit https://snapcraft.io/ and create an account (can use GitHub to sign in).

### 2. Install Snapcraft Tools

```bash
sudo snap install snapcraft --classic
```

### 3. Login to Snapcraft

```bash
snapcraft login
```

### 4. Register the Snap Name

```bash
snapcraft register pass-cli
```

This reserves the name "pass-cli" for your account.

### 5. Upload the Snap

After building with GoReleaser:

```bash
# Upload to Snap Store (initially to edge channel)
snapcraft upload dist/pass-cli_1.0.0_amd64.snap --release=edge

# If all looks good, promote to stable
snapcraft release pass-cli 1.0.0 stable
```

### 6. Verify in Snap Store

Check your snap at: https://snapcraft.io/pass-cli

## Installation (After Publishing)

Once published, users can install via:

```bash
# Install from stable channel
sudo snap install pass-cli

# Or install from edge (development) channel
sudo snap install pass-cli --edge
```

## Channels

Snap uses channels for versioning:

- **stable** - Production releases (v1.0.0, v1.1.0, etc.)
- **candidate** - Release candidates
- **beta** - Beta releases
- **edge** - Latest builds from main branch

## Confinement and Permissions

Pass-CLI uses **strict confinement** with these plugs:

- **home** - Access to user's home directory for vault files
- **network** - Network access (if needed for future features)
- **password-manager-service** - Access to system keychain/secret service

Users may need to connect interfaces manually:

```bash
sudo snap connect pass-cli:password-manager-service
```

## Updating the Snap

For future releases:

1. Update version in GoReleaser
2. Build new release
3. Upload new version:
   ```bash
   snapcraft upload dist/pass-cli_1.1.0_amd64.snap --release=stable
   ```

## Architecture Support

The Snap build targets:
- **amd64** (x86_64) - Most common
- **arm64** (ARM 64-bit) - Raspberry Pi, ARM servers

GoReleaser builds both automatically.

## References

- [Snapcraft Documentation](https://snapcraft.io/docs)
- [GoReleaser Snap Documentation](https://goreleaser.com/customization/snapcraft/)
- [Snap Confinement](https://snapcraft.io/docs/snap-confinement)
- [Snap Store Dashboard](https://snapcraft.io/snaps)
