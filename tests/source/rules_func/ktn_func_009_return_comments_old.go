package rules_func

// findMaxWithoutComments trouve le maximum sans commentaires sur les returns.
//
// Params:
//   - values: la liste de valeurs
//
// Returns:
//   - int: la valeur maximale
//   - bool: true si au moins une valeur existe
func findMaxWithoutComments(values []int) (int, bool) {
	if len(values) == 0 {
		return 0, false
	}

	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	return max, true
}

// isValidWithoutComments vérifie la validité sans commentaires.
//
// Params:
//   - value: la valeur à vérifier
//
// Returns:
//   - bool: true si valide
func isValidWithoutComments(value int) bool {
	if value < 0 {
		return false
	}

	if value > 100 {
		return false
	}

	return true
}
