// ✅ CORRIGÉ: Package name avec suffixe _test
package KTN_TEST_001_GOOD_test

import "testing"

// TestAddUser teste l'ajout d'utilisateur.
//
// Params:
//   - t: instance de test
func TestAddUser(t *testing.T) {
	// Note: Pour tester, on doit importer le package
	// Mais ici c'est juste pour montrer le bon format
	t.Log("Test avec bon package name")
}
