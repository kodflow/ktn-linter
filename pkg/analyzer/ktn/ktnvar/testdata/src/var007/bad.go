package var007

// Bad: Slices created without capacity when it could be known

// Constantes pour les tests
const (
	MAX_SIZE         int = 50
	SMALL_LOOP_COUNT int = 20
)

// badEmptySliceLiteral creates a slice without capacity
//
// Returns:
//   - []int: slice without preallocated capacity
func badEmptySliceLiteral() []int {
	// Bad: Empty literal without capacity
	itemList := []int{}
	// Retour de la fonction
	return itemList
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
	itemList := []string{}
	// Itération sur les éléments
	for index := 0; index < SMALL_LOOP_COUNT; index++ {
		itemList = append(itemList, "value")
	}
	// Retour de la fonction
	return itemList
}

// badSliceOfSlices creates a slice of slices without capacity
//
// Returns:
//   - [][]int: nested slice without capacity
func badSliceOfSlices() [][]int {
	// Bad: Nested slice without capacity
	matrixData := [][]int{}
	// Retour de la fonction
	return matrixData
}

// Item is a test struct
type Item struct {
	value int
}

// badSliceOfStructs creates a slice of structs without capacity
//
// Returns:
//   - []Item: slice of structs without capacity
func badSliceOfStructs() []Item {
	// Bad: Slice of structs without capacity
	itemArray := []Item{}
	// Retour de la fonction
	return itemArray
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
