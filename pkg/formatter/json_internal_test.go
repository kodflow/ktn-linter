// Package formatter provides tests for internal JSON formatter functions.
package formatter

import (
	"bytes"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Test_jsonFormatter_buildResult tests the buildResult internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_jsonFormatter_buildResult(t *testing.T) {
	// Create formatter
	f := &jsonFormatter{verbose: false}

	// Create fileset and position
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)

	// Create test diagnostic
	diag := analysis.Diagnostic{
		Pos:     file.Pos(10),
		Message: "KTN-VAR-001: test message",
	}

	// Build result
	result := f.buildResult(fset, diag)

	// Verify rule ID
	if result.RuleID != "KTN-VAR-001" {
		t.Errorf("expected ruleId KTN-VAR-001, got %s", result.RuleID)
	}

	// Verify level is error (KTN-VAR-001 is error severity)
	if result.Level != "error" {
		t.Errorf("expected level error, got %s", result.Level)
	}
}

// Test_jsonFormatter_severityToLevel tests severity level mapping.
//
// Params:
//   - t: testing object for running test cases
func Test_jsonFormatter_severityToLevel(t *testing.T) {
	// Create formatter
	f := &jsonFormatter{}

	// Test unknown severity defaults to warning
	result := f.severityToLevel(99)

	// Verify default
	if result != "warning" {
		t.Errorf("expected warning for unknown severity, got %s", result)
	}
}

// Test_jsonFormatter_buildReport tests the buildReport internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_jsonFormatter_buildReport(t *testing.T) {
	// Create formatter
	f := &jsonFormatter{verbose: false}

	// Create fileset
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)

	// Create test diagnostics
	diags := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "KTN-VAR-001: error message",
		},
		{
			Pos:     file.Pos(20),
			Message: "KTN-CONST-002: info message",
		},
	}

	// Build report
	report := f.buildReport(fset, diags)

	// Verify total issues
	if report.Summary.TotalIssues != 2 {
		t.Errorf("expected 2 total issues, got %d", report.Summary.TotalIssues)
	}

	// Verify results count
	if len(report.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(report.Results))
	}
}

// Test_jsonFormatter_Format tests the Format method.
//
// Params:
//   - t: testing object for running test cases
func Test_jsonFormatter_Format(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create formatter
	f := &jsonFormatter{
		writer:  &buf,
		verbose: false,
	}

	// Create fileset
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)

	// Create test diagnostics
	diags := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "KTN-VAR-001: test message",
		},
	}

	// Format diagnostics
	f.Format(fset, diags)

	// Verify output is not empty
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
