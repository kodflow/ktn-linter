package mock002

// want `\[KTN-MOCK-002\] L'interface 'UserService' n'a pas de mock correspondant dans 'mock\.go'`
type UserService interface {
	GetUser(id int) (string, error)
}

// want `\[KTN-MOCK-002\] L'interface 'DataStore' n'a pas de mock correspondant dans 'mock\.go'`
type DataStore interface {
	Read(key string) ([]byte, error)
}
