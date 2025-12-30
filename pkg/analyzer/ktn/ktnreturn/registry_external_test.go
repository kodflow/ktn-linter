// External tests for ktnreturn registry.
package ktnreturn_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn"
)

// TestAnalyzers tests the public Analyzers function.
func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"verify analyzers are returned"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			analyzers := ktnreturn.Analyzers()

			// Vérification de la longueur
			if len(analyzers) == 0 {
				t.Error("Analyzers() should return non-empty slice")
			}

			// Vérification des analyseurs
			for i, a := range analyzers {
				// Vérification nil
				if a == nil {
					t.Errorf("Analyzer at index %d is nil", i)
				}
			}

			// Vérifier que Analyzer001 est présent
			expectedNames := map[string]bool{
				"ktnreturn001": true,
			}

			// Vérification noms attendus
			for _, analyzer := range analyzers {
				// Vérification nom dans map
				if !expectedNames[analyzer.Name] {
					t.Errorf("Unexpected analyzer name: %s", analyzer.Name)
				}
			}
		})
	}
}
