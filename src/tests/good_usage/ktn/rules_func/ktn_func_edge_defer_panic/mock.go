//go:build test
// +build test

package rules_func

// MockResourceManager est le mock de ResourceManager.
type MockResourceManager struct {
	OpenFunc  func(name string) error
	CloseFunc func()
}

// Open implémente l'interface ResourceManager.
//
// Params:
//   - name: nom de la ressource à ouvrir
//
// Returns:
//   - error: erreur éventuelle
func (m *MockResourceManager) Open(name string) error {
	if m.OpenFunc != nil {
		// Early return from function.
		return m.OpenFunc(name)
	}
	// Early return from function.
	return nil
}

// Close implémente l'interface ResourceManager.
func (m *MockResourceManager) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}
