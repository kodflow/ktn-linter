package var007

// Good: Slices preallocated with capacity when known

const (
	// MAX_ITEMS defines the maximum number of items
	MAX_ITEMS int = 100
	// CAPACITY_TEN is capacity of ten
	CAPACITY_TEN int = 10
	// VALUE_ONE is value one
	VALUE_ONE int = 1
	// VALUE_TWO is value two
	VALUE_TWO int = 2
	// VALUE_THREE is value three
	VALUE_THREE int = 3
)

// goodPreallocatedSlice creates a slice with capacity
//
// Returns:
//   - []int: slice with preallocated capacity
func goodPreallocatedSlice() []int {
	// Good: Capacity specified
	items := make([]int, 0, MAX_ITEMS)
	// Continue traversing AST nodes.
	return items
}

// goodPreallocatedLoop creates a slice with capacity in loop
//
// Returns:
//   - []string: slice with preallocated capacity
func goodPreallocatedLoop() []string {
	// Good: Capacity specified based on known size
	result := make([]string, 0, CAPACITY_TEN)
	// Itération sur les éléments
	for i := 0; i < CAPACITY_TEN; i++ {
		result = append(result, "item")
	}
	// Retour de la fonction
	return result
}

// goodLiteralWithValues creates a slice with initial values
//
// Returns:
//   - []int: slice with initial values
func goodLiteralWithValues() []int {
	// Good: Literal with values is acceptable
	return []int{VALUE_ONE, VALUE_TWO, VALUE_THREE}
}

// goodNonSliceComposite creates a non-slice composite
//
// Returns:
//   - map[string]int: empty map
func goodNonSliceComposite() map[string]int {
	// Good: Maps don't need capacity
	return map[string]int{}
}

// goodDynamicSize creates a slice when size is unknown
//
// Params:
//   - data: input data
//
// Returns:
//   - []string: processed data
func goodDynamicSize(data []string) []string {
	// Good: Even when size is unknown, providing a reasonable capacity is better
	result := make([]string, 0, len(data))
	// Itération sur les données
	for _, item := range data {
		// Vérification d'une condition
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	// Retour de la fonction
	return result
}
