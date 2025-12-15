// Package func006 contains test cases for KTN rules.
package func006

import "context"

// processSixParams dépasse la limite avec 6 paramètres
//
// Params:
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - f: paramètre de test
func processSixParams(a, b, c, d, e, f int) {
	// Utilisation des paramètres
	_ = a + b + c + d + e + f
}

// calculateSevenParams dépasse largement avec 7 paramètres
//
// Params:
//   - a: entier
//   - b: chaîne
//   - c: booléen
//   - d: flottant
//   - e: slice d'entiers
//   - f: map
//   - g: canal
func calculateSevenParams(a int, b string, c bool, d float64, e []int, f map[string]int, g chan int) {
	// Utilisation des paramètres
	_, _, _, _, _, _, _ = a, b, c, d, e, f, g
}

// buildTenParams fonction avec beaucoup trop de paramètres
//
// Params:
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - f: paramètre de test
//   - g: paramètre de test
//   - h: paramètre de test
//   - i: paramètre de test
//   - j: paramètre de test
func buildTenParams(a, b, c, d, e, f, g, h, i, j int) {
	// Utilisation des paramètres
	_ = a + b + c + d + e + f + g + h + i + j
}

// createEightParams cas avec 8 paramètres groupés
//
// Params:
//   - a: entier
//   - b, c: chaînes
//   - d, e, f: booléens
//   - g, h: flottants
func createEightParams(a int, b, c string, d, e, f bool, g, h float64) {
	// Utilisation des paramètres
	_, _, _, _, _, _, _, _ = a, b, c, d, e, f, g, h
}

// formatSixParams exactement 6 paramètres (juste au-dessus de la limite)
//
// Params:
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - f: paramètre de test
func formatSixParams(a, b, c, d, e, f int) {
	// Utilisation des paramètres
	_ = a + b + c + d + e + f
}

// convertWithVariadicBad fonction variadique avec 6 paramètres au total
//
// Params:
//   - a: paramètre régulier
//   - b: paramètre régulier
//   - c: paramètre régulier
//   - d: paramètre régulier
//   - e: paramètre régulier
//   - f: paramètre variadique
func convertWithVariadicBad(a, b, c, d, e int, f ...string) {
	// Utilisation des paramètres
	_, _, _, _, _, _ = a, b, c, d, e, f
}

var (
	// badLiteralSix fonction littérale avec 6 paramètres
	badLiteralSix func(int, int, int, int, int, int) = func(a, b, c, d, e, f int) {
		// Utilisation des paramètres
		_ = a + b + c + d + e + f
	}

	// badLiteralUnnamed fonction littérale avec 6 paramètres non nommés
	badLiteralUnnamed func(int, string, bool, float64, []int, map[string]int) = func(int, string, bool, float64, []int, map[string]int) {
		// Fonction vide
	}

	// badLiteralSixUnnamed 6 paramètres non nommés identiques
	badLiteralSixUnnamed func(int, int, int, int, int, int) = func(int, int, int, int, int, int) {
		// Fonction vide
	}
)

// processWithContextBad 7 params, context exclu = 6 effectifs (ERREUR)
//
// Params:
//   - ctx: contexte (exclu du compte)
//   - a: paramètre de test
//   - b: paramètre de test
//   - c: paramètre de test
//   - d: paramètre de test
//   - e: paramètre de test
//   - f: paramètre de test
func processWithContextBad(ctx context.Context, a, b, c, d, e, f int) { // want "KTN-FUNC-006: la fonction 'processWithContextBad' a 6 paramètres \\(max: 5\\)"
	// Utilisation des paramètres
	_, _, _, _, _, _, _ = ctx, a, b, c, d, e, f
}

// init utilise les fonctions privées
func init() {
	// Appel des fonctions
	processSixParams(0, 0, 0, 0, 0, 0)
	calculateSevenParams(0, "", false, 0, nil, nil, nil)
	buildTenParams(0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	createEightParams(0, "", "", false, false, false, 0, 0)
	formatSixParams(0, 0, 0, 0, 0, 0)
	convertWithVariadicBad(0, 0, 0, 0, 0)
	processWithContextBad(nil, 0, 0, 0, 0, 0, 0)
}
