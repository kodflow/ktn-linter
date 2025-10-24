package func007

func TotallyNoDoc() string {
	// Retourne une chaîne vide
	return ""
}

// Missing documentation completely
func NoDoc() string {
	// Retourne une chaîne vide
	return ""
}

// Wrong format - doesn't start with function name
// This is wrong
func WrongFormat() {
	// Body
}

// BadMissingParams has params but no Params section
func BadMissingParams(x int) {
	// Body
}

// BadMissingReturns has returns but no Returns section
func BadMissingReturns() string {
	// Retourne une chaîne vide
	return ""
}

/* BlockCommentOnly uses block comment instead of line comment */
func BlockCommentOnly() {
	// Body
}

// BadEmptyParamsSection a une section Params: vide (sans items)
//
// Params:
//
// Returns:
//   - string: résultat
func BadEmptyParamsSection(x int) string {
	// Retourne une chaîne vide
	return ""
}

// BadEmptyReturnsSection a une section Returns: vide (sans items)
//
// Params:
//   - x: paramètre
//
// Returns:
func BadEmptyReturnsSection(x int) string {
	// Retourne une chaîne vide
	return ""
}
