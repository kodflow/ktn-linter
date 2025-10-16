package complete_mocks

// Service est une interface de service.
type Service interface {
	Process() error
}

// NewService crée une nouvelle instance de Service.
//
// Returns:
//   - Service: nouvelle instance
func NewService() Service {
	return nil
}

// Repository est une interface de repository.
type Repository interface {
	Save(data string) error
	Load() string
}

// NewRepository crée une nouvelle instance de Repository.
//
// Returns:
//   - Repository: nouvelle instance
func NewRepository() Repository {
	return nil
}
