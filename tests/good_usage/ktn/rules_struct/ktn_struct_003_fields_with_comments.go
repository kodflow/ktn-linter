package rules_struct

// ✅ CONFORME KTN-STRUCT-003
// Tous les champs exportés sont documentés

// commentedUserConfig contient la configuration utilisateur.
type commentedUserConfig struct {
	// Host est l'hôte du serveur
	Host string
	// Port est le port d'écoute
	Port int
	// privé non documenté (OK car privé)
	timeout int
}

// commentedAPIClient représente un client API.
type commentedAPIClient struct {
	// BaseURL est l'URL de base de l'API
	BaseURL string
	// Timeout est le délai d'expiration en secondes
	Timeout int
	// MaxRetries est le nombre maximum de tentatives
	MaxRetries int
	// interne non documenté (OK car privé)
	httpClient interface{}
}

// commentedDatabaseConnection représente une connexion à la base de données.
type commentedDatabaseConnection struct {
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
	// privé
	connected bool
}

// commentedHTTPRequestConfig contient la configuration d'une requête HTTP.
type commentedHTTPRequestConfig struct {
	// Method est la méthode HTTP
	Method string
	// URL est l'URL de la requête
	URL string
	// Headers contient les en-têtes HTTP
	Headers map[string]string
	// Body contient le corps de la requête
	Body []byte
	// Timeout est le délai d'expiration
	Timeout int
	// privé
	attempt int
}

// commentedCacheEntry représente une entrée dans le cache.
type commentedCacheEntry struct {
	// Key est la clé de l'entrée
	Key string
	// Value est la valeur stockée
	Value interface{}
	// ExpiresAt est la date d'expiration
	ExpiresAt int64
	// Priority est la priorité de l'entrée
	Priority int
	// interne
	lastAccessed int64
}

// commentedWorkerPoolConfig contient la configuration d'un pool de workers.
type commentedWorkerPoolConfig struct {
	// MinWorkers est le nombre minimum de workers
	MinWorkers int
	// MaxWorkers est le nombre maximum de workers
	MaxWorkers int
	// QueueSize est la taille de la queue
	QueueSize int
	// IdleTimeout est le délai avant arrêt d'un worker inactif
	IdleTimeout int
	// privé
	started bool
}

// commentedErrorResponse représente une réponse d'erreur.
type commentedErrorResponse struct {
	// Code est le code d'erreur HTTP
	Code int
	// Message est le message d'erreur
	Message string
	// Details contient les détails supplémentaires
	Details map[string]interface{}
	// Timestamp est l'horodatage de l'erreur
	Timestamp int64
}

// commentedAuthToken représente un token d'authentification.
type commentedAuthToken struct {
	// Token est le token JWT
	Token string
	// ExpiresAt est la date d'expiration
	ExpiresAt int64
	// UserID est l'identifiant de l'utilisateur
	UserID string
	// Scope contient les permissions du token
	Scope []string
	// interne
	refreshed bool
}

// commentedServerConfig contient la configuration du serveur.
type commentedServerConfig struct {
	// Host est l'hôte du serveur
	Host string
	// Port est le port d'écoute
	Port int
	// Timeout est le délai d'expiration
	Timeout int
	// TLS indique si TLS est activé
	TLS bool
	// CertPath est le chemin du certificat TLS
	CertPath string
	// privé
	listener interface{}
}

// commentedLogConfig contient la configuration des logs.
type commentedLogConfig struct {
	// Level est le niveau de log
	Level string
	// Output est la destination des logs
	Output string
	// Format est le format de sortie
	Format string
	// MaxSize est la taille maximale d'un fichier de log en MB
	MaxSize int
	// MaxAge est la durée de rétention en jours
	MaxAge int
	// interne
	writer interface{}
}

// commentedPrivateStruct est une struct privée avec champs privés non documentés (OK).
type commentedPrivateStruct struct {
	field1 string
	field2 int
	field3 bool
}
