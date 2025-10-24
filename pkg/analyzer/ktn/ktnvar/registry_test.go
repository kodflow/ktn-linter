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
	analyzers := Analyzers()

	// Vérifier que la liste n'est pas vide
	if len(analyzers) == 0 {
		t.Fatal("Analyzers() returned empty list")
	}

	// Vérifier que chaque analyseur a un nom et une doc
	for i, analyzer := range analyzers {
		// Vérification de la condition
		if analyzer == nil {
			t.Fatalf("Analyzer at index %d is nil", i)
		}

		// Vérification de la condition
		if analyzer.Name == "" {
			t.Errorf("Analyzer at index %d has empty name", i)
		}

		// Vérification de la condition
		if analyzer.Doc == "" {
			t.Errorf("Analyzer %s has empty documentation", analyzer.Name)
		}
	}

	// Vérifier le nombre attendu d'analyseurs (19 règles VAR)
	expectedCount := 19
	// Vérification de la condition
	if len(analyzers) != expectedCount {
		t.Errorf("Expected %d analyzers, got %d", expectedCount, len(analyzers))
	}
}
