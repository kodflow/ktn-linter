package var007

// Bad: Slices created without capacity when it could be known

// MAX_SIZE defines the maximum size
const MAX_SIZE int = 50

// badEmptySliceLiteral creates a slice without capacity
//
// Returns:
//   - []int: slice without preallocated capacity
func badEmptySliceLiteral() []int {
	// Bad: Empty literal without capacity
	items := []int{}
	// Retour de la fonction
	return items
}

// badMakeWithoutCapacity creates a slice using make without capacity
//
// Returns:
//   - []string: slice without preallocated capacity
func badMakeWithoutCapacity() []string {
	// Bad: make without capacity argument
	result := make([]string, 0)
	// Retour de la fonction
	return result
}

// badMakeInLoop creates a slice without capacity in loop
//
// Returns:
//   - []int: slice without preallocated capacity
func badMakeInLoop() []int {
	// Bad: Known size but no capacity
	numbers := make([]int, 0)
	// Itération sur les éléments
	for i := 0; i < MAX_SIZE; i++ {
		numbers = append(numbers, i)
	}
	// Retour de la fonction
	return numbers
}

// badEmptyLiteralInLoop creates a slice without capacity for loop
//
// Returns:
//   - []string: slice without preallocated capacity
func badEmptyLiteralInLoop() []string {
	// Bad: Empty literal for known-size loop
	items := []string{}
	// Itération sur les éléments
	for i := 0; i < 20; i++ {
		items = append(items, "value")
	}
	// Retour de la fonction
	return items
}

// badSliceOfSlices creates a slice of slices without capacity
//
// Returns:
//   - [][]int: nested slice without capacity
func badSliceOfSlices() [][]int {
	// Bad: Nested slice without capacity
	matrix := [][]int{}
	// Retour de la fonction
	return matrix
}

// badSliceOfStructs creates a slice of structs without capacity
type Item struct {
	value int
}

//
// Returns:
//   - []Item: slice of structs without capacity
func badSliceOfStructs() []Item {
	// Bad: Slice of structs without capacity
	items := []Item{}
	// Retour de la fonction
	return items
}

// badMakeByteSlice creates a byte slice without capacity
//
// Returns:
//   - []byte: byte slice without capacity
func badMakeByteSlice() []byte {
	// Bad: Byte slice without capacity
	buffer := make([]byte, 0)
	// Retour de la fonction
	return buffer
}
