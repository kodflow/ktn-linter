package goodinterfaces

// ServiceInterface est l'interface publique principale du package.
type ServiceInterface interface {
	// Process traite une requête.
	//
	// Params:
	//   - data: les données à traiter
	//
	// Returns:
	//   - error: une erreur si le traitement échoue
	Process(data string) error

	// Close ferme les ressources.
	//
	// Returns:
	//   - error: une erreur si la fermeture échoue
	Close() error
}

// HelperInterface est une interface auxiliaire publique.
type HelperInterface interface {
	// Help fournit de l'aide.
	//
	// Returns:
	//   - string: le message d'aide
	Help() string
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
