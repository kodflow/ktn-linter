// Bad examples for the return002 test case.
package return002

// badReturnNilSlice returns nil for a slice type.
//
// Returns:
//   - []string: slice of strings
func badReturnNilSlice() []string {
	// Retourne nil au lieu d'un slice vide
	return nil // want "KTN-RETURN-002"
}

// badReturnNilMap returns nil for a map type.
//
// Returns:
//   - map[string]int: map of strings to integers
func badReturnNilMap() map[string]int {
	// Retourne nil au lieu d'une map vide
	return nil // want "KTN-RETURN-002"
}

// badReturnNilSliceConditional returns nil conditionally.
//
// Params:
//   - x: input integer
//
// Returns:
//   - []int: slice of integers
func badReturnNilSliceConditional(x int) []int {
	// Vérifie si x est positif
	if x > 0 {
		// Retourne nil au lieu d'un slice vide
		return nil // want "KTN-RETURN-002"
	}
	// Retourne un slice avec x
	return []int{x}
}

// badReturnNilMapConditional returns nil conditionally.
//
// Params:
//   - key: input key string
//
// Returns:
//   - map[string]string: map of strings
func badReturnNilMapConditional(key string) map[string]string {
	// Vérifie si la clé est vide
	if key == "" {
		// Retourne nil au lieu d'une map vide
		return nil // want "KTN-RETURN-002"
	}
	// Retourne une map avec la clé
	return map[string]string{key: "value"}
}

// badMultipleReturnsWithNil has multiple return statements with nil.
//
// Params:
//   - _flag: boolean flag not used
//
// Returns:
//   - []byte: slice of bytes
func badMultipleReturnsWithNil(_flag bool) []byte {
	// Vérifie le flag
	if _flag {
		// Retourne nil au lieu d'un slice vide
		return nil // want "KTN-RETURN-002"
	}
	// Retourne nil au lieu d'un slice vide
	return nil // want "KTN-RETURN-002"
}

// init utilise les fonctions privées
func init() {
	// Appel de badReturnNilSlice
	badReturnNilSlice()
	// Appel de badReturnNilMap
	badReturnNilMap()
	// Appel de badReturnNilSliceConditional
	_ = badReturnNilSliceConditional(0)
	// Appel de badReturnNilMapConditional
	_ = badReturnNilMapConditional("")
	// Appel de badMultipleReturnsWithNil
	_ = badMultipleReturnsWithNil(false)
}
