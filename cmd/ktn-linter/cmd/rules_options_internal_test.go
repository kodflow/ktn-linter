// Internal tests for rules options.
package cmd

import (
	"testing"
)

// Test_rulesOptions tests the rulesOptions struct.
func Test_rulesOptions(t *testing.T) {
	tests := []struct {
		name       string
		opts       rulesOptions
		wantFormat string
		wantNoEx   bool
	}{
		{
			name: "default options",
			opts: rulesOptions{
				Format:     "text",
				NoExamples: false,
			},
			wantFormat: "text",
			wantNoEx:   false,
		},
		{
			name: "markdown format",
			opts: rulesOptions{
				Format:     "markdown",
				NoExamples: true,
			},
			wantFormat: "markdown",
			wantNoEx:   true,
		},
		{
			name: "json format",
			opts: rulesOptions{
				Format:     "json",
				NoExamples: false,
			},
			wantFormat: "json",
			wantNoEx:   false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Verify format
			if tt.opts.Format != tt.wantFormat {
				t.Errorf("Format = %q, want %q", tt.opts.Format, tt.wantFormat)
			}
			// Verify no examples flag
			if tt.opts.NoExamples != tt.wantNoEx {
				t.Errorf("NoExamples = %v, want %v", tt.opts.NoExamples, tt.wantNoEx)
			}
		})
	}
}
