// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-004: Fonctions de test dans fichier non-test
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//
//	Les fonctions Test*, Benchmark* et Example* doivent être uniquement
//	dans les fichiers *_test.go, jamais dans les fichiers .go normaux.
//
//	POURQUOI :
//	- Évite la compilation de code de test dans le binaire final
//	- Sépare clairement code production et code de test
//	- Évite les imports de testing dans le code production
//	- Prévient les bugs subtils avec go build vs go test
//
// ❌ CAS INCORRECT 1: Fonction TestEmail dans validator.go
// ERREUR ATTENDUE: KTN-TEST-004 sur TestEmail
//
// ❌ CAS INCORRECT 2: Fonction BenchmarkValidation dans validator.go
// ERREUR ATTENDUE: KTN-TEST-004 sur BenchmarkValidation
//
// ❌ CAS INCORRECT 3: Fonction ExampleValidator dans validator.go
// ERREUR ATTENDUE: KTN-TEST-004 sur ExampleValidator
//
// ✅ CAS PARFAIT (voir target/) :
//
//	Les tests sont dans validator_test.go, pas dans validator.go
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_004_test

import (
	"regexp"
	"testing"
)

// Validator valide des données.
type Validator struct {
	emailRegex *regexp.Regexp
}

// NewValidator crée un validateur.
//
// Returns:
//   - *Validator: nouvelle instance
func NewValidator() *Validator {
	// Early return from function.
	return &Validator{
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`),
	}
}

// ValidateEmail valide un email.
//
// Params:
//   - email: l'email à valider
//
// Returns:
//   - bool: true si valide
func (v *Validator) ValidateEmail(email string) bool {
	// Early return from function.
	return v.emailRegex.MatchString(email)
}

// ❌ VIOLATION: Fonction de test dans un fichier non-test
// TestEmail ne devrait PAS être ici, mais dans validator_test.go
func TestEmail(t *testing.T) {
	v := NewValidator()
	if !v.ValidateEmail("test@example.com") {
		t.Error("email should be valid")
	}
}

// ❌ VIOLATION: Benchmark dans un fichier non-test
// BenchmarkValidation ne devrait PAS être ici
func BenchmarkValidation(b *testing.B) {
	v := NewValidator()
	for i := 0; i < b.N; i++ {
		v.ValidateEmail("test@example.com")
	}
}

// ❌ VIOLATION: Example dans un fichier non-test
// ExampleValidator ne devrait PAS être ici
func ExampleValidator() {
	v := NewValidator()
	valid := v.ValidateEmail("test@example.com")
	if valid {
		println("Email is valid")
	}
	// Output: Email is valid
}
