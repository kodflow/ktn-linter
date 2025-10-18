package test002 // want `\[KTN_TEST_002\] Fichier 'notested.go' n'a pas de fichier de test correspondant`

// Mauvais : fichier sans _test.go correspondant
func DoSomething() string {
	return "something"
}

func Process(data string) error {
	return nil
}
