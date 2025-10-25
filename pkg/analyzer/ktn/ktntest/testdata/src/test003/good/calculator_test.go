package test003_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest/testdata/src/test003/good"
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

// TestMultiply teste la fonction Multiply.
//
// Params:
//   - t: contexte de test
func TestMultiply(t *testing.T) {
	result := test003.Multiply(4, 5)
	// Vérification du résultat
	if result != 20 {
		t.Errorf("Multiply(4, 5) = %d, want 20", result)
	}
}

// TestDivide teste la fonction Divide.
//
// Params:
//   - t: contexte de test
func TestDivide(t *testing.T) {
	result, ok := test003.Divide(20, 4)
	// Vérification pas d'erreur
	if !ok {
		t.Error("Divide(20, 4) ne devrait pas échouer")
	}
	// Vérification du résultat
	if result != 5 {
		t.Errorf("Divide(20, 4) = %d, want 5", result)
	}
}
