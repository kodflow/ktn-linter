// Good examples for the var009 test case.
package var008

const (
	// VALUE_TWO is constant value 2
	VALUE_TWO int = 2
	// VALUE_THREE is constant value 3
	VALUE_THREE int = 3
	// VALUE_FIVE is constant value 5
	VALUE_FIVE int = 5
	// VALUE_TEN is constant value 10
	VALUE_TEN int = 10
)

// goodLoopSliceReuse réutilise un slice déclaré avant la boucle.
func goodLoopSliceReuse() {
	// Déclaration avant la boucle
	data := make([]int, 0, VALUE_TEN)
	// Loop appends values to reused slice
	for i := 0; i < VALUE_TEN; i++ {
		// Append current iteration value
		data = append(data, i)
	}
	_ = data
}

// goodLoopMapReuse réutilise une map déclarée avant la boucle.
func goodLoopMapReuse() {
	// Déclaration avant la boucle avec capacité
	cache := make(map[string]int, VALUE_TEN)
	// Loop reuses map
	for i := 0; i < VALUE_TEN; i++ {
		// Store current value
		cache["key"] = i
	}
	_ = cache
}

// goodRangeSliceReuse réutilise un slice dans une boucle range.
func goodRangeSliceReuse() {
	items := []int{1, VALUE_TWO, VALUE_THREE}
	// Déclaration avant la boucle avec capacité
	buffer := make([]byte, 0, VALUE_THREE)
	// Range loop reuses buffer
	for _, item := range items {
		// Convert and append item
		buffer = append(buffer, byte(item))
	}
	_ = buffer
}

// goodNoLoopAlloc alloue hors d'une boucle.
func goodNoLoopAlloc() {
	// Pas de boucle, allocation OK avec array
	var data [VALUE_TEN]int
	_ = data
}

// goodNestedLoopReuse réutilise dans une boucle imbriquée.
func goodNestedLoopReuse() {
	// Déclaration avant la boucle avec array
	var temp [VALUE_TEN]int
	// Outer loop iterates
	for i := 0; i < VALUE_FIVE; i++ {
		// Inner loop modifies temp
		for j := 0; j < VALUE_FIVE; j++ {
			// Store multiplication result
			temp[j] = i * j
		}
	}
	_ = temp
}

// init utilise les fonctions privées
func init() {
	// Appel de goodLoopSliceReuse
	goodLoopSliceReuse()
	// Appel de goodLoopMapReuse
	goodLoopMapReuse()
	// Appel de goodRangeSliceReuse
	goodRangeSliceReuse()
	// Appel de goodNoLoopAlloc
	goodNoLoopAlloc()
	// Appel de goodNestedLoopReuse
	goodNestedLoopReuse()
}
