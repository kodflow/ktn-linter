package generic001

// goodContains cherche une valeur dans un slice avec comparable.
// Correct: utilise la contrainte comparable.
//
// Params:
//   - s: slice a parcourir
//   - v: valeur a chercher
//
// Returns:
//   - bool: true si trouve
func goodContains[T comparable](s []T, v T) bool {
	// Parcourir le slice
	for _, x := range s {
		// Comparaison == avec comparable OK
		if x == v {
			// Valeur trouvee
			return true
		}
	}
	// Valeur non trouvee
	return false
}

// goodMap applique une fonction sur chaque element.
// Correct: pas d'utilisation de == ou !=.
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

// goodFilter filtre un slice sans utiliser ==.
// Correct: utilise une fonction de filtre.
//
// Params:
//   - s: slice source
//   - pred: predicat de filtre
//
// Returns:
//   - []T: slice filtre
func goodFilter[T any](s []T, pred func(T) bool) []T {
	// Creer le slice resultat
	var result []T
	// Parcourir et filtrer
	for _, x := range s {
		// Tester le predicat
		if pred(x) {
			// Ajouter au resultat
			result = append(result, x)
		}
	}
	// Retourner le resultat
	return result
}

// goodReduce reduit un slice sans utiliser ==.
// Correct: pas de comparaison d'egalite.
//
// Params:
//   - s: slice source
//   - init: valeur initiale
//   - f: fonction de reduction
//
// Returns:
//   - U: valeur reduite
func goodReduce[T, U any](s []T, init U, f func(U, T) U) U {
	// Initialiser l'accumulateur
	acc := init
	// Parcourir et reduire
	for _, x := range s {
		// Appliquer la fonction
		acc = f(acc, x)
	}
	// Retourner le resultat
	return acc
}
