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
