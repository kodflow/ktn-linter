package rules_const

// Configuration correctement groupée et documentée.
const (
	// MaxRetries nombre maximum de tentatives.
	MaxRetries = 3
	// DefaultValue valeur par défaut du système.
	DefaultValue = "test"
	// APIKey clé d'authentification API.
	APIKey = "secret"
	// MinTimeout timeout minimum en secondes.
	MinTimeout = 10
	// HTTPStatusCode code de statut HTTP par défaut.
	HTTPStatusCode = 200
	// HTMLParser nom du parser HTML utilisé.
	HTMLParser = "parser"
	// URLEndpoint point d'entrée de l'API.
	URLEndpoint = "/api/test"
)

// Configuration de la base de données (groupée).
const (
	// DatabaseHost adresse du serveur de base de données.
	DatabaseHost = "localhost"
	// DatabasePort port du serveur de base de données.
	DatabasePort = 5432
	// DatabaseName nom de la base de données.
	DatabaseName = "mydb"
)

// Configuration de l'application (types organisés).
const (
	// ConfigVersion version de la configuration.
	ConfigVersion int = 1
	// ConfigName nom de l'application.
	ConfigName string = "app"
	// ConfigEnabled indique si la configuration est active.
	ConfigEnabled bool = true
	// ConfigTimeout timeout de configuration en secondes.
	ConfigTimeout int = 30
)
