// ✅ CORRIGÉ: Toutes les fonctions de test sont dans le fichier _test.go
package KTN_TEST_004_GOOD_test

import "testing"

// TestValidatorCreation teste la création du validateur.
//
// Params:
//   - t: instance de test
func TestValidatorCreation(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("validator should not be nil")
	}
}

// TestEmail teste la validation d'email.
// ✅ Cette fonction est maintenant dans validator_test.go (pas dans validator.go)
//
// Params:
//   - t: instance de test
func TestEmail(t *testing.T) {
	v := NewValidator()
	if !v.ValidateEmail("test@example.com") {
		t.Error("email should be valid")
	}
}

// BenchmarkValidation benchmark la validation.
// ✅ Cette fonction est maintenant dans validator_test.go (pas dans validator.go)
//
// Params:
//   - b: instance de benchmark
func BenchmarkValidation(b *testing.B) {
	v := NewValidator()
	for i := 0; i < b.N; i++ {
		v.ValidateEmail("test@example.com")
	}
}
