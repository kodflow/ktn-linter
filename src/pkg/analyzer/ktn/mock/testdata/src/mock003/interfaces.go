package mock003

// Ce fichier s'appelle interfaces.go mais ne contient PAS d'interfaces
// Il ne devrait donc PAS déclencher d'erreur MOCK-001

type ConcreteType struct {
	Value int
	Name  string
}

type AnotherType struct {
	Data []byte
}

func NewConcreteType() *ConcreteType {
	return &ConcreteType{}
}
