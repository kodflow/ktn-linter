// External tests for modernize registry.
package modernize_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/modernize"
)

// TestAnalyzers tests the public Analyzers function.
func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"verify analyzers are returned"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzers := modernize.Analyzers()

			// Vérification de la longueur (au moins quelques analyseurs après filtrage)
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

			// Vérifier que newexpr est bien exclu (désactivé)
			for _, analyzer := range analyzers {
				// Vérification que newexpr n'est pas présent
				if analyzer.Name == "newexpr" {
					t.Error("Analyzer 'newexpr' should be disabled")
				}
			}
		})
	}
}
