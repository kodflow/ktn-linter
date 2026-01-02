package generic003

import "golang.org/x/exp/constraints" // want "KTN-GENERIC-003:"

// badMax retourne le maximum de deux valeurs.
// Erreur: utilise le package obsolete x/exp/constraints.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: la plus grande valeur
func badMax[T constraints.Ordered](a, b T) T {
	// Comparer les valeurs
	if a > b {
		// a est plus grand
		return a
	}
	// b est plus grand ou egal
	return b
}

// badMin retourne le minimum de deux valeurs.
// Erreur: utilise le package obsolete x/exp/constraints.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: la plus petite valeur
func badMin[T constraints.Ordered](a, b T) T {
	// Comparer les valeurs
	if a < b {
		// a est plus petit
		return a
	}
	// b est plus petit ou egal
	return b
}
