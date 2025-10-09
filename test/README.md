# Testing Documentation

This directory contains end-to-end integration tests for Pass-CLI.

## Test Types

Pass-CLI uses two types of tests:

- **Unit Tests**: `*_test.go` files adjacent to source code in `cmd/` and `internal/` directories
  - Test individual functions and components in isolation
  - Fast execution, no external dependencies
  - Run during development for quick feedback

- **Integration Tests**: `*.go` files in this `test/` directory
  - Test complete workflows and end-to-end scenarios
  - Build actual binary and test real execution
  - Use build tag `//go:build integration` to separate from unit tests

## Running Tests

### Run All Tests
```bash
# Run both unit and integration tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
```

### Run Unit Tests Only
```bash
# Run unit tests in cmd/ and internal/
go test ./cmd/... ./internal/...
```

### Run Integration Tests Only
```bash
# Run integration tests in test/
go test -v -tags=integration ./test

# Or use make target
make test-integration
```

### Basic Integration Tests
```bash
# Run all integration tests
go test -v -tags=integration ./test

# Run with timeout (for slow systems)
go test -v -tags=integration ./test -timeout 5m

# Run specific test
go test -v -tags=integration ./test -run TestIntegration_CompleteWorkflow

# Run only keychain tests
go test -v -tags=integration ./test -run TestIntegration_Keychain

# Skip performance and stress tests (short mode)
go test -v -tags=integration ./test -short
```

### Make Targets
```bash
# Run integration tests
make test-integration

# Run all tests (unit + integration)
make test-all
```

## Test Coverage

### Complete Workflow Tests
- **TestIntegration_CompleteWorkflow**: Full E2E user journey
  - Init vault with master password
  - Add multiple credentials
  - List all credentials
  - Get specific credentials
  - Update credentials
  - Delete credentials
  - Generate passwords

### Error Handling Tests
- **TestIntegration_ErrorHandling**: Validates error scenarios
  - Wrong password rejection
  - Nonexistent credential handling
  - Duplicate vault initialization prevention

### Script-Friendly Output Tests
- **TestIntegration_ScriptFriendly**: CI/CD and automation support
  - Quiet mode (minimal output)
  - Field extraction (username, password, etc.)
  - No-clipboard mode

### Performance Tests
- **TestIntegration_Performance**: Validates performance targets
  - First unlock: <500ms target
  - Cached operations: <100ms target
  - Real-world timing measurements

### Stress Tests
- **TestIntegration_StressTest**: High-volume scenarios
  - Add 100 credentials
  - List performance with many entries
  - Random access patterns

### Version Test
- **TestIntegration_Version**: Basic version command validation

### Keychain Integration Tests
- **TestIntegration_KeychainWorkflow**: End-to-end keychain integration
  - Init vault with `--use-keychain` flag
  - Verify password stored in system keychain
  - Auto-unlock for add/get/list/update/delete commands (no password prompt)
  - Full workflow validation with native OS keychain
- **TestIntegration_KeychainFallback**: Keychain fallback behavior
  - Graceful degradation when keychain entry is deleted
  - Password prompt fallback when keychain unavailable
- **TestIntegration_KeychainUnavailable**: Unavailable keychain handling
  - Behavior when system keychain is not accessible
  - Error messages and warnings
- **TestIntegration_MultipleVaultsKeychain**: Multiple vault scenarios
  - Shared keychain entry behavior
  - Cross-vault operations
- **TestIntegration_KeychainVerboseOutput**: Verbose mode keychain feedback
  - Verification of keychain usage messages

## Test Architecture

### Build Tags
Tests use `//go:build integration` to separate from unit tests. This allows:
- Faster unit test runs during development
- Targeted integration test execution in CI/CD
- Resource-intensive tests only when needed

### Test Isolation
Each test:
- Creates temporary vault directories
- Uses unique vault paths
- Cleans up after completion
- Cleans up keychain entries after keychain tests

### Binary Building
`TestMain` automatically:
- Builds the `pass-cli` binary before tests
- Cleans up the binary after tests
- Sets up temporary test directories

## Performance Targets

Based on the spec requirements:
- **Cold start (first unlock)**: < 500ms ✅ (achieving ~95ms)
- **Cached operations**: < 100ms ✅ (achieving ~95ms)
- **Stress test**: Handle 100+ credentials efficiently ✅

## Notes

### Update Command
The update command test currently skips as it needs verification of the implementation. The test is designed to validate that updates persist correctly.

