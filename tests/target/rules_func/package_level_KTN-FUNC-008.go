package rules_func

import "errors"

// complexCalculationWithInternalComments effectue un calcul complexe avec commentaires internes.
//
// Params:
//   - value: la valeur d'entrée
//
// Returns:
//   - int: le résultat du calcul
//   - error: une erreur si la valeur est invalide
func complexCalculationWithInternalComments(value int) (int, error) {
	// Validation: rejet des valeurs négatives
	if value < 0 {
		return 0, errors.New("value must be non-negative")
	}

	result := 0

	// Calcul basé sur des règles métier spécifiques
	// Les multiples de 2, 3, 5, 7 ont des traitements différents
	for i := 0; i < value; i++ {
		if i%2 == 0 {
			// Traitement des nombres pairs
			if i%3 == 0 {
				if i%5 == 0 {
					// Bonus pour les multiples de 30 (2*3*5)
					result += i * 2
				} else {
					// Multiples de 6 seulement
					result += i
				}
			} else {
				// Pairs non multiples de 3
				result -= i
			}
		} else {
			// Traitement des nombres impairs
			if i%7 == 0 {
				// Bonus triple pour les multiples de 7
				result += i * 3
			} else {
				// Impairs standards: pénalité double
				result -= i * 2
			}
		}
	}

	return result, nil
}

// processDataWithComments traite des données avec explications claires.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - []int: les données traitées
func processDataWithComments(data []int) []int {
	processed := make([]int, 0, len(data))

	// Filtrage et transformation selon des règles métier
	for _, v := range data {
		// Ignorer les valeurs négatives ou nulles
		if v > 0 {
			if v%10 == 0 {
				// Doubler les multiples de 10
				processed = append(processed, v*2)
			} else if v%5 == 0 {
				// Ajouter 10 aux multiples de 5
				processed = append(processed, v+10)
			} else if v%2 == 0 {
				// Diviser par 2 les nombres pairs
				processed = append(processed, v/2)
			} else {
				// Conserver les impairs non multiples de 5
				processed = append(processed, v)
			}
		}
	}

	return processed
}
