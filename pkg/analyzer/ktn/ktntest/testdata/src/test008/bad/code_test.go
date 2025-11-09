package test008 // want "KTN-TEST-008: le fichier de test 'code_test.go' doit se terminer par '_internal_test.go' ou '_external_test.go'"

import "testing"

// TestSomething devrait être dans un fichier _internal_test.go ou _external_test.go
func TestSomething(t *testing.T) {
	t.Log("This file should have a proper suffix")
}
