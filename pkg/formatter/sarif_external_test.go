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

// TestNewSARIFFormatter tests the NewSARIFFormatter function.
//
// Params:
//   - t: testing object for running test cases
func TestNewSARIFFormatter(t *testing.T) {
	// Define test cases for NewSARIFFormatter
	tests := []struct {
		name     string
		verbose  bool
		expectNil bool
	}{
		{
			// Test creating SARIF formatter with verbose disabled
			name:      "create SARIF formatter with verbose disabled",
			verbose:   false,
			expectNil: false,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create SARIF formatter
			fmtr := formatter.NewSARIFFormatter(&buf, tt.verbose)

			// Verify formatter is not nil
			if (fmtr == nil) != tt.expectNil {
				t.Errorf("expected nil=%v, got nil=%v", tt.expectNil, fmtr == nil)
			}
		})
	}
}

// TestSARIFFormatter_Format tests the Format method of the SARIF formatter.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatter_Format(t *testing.T) {
	// Define test cases
	tests := []struct {
		name            string
		verbose         bool
		diagnostics     []analysis.Diagnostic
		expectNonEmpty  bool
		expectedVersion string
	}{
		{
			// Test formatting with diagnostics
			name:            "format with diagnostics produces valid SARIF",
			verbose:         false,
			expectedVersion: "2.1.0",
			expectNonEmpty:  true,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create SARIF formatter
			fmtr := formatter.NewSARIFFormatter(&buf, tt.verbose)

			// Create fileset and diagnostic
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

			// Verify output is not empty
			if tt.expectNonEmpty && buf.Len() == 0 {
				t.Error("expected non-empty output from Format")
			}

			// Parse and verify SARIF structure
			var report map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &report)
			// Verify JSON is valid
			if err != nil {
				t.Errorf("invalid SARIF output: %v", err)
				return
			}

			// Verify SARIF version
			version, ok := report["version"].(string)
			// Check version matches
			if !ok || version != tt.expectedVersion {
				t.Errorf("expected SARIF version %s, got %v", tt.expectedVersion, version)
			}
		})
	}
}

// TestSARIFFormatterEmptyDiagnostics tests SARIF output with no diagnostics.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatterEmptyDiagnostics(t *testing.T) {
	// Define test cases for empty diagnostics
	tests := []struct {
		name            string
		verbose         bool
		expectedVersion string
		expectRuns      bool
	}{
		{
			// Test empty diagnostics produces valid SARIF
			name:            "empty diagnostics produces valid SARIF 2.1.0",
			verbose:         false,
			expectedVersion: "2.1.0",
			expectRuns:      true,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create SARIF formatter
			fmtr := formatter.NewSARIFFormatter(&buf, tt.verbose)

			// Format empty diagnostics (use non-nil FileSet and empty slice)
			fmtr.Format(token.NewFileSet(), []analysis.Diagnostic{})

			// Parse the output JSON
			var report map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &report)
			// Fail fast on invalid JSON
			if err != nil {
				t.Fatalf("invalid SARIF output: %v", err)
			}

			// Verify SARIF version
			version, ok := report["version"].(string)
			// Check version exists
			if !ok || version != tt.expectedVersion {
				t.Errorf("expected SARIF version %s, got %v", tt.expectedVersion, version)
			}

			// Verify runs array exists
			runs, ok := report["runs"].([]interface{})
			// Check runs exists
			if tt.expectRuns && (!ok || len(runs) == 0) {
				t.Error("missing runs in SARIF output")
			}
		})
	}
}

// TestSARIFFormatterWithDiagnostics tests SARIF output with diagnostics.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatterWithDiagnostics(t *testing.T) {
	// Define test cases for diagnostics output
	tests := []struct {
		name            string
		verbose         bool
		message         string
		expectedResults int
	}{
		{
			// Test single diagnostic produces one result
			name:            "single diagnostic produces one result",
			verbose:         false,
			message:         "KTN-VAR-001: test message",
			expectedResults: 1,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create SARIF formatter
			fmtr := formatter.NewSARIFFormatter(&buf, tt.verbose)

			// Create fileset and diagnostic
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)
			diags := []analysis.Diagnostic{
				{
					Pos:     file.Pos(10),
					Message: tt.message,
				},
			}

			// Format diagnostics
			fmtr.Format(fset, diags)

			// Parse the output JSON
			var report map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &report)
			// Fail fast on invalid JSON
			if err != nil {
				t.Fatalf("invalid SARIF output: %v", err)
			}

			// Get runs array
			runs, ok := report["runs"].([]interface{})
			// Check runs exists - fail fast if missing
			if !ok || len(runs) == 0 {
				t.Fatal("missing runs in SARIF output")
			}

			// Get first run
			run, ok := runs[0].(map[string]interface{})
			// Check run exists
			if !ok {
				t.Error("invalid run structure")
				return
			}

			// Verify results exist
			results, ok := run["results"].([]interface{})
			// Check results exists
			if !ok || len(results) != tt.expectedResults {
				t.Errorf("expected %d result, got %v", tt.expectedResults, len(results))
			}
		})
	}
}

