// Package var007 contains test cases for KTN rules.
package var007

// badStringConcatInForLoop concatenates strings in a for loop.
//
// Params:
//   - items: slice of strings
//
// Returns:
//   - string: concatenated result
func badStringConcatInForLoop(items []string) string {
	result := ""

	// Bad: string concatenation in loop
	for _, item := range items {
		result += item
	}

	// Return the result
	return result
}

// badStringConcatWithSeparator concatenates with separator in loop.
//
// Params:
//   - items: slice of strings
//
// Returns:
//   - string: concatenated result
func badStringConcatWithSeparator(items []string) string {
	result := ""

	// Bad: string concatenation in loop
	for i, item := range items {
		// Check if separator is needed
		if i > 0 {
			result += ", "
		}
		result += item
	}

	// Return the result
	return result
}

// badNestedLoopConcat concatenates in nested loop.
//
// Params:
//   - matrix: 2D slice of strings
//
// Returns:
//   - string: concatenated result
func badNestedLoopConcat(matrix [][]string) string {
	result := ""

	// Iteration over rows
	for _, row := range matrix {
		// Bad: string concatenation in nested loop
		for _, cell := range row {
			result += cell
		}
	}

	// Return the result
	return result
}

// badClassicForLoop uses classic for loop with concatenation.
//
// Params:
//   - n: number of iterations
//
// Returns:
//   - string: concatenated result
func badClassicForLoop(n int) string {
	result := ""

	// Bad: string concatenation in classic for loop
	for i := 0; i < n; i++ {
		result += "x"
	}

	// Return the result
	return result
}

// init utilise les fonctions privées
func init() {
	// Appel de badStringConcatInForLoop
	_ = badStringConcatInForLoop(nil)
	// Appel de badStringConcatWithSeparator
	_ = badStringConcatWithSeparator(nil)
	// Appel de badNestedLoopConcat
	_ = badNestedLoopConcat(nil)
	// Appel de badClassicForLoop
	_ = badClassicForLoop(0)
}
