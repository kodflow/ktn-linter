// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-003: Fichier .go avec son _test.go (pas d'orphelin) (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_003_GOOD_test

// Helper fournit des fonctions utilitaires.
type Helper struct{}

// Format formate une chaîne.
//
// Params:
//   - input: la chaîne à formater
//
// Returns:
//   - string: la chaîne formatée
func (h *Helper) Format(input string) string {
	return "[" + input + "]"
}

// ✅ CORRIGÉ: helper_test.go existe (pas de fichier orphelin)
