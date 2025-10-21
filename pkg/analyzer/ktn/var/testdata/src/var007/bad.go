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
	items := []int{} // want "KTN-VAR-007"
	// Retour de la fonction
	return items
}

// badMakeWithoutCapacity creates a slice using make without capacity
//
// Returns:
//   - []string: slice without preallocated capacity
func badMakeWithoutCapacity() []string {
	// Bad: make without capacity argument
	result := make([]string, 0) // want "KTN-VAR-007"
	// Retour de la fonction
	return result
}

// badMakeInLoop creates a slice without capacity in loop
//
// Returns:
//   - []int: slice without preallocated capacity
func badMakeInLoop() []int {
	// Bad: Known size but no capacity
	numbers := make([]int, 0) // want "KTN-VAR-007"
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
	items := []string{} // want "KTN-VAR-007"
	// Itération sur les éléments
	for i := 0; i < 20; i++ {
		items = append(items, "value")
	}
	// Retour de la fonction
	return items
}
