package func007

func TotallyNoDoc() string { // want "KTN-FUNC-007"
	return ""
}

// Missing documentation completely
func NoDoc() string { // want "KTN-FUNC-007"
	return ""
}

// Wrong format - doesn't start with function name
// This is wrong
func WrongFormat() { // want "KTN-FUNC-007"
	// Body
}

// BadMissingParams has params but no Params section
func BadMissingParams(x int) { // want "KTN-FUNC-007"
	// Body
}

// BadMissingReturns has returns but no Returns section
func BadMissingReturns() string { // want "KTN-FUNC-007"
	return ""
}

/* BlockCommentOnly uses block comment instead of line comment */
func BlockCommentOnly() { // want "KTN-FUNC-007"
	// Body
}

// BadEmptyParamsSection a une section Params: vide (sans items)
//
// Params:
//
// Returns:
//   - string: résultat
func BadEmptyParamsSection(x int) string { // want "KTN-FUNC-007"
	return ""
}

// BadEmptyReturnsSection a une section Returns: vide (sans items)
//
// Params:
//   - x: paramètre
//
// Returns:
func BadEmptyReturnsSection(x int) string { // want "KTN-FUNC-007"
	return ""
}
