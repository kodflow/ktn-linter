package test002 // want `\[KTN_TEST_002\] Fichier 'types.go' n'a pas de fichier de test correspondant`

// Fichier qui contient uniquement des types struct - devrait nécessiter un test
// Config represents the struct.
type Config struct {
	Host string
	Port int
}
// Settings represents the struct.

type Settings struct {
	Debug bool
}
