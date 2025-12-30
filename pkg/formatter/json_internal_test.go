// Package formatter provides tests for internal JSON formatter functions.
package formatter

import (
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/severity"
	"golang.org/x/tools/go/analysis"
)

// Test_jsonFormatter_buildResult tests the buildResult internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_jsonFormatter_buildResult(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		verbose        bool
		message        string
		expectedRuleID string
		expectedLevel  string
	}{
		{
			// Test error severity rule
			name:           "error severity rule",
			verbose:        false,
			message:        "KTN-VAR-001: test message",
			expectedRuleID: "KTN-VAR-001",
			expectedLevel:  "error",
		},
		{
			// Test warning severity rule
			name:           "warning severity rule",
			verbose:        false,
			message:        "KTN-VAR-003: test message",
			expectedRuleID: "KTN-VAR-003",
			expectedLevel:  "warning",
		},
		{
			// Test info severity rule
			name:           "info severity rule",
			verbose:        false,
			message:        "KTN-CONST-002: test message",
			expectedRuleID: "KTN-CONST-002",
			expectedLevel:  "info",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create formatter
			f := &jsonFormatter{verbose: tt.verbose}

			// Create fileset and position
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			// Create test diagnostic
			diag := analysis.Diagnostic{
				Pos:     file.Pos(10),
				Message: tt.message,
			}

			// Build result
			result := f.buildResult(fset, diag)

			// Verify rule ID
			if result.RuleID != tt.expectedRuleID {
				t.Errorf("expected ruleId %s, got %s", tt.expectedRuleID, result.RuleID)
			}

			// Verify level
			if result.Level != tt.expectedLevel {
				t.Errorf("expected level %s, got %s", tt.expectedLevel, result.Level)
			}
		})
	}
}

// Test_jsonFormatter_severityToLevel tests severity level mapping.
//
// Params:
//   - t: testing object for running test cases
func Test_jsonFormatter_severityToLevel(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		severityLevel  severity.Level
		expectedResult string
	}{
		{
			// Test error severity
			name:           "error severity returns error",
			severityLevel:  severity.SeverityError,
			expectedResult: "error",
		},
		{
			// Test warning severity
			name:           "warning severity returns warning",
			severityLevel:  severity.SeverityWarning,
			expectedResult: "warning",
		},
		{
			// Test info severity
			name:           "info severity returns info",
			severityLevel:  severity.SeverityInfo,
			expectedResult: "info",
		},
		{
			// Test unknown severity defaults to warning
			name:           "unknown severity defaults to warning",
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
			f := &jsonFormatter{}

			// Call severityToLevel
			result := f.severityToLevel(tt.severityLevel)

			// Verify result
			if result != tt.expectedResult {
				t.Errorf("severityToLevel(%v) = %s, want %s", tt.severityLevel, result, tt.expectedResult)
			}
		})
	}
}

// Test_jsonFormatter_buildReport tests the buildReport internal method.
//
// Params:
//   - t: testing object for running test cases
func Test_jsonFormatter_buildReport(t *testing.T) {
	// Define test cases
	tests := []struct {
		name              string
		verbose           bool
		diagnosticCount   int
		expectedTotal     int
		expectedResultLen int
	}{
		{
			// Test empty diagnostics
			name:              "empty diagnostics",
			verbose:           false,
			diagnosticCount:   0,
			expectedTotal:     0,
			expectedResultLen: 0,
		},
		{
			// Test single diagnostic
			name:              "single diagnostic",
			verbose:           false,
			diagnosticCount:   1,
			expectedTotal:     1,
			expectedResultLen: 1,
		},
		{
			// Test multiple diagnostics
			name:              "multiple diagnostics",
			verbose:           false,
			diagnosticCount:   3,
			expectedTotal:     3,
			expectedResultLen: 3,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Run individual test case
		t.Run(tt.name, func(t *testing.T) {
			// Create formatter
			f := &jsonFormatter{verbose: tt.verbose}

			// Create fileset
			fset := token.NewFileSet()
			file := fset.AddFile("test.go", -1, 100)

			// Create test diagnostics
			diags := make([]analysis.Diagnostic, tt.diagnosticCount)
			for i := range diags {
				diags[i] = analysis.Diagnostic{
					Pos:     file.Pos(10 + i),
					Message: "KTN-VAR-001: test message",
				}
			}

			// Build report
			report := f.buildReport(fset, diags)

			// Verify total issues
			if report.Summary.TotalIssues != tt.expectedTotal {
				t.Errorf("expected %d total issues, got %d", tt.expectedTotal, report.Summary.TotalIssues)
			}

			// Verify results count
			if len(report.Results) != tt.expectedResultLen {
				t.Errorf("expected %d results, got %d", tt.expectedResultLen, len(report.Results))
			}
		})
	}
}
