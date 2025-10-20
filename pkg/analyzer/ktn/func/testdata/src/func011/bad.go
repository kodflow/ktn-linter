package func011

// Bad: Missing comments on branches and returns

// BadIfNoComment has an if statement without a comment
func BadIfNoComment(x int) bool {
	if x > 0 {
		return true
	}
	return false
}

// BadElseNoComment has an else statement without a comment
func BadElseNoComment(x int) bool {
	// Check if x is positive
	if x > 0 {
		// Return true for positive
		return true
	} else {
		return false
	}
}

// BadSwitchNoComment has a switch without a comment
func BadSwitchNoComment(x int) string {
	switch x {
	case 0:
		return "zero"
	case 1:
		return "one"
	default:
		return "other"
	}
}

// BadLoopNoComment has a loop without a comment
func BadLoopNoComment(items []int) int {
	sum := 0
	for _, item := range items {
		sum += item
	}
	return sum
}

// BadReturnNoComment has returns without comments
func BadReturnNoComment(x int) int {
	return x * 2
}

// BadMixedComments has some comments but not all
func BadMixedComments(x int) bool {
	// Check if x is positive
	if x > 0 {
		return true
	}
	// Check if x is zero
	if x == 0 {
		// Return false for zero
		return false
	}
	return true
}

// BadTypeSwitchNoComment has a type switch without comments
func BadTypeSwitchNoComment(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case int:
		return "int"
	default:
		return "unknown"
	}
}

// BadElseIfNoComment has an else if without a comment
func BadElseIfNoComment(x int) string {
	// Check if x is negative
	if x < 0 {
		// Return "negative" for negative values
		return "negative"
	} else if x > 0 {
		return "positive"
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
	} else {
		return "non-positive"
	}
}

// BadElseEmptyBlock has an empty else block without comment
func BadElseEmptyBlock(x int) {
	// Check if x is positive
	if x > 0 {
		// Return for positive values
		return
	} else {
	}
}

// BadElseBlockNoCommentBeforeFirstStmt has else block without comment before or inside
func BadElseBlockNoCommentBeforeFirstStmt(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else {
		x = 0
		return "zero"
	}
}
