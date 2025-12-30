// Package var009 contains test cases for KTN-VAR-009.
package var009

// LargeStruct est une structure avec plus de 3 champs.
type LargeStruct struct {
	Field1 int
	Field2 string
	Field3 bool
	Field4 float64
}

// badProcessByValue prend une grande struct par valeur.
func badProcessByValue(data LargeStruct) { // want "KTN-VAR-009"
	// Utilise la structure
	_ = data.Field1
}

// badMultipleParams prend plusieurs grandes structs par valeur.
func badMultipleParams(a LargeStruct, b LargeStruct) { // want "KTN-VAR-009" "KTN-VAR-009"
	// Utilise les structures
	_, _ = a.Field1, b.Field1
}

// badMixedParams mélange pointeur et valeur.
func badMixedParams(good *LargeStruct, bad LargeStruct) { // want "KTN-VAR-009"
	// Le pointeur est OK, la valeur non
	_ = good.Field1 + bad.Field1
}
