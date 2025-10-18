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
		// Early return from function.
		return m.ProcessFunc()
	}
	// Early return from function.
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
		// Early return from function.
		return m.SaveFunc(data)
	}
	// Early return from function.
	return nil
}

// Load implémente Repository.Load.
//
// Returns:
//   - string: données chargées
func (m *MockRepository) Load() string {
	if m.LoadFunc != nil {
		// Early return from function.
		return m.LoadFunc()
	}
	// Early return from function.
	return ""
}
