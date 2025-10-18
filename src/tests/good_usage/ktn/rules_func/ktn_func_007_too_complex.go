package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-007: Fonction avec complexité cyclomatique < 10
// ════════════════════════════════════════════════════════════════════════════

// ValidateInputF007Good valide les données d'entrée (✅ complexité faible).
//
// Params:
//   - input: la chaîne de caractères à valider
//
// Returns:
//   - error: une erreur si la validation échoue
func ValidateInputF007Good(input string) error {
	if input == "" {
		// Retourne une erreur car l'input est vide
		return errors.New("input vide")
	}
	if len(input) < 3 {
		// Retourne une erreur car l'input est trop court
		return errors.New("input trop court")
	}
	if len(input) > 100 {
		// Retourne une erreur car l'input est trop long
		return errors.New("input trop long")
	}
	// Retourne nil car l'input est valide
	return nil
}
