package rules_func_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/tests/target/rules_func"
)

// TestComplexFunctionWithComments teste la fonction avec commentaires internes profonds.
//
// Params:
//   - t: contexte de test
func TestComplexFunctionWithComments(t *testing.T) {
	result := rules_func.ComplexFunctionWithComments()
	if result != 0 {
		t.Logf("Résultat: %d", result)
	}
}
