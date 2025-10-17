//go:build test
// +build test

package goodinterfaces

// MockServiceInterface est le mock de ServiceInterface.
type MockServiceInterface struct {
	ProcessFunc func(data string) error
	CloseFunc   func() error
}

// Process implémente l'interface ServiceInterface.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: une erreur si le traitement échoue
func (m *MockServiceInterface) Process(data string) error {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(data)
	}
	return nil
}

// Close implémente l'interface ServiceInterface.
//
// Returns:
//   - error: une erreur si la fermeture échoue
func (m *MockServiceInterface) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockHelperInterface est le mock de HelperInterface.
type MockHelperInterface struct {
	HelpFunc func() string
}

// Help implémente l'interface HelperInterface.
//
// Returns:
//   - string: le message d'aide
func (m *MockHelperInterface) Help() string {
	if m.HelpFunc != nil {
		return m.HelpFunc()
	}
	return ""
}
