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
		// Retourne une erreur car les données sont vides
		return errors.New("empty data")
	}
	// Retourne nil car le traitement est terminé avec succès
	return nil
}

// Close ferme les ressources.
//
// Returns:
//   - error: une erreur si la fermeture échoue
func (s *serviceImpl) Close() error {
	// Retourne nil car la fermeture est réussie
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
	// Retourne une nouvelle instance de l'interface
	return &serviceImpl{name: name}
}

// helperImpl implémente HelperInterface.
type helperImpl struct{}

// Help fournit de l'aide.
//
// Returns:
//   - string: le message d'aide
func (h *helperImpl) Help() string {
	// Retourne le message d'aide
	return "Help message"
}

// NewHelper crée une nouvelle instance de HelperInterface.
//
// Returns:
//   - HelperInterface: une nouvelle instance
func NewHelper() HelperInterface {
	// Retourne une nouvelle instance de l'interface
	return &helperImpl{}
}

// NewServiceInterface crée une nouvelle instance de ServiceInterface.
//
// Returns:
//   - ServiceInterface: l'instance créée
func NewServiceInterface() ServiceInterface {
	// Retourne nil comme placeholder
	return nil // Placeholder
}

// NewHelperInterface crée une nouvelle instance de HelperInterface.
//
// Returns:
//   - HelperInterface: l'instance créée
func NewHelperInterface() HelperInterface {
	// Retourne nil comme placeholder
	return nil // Placeholder
}
