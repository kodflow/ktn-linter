package func011

// Good: All branches and returns have comments

// GoodExample demonstrates proper commenting
func GoodExample(x int) bool {
	// Check if x is positive
	if x > 0 {
		// Return true because x is positive
		return true
	} else {
		// Return false because x is not positive
		return false
	}
}

// GoodSwitch demonstrates proper switch commenting
func GoodSwitch(x int) string {
	// Determine result based on x value
	switch x {
	// Handle zero case
	case 0:
		// Return "zero" for zero value
		return "zero"
	// Handle positive case
	case 1:
		// Return "one" for one value
		return "one"
	// Handle all other cases
	default:
		// Return "other" for all other values
		return "other"
	}
}

// GoodLoop demonstrates proper loop commenting
func GoodLoop(items []int) int {
	sum := 0
	// Iterate through all items to calculate sum
	for _, item := range items {
		sum += item
	}
	// Return the total sum
	return sum
}

// GoodInlineComment demonstrates inline comments on returns
func GoodInlineComment(x int) int {
	return x * 2 // Return double of x
}

// GoodMultipleReturns demonstrates multiple returns with comments
func GoodMultipleReturns(x int) (int, error) {
	// Check for negative input
	if x < 0 {
		// Return error for negative input
		return 0, nil
	}
	// Return success with doubled value
	return x * 2, nil
}

// GoodTypeSwitch demonstrates proper type switch commenting
func GoodTypeSwitch(v interface{}) string {
	// Determine type of value
	switch v.(type) {
	// Handle string type
	case string:
		// Return "string" for string type
		return "string"
	// Handle int type
	case int:
		// Return "int" for int type
		return "int"
	// Handle all other types
	default:
		// Return "unknown" for unknown types
		return "unknown"
	}
}
