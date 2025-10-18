package struct003

import "fmt"

// Variables et constantes (pas d'erreur)
var globalVar = 42

const globalConst = "test"

// CORRECT: Champs exportés documentés

// GoodConfig représente une configuration.
type GoodConfig struct {
	// Name est le nom de la configuration.
	Name string

	// Port est le port d'écoute.
	Port int

	// Les champs privés n'ont pas besoin de documentation
	internal string
}

// EmptyStruct est une struct vide (pas d'erreur).
type EmptyStruct struct {
}

// OnlyPrivateFields contient uniquement des champs privés.
type OnlyPrivateFields struct {
	field1 string
	field2 int
	field3 bool
}

// ComplexStruct avec plusieurs champs bien documentés.
type ComplexStruct struct {
	// Field1 est documenté.
	Field1 string
	// Field2 est documenté.
	Field2 int
	// Field3 est documenté.
	Field3 bool
	// Field4 est documenté.
	Field4 float64
	private string
}

// BAD: Champs exportés non documentés

// BadConfig représente une configuration.
type BadConfig struct {
	Name string // want "KTN-STRUCT-003.*sans commentaire"

	Port int // want "KTN-STRUCT-003.*sans commentaire"

	// Champ privé (pas d'erreur)
	internal string
}

// AnotherBadStruct avec champs non documentés.
type AnotherBadStruct struct {
	ID   int    // want "KTN-STRUCT-003.*sans commentaire"
	Data string // want "KTN-STRUCT-003.*sans commentaire"
}

// MixedFields avec un champ exporté non documenté.
type MixedFields struct {
	ExportedField string // want "KTN-STRUCT-003.*sans commentaire"
	privateField  int
}

// ThirdBadStruct avec plusieurs champs non documentés.
type ThirdBadStruct struct {
	Field1 string // want "KTN-STRUCT-003.*sans commentaire"
	Field2 int    // want "KTN-STRUCT-003.*sans commentaire"
	Field3 bool   // want "KTN-STRUCT-003.*sans commentaire"
	private int
}

// FourthBadStruct avec un seul champ non documenté.
type FourthBadStruct struct {
	UndocField string // want "KTN-STRUCT-003.*sans commentaire"
}

// Types non-struct (pas d'erreur)
type MyInt int
type MyString string
type MyFloat float64

// MyFunc est une fonction (pas d'erreur).
func MyFunc() {
	fmt.Println("test")
}

func anotherFunc() {
	_ = 1
}
