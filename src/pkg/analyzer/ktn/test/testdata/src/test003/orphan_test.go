package test003_test // want `\[KTN_TEST_003\] Fichier de test 'orphan_test.go' n'a pas de fichier source correspondant`

import "testing"

// Mauvais : fichier de test sans fichier source correspondant
func TestOrphan(t *testing.T) {
	t.Log("This test file has no corresponding source file")
}
