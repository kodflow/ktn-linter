package KTN_TEST_004_test

import "testing"

// TestValidatorCreation teste la création du validateur.
func TestValidatorCreation(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("validator should not be nil")
	}
}

// NOTE: Les vraies violations (TestEmail, BenchmarkValidation)
// sont dans validator.go (fichier non-test) - c'est la violation KTN-TEST-004
