// want `\[KTN-MOCK-001\] Le fichier 'interfaces\.go' contient des interfaces mais 'mock\.go' n'existe pas`
package mock001

// UserService defines the interface.
type UserService interface {
	GetUser(id int) (string, error)
	SaveUser(name string) error
}
// DataStore defines the interface.

// DataStore defines the interface.
type DataStore interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
// Cache defines the interface.
}

// Cache defines the interface.
type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
// Logger defines the interface.
	Delete(key string)
}

// Logger defines the interface.
type Logger interface {
// Processor defines the interface.
	Log(message string)
	LogError(err error)
}

// Processor defines the interface.
type Processor interface {
	Process(data []byte) error
	Validate(data []byte) bool
}
