package func002

// processNoParams fonction sans paramètre
func processNoParams() {
	// Fonction vide
}

// calculateOneParam fonction avec 1 paramètre
//
// Params:
//   - a: paramètre de test
func calculateOneParam(a int) {
	// Utilisation du paramètre
	_ = a
}

// buildFiveParams exactement 5 paramètres (à la limite)
//
// Params:
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
func buildFiveParams(a, b, c, d, e int) {
	// Utilisation des paramètres
	_ = a + b + c + d + e
}

// createFiveParamsMixed 5 paramètres de types différents
//
// Params:
//   - a: entier
//   - b: chaîne
//   - c: booléen
//   - d: flottant
//   - e: slice
func createFiveParamsMixed(a int, b string, c bool, d float64, e []int) {
	// Utilisation des paramètres
	_, _, _, _, _ = a, b, c, d, e
}

// MyType structure de test
type MyType struct{}

// processMethodFourParams méthode avec 4 paramètres (le receiver ne compte pas)
//
// Params:
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
func (m MyType) processMethodFourParams(a, b, c, d int) {
	// Utilisation des paramètres
	_ = a + b + c + d
}

// TestWithManyParams les fonctions de test sont exemptées
//
// Params:
//   - t: paramètre de test
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - f: paramètre de test
func TestWithManyParams(t, a, b, c, d, e, f int) {
	// Utilisation des paramètres
	_ = t + a + b + c + d + e + f
}

// BenchmarkWithManyParams les fonctions de benchmark sont exemptées
//
// Params:
//   - b: paramètre de test
//   - a: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - f: paramètre de test
//   - g: paramètre de test
func BenchmarkWithManyParams(b, a, c, d, e, f, g int) {
	// Utilisation des paramètres
	_ = b + a + c + d + e + f + g
}

// ExampleWithManyParams les fonctions d'exemple sont exemptées
//
// Params:
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - f: paramètre de test
//   - g: paramètre de test
func ExampleWithManyParams(a, b, c, d, e, f, g int) {
	// Utilisation des paramètres
	_ = a + b + c + d + e + f + g
}

// FuzzWithManyParams les fonctions de fuzzing sont exemptées
//
// Params:
//   - f: paramètre de test
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - g: paramètre de test
func FuzzWithManyParams(f, a, b, c, d, e, g int) {
	// Utilisation des paramètres
	_ = f + a + b + c + d + e + g
}

// formatThreeParams fonction avec 3 paramètres
//
// Params:
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
func formatThreeParams(a, b, c int) {
	// Utilisation des paramètres
	_ = a + b + c
}

// convertTwoParamsMixed fonction avec 2 paramètres de types différents
//
// Params:
//   - a: entier
//   - b: chaîne
func convertTwoParamsMixed(a int, b string) {
	// Utilisation des paramètres
	_, _ = a, b
}

// validateFourParamsGrouped fonction avec 4 paramètres groupés par type
//
// Params:
//   - a: entier
//   - b: entier
//   - c: chaîne
//   - d: chaîne
func validateFourParamsGrouped(a, b int, c, d string) {
	// Utilisation des paramètres
	_, _, _, _ = a, b, c, d
}

// convertWithVariadic fonction variadique avec 5 paramètres (variadique compte pour 1)
//
// Params:
//   - a: paramètre régulier
//   - b: paramètre régulier
//   - c: paramètre régulier
//   - d: paramètre régulier
//   - e: paramètre variadique
func convertWithVariadic(a, b, c, d int, e ...string) {
	// Utilisation des paramètres
	_, _, _, _, _ = a, b, c, d, e
}

var (
	// goodLiteralUnnamed fonction littérale avec exactement 5 paramètres non nommés
	goodLiteralUnnamed func(int, string, bool, float64, []int) = func(int, string, bool, float64, []int) {
		// Fonction vide
	}

	// goodLiteralFourUnnamed fonction littérale avec 4 paramètres non nommés
	goodLiteralFourUnnamed func(int, int, int, int) = func(int, int, int, int) {
		// Fonction vide
	}

	// goodLiteralOneUnnamed fonction littérale avec 1 paramètre non nommé
	goodLiteralOneUnnamed func(int) = func(int) {
		// Fonction vide
	}

	// goodLiteralNoParams fonction littérale sans paramètre
	goodLiteralNoParams func() = func() {
		// Fonction vide
	}

	// goodLiteralThreeUnnamed fonction littérale avec 3 paramètres non nommés
	goodLiteralThreeUnnamed func(int, string, bool) = func(int, string, bool) {
		// Fonction vide
	}

	// goodLiteralTwoUnnamed fonction littérale avec 2 paramètres non nommés
	goodLiteralTwoUnnamed func(int, string) = func(int, string) {
		// Fonction vide
	}
)
