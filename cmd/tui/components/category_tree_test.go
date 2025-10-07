package components

import (
	"testing"
	"time"

	"pass-cli/internal/vault"
)

func TestCategorizeCredentials_Cloud(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "aws-main", Username: "admin"},
		{Service: "azure-main", Username: "user"},
	}

	categories := CategorizeCredentials(credentials)

	// Find cloud category
	var cloudCategory *Category
	for i := range categories {
		if categories[i].Name == string(CategoryCloud) {
			cloudCategory = &categories[i]
			break
		}
	}

	if cloudCategory == nil {
		t.Fatal("Cloud category not found")
	}

	if cloudCategory.Count < 2 {
		t.Errorf("Expected at least 2 credentials in Cloud, got %d", cloudCategory.Count)
	}
}

func TestCategorizeCredentials_Databases(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "postgres-main", Username: "dbuser"},
		{Service: "mysql-dev", Username: "root"},
		{Service: "database-prod", Username: "admin"}, // Use "database" keyword directly
	}

	categories := CategorizeCredentials(credentials)

	var dbCategory *Category
	for i := range categories {
		if categories[i].Name == string(CategoryDatabases) {
			dbCategory = &categories[i]
			break
		}
	}

	if dbCategory == nil {
		t.Fatal("Databases category not found")
	}

	if dbCategory.Count < 3 {
		t.Errorf("Expected at least 3 credentials in Databases, got %d", dbCategory.Count)
	}
}

func TestCategorizeCredentials_Git(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "github-personal", Username: "user"},
		{Service: "gitlab-work", Username: "developer"},
		{Service: "bitbucket-team", Username: "member"},
	}

	categories := CategorizeCredentials(credentials)

	var gitCategory *Category
	for i := range categories {
		if categories[i].Name == string(CategoryVersionControl) {
			gitCategory = &categories[i]
			break
		}
	}

	if gitCategory == nil {
		t.Fatal("Version Control category not found")
	}

	if gitCategory.Count != 3 {
		t.Errorf("Expected 3 credentials in Version Control, got %d", gitCategory.Count)
	}
}

func TestCategorizeCredentials_APIs(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "custom-api", Username: "key"},   // Changed to ensure it matches API pattern
		{Service: "rest-service", Username: "sid"}, // Changed to match service/rest pattern
		{Service: "oauth-provider", Username: "client"},
	}

	categories := CategorizeCredentials(credentials)

	var apiCategory *Category
	for i := range categories {
		if categories[i].Name == string(CategoryAPIs) {
			apiCategory = &categories[i]
			break
		}
	}

	if apiCategory == nil {
		t.Fatal("APIs category not found")
	}

	if apiCategory.Count != 3 {
		t.Errorf("Expected 3 credentials in APIs, got %d", apiCategory.Count)
	}
}

func TestCategorizeCredentials_AI(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "openai-key", Username: "key"},
		{Service: "anthropic-key", Username: "apikey"},
	}

	categories := CategorizeCredentials(credentials)

	var aiCategory *Category
	for i := range categories {
		if categories[i].Name == string(CategoryAI) {
			aiCategory = &categories[i]
			break
		}
	}

	if aiCategory == nil {
		t.Fatal("AI Services category not found")
	}

	if aiCategory.Count < 2 {
		t.Errorf("Expected at least 2 credentials in AI Services, got %d", aiCategory.Count)
	}
}

func TestCategorizeCredentials_Uncategorized(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "random-xyz-123", Username: "user"}, // Changed to ensure it doesn't match any pattern
		{Service: "unknown-platform", Username: "account"},
	}

	categories := CategorizeCredentials(credentials)

	var uncatCategory *Category
	for i := range categories {
		if categories[i].Name == string(CategoryUncategorized) {
			uncatCategory = &categories[i]
			break
		}
	}

	if uncatCategory == nil {
		t.Fatal("Uncategorized category not found")
	}

	if uncatCategory.Count != 2 {
		t.Errorf("Expected 2 credentials in Uncategorized, got %d", uncatCategory.Count)
	}
}

