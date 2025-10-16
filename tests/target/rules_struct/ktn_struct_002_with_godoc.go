package rules_struct

// ✅ CONFORME KTN-STRUCT-002
// Toutes les structs ont un commentaire godoc

// GoodDocUserConfig contient la configuration utilisateur pour l'application.
type GoodDocUserConfig struct {
	// Host est l'hôte du serveur
	Host string
	// Port est le port d'écoute
	Port int
}

// GoodDocAPIClient représente un client pour interagir avec l'API REST.
type GoodDocAPIClient struct {
	// BaseURL est l'URL de base de l'API
	BaseURL string
	// Timeout est le délai d'expiration des requêtes en secondes
	Timeout int
}

// GoodDocDatabaseConnection représente une connexion active à la base de données.
type GoodDocDatabaseConnection struct {
	// Host est l'hôte de la base de données
	Host string
	// Port est le port de connexion
	Port int
	// Database est le nom de la base
	Database string
	// Username est le nom d'utilisateur pour l'authentification
	Username string
	// Password est le mot de passe
	Password string
}

// GoodDocHTTPRequestConfig contient tous les paramètres nécessaires pour effectuer une requête HTTP.
type GoodDocHTTPRequestConfig struct {
	// Method est la méthode HTTP (GET, POST, etc.)
	Method string
	// URL est l'URL complète de la requête
	URL string
	// Headers contient les en-têtes HTTP personnalisés
	Headers map[string]string
	// Body contient le corps de la requête
	Body []byte
}

// GoodDocCacheEntry représente une entrée unique dans le système de cache.
type GoodDocCacheEntry struct {
	// Key est l'identifiant unique de l'entrée
	Key string
	// Value est la valeur stockée dans le cache
	Value interface{}
	// ExpiresAt est le timestamp d'expiration en Unix time
	ExpiresAt int64
}

// GoodDocWorkerPoolConfig définit la configuration d'un pool de workers concurrents.
type GoodDocWorkerPoolConfig struct {
	// MinWorkers est le nombre minimum de workers à maintenir actifs
	MinWorkers int
	// MaxWorkers est le nombre maximum de workers autorisés
	MaxWorkers int
	// QueueSize est la taille maximale de la queue de tâches
	QueueSize int
}

// GoodDocErrorResponse représente une réponse d'erreur standardisée de l'API.
type GoodDocErrorResponse struct {
	// Code est le code d'erreur HTTP
	Code int
	// Message est un message d'erreur lisible par l'utilisateur
	Message string
	// Details contient des informations détaillées sur l'erreur
	Details map[string]interface{}
}

// GoodDocAuthToken représente un token d'authentification JWT valide.
type GoodDocAuthToken struct {
	// Token est le token JWT encodé
	Token string
	// ExpiresAt est la date d'expiration du token
	ExpiresAt int64
	// UserID est l'identifiant de l'utilisateur associé
	UserID string
}

// GoodDocServerConfig contient tous les paramètres de configuration du serveur.
type GoodDocServerConfig struct {
	// Host est l'adresse d'écoute du serveur
	Host string
	// Port est le port d'écoute
	Port int
	// Timeout est le délai d'expiration global en secondes
	Timeout int
}

// GoodDocLogConfig définit la configuration du système de logging.
type GoodDocLogConfig struct {
	// Level est le niveau de log (debug, info, warn, error)
	Level string
	// Output est la destination des logs (stdout, file, etc.)
	Output string
	// Format est le format de sortie (json, text)
	Format string
}

// goodDocPrivateConfig est une configuration privée interne au package.
type goodDocPrivateConfig struct {
	// apiKey est la clé API secrète
	apiKey string
	// secret est le secret partagé
	secret string
}
