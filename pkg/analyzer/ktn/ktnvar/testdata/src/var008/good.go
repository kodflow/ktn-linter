package var008

// Good: Proper use of make with capacity or without length

const (
	// BUFFER_SIZE defines the buffer size
	BUFFER_SIZE int = 100
	// LOOP_COUNT is the loop count
	LOOP_COUNT int = 10
	// VALUE_ONE is value one
	VALUE_ONE int = 1
	// VALUE_TWO is value two
	VALUE_TWO int = 2
	// VALUE_THREE is value three
	VALUE_THREE int = 3
	// CAPACITY_FIFTY is capacity of 50
	CAPACITY_FIFTY int = 50
)

// goodMakeWithCapacity creates a slice with capacity for append
//
// Returns:
//   - []int: slice with preallocated capacity
func goodMakeWithCapacity() []int {
	// Good: Capacity specified, length is 0
	items := make([]int, 0, BUFFER_SIZE)
	// Itération sur les éléments
	for i := 0; i < LOOP_COUNT; i++ {
		items = append(items, i)
	}
	// Retour de la fonction
	return items
}

// goodMakeWithIndexing creates a slice using indexing instead of append
//
// Returns:
//   - []string: slice with initial zero values
func goodMakeWithIndexing() []string {
	// Good: Using make([]T, 0, cap) with append is the proper way
	items := make([]string, 0, BUFFER_SIZE)
	// Itération sur les éléments
	for i := 0; i < BUFFER_SIZE; i++ {
		items = append(items, "value")
	}
	// Retour de la fonction
	return items
}

// goodLiteralWithValues creates a slice with initial values
//
// Returns:
//   - []int: slice with initial values
func goodLiteralWithValues() []int {
	// Good: Literal with values
	return []int{VALUE_ONE, VALUE_TWO, VALUE_THREE}
}

// goodEmptySliceForDynamic creates a slice for dynamic append
//
// Params:
//   - data: input data
//
// Returns:
//   - []string: filtered data
func goodEmptySliceForDynamic(data []string) []string {
	// Good: Use make with capacity based on input length
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

// goodMakeZeroLength creates a slice with explicit zero length
//
// Returns:
//   - []int: slice with capacity for append
func goodMakeZeroLength() []int {
	// Good: Length is 0, capacity is specified
	return make([]int, 0, CAPACITY_FIFTY)
}
