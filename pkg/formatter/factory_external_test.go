// Package formatter_test provides tests for the formatter package.
package formatter_test

import (
	"bytes"
	"encoding/json"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/formatter"
	"golang.org/x/tools/go/analysis"
)

// TestNewFormatterByFormat tests the NewFormatterByFormat factory function.
//
// Params:
//   - t: testing object for running test cases
func TestNewFormatterByFormat(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		format         formatter.OutputFormat
		opts           formatter.FormatterOptions
		expectNonNil   bool
		validateOutput func(t *testing.T, output string)
	}{
		{
			// Test text format
			name:   "text format returns text formatter",
			format: formatter.FormatText,
			opts: formatter.FormatterOptions{
				AIMode:      false,
				NoColor:     true,
				SimpleMode:  true,
				VerboseMode: false,
			},
			expectNonNil: true,
			validateOutput: func(t *testing.T, output string) {
				// Text output should not start with JSON
				if strings.HasPrefix(strings.TrimSpace(output), "{") {
					t.Error("expected text output, got JSON")
				}
			},
		},
		{
			// Test JSON format
			name:   "json format returns json formatter",
			format: formatter.FormatJSON,
			opts: formatter.FormatterOptions{
				VerboseMode: false,
			},
			expectNonNil: true,
			validateOutput: func(t *testing.T, output string) {
				// Parse output as JSON
				var report map[string]interface{}
				err := json.Unmarshal([]byte(output), &report)
				// Verify JSON is valid
				if err != nil {
					t.Errorf("expected valid JSON output: %v", err)
				}
				// Verify JSON has expected fields
				if _, ok := report["$schema"]; !ok {
					t.Error("missing $schema in JSON output")
				}
			},
		},
		{
			// Test SARIF format
			name:   "sarif format returns sarif formatter",
			format: formatter.FormatSARIF,
			opts: formatter.FormatterOptions{
				VerboseMode: false,
			},
			expectNonNil: true,
			validateOutput: func(t *testing.T, output string) {
				// Parse output as JSON
				var report map[string]interface{}
				err := json.Unmarshal([]byte(output), &report)
				// Verify JSON is valid
				if err != nil {
					t.Errorf("expected valid SARIF output: %v", err)
				}
				// Verify SARIF version
				version, ok := report["version"].(string)
				// Check version is SARIF 2.1.0
				if !ok || version != "2.1.0" {
					t.Errorf("expected SARIF version 2.1.0, got %v", version)
				}
			},
		},
		{
			// Test unknown format defaults to text
			name:   "unknown format defaults to text formatter",
			format: formatter.OutputFormat("xml"),
			opts: formatter.FormatterOptions{
				NoColor:    true,
				SimpleMode: true,
			},
			expectNonNil: true,
			validateOutput: func(t *testing.T, output string) {
				// Text output should not start with JSON
				if strings.HasPrefix(strings.TrimSpace(output), "{") {
					t.Error("expected text output for unknown format, got JSON")
				}
			},
		},
		{
			// Test options passthrough with verbose mode
			name:   "options are passed to formatter",
			format: formatter.FormatJSON,
			opts: formatter.FormatterOptions{
				AIMode:      false,
				NoColor:     true,
				SimpleMode:  true,
				VerboseMode: true,
			},
			expectNonNil: true,
			validateOutput: func(t *testing.T, output string) {
				// Just verify output is valid JSON
				var report map[string]interface{}
				err := json.Unmarshal([]byte(output), &report)
				// Verify JSON is valid
				if err != nil {
					t.Errorf("expected valid JSON output: %v", err)
				}
			},
		},
	}

	// Run all test cases
	for _, tt := range tests {
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create formatter by format
			fmtr := formatter.NewFormatterByFormat(tt.format, &buf, tt.opts)

			// Verify formatter is not nil
			if tt.expectNonNil && fmtr == nil {
				t.Errorf("NewFormatterByFormat returned nil for format %v", tt.format)
				return
			}

			// Create fileset and diagnostic for output validation
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)
			diags := []analysis.Diagnostic{
				{
					Pos:     file.Pos(10),
					Message: "KTN-VAR-001: test message",
				},
			}

			// Format diagnostics
			fmtr.Format(fset, diags)

			// Validate output if validator provided
			if tt.validateOutput != nil {
				tt.validateOutput(t, buf.String())
			}
		})
	}
}
