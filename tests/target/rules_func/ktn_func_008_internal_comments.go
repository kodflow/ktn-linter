package rules_func

import "errors"

// ComplexCalculationWithInternalComments effectue un calcul complexe avec commentaires internes.
//
// Params:
//   - value: la valeur d'entrée
//
// Returns:
//   - int: le résultat du calcul
//   - error: une erreur si la valeur est invalide
func ComplexCalculationWithInternalComments(value int) (int, error) {
	// Validation: rejet des valeurs négatives
	if value < 0 {
		return 0, errors.New("value must be non-negative")
	}

	result := 0

	// Calcul basé sur des règles métier spécifiques
	// Les multiples de 2, 3, 5, 7 ont des traitements différents
	for i := 0; i < value; i++ {
		result += calculateValueContribution(i)
	}

	return result, nil
}

// calculateValueContribution calcule la contribution d'une valeur au résultat.
//
// Params:
//   - i: la valeur à évaluer
//
// Returns:
//   - int: la contribution au résultat
func calculateValueContribution(i int) int {
	if i%2 == 0 {
		return calculateEvenContribution(i)
	}
	return calculateOddContribution(i)
}

// calculateEvenContribution calcule la contribution pour les nombres pairs.
//
// Params:
//   - i: la valeur paire
//
// Returns:
//   - int: la contribution
func calculateEvenContribution(i int) int {
	// Traitement des nombres pairs
	if i%3 == 0 && i%5 == 0 {
		// Bonus pour les multiples de 30 (2*3*5)
		return i * 2
	}
	if i%3 == 0 {
		// Multiples de 6 seulement
		return i
	}
	// Pairs non multiples de 3
	return -i
}

// calculateOddContribution calcule la contribution pour les nombres impairs.
//
// Params:
//   - i: la valeur impaire
//
// Returns:
//   - int: la contribution
func calculateOddContribution(i int) int {
	// Traitement des nombres impairs
	if i%7 == 0 {
		// Bonus triple pour les multiples de 7
		return i * 3
	}
	// Impairs standards: pénalité double
	return -i * 2
}

// ProcessDataWithComments traite des données avec explications claires.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - []int: les données traitées
func ProcessDataWithComments(data []int) []int {
	processed := make([]int, 0, len(data))

	// Filtrage et transformation selon des règles métier
	for _, v := range data {
		if transformed, ok := transformValue(v); ok {
			processed = append(processed, transformed)
		}
	}

	return processed
}

// transformValue transforme une valeur selon les règles métier.
//
// Params:
//   - v: la valeur à transformer
//
// Returns:
//   - int: la valeur transformée
//   - bool: true si la valeur doit être incluse
func transformValue(v int) (int, bool) {
	// Ignorer les valeurs négatives ou nulles
	if v <= 0 {
		return 0, false
	}

	if v%10 == 0 {
		// Doubler les multiples de 10
		return v * 2, true
	}
	if v%5 == 0 {
		// Ajouter 10 aux multiples de 5
		return v + 10, true
	}
	if v%2 == 0 {
		// Diviser par 2 les nombres pairs
		return v / 2, true
	}
	// Conserver les impairs non multiples de 5
	return v, true
}
