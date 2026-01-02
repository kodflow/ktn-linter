// Package var025 contains good test cases for KTN-VAR-025.
package var025

const (
	// HalfDivisor is used to calculate half length
	HalfDivisor int = 2
	// ValueOne represents 1
	ValueOne int = 1
	// ValueThree represents 3
	ValueThree int = 3
)

// goodClearMap utilise clear() pour vider une map.
//
// Params:
//   - mapData: map a vider
func goodClearMap(mapData map[string]int) {
	// Pattern correct: clear built-in
	clear(mapData)
}

// goodZeroSlice utilise clear() pour remettre a zero une slice.
//
// Params:
//   - sliceData: slice a remettre a zero
func goodZeroSlice(sliceData []int) {
	// Pattern correct: clear built-in
	clear(sliceData)
}

// goodPartialClear efface partiellement une slice (non detectable).
//
// Params:
//   - sliceData: slice a effacer partiellement
func goodPartialClear(sliceData []int) {
	// OK: effacement partiel, pas un pattern clear complet
	halfLen := len(sliceData) / HalfDivisor
	// Boucle pour effacement partiel
	for idx := range halfLen {
		// Effacement partiel seulement
		sliceData[idx] = 0
	}
}

// goodMultipleStatements a plusieurs instructions dans la boucle.
//
// Params:
//   - mapData: map a traiter
func goodMultipleStatements(mapData map[string]int) {
	// OK: plusieurs instructions, pas un pattern clear simple
	for keyItem := range mapData {
		// Log avant suppression
		_ = keyItem
		// Suppression
		delete(mapData, keyItem)
	}
}

// goodDifferentKey utilise une cle differente pour delete.
//
// Params:
//   - mapData: map a traiter
//   - keyToDelete: cle specifique a supprimer
func goodDifferentKey(mapData map[string]int, keyToDelete string) {
	// OK: cle differente de celle du range
	for keyItem := range mapData {
		// Utilisation de la cle de range
		_ = keyItem
		// Suppression avec cle differente
		delete(mapData, keyToDelete)
	}
}

// goodNonZeroValue assigne une valeur non-zero.
//
// Params:
//   - sliceData: slice a modifier
func goodNonZeroValue(sliceData []int) {
	// OK: valeur non-zero
	for idx := range sliceData {
		// Affectation de valeur non-zero
		sliceData[idx] = ValueOne
	}
}

// init utilise les fonctions
func init() {
	// Test goodClearMap
	testMap := map[string]int{"one": ValueOne}
	goodClearMap(testMap)

	// Test goodZeroSlice
	testSlice := []int{ValueOne, HalfDivisor, ValueThree}
	goodZeroSlice(testSlice)

	// Test goodPartialClear
	goodPartialClear(testSlice)

	// Test goodMultipleStatements
	goodMultipleStatements(testMap)

	// Test goodDifferentKey
	goodDifferentKey(testMap, "key")

	// Test goodNonZeroValue
	goodNonZeroValue(testSlice)
}
