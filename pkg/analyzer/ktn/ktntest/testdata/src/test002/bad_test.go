package test002_test

import "testing"

// TestOrphan est un test orphelin (pas de fichier bad.go correspondant)
func TestOrphan(t *testing.T) {
	t.Log("Ce fichier de test n'a pas de fichier source correspondant bad.go")
}
