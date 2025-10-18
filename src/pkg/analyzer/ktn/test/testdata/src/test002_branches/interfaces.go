package test002_branches

// Fichier interfaces.go avec seulement des interfaces - devrait être ignoré
// Reader defines the interface.
type Reader interface {
	Read() string
}
// Writer defines the interface.

type Writer interface {
	Write(data string) error
}
