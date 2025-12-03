package var007

// Bad: Slices created without capacity when it could be known

// Constantes pour les tests
const (
	MAX_SIZE         int = 50
	SMALL_LOOP_COUNT int = 20
)

// badMakeStringSlice creates a string slice without capacity
//
// Returns:
//   - []string: slice without preallocated capacity
func badMakeStringSlice() []string {
	// Bad: make without capacity argument
	items := make([]string, 0)
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

// badMakeFloatSlice creates a float slice without capacity
//
// Returns:
//   - []float64: slice without preallocated capacity
func badMakeFloatSlice() []float64 {
	// Bad: make for float64 without capacity
	values := make([]float64, 0)
	// Retour de la fonction
	return values
}

// badMakeInterfaceSlice creates an interface slice without capacity
//
// Returns:
//   - []interface{}: interface slice without capacity
func badMakeInterfaceSlice() []interface{} {
	// Bad: make for interface{} without capacity
	items := make([]interface{}, 0)
	// Retour de la fonction
	return items
}

// Item is a test struct
type Item struct {
	value int
}

// badMakeStructSlice creates a slice of structs without capacity
//
// Returns:
//   - []Item: slice of structs without capacity
func badMakeStructSlice() []Item {
	// Bad: make for struct slice without capacity
	items := make([]Item, 0)
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

// badEmptyLiteralWithAppend creates empty literal then appends
//
// Returns:
//   - []int: slice that should use make with capacity
func badEmptyLiteralWithAppend() []int {
	// Bad: []T{} with subsequent append should use make
	numbers := []int{}
	// Ajout d'éléments
	numbers = append(numbers, SMALL_LOOP_COUNT)
	numbers = append(numbers, MAX_SIZE)
	// Retour de la fonction
	return numbers
}
