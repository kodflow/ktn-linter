// Good examples for the var012 test case.
package var012

import "strings"

const (
	// AVG_ITEM_LENGTH is the average item length estimate
	AVG_ITEM_LENGTH int = 10
)

// goodStringsBuilder uses strings.Builder for concatenation.
//
// Params:
//   - items: slice of strings
//
// Returns:
//   - string: concatenated result
func goodStringsBuilder(items []string) string {
	// sb is the strings builder
	var sb strings.Builder

	// Good: using strings.Builder
	for _, item := range items {
		sb.WriteString(item)
	}

	// Return the result
	return sb.String()
}

// goodStringsBuilderWithGrow uses strings.Builder with Grow.
//
// Params:
//   - items: slice of strings
//
// Returns:
//   - string: concatenated result
func goodStringsBuilderWithGrow(items []string) string {
	// sb is the strings builder
	var sb strings.Builder
	sb.Grow(len(items) * AVG_ITEM_LENGTH)

	// Good: using strings.Builder with preallocated size
	for _, item := range items {
		sb.WriteString(item)
	}

	// Return the result
	return sb.String()
}

// goodStringsJoin uses strings.Join for simple concatenation.
//
// Params:
//   - items: slice of strings
//
// Returns:
//   - string: concatenated result
func goodStringsJoin(items []string) string {
	// Good: using strings.Join for simple cases
	return strings.Join(items, "")
}

// goodSingleConcat performs single concatenation outside loop.
//
// Params:
//   - a: first string
//   - b: second string
//
// Returns:
//   - string: concatenated result
func goodSingleConcat(a string, b string) string {
	// Good: single concatenation, not in a loop
	return a + b
}

// goodNoStringConcat uses int concatenation in loop (allowed).
//
// Params:
//   - n: number of iterations
//
// Returns:
//   - int: sum result
func goodNoStringConcat(n int) int {
	result := 0

	// Good: not string concatenation
	for i := 0; i < n; i++ {
		result += i
	}

	// Return the result
	return result
}

// init utilise les fonctions privées
func init() {
	// Appel de goodStringsBuilder
	_ = goodStringsBuilder(nil)
	// Appel de goodStringsBuilderWithGrow
	_ = goodStringsBuilderWithGrow(nil)
	// Appel de goodStringsJoin
	_ = goodStringsJoin(nil)
	// Appel de goodSingleConcat
	_ = goodSingleConcat("", "")
	// Appel de goodNoStringConcat
	_ = goodNoStringConcat(0)
}
