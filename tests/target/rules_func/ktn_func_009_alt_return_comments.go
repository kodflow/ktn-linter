package rules_func

// findMaxWithComments trouve le maximum avec commentaires sur tous les returns.
//
// Params:
//   - values: la liste de valeurs
//
// Returns:
//   - int: la valeur maximale
//   - bool: true si au moins une valeur existe
func findMaxWithComments(values []int) (int, bool) {
	if len(values) == 0 {
		// Retourne 0 et false car la liste est vide
		return 0, false
	}

	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	// Retourne la valeur maximale trouvée et true
	return max, true
}

// isValidWithComments vérifie la validité avec commentaires sur tous les returns.
//
// Params:
//   - value: la valeur à vérifier
//
// Returns:
//   - bool: true si valide
func isValidWithComments(value int) bool {
	if value < 0 {
		// Retourne false car la valeur est négative
		return false
	}

	if value > 100 {
		// Retourne false car la valeur dépasse 100
		return false
	}

	// Retourne true car la valeur est dans la plage valide [0, 100]
	return true
}
