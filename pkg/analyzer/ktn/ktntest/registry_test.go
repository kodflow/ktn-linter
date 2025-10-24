package ktntest

import (
	"testing"
)

// TestAnalyzers tests the functionality of Analyzers
func TestAnalyzers(t *testing.T) {
	analyzers := Analyzers()

	// Vérifier que la liste n'est pas vide
	if len(analyzers) == 0 {
		t.Error("Analyzers should return at least one analyzer")
	}

	// Vérifier que chaque analyseur est valide
	for _, analyzer := range analyzers {
		// Vérifier que l'analyseur n'est pas nil
		if analyzer == nil {
			t.Error("Analyzers returned nil analyzer")
		}

		// Vérifier que l'analyseur a un nom
		if analyzer.Name == "" {
			t.Error("Analyzer has empty name")
		}

		// Vérifier que l'analyseur a une doc
		if analyzer.Doc == "" {
			t.Error("Analyzer has empty documentation")
		}

		// Vérifier que l'analyseur a une fonction Run
		if analyzer.Run == nil {
			t.Errorf("Analyzer %s has nil Run function", analyzer.Name)
		}
	}
}
