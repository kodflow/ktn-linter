package rules_func

// deeplyNestedBad a une profondeur d'imbrication trop élevée (4 niveaux).
//
// Params:
//   - value: la valeur d'entrée
//
// Returns:
//   - int: le résultat du calcul
func deeplyNestedBad(value int) int {
	result := 0

	for i := 0; i < value; i++ { // Niveau 1
		if i%2 == 0 { // Niveau 2
			for j := 0; j < 5; j++ { // Niveau 3
				if j%2 == 0 { // Niveau 4 - TROP PROFOND
					result += i + j
				}
			}
		}
	}

	// Early return from function.
	return result
}

// extremelyNestedBad a une profondeur d'imbrication extrême (5 niveaux).
//
// Params:
//   - x: première valeur
//   - y: deuxième valeur
//
// Returns:
//   - int: le résultat
func extremelyNestedBad(x, y int) int {
	count := 0

	for i := 0; i < x; i++ { // Niveau 1
		if i > 0 { // Niveau 2
			for j := 0; j < y; j++ { // Niveau 3
				if j > 0 { // Niveau 4
					for k := 0; k < 3; k++ { // Niveau 5 - TRÈS MAUVAIS
						count++
					}
				}
			}
		}
	}

	// Early return from function.
	return count
}

// complexNestedBad mélange plusieurs types d'imbrication.
//
// Params:
//   - values: liste de valeurs
//
// Returns:
//   - []int: valeurs filtrées
func complexNestedBad(values []int) []int {
	result := make([]int, 0)

	for _, v := range values { // Niveau 1
		if v > 0 { // Niveau 2
			switch v % 3 { // Niveau 3
			case 0:
				if v < 100 { // Niveau 4 - TROP PROFOND
					result = append(result, v)
				}
			case 1:
				if v > 10 { // Niveau 4 - TROP PROFOND
					result = append(result, v*2)
				}
			}
		}
	}

	// Early return from function.
	return result
}
