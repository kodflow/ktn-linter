package test002

// Fichier interfaces.go qui contient uniquement des interfaces - devrait être ignoré
type Reader interface {
	Read() string
}

type Writer interface {
	Write(data string) error
}
