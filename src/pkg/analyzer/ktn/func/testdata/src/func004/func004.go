package func004

// goodFunc a une documentation Returns valide.
//
// Returns:
//   - int: la valeur calculée
func goodFunc() int {
	return 42
}

// goodFuncWithParams a des params et returns documentés.
//
// Params:
//   - x: la valeur d'entrée
//
// Returns:
//   - int: le résultat
func goodFuncWithParams(x int) int {
	return x * 2
}

// goodMultipleReturns documente plusieurs retours.
//
// Params:
//   - x: premier paramètre
//   - y: second paramètre
//
// Returns:
//   - int: somme
//   - bool: succès
//   - error: erreur éventuelle
func goodMultipleReturns(x, y int) (int, bool, error) {
	return x + y, true, nil
}

// goodNoReturns n'a pas de retour donc pas besoin de Returns:.
func goodNoReturns() {
}

// badMissingReturns teste l'absence de Returns. // want "KTN_FUNC_004.*section 'Returns:'"
func badMissingReturns() int {
	return 0
}

// badMissingReturnsWithError manque Returns. // want "KTN_FUNC_004.*section 'Returns:'"
//
// Params:
//   - x: valeur
func badMissingReturnsWithError(x int) error {
	return nil
}

// badMultipleReturnsMissing manque Returns. // want "KTN_FUNC_004.*section 'Returns:'"
//
// Params:
//   - a: premier
//   - b: second
func badMultipleReturnsMissing(a, b int) (int, error) {
	return a + b, nil
}
