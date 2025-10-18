// Package ktninterface004 démontre le respect de KTN-INTERFACE-004.
// Ce fichier contient UNIQUEMENT des interfaces (pas de structs).
package ktninterface004

// Service définit l'interface du service principal.
type Service interface {
	// Process traite les données.
	//
	// Params:
	//   - data: les données à traiter
	//
	// Returns:
	//   - error: erreur si le traitement échoue
	Process(data string) error

	// GetStatus retourne le statut du service.
	//
	// Returns:
	//   - string: le statut actuel
	GetStatus() string
}

// Repository définit l'interface du repository.
type Repository interface {
	// Save sauvegarde les données.
	//
	// Params:
	//   - data: les données à sauvegarder
	//
	// Returns:
	//   - error: erreur si la sauvegarde échoue
	Save(data string) error

	// Load charge les données.
	//
	// Params:
	//   - id: l'identifiant des données
	//
	// Returns:
	//   - string: les données chargées
	//   - error: erreur si le chargement échoue
	Load(id string) (string, error)
}