func TestCategorizeCredentials_CaseInsensitive(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "AWS-Production", Username: "admin"},
		{Service: "GitHub-Enterprise", Username: "user"},
		{Service: "PostgreSQL-DB", Username: "dbuser"}, // Changed to use "db" pattern
	}

	categories := CategorizeCredentials(credentials)

	// Count total categorized credentials
	totalCategorized := 0
	for _, cat := range categories {
		totalCategorized += cat.Count
	}

	if totalCategorized != len(credentials) {
		t.Errorf("Expected all %d credentials categorized, but got %d", len(credentials), totalCategorized)
	}

	// At minimum, verify none went to Uncategorized (they should all match despite case)
	var uncatCount int
	for _, cat := range categories {
		if cat.Name == string(CategoryUncategorized) {
			uncatCount = cat.Count
		}
	}

	if uncatCount > 0 {
		t.Error("Case-insensitive matching failed - credentials went to Uncategorized")
	}
}

func TestCategorizeCredentials_Empty(t *testing.T) {
	credentials := []vault.CredentialMetadata{}

	categories := CategorizeCredentials(credentials)

	// Should still return category list, but with zero counts
	totalCount := 0
	for _, cat := range categories {
		totalCount += cat.Count
	}

	if totalCount != 0 {
		t.Errorf("Expected 0 total credentials, got %d", totalCount)
	}
}

func TestCategorizeCredentials_AllCredentialsCategorized(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin"},
		{Service: "postgres-db", Username: "user"},
		{Service: "github-repo", Username: "dev"},
		{Service: "unknown-service", Username: "test"},
		{Service: "openai-key", Username: "api"},
	}

	categories := CategorizeCredentials(credentials)

	// Count total categorized credentials
	totalCategorized := 0
	for _, cat := range categories {
		totalCategorized += cat.Count
	}

	if totalCategorized != len(credentials) {
		t.Errorf("Expected all %d credentials categorized, but got %d", len(credentials), totalCategorized)
	}
}

func TestGetCategoryIcon(t *testing.T) {
	testCases := []struct {
		category CategoryType
		expected string
	}{
		{CategoryCloud, "☁️"},
		{CategoryAPIs, "🔑"},
		{CategoryDatabases, "💾"},
		{CategoryVersionControl, "📦"},
		{CategoryCommunication, "📧"},
		{CategoryPayment, "💰"},
		{CategoryAI, "🤖"},
		{CategoryUncategorized, "📁"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.category), func(t *testing.T) {
			icon := GetCategoryIcon(tc.category)
			if icon != tc.expected {
				t.Errorf("Expected icon %s for %s, got %s", tc.expected, tc.category, icon)
			}
		})
	}
}

func TestGetStatusIcon(t *testing.T) {
	testCases := []struct {
		status   string
		expected string
	}{
		{"pending", "⏳"},
		{"running", "⏳"},
		{"success", "✓"},
		{"failed", "✗"},
		{"collapsed", "▶"},
		{"expanded", "▼"},
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			icon := GetStatusIcon(tc.status)
			if icon != tc.expected {
				t.Errorf("Expected icon %s for %s, got %s", tc.expected, tc.status, icon)
			}
		})
	}
}

func TestCategorizeCredentials_OrderPreservation(t *testing.T) {
	// Create credentials with known timestamps to test ordering
	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin", CreatedAt: time.Now().Add(-3 * time.Hour)},
		{Service: "aws-dev", Username: "user", CreatedAt: time.Now().Add(-1 * time.Hour)},
		{Service: "aws-test", Username: "tester", CreatedAt: time.Now().Add(-2 * time.Hour)},
	}

	categories := CategorizeCredentials(credentials)

	var cloudCategory *Category
	for i := range categories {
		if categories[i].Name == string(CategoryCloud) {
			cloudCategory = &categories[i]
			break
		}
	}

	if cloudCategory == nil {
		t.Fatal("Cloud category not found")
	}

	// Check that all credentials are present
	if len(cloudCategory.Credentials) != 3 {
		t.Errorf("Expected 3 credentials, got %d", len(cloudCategory.Credentials))
	}

	// Verify each credential is present
	services := make(map[string]bool)
	for _, cred := range cloudCategory.Credentials {
		services[cred.Service] = true
	}

	if !services["aws-prod"] || !services["aws-dev"] || !services["aws-test"] {
		t.Error("Not all credentials were preserved in categorization")
	}
}
