package predecl001

// L'analyzer détecte uniquement les déclarations var/const/type
// pas les short variable declarations (:=)

func BadShadowError() {
	var error = "custom error" // want `\[KTN-PREDECL-001\].*`
	_ = error
}

func BadShadowLen() {
	var len = 10 // want `\[KTN-PREDECL-001\].*`
	_ = len
}

func BadShadowTrue() {
	var true = false // want `\[KTN-PREDECL-001\].*`
	_ = true
}

func GoodCustomNames() {
	err := "custom error"
	length := 10
	isTrue := false
	_, _, _ = err, length, isTrue
}
