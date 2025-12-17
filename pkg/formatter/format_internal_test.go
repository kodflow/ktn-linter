// Package formatter provides tests for internal format functions.
package formatter

import (
	"testing"
)

// TestOutputFormat_IsValid tests the IsValid method on OutputFormat.
//
// Params:
//   - t: testing object for running test cases
func TestOutputFormat_IsValid(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		format   OutputFormat
		expected bool
	}{
		{
			// Test FormatText is valid
			name:     "text is valid",
			format:   FormatText,
			expected: true,
		},
		{
			// Test FormatJSON is valid
			name:     "json is valid",
			format:   FormatJSON,
			expected: true,
		},
		{
			// Test FormatSARIF is valid
			name:     "sarif is valid",
			format:   FormatSARIF,
			expected: true,
		},
		{
			// Test unknown format is invalid
			name:     "xml is invalid",
			format:   OutputFormat("xml"),
			expected: false,
		},
		{
			// Test empty format is invalid
			name:     "empty is invalid",
			format:   OutputFormat(""),
			expected: false,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Check if format is valid
			result := tt.format.IsValid()

			// Verify result matches expected
			if result != tt.expected {
				// Report error with details
				t.Errorf("OutputFormat(%q).IsValid() = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}
