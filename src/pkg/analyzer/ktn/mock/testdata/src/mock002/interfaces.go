package mock002

// UserService a un mock (MockUserService dans mock.go) - OK
type UserService interface {
	GetUser(id int) (string, error)
}

// DataStore n'a PAS de mock - devrait générer un diagnostic
type DataStore interface { // want `\[KTN-MOCK-002\] L'interface 'DataStore' n'a pas de mock correspondant dans 'mock\.go'`
	Read(key string) ([]byte, error)
}
