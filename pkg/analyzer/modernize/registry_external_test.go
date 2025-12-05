// External tests for registry.go - modernize package.
package modernize_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/modernize"
)

// TestAnalyzers tests the public Analyzers function.
//
// Params:
//   - t: testing context
func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name      string
		wantEmpty bool
	}{
		{
			name:      "returns non-empty list",
			wantEmpty: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzers := modernize.Analyzers()
			// Vérification résultat
			if tt.wantEmpty && len(analyzers) > 0 {
				t.Error("Analyzers() returned non-empty when expected empty")
			}
			// Vérification non-vide
			if !tt.wantEmpty && len(analyzers) == 0 {
				t.Error("Analyzers() returned empty when expected non-empty")
			}
		})
	}
}
