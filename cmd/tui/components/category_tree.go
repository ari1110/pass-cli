package components

import (
	"strings"

	"pass-cli/internal/vault"
)

// Category represents a credential category with its contents
type Category struct {
	Name        string
	Icon        string
	Count       int
	Expanded    bool
	Credentials []vault.CredentialMetadata
}

// CategoryType represents the type of category
type CategoryType string

const (
	CategoryAPIs           CategoryType = "APIs & Services"
	CategoryCloud          CategoryType = "Cloud Infrastructure"
	CategoryDatabases      CategoryType = "Databases"
	CategoryVersionControl CategoryType = "Version Control"
	CategoryCommunication  CategoryType = "Communication"
	CategoryPayment        CategoryType = "Payment Processing"
	CategoryAI             CategoryType = "AI Services"
	CategoryUncategorized  CategoryType = "Uncategorized"
)

// categoryPatterns maps category types to service name patterns
var categoryPatterns = map[CategoryType][]string{
	CategoryAPIs: {
		"api", "rest", "graphql", "webhook", "endpoint",
		"service", "oauth", "auth0", "jwt", "token",
	},
	CategoryCloud: {
		"aws", "amazon", "s3", "ec2", "lambda",
		"azure", "gcp", "google cloud", "digitalocean", "heroku",
		"vercel", "netlify", "cloudflare", "linode",
	},
	CategoryDatabases: {
		"postgres", "postgresql", "mysql", "mariadb", "mongodb",
		"redis", "elasticsearch", "database", "db", "sql",
		"cassandra", "dynamodb", "sqlite", "oracle",
	},
	CategoryVersionControl: {
		"github", "gitlab", "bitbucket", "git",
		"sourceforge", "mercurial", "svn",
	},
	CategoryCommunication: {
		"slack", "discord", "telegram", "whatsapp",
		"email", "smtp", "mail", "sendgrid", "mailgun",
		"twilio", "zoom", "teams", "skype",
	},
	CategoryPayment: {
		"stripe", "paypal", "square", "braintree",
		"payment", "billing", "merchant", "checkout",
	},
	CategoryAI: {
		"openai", "anthropic", "claude", "gpt",
		"gemini", "ai", "ml", "machine learning",
		"huggingface", "replicate",
	},
}

// CategorizeCredentials organizes credentials into categories based on service name patterns
func CategorizeCredentials(credentials []vault.CredentialMetadata) []Category {
	// Initialize all categories
	categoryMap := make(map[CategoryType]*Category)
	for _, catType := range []CategoryType{
		CategoryAPIs,
		CategoryCloud,
		CategoryDatabases,
		CategoryVersionControl,
		CategoryCommunication,
		CategoryPayment,
		CategoryAI,
		CategoryUncategorized,
	} {
		categoryMap[catType] = &Category{
			Name:        string(catType),
			Icon:        GetCategoryIcon(catType),
			Expanded:    false,
			Credentials: []vault.CredentialMetadata{},
		}
	}

	// Categorize each credential
	for _, cred := range credentials {
		category := categorizeCredential(cred)
		categoryMap[category].Credentials = append(categoryMap[category].Credentials, cred)
	}

	// Build result slice with non-empty categories first, then uncategorized
	result := []Category{}
	for _, catType := range []CategoryType{
		CategoryCloud,
		CategoryDatabases,
		CategoryVersionControl,
		CategoryAPIs,
		CategoryAI,
		CategoryPayment,
		CategoryCommunication,
	} {
		cat := categoryMap[catType]
		cat.Count = len(cat.Credentials)
		if cat.Count > 0 {
			result = append(result, *cat)
		}
	}

	// Add uncategorized last if it has items
	uncategorized := categoryMap[CategoryUncategorized]
	uncategorized.Count = len(uncategorized.Credentials)
	if uncategorized.Count > 0 {
		result = append(result, *uncategorized)
	}

	return result
}

// categorizeCredential determines which category a credential belongs to
func categorizeCredential(cred vault.CredentialMetadata) CategoryType {
	serviceLower := strings.ToLower(cred.Service)

	// Check categories in specific order to avoid randomness from map iteration
	// More specific categories should be checked first
	categoryOrder := []CategoryType{
		CategoryCloud,
		CategoryDatabases,
		CategoryVersionControl,
		CategoryAI,
		CategoryPayment,
		CategoryCommunication,
		CategoryAPIs, // Check this last as it has generic patterns like "service"
	}

	for _, category := range categoryOrder {
		patterns := categoryPatterns[category]
		for _, pattern := range patterns {
			if strings.Contains(serviceLower, strings.ToLower(pattern)) {
				return category
			}
		}
	}

	return CategoryUncategorized
}

// GetCategoryIcon returns the icon for a given category type
func GetCategoryIcon(category CategoryType) string {
	icons := map[CategoryType]string{
		CategoryAPIs:           "🔑",
		CategoryCloud:          "☁️",
		CategoryDatabases:      "💾",
		CategoryVersionControl: "📦",
		CategoryCommunication:  "📧",
		CategoryPayment:        "💰",
		CategoryAI:             "🤖",
		CategoryUncategorized:  "📁",
	}

	if icon, ok := icons[category]; ok {
		return icon
	}
	return "📁"
}

// GetStatusIcon returns the icon for a given status
func GetStatusIcon(status string) string {
	icons := map[string]string{
		"pending":   "⏳",
		"running":   "⏳",
		"success":   "✓",
		"failed":    "✗",
		"collapsed": "▶",
		"expanded":  "▼",
	}

	if icon, ok := icons[status]; ok {
		return icon
	}
	return ""
}

// GetCategoryIconASCII returns ASCII fallback for category icons
func GetCategoryIconASCII(category CategoryType) string {
	icons := map[CategoryType]string{
		CategoryAPIs:           "[API]",
		CategoryCloud:          "[CLD]",
		CategoryDatabases:      "[DB]",
		CategoryVersionControl: "[GIT]",
		CategoryCommunication:  "[MSG]",
		CategoryPayment:        "[PAY]",
		CategoryAI:             "[AI]",
		CategoryUncategorized:  "[???]",
	}

	if icon, ok := icons[category]; ok {
		return icon
	}
	return "[???]"
}

// GetStatusIconASCII returns ASCII fallback for status icons
func GetStatusIconASCII(status string) string {
	icons := map[string]string{
		"pending":   "[.]",
		"running":   "[.]",
		"success":   "[+]",
		"failed":    "[X]",
		"collapsed": "[>]",
		"expanded":  "[v]",
	}

	if icon, ok := icons[status]; ok {
		return icon
	}
	return ""
}
