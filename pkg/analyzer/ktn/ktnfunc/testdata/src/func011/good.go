package func011

// Good: All branches and returns have comments

const (
	// DOUBLE_MULTIPLIER is the multiplier used for doubling values
	DOUBLE_MULTIPLIER int = 2
	// BENCHMARK_LIMIT is the iteration limit for benchmark tests
	BENCHMARK_LIMIT int = 100
)

// GoodExample demonstrates proper commenting.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - bool: true if x is positive, false otherwise
func GoodExample(x int) bool {
	// Check if x is positive
	if x > 0 {
		// Return true because x is positive
		return true
	}
	// Return false because x is not positive
	return false
}

// GoodSwitch demonstrates proper switch commenting.
//
// Params:
//   - x: the integer to determine the result
//
// Returns:
//   - string: "zero", "one", or "other"
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

// GoodLoop demonstrates proper loop commenting.
//
// Params:
//   - items: the slice of integers to sum
//
// Returns:
//   - int: the sum of all items
func GoodLoop(items []int) int {
	sum := 0
	// Iterate through all items to calculate sum
	for _, item := range items {
		sum += item
	}
	// Return the total sum
	return sum
}

// GoodInlineComment demonstrates inline comments on returns.
//
// Params:
//   - x: the integer to double
//
// Returns:
//   - int: double of x
func GoodInlineComment(x int) int {
	// Retour de la fonction
	return x * DOUBLE_MULTIPLIER // Return double of x
}

// GoodMultipleReturns demonstrates multiple returns with comments.
//
// Params:
//   - x: the integer to process
//
// Returns:
//   - int: doubled value or 0
//   - error: nil always
func GoodMultipleReturns(x int) (int, error) {
	// Check for negative input
	if x < 0 {
		// Return error for negative input
		return 0, nil
	}
	// Return success with doubled value
	return x * DOUBLE_MULTIPLIER, nil
}

// GoodTypeSwitch demonstrates proper type switch commenting.
//
// Params:
//   - v: the interface value to check
//
// Returns:
//   - string: the type name or "unknown"
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

// GoodElseIf demonstrates proper else if commenting with early return.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - string: "negative", "positive", or "zero"
func GoodElseIf(x int) string {
	// Check if x is negative
	if x < 0 {
		// Return "negative" for negative values
		return "negative"
	}
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	}
	// Return "zero" for zero value
	return "zero"
}

// GoodElseCommentInside demonstrates else with comment inside block.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - string: "positive" or "non-positive"
func GoodElseCommentInside(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	}
	// Return "non-positive" for non-positive values
	return "non-positive"
}

// GoodInterface defines an interface method (no body needed).
type GoodInterface interface {
	GoodNoBody(x int) bool
}

// GoodElseIfCommentBefore demonstrates else if with comment before the if part.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - string: "negative", "positive", or "zero"
func GoodElseIfCommentBefore(x int) string {
	// Check if x is negative
	if x < 0 {
		// Return "negative" for negative values
		return "negative"
	}
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	}
	// Return "zero" for zero value
	return "zero"
}

// GoodElseBlockWithCommentInside demonstrates else block with comment at the start.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - string: "positive" or "non-positive"
func GoodElseBlockWithCommentInside(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	}
	// Return "non-positive" for non-positive values
	return "non-positive"
}

// TestGoodFunction demonstrates that test functions are skipped by the analyzer.
//
// Params:
//   - t: the test interface
func TestGoodFunction(t interface{}) {
	// This is a test function, it should be skipped
	if true {
		// Retour de la fonction
		return
	}
	// Retour de la fonction
	return
}

// BenchmarkGoodFunction demonstrates that benchmark functions are skipped.
//
// Params:
//   - b: the benchmark interface
func BenchmarkGoodFunction(b interface{}) {
	// This is a benchmark function, it should be skipped
	for i := 0; i < BENCHMARK_LIMIT; i++ {
		_ = i * DOUBLE_MULTIPLIER
	}
	// Retour de la fonction
	return
}

// ExampleGoodFunction demonstrates that example functions are skipped.
func ExampleGoodFunction() {
	// This is an example function, it should be skipped
	if true {
		// Retour de la fonction
		return
	}
	// Retour de la fonction
	return
}

// FuzzGoodFunction demonstrates that fuzz functions are skipped.
//
// Params:
//   - f: the fuzz interface
func FuzzGoodFunction(f interface{}) {
	// This is a fuzz function, it should be skipped
	if true {
		// Retour de la fonction
		return
	}
	// Retour de la fonction
	return
}

// GoodElseEmptyBlock demonstrates empty else block with comment.
//
// Params:
//   - x: the integer to check
func GoodElseEmptyBlock(x int) {
	// Check if x is positive
	if x > 0 {
		// Return for positive values
		return
	}
	// Handle non-positive case with empty block
}

// GoodElseWithCommentBefore demonstrates else with comment before.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - string: classification result
func GoodElseWithCommentBefore(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
		// Handle non-positive case
	} else {
		// Return "non-positive" for non-positive values
		return "non-positive"
	}
}

// GoodElseBlockWithCommentAtStart demonstrates else block with comment at start.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - string: classification result
func GoodElseBlockWithCommentAtStart(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	} else {
		// Handle non-positive case - comment at start of block
		// Return "non-positive" for non-positive values
		return "non-positive"
	}
}

// GoodInlineComment demonstrates inline comments on same line as code.
//
// Params:
//   - x: the integer to check
//
// Returns:
//   - int: doubled value
func GoodInlineComment(x int) int {
	// Check if x is positive
	if x > 0 {
		return x * DOUBLE_MULTIPLIER // Inline comment: double the positive value
	}
	// Return zero for non-positive values
	return 0
}
