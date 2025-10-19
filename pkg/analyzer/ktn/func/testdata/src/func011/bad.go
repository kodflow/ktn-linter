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
