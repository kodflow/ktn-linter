package utils_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
)

// TestIsAllCaps tests the functionality of the corresponding implementation.
func TestIsAllCaps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"all caps", "HTTP", true},
		{"mixed case", "Http", false},
		{"lowercase", "http", false},
		{"empty string", "", false},
		{"numbers only", "123", false},
		{"caps with numbers", "HTTP2", true},
		{"underscore with caps", "HTTP_URL", true},
		{"single lowercase", "h", false},
		{"single uppercase", "H", true},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsAllCaps(tt.input)
			if got != tt.expected {
				t.Errorf("utils.IsAllCaps(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestIsMixedCaps tests the functionality of the corresponding implementation.
func TestIsMixedCaps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid MixedCaps", "MyVariable", true},
		{"valid mixedCaps", "myVariable", true},
		{"snake_case", "my_variable", false},
		{"ALL_CAPS non-initialism", "MY_VAR", false},
		{"valid initialism", "HTTPServer", true},
		{"empty string", "", false},
		{"single char", "x", true},
		{"numbers", "var123", true},
		{"all caps valid", "HTTP", true},
		{"all caps invalid", "ABC", false},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsMixedCaps(tt.input)
			if got != tt.expected {
				t.Errorf("utils.IsMixedCaps(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestIsValidInitialism tests the functionality of the corresponding implementation.
func TestIsValidInitialism(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid single initialism", "HTTP", true},
		{"valid double initialism", "HTTPURL", true},
		{"valid triple initialism", "HTTPURLID", true},
		{"invalid with underscore", "HTTP_URL", false},
		{"invalid mixed", "HTTPserver", false}, // server is not an initialism
		{"not an initialism", "NOTINIT", false},
		{"empty", "", false},
		{"partial match", "HTTPS", true},
		{"lowercase", "http", false},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsValidInitialism(tt.input)
			if got != tt.expected {
				t.Errorf("utils.IsValidInitialism(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
