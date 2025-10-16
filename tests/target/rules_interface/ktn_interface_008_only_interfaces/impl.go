// Package ktninterface008 fournit les implémentations des interfaces.
package ktninterface008

// service est l'implémentation privée de Service.
type service struct {
	status string
}

// NewService crée une nouvelle instance de Service.
//
// Returns:
//   - Service: l'instance du service
func NewService() Service {
	return &service{status: "ready"}
}

// Process implémente Service.Process.
func (s *service) Process(data string) error {
	return nil
}

// GetStatus implémente Service.GetStatus.
func (s *service) GetStatus() string {
	return s.status
}

// repository est l'implémentation privée de Repository.
type repository struct {
	data map[string]string
}

// NewRepository crée une nouvelle instance de Repository.
//
// Returns:
//   - Repository: l'instance du repository
func NewRepository() Repository {
	return &repository{data: make(map[string]string)}
}

// Save implémente Repository.Save.
func (r *repository) Save(data string) error {
	r.data["default"] = data
	return nil
}

// Load implémente Repository.Load.
func (r *repository) Load(id string) (string, error) {
	return r.data[id], nil
}
