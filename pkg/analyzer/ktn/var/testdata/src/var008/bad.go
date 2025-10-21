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
	items := make([]int, SIZE_LIMIT) // want "KTN-VAR-008"
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
	result := make([]string, 10, 20) // want "KTN-VAR-008"
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
	numbers := make([]int, 5) // want "KTN-VAR-008"
	// Itération sur les éléments
	for i := 0; i < 3; i++ {
		numbers = append(numbers, i*2)
	}
	// Retour de la fonction
	return numbers
}
