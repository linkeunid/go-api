package util

import (
	"testing"
)

func TestMaskDsn(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard DSN",
			input:    "root:password@tcp(localhost:3306)/dbname",
			expected: "root:******@tcp(localhost:3306)/dbname",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "No password",
			input:    "root@localhost",
			expected: "root@localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskDsn(tt.input)
			if result != tt.expected {
				t.Errorf("MaskDsn(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaskSensitive(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		visiblePrefixChars int
		visibleSuffixChars int
		expected           string
	}{
		{
			name:               "Regular string",
			input:              "supersecretpassword",
			visiblePrefixChars: 2,
			visibleSuffixChars: 2,
			expected:           "su***************rd",
		},
		{
			name:               "Short string",
			input:              "secret",
			visiblePrefixChars: 2,
			visibleSuffixChars: 2,
			expected:           "se**et",
		},
		{
			name:               "Very short string",
			input:              "pwd",
			visiblePrefixChars: 2,
			visibleSuffixChars: 2,
			expected:           "***",
		},
		{
			name:               "Empty string",
			input:              "",
			visiblePrefixChars: 2,
			visibleSuffixChars: 2,
			expected:           "",
		},
		{
			name:               "No suffix visible",
			input:              "apikey12345",
			visiblePrefixChars: 3,
			visibleSuffixChars: 0,
			expected:           "api********",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskSensitive(tt.input, tt.visiblePrefixChars, tt.visibleSuffixChars)
			if result != tt.expected {
				t.Errorf("MaskSensitive(%q, %d, %d) = %q, want %q",
					tt.input, tt.visiblePrefixChars, tt.visibleSuffixChars, result, tt.expected)
			}
		})
	}
}

func TestMaskCredential(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard credential",
			input:    "supersecretpassword",
			expected: "su***************rd",
		},
		{
			name:     "Short credential",
			input:    "secret",
			expected: "se**et",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskCredential(tt.input)
			if result != tt.expected {
				t.Errorf("MaskCredential(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard email",
			input:    "john.doe@example.com",
			expected: "jo******@example.com",
		},
		{
			name:     "Short local part",
			input:    "jo@example.com",
			expected: "jo@example.com", // No masking for very short local parts
		},
		{
			name:     "Invalid email format",
			input:    "not-an-email",
			expected: "no********il",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskEmail(tt.input)
			if result != tt.expected {
				t.Errorf("MaskEmail(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaskJWT(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard JWT",
			input:    "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyMTIzIn0.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ",
			expected: "eyJh******.eyJz******.TJVA9******",
		},
		{
			name:     "Invalid JWT format",
			input:    "invalid-token",
			expected: "inva******ken",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskJWT(tt.input)
			if result != tt.expected {
				t.Errorf("MaskJWT(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaskURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with basic auth",
			input:    "https://username:password@example.com/api",
			expected: "https://username:******@example.com/api",
		},
		{
			name:     "URL with token",
			input:    "https://example.com/api?token=12345abcde",
			expected: "https://example.com/api?token=******",
		},
		{
			name:     "URL with multiple sensitive params",
			input:    "https://example.com/api?token=12345abcde&user=john&api_key=secret",
			expected: "https://example.com/api?token=******&user=john&api_key=******",
		},
		{
			name:     "Regular URL",
			input:    "https://example.com/api",
			expected: "https://example.com/api",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskURL(tt.input)
			if result != tt.expected {
				t.Errorf("MaskURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
