#!/bin/bash
# ========================================
# Test Data Setup for tview TUI Testing
# ========================================
#
# This script creates a test vault with comprehensive
# test data for manual testing of the tview TUI implementation.
#
# Usage: ./test/setup-tview-test-data.sh
#

set -e  # Exit on error

echo ""
echo "========================================"
echo "tview TUI Test Data Setup"
echo "========================================"
echo ""

# Configuration
TEST_VAULT="test-vault-tview/vault.enc"
TEST_PASSWORD="test123456"
BINARY="./pass-cli"

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo "[Step 1/4] Building pass-cli..."
    go build -o "$BINARY" .
    if [ $? -ne 0 ]; then
        echo "ERROR: Build failed!"
        exit 1
    fi
    echo "Build successful!"
else
    echo "[Step 1/4] Binary exists, skipping build"
fi
echo ""

# Clean up old test vault if exists
if [ -d "test-vault-tview" ]; then
    echo "[Step 2/4] Cleaning up old test vault..."
    rm -rf test-vault-tview
    echo "Old vault removed"
else
    echo "[Step 2/4] No old vault to clean"
fi
echo ""

# Initialize test vault
echo "[Step 3/4] Initializing test vault..."
echo "Password: $TEST_PASSWORD"
echo ""
echo "$TEST_PASSWORD" | $BINARY --vault "$TEST_VAULT" init
if [ $? -ne 0 ]; then
    echo "ERROR: Vault initialization failed!"
    exit 1
fi
echo "Vault initialized successfully!"
echo ""

# Add comprehensive test credentials
echo "[Step 4/4] Adding test credentials..."
echo ""

# Helper function to add credential
add_cred() {
    local num=$1
    local service=$2
    local username=$3
    local password=$4
    local category=$5
    local notes=$6

    echo "[$num/15] Adding $service..."

    if [ -n "$category" ]; then
        echo "$TEST_PASSWORD" | $BINARY --vault "$TEST_VAULT" add "$service" -u "$username" -p "$password" -c "$category" -n "$notes" 2>/dev/null
    else
        echo "$TEST_PASSWORD" | $BINARY --vault "$TEST_VAULT" add "$service" -u "$username" -p "$password" -n "$notes" 2>/dev/null
    fi

    if [ $? -eq 0 ]; then
        echo "  ✓ $service added"
    else
        echo "  ✗ Failed to add $service"
    fi
}

# AWS Credentials (Cloud category)
add_cred 1 "aws-production" "admin" "AWSProd@2024!SecureKey" "Cloud" "Production AWS account"
add_cred 2 "aws-dev" "developer" "DevAWS#123Secure" "Cloud" "Development AWS account"

# GitHub Credentials (Version Control)
add_cred 3 "github-personal" "myusername" "GitHub@Personal2024!" "" "Personal GitHub account"
add_cred 4 "github-work" "work.user" "WorkGitHub#Secure99" "" "Work GitHub account"

# Database Credentials
add_cred 5 "postgres-main" "dbadmin" "PostgresMain@2024Secure!" "Databases" "Main PostgreSQL database"
add_cred 6 "mysql-dev" "devuser" "MySQLDev#SecurePass123" "Databases" "Development MySQL instance"

# API Keys
add_cred 7 "stripe-api" "sk_live_key" "sk_live_XXXXXXXXXXXXXXXXXXXXXXXX" "APIs" "Stripe production API key"
add_cred 8 "openai-api" "project-key" "sk-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" "AI Services" "OpenAI API key for project"

# Email/Communication
add_cred 9 "gmail-main" "user@gmail.com" "Gmail@SecurePass2024!" "Communication" "Main Gmail account"

# Payment Services
add_cred 10 "paypal-business" "business@example.com" "PayPal@Business2024Secure" "Payment" "Business PayPal account"

# Additional Cloud Services
add_cred 11 "azure-storage" "azureuser" "Azure@Storage2024Secure!" "Cloud" "Azure storage account"

# GitLab
add_cred 12 "gitlab-work" "gitlab.user" "GitLab@Work2024Secure" "" "Work GitLab instance"

# MongoDB
add_cred 13 "mongodb-cluster" "mongoadmin" "MongoDB@Cluster2024!" "Databases" "MongoDB production cluster"

# Test credential with special characters
add_cred 14 "special-chars-test" "user@domain.com" "P@\$\$w0rd!#^&*()_+-=[]{}|;:',.<>?/" "" "Testing special characters in password"

# Uncategorized credential
add_cred 15 "random-service" "randomuser" "RandomService@2024Secure" "" "Uncategorized test service"

echo ""
echo "========================================"
echo "Test Data Setup Complete!"
echo "========================================"
echo ""
echo "Vault Location: $TEST_VAULT"
echo "Master Password: $TEST_PASSWORD"
echo "Credentials Added: 15"
echo ""
echo "Categories Created:"
echo "  - Cloud (3 credentials)"
echo "  - Databases (3 credentials)"
echo "  - APIs (1 credential)"
echo "  - AI Services (1 credential)"
echo "  - Communication (1 credential)"
echo "  - Payment (1 credential)"
echo "  - Version Control (3 credentials)"
echo "  - Uncategorized (2 credentials)"
echo ""
echo "Ready for Testing!"
echo ""
echo "To launch the tview TUI:"
echo "  $BINARY tui --vault $TEST_VAULT"
echo ""
echo "Or use default vault location:"
echo "  $BINARY tui"
echo ""
echo "Use password: $TEST_PASSWORD"
echo ""
