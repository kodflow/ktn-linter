package test002_interfunc // want `\[KTN_TEST_002\] Fichier 'interfaces.go' n'a pas de fichier de test correspondant`

// Fichier interfaces.go avec fonction - devrait nécessiter test
// containsOnlyInterfaces002 retourne false car il y a une fonction
// Service defines the interface.
type Service interface {
	Method() error
}

func Helper() {
	// Cette fonction fait que containsOnlyInterfaces002 retourne false
}
