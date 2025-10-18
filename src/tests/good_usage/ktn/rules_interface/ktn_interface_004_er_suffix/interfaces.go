// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-004: Interfaces avec constructeurs (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════

package KTN_INTERFACE_006

// CacheManagerI006Good définit le contrat de gestion du cache.
type CacheManagerI006Good interface {
	// GetI006Good récupère une valeur du cache.
	//
	// Params:
	//   - key: la clé à rechercher
	//
	// Returns:
	//   - string: la valeur trouvée
	//   - bool: true si la clé existe
	GetI006Good(key string) (string, bool)

	// SetI006Good définit une valeur dans le cache.
	//
	// Params:
	//   - key: la clé
	//   - value: la valeur
	SetI006Good(key string, value string)
}

// LoggerI006Good définit le contrat de journalisation.
type LoggerI006Good interface {
	// InfoI006Good enregistre un message d'information.
	//
	// Params:
	//   - msg: le message à enregistrer
	InfoI006Good(msg string)

	// ErrorI006Good enregistre un message d'erreur.
	//
	// Params:
	//   - msg: le message d'erreur
	ErrorI006Good(msg string)
}

// MarkerInterfaceI006Good est une interface marqueur.
type MarkerInterfaceI006Good interface{}
