package good

// Bon : interfaces.go contient au moins une interface publique
// Service defines the interface.
type Service interface {
	DoSomething() string
}
// Repository defines the interface.

type Repository interface {
	Save(data string) error
}