// TestSARIFFormatterToolInfo tests SARIF tool information.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatterToolInfo(t *testing.T) {
	// Define test cases for tool information
	tests := []struct {
		name             string
		verbose          bool
		message          string
		expectedToolName string
	}{
		{
			// Test tool name is ktn-linter
			name:             "tool name is ktn-linter",
			verbose:          false,
			message:          "KTN-VAR-001: test message",
			expectedToolName: "ktn-linter",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create SARIF formatter
			fmtr := formatter.NewSARIFFormatter(&buf, tt.verbose)

			// Create fileset and diagnostic
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)
			diags := []analysis.Diagnostic{
				{
					Pos:     file.Pos(10),
					Message: tt.message,
				},
			}

			// Format diagnostics
			fmtr.Format(fset, diags)

			// Parse the output JSON
			var report map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &report)
			// Fail fast on invalid JSON
			if err != nil {
				t.Fatalf("invalid SARIF output: %v", err)
			}

			// Get runs array
			runs, ok := report["runs"].([]interface{})
			// Check runs exists - fail fast if missing
			if !ok || len(runs) == 0 {
				t.Fatal("missing runs in SARIF output")
			}

			// Get first run
			run, ok := runs[0].(map[string]interface{})
			// Check run exists
			if !ok {
				t.Error("invalid run structure")
				return
			}

			// Verify tool information
			tool, ok := run["tool"].(map[string]interface{})
			// Check tool exists
			if !ok {
				t.Error("missing tool in SARIF output")
				return
			}

			// Verify driver information
			driver, ok := tool["driver"].(map[string]interface{})
			// Check driver exists
			if !ok {
				t.Error("missing driver in tool")
				return
			}

			// Verify tool name
			name, ok := driver["name"].(string)
			// Check name matches
			if !ok || name != tt.expectedToolName {
				t.Errorf("expected tool name '%s', got %v", tt.expectedToolName, name)
			}
		})
	}
}

// TestSARIFFormatterSeverityLevels tests SARIF severity mapping.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatterSeverityLevels(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		ruleCode      string
		expectedLevel string
	}{
		{
			// Test error level
			name:          "error severity",
			ruleCode:      "KTN-VAR-001",
			expectedLevel: "error",
		},
		{
			// Test warning level
			name:          "warning severity",
			ruleCode:      "KTN-VAR-003",
			expectedLevel: "warning",
		},
		{
			// Test info level
			name:          "info severity",
			ruleCode:      "KTN-CONST-002",
			expectedLevel: "note",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer for output
			var buf bytes.Buffer

			// Create SARIF formatter
			fmtr := formatter.NewSARIFFormatter(&buf, false)

			// Create fileset and diagnostic
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)
			diags := []analysis.Diagnostic{
				{
					Pos:     file.Pos(10),
					Message: tt.ruleCode + ": test message",
				},
			}

			// Format diagnostics
			fmtr.Format(fset, diags)

			// Parse the output JSON
			var report map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &report)
			// Fail fast on invalid JSON
			if err != nil {
				t.Fatalf("invalid SARIF output: %v", err)
			}

			// Get runs array
			runs, ok := report["runs"].([]interface{})
			// Check runs exists - fail fast if missing
			if !ok || len(runs) == 0 {
				t.Fatal("missing runs in SARIF output")
			}

			// Get first run
			run, ok := runs[0].(map[string]interface{})
			// Check run exists
			if !ok {
				t.Error("invalid run structure")
				return
			}

			// Get results
			results, ok := run["results"].([]interface{})
			// Check results exists
			if !ok || len(results) == 0 {
				t.Error("missing results in SARIF output")
				return
			}

			// Get first result
			result, ok := results[0].(map[string]interface{})
			// Check result exists
			if !ok {
				t.Error("invalid result structure")
				return
			}

			// Verify level
			level, ok := result["level"].(string)
			// Check level matches
			if !ok || level != tt.expectedLevel {
				t.Errorf("expected level %q, got %q", tt.expectedLevel, level)
			}
		})
	}
}
