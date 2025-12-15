// Package var012 provides good test cases.
package var012

import "bytes"

// init demonstrates correct usage patterns
func init() {
	// Conversion une seule fois hors de la boucle
	data := [][]byte{[]byte("hello"), []byte("world")}
	target := "hello"
	targetBytes := []byte(target) // OK: Conversion une seule fois
	count := 0
	// Parcours des éléments
	for _, item := range data {
		// Vérification de la condition
		if bytes.Equal(item, targetBytes) {
			count++
		}
	}
	_ = count

	// Conversion une seule fois puis réutilisation
	byteData := []byte("test")
	str := string(byteData) // OK: Conversion une seule fois
	// Vérification de la condition
	if str == "hello" {
		println("found hello")
	}
	// Vérification de la condition
	if str == "world" {
		println("found world")
	}
	println(str)

	// Utilisation de bytes.Equal au lieu de string()
	items := [][]byte{[]byte("test")}
	targetTest := []byte("test")
	// Parcours des éléments
	for i := range len(items) {
		// Vérification de la condition
		if bytes.Equal(items[i], targetTest) { // OK: Pas de conversion string
			println("found")
		}
	}

	// Une seule utilisation
	singleData := []byte("unique")
	// Vérification de la condition
	if string(singleData) == "unique" { // OK: Une seule utilisation
		println("found unique")
	}

	// Chaque variable convertie une seule fois
	a := []byte("a")
	b := []byte("b")
	c := []byte("c")
	println(string(a))
	println(string(b))
	println(string(c))
}
