// Tests for KTN-INTERFACE-006 (version corrigée)
package rules_interface

// KTN-INTERFACE-006 GOOD: Les interfaces ont maintenant des constructeurs

// cacheManagerImplI006Good est l'implémentation du CacheManagerI006Good.
type cacheManagerImplI006Good struct {
	data map[string]string
}

// GetI006Good implémente l'interface.
//
// Params:
//   - key: la clé à rechercher
//
// Returns:
//   - string: la valeur trouvée
//   - bool: true si la clé existe
func (c *cacheManagerImplI006Good) GetI006Good(key string) (string, bool) {
	val, ok := c.data[key]
	return val, ok
}

// SetI006Good implémente l'interface.
//
// Params:
//   - key: la clé
//   - value: la valeur
func (c *cacheManagerImplI006Good) SetI006Good(key string, value string) {
	c.data[key] = value
}

// NewCacheManagerI006Good crée une nouvelle instance du cache.
//
// Returns:
//   - CacheManagerI006Good: une nouvelle instance
func NewCacheManagerI006Good() CacheManagerI006Good {
	return &cacheManagerImplI006Good{
		data: make(map[string]string),
	}
}

// loggerImplI006Good est l'implémentation du LoggerI006Good.
type loggerImplI006Good struct {
	prefix string
}

// InfoI006Good implémente l'interface.
//
// Params:
//   - msg: le message à enregistrer
func (l *loggerImplI006Good) InfoI006Good(msg string) {
	// Implementation
}

// ErrorI006Good implémente l'interface.
//
// Params:
//   - msg: le message d'erreur
func (l *loggerImplI006Good) ErrorI006Good(msg string) {
	// Implementation
}

// NewLoggerI006Good crée une nouvelle instance du logger.
//
// Params:
//   - prefix: le préfixe pour les messages
//
// Returns:
//   - LoggerI006Good: une nouvelle instance
func NewLoggerI006Good(prefix string) LoggerI006Good {
	return &loggerImplI006Good{prefix: prefix}
}
