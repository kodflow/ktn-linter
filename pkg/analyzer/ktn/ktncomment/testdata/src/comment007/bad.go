// Bad examples for the func009 test case.
package comment007

// Bad: Missing comments on branches and returns

const (
	// MULTIPLIER_TWO is used for doubling values
	MULTIPLIER_TWO int = 2
)

// BadIfNoComment has an if statement without a comment.
//
// Params:
//   - x: the number to check
//
// Returns:
//   - bool: true if x is positive, false otherwise
func BadIfNoComment(x int) bool {
	if x > 0 {
		return true
	}
	return false
}

// BadSwitchNoComment has a switch without a comment.
//
// Params:
//   - x: the number to classify
//
// Returns:
//   - string: textual representation of the number
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

// BadLoopNoComment has a loop without a comment.
//
// Params:
//   - items: slice of integers to sum
//
// Returns:
//   - int: sum of all items
func BadLoopNoComment(items []int) int {
	sum := 0
	for _, item := range items {
		sum += item
	}
	return sum
}

// BadReturnNoComment has returns without comments.
//
// Params:
//   - x: the number to double
//
// Returns:
//   - int: x multiplied by 2
func BadReturnNoComment(x int) int {
	return x * MULTIPLIER_TWO
}

// BadMixedComments has some comments but not all.
//
// Params:
//   - x: the number to check
//
// Returns:
//   - bool: result based on value of x
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

// BadTypeSwitchNoComment has a type switch without comments.
//
// Params:
//   - v: the interface value to check type
//
// Returns:
//   - string: name of the type
func BadTypeSwitchNoComment(v any) string {
	switch v.(type) {
	case string:
		return "string"
	case int:
		return "int"
	default:
		return "unknown"
	}
}

// BadElseIfNoComment has an else if without a comment.
//
// Params:
//   - x: the number to classify
//
// Returns:
//   - string: classification of the number
func BadElseIfNoComment(x int) string {
	// Check if x is negative
	if x < 0 {
		// Return "negative" for negative values
		return "negative"
	}
	if x > 0 {
		return "positive"
	}
	// Return "zero" for zero value
	return "zero"
}

// BadIfNested has nested if without comment.
//
// Params:
//   - x: the number to check
//
// Returns:
//   - string: classification result
func BadIfNested(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	}
	return "non-positive"
}

// BadEmptyAfterIf has code after if without else.
//
// Params:
//   - x: the number to check
func BadEmptyAfterIf(x int) {
	// Check if x is positive
	if x > 0 {
		// Return for positive values
		return
	}
}

// BadBlockAfterIf has block after if without comment.
//
// Params:
//   - x: the number to check
//
// Returns:
//   - string: result string
func BadBlockAfterIf(x int) string {
	// Check if x is positive
	if x > 0 {
		// Return "positive" for positive values
		return "positive"
	}
	x = 0
	return "zero"
}

// BadElseNoComment has else without comment.
//
// Params:
//   - x: the number to check
//
// Returns:
//   - string: classification result
func BadElseNoComment(x int) string {
	var result string
	// Check if x is positive
	if x > 0 {
		result = "positive"
	} else {
		result = "non-positive"
	}
	return result
}

// BadComplexReturn has complex return without comment.
//
// Params:
//   - x: the number to process
//
// Returns:
//   - int: processed value
func BadComplexReturn(x int) int {
	// Check if x is positive
	if x > 0 {
		return x * MULTIPLIER_TWO + 1
	}
	return x - 1
}

// BadMultipleNonTrivialReturns has multiple non-trivial returns without comments.
//
// Params:
//   - x: the number to check
//
// Returns:
//   - string: status message
func BadMultipleNonTrivialReturns(x int) string {
	// Check if x is negative
	if x < 0 {
		return "error: negative"
	}
	// Check if x is zero
	if x == 0 {
		return "warning: zero"
	}
	return "success"
}

// Result contient le résultat d'une opération.
// Utilisé pour tester les returns de struct sans commentaire.
type Result struct {
	Value int
}

// BadReturnStruct has struct return without comment.
//
// Params:
//   - x: the number to wrap
//
// Returns:
//   - Result: result struct
func BadReturnStruct(x int) Result {
	// Check if x is positive
	if x > 0 {
		return Result{Value: x}
	}
	return Result{Value: 0}
}

// BadReturnPointer has pointer return without comment.
//
// Params:
//   - x: the number to wrap
//
// Returns:
//   - *int: pointer to result
func BadReturnPointer(x int) *int {
	result := x * MULTIPLIER_TWO
	// Check if result is positive
	if result > 0 {
		return &result
	}
	zero := 0
	return &zero
}
