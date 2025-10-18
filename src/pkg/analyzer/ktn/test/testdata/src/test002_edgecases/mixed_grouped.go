package test002_edgecases // want `\[KTN_TEST_002\] Fichier 'mixed_grouped.go' n'a pas de fichier de test correspondant`

// Fichier avec interface et type alias groupés - devrait nécessiter test
// Service defines the interface.
type (
	Service interface {
		Method() error
	}
	MyAlias int // Ce n'est pas une interface donc containsOnlyInterfaces002 retourne false
)
