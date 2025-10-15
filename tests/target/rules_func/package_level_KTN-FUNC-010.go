package rules_func

// deeplyNestedGood a une profondeur d'imbrication acceptable (max 3 niveaux).
//
// Params:
//   - value: la valeur d'entrée
//
// Returns:
//   - int: le résultat du calcul
func deeplyNestedGood(value int) int {
	result := 0

	for i := 0; i < value; i++ { // Niveau 1
		if i%2 == 0 { // Niveau 2
			// Extraction en sous-fonction pour éviter niveau 4
			result += processInnerLoop(i) // Niveau 3
		}
	}

	return result
}

// processInnerLoop sous-fonction extraite pour réduire l'imbrication.
//
// Params:
//   - i: valeur à traiter
//
// Returns:
//   - int: résultat du traitement
func processInnerLoop(i int) int {
	sum := 0
	for j := 0; j < 5; j++ {
		if j%2 == 0 {
			sum += i + j
		}
	}
	return sum
}

// extremelyNestedGood limite l'imbrication avec extraction de fonctions.
//
// Params:
//   - x: première valeur
//   - y: deuxième valeur
//
// Returns:
//   - int: le résultat
func extremelyNestedGood(x, y int) int {
	count := 0

	for i := 0; i < x; i++ { // Niveau 1
		if i > 0 { // Niveau 2
			// Extraction de la logique imbriquée
			count += countInner(y) // Niveau 3
		}
	}

	return count
}

// countInner sous-fonction extraite pour réduire l'imbrication.
//
// Params:
//   - y: limite de la boucle
//
// Returns:
//   - int: nombre d'itérations
func countInner(y int) int {
	count := 0
	for j := 0; j < y; j++ {
		if j > 0 {
			for k := 0; k < 3; k++ {
				count++
			}
		}
	}
	return count
}

// complexNestedGood utilise des fonctions helper pour limiter l'imbrication.
//
// Params:
//   - values: liste de valeurs
//
// Returns:
//   - []int: valeurs filtrées
func complexNestedGood(values []int) []int {
	result := make([]int, 0)

	for _, v := range values { // Niveau 1
		if v > 0 { // Niveau 2
			// Extraction de la logique du switch
			if processed, ok := processValue(v); ok { // Niveau 3
				result = append(result, processed)
			}
		}
	}

	return result
}

// processValue sous-fonction pour le traitement basé sur modulo.
//
// Params:
//   - v: valeur à traiter
//
// Returns:
//   - int: valeur traitée
//   - bool: true si la valeur doit être conservée
func processValue(v int) (int, bool) {
	switch v % 3 {
	case 0:
		if v < 100 {
			return v, true
		}
	case 1:
		if v > 10 {
			return v * 2, true
		}
	}
	return 0, false
}
