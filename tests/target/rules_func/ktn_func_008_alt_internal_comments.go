package rules_func

// complexFunctionWithComments est une fonction complexe avec commentaires internes.
//
// Returns:
//   - int: résultat du calcul
func complexFunctionWithComments() int {
	// Initialiser le résultat
	result := 0

	// Parcourir tous les nombres de 0 à 99
	for i := 0; i < 100; i++ {
		// Traiter les nombres pairs
		if i%2 == 0 {
			// Vérifier si divisible par 3
			if i%3 == 0 {
				// Cas spécial: divisible par 5
				if i%5 == 0 {
					result += i * 2
				} else {
					result += i
				}
			} else {
				result -= i
			}
		} else {
			// Traiter les nombres impairs
			if i%7 == 0 {
				// Divisible par 7: multiplier par 3
				result += i * 3
			} else {
				// Autres impairs: soustraire le double
				result -= i * 2
			}
		}
	}

	// Retourne le résultat calculé
	return result
}
