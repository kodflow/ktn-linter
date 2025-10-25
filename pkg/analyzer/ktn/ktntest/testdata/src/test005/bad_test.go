package test005_test

import "testing"

// TestWithoutTableDriven a plusieurs assertions sans table-driven (PAS BIEN)
func TestWithoutTableDriven(t *testing.T) {
	// Première assertion
	if false {
		t.Error("error 1")
	}

	// Deuxième assertion
	if false {
		t.Error("error 2")
	}

	// Troisième assertion - déclenche la règle (>= 2 assertions)
	if false {
		t.Error("error 3")
	}
}
