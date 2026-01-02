package generic006

import "cmp"

// goodMin retourne le minimum de deux valeurs.
// Correct: utilise la contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: valeur minimum
func goodMin[T cmp.Ordered](a, b T) T {
	// Comparaison < avec cmp.Ordered OK
	if a < b {
		// Retourne a
		return a
	}
	// Retourne b
	return b
}

// goodMax retourne le maximum de deux valeurs.
// Correct: utilise la contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: valeur maximum
func goodMax[T cmp.Ordered](a, b T) T {
	// Comparaison > avec cmp.Ordered OK
	if a > b {
		// Retourne a
		return a
	}
	// Retourne b
	return b
}

// goodSum calcule la somme d'un slice.
// Correct: utilise la contrainte cmp.Ordered.
//
// Params:
//   - values: slice de valeurs
//
// Returns:
//   - T: somme des valeurs
func goodSum[T cmp.Ordered](values ...T) T {
	// Initialise la somme
	var sum T
	// Parcourt les valeurs
	for _, v := range values {
		// Addition avec cmp.Ordered OK
		sum = sum + v
	}
	// Retourne la somme
	return sum
}

// goodMap applique une fonction sur chaque element.
// Correct: pas d'utilisation d'operateurs ordered.
//
// Params:
//   - s: slice source
//   - f: fonction de transformation
//
// Returns:
//   - []U: slice transforme
func goodMap[T, U any](s []T, f func(T) U) []U {
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

// goodCompareWithFunc compare avec une fonction personnalisee.
// Correct: utilise une fonction de comparaison au lieu d'operateurs.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//   - less: fonction de comparaison
//
// Returns:
//   - T: valeur minimum selon less
func goodCompareWithFunc[T any](a, b T, less func(T, T) bool) T {
	// Utilise la fonction de comparaison
	if less(a, b) {
		// Retourne a
		return a
	}
	// Retourne b
	return b
}

// goodClamp limite une valeur entre min et max.
// Correct: utilise la contrainte cmp.Ordered.
//
// Params:
//   - value: valeur a limiter
//   - minVal: valeur minimum
//   - maxVal: valeur maximum
//
// Returns:
//   - T: valeur limitee
func goodClamp[T cmp.Ordered](value, minVal, maxVal T) T {
	// Verifier si en dessous du minimum
	if value < minVal {
		// Retourne le minimum
		return minVal
	}
	// Verifier si au dessus du maximum
	if value > maxVal {
		// Retourne le maximum
		return maxVal
	}
	// Retourne la valeur
	return value
}
