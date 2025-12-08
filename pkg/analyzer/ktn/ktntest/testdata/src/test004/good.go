// Package test004 provides resource management utilities.
package test004

// GoodResourceInterface defines the public API for GoodResource.
type GoodResourceInterface interface {
	Metadata() string
	Schema() map[string]string
	Configure(config string) error
	Name() string
}

// GoodResource représente une ressource avec des méthodes publiques.
// Toutes les méthodes ont des tests correspondants.
type GoodResource struct {
	name string
}

// NewGoodResource crée une nouvelle instance.
//
// Returns:
//   - *GoodResource: nouvelle instance
func NewGoodResource() *GoodResource {
	r := &GoodResource{}
	// Use private functions to avoid dead code
	r.name = r.sanitize("init")
	_ = validateConfig(r.name)
	// Retour de la nouvelle instance
	return r
}

// Metadata retourne les métadonnées.
//
// Returns:
//   - string: nom de la ressource
func (r *GoodResource) Metadata() string {
	// Retour des métadonnées
	return "good_resource"
}

// Schema retourne le schéma.
//
// Returns:
//   - map[string]string: schéma de la ressource
func (r *GoodResource) Schema() map[string]string {
	// Retour du schéma
	return map[string]string{"type": "test"}
}

// Configure configure la ressource.
//
// Params:
//   - config: configuration à appliquer
//
// Returns:
//   - error: erreur éventuelle
func (r *GoodResource) Configure(config string) error {
	r.name = config
	// Retour succès
	return nil
}

// validateConfig valide la configuration (fonction privée avec test).
//
// Params:
//   - config: configuration à valider
//
// Returns:
//   - bool: true si valide
func validateConfig(config string) bool {
	// Validation
	return len(config) > 0
}

// sanitize nettoie les données (fonction privée avec test).
//
// Params:
//   - data: données à nettoyer
//
// Returns:
//   - string: données nettoyées
func (r *GoodResource) sanitize(data string) string {
	// Nettoyage
	return data + "_sanitized"
}

// Name retourne le nom de la ressource.
//
// Returns:
//   - string: nom de la ressource
func (r *GoodResource) Name() string {
	// Retour du nom
	return r.name
}
