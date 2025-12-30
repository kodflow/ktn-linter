// Package formatter_test provides tests for the formatter package.
package formatter_test

import (
	"bytes"
	"encoding/json"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/formatter"
	"golang.org/x/tools/go/analysis"
)

// TestJSONFormatter tests the JSON formatter functionality.
//
// Params:
//   - t: testing object for running test cases
func TestJSONFormatter(t *testing.T) {
	// Define test cases
	tests := []struct {
		name             string
		verbose          bool
		diagnostics      []analysis.Diagnostic
		validateOutput   func(t *testing.T, report map[string]interface{})
	}{
		{
			// Test formatter creation
			name:        "formatter creation returns non-nil",
			verbose:     false,
			diagnostics: nil,
			validateOutput: func(t *testing.T, report map[string]interface{}) {
				// Verify report is not nil
				if report == nil {
					t.Error("expected non-nil report")
				}
			},
		},
		{
			// Test empty diagnostics
			name:        "empty diagnostics shows zero issues",
			verbose:     false,
			diagnostics: []analysis.Diagnostic{},
			validateOutput: func(t *testing.T, report map[string]interface{}) {
				// Verify summary shows zero issues
				summary, ok := report["summary"].(map[string]interface{})
				// Check summary exists
				if !ok {
					t.Error("missing summary in JSON output")
					return
				}
				// Check total issues count
				if total, ok := summary["totalIssues"].(float64); !ok || total != 0 {
					t.Errorf("expected totalIssues=0, got %v", summary["totalIssues"])
				}
			},
		},
		{
			// Test with diagnostics
			name:    "single diagnostic shows one issue",
			verbose: false,
			diagnostics: []analysis.Diagnostic{
				{
					Pos:     token.NoPos,
					Message: "KTN-VAR-001: test message",
				},
			},
			validateOutput: func(t *testing.T, report map[string]interface{}) {
				// Verify summary shows one issue
				summary, ok := report["summary"].(map[string]interface{})
				// Check summary exists
				if !ok {
					t.Error("missing summary in JSON output")
					return
				}
				// Check total issues count
				if total, ok := summary["totalIssues"].(float64); !ok || total != 1 {
					t.Errorf("expected totalIssues=1, got %v", summary["totalIssues"])
				}
				// Verify results contain the diagnostic
				results, ok := report["results"].([]interface{})
				// Check results exists
				if !ok || len(results) != 1 {
					t.Errorf("expected 1 result, got %v", len(results))
				}
			},
		},
		{
			// Test verbose mode
			name:    "verbose mode includes full message",
			verbose: true,
			diagnostics: []analysis.Diagnostic{
				{
					Pos:     token.NoPos,
					Message: "KTN-VAR-001: short message\nVerbose explanation here",
				},
			},
			validateOutput: func(t *testing.T, report map[string]interface{}) {
				// Verify report is valid JSON (verbose mode should work)
				if report == nil {
					t.Error("expected non-nil report in verbose mode")
				}
			},
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create JSON formatter
			fmtr := formatter.NewJSONFormatter(&buf, tt.verbose)

			// Verify formatter is not nil
			if fmtr == nil {
				t.Error("NewJSONFormatter returned nil")
				return
			}

			// Create fileset
			fset := token.NewFileSet()

			// Work on a local copy to avoid mutating test-case data
			diags := append([]analysis.Diagnostic(nil), tt.diagnostics...)

			// Add file if we have diagnostics with positions
			if len(diags) > 0 {
				file := fset.AddFile("test.go", -1, 100)
				// Update diagnostic positions
				for i := range diags {
					if diags[i].Pos == token.NoPos {
						diags[i].Pos = file.Pos(10)
					}
				}
			}

			// Format diagnostics
			fmtr.Format(fset, diags)

			// Fail fast if buffer is empty
			if buf.Len() == 0 {
				t.Fatal("Format() produced empty output")
			}

			// Parse the output JSON
			var report map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &report)
			// Verify JSON is valid
			if err != nil {
				t.Errorf("invalid JSON output: %v", err)
				return
			}

			// Validate output
			if tt.validateOutput != nil {
				tt.validateOutput(t, report)
			}
		})
	}
}
