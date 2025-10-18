package rules_func

// sixParamsIsBad a exactement 6 paramètres (limite dépassée).
//
// Params:
//   - a: premier paramètre
//   - b: deuxième paramètre
//   - c: troisième paramètre
//   - d: quatrième paramètre
//   - e: cinquième paramètre
//   - f: sixième paramètre (TROP)
//
// Returns:
//   - int: résultat
func sixParamsIsBad(a, b, c, d, e, f int) int {
	return a + b + c + d + e + f
}

// sevenParamsIsWorse a 7 paramètres (bien au-delà de la limite).
//
// Params:
//   - a: paramètre 1
//   - b: paramètre 2
//   - c: paramètre 3
//   - d: paramètre 4
//   - e: paramètre 5
//   - f: paramètre 6
//   - g: paramètre 7
//
// Returns:
//   - int: résultat
func sevenParamsIsWorse(a, b, c, d, e, f, g int) int {
	return a + b + c + d + e + f + g
}
