package test003_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test003/bad"
)

// TestAdd teste la fonction Add.
//
// Params:
//   - t: contexte de test
func TestAdd(t *testing.T) {
	result := test003.Add(2, 3)
	// Vérification du résultat
	if result != 5 {
		t.Errorf("Add(2, 3) = %d, want 5", result)
	}
}

// TestSubtract teste la fonction Subtract.
//
// Params:
//   - t: contexte de test
func TestSubtract(t *testing.T) {
	result := test003.Subtract(10, 3)
	// Vérification du résultat
	if result != 7 {
		t.Errorf("Subtract(10, 3) = %d, want 7", result)
	}
}

// Pas de TestMultiply - devrait générer une erreur
// Pas de TestDivide - devrait générer une erreur
