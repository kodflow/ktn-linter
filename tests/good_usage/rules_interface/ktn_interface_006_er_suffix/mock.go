//go:build test
// +build test

package KTN_INTERFACE_006

// MockCacheManagerI006Good est le mock de CacheManagerI006Good.
type MockCacheManagerI006Good struct {
	GetI006GoodFunc func(key string) (string, bool)
	SetI006GoodFunc func(key string, value string)
}

// GetI006Good implémente l'interface CacheManagerI006Good.
//
// Params:
//   - key: la clé à rechercher
//
// Returns:
//   - string: la valeur trouvée
//   - bool: true si la clé existe
func (m *MockCacheManagerI006Good) GetI006Good(key string) (string, bool) {
	if m.GetI006GoodFunc != nil {
		return m.GetI006GoodFunc(key)
	}
	return "", false
}

// SetI006Good implémente l'interface CacheManagerI006Good.
//
// Params:
//   - key: la clé
//   - value: la valeur
func (m *MockCacheManagerI006Good) SetI006Good(key string, value string) {
	if m.SetI006GoodFunc != nil {
		m.SetI006GoodFunc(key, value)
	}
}

// MockLoggerI006Good est le mock de LoggerI006Good.
type MockLoggerI006Good struct {
	InfoI006GoodFunc  func(msg string)
	ErrorI006GoodFunc func(msg string)
}

// InfoI006Good implémente l'interface LoggerI006Good.
//
// Params:
//   - msg: le message à enregistrer
func (m *MockLoggerI006Good) InfoI006Good(msg string) {
	if m.InfoI006GoodFunc != nil {
		m.InfoI006GoodFunc(msg)
	}
}

// ErrorI006Good implémente l'interface LoggerI006Good.
//
// Params:
//   - msg: le message d'erreur
func (m *MockLoggerI006Good) ErrorI006Good(msg string) {
	if m.ErrorI006GoodFunc != nil {
		m.ErrorI006GoodFunc(msg)
	}
}

// MockMarkerInterfaceI006Good est le mock de MarkerInterfaceI006Good.
type MockMarkerInterfaceI006Good struct {
}
