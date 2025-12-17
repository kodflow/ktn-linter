// Package formatter provides tests for internal SARIF formatter functions.
package formatter

import (
	"bytes"
	"go/token"
	"testing"

	sarif "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	"golang.org/x/tools/go/analysis"
)

// Test_sarifFormatter_severityToSARIF tests severity level mapping.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_severityToSARIF(t *testing.T) {
	// Create formatter
	f := &sarifFormatter{}

	// Test unknown severity defaults to warning
	result := f.severityToSARIF(99)

	// Verify default
	if result != "warning" {
		t.Errorf("expected warning for unknown severity, got %s", result)
	}
}

// Test_sarifFormatter_addRule tests the addRule internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_addRule(t *testing.T) {
	// Create formatter
	f := &sarifFormatter{verbose: false}

	// Create a new run
	run := sarif.NewRunWithInformationURI("test", "http://test.com")

	// Add a rule
	f.addRule(run, "KTN-VAR-001")

	// Verify rule was added
	if len(run.Tool.Driver.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(run.Tool.Driver.Rules))
	}

	// Verify rule ID
	ruleID := run.Tool.Driver.Rules[0].ID
	// Check rule ID is set
	if ruleID == nil || *ruleID != "KTN-VAR-001" {
		t.Errorf("expected rule ID KTN-VAR-001, got %v", ruleID)
	}
}

// Test_sarifFormatter_addResults tests the addResults internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_addResults(t *testing.T) {
	// Create formatter
	f := &sarifFormatter{verbose: false}

	// Create a new run
	run := sarif.NewRunWithInformationURI("test", "http://test.com")

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

	// Add results
	f.addResults(run, fset, diags)

	// Verify results were added
	if len(run.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(run.Results))
	}

	// Verify rules were added
	if len(run.Tool.Driver.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(run.Tool.Driver.Rules))
	}
}

// Test_sarifFormatter_Format tests the Format method.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_Format(t *testing.T) {
	// Create buffer for output
	var buf bytes.Buffer

	// Create formatter
	f := &sarifFormatter{
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

// Test_sarifFormatter_RuleDeduplication tests that duplicate rules are not added.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_RuleDeduplication(t *testing.T) {
	// Create formatter
	f := &sarifFormatter{verbose: false}

	// Create a new run
	run := sarif.NewRunWithInformationURI("test", "http://test.com")

	// Create fileset
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)

	// Create test diagnostics with same rule
	diags := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "KTN-VAR-001: first message",
		},
		{
			Pos:     file.Pos(20),
			Message: "KTN-VAR-001: second message",
		},
	}

	// Add results
	f.addResults(run, fset, diags)

	// Verify only one rule was added
	if len(run.Tool.Driver.Rules) != 1 {
		t.Errorf("expected 1 rule (deduped), got %d", len(run.Tool.Driver.Rules))
	}

	// Verify both results were added
	if len(run.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(run.Results))
	}
}
