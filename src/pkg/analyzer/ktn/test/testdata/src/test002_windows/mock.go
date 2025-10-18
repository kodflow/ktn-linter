package test002_windows

// Fichier mock.go - devrait être ignoré même avec des fonctions testables
// MockService represents the struct.
type MockService struct {
	Name string
}

func (m *MockService) Process() error {
	return nil
}
