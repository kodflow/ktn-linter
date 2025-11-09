package test008 // want "KTN-TEST-008: le fichier 'wrong_external_test.go' doit utiliser 'package xxx_test' \\(pas 'test008'\\) pour les tests externes"

import "testing"

// TestBadExternal utilise le mauvais package pour un test externe
func TestBadExternal(t *testing.T) {
	t.Log("This should use package test008_test, not test008")
}
