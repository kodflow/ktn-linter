// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-004: Pas de fonctions de test dans fichier .go (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_004_GOOD_test

import "regexp"

// Validator valide des données.
type Validator struct {
	emailRegex *regexp.Regexp
}

// NewValidator crée un validateur.
//
// Returns:
//   - *Validator: nouvelle instance
func NewValidator() *Validator {
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
	return v.emailRegex.MatchString(email)
}

// ✅ CORRIGÉ: Pas de fonctions Test*, Benchmark* ou Example* ici
// Elles sont toutes dans validator_test.go
