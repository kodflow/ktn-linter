package var008

// Bad: Using make with length > 0 when append is used

// SIZE_LIMIT defines the size limit
const SIZE_LIMIT int = 50

// badMakeWithLength creates a slice with length then uses append
//
// Returns:
//   - []int: slice with wasted zero values
func badMakeWithLength() []int {
	// Bad: Length > 0 creates zero values that will be unused
	items := make([]int, SIZE_LIMIT)
	// Itération sur les éléments
	for i := 0; i < 10; i++ {
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
	result := make([]string, 10, 20)
	// Itération sur les éléments
	for i := 0; i < 5; i++ {
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
	numbers := make([]int, 5)
	// Itération sur les éléments
	for i := 0; i < 3; i++ {
		numbers = append(numbers, i*2)
	}
	// Retour de la fonction
	return numbers
}

// TEN_ITEMS defines the number of items
const TEN_ITEMS int = 10

// badMakeWithConstLength creates a slice with const length then append
//
// Returns:
//   - []string: slice with wasted zero values
func badMakeWithConstLength() []string {
	// Bad: Length from constant > 0
	items := make([]string, TEN_ITEMS)
	// Itération sur les éléments
	for i := 0; i < 5; i++ {
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
	buffer := make([]byte, 100)
	// Ajout de données
	buffer = append(buffer, "test"...)
	// Retour de la fonction
	return buffer
}

// badMakeInterfaceSlice creates a slice of interfaces with length then append
//
// Returns:
//   - []interface{}: slice with wasted zero values
func badMakeInterfaceSlice() []interface{} {
	// Bad: Interface slice with length > 0
	items := make([]interface{}, 5)
	// Itération sur les éléments
	for i := 0; i < 3; i++ {
		items = append(items, i)
	}
	// Retour de la fonction
	return items
}

// THREE_ITEMS defines three items
const THREE_ITEMS int = 3

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
	data = append(data, 1, 2, 3)
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
	nums := make([]int, THREE_ITEMS+2)
	// Ajout de données
	nums = append(nums, 99)
	// Retour de la fonction
	return nums
}
