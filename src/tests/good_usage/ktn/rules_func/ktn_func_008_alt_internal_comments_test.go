package rules_func

import (
	"testing"
)

// TestComplexFunctionWithComments teste la fonction avec commentaires internes profonds.
//
// Params:
//   - t: contexte de test
func TestComplexFunctionWithComments(t *testing.T) {
	result := ComplexFunctionWithComments()
	if result != 0 {
		t.Logf("Résultat: %d", result)
	}
}
