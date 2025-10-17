//go:build test
// +build test

package KTN_INTERFACE_001

// MockUserServiceI001Good est le mock de UserServiceI001Good.
type MockUserServiceI001Good struct {
	GetUserFunc func(id string) (string, error)
}

// GetUser implémente l'interface UserServiceI001Good.
//
// Params:
//   - id: l'identifiant de l'utilisateur
//
// Returns:
//   - string: le nom de l'utilisateur
//   - error: une erreur si l'opération échoue
func (m *MockUserServiceI001Good) GetUser(id string) (string, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(id)
	}
	return "", nil
}
