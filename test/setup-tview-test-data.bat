@echo off
REM ========================================
REM Test Data Setup for tview TUI Testing
REM ========================================
REM
REM This script creates a test vault with comprehensive
REM test data for manual testing of the tview TUI implementation.
REM
REM Usage: test\setup-tview-test-data.bat
REM

echo.
echo ========================================
echo tview TUI Test Data Setup
echo ========================================
echo.

REM Configuration
set TEST_VAULT=test-vault-tview\vault.enc
set TEST_PASSWORD=test123456
set BINARY=pass-cli.exe

REM Check if binary exists
if not exist %BINARY% (
    echo [Step 1/4] Building pass-cli...
    go build -o %BINARY% .
    if %errorlevel% neq 0 (
        echo ERROR: Build failed!
        exit /b 1
    )
    echo Build successful!
) else (
    echo [Step 1/4] Binary exists, skipping build
)
echo.

REM Clean up old test vault if exists
if exist test-vault-tview (
    echo [Step 2/4] Cleaning up old test vault...
    rmdir /s /q test-vault-tview 2>nul
    echo Old vault removed
) else (
    echo [Step 2/4] No old vault to clean
)
echo.

REM Initialize test vault
echo [Step 3/4] Initializing test vault...
echo Password: %TEST_PASSWORD%
echo.
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% init
if %errorlevel% neq 0 (
    echo ERROR: Vault initialization failed!
    exit /b 1
)
echo Vault initialized successfully!
echo.

REM Add comprehensive test credentials
echo [Step 4/4] Adding test credentials...
echo.

REM AWS Credentials (Cloud category)
echo [1/15] Adding aws-production...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add aws-production -u admin -p "AWSProd@2024!SecureKey" -c Cloud -n "Production AWS account" 2>nul
if %errorlevel% equ 0 (echo   ✓ aws-production added) else (echo   ✗ Failed to add aws-production)

echo [2/15] Adding aws-dev...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add aws-dev -u developer -p "DevAWS#123Secure" -c Cloud -n "Development AWS account" 2>nul
if %errorlevel% equ 0 (echo   ✓ aws-dev added) else (echo   ✗ Failed to add aws-dev)

REM GitHub Credentials (Version Control)
echo [3/15] Adding github-personal...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add github-personal -u myusername -p "GitHub@Personal2024!" -n "Personal GitHub account" 2>nul
if %errorlevel% equ 0 (echo   ✓ github-personal added) else (echo   ✗ Failed to add github-personal)

echo [4/15] Adding github-work...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add github-work -u work.user -p "WorkGitHub#Secure99" -n "Work GitHub account" 2>nul
if %errorlevel% equ 0 (echo   ✓ github-work added) else (echo   ✗ Failed to add github-work)

REM Database Credentials
echo [5/15] Adding postgres-main...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add postgres-main -u dbadmin -p "PostgresMain@2024Secure!" -c Databases -n "Main PostgreSQL database" 2>nul
if %errorlevel% equ 0 (echo   ✓ postgres-main added) else (echo   ✗ Failed to add postgres-main)

echo [6/15] Adding mysql-dev...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add mysql-dev -u devuser -p "MySQLDev#SecurePass123" -c Databases -n "Development MySQL instance" 2>nul
if %errorlevel% equ 0 (echo   ✓ mysql-dev added) else (echo   ✗ Failed to add mysql-dev)

REM API Keys
echo [7/15] Adding stripe-api...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add stripe-api -u sk_live_key -p "sk_live_XXXXXXXXXXXXXXXXXXXXXXXX" -c APIs -n "Stripe production API key" 2>nul
if %errorlevel% equ 0 (echo   ✓ stripe-api added) else (echo   ✗ Failed to add stripe-api)

echo [8/15] Adding openai-api...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add openai-api -u project-key -p "sk-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" -c "AI Services" -n "OpenAI API key for project" 2>nul
if %errorlevel% equ 0 (echo   ✓ openai-api added) else (echo   ✗ Failed to add openai-api)

REM Email/Communication
echo [9/15] Adding gmail-main...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add gmail-main -u user@gmail.com -p "Gmail@SecurePass2024!" -c Communication -n "Main Gmail account" 2>nul
if %errorlevel% equ 0 (echo   ✓ gmail-main added) else (echo   ✗ Failed to add gmail-main)

REM Payment Services
echo [10/15] Adding paypal-business...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add paypal-business -u business@example.com -p "PayPal@Business2024Secure" -c Payment -n "Business PayPal account" 2>nul
if %errorlevel% equ 0 (echo   ✓ paypal-business added) else (echo   ✗ Failed to add paypal-business)

REM Additional Cloud Services
echo [11/15] Adding azure-storage...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add azure-storage -u azureuser -p "Azure@Storage2024Secure!" -c Cloud -n "Azure storage account" 2>nul
if %errorlevel% equ 0 (echo   ✓ azure-storage added) else (echo   ✗ Failed to add azure-storage)

REM GitLab
echo [12/15] Adding gitlab-work...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add gitlab-work -u gitlab.user -p "GitLab@Work2024Secure" -n "Work GitLab instance" 2>nul
if %errorlevel% equ 0 (echo   ✓ gitlab-work added) else (echo   ✗ Failed to add gitlab-work)

REM MongoDB
echo [13/15] Adding mongodb-cluster...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add mongodb-cluster -u mongoadmin -p "MongoDB@Cluster2024!" -c Databases -n "MongoDB production cluster" 2>nul
if %errorlevel% equ 0 (echo   ✓ mongodb-cluster added) else (echo   ✗ Failed to add mongodb-cluster)

REM Test credential with special characters
echo [14/15] Adding special-chars-test...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add "special-chars-test" -u "user@domain.com" -p "P@$$w0rd!#^&*()_+-=[]{}|;:',.<>?/" -n "Testing special characters in password" 2>nul
if %errorlevel% equ 0 (echo   ✓ special-chars-test added) else (echo   ✗ Failed to add special-chars-test)

REM Uncategorized credential
echo [15/15] Adding random-service...
echo %TEST_PASSWORD% | %BINARY% --vault %TEST_VAULT% add random-service -u randomuser -p "RandomService@2024Secure" -n "Uncategorized test service" 2>nul
if %errorlevel% equ 0 (echo   ✓ random-service added) else (echo   ✗ Failed to add random-service)

echo.
echo ========================================
echo Test Data Setup Complete!
echo ========================================
echo.
echo Vault Location: %TEST_VAULT%
echo Master Password: %TEST_PASSWORD%
echo Credentials Added: 15
echo.
echo Categories Created:
echo   - Cloud (3 credentials)
echo   - Databases (3 credentials)
echo   - APIs (1 credential)
echo   - AI Services (1 credential)
echo   - Communication (1 credential)
echo   - Payment (1 credential)
echo   - Version Control (3 credentials)
echo   - Uncategorized (2 credentials)
echo.
echo Ready for Testing!
echo.
echo To launch the tview TUI:
echo   %BINARY% tui --vault %TEST_VAULT%
echo.
echo Or use default vault location:
echo   %BINARY% tui
echo.
echo Use password: %TEST_PASSWORD%
echo.
pause
