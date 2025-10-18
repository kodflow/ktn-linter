package predecl002

func BadShadowError() {
	// want `\[KTN-OPS-PREDECL-002\] Shadowing de l'identifiant prédéclaré 'error'`
	error := "custom error"
	_ = error
}

func BadShadowLen() {
	// want `\[KTN-OPS-PREDECL-002\] Shadowing de l'identifiant prédéclaré 'len'`
	len := 10
	_ = len
}

func BadShadowTrue() {
	// want `\[KTN-OPS-PREDECL-002\] Shadowing de l'identifiant prédéclaré 'true'`
	true := false
	_ = true
}

func GoodCustomNames() {
	err := "custom error"
	length := 10
	isTrue := false
	_, _, _ = err, length, isTrue
}
