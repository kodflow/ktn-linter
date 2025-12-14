package ktnstruct_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
)

// TestGetAnalyzers teste le registre des analyseurs de structures.
func TestGetAnalyzers(t *testing.T) {
	analyzers := ktnstruct.GetAnalyzers()

	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "returns_non_empty_list",
			check: func(t *testing.T) {
				// Vérification liste non vide
				if len(analyzers) == 0 {
					t.Error("GetAnalyzers() returned empty list")
				}
			},
		},
		{
			name: "returns_expected_count",
			check: func(t *testing.T) {
				const expectedCount int = 7 // 5 + STRUCT-006 + STRUCT-007
				// Vérification nombre d'analyseurs
				if len(analyzers) != expectedCount {
					t.Errorf("Expected %d analyzers, got %d", expectedCount, len(analyzers))
				}
			},
		},
		{
			name: "all_analyzers_are_non_nil",
			check: func(t *testing.T) {
				// Vérification analyseurs non nil
				for i, analyzer := range analyzers {
					// Vérification nil
					if analyzer == nil {
						t.Errorf("Analyzer at index %d is nil", i)
					}
				}
			},
		},
		{
			name: "analyzers_have_expected_names",
			check: func(t *testing.T) {
				expectedNames := []string{
					"ktnstruct001",
					"ktnstruct002",
					"ktnstruct003",
					"ktnstruct004",
					"ktnstruct005",
					"ktnstruct006",
					"ktnstruct007",
				}

				// Vérification noms
				for i, analyzer := range analyzers {
					// Vérification nom
					if analyzer.Name != expectedNames[i] {
						t.Errorf("Expected analyzer name '%s' at index %d, got '%s'",
							expectedNames[i], i, analyzer.Name)
					}
				}
			},
		},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
