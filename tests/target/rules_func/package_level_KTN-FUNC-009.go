package rules_func

import "errors"

// findMaxValueWithReturnComments trouve le maximum avec commentaires sur les returns.
//
// Params:
//   - values: la liste de valeurs
//
// Returns:
//   - int: la valeur maximale
//   - bool: true si au moins une valeur existe
func findMaxValueWithReturnComments(values []int) (int, bool) {
	if len(values) == 0 {
		// Retourne valeur par défaut si la liste est vide
		return 0, false
	}

	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	// Retourne le maximum trouvé avec succès
	return max, true
}

// validateInputWithReturnComments vérifie la validité avec commentaires.
//
// Params:
//   - value: la valeur à vérifier
//
// Returns:
//   - error: une erreur si invalide
func validateInputWithReturnComments(value int) error {
	if value < 0 {
		// Retourne erreur si la valeur est négative
		return errors.New("value cannot be negative")
	}

	if value > 100 {
		// Retourne erreur si la valeur dépasse la limite
		return errors.New("value cannot exceed 100")
	}

	// Retourne nil si la validation réussit
	return nil
}

// divideNumbersWithReturnComments effectue une division.
//
// Params:
//   - a: le numérateur
//   - b: le dénominateur
//
// Returns:
//   - float64: le résultat de la division
//   - error: une erreur si le dénominateur est zéro
func divideNumbersWithReturnComments(a, b int) (float64, error) {
	if b == 0 {
		// Retourne erreur car division par zéro impossible
		return 0, errors.New("division by zero")
	}

	// Retourne le résultat de la division
	return float64(a) / float64(b), nil
}

// processWithMultipleExitsWithComments a plusieurs points de sortie commentés.
//
// Params:
//   - value: la valeur à traiter
//
// Returns:
//   - string: le résultat du traitement
//   - error: une erreur si le traitement échoue
func processWithMultipleExitsWithComments(value int) (string, error) {
	if value < 0 {
		// Retourne erreur pour valeur négative
		return "", errors.New("negative value")
	}

	if value == 0 {
		// Retourne "zero" pour valeur nulle
		return "zero", nil
	}

	if value < 10 {
		// Retourne "small" pour petites valeurs (1-9)
		return "small", nil
	}

	if value < 100 {
		// Retourne "medium" pour valeurs moyennes (10-99)
		return "medium", nil
	}

	// Retourne "large" pour grandes valeurs (>=100)
	return "large", nil
}
