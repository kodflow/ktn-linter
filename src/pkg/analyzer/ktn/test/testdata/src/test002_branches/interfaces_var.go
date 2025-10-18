package test002_branches // want `\[KTN_TEST_002\] Fichier 'interfaces_var.go' n'a pas de fichier de test correspondant`

// Fichier avec GenDecl non-TYPE (var) et interface - devrait nécessiter test
// Ce fichier teste que la présence de var + interface nécessite un test
var GlobalVar = 10

type Service interface {
	Method() error
}
