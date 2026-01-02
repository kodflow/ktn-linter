package generic006

// badMin retourne le minimum de deux valeurs.
// Erreur: utilise < sans contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: valeur minimum
func badMin[T any](a, b T) T {
	// Comparaison < sans cmp.Ordered
	if a < b {
		// Retourne a
		return a
	}
	// Retourne b
	return b
}

// badMax retourne le maximum de deux valeurs.
// Erreur: utilise > sans contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: valeur maximum
func badMax[T any](a, b T) T {
	// Comparaison > sans cmp.Ordered
	if a > b {
		// Retourne a
		return a
	}
	// Retourne b
	return b
}

// badSum calcule la somme d'un slice.
// Erreur: utilise + sans contrainte cmp.Ordered.
//
// Params:
//   - values: slice de valeurs
//
// Returns:
//   - T: somme des valeurs
func badSum[T any](values ...T) T {
	// Initialise la somme
	var sum T
	// Parcourt les valeurs
	for _, v := range values {
		// Addition sans contrainte
		sum = sum + v
	}
	// Retourne la somme
	return sum
}

// badAverage calcule la moyenne simplifiee.
// Erreur: utilise + sans contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: somme des valeurs
func badAverage[T any](a, b T) T {
	// Addition sans contrainte
	result := a + b
	// Retourne le resultat
	return result
}

// badDiff calcule la difference de deux valeurs.
// Erreur: utilise - sans contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: difference des valeurs
func badDiff[T any](a, b T) T {
	// Soustraction sans contrainte
	return a - b
}

// badProduct calcule le produit de deux valeurs.
// Erreur: utilise * sans contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - T: produit des valeurs
func badProduct[T any](a, b T) T {
	// Multiplication sans contrainte
	return a * b
}

// badModulo calcule le reste de la division.
// Erreur: utilise % sans contrainte cmp.Ordered.
//
// Params:
//   - a: dividende
//   - b: diviseur
//
// Returns:
//   - T: reste de la division
func badModulo[T any](a, b T) T {
	// Modulo sans contrainte
	return a % b
}

// badLessOrEqual compare deux valeurs.
// Erreur: utilise <= sans contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - bool: true si a <= b
func badLessOrEqual[T any](a, b T) bool {
	// Comparaison <= sans contrainte
	return a <= b
}

// badGreaterOrEqual compare deux valeurs.
// Erreur: utilise >= sans contrainte cmp.Ordered.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - bool: true si a >= b
func badGreaterOrEqual[T any](a, b T) bool {
	// Comparaison >= sans contrainte
	return a >= b
}
