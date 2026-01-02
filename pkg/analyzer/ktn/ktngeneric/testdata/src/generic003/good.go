package generic003

import "cmp"

// goodMax retourne le maximum de deux valeurs.
// Correct: utilise le package standard cmp.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: la plus grande valeur
func goodMax[T cmp.Ordered](a, b T) T {
	// Comparer les valeurs
	if a > b {
		// a est plus grand
		return a
	}
	// b est plus grand ou egal
	return b
}

// goodMin retourne le minimum de deux valeurs.
// Correct: utilise le package standard cmp.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: la plus petite valeur
func goodMin[T cmp.Ordered](a, b T) T {
	// Comparer les valeurs
	if a < b {
		// a est plus petit
		return a
	}
	// b est plus petit ou egal
	return b
}

// goodCompare compare deux valeurs avec cmp.Compare.
// Correct: utilise la fonction cmp.Compare standard.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - int: -1, 0 ou 1
func goodCompare[T cmp.Ordered](a, b T) int {
	// Utiliser cmp.Compare
	return cmp.Compare(a, b)
}

// goodNoConstraints n'utilise pas de contraintes ordonnees.
// Correct: pas besoin de constraints.Ordered.
//
// Params:
//   - s: slice source
//   - f: fonction de transformation
//
// Returns:
//   - []U: slice transforme
func goodNoConstraints[T, U any](s []T, f func(T) U) []U {
	// Creer le slice resultat
	result := make([]U, len(s))
	// Parcourir et transformer
	for i, x := range s {
		// Appliquer la fonction
		result[i] = f(x)
	}
	// Retourner le resultat
	return result
}
