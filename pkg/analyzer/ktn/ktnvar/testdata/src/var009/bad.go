// Package var009 contains test cases for KTN-VAR-009.
package var009

// LargeStruct est une structure avec >64 bytes (9 * 8 = 72 bytes).
type LargeStruct struct {
	Field1 int64
	Field2 int64
	Field3 int64
	Field4 int64
	Field5 int64
	Field6 int64
	Field7 int64
	Field8 int64
	Field9 int64
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

// badMethodReceiver utilise receiver par valeur sur grande struct.
func (s LargeStruct) badMethodReceiver() int64 { // want "KTN-VAR-009"
	// Le receiver est copié
	return s.Field1
}
