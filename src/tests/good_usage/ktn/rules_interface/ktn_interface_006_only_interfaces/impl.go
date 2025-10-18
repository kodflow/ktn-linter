// Package ktninterface004 fournit les implémentations des interfaces.
package ktninterface004

// service est l'implémentation privée de Service.
type service struct {
	status string
}

// NewService crée une nouvelle instance de Service.
//
// Returns:
//   - Service: l'instance du service
func NewService() Service {
	// Retourne une nouvelle instance du service
	return &service{status: "ready"}
}

// Process implémente Service.Process.
//
// Params:
//   - data: les données à traiter
//
// Returns:
//   - error: erreur si le traitement échoue
func (s *service) Process(data string) error {
	// Retourne nil car le traitement est terminé avec succès
	return nil
}

// GetStatus implémente Service.GetStatus.
//
// Returns:
//   - string: le statut actuel
func (s *service) GetStatus() string {
	// Retourne le statut actuel du service
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
	// Retourne une nouvelle instance du repository
	return &repository{data: make(map[string]string)}
}

// Save implémente Repository.Save.
//
// Params:
//   - data: les données à sauvegarder
//
// Returns:
//   - error: erreur si la sauvegarde échoue
func (r *repository) Save(data string) error {
	r.data["default"] = data
	// Retourne nil car la sauvegarde est réussie
	return nil
}

// Load implémente Repository.Load.
//
// Params:
//   - id: l'identifiant des données
//
// Returns:
//   - string: les données chargées
//   - error: erreur si le chargement échoue
func (r *repository) Load(id string) (string, error) {
	// Retourne les données chargées et nil pour l'erreur
	return r.data[id], nil
}
