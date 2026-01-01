// Package var025 contains bad test cases for KTN-VAR-025.
package var025

// badClearMap utilise une boucle au lieu de clear() pour vider une map.
//
// Params:
//   - mapData: map a vider
func badClearMap(mapData map[string]int) {
	// Pattern inefficace: boucle delete
	for keyItem := range mapData {
		delete(mapData, keyItem) // want "KTN-VAR-025"
	}
}

// badZeroSlice utilise une boucle au lieu de clear() pour remettre a zero.
//
// Params:
//   - sliceData: slice a remettre a zero
func badZeroSlice(sliceData []int) {
	// Pattern inefficace: boucle zero
	for idx := range sliceData {
		sliceData[idx] = 0 // want "KTN-VAR-025"
	}
}

// init utilise les fonctions privees
func init() {
	// Appel de badClearMap
	testMap := map[string]int{"one": 1}
	badClearMap(testMap)

	// Appel de badZeroSlice
	testSlice := []int{1}
	badZeroSlice(testSlice)
}
