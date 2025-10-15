package goodinterfaces

import "errors"

// serviceImpl implémente ServiceInterface.
type serviceImpl struct {
	name string
}

// Process traite une requête.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: une erreur si le traitement échoue
func (s *serviceImpl) Process(data string) error {
	if data == "" {
		return errors.New("empty data")
	}
	return nil
}

// Close ferme les ressources.
//
// Returns:
//   - error: une erreur si la fermeture échoue
func (s *serviceImpl) Close() error {
	return nil
}

// NewService crée une nouvelle instance de ServiceInterface.
//
// Params:
//   - name: le nom du service
//
// Returns:
//   - ServiceInterface: une nouvelle instance
func NewService(name string) ServiceInterface {
	return &serviceImpl{name: name}
}

// helperImpl implémente HelperInterface.
type helperImpl struct{}

// Help fournit de l'aide.
//
// Returns:
//   - string: le message d'aide
func (h *helperImpl) Help() string {
	return "Help message"
}

// NewHelper crée une nouvelle instance de HelperInterface.
//
// Returns:
//   - HelperInterface: une nouvelle instance
func NewHelper() HelperInterface {
	return &helperImpl{}
}
