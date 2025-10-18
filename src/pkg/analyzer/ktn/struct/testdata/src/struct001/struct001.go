package struct001

import "fmt"

// Variables et constantes (pas d'erreur)
var globalVar = 42

const globalConst = "test"

var (
	multiVar1 = 1
	multiVar2 = 2
)

const (
	multiConst1 = "a"
	multiConst2 = "b"
)

// CORRECT: Structs en MixedCaps

// UserConfig utilise MixedCaps.
type UserConfig struct {
	Name string
}

// httpClient utilise mixedCaps (privé).
type httpClient struct {
	URL string
}

// HTTPServer utilise un initialisme.
type HTTPServer struct {
	Port int
}

// BAD: Structs avec underscores

// user_config utilise snake_case (incorrect).
type user_config struct { // want "KTN-STRUCT-001.*MixedCaps"
	name string
}

// User_Profile mélange majuscules et underscores (incorrect).
type User_Profile struct { // want "KTN-STRUCT-001.*MixedCaps"
	ID int
}

// api_response utilise snake_case (incorrect).
type api_response struct { // want "KTN-STRUCT-001.*MixedCaps"
	data string
}

// test_struct_bad utilise snake_case (incorrect).
type test_struct_bad struct { // want "KTN-STRUCT-001.*MixedCaps"
	value int
}

// another_bad_name utilise snake_case (incorrect).
type another_bad_name struct { // want "KTN-STRUCT-001.*MixedCaps"
	field string
}

// Types non-struct (pas d'erreur)
type MyInt int
type MyString string
type MyFloat float64

// Déclarations groupées
type (
	// GroupedGood utilise MixedCaps.
	GroupedGood struct {
		Value int
	}

	grouped_bad struct { // want "KTN-STRUCT-001.*MixedCaps"
		value int
	}

	// AnotherType est un alias.
	AnotherType = string
)

// MyFunc est une fonction (pas d'erreur).
func MyFunc() {
	fmt.Println("test")
}

func anotherFunc() {
	_ = 1
}
