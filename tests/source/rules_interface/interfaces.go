// interfaces.go pour les tests KTN-INTERFACE-006
package rules_interface

// KTN-INTERFACE-006: Ces interfaces n'ont pas de constructeurs

// CacheManagerI006 définit le contrat de gestion du cache.
// Cette interface devrait avoir un constructeur NewCacheManagerI006.
type CacheManagerI006 interface {
	// GetI006 récupère une valeur du cache.
	//
	// Params:
	//   - key: la clé à rechercher
	//
	// Returns:
	//   - string: la valeur trouvée
	//   - bool: true si la clé existe
	GetI006(key string) (string, bool)

	// SetI006 définit une valeur dans le cache.
	//
	// Params:
	//   - key: la clé
	//   - value: la valeur
	SetI006(key string, value string)
}

// LoggerI006 définit le contrat de journalisation.
// Cette interface devrait avoir un constructeur NewLoggerI006.
type LoggerI006 interface {
	// InfoI006 enregistre un message d'information.
	//
	// Params:
	//   - msg: le message à enregistrer
	InfoI006(msg string)

	// ErrorI006 enregistre un message d'erreur.
	//
	// Params:
	//   - msg: le message d'erreur
	ErrorI006(msg string)
}

// MarkerInterfaceI006 est une interface marqueur (pas de méthodes).
// Pas besoin de constructeur pour celle-ci car elle n'a pas de méthodes.
type MarkerInterfaceI006 interface{}
