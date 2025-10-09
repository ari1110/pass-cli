# tview TUI Testing - Quick Start Guide

**Purpose**: Get started testing the tview TUI implementation in 5 minutes
**Audience**: QA testers, developers validating the implementation

---

## Prerequisites

- ✅ Go 1.25+ installed
- ✅ Git repository cloned: `pass-cli`
- ✅ Terminal emulator (Windows Terminal, iTerm2, gnome-terminal, etc.)
- ✅ ~10 minutes for testing

---

## Step 1: Build the Binary

```bash
# Navigate to project root
cd pass-cli

# Build binary
go build -o pass-cli.exe .    # Windows
# OR
go build -o pass-cli .        # macOS/Linux

# Verify build
./pass-cli version
```

**Expected**: Version information displays

---

## Step 2: Set Up Test Data

### Option A: Automated Setup (Recommended)

**Windows**:
```cmd
test\setup-tview-test-data.bat
```

**macOS/Linux**:
```bash
chmod +x test/setup-tview-test-data.sh
./test/setup-tview-test-data.sh
```

**What it does**:
- Creates test vault at `test-vault-tview/vault.enc`
- Password: `test123456`
- Adds 15 test credentials across multiple categories

### Option B: Manual Setup

```bash
# Initialize vault
./pass-cli --vault test-vault-tview/vault.enc init
# Enter password: test123456

# Add test credentials
echo "test123456" | ./pass-cli --vault test-vault-tview/vault.enc add aws-production -u admin -p "SecurePass123" -c Cloud
echo "test123456" | ./pass-cli --vault test-vault-tview/vault.enc add github-personal -u myuser -p "GitHubPass456"
# Add more as needed...
```

---

## Step 3: Launch the TUI

```bash
./pass-cli tui --vault test-vault-tview/vault.enc
```

**Enter password**: `test123456`

**Expected**: Dashboard displays with sidebar, table, and status bar

---

## Step 4: Execute Critical Tests

### Test 1: Form Input Protection (CRITICAL)

**This is the most important test - previous implementations failed this.**

1. Press `n` to open Add Credential form
2. Click into the "Service" field
3. **Type the letters**: `e`, `n`, `d`

**✅ PASS**: Letters appear in the input field
**❌ FAIL**: Shortcuts trigger (edit dialog opens, new form opens, delete dialog opens)

**Why Critical**: Previous implementations intercepted form input, making it impossible to type certain letters in credentials.

### Test 2: Navigation and Selection

1. Use **arrow keys** (↑↓) to navigate credentials in table
2. Press **Enter** to view credential details
3. Observe detail view updates on right panel
4. Press **Tab** to cycle focus between panels
5. Observe border colors change (active = bright cyan, inactive = gray)

**✅ PASS**: Smooth navigation, clear focus indication

### Test 3: CRUD Operations

**Add**:
1. Press `n`
2. Fill form: Service=`test-service`, Username=`testuser`, Password=`password123`
3. Press Enter or click Add
4. **✅ Expected**: Credential appears in table

**Edit**:
1. Select `test-service`
2. Press `e`
3. Change username to `modified-user`
4. Submit
5. **✅ Expected**: Table shows updated username

**Delete**:
1. Select `test-service`
2. Press `d`
3. Confirm deletion
4. **✅ Expected**: Credential removed from table

### Test 4: Responsive Layout

1. **Start large**: Resize terminal to 150 columns
   - **✅ Expected**: Sidebar + Table + Detail visible
2. **Go medium**: Resize to 100 columns
   - **✅ Expected**: Sidebar + Table (Detail hidden)
3. **Go small**: Resize to 70 columns
   - **✅ Expected**: Table only (Sidebar hidden)

### Test 5: Keyboard Shortcuts

Quick test of all major shortcuts:

- `q` - Quit ✅
- `Ctrl+C` - Force quit ✅
- `n` - New credential ✅
- `e` - Edit (with credential selected) ✅
- `d` - Delete (with credential selected) ✅
- `c` - Copy password ✅
- `p` - Toggle password visibility ✅
- `?` - Help screen ✅
- `Tab` - Cycle focus ✅

---

## Step 5: Document Results

Use the provided templates:

### Quick Test (5 minutes):
- **Checklist**: `docs/development/TVIEW_MANUAL_TESTING_CHECKLIST.md`
- Mark critical tests only

### Full Test (30-60 minutes):
- **Full Checklist**: Complete all sections
- **Test Report**: `docs/development/TVIEW_TEST_REPORT_TEMPLATE.md`
- **Expected Results**: `docs/development/TVIEW_EXPECTED_RESULTS.md` (reference)

---

## Common Issues and Troubleshooting

### Issue: Binary won't build

**Solution**:
```bash
# Check Go version
go version    # Should be 1.25+

# Clean and rebuild
go clean
go mod tidy
go build -o pass-cli.exe .
```

