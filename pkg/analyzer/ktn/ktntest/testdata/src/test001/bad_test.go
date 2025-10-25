package test001

import "testing"

// TestBadPackageName utilise le mauvais nom de package (test001 au lieu de test001_test)
func TestBadPackageName(t *testing.T) {
	t.Log("Devrait utiliser test001_test")
}
