//go:build test
// +build test

package rules_var

// MockValidator est le mock de Validator.
type MockValidator struct {
	ValidateFunc func(value interface{}) error
}

// Validate implémente l'interface Validator.
//
// Params:
//   - value: la valeur à valider
//
// Returns:
//   - error: l'erreur de validation
func (m *MockValidator) Validate(value interface{}) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(value)
	}
	return nil
}
