// Bad examples for the comment006 test case.
package comment006

// badNoDoc n'a aucune documentation.
func badNoDoc() string { // want "KTN-COMMENT-006: la fonction 'badNoDoc' doit avoir une documentation"
	// Retourne une chaîne vide
	return ""
}

// Missing documentation - doesn't start with function name
// This is wrong format
func badWrongFormat() { // want "KTN-COMMENT-006: la description doit commencer par '// badWrongFormat '"
	// Body
}

// badMissingParams has params but no Params section.
func badMissingParams(_x int) { // want "KTN-COMMENT-006: section 'Params:' manquante"
	// Paramètre non utilisé
}

// badMissingReturns has returns but no Returns section.
func badMissingReturns() string { // want "KTN-COMMENT-006: section 'Returns:' manquante"
	// Retourne une chaîne vide
	return ""
}

// badEmptyParamsSection a une section Params: vide (sans items).
//
// Params:
//
// Returns:
//   - string: résultat
func badEmptyParamsSection(_x int) string { // want "KTN-COMMENT-006: au moins un paramètre doit être documenté dans 'Params:'"
	// Retourne une chaîne vide
	return ""
}

// badEmptyReturnsSection a une section Returns: vide (sans items).
//
// Params:
//   - _x: paramètre non utilisé
//
// Returns:
func badEmptyReturnsSection(_x int) string { // want "KTN-COMMENT-006: au moins une valeur de retour doit être documentée dans 'Returns:'"
	// Retourne une chaîne vide
	return ""
}

// badMissingParamInDoc documente seulement un paramètre sur deux.
//
// Params:
//   - x: premier paramètre
//
// Returns:
//   - int: résultat
func badMissingParamInDoc(x, y int) int { // Note: La règle actuelle ne vérifie pas le nombre de params
	// Utilisation des paramètres
	return x + y
}

// badMissingReturnInDoc documente seulement un retour sur deux.
//
// Params:
//   - x: paramètre
//
// Returns:
//   - int: résultat
func badMissingReturnInDoc(x int) (int, error) { // Note: La règle actuelle ne vérifie pas le nombre de retours
	// Retourne le résultat
	return x, nil
}

// init appelle toutes les fonctions bad pour éviter le dead code.
func init() {
	// Appels des fonctions bad
	_ = badNoDoc()
	badWrongFormat()
	badMissingParams(1)
	_ = badMissingReturns()
	_ = badEmptyParamsSection(1)
	_ = badEmptyReturnsSection(1)
	_ = badMissingParamInDoc(1, 1)
	_, _ = badMissingReturnInDoc(1)
}
