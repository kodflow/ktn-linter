package test008_test // want "KTN-TEST-008: le fichier 'wrong_internal_test.go' doit utiliser 'package wrong_internal' \\(pas 'test008_test'\\) pour les tests internes"

import "testing"

// TestBadInternal utilise le mauvais package pour un test interne
func TestBadInternal(t *testing.T) {
	t.Log("This should use package test008, not test008_test")
}
