// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-006: Interface sans constructeur
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Chaque interface publique avec méthodes doit avoir un constructeur New*
//    qui retourne l'interface (pas l'implémentation concrète).
//
//    POURQUOI :
//    - Encapsule la création et cache l'implémentation
//    - Facilite le changement d'implémentation sans changer l'API
//    - Pattern standard Go pour dependency injection
//    - Rend les tests plus simples avec injection de mocks
//
// ❌ CAS INCORRECT 1: Interface sans constructeur (SEULE ERREUR: KTN-INTERFACE-006)
// NOTE: CacheManagerI006 n'a pas de NewCacheManagerI006
// ERREUR ATTENDUE: KTN-INTERFACE-006 sur CacheManagerI006
//
// ❌ CAS INCORRECT 2: Interface sans constructeur (SEULE ERREUR: KTN-INTERFACE-006)
// NOTE: LoggerI006 n'a pas de NewLoggerI006
// ERREUR ATTENDUE: KTN-INTERFACE-006 sur LoggerI006
//
// ✅ CAS PARFAIT (voir target/) :
//    // interfaces.go
//    type CacheManager interface {
//        Get(key string) (string, bool)
//        Set(key string, value string)
//    }
//
//    func NewCacheManager() CacheManager {
//        return &cacheManagerImpl{...}
//    }
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_INTERFACE_006

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
