//go:build test
// +build test

package KTN_INTERFACE_002

// MockUserServiceI002Good est le mock de UserServiceI002Good.
type MockUserServiceI002Good struct {
	ProcessI002GoodFunc func(data string) error
}

// ProcessI002Good implémente l'interface UserServiceI002Good.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: une erreur si l'opération échoue
func (m *MockUserServiceI002Good) ProcessI002Good(data string) error {
	if m.ProcessI002GoodFunc != nil {
		// Early return from function.
		return m.ProcessI002GoodFunc(data)
	}
	// Early return from function.
	return nil
}

// MockOrderManagerI002Good est le mock de OrderManagerI002Good.
type MockOrderManagerI002Good struct {
	CreateOrderI002GoodFunc func(order string) error
}

// CreateOrderI002Good implémente l'interface OrderManagerI002Good.
//
// Params:
//   - order: la commande à créer
//
// Returns:
//   - error: une erreur si l'opération échoue
func (m *MockOrderManagerI002Good) CreateOrderI002Good(order string) error {
	if m.CreateOrderI002GoodFunc != nil {
		// Early return from function.
		return m.CreateOrderI002GoodFunc(order)
	}
	// Early return from function.
	return nil
}
