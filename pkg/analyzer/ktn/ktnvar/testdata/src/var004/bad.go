// Bad examples for the var005 test case.
package var004

// Bad: Slices created without capacity when it could be known

// Constantes pour les tests
const (
	MaxSize        int = 50
	SmallLoopCount int = 20
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
	for i := range MaxSize {
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
//   - []any: interface slice without capacity
func badMakeInterfaceSlice() []any {
	// Bad: make for any without capacity
	items := make([]any, 0)
	// Retour de la fonction
	return items
}

// Item represents a test item with a single value field.
// This struct is used for testing slice allocation patterns.
type Item struct {
	value int
}

// ItemInterface définit les méthodes de Item.
type ItemInterface interface {
	Value() int
}

// NewItem crée une nouvelle instance de Item.
//
// Returns:
//   - *Item: nouvelle instance
func NewItem() *Item {
	// Retourne une nouvelle instance
	return &Item{}
}

// Value retourne la valeur.
//
// Returns:
//   - int: valeur du champ
func (i *Item) Value() int {
	// Retourne le champ value
	return i.value
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
	// Bad: make([]T, 0) without capacity when size could be known
	numbers := make([]int, 0)
	// Ajout d'éléments
	numbers = append(numbers, SmallLoopCount)
	numbers = append(numbers, MaxSize)
	// Retour de la fonction
	return numbers
}

// init utilise les fonctions privées
func init() {
	// Appel de badMakeStringSlice
	badMakeStringSlice()
	// Appel de badMakeWithoutCapacity
	badMakeWithoutCapacity()
	// Appel de badMakeInLoop
	badMakeInLoop()
	// Appel de badMakeFloatSlice
	badMakeFloatSlice()
	// Appel de badMakeInterfaceSlice
	badMakeInterfaceSlice()
	// Appel de badMakeStructSlice
	badMakeStructSlice()
	// Appel de badMakeByteSlice
	badMakeByteSlice()
	// Appel de badEmptyLiteralWithAppend
	badEmptyLiteralWithAppend()
}
