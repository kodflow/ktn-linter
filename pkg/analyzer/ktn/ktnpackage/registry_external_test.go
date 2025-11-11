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
	// Récupérer les analyseurs
	analyzers := ktnpackage.Analyzers()

	// Vérifier qu'on a au moins 1 analyseur
	if len(analyzers) < 1 {
		t.Errorf("expected at least 1 analyzer, got %d", len(analyzers))
	}

	// Vérifier que tous les analyseurs sont non-nil
	for i, analyzer := range analyzers {
		// Vérification non-nil
		if analyzer == nil {
			t.Errorf("analyzer at index %d is nil", i)
		}
	}
}
