// Good examples for the var007 test case.
package var007

// Good: Slices preallocated with capacity when known

const (
	// MAX_ITEMS defines the maximum number of items
	MAX_ITEMS int = 100
	// CAPACITY_TEN is capacity of ten
	CAPACITY_TEN int = 10
	// VALUE_ONE is value one
	VALUE_ONE int = 1
	// VALUE_TWO is value two
	VALUE_TWO int = 2
	// VALUE_THREE is value three
	VALUE_THREE int = 3
)

// goodPreallocatedSlice creates a slice with capacity
//
// Returns:
//   - []int: slice with preallocated capacity
func goodPreallocatedSlice() []int {
	// Good: Capacity specified
	items := make([]int, 0, MAX_ITEMS)
	// Continue traversing AST nodes.
	return items
}

// goodPreallocatedLoop creates a slice with capacity in loop
//
// Returns:
//   - []string: slice with preallocated capacity
func goodPreallocatedLoop() []string {
	// Good: Capacity specified based on known size
	result := make([]string, 0, CAPACITY_TEN)
	// Itération sur les éléments
	for i := range CAPACITY_TEN {
		result = append(result, "item")
		// Utilisation de i pour éviter le warning
		_ = i
	}
	// Retour de la fonction
	return result
}

// goodLiteralWithValues creates a slice with initial values
//
// Returns:
//   - []int: slice with initial values
func goodLiteralWithValues() []int {
	// Good: Literal with values is acceptable
	return []int{VALUE_ONE, VALUE_TWO, VALUE_THREE}
}

// goodNonSliceComposite creates a non-slice composite
//
// Returns:
//   - map[string]int: empty map
func goodNonSliceComposite() map[string]int {
	// Good: Maps don't need capacity
	return map[string]int{}
}

// goodDynamicSize creates a slice when size is unknown
//
// Params:
//   - data: input data
//
// Returns:
//   - []string: processed data
func goodDynamicSize(data []string) []string {
	// Good: Even when size is unknown, providing a reasonable capacity is better
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

// goodEmptySliceLiteral returns empty slice literal (acceptable pattern)
//
// Returns:
//   - []int: empty slice
func goodEmptySliceLiteral() []int {
	// Good: []T{} is acceptable when capacity unknown
	// Retour de la fonction
	return []int{}
}

// goodEmptyLiteralForReturn uses empty literal in direct return
//
// Returns:
//   - []string: empty slice
func goodEmptyLiteralForReturn() []string {
	// Good: Direct return with []T{} is a common pattern
	items := []string{}
	// Retour de la fonction
	return items
}

// goodSliceOfSlices creates empty slice of slices
//
// Returns:
//   - [][]int: empty nested slice
func goodSliceOfSlices() [][]int {
	// Good: [][]T{} acceptable when capacity unknown
	// Retour de la fonction
	return [][]int{}
}

// GoodItem is a test struct.
// Utilisé pour démontrer les slices de structs.
type GoodItem struct {
	value int
}

// GoodContainer is a test struct with slice fields.
// Contient des slices pour démontrer l'initialisation.
type GoodContainer struct {
	items  []int
	values []string
}

// goodSliceOfStructs creates empty slice of structs
//
// Returns:
//   - []GoodItem: empty slice of structs
func goodSliceOfStructs() []GoodItem {
	// Good: []T{} for structs is acceptable
	// Retour de la fonction
	return []GoodItem{}
}

// goodEmptySliceInStructInit initializes struct with empty slices
//
// Returns:
//   - GoodContainer: struct with empty slices
func goodEmptySliceInStructInit() GoodContainer {
	// Good: []T{} in struct initialization is acceptable
	c := GoodContainer{
		items:  []int{},
		values: []string{},
	}
	// Retour de la fonction
	return c
}

// goodEmptySliceNotAppended creates slice but never appends
//
// Returns:
//   - []int: empty slice (no append = no need for capacity)
func goodEmptySliceNotAppended() []int {
	// Good: []T{} without append is acceptable
	numbers := []int{}
	// Retour de la fonction
	return numbers
}

// init utilise les fonctions privées
func init() {
	// Appel de goodPreallocatedSlice
	goodPreallocatedSlice()
	// Appel de goodPreallocatedLoop
	goodPreallocatedLoop()
	// Appel de goodLiteralWithValues
	goodLiteralWithValues()
	// Appel de goodNonSliceComposite
	goodNonSliceComposite()
	// Appel de goodDynamicSize
	_ = goodDynamicSize(nil)
	// Appel de goodEmptySliceLiteral
	goodEmptySliceLiteral()
	// Appel de goodEmptyLiteralForReturn
	goodEmptyLiteralForReturn()
	// Appel de goodSliceOfSlices
	goodSliceOfSlices()
	// Appel de goodSliceOfStructs
	goodSliceOfStructs()
	// Appel de goodEmptySliceInStructInit
	goodEmptySliceInStructInit()
	// Appel de goodEmptySliceNotAppended
	goodEmptySliceNotAppended()
}
