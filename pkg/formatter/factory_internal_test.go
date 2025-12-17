// Package formatter provides tests for internal factory functions.
package formatter

import (
	"bytes"
	"testing"
)

// Test_NewFormatterByFormat tests internal factory behavior.
//
// Params:
//   - t: testing object for running test cases
func Test_NewFormatterByFormat(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		format OutputFormat
	}{
		{
			// Test JSON format
			name:   "json format",
			format: FormatJSON,
		},
		{
			// Test SARIF format
			name:   "sarif format",
			format: FormatSARIF,
		},
		{
			// Test text format
			name:   "text format",
			format: FormatText,
		},
		{
			// Test default format
			name:   "default format",
			format: OutputFormat("unknown"),
		},
	}

	// Run all test cases
	for _, tt := range tests {
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create formatter options
			opts := FormatterOptions{
				AIMode:      false,
				NoColor:     true,
				SimpleMode:  true,
				VerboseMode: false,
			}

			// Create formatter
			fmtr := NewFormatterByFormat(tt.format, &buf, opts)

			// Verify formatter is not nil
			if fmtr == nil {
				t.Errorf("NewFormatterByFormat returned nil for format %v", tt.format)
			}
		})
	}
}
