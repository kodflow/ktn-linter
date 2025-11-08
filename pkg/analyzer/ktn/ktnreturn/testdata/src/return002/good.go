package return002

const (
	// DEFAULT_MAP_CAPACITY capacité par défaut des maps
	DEFAULT_MAP_CAPACITY int = 10
)

// goodReturnEmptySlice returns an empty slice instead of nil.
//
// Returns:
//   - []string: slice vide
func goodReturnEmptySlice() []string {
	// Retourne un slice vide
	return []string{}
}

// goodReturnEmptyMap returns an empty map instead of nil.
//
// Returns:
//   - map[string]int: map vide
func goodReturnEmptyMap() map[string]int {
	// Retourne une map vide
	return map[string]int{}
}

// goodReturnEmptySliceConditional returns empty slice conditionally.
//
// Params:
//   - x: valeur de test
//
// Returns:
//   - []int: slice vide ou avec x
func goodReturnEmptySliceConditional(x int) []int {
	// Vérifie si x est positif
	if x > 0 {
		// Retourne slice vide pour valeur positive
		return []int{}
	}
	// Retourne slice contenant x
	return []int{x}
}

// goodReturnPointerNil can return nil for pointer types (allowed).
//
// Returns:
//   - *string: pointeur nil
func goodReturnPointerNil() *string {
	// Retourne nil pour pointer
	return nil
}

// goodReturnInterfaceNil can return nil for interface types (allowed).
//
// Returns:
//   - error: nil
func goodReturnInterfaceNil() error {
	// Retourne nil pour interface
	return nil
}

// goodReturnError returns nil for error type (allowed).
//
// Returns:
//   - error: nil
func goodReturnError() error {
	// Retourne nil pour error
	return nil
}

// goodReturnMakeSlice returns a slice created with make.
//
// Params:
//   - size: taille du slice
//
// Returns:
//   - []int: slice créé avec make
func goodReturnMakeSlice(size int) []int {
	// Retourne slice de taille spécifiée
	return make([]int, 0, size)
}

// goodReturnMakeMap returns a map created with make.
//
// Returns:
//   - map[string]bool: map créée avec make
func goodReturnMakeMap() map[string]bool {
	// Retourne map avec capacité par défaut
	return make(map[string]bool, DEFAULT_MAP_CAPACITY)
}
