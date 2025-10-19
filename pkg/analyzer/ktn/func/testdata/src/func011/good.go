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

// GoodElseIf demonstrates proper else if commenting with comment before else if
func GoodElseIf(x int) string {
	// Check if x is negative
	if x < 0 {
		// Return "negative" for negative values
		return "negative"
	// Check if x is positive (comment before else if)
	} else if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else {
		// Return "zero" for zero value
		return "zero"
	}
}

// GoodElseCommentInside demonstrates else with comment inside block
func GoodElseCommentInside(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else {
		// Handle non-positive values
		// Return "non-positive" for non-positive values
		return "non-positive"
	}
}

// GoodNoBody is an interface method declaration (no body, should be skipped)
type GoodInterface interface {
	GoodNoBody(x int) bool
}

// GoodExternalFuncDeclaration demonstrates external function declaration (no body)
// This is used in CGO or external linkage
// Declaration externe sans body
//
//go:noescape
//go:linkname externalFunc runtime.externalFunc
func externalFunc(x int) int

// GoodElseIfCommentBefore demonstrates else if with comment before the if part
func GoodElseIfCommentBefore(x int) string {
	// Check if x is negative
	if x < 0 {
		// Return "negative" for negative values
		return "negative"
	} else
	// Check if x is positive (comment right before if keyword in else if)
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else {
		// Return "zero" for zero value
		return "zero"
	}
}

// GoodElseBlockWithCommentInside demonstrates else block with comment at the start
func GoodElseBlockWithCommentInside(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else {
		// Handle non-positive case
		// Return "non-positive" for non-positive values
		return "non-positive"
	}
}

// TestGoodFunction demonstrates that test functions are skipped by the analyzer
func TestGoodFunction(t interface{}) {
	// This is a test function, it should be skipped
	if true {
		return
	}
	return
}

// BenchmarkGoodFunction demonstrates that benchmark functions are skipped
func BenchmarkGoodFunction(b interface{}) {
	// This is a benchmark function, it should be skipped
	for i := 0; i < 100; i++ {
		_ = i * 2
	}
	return
}

// ExampleGoodFunction demonstrates that example functions are skipped
func ExampleGoodFunction() {
	// This is an example function, it should be skipped
	if true {
		return
	}
	return
}

// FuzzGoodFunction demonstrates that fuzz functions are skipped
func FuzzGoodFunction(f interface{}) {
	// This is a fuzz function, it should be skipped
	if true {
		return
	}
	return
}

// GoodElseEmptyBlock demonstrates empty else block with comment
func GoodElseEmptyBlock(x int) {
	// Check if x is positive
	if x > 0 {
		// Return for positive values
		return
	// Handle non-positive case with empty block
	} else {
	}
}
