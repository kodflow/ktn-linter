package KTN_INTERFACE_006

// cacheManagerImplI006 est l'implémentation du CacheManagerI006.
type cacheManagerImplI006 struct {
	data map[string]string
}

// GetI006 implémente l'interface.
//
// Params:
//   - key: la clé à rechercher
//
// Returns:
//   - string: la valeur trouvée
//   - bool: true si la clé existe
func (c *cacheManagerImplI006) GetI006(key string) (string, bool) {
	val, ok := c.data[key]
	// Early return from function.
	return val, ok
}

// SetI006 implémente l'interface.
//
// Params:
//   - key: la clé
//   - value: la valeur
func (c *cacheManagerImplI006) SetI006(key string, value string) {
	c.data[key] = value
}

// loggerImplI006 est l'implémentation du LoggerI006.
type loggerImplI006 struct {
	prefix string
}

// InfoI006 implémente l'interface.
//
// Params:
//   - msg: le message à enregistrer
func (l *loggerImplI006) InfoI006(msg string) {
	// Implementation
}

// ErrorI006 implémente l'interface.
//
// Params:
//   - msg: le message d'erreur
func (l *loggerImplI006) ErrorI006(msg string) {
	// Implementation
}

// NOTE: Pas de constructeurs NewCacheManagerI006 ni NewLoggerI006
// C'est la violation KTN-INTERFACE-004
