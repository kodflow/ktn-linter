package test002 // want `\[KTN_TEST_002\] Fichier 'types.go' n'a pas de fichier de test correspondant`

// Fichier qui contient uniquement des types struct - devrait nécessiter un test
type Config struct {
	Host string
	Port int
}

type Settings struct {
	Debug bool
}