### Quiet Mode
Quiet mode outputs are logged for verification. The implementation appears to work but may need refinement for truly "script-friendly" output (just the value, no formatting).

### Cross-Platform
Tests are designed to run on Windows, macOS, and Linux. Platform-specific features (like keychain integration) are handled gracefully.

### Keychain Tests
Keychain integration tests interact with real OS keychains:
- **Windows**: Windows Credential Manager
- **macOS**: Keychain Access
- **Linux**: Secret Service (D-Bus)

**Important Notes:**
- Tests automatically skip if system keychain is unavailable
- Tests clean up keychain entries in `defer` blocks
- Safe to run locally - won't interfere with other apps
- On CI/CD, keychain may not be available (tests will skip gracefully)
- Use the keychain service name "pass-cli" and account "master-password"

## Test Utilities

### test-tui.bat
Manual TUI testing script for interactive testing (Windows batch script) - **LEGACY Bubble Tea TUI**

**Location**: `test/test-tui.bat`

**Usage** (run from project root):
```bash
# Windows
test\test-tui.bat

# The script will:
# 1. Build the pass-cli binary
# 2. Initialize a test vault (test-tui-vault/vault.enc)
# 3. Add sample credentials (github.com, gitlab.com, aws.com)
# 4. Provide instructions to launch the TUI manually
```

**Purpose**: Quickly set up a test environment for manual TUI testing and exploration (legacy Bubble Tea implementation).

### setup-tview-test-data.bat / .sh
**NEW**: Comprehensive test data setup for tview TUI implementation manual testing

**Location**: `test/setup-tview-test-data.bat` (Windows), `test/setup-tview-test-data.sh` (Unix)

**Usage** (run from project root):
```bash
# Windows
test\setup-tview-test-data.bat

# macOS/Linux (make executable first)
chmod +x test/setup-tview-test-data.sh
./test/setup-tview-test-data.sh

# The script will:
# 1. Build the pass-cli binary (if needed)
# 2. Initialize test vault (test-vault-tview/vault.enc)
# 3. Add 15 comprehensive test credentials across 8 categories
# 4. Provide launch instructions
```

**Test Data Created**:
- **Vault**: `test-vault-tview/vault.enc`
- **Password**: `test123456`
- **Credentials**: 15 credentials across Cloud, Databases, APIs, AI Services, Communication, Payment, Version Control, and Uncategorized categories
- **Special Test Cases**: Long names, special characters, Unicode support testing

**Purpose**: Create comprehensive test data for validating tview TUI implementation against all requirements (Task 17 of tui-tview-implementation spec).

**Documentation**:
- **Quick Start**: `docs/development/TVIEW_TESTING_QUICKSTART.md` - 5-minute setup guide
- **Full Checklist**: `docs/development/TVIEW_MANUAL_TESTING_CHECKLIST.md` - Complete testing checklist
- **Expected Results**: `docs/development/TVIEW_EXPECTED_RESULTS.md` - Detailed expected behavior
- **Test Report Template**: `docs/development/TVIEW_TEST_REPORT_TEMPLATE.md` - Formal test report

## Test Data

### test-vault/
Integration test fixture directory containing encrypted vault for testing.

**Location**: `test-vault/`

**Contents**:
- `vault.enc` - Pre-encrypted test vault with known password
- Used by integration tests to validate vault operations
- Contains sample credentials for testing purposes

**Security Note**: This is a TEST vault only. Never use for real credentials.

## CI/CD Integration

```yaml
# Example GitHub Actions workflow
- name: Run Integration Tests
  run: go test -v -tags=integration ./test -timeout 5m

# With coverage
- name: Run Integration Tests with Coverage
  run: go test -v -tags=integration -coverprofile=integration-coverage.out ./test

# Note: Keychain tests will automatically skip in CI environments
# where system keychain is unavailable. To run keychain tests in CI:
# - macOS: Use macOS runners (keychain available by default)
# - Linux: Install and configure gnome-keyring or similar
# - Windows: Use Windows runners (Credential Manager available)
```

## Future Enhancements

Potential additions for comprehensive testing:
- Per-vault keychain entries (currently uses single shared entry)
- Concurrent access tests (multiple processes)
- Backup/restore workflow tests
- Import/export functionality tests (if implemented)
- Migration tests for vault format changes
- Keychain permission tests (verify proper OS-level isolation)