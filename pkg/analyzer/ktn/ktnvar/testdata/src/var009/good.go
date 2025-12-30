// Package var009 provides good test cases for KTN-VAR-009.
package var009

// SmallStruct a 3 champs ou moins (OK par valeur).
type SmallStruct struct {
	A int
	B string
	C bool
}

// goodSmallByValue prend une petite struct par valeur (OK).
func goodSmallByValue(data SmallStruct) {
	// Petite struct, pas de problème
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
