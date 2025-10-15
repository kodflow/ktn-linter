package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-004 : Variable avec type explicite
// ════════════════════════════════════════════════════════════════════════════

// Integer configuration
// Ces variables configurent les limites entières (ajustables à runtime)
var (
	// MaxConnectionsV004Good définit le nombre maximum de connexions simultanées
	MaxConnectionsV004Good int = 1000
	// DefaultPortV004Good est le port par défaut de l'application
	DefaultPortV004Good int = 8080
	// maxRetriesV004Good définit le nombre maximum de tentatives
	maxRetriesV004Good int = 3
)

// File and timeout limits
// Ces variables définissent les limites de fichier et timeout
var (
	// MaxFileSizeV004Good est la taille maximale d'un fichier en octets
	MaxFileSizeV004Good int32 = 104857600
	// DefaultTimeoutV004Good est le timeout par défaut en millisecondes
	DefaultTimeoutV004Good int32 = 30000
	// maxRequestsPerMinuteV004Good limite le nombre de requêtes par minute
	maxRequestsPerMinuteV004Good int32 = 1000
)

// Default tags
// Ces variables définissent les tags par défaut (modifiables)
var (
	// DefaultTagsV004Good est la liste des tags par défaut appliqués
	DefaultTagsV004Good []string = []string{"production", "main"}
	// AllowedMethodsV004Good liste les méthodes HTTP autorisées
	AllowedMethodsV004Good []string = []string{"GET", "POST"}
	// errorCodesV004Good liste les codes d'erreur HTTP standards
	errorCodesV004Good []int = []int{400, 401, 403, 404, 500}
)

// Configuration maps
// Ces variables contiennent les configurations sous forme de maps
var (
	// ConfigDefaultsV004Good contient les valeurs de configuration par défaut
	ConfigDefaultsV004Good map[string]string = map[string]string{
		"timeout": "30s",
		"retry":   "3",
	}
	// headerDefaultsV004Good contient les en-têtes HTTP par défaut
	headerDefaultsV004Good map[string]string = map[string]string{
		"Content-Type": "application/json",
	}
)

// Initialized variables
// Ces variables sont initialisées via des fonctions à l'initialisation du package
var (
	// CurrentTimeV004Good contient l'heure de démarrage de l'application
	CurrentTimeV004Good string = getCurrentTimeV004Good()
	// defaultLoggerV004Good est l'instance de logger par défaut
	defaultLoggerV004Good interface{} = createLoggerV004Good()
)

// Config struct
// DefaultConfigV004Good contient la configuration par défaut de l'application
var (
	// DefaultConfigV004Good définit les valeurs de timeout et retries
	DefaultConfigV004Good struct {
		Timeout int
		Retries int
	} = struct {
		Timeout int
		Retries int
	}{
		Timeout: 30,
		Retries: 3,
	}
)

// Pointer variables
// Ces variables sont des pointeurs vers des structures (partagées)
var (
	// GlobalContextV004Good est le contexte global de l'application
	GlobalContextV004Good *ContextV004GoodData = &ContextV004GoodData{}
	// defaultUserV004Good est l'utilisateur par défaut (anonyme)
	defaultUserV004Good *UserV004GoodData = &UserV004GoodData{Name: "anonymous"}
)

// Counter variables with explicit zero value
// Ces variables comptent les événements (zero value intentionnelle)
var (
	// RequestCountV004Good compte le nombre total de requêtes
	RequestCountV004Good int = 0
	// ErrorCountV004Good compte le nombre total d'erreurs
	ErrorCountV004Good int = 0
	// warningCountV004Good compte le nombre total d'avertissements
	warningCountV004Good int = 0
)

// getCurrentTimeV004Good retourne l'heure actuelle en string.
//
// Returns:
//   - string: l'heure actuelle sous forme de string
func getCurrentTimeV004Good() string { return "" }

// createLoggerV004Good crée une instance de logger.
//
// Returns:
//   - interface{}: l'instance du logger créée
func createLoggerV004Good() interface{} { return nil }

// updateConfigV004Good modifie les configurations à runtime.
func updateConfigV004Good() {
	MaxConnectionsV004Good = 2000
	DefaultPortV004Good = 9090
	maxRetriesV004Good = 5
	MaxFileSizeV004Good = 209715200
	DefaultTimeoutV004Good = 60000
	maxRequestsPerMinuteV004Good = 2000
	RequestCountV004Good = 100
	ErrorCountV004Good = 10
	warningCountV004Good = 5
}

type ContextV004GoodData struct{}
type UserV004GoodData struct{ Name string }
