package struct002

import "fmt"

// Variables et constantes (pas d'erreur)
var globalVar = 42

const globalConst = "test"

// CORRECT: Structs avec documentation godoc

// UserConfig représente la configuration utilisateur.
type UserConfig struct {
	Name string
}

// HTTPClient représente un client HTTP.
type HTTPClient struct {
	URL string
}

// ServerConfig représente la configuration du serveur.
type ServerConfig struct {
	Port int
	Host string
}

// MinimalStruct est minimale.
type MinimalStruct struct {
	X int
}

// BAD: Structs sans documentation godoc

type BadNoDoc struct { // want "KTN-STRUCT-002.*commentaire godoc"
	Field string
}

type AnotherBadStruct struct { // want "KTN-STRUCT-002.*commentaire godoc"
	ID   int
	Name string
}

type ThirdBadStruct struct { // want "KTN-STRUCT-002.*commentaire godoc"
	Value float64
}

type FourthBadStruct struct { // want "KTN-STRUCT-002.*commentaire godoc"
	Data string
}

type FifthBadStruct struct { // want "KTN-STRUCT-002.*commentaire godoc"
	Count int
}

// Types non-struct (pas d'erreur)
type MyInt int
type MyString string
type MyFloat float64

// MyInterface est une interface (pas d'erreur).
type MyInterface interface {
	Method()
}

// MyFunc est une fonction (pas d'erreur).
func MyFunc() {
	fmt.Println("test")
}

// GroupedGood est documenté.
type GroupedGood struct {
	Value int
}

type GroupedBad struct { // want "KTN-STRUCT-002.*commentaire godoc"
	value int
}
