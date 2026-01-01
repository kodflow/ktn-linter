// Package var026 contains test cases for KTN-VAR-026.
package var026

import "math"

// badMaxFloat utilise math.Max au lieu de max().
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - float64: valeur maximum
func badMaxFloat(a, b float64) float64 {
	// Utilisation de math.Max au lieu de max() built-in
	return math.Max(a, b) // want "KTN-VAR-026"
}

// badMinFloat utilise math.Min au lieu de min().
//
// Params:
//   - a: premiere valeur
//   - b: deuxieme valeur
//
// Returns:
//   - float64: valeur minimum
func badMinFloat(a, b float64) float64 {
	// Utilisation de math.Min au lieu de min() built-in
	return math.Min(a, b) // want "KTN-VAR-026"
}

// init utilise les fonctions pour eviter KTN-FUNC-004.
func init() {
	// Appel des fonctions
	_ = badMaxFloat(1.0, 2.0)
	_ = badMinFloat(1.0, 2.0)
}
