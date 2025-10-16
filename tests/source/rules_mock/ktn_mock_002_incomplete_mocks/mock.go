//go:build test
// +build test

package incomplete_mocks

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

// MockRepository manque intentionnellement pour tester KTN-MOCK-002
