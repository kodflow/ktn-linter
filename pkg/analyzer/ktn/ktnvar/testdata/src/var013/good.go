// Good examples for the var013 test case.
package var013

import "bytes"

// goodPreallocatedConversion préalloue la conversion hors de la boucle.
//
// Params:
//   - data: données à traiter
//   - target: cible à rechercher
//
// Returns:
//   - int: nombre d'occurrences
func goodPreallocatedConversion(data [][]byte, target string) int {
	count := 0
	targetBytes := []byte(target) // OK: Conversion une seule fois
	// Parcours des éléments
	for _, item := range data {
		// Vérification de la condition
		if bytes.Equal(item, targetBytes) {
			count++
		}
	}
	// Retour de la fonction
	return count
}

// goodSingleConversion convertit une seule fois.
//
// Params:
//   - data: données à convertir
func goodSingleConversion(data []byte) {
	str := string(data) // OK: Conversion une seule fois puis réutilisation
	// Vérification de la condition
	if str == "hello" {
		println("found hello")
	}
	// Vérification de la condition
	if str == "world" {
		println("found world")
	}
	println(str)
}

// goodBytesEqual utilise bytes.Equal au lieu de string().
//
// Params:
//   - items: éléments à comparer
func goodBytesEqual(items [][]byte) {
	target := []byte("test")
	// Parcours des éléments
	for i := range len(items) {
		// Vérification de la condition
		if bytes.Equal(items[i], target) { // OK: Pas de conversion string
			println("found")
		}
	}
}

// goodNestedWithPrealloc préalloue dans la boucle externe.
//
// Params:
//   - matrix: matrice à parcourir
func goodNestedWithPrealloc(matrix [][][]byte) {
	target := []byte("x")
	// Parcours des éléments
	for _, row := range matrix {
		// Parcours des éléments
		for _, cell := range row {
			// Vérification de la condition
			if bytes.Equal(cell, target) { // OK: Pas de conversion
				println("found x")
			}
		}
	}
}

// goodSingleUseConversion utilise string() une seule fois.
//
// Params:
//   - data: données à vérifier
func goodSingleUseConversion(data []byte) {
	// Vérification de la condition
	if string(data) == "unique" { // OK: Une seule utilisation
		println("found unique")
	}
}

// goodDifferentVariables convertit des variables différentes.
//
// Params:
//   - a: première variable
//   - b: deuxième variable
//   - c: troisième variable
func goodDifferentVariables(a []byte, b []byte, c []byte) {
	println(string(a)) // OK: Chaque variable convertie une seule fois
	println(string(b))
	println(string(c))
}

// init utilise les fonctions privées
func init() {
	// Appel de goodPreallocatedConversion
	_ = goodPreallocatedConversion(nil, "")
	// Appel de goodSingleConversion
	goodSingleConversion(nil)
	// Appel de goodBytesEqual
	goodBytesEqual(nil)
	// Appel de goodNestedWithPrealloc
	goodNestedWithPrealloc(nil)
	// Appel de goodSingleUseConversion
	goodSingleUseConversion(nil)
	// Appel de goodDifferentVariables
	goodDifferentVariables(nil, nil, nil)
}
