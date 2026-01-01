// Package var005 tests make with length then append detection.
package var005

// Bad: Using make with length > 0 when append is used

const (
	// SizeLimit defines the size limit
	SizeLimit int = 50
	// LoopIterations defines loop count
	LoopIterations int = 10
	// FiveItems defines five items
	FiveItems int = 5
	// SmallCount defines small count
	SmallCount int = 3
	// TwoValue defines two value
	TwoValue int = 2
	// BadBufferSize defines bad buffer size
	BadBufferSize int = 100
	// Capacity defines capacity
	Capacity int = 200
)

// badMakeWithLengthAndCapacity creates a slice with length then append.
// Bad: has capacity but also has length > 0, which wastes zero values.
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithLengthAndCapacity() []string {
	// Bad: Length > 0 even with capacity creates zero values before append
	result := make([]string, LoopIterations, Capacity)
	// Itération sur les éléments
	for range FiveItems {
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
	numbers := make([]int, FiveItems, Capacity)
	// Itération sur les éléments
	for i := range SmallCount {
		numbers = append(numbers, i*TwoValue)
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
	items := make([]string, LoopIterations, Capacity)
	// Itération sur les éléments
	for range FiveItems {
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
	buffer := make([]byte, BadBufferSize, Capacity)
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
	items := make([]any, FiveItems, Capacity)
	// Itération sur les éléments
	for i := range SmallCount {
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
	data := make([]int, n, Capacity)
	// Ajout de données
	data = append(data, LoopIterations, TwoValue, SmallCount)
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
	result := make([]string, len(items), Capacity)
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
	nums := make([]int, SmallCount+TwoValue, Capacity)
	// Ajout de données
	nums = append(nums, SizeLimit)
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
