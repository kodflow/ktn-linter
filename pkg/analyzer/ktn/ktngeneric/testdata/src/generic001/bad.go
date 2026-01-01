package generic001

// badContains cherche une valeur dans un slice.
// Erreur: utilise == sans contrainte comparable.
//
// Params:
//   - s: slice a parcourir
//   - v: valeur a chercher
//
// Returns:
//   - bool: true si trouve
func badContains[T any](s []T, v T) bool {
	// Parcourir le slice
	for _, x := range s {
		// Comparaison == sans comparable
		if x == v {
			// Valeur trouvee
			return true
		}
	}
	// Valeur non trouvee
	return false
}

// badIndex cherche l'index d'une valeur.
// Erreur: utilise == sans contrainte comparable.
//
// Params:
//   - s: slice a parcourir
//   - v: valeur a chercher
//
// Returns:
//   - int: index ou -1 si non trouve
func badIndex[T any](s []T, v T) int {
	// Parcourir le slice
	for i, x := range s {
		// Comparaison == sans comparable
		if x == v {
			// Retourne l'index
			return i
		}
	}
	// Non trouve
	return -1
}
