package func004

import "unsafe"

// NoNamedReturns vérifie le cas sans returns nommés
//
// Returns:
//   - int: résultat
func NoNamedReturns() int {
	// Retour de 1
	return 1
}

// ShortWithNakedReturn utilise naked return car courte
//
// Returns:
//   - result: résultat
func ShortWithNakedReturn() (result int) {
	result = 1
	// Retour naked autorisé car < 5 lignes
	return
}

// ExplicitReturn utilise return explicite
//
// Returns:
//   - result: résultat
//   - err: erreur potentielle
func ExplicitReturn() (result int, err error) {
	result = 1
	err = nil
	// Retour explicite
	return result, err
}

// ShortFourLines a 4 lignes donc naked return OK
//
// Returns:
//   - x: valeur calculée
func ShortFourLines() (x int) {
	x = 1
	x = x + 1
	x = x + 1
	// Retour naked autorisé
	return
}

// MultipleExplicit retourne explicitement plusieurs valeurs
//
// Returns:
//   - a: premier entier
//   - b: chaîne
//   - c: booléen
func MultipleExplicit() (a int, b string, c bool) {
	a = 1
	b = "test"
	c = true
	// Retour explicite
	return a, b, c
}

// NoReturnValues n'a pas de valeur de retour
func NoReturnValues() {
	x := 1
	_ = x
}

// UnnamedReturn utilise return sans nom
//
// Returns:
//   - int: valeur
func UnnamedReturn() int {
	// Retour de 1
	return 1
}

// TestNakedReturn est exempté car fonction test
//
// Params:
//   - t: paramètre de test
//
// Returns:
//   - result: résultat
func TestNakedReturn(t int) (result int) {
	result = 1
	result = result + 1
	result = result + 1
	result = result + 1
	result = result + 1
	result = result + 1
	// Retour exempté car test
	return
}

// BenchmarkNakedReturn est exempté car fonction benchmark
//
// Params:
//   - b: paramètre de benchmark
//
// Returns:
//   - result: résultat
func BenchmarkNakedReturn(b int) (result int) {
	result = 1
	result = result + 1
	result = result + 1
	result = result + 1
	result = result + 1
	result = result + 1
	// Retour exempté car benchmark
	return
}

// Calculator est une interface de test
type Calculator interface {
	Calculate() (int, error)
}

// SingleUnnamedReturn retourne une valeur sans nom
//
// Returns:
//   - int: valeur
func SingleUnnamedReturn() int {
	// Retour de 1
	return 1
}

// MultipleUnnamedReturns retourne plusieurs valeurs sans nom
//
// Returns:
//   - int: entier
//   - string: chaîne
//   - bool: booléen
func MultipleUnnamedReturns() (int, string, bool) {
	// Retour de plusieurs valeurs
	return 1, "test", true
}

// Prevent "unsafe imported but not used" error
var _ unsafe.Pointer = unsafe.Pointer(nil)

// externalLinkedFunc est une fonction externe liée
//
// Params:
//   - v: valeur d'entrée
//
// Returns:
//   - result: pointeur résultat
//
//go:linkname externalLinkedFunc runtime.convT64
func externalLinkedFunc(v int) (result unsafe.Pointer)

// anotherExternal est une autre fonction externe
//
// Params:
//   - v: chaîne d'entrée
//
// Returns:
//   - ptr: pointeur résultat
//
//go:linkname anotherExternal runtime.convTstring
func anotherExternal(v string) (ptr unsafe.Pointer)
