// External tests for the selector package.
package orchestrator_test

import (
	"bytes"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/orchestrator"
)

// TestNewAnalyzerSelector tests creating a new analyzer selector.
func TestNewAnalyzerSelector(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{
			name:    "create selector without verbose",
			verbose: false,
		},
		{
			name:    "create selector with verbose",
			verbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			selector := orchestrator.NewAnalyzerSelector(&buf, tt.verbose)

			// Check selector created
			if selector == nil {
				t.Error("expected non-nil selector")
			}
		})
	}
}

// TestAnalyzerSelector_Select tests the Select method.
func TestAnalyzerSelector_Select(t *testing.T) {
	tests := []struct {
		name      string
		opts      orchestrator.Options
		expectErr bool
		minCount  int
	}{
		{
			name:      "select all rules",
			opts:      orchestrator.Options{},
			expectErr: false,
			minCount:  1,
		},
		{
			name:      "select by category",
			opts:      orchestrator.Options{Category: "func"},
			expectErr: false,
			minCount:  1,
		},
		{
			name:      "select single rule",
			opts:      orchestrator.Options{OnlyRule: "KTN-FUNC-001"},
			expectErr: false,
			minCount:  1,
		},
		{
			name:      "unknown category error",
			opts:      orchestrator.Options{Category: "nonexistent"},
			expectErr: true,
			minCount:  0,
		},
		{
			name:      "unknown rule error",
			opts:      orchestrator.Options{OnlyRule: "KTN-INVALID-999"},
			expectErr: true,
			minCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			selector := orchestrator.NewAnalyzerSelector(&buf, false)

			analyzers, err := selector.Select(tt.opts)

			// Check error expectation
			if tt.expectErr && err == nil {
				t.Error("expected error but got nil")
			}
			// Check no error expectation
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			// Check count
			if !tt.expectErr && len(analyzers) < tt.minCount {
				t.Errorf("expected at least %d analyzers, got %d", tt.minCount, len(analyzers))
			}
		})
	}
}
