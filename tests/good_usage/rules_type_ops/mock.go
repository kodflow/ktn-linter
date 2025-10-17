//go:build test
// +build test

package rules_type_ops

// MockMyInterfaceGood est le mock de MyInterfaceGood.
type MockMyInterfaceGood struct {
	// DoFunc permet de mocker la méthode Do.
	DoFunc func()
}

// Do implémente l'interface MyInterfaceGood.
func (m *MockMyInterfaceGood) Do() {
	if m.DoFunc != nil {
		m.DoFunc()
	}
}
