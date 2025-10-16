package incomplete_mocks

// Service est une interface de service.
type Service interface {
	Process() error
}

// Repository est une interface de repository.
type Repository interface {
	Save(data string) error
	Load() string
}
