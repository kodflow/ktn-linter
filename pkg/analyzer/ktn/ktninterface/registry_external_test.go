// External tests for ktninterface registry.
package ktninterface_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
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
			analyzers := ktninterface.Analyzers()

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

			// Vérifier que les analyseurs sont présents
			expectedNames := map[string]bool{
				"ktninterface001": true,
				"ktninterface003": true,
				"ktninterface004": true,
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
