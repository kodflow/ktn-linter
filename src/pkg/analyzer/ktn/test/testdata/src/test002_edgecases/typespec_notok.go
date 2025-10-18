package test002_edgecases

// Fichier pour tester la branche typeSpec, ok := spec.(*ast.TypeSpec); !ok
// Cependant GenDecl.Specs contient toujours des TypeSpec pour token.TYPE
// Donc ce cas est difficile à reproduire naturellement
const (
	A = 1
	B = 2
)
