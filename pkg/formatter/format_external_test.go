package formatter_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/formatter"
)

// TestParseOutputFormat tests the ParseOutputFormat function with various inputs.
//
// Params:
//   - t: testing object for running test cases
func TestParseOutputFormat(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    string
		expected formatter.OutputFormat
	}{
		{
			// Test valid text format
			name:     "valid text format",
			input:    "text",
			expected: formatter.FormatText,
		},
		{
			// Test valid JSON format
			name:     "valid json format",
			input:    "json",
			expected: formatter.FormatJSON,
		},
		{
			// Test valid SARIF format
			name:     "valid sarif format",
			input:    "sarif",
			expected: formatter.FormatSARIF,
		},
		{
			// Test unknown format defaults to text
			name:     "unknown format defaults to text",
			input:    "xml",
			expected: formatter.FormatText,
		},
		{
			// Test empty string defaults to text
			name:     "empty string defaults to text",
			input:    "",
			expected: formatter.FormatText,
		},
		{
			// Test case sensitivity
			name:     "uppercase JSON defaults to text",
			input:    "JSON",
			expected: formatter.FormatText,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Parse the format
			result := formatter.ParseOutputFormat(tt.input)

			// Verify the result matches expected
			if result != tt.expected {
				// Report error with details
				t.Errorf("ParseOutputFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestOutputFormatIsValid tests the IsValid method on OutputFormat.
//
// Params:
//   - t: testing object for running test cases
func TestOutputFormatIsValid(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		format   formatter.OutputFormat
		expected bool
	}{
		{
			// Test FormatText is valid
			name:     "FormatText is valid",
			format:   formatter.FormatText,
			expected: true,
		},
		{
			// Test FormatJSON is valid
			name:     "FormatJSON is valid",
			format:   formatter.FormatJSON,
			expected: true,
		},
		{
			// Test FormatSARIF is valid
			name:     "FormatSARIF is valid",
			format:   formatter.FormatSARIF,
			expected: true,
		},
		{
			// Test unknown format is invalid
			name:     "unknown format is invalid",
			format:   formatter.OutputFormat("xml"),
			expected: false,
		},
		{
			// Test empty format is invalid
			name:     "empty format is invalid",
			format:   formatter.OutputFormat(""),
			expected: false,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Check if format is valid
			result := tt.format.IsValid()

			// Verify the result matches expected
			if result != tt.expected {
				// Report error with details
				t.Errorf("OutputFormat(%q).IsValid() = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}

// TestFormatConstants tests that all format constants are valid.
//
// Params:
//   - t: testing object for running test cases
func TestFormatConstants(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		format formatter.OutputFormat
	}{
		{
			// Test FormatText constant
			name:   "FormatText is valid constant",
			format: formatter.FormatText,
		},
		{
			// Test FormatJSON constant
			name:   "FormatJSON is valid constant",
			format: formatter.FormatJSON,
		},
		{
			// Test FormatSARIF constant
			name:   "FormatSARIF is valid constant",
			format: formatter.FormatSARIF,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Check if format constant is valid
			if !tt.format.IsValid() {
				// Report error if constant is invalid
				t.Errorf("Format constant %q should be valid", tt.format)
			}
		})
	}
}

// TestParseOutputFormatRoundTrip tests parsing and validation together.
//
// Params:
//   - t: testing object for running test cases
func TestParseOutputFormatRoundTrip(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		formatStr string
	}{
		{
			// Test text format round trip
			name:      "text format round trip",
			formatStr: "text",
		},
		{
			// Test json format round trip
			name:      "json format round trip",
			formatStr: "json",
		},
		{
			// Test sarif format round trip
			name:      "sarif format round trip",
			formatStr: "sarif",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Parse the format string
			format := formatter.ParseOutputFormat(tt.formatStr)

			// Verify the parsed format is valid
			if !format.IsValid() {
				// Report error if valid format becomes invalid
				t.Errorf("ParseOutputFormat(%q).IsValid() = false, want true", tt.formatStr)
			}

			// Verify the format matches the input
			if string(format) != tt.formatStr {
				// Report error if format string doesn't match
				t.Errorf("ParseOutputFormat(%q) = %q, want %q", tt.formatStr, format, tt.formatStr)
			}
		})
	}
}
