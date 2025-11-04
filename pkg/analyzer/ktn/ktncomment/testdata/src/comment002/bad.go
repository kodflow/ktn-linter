package comment002

// badFunctionWithLongComment demonstrates a function with comments.
func badFunctionWithLongComment() {
	// This is a very long inline comment that exceeds the maximum allowed length of 80 characters and should trigger the rule // want "KTN-COMMENT-002"
	x := 42
	_ = x
}

// badMultipleLongComments has multiple violations.
func badMultipleLongComments() {
	// This comment is way too long and contains unnecessary verbose explanations about what the code does below // want "KTN-COMMENT-002"
	y := "test"

	// Another extremely long comment that provides excessive detail about a simple variable assignment operation here // want "KTN-COMMENT-002"
	z := 100
	_, _ = y, z
}

// badNestedLongComment has long comment in nested block.
func badNestedLongComment(flag bool) {
	if flag {
		// This is an excessively long comment inside a conditional block that should be shortened significantly // want "KTN-COMMENT-002"
		return
	}
}
