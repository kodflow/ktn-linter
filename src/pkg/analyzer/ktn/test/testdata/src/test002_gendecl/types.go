package test002_gendecl // want `\[KTN_TEST_002\] Fichier 'types.go' n'a pas de fichier de test correspondant`

// Fichier pour tester les différentes branches de GenDecl
// Writer defines the interface.
// Reader defines the interface.
type (
	Reader interface {
		Read() string
	}
	Writer interface {
		Write(data string) error
	}
)
