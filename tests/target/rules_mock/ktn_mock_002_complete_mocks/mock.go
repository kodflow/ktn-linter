//go:build test
// +build test

package complete_mocks

// MockService est le mock de Service.
type MockService struct {
	ProcessFunc func() error
}

// Process implémente Service.Process.
//
// Returns:
//   - error: erreur éventuelle
func (m *MockService) Process() error {
	if m.ProcessFunc != nil {
		return m.ProcessFunc()
	}
	return nil
}

// MockRepository est le mock de Repository.
type MockRepository struct {
	SaveFunc func(data string) error
	LoadFunc func() string
}

// Save implémente Repository.Save.
//
// Params:
//   - data: données à sauvegarder
//
// Returns:
//   - error: erreur éventuelle
func (m *MockRepository) Save(data string) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(data)
	}
	return nil
}

// Load implémente Repository.Load.
//
// Returns:
//   - string: données chargées
func (m *MockRepository) Load() string {
	if m.LoadFunc != nil {
		return m.LoadFunc()
	}
	return ""
}
