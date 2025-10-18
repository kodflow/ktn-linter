package test002_func // want `\[KTN_TEST_002\] Fichier 'withfunc.go' n'a pas de fichier de test correspondant`

// Fichier pour tester la branche FuncDecl dans containsOnlyInterfaces002
// Ce n'est pas interfaces.go donc il faut un test
type Service interface {
	Method() error
}

func Helper() {
	// Function
}
