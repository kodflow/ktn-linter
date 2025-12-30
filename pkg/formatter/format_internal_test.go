// Package formatter provides tests for internal format functions.
package formatter

import (
	"testing"
)

// Test_outputFormatString tests the string representation of OutputFormat.
//
// Params:
//   - t: testing object for running test cases
func Test_outputFormatString(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		format   OutputFormat
		expected string
	}{
		{
			// Test text format string
			name:     "text format string",
			format:   FormatText,
			expected: "text",
		},
		{
			// Test json format string
			name:     "json format string",
			format:   FormatJSON,
			expected: "json",
		},
		{
			// Test sarif format string
			name:     "sarif format string",
			format:   FormatSARIF,
			expected: "sarif",
		},
		{
			// Test custom format string
			name:     "custom format preserves value",
			format:   OutputFormat("custom"),
			expected: "custom",
		},
		{
			// Test empty format string
			name:     "empty format string",
			format:   OutputFormat(""),
			expected: "",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Get string representation
			result := string(tt.format)

			// Verify string matches expected
			if result != tt.expected {
				// Report error with details
				t.Errorf("string(%v) = %q, want %q", tt.format, result, tt.expected)
			}
		})
	}
}
