package func002

// processSixParams dépasse la limite avec 6 paramètres
//
// Params:
//   - a, b, c, d, e, f: paramètres de test
func processSixParams(a, b, c, d, e, f int) { // want "KTN-FUNC-002: la fonction 'processSixParams' a 6 paramètres \\(max: 5\\)"
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
func calculateSevenParams(a int, b string, c bool, d float64, e []int, f map[string]int, g chan int) { // want "KTN-FUNC-002: la fonction 'calculateSevenParams' a 7 paramètres \\(max: 5\\)"
	// Utilisation des paramètres
	_, _, _, _, _, _, _ = a, b, c, d, e, f, g
}

// buildTenParams fonction avec beaucoup trop de paramètres
//
// Params:
//   - a à j: paramètres de test
func buildTenParams(a, b, c, d, e, f, g, h, i, j int) { // want "KTN-FUNC-002: la fonction 'buildTenParams' a 10 paramètres \\(max: 5\\)"
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
func createEightParams(a int, b, c string, d, e, f bool, g, h float64) { // want "KTN-FUNC-002: la fonction 'createEightParams' a 8 paramètres \\(max: 5\\)"
	// Utilisation des paramètres
	_, _, _, _, _, _, _, _ = a, b, c, d, e, f, g, h
}

// formatSixParams exactement 6 paramètres (juste au-dessus de la limite)
//
// Params:
//   - a, b, c, d, e, f: paramètres de test
func formatSixParams(a, b, c, d, e, f int) { // want "KTN-FUNC-002: la fonction 'formatSixParams' a 6 paramètres \\(max: 5\\)"
	// Utilisation des paramètres
	_ = a + b + c + d + e + f
}

// convertWithVariadicBad fonction variadique avec 6 paramètres au total
//
// Params:
//   - a, b, c, d, e: paramètres réguliers
//   - f: paramètre variadique
func convertWithVariadicBad(a, b, c, d, e int, f ...string) { // want "KTN-FUNC-002: la fonction 'convertWithVariadicBad' a 6 paramètres \\(max: 5\\)"
	// Utilisation des paramètres
	_, _, _, _, _, _ = a, b, c, d, e, f
}

// badLiteralSix fonction littérale avec 6 paramètres
var badLiteralSix = func(a, b, c, d, e, f int) { // want "KTN-FUNC-002: la fonction 'function literal' a 6 paramètres \\(max: 5\\)"
	// Utilisation des paramètres
	_ = a + b + c + d + e + f
}

// badLiteralUnnamed fonction littérale avec 6 paramètres non nommés
var badLiteralUnnamed = func(int, string, bool, float64, []int, map[string]int) { // want "KTN-FUNC-002: la fonction 'function literal' a 6 paramètres \\(max: 5\\)"
	// Fonction vide
}

// badLiteralSixUnnamed fonction littérale avec 6 paramètres non nommés identiques
var badLiteralSixUnnamed = func(int, int, int, int, int, int) { // want "KTN-FUNC-002: la fonction 'function literal' a 6 paramètres \\(max: 5\\)"
	// Fonction vide
}
