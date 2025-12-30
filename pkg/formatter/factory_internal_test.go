// Package formatter provides tests for internal factory functions.
package formatter

import (
	"testing"
)

// Test_FormatterOptions tests the FormatterOptions struct behavior.
//
// Params:
//   - t: testing object for running test cases
func Test_FormatterOptions(t *testing.T) {
	// Define test cases
	tests := []struct {
		name            string
		opts            FormatterOptions
		expectAIMode    bool
		expectNoColor   bool
		expectSimple    bool
		expectVerbose   bool
	}{
		{
			// Test default options
			name:            "default options all false",
			opts:            FormatterOptions{},
			expectAIMode:    false,
			expectNoColor:   false,
			expectSimple:    false,
			expectVerbose:   false,
		},
		{
			// Test AI mode enabled
			name:            "ai mode enabled",
			opts:            FormatterOptions{AIMode: true},
			expectAIMode:    true,
			expectNoColor:   false,
			expectSimple:    false,
			expectVerbose:   false,
		},
		{
			// Test no color enabled
			name:            "no color enabled",
			opts:            FormatterOptions{NoColor: true},
			expectAIMode:    false,
			expectNoColor:   true,
			expectSimple:    false,
			expectVerbose:   false,
		},
		{
			// Test simple mode enabled
			name:            "simple mode enabled",
			opts:            FormatterOptions{SimpleMode: true},
			expectAIMode:    false,
			expectNoColor:   false,
			expectSimple:    true,
			expectVerbose:   false,
		},
		{
			// Test verbose mode enabled
			name:            "verbose mode enabled",
			opts:            FormatterOptions{VerboseMode: true},
			expectAIMode:    false,
			expectNoColor:   false,
			expectSimple:    false,
			expectVerbose:   true,
		},
		{
			// Test all options enabled
			name: "all options enabled",
			opts: FormatterOptions{
				AIMode:      true,
				NoColor:     true,
				SimpleMode:  true,
				VerboseMode: true,
			},
			expectAIMode:    true,
			expectNoColor:   true,
			expectSimple:    true,
			expectVerbose:   true,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Verify AIMode
			if tt.opts.AIMode != tt.expectAIMode {
				t.Errorf("AIMode = %v, want %v", tt.opts.AIMode, tt.expectAIMode)
			}
			// Verify NoColor
			if tt.opts.NoColor != tt.expectNoColor {
				t.Errorf("NoColor = %v, want %v", tt.opts.NoColor, tt.expectNoColor)
			}
			// Verify SimpleMode
			if tt.opts.SimpleMode != tt.expectSimple {
				t.Errorf("SimpleMode = %v, want %v", tt.opts.SimpleMode, tt.expectSimple)
			}
			// Verify VerboseMode
			if tt.opts.VerboseMode != tt.expectVerbose {
				t.Errorf("VerboseMode = %v, want %v", tt.opts.VerboseMode, tt.expectVerbose)
			}
		})
	}
}
