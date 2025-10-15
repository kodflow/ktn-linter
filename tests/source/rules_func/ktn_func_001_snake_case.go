package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-001: Nom pas en MixedCaps/mixedCaps (snake_case interdit)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1: snake_case (SEULE ERREUR: KTN-FUNC-001)
// NOTE: Tout est parfait (commentaire + params OK) SAUF nom en snake_case
// ERREUR ATTENDUE: KTN-FUNC-001 sur parse_http_requestF001

// parse_http_requestF001 parse une requête HTTP.
//
// Params:
//   - data: la chaîne de données à parser
//
// Returns:
//   - error: une erreur si les données sont vides
func parse_http_requestF001(data string) error {
	if data == "" {
		return errors.New("data vide")
	}
	return nil
}

// ❌ CAS INCORRECT 2: Snake_Case mixte (SEULE ERREUR: KTN-FUNC-001)
// ERREUR ATTENDUE: KTN-FUNC-001 sur Calculate_TotalF001

// Calculate_TotalF001 calcule le total.
//
// Params:
//   - values: le tableau d'entiers à sommer
//
// Returns:
//   - int: la somme de tous les entiers du tableau
func Calculate_TotalF001(values []int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}
