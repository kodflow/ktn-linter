package ktnvar

import (
	"testing"
)

// TestAnalyzers vérifie que la fonction Analyzers retourne
// tous les analyseurs de la catégorie VAR.
//
// Returns:
//   - (voir code)
func TestAnalyzers(t *testing.T) {
	const EXPECTED_COUNT int = 19

	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "returns non-empty list",
			check: func(t *testing.T) {
				analyzers := Analyzers()
				// Vérification liste non vide
				if len(analyzers) == 0 {
					t.Fatal("Analyzers() returned empty list")
				}
			},
		},
		{
			name: "all analyzers have name and doc",
			check: func(t *testing.T) {
				analyzers := Analyzers()
				// Vérification chaque analyzer
				for i, analyzer := range analyzers {
					// Vérification non-nil
					if analyzer == nil {
						t.Fatalf("Analyzer at index %d is nil", i)
					}
					// Vérification nom
					if analyzer.Name == "" {
						t.Errorf("Analyzer at index %d has empty name", i)
					}
					// Vérification doc
					if analyzer.Doc == "" {
						t.Errorf("Analyzer %s has empty documentation", analyzer.Name)
					}
				}
			},
		},
		{
			name: "returns expected count",
			check: func(t *testing.T) {
				analyzers := Analyzers()
				// Vérification nombre
				if len(analyzers) != EXPECTED_COUNT {
					t.Errorf("Expected %d analyzers, got %d", EXPECTED_COUNT, len(analyzers))
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
