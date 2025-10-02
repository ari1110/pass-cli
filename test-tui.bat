@echo off
echo ========================================
echo Testing pass-cli TUI
echo ========================================
echo.

REM Set up test vault path
set TEST_VAULT=test-tui-vault\vault.enc

echo Step 1: Building pass-cli...
go build -o pass-cli.exe .
if %errorlevel% neq 0 (
    echo Build failed!
    exit /b 1
)
echo Build successful!
echo.

echo Step 2: Initializing test vault...
echo You'll be prompted for a password. Use: test123456
echo.
pass-cli.exe --vault %TEST_VAULT% init
if %errorlevel% neq 0 (
    echo Init failed!
    exit /b 1
)
echo.

echo Step 3: Adding sample credentials...
echo Adding github.com...
echo test123456 | pass-cli.exe --vault %TEST_VAULT% add github.com --username myuser --password mypass123
echo.
echo Adding gitlab.com...
echo test123456 | pass-cli.exe --vault %TEST_VAULT% add gitlab.com --username devuser --password devpass456
echo.
echo Adding aws.com...
echo test123456 | pass-cli.exe --vault %TEST_VAULT% add aws.com --username admin --password awspass789
echo.

echo ========================================
echo Test vault ready! Now launching TUI...
echo ========================================
echo.
echo Press Ctrl+C to stop this script, then run:
echo   pass-cli.exe --vault %TEST_VAULT%
echo.
echo Or just run: pass-cli.exe (uses default vault)
echo.
pause
