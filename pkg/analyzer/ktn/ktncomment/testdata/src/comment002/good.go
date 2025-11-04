package comment002

// goodFunctionWithShortComments demonstrates proper comment length.
func goodFunctionWithShortComments() {
	// Short comment within limit
	x := 42
	_ = x
}

// goodMultipleShortComments has acceptable comments.
func goodMultipleShortComments() {
	// Initialize variable
	y := "test"

	// Set counter
	z := 100
	_, _ = y, z
}

// goodNoInlineComments has no inline comments.
func goodNoInlineComments() {
	x := 1
	y := 2
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
