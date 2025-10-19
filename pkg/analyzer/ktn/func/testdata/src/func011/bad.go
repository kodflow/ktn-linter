package func011

// Bad: Missing comments on branches and returns

// BadIfNoComment has an if statement without a comment
func BadIfNoComment(x int) bool {
	if x > 0 { // want "KTN-FUNC-011.*if.*commentaire"
		return true // want "KTN-FUNC-011.*return.*commentaire"
	}
	return false // want "KTN-FUNC-011.*return.*commentaire"
}

// BadElseNoComment has an else statement without a comment
func BadElseNoComment(x int) bool {
	// Check if x is positive
	if x > 0 {
		// Return true for positive
		return true
	} else { // want "KTN-FUNC-011.*else.*commentaire"
		return false // want "KTN-FUNC-011.*return.*commentaire"
	}
}

// BadSwitchNoComment has a switch without a comment
func BadSwitchNoComment(x int) string {
	switch x { // want "KTN-FUNC-011.*switch.*commentaire"
	case 0: // want "KTN-FUNC-011.*case.*commentaire"
		return "zero" // want "KTN-FUNC-011.*return.*commentaire"
	case 1: // want "KTN-FUNC-011.*case.*commentaire"
		return "one" // want "KTN-FUNC-011.*return.*commentaire"
	default: // want "KTN-FUNC-011.*case.*commentaire"
		return "other" // want "KTN-FUNC-011.*return.*commentaire"
	}
}

// BadLoopNoComment has a loop without a comment
func BadLoopNoComment(items []int) int {
	sum := 0
	for _, item := range items { // want "KTN-FUNC-011.*boucle.*commentaire"
		sum += item
	}
	return sum // want "KTN-FUNC-011.*return.*commentaire"
}

// BadReturnNoComment has returns without comments
func BadReturnNoComment(x int) int {
	return x * 2 // want "KTN-FUNC-011.*return.*commentaire"
}

// BadMixedComments has some comments but not all
func BadMixedComments(x int) bool {
	// Check if x is positive
	if x > 0 {
		return true // want "KTN-FUNC-011.*return.*commentaire"
	}
	// Check if x is zero
	if x == 0 {
		// Return false for zero
		return false
	}
	return true // want "KTN-FUNC-011.*return.*commentaire"
}

// BadTypeSwitchNoComment has a type switch without comments
func BadTypeSwitchNoComment(v interface{}) string {
	switch v.(type) { // want "KTN-FUNC-011.*switch.*commentaire"
	case string: // want "KTN-FUNC-011.*case.*commentaire"
		return "string" // want "KTN-FUNC-011.*return.*commentaire"
	case int: // want "KTN-FUNC-011.*case.*commentaire"
		return "int" // want "KTN-FUNC-011.*return.*commentaire"
	default: // want "KTN-FUNC-011.*case.*commentaire"
		return "unknown" // want "KTN-FUNC-011.*return.*commentaire"
	}
}

// BadElseIfNoComment has an else if without a comment
func BadElseIfNoComment(x int) string {
	// Check if x is negative
	if x < 0 {
		// Return "negative" for negative values
		return "negative"
	} else if x > 0 { // want "KTN-FUNC-011.*if.*commentaire" "KTN-FUNC-011.*else.*commentaire"
		return "positive" // want "KTN-FUNC-011.*return.*commentaire"
	} else {
		// Return "zero" for zero value
		return "zero"
	}
}

// BadElseNoCommentInside has an else with comment before but checking inside path
func BadElseNoCommentInside(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else { // want "KTN-FUNC-011.*else.*commentaire"
		return "non-positive" // want "KTN-FUNC-011.*return.*commentaire"
	}
}

// BadElseEmptyBlock has an empty else block without comment
func BadElseEmptyBlock(x int) {
	// Check if x is positive
	if x > 0 {
		// Return for positive values
		return
	} else { // want "KTN-FUNC-011.*else.*commentaire"
	}
}

// BadElseBlockNoCommentBeforeFirstStmt has else block without comment before or inside
func BadElseBlockNoCommentBeforeFirstStmt(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else { // want "KTN-FUNC-011.*else.*commentaire"
		x = 0
		return "zero" // want "KTN-FUNC-011.*return.*commentaire"
	}
}
