package interface002

// Bon : struct privée (implémentation)
type serviceImpl struct {
	name string
}

func (s *serviceImpl) DoSomething() string {
	return s.name
}

// Exception : types de données publics autorisés
// UserConfig represents the struct.
type UserConfig struct {
	Name string
	Age  int
}

// StatusType is a custom type.
type StatusType string
