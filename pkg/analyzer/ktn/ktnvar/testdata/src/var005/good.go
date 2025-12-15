// Package var005 provides good test cases.
package var005

// Good: Proper use of make with capacity or without length

const (
	// BufferSize defines the buffer size
	BufferSize int = 100
	// LoopCount is the loop count
	LoopCount int = 10
	// ValueOne is value one
	ValueOne int = 1
	// ValueTwo is value two
	ValueTwo int = 2
	// ValueThree is value three
	ValueThree int = 3
	// CapacityFifty is capacity of 50
	CapacityFifty int = 50
)

// init demonstrates proper make usage
func init() {
	// Good: Capacity specified, length is 0
	items := make([]int, 0, BufferSize)
	// Itération sur les éléments
	for i := range LoopCount {
		items = append(items, i)
	}
	_ = items

	// Good: Using make([]T, 0, cap) with append is the proper way
	items2 := make([]string, 0, BufferSize)
	// Itération sur les éléments
	for i := range BufferSize {
		items2 = append(items2, "value")
		// Utilisation de i pour éviter le warning
		_ = i
	}
	_ = items2

	// Good: Literal with values
	values := []int{ValueOne, ValueTwo, ValueThree}
	_ = values

	// Good: Use make with capacity based on input length
	data := []string{"a", "b"}
	result := make([]string, 0, len(data))
	// Itération sur les données
	for _, item := range data {
		// Vérification d'une condition
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	_ = result

	// Good: Length is 0, capacity is specified
	slice := make([]int, 0, CapacityFifty)
	_ = slice
}
