package mock003

// Ce fichier s'appelle interfaces.go mais ne contient PAS d'interfaces
// Il ne devrait donc PAS déclencher d'erreur MOCK-001

// ConcreteType represents the struct.
type ConcreteType struct {
	Value int
	Name  string
}
// AnotherType represents the struct.

// AnotherType represents the struct.
type AnotherType struct {
	Data []byte
}

func NewConcreteType() *ConcreteType {
	return &ConcreteType{}
}
