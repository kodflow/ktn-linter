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
