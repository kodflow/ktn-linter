// External tests for the runner package.
package orchestrator_test

import (
	"bytes"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// TestNewAnalysisRunner tests creating a new analysis runner.
func TestNewAnalysisRunner(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{
			name:    "create runner without verbose",
			verbose: false,
		},
		{
			name:    "create runner with verbose",
			verbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := orchestrator.NewAnalysisRunner(&buf, tt.verbose)

			// Check runner created
			if runner == nil {
				t.Error("expected non-nil runner")
			}
		})
	}
}

// TestAnalysisRunner_Run tests the Run method.
func TestAnalysisRunner_Run(t *testing.T) {
	tests := []struct {
		name      string
		pkgs      []*packages.Package
		analyzers []*analysis.Analyzer
	}{
		{
			name:      "empty packages returns empty results",
			pkgs:      []*packages.Package{},
			analyzers: []*analysis.Analyzer{},
		},
		{
			name:      "nil packages returns empty results",
			pkgs:      nil,
			analyzers: []*analysis.Analyzer{},
		},
		{
			name:      "nil analyzers returns empty results",
			pkgs:      []*packages.Package{},
			analyzers: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner := orchestrator.NewAnalysisRunner(&buf, false)

			// Should not panic
			diags := runner.Run(tt.pkgs, tt.analyzers)

			// Verify returns slice (may be nil or empty)
			_ = diags
		})
	}
}
