package rules_func

// complexCalculationWithoutInternalComments effectue un calcul complexe sans commentaires internes.
//
// Params:
//   - value: la valeur d'entrée
//
// Returns:
//   - int: le résultat du calcul
//   - error: une erreur si la valeur est invalide
func complexCalculationWithoutInternalComments(value int) (int, error) {
	if value < 0 {
		// Early return from function.
		return 0, nil
	}

	result := 0

	for i := 0; i < value; i++ {
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

	// Early return from function.
	return result, nil
}

// processDataWithoutComments traite des données sans expliquer la logique.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - []int: les données traitées
func processDataWithoutComments(data []int) []int {
	processed := make([]int, 0, len(data))

	for _, v := range data {
		if v > 0 {
			if v%10 == 0 {
				processed = append(processed, v*2)
			} else if v%5 == 0 {
				processed = append(processed, v+10)
			} else if v%2 == 0 {
				processed = append(processed, v/2)
			} else {
				processed = append(processed, v)
			}
		}
	}

	// Early return from function.
	return processed
}
