// Package comment001 provides good test cases.
package comment001

const (
	// defaultValue valeur par défaut pour les exemples
	defaultValue int = 42
	// maxCounter valeur maximale du compteur
	maxCounter int = 100
	// multiplier multiplicateur pour les calculs
	multiplier int = 2
)

// goodFunctionWithShortComments demonstrates proper comment length.
func goodFunctionWithShortComments() {
	// Short comment within limit
	x := defaultValue
	_ = x
}

// goodMultipleShortComments has acceptable comments.
func goodMultipleShortComments() {
	// Initialize variable
	y := "test"

	// Set counter
	z := maxCounter
	_, _ = y, z
}

// goodNoInlineComments has no inline comments.
func goodNoInlineComments() {
	x := 1
	y := multiplier
	_ = x + y
}

// goodExactly80Chars has comment exactly at limit.
func goodExactly80Chars() {
	// This comment is exactly eighty characters long including spaces and punct.
	value := true
	_ = value
}

// goodMultiLineDocComment has proper doc comment.
// This is a multi-line documentation comment which is allowed
// to span multiple lines without triggering the rule since
// each individual line stays within the 80 character limit.
func goodMultiLineDocComment() {
	// Retour de la fonction
	return
}

// init utilise les fonctions privées
func init() {
	// Appel de goodFunctionWithShortComments
	goodFunctionWithShortComments()
	// Appel de goodMultipleShortComments
	goodMultipleShortComments()
	// Appel de goodNoInlineComments
	goodNoInlineComments()
	// Appel de goodExactly80Chars
	goodExactly80Chars()
	// Appel de goodMultiLineDocComment
	goodMultiLineDocComment()
}
