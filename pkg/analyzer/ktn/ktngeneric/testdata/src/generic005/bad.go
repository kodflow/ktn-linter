// bad.go contient des exemples de code non-conforme pour KTN-GENERIC-005.
// Ces fonctions utilisent des identifiants predeclares comme type parameters.
package generic005

// badProcessString utilise "string" comme type parameter.
// Erreur: "string" est un identifiant predeclare.
//
// Params:
//   - s: valeur a traiter
//
// Returns:
//   - string: valeur traitee
func badProcessString[string any](s string) string { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return s
}

// badHandleError utilise "error" comme type parameter.
// Erreur: "error" est un identifiant predeclare.
//
// Params:
//   - e: valeur a traiter
//
// Returns:
//   - error: valeur traitee
func badHandleError[error any](e error) error { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return e
}

// badProcessInt utilise "int" comme type parameter.
// Erreur: "int" est un identifiant predeclare.
//
// Params:
//   - i: valeur a traiter
//
// Returns:
//   - int: valeur traitee
func badProcessInt[int any](i int) int { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return i
}

// badCheckBool utilise "bool" comme type parameter.
// Erreur: "bool" est un identifiant predeclare.
//
// Params:
//   - b: valeur a verifier
//
// Returns:
//   - bool: resultat
func badCheckBool[bool any](b bool) bool { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return b
}

// badProcessByte utilise "byte" comme type parameter.
// Erreur: "byte" est un identifiant predeclare.
//
// Params:
//   - b: valeur a traiter
//
// Returns:
//   - byte: valeur traitee
func badProcessByte[byte any](b byte) byte { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return b
}

// badProcessRune utilise "rune" comme type parameter.
// Erreur: "rune" est un identifiant predeclare.
//
// Params:
//   - r: valeur a traiter
//
// Returns:
//   - rune: valeur traitee
func badProcessRune[rune any](r rune) rune { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return r
}

// badUseLen utilise "len" comme type parameter.
// Erreur: "len" est une fonction predeclaree.
//
// Params:
//   - l: valeur a traiter
//
// Returns:
//   - len: valeur traitee
func badUseLen[len any](l len) len { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return l
}

// badUseMake utilise "make" comme type parameter.
// Erreur: "make" est une fonction predeclaree.
//
// Params:
//   - m: valeur a traiter
//
// Returns:
//   - make: valeur traitee
func badUseMake[make any](m make) make { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return m
}

// badUseNil utilise "nil" comme type parameter.
// Erreur: "nil" est une constante predeclaree.
//
// Params:
//   - n: valeur a traiter
//
// Returns:
//   - nil: valeur traitee
func badUseNil[nil any](n nil) nil { // want "KTN-GENERIC-005"
	// Retour de la valeur
	return n
}

// BadTypeString est un type generique avec "string" comme type parameter.
// Erreur: "string" est un identifiant predeclare.
type BadTypeString[string any] struct { // want "KTN-GENERIC-005"
	// value est la valeur stockee
	value string
}
