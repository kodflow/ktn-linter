// Package rules_func contient les fonctions de test edge defer/panic.
package rules_func

// ResourceManager gère les ressources.
type ResourceManager interface {
	Open(name string) error
	Close()
}

// NewResourceManager crée un nouveau gestionnaire de ressources.
//
// Returns:
//   - ResourceManager: instance du gestionnaire
func NewResourceManager() ResourceManager {
	// Retourne nil comme placeholder
	return nil // Placeholder
}
