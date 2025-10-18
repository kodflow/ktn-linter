package test002 // want `\[KTN_TEST_002\] Fichier 'mixed.go' n'a pas de fichier de test correspondant`

// Fichier qui mélange interface et struct - devrait nécessiter test (pas que des interfaces)
type ServiceInterface interface {
	Process() error
}

type ServiceImpl struct {
	Name string
}
