// Package var005 tests make with length then append detection.
package var005

// Bad: Using make with length > 0 when append is used

const (
	// SIZE_LIMIT defines the size limit
	SIZE_LIMIT int = 50
	// LOOP_ITERATIONS defines loop count
	LOOP_ITERATIONS int = 10
	// FIVE_ITEMS defines five items
	FIVE_ITEMS int = 5
	// SMALL_COUNT defines small count
	SMALL_COUNT int = 3
	// TWO_VALUE defines two value
	TWO_VALUE int = 2
	// BAD_BUFFER_SIZE defines bad buffer size
	BAD_BUFFER_SIZE int = 100
	// CAPACITY defines capacity
	CAPACITY int = 200
)

// badMakeWithLengthAndCapacity creates a slice with length then append.
// Bad: has capacity but also has length > 0, which wastes zero values.
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithLengthAndCapacity() []string {
	// Bad: Length > 0 even with capacity creates zero values before append
	result := make([]string, LOOP_ITERATIONS, CAPACITY)
	// Itération sur les éléments
	for range FIVE_ITEMS {
		result = append(result, "value")
	}
	// Retour de la fonction
	return result
}

// badMakeWithSmallLengthAndCapacity creates a slice with small length.
// Bad: even small length > 0 with capacity is wasteful.
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithSmallLengthAndCapacity() []int {
	// Bad: Even small length > 0 wastes zero values
	numbers := make([]int, FIVE_ITEMS, CAPACITY)
	// Itération sur les éléments
	for i := range SMALL_COUNT {
		numbers = append(numbers, i*TWO_VALUE)
	}
	// Retour de la fonction
	return numbers
}

// badMakeWithConstLengthAndCapacity creates a slice with const length.
// Bad: using const for length creates unused zero values.
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithConstLengthAndCapacity() []string {
	// Bad: Length from constant > 0 with capacity
	items := make([]string, LOOP_ITERATIONS, CAPACITY)
	// Itération sur les éléments
	for range FIVE_ITEMS {
		items = append(items, "value")
	}
	// Retour de la fonction
	return items
}

// badMakeByteSliceWithLengthAndCapacity creates a byte slice with length.
// Bad: byte slice with length creates zero bytes before append.
//
// Returns:
//   - []byte: byte slice with wasted zero values
func badMakeByteSliceWithLengthAndCapacity() []byte {
	// Bad: Byte slice with length > 0 and capacity
	buffer := make([]byte, BAD_BUFFER_SIZE, CAPACITY)
	// Ajout de données
	buffer = append(buffer, "test"...)
	// Retour de la fonction
	return buffer
}

// badMakeInterfaceSliceWithCapacity creates a slice of interfaces.
// Bad: any slice with length creates nil values before append.
//
// Returns:
//   - []any: slice with wasted zero values
func badMakeInterfaceSliceWithCapacity() []any {
	// Bad: any slice with length > 0 and capacity
	items := make([]any, FIVE_ITEMS, CAPACITY)
	// Itération sur les éléments
	for i := range SMALL_COUNT {
		items = append(items, i)
	}
	// Retour de la fonction
	return items
}

// badMakeWithVariableAndCapacity creates a slice with variable length.
// Bad: using variable length with capacity still wastes space.
//
// Params:
//   - n: size parameter
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithVariableAndCapacity(n int) []int {
	// Bad: Using variable length with capacity
	data := make([]int, n, CAPACITY)
	// Ajout de données
	data = append(data, LOOP_ITERATIONS, TWO_VALUE, SMALL_COUNT)
	// Retour de la fonction
	return data
}

// badMakeWithExpressionAndCapacity creates a slice with expression length.
// Bad: using len() with capacity still creates unused zero values.
//
// Params:
//   - items: input items
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithExpressionAndCapacity(items []string) []string {
	// Bad: Using len() expression with capacity
	result := make([]string, len(items), CAPACITY)
	// Ajout de données
	result = append(result, "extra")
	// Retour de la fonction
	return result
}

// badMakeWithOperationAndCapacity creates a slice with arithmetic expression.
// Bad: using arithmetic with capacity creates unused zero values.
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithOperationAndCapacity() []int {
	// Bad: Using arithmetic expression with capacity
	nums := make([]int, SMALL_COUNT+TWO_VALUE, CAPACITY)
	// Ajout de données
	nums = append(nums, SIZE_LIMIT)
	// Retour de la fonction
	return nums
}

// init utilise les fonctions privées
func init() {
	// Appel de badMakeWithLengthAndCapacity
	badMakeWithLengthAndCapacity()
	// Appel de badMakeWithSmallLengthAndCapacity
	badMakeWithSmallLengthAndCapacity()
	// Appel de badMakeWithConstLengthAndCapacity
	badMakeWithConstLengthAndCapacity()
	// Appel de badMakeByteSliceWithLengthAndCapacity
	badMakeByteSliceWithLengthAndCapacity()
	// Appel de badMakeInterfaceSliceWithCapacity
	badMakeInterfaceSliceWithCapacity()
	// Appel de badMakeWithVariableAndCapacity
	_ = badMakeWithVariableAndCapacity(0)
	// Appel de badMakeWithExpressionAndCapacity
	_ = badMakeWithExpressionAndCapacity(nil)
	// Appel de badMakeWithOperationAndCapacity
	badMakeWithOperationAndCapacity()
}
