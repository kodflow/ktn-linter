package rules_func_test

import "testing"

// TestReturnCommentsViolations vérifie la détection des violations KTN-FUNC-008.
//
// Params:
//   - t: contexte de test
func TestReturnCommentsViolations(t *testing.T) {
	// Succès car ce fichier source doit générer des violations
	t.Log("KTN-FUNC-008: Ce fichier doit générer des violations pour tous les returns sans commentaire")
}
