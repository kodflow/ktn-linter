// Package rules_func contient les fonctions de test edge defer/panic.
package rules_func

// ResourceManager gère les ressources.
type ResourceManager interface {
	Open(name string) error
	Close()
}
