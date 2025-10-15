package rules_func

import (
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-001: Nom en MixedCaps/mixedCaps (pas de snake_case)
// ════════════════════════════════════════════════════════════════════════════

// ParseHTTPRequestF001Good parse une requête HTTP.
//
// Params:
//   - data: la chaîne de données à parser
//
// Returns:
//   - error: une erreur si les données sont vides
func ParseHTTPRequestF001Good(data string) error {
	if data == "" {
		return errors.New("data vide")
	}
	return nil
}

// CalculateTotalF001Good calcule le total.
//
// Params:
//   - values: le tableau d'entiers à sommer
//
// Returns:
//   - int: la somme de tous les entiers du tableau
func CalculateTotalF001Good(values []int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}
