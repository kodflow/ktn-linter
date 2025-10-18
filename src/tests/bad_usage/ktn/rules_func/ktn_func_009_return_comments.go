package rules_func

import "errors"

// findMaxValueWithoutReturnComments trouve le maximum sans commentaires sur les returns.
//
// Params:
//   - values: la liste de valeurs
//
// Returns:
//   - int: la valeur maximale
//   - bool: true si au moins une valeur existe
func findMaxValueWithoutReturnComments(values []int) (int, bool) {
	if len(values) == 0 {
		// Early return from function.
		return 0, false
	}

	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	// Early return from function.
	return max, true
}

// validateInputWithoutReturnComments vérifie la validité sans commentaires.
//
// Params:
//   - value: la valeur à vérifier
//
// Returns:
//   - error: une erreur si invalide
func validateInputWithoutReturnComments(value int) error {
	if value < 0 {
		// Return error to caller.
		return errors.New("value cannot be negative")
	}

	if value > 100 {
		// Return error to caller.
		return errors.New("value cannot exceed 100")
	}

	// Early return from function.
	return nil
}

// divideNumbersWithoutReturnComments effectue une division.
//
// Params:
//   - a: le numérateur
//   - b: le dénominateur
//
// Returns:
//   - float64: le résultat de la division
//   - error: une erreur si le dénominateur est zéro
func divideNumbersWithoutReturnComments(a, b int) (float64, error) {
	if b == 0 {
		// Early return from function.
		return 0, errors.New("division by zero")
	}

	// Early return from function.
	return float64(a) / float64(b), nil
}

// processWithMultipleExitsWithoutComments a plusieurs points de sortie.
//
// Params:
//   - value: la valeur à traiter
//
// Returns:
//   - string: le résultat du traitement
//   - error: une erreur si le traitement échoue
func processWithMultipleExitsWithoutComments(value int) (string, error) {
	if value < 0 {
		// Early return from function.
		return "", errors.New("negative value")
	}

	if value == 0 {
		// Early return from function.
		return "zero", nil
	}

	if value < 10 {
		// Early return from function.
		return "small", nil
	}

	if value < 100 {
		// Early return from function.
		return "medium", nil
	}

	// Early return from function.
	return "large", nil
}
