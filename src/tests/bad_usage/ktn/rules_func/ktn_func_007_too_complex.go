package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-007: Complexité cyclomatique trop élevée (≥ 10)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1: Complexité ≥ 10 (SEULE ERREUR: KTN-FUNC-007)
// NOTE: Tout est parfait (nom + commentaire + params + longueur OK) SAUF complexité
// ERREUR ATTENDUE: KTN-FUNC-007 sur ValidateComplexInputF007

// ValidateComplexInputF007 valide des données complexes avec de nombreuses conditions.
//
// Params:
//   - input: la chaîne de caractères à valider
//   - level: le niveau de validation à appliquer
//
// Returns:
//   - error: une erreur si la validation échoue
func ValidateComplexInputF007(input string, level int) error {
	if input == "" {
		return errors.New("vide")
	}
	if level > 0 && len(input) < 3 {
		return errors.New("court")
	}
	if level > 1 && len(input) > 100 {
		return errors.New("long")
	}
	if level > 2 && input[0] == ' ' {
		return errors.New("espace début")
	}
	if level > 3 && input[len(input)-1] == ' ' {
		return errors.New("espace fin")
	}
	if level > 4 {
		for _, c := range input {
			if c == '\n' {
				return errors.New("newline")
			}
		}
	}
	if level > 5 {
		for _, c := range input {
			if c == '\t' {
				return errors.New("tab")
			}
		}
	}
	return nil
}
