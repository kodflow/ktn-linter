// want `\[KTN-MOCK-001\] Le fichier 'interfaces\.go' contient des interfaces mais 'mock\.go' n'existe pas`
package mock001

type UserService interface {
	GetUser(id int) (string, error)
	SaveUser(name string) error
}

type DataStore interface {
	Read(key string) ([]byte, error)
	Write(key string, data []byte) error
}

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	Delete(key string)
}

type Logger interface {
	Log(message string)
	LogError(err error)
}

type Processor interface {
	Process(data []byte) error
	Validate(data []byte) bool
}
