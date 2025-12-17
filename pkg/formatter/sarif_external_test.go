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
	// Create buffer for output
	var buf bytes.Buffer

	// Create SARIF formatter
	fmtr := formatter.NewSARIFFormatter(&buf, false)

	// Verify formatter is not nil
	if fmtr == nil {
		t.Error("NewSARIFFormatter returned nil")
	}
}

// TestSARIFFormatterEmptyDiagnostics tests SARIF output with no diagnostics.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatterEmptyDiagnostics(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create SARIF formatter
	fmtr := formatter.NewSARIFFormatter(&buf, false)

	// Format empty diagnostics
	fmtr.Format(nil, nil)

	// Parse the output JSON
	var report map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &report)
	// Verify JSON is valid
	if err != nil {
		t.Errorf("invalid SARIF output: %v", err)
	}

	// Verify SARIF version
	version, ok := report["version"].(string)
	// Check version exists
	if !ok || version != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %v", version)
	}

	// Verify runs array exists
	runs, ok := report["runs"].([]interface{})
	// Check runs exists
	if !ok || len(runs) == 0 {
		t.Error("missing runs in SARIF output")
	}
}

// TestSARIFFormatterWithDiagnostics tests SARIF output with diagnostics.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatterWithDiagnostics(t *testing.T) {
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
			Message: "KTN-VAR-001: test message",
		},
	}

	// Format diagnostics
	fmtr.Format(fset, diags)

	// Parse the output JSON
	var report map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &report)
	// Verify JSON is valid
	if err != nil {
		t.Errorf("invalid SARIF output: %v", err)
	}

	// Get runs array
	runs, ok := report["runs"].([]interface{})
	// Check runs exists
	if !ok || len(runs) == 0 {
		t.Error("missing runs in SARIF output")
		return
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
	if !ok || len(results) != 1 {
		t.Errorf("expected 1 result, got %v", len(results))
	}
}

// TestSARIFFormatterToolInfo tests SARIF tool information.
//
// Params:
//   - t: testing object for running test cases
func TestSARIFFormatterToolInfo(t *testing.T) {
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
			Message: "KTN-VAR-001: test message",
		},
	}

	// Format diagnostics
	fmtr.Format(fset, diags)

	// Parse the output JSON
	var report map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &report)
	// Verify JSON is valid
	if err != nil {
		t.Errorf("invalid SARIF output: %v", err)
	}

	// Get runs array
	runs, ok := report["runs"].([]interface{})
	// Check runs exists
	if !ok || len(runs) == 0 {
		t.Error("missing runs in SARIF output")
		return
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
	if !ok || name != "ktn-linter" {
		t.Errorf("expected tool name 'ktn-linter', got %v", name)
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
			// Verify JSON is valid
			if err != nil {
				t.Errorf("invalid SARIF output: %v", err)
			}

			// Get runs array
			runs, ok := report["runs"].([]interface{})
			// Check runs exists
			if !ok || len(runs) == 0 {
				t.Error("missing runs in SARIF output")
				return
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
