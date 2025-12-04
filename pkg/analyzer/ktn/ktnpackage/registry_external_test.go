package ktnpackage_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnpackage"
)

// TestAnalyzers teste que le registry retourne les analyseurs.
//
// Params:
//   - t: contexte de test
func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "analyzers should be valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzers := ktnpackage.Analyzers()

			// Vérifier qu'on a au moins 1 analyseur et tous non-nil
			if len(analyzers) < 1 {
				t.Errorf("expected at least 1 analyzer, got %d", len(analyzers))
				return
			}

			for i, analyzer := range analyzers {
				if analyzer == nil {
					t.Errorf("analyzer at index %d is nil", i)
				}
			}
		})
	}
}
