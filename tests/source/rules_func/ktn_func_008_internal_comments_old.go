package rules_func

// complexFunctionWithoutComments est une fonction complexe sans commentaires internes.
//
// Returns:
//   - int: résultat du calcul
func complexFunctionWithoutComments() int {
	result := 0

	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			if i%3 == 0 {
				if i%5 == 0 {
					result += i * 2
				} else {
					result += i
				}
			} else {
				result -= i
			}
		} else {
			if i%7 == 0 {
				result += i * 3
			} else {
				result -= i * 2
			}
		}
	}

	// Retourne le résultat calculé
	return result
}
