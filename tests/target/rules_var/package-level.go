package rules_var

// Boolean variables
// Ces variables représentent des valeurs booléennes pour la configuration
var (
	// EnableFeatureX active la fonctionnalité X
	EnableFeatureX bool = true
	// EnableDebug active le mode debug
	EnableDebug bool = false
	// isProduction indique si l'environnement est en production
	isProduction bool = true
)

// Application metadata
// Ces variables contiennent les métadonnées de l'application
var (
	// ApplicationName est le nom de l'application
	ApplicationName string = "MyApp"
	// Version est la version actuelle de l'application
	Version string = "1.0.0"
	// defaultEncoding est l'encodage par défaut utilisé
	defaultEncoding string = "UTF-8"
)

// Integer variables (int)
// Ces variables utilisent le type int pour les valeurs entières standards
var (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections int = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort int = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries int = 3
)

// Age limits
// Ces variables définissent les limites d'âge avec int8 (-128 à 127)
var (
	// MinAge est l'âge minimum autorisé
	MinAge int8 = 18
	// MaxAge est l'âge maximum autorisé
	MaxAge int8 = 120
	// defaultPriority est la priorité par défaut
	defaultPriority int8 = 5
)

// Queue configuration
// Ces variables configurent les tailles de queue avec int16
var (
	// MaxQueueSize est la taille maximale de la queue
	MaxQueueSize int16 = 10000
	// DefaultBufferSize est la taille par défaut du buffer
	DefaultBufferSize int16 = 4096
	// minCacheSize est la taille minimale du cache
	minCacheSize int16 = 512
)

// File and timeout limits
// Ces variables définissent les limites de fichier et timeout en int32
var (
	// MaxFileSize est la taille maximale d'un fichier en octets
	MaxFileSize int32 = 104857600
	// DefaultTimeout est le timeout par défaut en millisecondes
	DefaultTimeout int32 = 30000
	// maxRequestsPerMinute limite le nombre de requêtes par minute
	maxRequestsPerMinute int32 = 1000
)

// Disk and time constants
// Ces variables utilisent int64 pour les grandes valeurs
var (
	// MaxDiskSpace est l'espace disque maximum en octets
	MaxDiskSpace int64 = 1099511627776
	// UnixEpoch représente le timestamp Unix epoch (intentionnellement 0)
	UnixEpoch int64 = 0
	// nanosPerSecond est le nombre de nanosecondes par seconde
	nanosPerSecond int64 = 1000000000
)

// User and pool limits
// Ces variables définissent les limites utilisateur et pool en uint
var (
	// MaxUserID est l'ID utilisateur maximum
	MaxUserID uint = 4294967295
	// DefaultPoolSize est la taille par défaut du pool
	DefaultPoolSize uint = 100
	// minWorkers est le nombre minimum de workers
	minWorkers uint = 4
)

// Default tags
// Ces variables définissent les tags par défaut de l'application
var (
	// DefaultTags est la liste des tags par défaut appliqués
	DefaultTags []string = []string{"production", "main"}
	// AllowedMethods liste les méthodes HTTP autorisées
	AllowedMethods []string = []string{"GET", "POST"}
	// errorCodes liste les codes d'erreur HTTP standards
	errorCodes []int = []int{400, 401, 403, 404, 500}
)

// Configuration maps
// Ces variables contiennent les configurations sous forme de maps
var (
	// ConfigDefaults contient les valeurs de configuration par défaut
	ConfigDefaults map[string]string = map[string]string{
		"timeout": "30s",
		"retry":   "3",
	}
	// headerDefaults contient les en-têtes HTTP par défaut
	headerDefaults map[string]string = map[string]string{
		"Content-Type": "application/json",
	}
)

// HTTP status codes
// Ces variables représentent les codes de statut HTTP
var (
	// HTTPOK représente le code HTTP 200
	HTTPOK int = 200
	// HTTPNotFound représente le code HTTP 404
	HTTPNotFound int = 404
)

// Network settings
// Ces variables configurent les paramètres réseau
var (
	// HostName est le nom d'hôte par défaut
	HostName string = "localhost"
	// Port est le port réseau par défaut
	Port int = 8080
)

// Counter variables
// Ces variables comptent les événements (zero value intentionnelle)
var (
	// RequestCount compte le nombre total de requêtes
	RequestCount int = 0
	// ErrorCount compte le nombre total d'erreurs
	ErrorCount int = 0
	// warningCount compte le nombre total d'avertissements
	warningCount int = 0
)

// Initialized variables
// Ces variables sont initialisées via des fonctions
var (
	// CurrentTime contient l'heure de démarrage de l'application
	CurrentTime string = getCurrentTime()
	// defaultLogger est l'instance de logger par défaut
	defaultLogger interface{} = createLogger()
)

// Channel variables
// Ces variables sont des channels pour la communication inter-goroutines
var (
	// MessageQueue est le channel pour les messages (buffer=100)
	MessageQueue chan string = make(chan string, 100)
	// ErrorQueue est le channel pour les erreurs (buffer=50)
	ErrorQueue chan error = make(chan error, 50)
	// doneSignal signale la fin d'exécution (unbuffered intentionnel)
	doneSignal chan bool = make(chan bool)
)

// Config struct
// DefaultConfig contient la configuration par défaut de l'application
var (
	// DefaultConfig définit les valeurs de timeout et retries
	DefaultConfig struct {
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
// Ces variables sont des pointeurs vers des structures
var (
	// GlobalContext est le contexte global de l'application
	GlobalContext *Context = &Context{}
	// defaultUser est l'utilisateur par défaut (anonyme)
	defaultUser *User = &User{Name: "anonymous"}
)

// Theme configuration
// Ces variables définissent les thèmes disponibles
var (
	// ThemeLight est l'identifiant du thème clair
	ThemeLight string = "light"
	// ThemeDark est l'identifiant du thème sombre
	ThemeDark string = "dark"
	// ThemeAuto est l'identifiant du thème automatique
	ThemeAuto string = "auto"
	// ThemeCustom est l'identifiant du thème personnalisé
	ThemeCustom string = "custom"
)

// Types factices pour les exemples
func getCurrentTime() string { return "" }
func createLogger() interface{} { return nil }

type Context struct{}
type User struct{ Name string }
