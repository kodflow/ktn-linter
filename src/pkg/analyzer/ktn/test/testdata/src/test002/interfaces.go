package test002

// Fichier interfaces.go qui contient uniquement des interfaces - devrait être ignoré
// Reader defines the interface.
type Reader interface {
	Read() string
}
// Writer defines the interface.

// Writer defines the interface.
type Writer interface {
	Write(data string) error
}
