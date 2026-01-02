package ktnvar

import (
	"testing"
)

// TestHasUnderscore003 tests the hasUnderscore003 function.
func TestHasUnderscore003(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "blank identifier",
			input:    "_",
			expected: false,
		},
		{
			name:     "camelCase",
			input:    "myVariable",
			expected: false,
		},
		{
			name:     "PascalCase",
			input:    "MyVariable",
			expected: false,
		},
		{
			name:     "snake_case",
			input:    "my_variable",
			expected: true,
		},
		{
			name:     "SCREAMING_SNAKE_CASE",
			input:    "MY_VARIABLE",
			expected: true,
		},
		{
			name:     "mixed_Case",
			input:    "My_Variable",
			expected: true,
		},
		{
			name:     "single underscore prefix",
			input:    "_private",
			expected: true,
		},
		{
			name:     "acronym",
			input:    "HTTPStatus",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasUnderscore003(tt.input)
			if result != tt.expected {
				t.Errorf("hasUnderscore003(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