### Issue: Rounded borders show as `?` or boxes

**Cause**: Terminal doesn't support Unicode

**Solutions**:
- **Windows**: Use Windows Terminal (not CMD)
- **macOS**: Use iTerm2 or Terminal.app
- **Linux**: Use gnome-terminal or Alacritty

**Expected**: ASCII fallback (straight lines `+-|`) is acceptable

### Issue: Colors look wrong

**Cause**: Limited terminal color support

**Solutions**:
- Check terminal settings for true color/256-color support
- Modern terminals (Windows Terminal, iTerm2, Alacritty) support true color
- 256-color is acceptable fallback

### Issue: Password prompt doesn't appear

**Cause**: Keychain might have stored password

**Solution**:
```bash
# Windows - Delete keychain entry
cmdkey /delete:pass-cli

# macOS - Delete from Keychain Access
# Open Keychain Access app, search "pass-cli", delete entry

# Linux - Delete from Secret Service
secret-tool clear service pass-cli
```

### Issue: Can't type 'e' 'n' 'd' in forms

**THIS IS A BUG** - Report immediately

**Expected**: You MUST be able to type any character in form input fields

### Issue: Application hangs or freezes

**THIS IS A CRITICAL BUG** (deadlock)

**To report**:
1. Note exactly what you were doing
2. Terminal size
3. Platform and terminal emulator
4. Whether it's reproducible

**Workaround**: Ctrl+C to force quit

### Issue: Terminal corrupted after quit

**THIS IS A BUG** (terminal restoration failure)

**Temporary fix**:
```bash
reset    # Unix/macOS
# Or close and reopen terminal
```

---

## Test Data Reference

### Test Vault Credentials (15 total)

**Cloud (3)**:
- aws-production
- aws-dev
- azure-storage

**Databases (3)**:
- postgres-main
- mysql-dev
- mongodb-cluster

**APIs (1)**:
- stripe-api

**AI Services (1)**:
- openai-api

**Communication (1)**:
- gmail-main

**Payment (1)**:
- paypal-business

**Version Control (3)**:
- github-personal
- github-work
- gitlab-work

**Uncategorized (2)**:
- special-chars-test (tests special characters in password)
- random-service

---

## Testing Priorities

### Critical (MUST PASS for release):
1. ✅ Form input protection (can type 'e', 'n', 'd' in forms)
2. ✅ No deadlocks or hangs
3. ✅ Terminal restoration on quit
4. ✅ All CRUD operations work
5. ✅ Basic navigation (arrow keys, Tab, Enter)

### High Priority (should work):
1. ✅ Keyboard shortcuts
2. ✅ Responsive layout (all breakpoints)
3. ✅ Focus management
4. ✅ Visual styling (rounded borders, colors)
5. ✅ Error handling

### Nice to Have (document if broken):
1. ✅ Search/filter (if implemented)
2. ✅ Mouse support (if enabled)
3. ✅ Advanced features

---

## Quick Test Checklist (5 minutes)

Absolute minimum to verify basic functionality:

- [ ] Build succeeds
- [ ] TUI launches with password prompt
- [ ] Dashboard displays (sidebar, table, status bar visible)
- [ ] Can navigate with arrow keys
- [ ] Can add credential (press 'n', fill form, submit)
- [ ] **CRITICAL**: Can type 'e' 'n' 'd' in form fields
- [ ] Can edit credential (select, press 'e', modify, submit)
- [ ] Can delete credential (select, press 'd', confirm)
- [ ] Tab cycles focus with visual border change
- [ ] Quit with 'q' restores terminal
- [ ] No crashes or hangs

**If all pass**: Proceed with full testing
**If any fail**: Report critical issue

---

## Test Reports Location

Save your reports in:
```
docs/development/
├── TVIEW_MANUAL_TESTING_CHECKLIST.md    (Your working checklist)
├── TVIEW_TEST_REPORT_TEMPLATE.md        (Formal test report)
├── TVIEW_TEST_REPORT_[DATE]_[TESTER].md (Your completed report)
└── TVIEW_EXPECTED_RESULTS.md            (Reference guide)
```

---

## Getting Help

**Questions about expected behavior**: See `TVIEW_EXPECTED_RESULTS.md`

**Bug found**: Document in test report using bug template

**Need more test data**: Run setup script again (creates fresh vault)

**Stuck**:
1. Check troubleshooting section above
2. Try with clean vault: `rm -rf test-vault-tview/`
3. Rebuild binary: `go clean && go build`

---

## Next Steps After Testing

1. **Fill out test report** using template
2. **Document all bugs** with reproduction steps
3. **Take screenshots** of visual issues
4. **Note terminal compatibility** for each platform tested
5. **Submit report** to development team

---

**Happy Testing!** 🚀

Remember: **Form input protection** is the most critical test. If you can't type 'e', 'n', 'd' in form fields, that's a critical blocker.
