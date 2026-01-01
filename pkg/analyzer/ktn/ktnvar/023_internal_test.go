package ktnvar

import (
	"testing"
)

// TestIsSecurityName tests the isSecurityName function.
func TestIsSecurityName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "contains key",
			input:    "generateKey",
			expected: true,
		},
		{
			name:     "contains token",
			input:    "createToken",
			expected: true,
		},
		{
			name:     "contains secret",
			input:    "badSecretKey",
			expected: true,
		},
		{
			name:     "contains password",
			input:    "hashPassword",
			expected: true,
		},
		{
			name:     "contains salt",
			input:    "generateSalt",
			expected: true,
		},
		{
			name:     "contains nonce",
			input:    "createNonce",
			expected: true,
		},
		{
			name:     "contains crypt",
			input:    "encryptData",
			expected: true,
		},
		{
			name:     "contains auth",
			input:    "authHandler",
			expected: true,
		},
		{
			name:     "contains credential",
			input:    "getCredentials",
			expected: true,
		},
		{
			name:     "no security keyword",
			input:    "shuffleItems",
			expected: false,
		},
		{
			name:     "random index",
			input:    "randomIndex",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "case insensitive KEY",
			input:    "generateKEY",
			expected: true,
		},
		{
			name:     "case insensitive Token",
			input:    "CreateToken",
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isSecurityName(tt.input)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isSecurityName(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
