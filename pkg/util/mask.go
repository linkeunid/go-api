package util

import (
	"regexp"
	"strings"
)

// MaskDsn replaces password in DSN with asterisks
func MaskDsn(dsn string) string {
	// Match password pattern in MySQL DSN: user:password@tcp(...)
	re := regexp.MustCompile(`([^:]+):([^@]+)@`)
	return re.ReplaceAllString(dsn, "$1:******@")
}

// MaskSensitive masks a sensitive string value with asterisks,
// keeping a few characters at the beginning and end visible
func MaskSensitive(value string, visiblePrefixChars, visibleSuffixChars int) string {
	if value == "" {
		return ""
	}

	length := len(value)

	// For very short strings, just return all asterisks
	if length <= visiblePrefixChars+visibleSuffixChars {
		return strings.Repeat("*", length)
	}

	// Extract visible parts
	prefix := value[:visiblePrefixChars]
	suffix := ""
	if visibleSuffixChars > 0 {
		suffix = value[length-visibleSuffixChars:]
	}

	// Calculate number of asterisks needed for mask
	maskLength := length - visiblePrefixChars - visibleSuffixChars

	// Generate mask with exact length
	mask := strings.Repeat("*", maskLength)

	return prefix + mask + suffix
}

// MaskCredential masks a credential string, showing only the first 2 and last 2 characters
func MaskCredential(credential string) string {
	return MaskSensitive(credential, 2, 2)
}

// MaskEmail masks the local part of an email address, showing only the first 2 characters
// Example: jo******@example.com
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return MaskSensitive(email, 2, 2) // Not a valid email format, mask it differently
	}

	local := parts[0]
	domain := parts[1]

	// Special case for very short local parts
	if len(local) <= 2 {
		return email // Don't mask very short local parts
	}

	maskedLocal := local[:2] + strings.Repeat("*", len(local)-2)
	return maskedLocal + "@" + domain
}

// MaskJWT masks a JWT token, revealing only signature type and a few chars
// Example: eyJh******.eyJs******.SflKx******
func MaskJWT(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return MaskSensitive(token, 4, 3) // Not a valid JWT format, mask it differently
	}

	// Mask each part of the JWT separately
	maskedHeader := parts[0][:4] + strings.Repeat("*", 6)
	maskedPayload := parts[1][:4] + strings.Repeat("*", 6)
	maskedSignature := parts[2][:5] + strings.Repeat("*", 6)

	return maskedHeader + "." + maskedPayload + "." + maskedSignature
}

// MaskURL masks sensitive parts of a URL like username, password, and access tokens
// Example: https://user:****@example.com/api?token=****
func MaskURL(url string) string {
	// Mask basic auth credentials
	credRegex := regexp.MustCompile(`(https?:\/\/)([^:]+):([^@]+)@`)
	masked := credRegex.ReplaceAllString(url, "${1}${2}:******@")

	// Mask tokens, api keys, etc. in query params
	sensitiveParams := []string{"token", "key", "secret", "password", "access_token", "api_key"}

	for _, param := range sensitiveParams {
		paramRegex := regexp.MustCompile(`([\?&]` + param + `=)([^&]+)`)
		masked = paramRegex.ReplaceAllString(masked, "${1}******")
	}

	return masked
}
