// Package var026 contains test cases for KTN-VAR-026.
package var026

// goodMinMax utilise les built-ins min() et max().
func goodMinMax() {
	// Utilisation des built-ins
	x := min(1, 2)
	y := max(3, 4, 5)
	z := min(1, 2, 3, 4, 5)
	// Utilisation des variables
	_ = x + y + z
}

// goodMinInt utilise min() built-in.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - int: valeur minimum
func goodMinInt(a, b int) int {
	// Utilisation du built-in min()
	return min(a, b)
}

// goodMaxInt utilise max() built-in.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - int: valeur maximum
func goodMaxInt(a, b int) int {
	// Utilisation du built-in max()
	return max(a, b)
}

// goodMinFloat utilise min() built-in pour float64.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - float64: valeur minimum
func goodMinFloat(a, b float64) float64 {
	// Utilisation du built-in min()
	return min(a, b)
}

// goodMaxFloat utilise max() built-in pour float64.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - float64: valeur maximum
func goodMaxFloat(a, b float64) float64 {
	// Utilisation du built-in max()
	return max(a, b)
}

// goodCompareLogic utilise une comparaison normale sans pattern min/max.
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - bool: true si a est inferieur a b
func goodCompareLogic(a, b int) bool {
	// Comparaison simple, pas un pattern min/max
	if a < b {
		// a est plus petit
		return true
	}
	// b est plus grand ou egal
	return false
}

// init utilise les fonctions pour eviter KTN-FUNC-004.
func init() {
	// Appel des fonctions
	goodMinMax()
	_ = goodMinInt(1, 2)
	_ = goodMaxInt(1, 2)
	_ = goodMinFloat(1.0, 2.0)
	_ = goodMaxFloat(1.0, 2.0)
	_ = goodCompareLogic(1, 2)
}
