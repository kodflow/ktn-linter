package test002_branches // want `\[KTN_TEST_002\] Fichier 'interfaces_var.go' n'a pas de fichier de test correspondant`

// Fichier avec GenDecl non-TYPE (const) et interface - devrait nécessiter test
// Ce fichier teste que la présence de const + interface nécessite un test
const GlobalVar = 10

// Service defines the interface.
type Service interface {
	Method() error
}
