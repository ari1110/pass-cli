package security

import (
	"errors"
	"fmt"
	"unicode"
)

// T041 [US3]: PasswordPolicy struct defines password requirements
type PasswordPolicy struct {
	MinLength         int
	RequireUppercase  bool
	RequireLowercase  bool
	RequireDigit      bool
	RequireSymbol     bool
}

// T042 [US3]: DefaultPasswordPolicy constant (12 chars, all requirements true)
// FR-016: Minimum 12 characters with uppercase, lowercase, digit, and symbol
var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:         12,
	RequireUppercase:  true,
	RequireLowercase:  true,
	RequireDigit:      true,
	RequireSymbol:     true,
}

// PasswordStrength represents the strength level of a password
type PasswordStrength int

const (
	PasswordStrengthWeak PasswordStrength = iota
	PasswordStrengthMedium
	PasswordStrengthStrong
)

func (s PasswordStrength) String() string {
	switch s {
	case PasswordStrengthWeak:
		return "Weak"
	case PasswordStrengthMedium:
		return "Medium"
	case PasswordStrengthStrong:
		return "Strong"
	default:
		return "Unknown"
	}
}

// T043 [US3]: Validate method validates password against policy
// FR-016: Return descriptive error messages for each failed requirement
func (p *PasswordPolicy) Validate(password []byte) error {
	if password == nil {
		return fmt.Errorf("password cannot be empty (must be at least %d characters)", p.MinLength)
	}

	// Convert to rune slice for proper Unicode handling
	runes := []rune(string(password))

	// Check minimum length (count runes, not bytes)
	if len(runes) < p.MinLength {
		return fmt.Errorf("password must be at least %d characters long (got %d)", p.MinLength, len(runes))
	}

	// Track which requirements are met
	var hasUpper, hasLower, hasDigit, hasSymbol bool

	for _, r := range runes {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsLower(r) {
			hasLower = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			hasSymbol = true
		}
	}

	// Check requirements and return descriptive errors
	if p.RequireUppercase && !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if p.RequireLowercase && !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if p.RequireDigit && !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if p.RequireSymbol && !hasSymbol {
		return errors.New("password must contain at least one special character or symbol")
	}

	return nil
}

// T044 [US3]: Strength method calculates password strength
// FR-017: Calculate weak/medium/strong based on length and character variety
// Algorithm per data-model.md:186-238
func (p *PasswordPolicy) Strength(password []byte) PasswordStrength {
	if password == nil || len(password) == 0 {
		return PasswordStrengthWeak
	}

	// Convert to rune slice for proper Unicode handling
	runes := []rune(string(password))
	length := len(runes)

	// Calculate character variety score
	var hasUpper, hasLower, hasDigit, hasSymbol bool
	symbolCount := 0

	for _, r := range runes {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsLower(r) {
			hasLower = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			hasSymbol = true
			symbolCount++
		}
	}

	// Count character types present
	typeCount := 0
	if hasUpper {
		typeCount++
	}
	if hasLower {
		typeCount++
	}
	if hasDigit {
		typeCount++
	}
	if hasSymbol {
		typeCount++
	}

	// Strength calculation based on length and variety
	// Weak: < 12 characters OR missing required types
	// Medium: 12-19 characters with all required types
	// Strong: 20+ characters with all required types OR exceptional variety

	if length < 12 || typeCount < 4 {
		return PasswordStrengthWeak
	}

	if length >= 25 || (length >= 20 && symbolCount >= 3) {
		return PasswordStrengthStrong
	}

	if length >= 16 && typeCount == 4 {
		return PasswordStrengthMedium
	}

	return PasswordStrengthWeak
}
