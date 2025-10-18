package test002_interstruct // want `\[KTN_TEST_002\] Fichier 'interfaces.go' n'a pas de fichier de test correspondant`

// Fichier interfaces.go avec struct - devrait nécessiter test
// containsOnlyInterfaces002 retourne false car il y a un type non-interface
type Service interface {
	Method() error
}

type Config struct {
	Name string
}
