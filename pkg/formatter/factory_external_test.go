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

// TestNewFormatterByFormatText tests factory with text format.
//
// Params:
//   - t: testing object for running test cases
func TestNewFormatterByFormatText(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create formatter options
	opts := formatter.FormatterOptions{
		AIMode:      false,
		NoColor:     true,
		SimpleMode:  true,
		VerboseMode: false,
	}

	// Create formatter by format
	fmtr := formatter.NewFormatterByFormat(formatter.FormatText, &buf, opts)

	// Verify formatter is not nil
	if fmtr == nil {
		t.Error("NewFormatterByFormat returned nil for text format")
	}

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

	// Verify output is text format (not JSON)
	output := buf.String()
	// Check output is not JSON
	if strings.HasPrefix(strings.TrimSpace(output), "{") {
		t.Error("expected text output, got JSON")
	}
}

// TestNewFormatterByFormatJSON tests factory with JSON format.
//
// Params:
//   - t: testing object for running test cases
func TestNewFormatterByFormatJSON(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create formatter options
	opts := formatter.FormatterOptions{
		VerboseMode: false,
	}

	// Create formatter by format
	fmtr := formatter.NewFormatterByFormat(formatter.FormatJSON, &buf, opts)

	// Verify formatter is not nil
	if fmtr == nil {
		t.Error("NewFormatterByFormat returned nil for JSON format")
	}

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

	// Parse output as JSON
	var report map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &report)
	// Verify JSON is valid
	if err != nil {
		t.Errorf("expected valid JSON output: %v", err)
	}

	// Verify JSON has expected fields
	if _, ok := report["$schema"]; !ok {
		t.Error("missing $schema in JSON output")
	}
}

// TestNewFormatterByFormatSARIF tests factory with SARIF format.
//
// Params:
//   - t: testing object for running test cases
func TestNewFormatterByFormatSARIF(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create formatter options
	opts := formatter.FormatterOptions{
		VerboseMode: false,
	}

	// Create formatter by format
	fmtr := formatter.NewFormatterByFormat(formatter.FormatSARIF, &buf, opts)

	// Verify formatter is not nil
	if fmtr == nil {
		t.Error("NewFormatterByFormat returned nil for SARIF format")
	}

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

	// Parse output as JSON
	var report map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &report)
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
}

// TestNewFormatterByFormatUnknown tests factory with unknown format.
//
// Params:
//   - t: testing object for running test cases
func TestNewFormatterByFormatUnknown(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create formatter options
	opts := formatter.FormatterOptions{
		NoColor:    true,
		SimpleMode: true,
	}

	// Create formatter by format with unknown format
	fmtr := formatter.NewFormatterByFormat(formatter.OutputFormat("xml"), &buf, opts)

	// Verify formatter is not nil (should default to text)
	if fmtr == nil {
		t.Error("NewFormatterByFormat returned nil for unknown format")
	}

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

	// Verify output is text format (not JSON)
	output := buf.String()
	// Check output is not JSON
	if strings.HasPrefix(strings.TrimSpace(output), "{") {
		t.Error("expected text output for unknown format, got JSON")
	}
}

// TestFormatterOptionsPassthrough tests that options are passed correctly.
//
// Params:
//   - t: testing object for running test cases
func TestFormatterOptionsPassthrough(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create formatter options with verbose mode
	opts := formatter.FormatterOptions{
		AIMode:      false,
		NoColor:     true,
		SimpleMode:  true,
		VerboseMode: true,
	}

	// Create JSON formatter by format
	fmtr := formatter.NewFormatterByFormat(formatter.FormatJSON, &buf, opts)

	// Verify formatter is not nil
	if fmtr == nil {
		t.Error("NewFormatterByFormat returned nil")
	}
}
