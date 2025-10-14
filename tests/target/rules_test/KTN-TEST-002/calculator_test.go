// ✅ CORRIGÉ: Ce fichier existe maintenant (relation 1:1 avec calculator.go)
package KTN_TEST_002_GOOD_test

import (
	"testing"

	. "github.com/kodflow/ktn-linter/tests/target/rules_test/KTN-TEST-002"
)

// TestAdd teste l'addition.
//
// Params:
//   - t: instance de test
func TestAdd(t *testing.T) {
	calc := &CalculatorData{}
	result := calc.Add(2, 3)
	if result != 5 {
		t.Errorf("expected 5, got %d", result)
	}
}

// TestMultiply teste la multiplication.
//
// Params:
//   - t: instance de test
func TestMultiply(t *testing.T) {
	calc := &CalculatorData{}
	result := calc.Multiply(2, 3)
	if result != 6 {
		t.Errorf("expected 6, got %d", result)
	}
}
