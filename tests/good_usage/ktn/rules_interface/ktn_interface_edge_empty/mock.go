//go:build test
// +build test

package goodempty

// MockProcessor est le mock de Processor.
type MockProcessor struct {
	ProcessFunc func(data string) string
}

// Process implémente l'interface Processor.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat du traitement
func (m *MockProcessor) Process(data string) string {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(data)
	}
	return ""
}

// MockCache est le mock de Cache.
type MockCache struct {
	GetFunc func(key string) (string, bool)
	SetFunc func(key string, value string)
}

// Get implémente l'interface Cache.
//
// Params:
//   - key: clé de la valeur
//
// Returns:
//   - string: valeur associée à la clé
//   - bool: true si la clé existe
func (m *MockCache) Get(key string) (string, bool) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return "", false
}

// Set implémente l'interface Cache.
//
// Params:
//   - key: clé de la valeur
//   - value: valeur à stocker
func (m *MockCache) Set(key string, value string) {
	if m.SetFunc != nil {
		m.SetFunc(key, value)
	}
}

// MockContainer est le mock de Container.
type MockContainer[T any] struct {
	AddFunc    func(item T)
	GetAllFunc func() []T
}

// Add implémente l'interface Container.
//
// Params:
//   - item: élément à ajouter
func (m *MockContainer[T]) Add(item T) {
	if m.AddFunc != nil {
		m.AddFunc(item)
	}
}

// GetAll implémente l'interface Container.
//
// Returns:
//   - []T: tous les éléments du conteneur
func (m *MockContainer[T]) GetAll() []T {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil
}
