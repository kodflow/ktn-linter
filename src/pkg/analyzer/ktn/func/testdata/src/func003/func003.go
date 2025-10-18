package func003

// goodNoParams n'a pas de paramètres.
func goodNoParams() {}

// goodWithParams a des paramètres documentés.
//
// Params:
//   - x: premier paramètre
//   - y: second paramètre
func goodWithParams(x, y int) {}

// goodFullDoc a une documentation complète.
//
// Params:
//   - data: les données à traiter
//   - count: nombre d'éléments
//
// Returns:
//   - bool: succès
func goodFullDoc(data string, count int) bool {
	return true
}

// badMissingParams manque la section Params. // want "KTN-FUNC-003.*section 'Params:'"
func badMissingParams(x int) {}

// badUndocumentedParam ne documente pas tous les params. // want "KTN-FUNC-003.*non documentés"
//
// Params:
//   - x: premier paramètre
func badUndocumentedParam(x, y int) {}

// badMissingMultiple manque plusieurs params. // want "KTN-FUNC-003.*non documentés"
//
// Params:
//   - x: seulement x
func badMissingMultiple(x, y, z int) {}
