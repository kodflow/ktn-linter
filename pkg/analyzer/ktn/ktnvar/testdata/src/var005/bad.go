// Bad examples for the var006 test case.
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
	// FINAL_VALUE defines final value
	FINAL_VALUE int = 99
)

// badMakeWithLength creates a slice with length then uses append
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithLength() []int {
	// Bad: Length > 0 creates zero values that will be unused
	items := make([]int, SIZE_LIMIT)
	// Itération sur les éléments
	for i := range LOOP_ITERATIONS {
		items = append(items, i)
	}
	// Retour de la fonction
	return items
}

// badMakeWithLengthAndCapacity creates a slice with length then append
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithLengthAndCapacity() []string {
	// Bad: Length > 0 even with capacity
	result := make([]string, LOOP_ITERATIONS, LOOP_ITERATIONS*TWO_VALUE)
	// Itération sur les éléments
	for range FIVE_ITEMS {
		result = append(result, "value")
	}
	// Retour de la fonction
	return result
}

// badMakeWithSmallLength creates a slice with small length then append
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithSmallLength() []int {
	// Bad: Even small length > 0 is wasteful
	numbers := make([]int, FIVE_ITEMS)
	// Itération sur les éléments
	for i := range SMALL_COUNT {
		numbers = append(numbers, i*TWO_VALUE)
	}
	// Retour de la fonction
	return numbers
}

// badMakeWithConstLength creates a slice with const length then append
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithConstLength() []string {
	// Bad: Length from constant > 0
	items := make([]string, LOOP_ITERATIONS)
	// Itération sur les éléments
	for range FIVE_ITEMS {
		items = append(items, "value")
	}
	// Retour de la fonction
	return items
}

// badMakeByteSliceWithLength creates a byte slice with length then append
//
// Returns:
//   - []byte: byte slice with wasted zero values
func badMakeByteSliceWithLength() []byte {
	// Bad: Byte slice with length > 0
	buffer := make([]byte, BAD_BUFFER_SIZE)
	// Ajout de données
	buffer = append(buffer, "test"...)
	// Retour de la fonction
	return buffer
}

// badMakeInterfaceSlice creates a slice of interfaces with length then append
//
// Returns:
//   - []any: slice with wasted zero values
func badMakeInterfaceSlice() []any {
	// Bad: any slice with length > 0
	items := make([]any, FIVE_ITEMS)
	// Itération sur les éléments
	for i := range SMALL_COUNT {
		items = append(items, i)
	}
	// Retour de la fonction
	return items
}

// badMakeWithVariable creates a slice with variable length
//
// Params:
//   - n: size parameter
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithVariable(n int) []int {
	// Bad: Using variable length (assumed positive)
	data := make([]int, n)
	// Ajout de données
	data = append(data, LOOP_ITERATIONS, TWO_VALUE, SMALL_COUNT)
	// Retour de la fonction
	return data
}

// badMakeWithExpression creates a slice with expression length
//
// Params:
//   - items: input items
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithExpression(items []string) []string {
	// Bad: Using len() expression (assumed positive)
	result := make([]string, len(items))
	// Ajout de données
	result = append(result, "extra")
	// Retour de la fonction
	return result
}

// badMakeWithOperation creates a slice with arithmetic expression
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithOperation() []int {
	// Bad: Using arithmetic expression
	nums := make([]int, SMALL_COUNT+TWO_VALUE)
	// Ajout de données
	nums = append(nums, FINAL_VALUE)
	// Retour de la fonction
	return nums
}

// init utilise les fonctions privées
func init() {
	// Appel de badMakeWithLength
	badMakeWithLength()
	// Appel de badMakeWithLengthAndCapacity
	badMakeWithLengthAndCapacity()
	// Appel de badMakeWithSmallLength
	badMakeWithSmallLength()
	// Appel de badMakeWithConstLength
	badMakeWithConstLength()
	// Appel de badMakeByteSliceWithLength
	badMakeByteSliceWithLength()
	// Appel de badMakeInterfaceSlice
	badMakeInterfaceSlice()
	// Appel de badMakeWithVariable
	_ = badMakeWithVariable(0)
	// Appel de badMakeWithExpression
	_ = badMakeWithExpression(nil)
	// Appel de badMakeWithOperation
	badMakeWithOperation()
}
