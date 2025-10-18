package test002_edgecases // want `\[KTN_TEST_002\] Fichier 'mixed_grouped.go' n'a pas de fichier de test correspondant`

// Fichier avec interface et type alias groupés - devrait nécessiter test
type (
	// Service defines the interface.
	Service interface {
		Method() error
	}
	// MyAlias is a custom type.
	MyAlias int // Ce n'est pas une interface donc containsOnlyInterfaces002 retourne false
)
