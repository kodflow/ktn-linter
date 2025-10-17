package rules_const

// Configuration correctement groupée et documentée avec types explicites.
const (
	// MaxRetries nombre maximum de tentatives.
	MaxRetries int = 3
	// DefaultValue valeur par défaut du système.
	DefaultValue string = "test"
	// APIKey clé d'authentification API.
	APIKey string = "secret"
	// MinTimeout timeout minimum en secondes.
	MinTimeout int = 10
	// HTTPStatusCode code de statut HTTP par défaut.
	HTTPStatusCode int = 200
	// HTMLParser nom du parser HTML utilisé.
	HTMLParser string = "parser"
	// URLEndpoint point d'entrée de l'API.
	URLEndpoint string = "/api/test"
)

// Configuration de la base de données (groupée avec types explicites).
const (
	// DatabaseHost adresse du serveur de base de données.
	DatabaseHost string = "localhost"
	// DatabasePort port du serveur de base de données.
	DatabasePort int = 5432
	// DatabaseName nom de la base de données.
	DatabaseName string = "mydb"
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
