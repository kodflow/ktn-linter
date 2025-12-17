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

// TestNewJSONFormatter tests the NewJSONFormatter function.
//
// Params:
//   - t: testing object for running test cases
func TestNewJSONFormatter(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create JSON formatter
	fmtr := formatter.NewJSONFormatter(&buf, false)

	// Verify formatter is not nil
	if fmtr == nil {
		t.Error("NewJSONFormatter returned nil")
	}
}

// TestJSONFormatterEmptyDiagnostics tests JSON output with no diagnostics.
//
// Params:
//   - t: testing object for running test cases
func TestJSONFormatterEmptyDiagnostics(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create JSON formatter
	fmtr := formatter.NewJSONFormatter(&buf, false)

	// Format empty diagnostics
	fmtr.Format(nil, nil)

	// Parse the output JSON
	var report map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &report)
	// Verify JSON is valid
	if err != nil {
		t.Errorf("invalid JSON output: %v", err)
	}

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
}

// TestJSONFormatterWithDiagnostics tests JSON output with diagnostics.
//
// Params:
//   - t: testing object for running test cases
func TestJSONFormatterWithDiagnostics(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create JSON formatter
	fmtr := formatter.NewJSONFormatter(&buf, false)

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
		t.Errorf("invalid JSON output: %v", err)
	}

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
}

// TestJSONFormatterVerboseMode tests JSON output in verbose mode.
//
// Params:
//   - t: testing object for running test cases
func TestJSONFormatterVerboseMode(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create JSON formatter with verbose mode
	fmtr := formatter.NewJSONFormatter(&buf, true)

	// Create fileset and diagnostic with verbose message
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)
	diags := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "KTN-VAR-001: short message\nVerbose explanation here",
		},
	}

	// Format diagnostics
	fmtr.Format(fset, diags)

	// Parse the output JSON
	var report map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &report)
	// Verify JSON is valid
	if err != nil {
		t.Errorf("invalid JSON output: %v", err)
	}
}
