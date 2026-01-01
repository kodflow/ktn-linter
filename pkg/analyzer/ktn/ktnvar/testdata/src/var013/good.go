// Package var013 provides good test cases for KTN-VAR-013.
package var013

// SmallStruct a exactement 64 bytes (8 * 8 = 64 bytes, seuil OK).
type SmallStruct struct {
	A int64
	B int64
	C int64
	D int64
	E int64
	F int64
	G int64
	H int64
}

// goodSmallByValue prend une struct ≤64 bytes par valeur (OK).
func goodSmallByValue(data SmallStruct) {
	// Struct ≤64 bytes, pas de problème
	_ = data.A
}

// goodLargeByPointer prend une grande struct par pointeur (correct).
func goodLargeByPointer(data *LargeStruct) {
	// Pointeur, pas de copie
	_ = data.Field1
}

// goodMultiplePointers prend plusieurs pointeurs.
func goodMultiplePointers(a *LargeStruct, b *LargeStruct) {
	// Tous les pointeurs
	_, _ = a.Field1, b.Field1
}

// goodMixedSmall mélange petite struct et pointeur.
func goodMixedSmall(small SmallStruct, large *LargeStruct) {
	// Petite par valeur OK, grande par pointeur OK
	_ = small.A + large.Field1
}

// goodPrimitives utilise des types primitifs.
func goodPrimitives(a int, b string, c bool) {
	// Primitifs par valeur, pas de problème
	_, _, _ = a, b, c
}

// goodMethodReceiver utilise receiver par pointeur.
func (s *LargeStruct) goodMethodReceiver() int64 {
	// Pointeur, pas de copie
	return s.Field1
}
