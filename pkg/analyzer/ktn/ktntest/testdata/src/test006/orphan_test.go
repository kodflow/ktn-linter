package test006 // want "KTN-TEST-006: fichier de test 'orphan_test.go' n'a pas de fichier source correspondant 'orphan.go' dans le même package. Dispatcher son contenu dans les fichiers de test appropriés puis le supprimer"

import "testing"

// TestOrphanFunction tests a function that doesn't exist in a corresponding source file
// This file is orphan because there's no "orphan.go" file
func TestOrphanFunction(t *testing.T) {
	// This test file has no corresponding source file
	t.Log("This is an orphan test file")
}
