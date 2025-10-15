package rules_func

// ComplexFunctionWithComments est une fonction complexe avec commentaires internes.
//
// Returns:
//   - int: résultat du calcul
func ComplexFunctionWithComments() int {
	// Initialiser le résultat
	result := 0

	// Parcourir tous les nombres de 0 à 99
	for i := 0; i < 100; i++ {
		// Traiter les nombres pairs
		if i%2 == 0 {
			result = processEvenNumber(i, result)
		} else {
			result = processOddNumber(i, result)
		}
	}

	// Retourne le résultat calculé
	return result
}

// processEvenNumber traite un nombre pair.
//
// Params:
//   - i: le nombre à traiter
//   - result: le résultat accumulé
//
// Returns:
//   - int: le nouveau résultat
func processEvenNumber(i int, result int) int {
	// Vérifier si divisible par 3
	if i%3 == 0 {
		// Cas spécial: divisible par 5
		if i%5 == 0 {
			return result + i*2
		}
		return result + i
	}
	return result - i
}

// processOddNumber traite un nombre impair.
//
// Params:
//   - i: le nombre à traiter
//   - result: le résultat accumulé
//
// Returns:
//   - int: le nouveau résultat
func processOddNumber(i int, result int) int {
	// Divisible par 7: multiplier par 3
	if i%7 == 0 {
		return result + i*3
	}
	// Autres impairs: soustraire le double
	return result - i*2
}
