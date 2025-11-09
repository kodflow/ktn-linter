package test003

import "testing"

// TestSomethingElse teste autre chose (pas les fonctions publiques de bad.go)
func TestSomethingElse(t *testing.T) {
	// Ce test ne correspond à aucune fonction publique
	if 1+1 != 2 {
		t.Error("Math broken")
	}
}
