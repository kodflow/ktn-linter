package comment002

const (
	// DEFAULT_VALUE valeur par défaut pour les exemples
	DEFAULT_VALUE int = 42
	// MAX_COUNTER valeur maximale du compteur
	MAX_COUNTER int = 100
	// MULTIPLIER multiplicateur pour les calculs
	MULTIPLIER int = 2
)

// goodFunctionWithShortComments demonstrates proper comment length.
func goodFunctionWithShortComments() {
	// Short comment within limit
	x := DEFAULT_VALUE
	_ = x
}

// goodMultipleShortComments has acceptable comments.
func goodMultipleShortComments() {
	// Initialize variable
	y := "test"

	// Set counter
	z := MAX_COUNTER
	_, _ = y, z
}

// goodNoInlineComments has no inline comments.
func goodNoInlineComments() {
	x := 1
	y := MULTIPLIER
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
	return
}
