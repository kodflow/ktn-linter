package good

// Bon : interfaces.go contient au moins une interface publique
type Service interface {
	DoSomething() string
}

type Repository interface {
	Save(data string) error
}
