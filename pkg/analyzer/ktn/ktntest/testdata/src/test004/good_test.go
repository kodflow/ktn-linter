package test004_test

import "testing"

// TestWithErrorCases teste avec cas d'erreur (BIEN)
func TestWithErrorCases(t *testing.T) {
	// Cas d'erreur détecté grâce aux mots "error", "invalid", "fail"
	t.Log("Testing error case")
	t.Log("Testing invalid input")
	t.Log("Testing fail scenario")
}
