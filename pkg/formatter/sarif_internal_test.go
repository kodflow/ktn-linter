// Package formatter provides tests for internal SARIF formatter functions.
package formatter

import (
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/severity"
	sarif "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
	"golang.org/x/tools/go/analysis"
)

// Test_sarifFormatter_severityToSARIF tests severity level mapping.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_severityToSARIF(t *testing.T) {
	// Define test cases for severity mapping
	tests := []struct {
		name           string
		severityLevel  severity.Level
		expectedResult string
	}{
		{
			// Test unknown severity defaults to warning
			name:           "unknown severity returns warning",
			severityLevel:  severity.Level(99),
			expectedResult: "warning",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create formatter
			f := &sarifFormatter{}

			// Execute severity mapping
			result := f.severityToSARIF(tt.severityLevel)

			// Verify result
			if result != tt.expectedResult {
				t.Errorf("expected %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

// Test_sarifFormatter_addRule tests the addRule internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_addRule(t *testing.T) {
	// Define test cases for addRule method
	tests := []struct {
		name           string
		ruleID         string
		verbose        bool
		expectedRules  int
		expectedRuleID string
	}{
		{
			// Test adding a single rule
			name:           "add single rule",
			ruleID:         "KTN-VAR-001",
			verbose:        false,
			expectedRules:  1,
			expectedRuleID: "KTN-VAR-001",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create formatter
			f := &sarifFormatter{verbose: tt.verbose}

			// Create a new run
			run := sarif.NewRunWithInformationURI("test", "http://test.com")

			// Create adapter for interface compliance
			adapter := &sarifRunAdapter{run: run}

			// Add a rule
			f.addRule(adapter, tt.ruleID)

			// Verify rule was added
			if len(run.Tool.Driver.Rules) != tt.expectedRules {
				t.Fatalf("expected %d rule, got %d", tt.expectedRules, len(run.Tool.Driver.Rules))
			}

			// Guard clause for empty slice access
			if len(run.Tool.Driver.Rules) == 0 {
				t.Fatal("no rules found, cannot verify rule ID")
			}

			// Verify rule ID
			ruleID := run.Tool.Driver.Rules[0].ID
			// Check rule ID is set
			if ruleID == nil || *ruleID != tt.expectedRuleID {
				t.Errorf("expected rule ID %s, got %v", tt.expectedRuleID, ruleID)
			}
		})
	}
}

// Test_sarifFormatter_addResults tests the addResults internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_addResults(t *testing.T) {
	// Define test cases for addResults method
	tests := []struct {
		name            string
		verbose         bool
		message         string
		expectedResults int
		expectedRules   int
	}{
		{
			// Test adding diagnostic results
			name:            "add single diagnostic result",
			verbose:         false,
			message:         "KTN-VAR-001: test message",
			expectedResults: 1,
			expectedRules:   1,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create formatter
			f := &sarifFormatter{verbose: tt.verbose}

			// Create a new run
			run := sarif.NewRunWithInformationURI("test", "http://test.com")

			// Create fileset
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			// Create test diagnostics
			diags := []analysis.Diagnostic{
				{
					Pos:     file.Pos(10),
					Message: tt.message,
				},
			}

			// Create adapter for interface compliance
			adapter := &sarifRunAdapter{run: run}

			// Add results
			f.addResults(adapter, fset, diags)

			// Verify results were added
			if len(run.Results) != tt.expectedResults {
				t.Errorf("expected %d result, got %d", tt.expectedResults, len(run.Results))
			}

			// Verify rules were added
			if len(run.Tool.Driver.Rules) != tt.expectedRules {
				t.Errorf("expected %d rule, got %d", tt.expectedRules, len(run.Tool.Driver.Rules))
			}
		})
	}
}

// Test_sarifFormatter_RuleDeduplication tests that duplicate rules are not added.
//
// Params:
//   - t: testing object for running test cases
func Test_sarifFormatter_RuleDeduplication(t *testing.T) {
	// Define test cases for rule deduplication
	tests := []struct {
		name            string
		verbose         bool
		messages        []string
		expectedRules   int
		expectedResults int
	}{
		{
			// Test duplicate rules are deduplicated
			name:    "duplicate rules are deduplicated",
			verbose: false,
			messages: []string{
				"KTN-VAR-001: first message",
				"KTN-VAR-001: second message",
			},
			expectedRules:   1,
			expectedResults: 2,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create formatter
			f := &sarifFormatter{verbose: tt.verbose}

			// Create a new run
			run := sarif.NewRunWithInformationURI("test", "http://test.com")

			// Create fileset
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			// Create test diagnostics with same rule
			diags := make([]analysis.Diagnostic, len(tt.messages))
			// Create diagnostics from messages
			for i, msg := range tt.messages {
				diags[i] = analysis.Diagnostic{
					Pos:     file.Pos(10 + i*10),
					Message: msg,
				}
			}

			// Create adapter for interface compliance
			adapter := &sarifRunAdapter{run: run}

			// Add results
			f.addResults(adapter, fset, diags)

			// Verify only one rule was added
			if len(run.Tool.Driver.Rules) != tt.expectedRules {
				t.Errorf("expected %d rule (deduped), got %d", tt.expectedRules, len(run.Tool.Driver.Rules))
			}

			// Verify both results were added
			if len(run.Results) != tt.expectedResults {
				t.Errorf("expected %d results, got %d", tt.expectedResults, len(run.Results))
			}
		})
	}
}
