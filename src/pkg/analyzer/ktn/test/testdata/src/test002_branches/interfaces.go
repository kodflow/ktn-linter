package test002_branches

// Fichier interfaces.go avec seulement des interfaces - devrait être ignoré
type Reader interface {
	Read() string
}

type Writer interface {
	Write(data string) error
}
