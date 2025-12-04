// Bad examples for the comment002 test case.
package comment002

const (
	// BAD_DEFAULT_VALUE est la valeur par défaut pour bad examples
	BAD_DEFAULT_VALUE int = 42
	// BAD_MAX_THRESHOLD est le seuil maximum pour bad examples
	BAD_MAX_THRESHOLD int = 100
)

// badFunctionWithLongComment demonstrates a function with comments.
//
// Params:
//   - none
//
// Returns:
//   - none
func badFunctionWithLongComment() {
	// This is a very long inline comment that exceeds the maximum allowed length of 80 characters and should trigger the rule // want "KTN-COMMENT-002"
	x := BAD_DEFAULT_VALUE
	_ = x
}

// badMultipleLongComments has multiple violations.
//
// Params:
//   - none
//
// Returns:
//   - none
func badMultipleLongComments() {
	// This comment is way too long and contains unnecessary verbose explanations about what the code does below // want "KTN-COMMENT-002"
	y := "test"

	// Another extremely long comment that provides excessive detail about a simple variable assignment operation here // want "KTN-COMMENT-002"
	z := BAD_MAX_THRESHOLD
	_, _ = y, z
}

// badNestedLongComment has long comment in nested block.
//
// Params:
//   - _flag: boolean flag not used
//
// Returns:
//   - none
func badNestedLongComment(_flag bool) {
	// Vérifie le flag
	if _flag {
		// This is an excessively long comment inside a conditional block that should be shortened significantly // want "KTN-COMMENT-002"
		// Retourne si le flag est true
		return
	}
}

// init utilise les fonctions privées
func init() {
	// Appel de badFunctionWithLongComment
	badFunctionWithLongComment()
	// Appel de badMultipleLongComments
	badMultipleLongComments()
	// Appel de badNestedLongComment
	badNestedLongComment(false)
}
