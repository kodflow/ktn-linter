// good.go contient des exemples de code conforme pour KTN-GENERIC-005.
// Ces fonctions utilisent des noms conventionnels pour les type parameters.
package generic005

// goodProcessT utilise "T" comme type parameter.
// Correct: "T" est un nom conventionnel pour type parameter.
//
// Params:
//   - s: valeur a traiter
//
// Returns:
//   - T: valeur traitee
func goodProcessT[T any](s T) T {
	// Retour de la valeur
	return s
}

// goodHandleE utilise "E" comme type parameter.
// Correct: "E" est un nom conventionnel pour type parameter.
//
// Params:
//   - e: valeur a traiter
//
// Returns:
//   - E: valeur traitee
func goodHandleE[E any](e E) E {
	// Retour de la valeur
	return e
}

// goodProcessN utilise "N" comme type parameter.
// Correct: "N" est un nom conventionnel pour type parameter.
//
// Params:
//   - n: valeur a traiter
//
// Returns:
//   - N: valeur traitee
func goodProcessN[N any](n N) N {
	// Retour de la valeur
	return n
}

// goodCheckB utilise "B" comme type parameter.
// Correct: "B" est un nom conventionnel pour type parameter.
//
// Params:
//   - b: valeur a verifier
//
// Returns:
//   - B: resultat
func goodCheckB[B any](b B) B {
	// Retour de la valeur
	return b
}

// goodMapKV utilise "K" et "V" comme type parameters.
// Correct: "K" et "V" sont des noms conventionnels pour key/value.
//
// Params:
//   - k: cle
//   - v: valeur
//
// Returns:
//   - K: cle retournee
//   - V: valeur retournee
func goodMapKV[K comparable, V any](k K, v V) (K, V) {
	// Retour des valeurs
	return k, v
}

// goodTransform utilise "Input" et "Output" comme type parameters.
// Correct: noms descriptifs qui ne sont pas des identifiants predeclares.
//
// Params:
//   - i: valeur d'entree
//   - f: fonction de transformation
//
// Returns:
//   - Output: valeur transformee
func goodTransform[Input, Output any](i Input, f func(Input) Output) Output {
	// Appliquer la transformation
	return f(i)
}

// goodSliceOp utilise "Element" comme type parameter.
// Correct: nom descriptif qui n'est pas un identifiant predeclare.
//
// Params:
//   - s: slice a traiter
//
// Returns:
//   - []Element: slice traite
func goodSliceOp[Element any](s []Element) []Element {
	// Retour du slice
	return s
}

// GoodTypeT est un type generique avec "T" comme type parameter.
// Correct: "T" est un nom conventionnel pour type parameter.
type GoodTypeT[T any] struct {
	// value est la valeur stockee
	value T
}

// GoodTypeKV est un type generique avec "K" et "V" comme type parameters.
// Correct: "K" et "V" sont des noms conventionnels pour key/value.
type GoodTypeKV[K comparable, V any] struct {
	// key est la cle
	key K
	// value est la valeur
	value V
}
