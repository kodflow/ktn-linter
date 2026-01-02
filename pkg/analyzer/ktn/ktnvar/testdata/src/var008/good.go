// Package var004 provides good test cases.
package var004

// Good: Slices preallocated with capacity when known

const (
	// MaxItems defines the maximum number of items
	MaxItems int = 100
	// CapacityTen is capacity of ten
	CapacityTen int = 10
	// ValueOne is value one
	ValueOne int = 1
	// ValueTwo is value two
	ValueTwo int = 2
	// ValueThree is value three
	ValueThree int = 3
)

// init demonstrates good slice preallocation patterns
func init() {
	// Good: Capacity specified
	items := make([]int, 0, MaxItems)
	_ = items

	// Good: Capacity specified based on known size
	result := make([]string, 0, CapacityTen)
	// Itération sur les éléments
	for i := range CapacityTen {
		result = append(result, "item")
		// Utilisation de i pour éviter le warning
		_ = i
	}

	// Good: Literal with values is acceptable
	values := []int{ValueOne, ValueTwo, ValueThree}
	_ = values

	// Good: Maps don't need capacity
	m := map[string]int{}
	_ = m

	// Good: Even when size is unknown, providing a reasonable capacity is better
	data := []string{"a", "b"}
	filtered := make([]string, 0, len(data))
	// Itération sur les données
	for _, item := range data {
		// Vérification d'une condition
		if len(item) > 0 {
			filtered = append(filtered, item)
		}
	}
	_ = filtered

	// Good: []T{} is acceptable when capacity unknown
	empty := []int{}
	_ = empty

	// Good: Direct return with []T{} is a common pattern
	items2 := []string{}
	_ = items2

	// Good: [][]T{} acceptable when capacity unknown
	nested := [][]int{}
	_ = nested
}
