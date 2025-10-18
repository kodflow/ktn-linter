package test002 // want `\[KTN_TEST_002\] Fichier 'funconly.go' n'a pas de fichier de test correspondant`

// Fichier avec seulement des fonctions - devrait nécessiter test
func Helper() string {
	return "help"
}

func Process2(x int) int {
	return x * 2
}
