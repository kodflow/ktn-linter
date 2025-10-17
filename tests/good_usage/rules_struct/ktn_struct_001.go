package rules_struct

// ✅ CONFORME KTN-STRUCT-001
// Les structs utilisent MixedCaps (ou mixedCaps pour les privées)

// userConfig contient la configuration utilisateur.
type userConfig struct {
	// Host est l'hôte du serveur
	Host string
	// Port est le port d'écoute
	Port int
}

// apiClient représente un client API.
type apiClient struct {
	// BaseURL est l'URL de base de l'API
	BaseURL string
	// Timeout est le délai d'expiration
	Timeout int
}

// databaseConnection représente une connexion à la base de données.
type databaseConnection struct {
	// Host est l'hôte de la base de données
	Host string
	// Port est le port de connexion
	Port int
	// Database est le nom de la base
	Database string
	// Username est le nom d'utilisateur
	Username string
	// Password est le mot de passe
	Password string
}

// httpRequestConfig contient la configuration d'une requête HTTP.
type httpRequestConfig struct {
	// Method est la méthode HTTP
	Method string
	// URL est l'URL de la requête
	URL string
	// Headers contient les en-têtes HTTP
	Headers map[string]string
	// Body contient le corps de la requête
	Body []byte
}

// cacheEntry représente une entrée dans le cache.
type cacheEntry struct {
	// Key est la clé de l'entrée
	Key string
	// Value est la valeur stockée
	Value interface{}
	// ExpiresAt est la date d'expiration
	ExpiresAt int64
}

// workerPoolConfig contient la configuration d'un pool de workers.
type workerPoolConfig struct {
	// MinWorkers est le nombre minimum de workers
	MinWorkers int
	// MaxWorkers est le nombre maximum de workers
	MaxWorkers int
	// QueueSize est la taille de la queue
	QueueSize int
}

// errorResponse représente une réponse d'erreur.
type errorResponse struct {
	// Code est le code d'erreur HTTP
	Code int
	// Message est le message d'erreur
	Message string
	// Details contient les détails supplémentaires
	Details map[string]interface{}
}

// authToken représente un token d'authentification.
type authToken struct {
	// Token est le token JWT
	Token string
	// ExpiresAt est la date d'expiration
	ExpiresAt int64
	// UserID est l'identifiant de l'utilisateur
	UserID string
}

// internalCache est une struct privée bien nommée.
type internalCache struct {
	// entries contient les entrées du cache
	entries map[string]interface{}
	// maxSize est la taille maximale
	maxSize int
}

// httpClient est un client HTTP privé.
type httpClient struct {
	// baseURL est l'URL de base
	baseURL string
	// timeout est le délai d'expiration
	timeout int
}
