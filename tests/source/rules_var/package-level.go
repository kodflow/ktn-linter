package rules_var

// ERREUR 1: Variables non groupées
var EnableFeatureX bool = true
var EnableDebug bool = false
var isProduction bool = true

// ERREUR 2: Pas de commentaire de groupe ET pas de commentaires individuels
var (
	ApplicationName string = "MyApp"
	Version         string = "1.0.0"
	defaultEncoding string = "UTF-8"
)

// ERREUR 3: Type non explicite (inféré)
// Ces variables n'ont pas de type explicite
var (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries = 3
)

// ERREUR 4: Commentaire de groupe mais pas de commentaires individuels
// Ces variables utilisent des entiers 8 bits (-128 à 127)
var (
	MinAge          int8 = 18
	MaxAge          int8 = 120
	defaultPriority int8 = 5
)

// ERREUR 5: Pas de regroupement + pas de commentaire
var MaxQueueSize int16 = 10000
var DefaultBufferSize int16 = 4096
var minCacheSize int16 = 512

// ERREUR 6: Type non explicite + pas de commentaire individuel
// Integer 32-bit variables
var (
	MaxFileSize          = 104857600
	DefaultTimeout       = 30000
	maxRequestsPerMinute = 1000
)

// ERREUR 7: Pas de commentaire du tout (ni groupe ni individuel)
var (
	MaxDiskSpace   int64 = 1099511627776
	UnixEpoch      int64 = 0
	nanosPerSecond int64 = 1000000000
)

// ERREUR 8: Pas de regroupement + type non explicite
var MaxUserID = 4294967295
var DefaultPoolSize = 100
var minWorkers = 4

// ERREUR 9: Variables mutables (slice) sans initialisation claire
// Slice variables
var (
	// DefaultTags est la liste des tags par défaut
	DefaultTags = []string{"production", "main"}
	// AllowedMethods liste les méthodes HTTP autorisées
	AllowedMethods = []string{"GET", "POST"}
	// errorCodes liste les codes d'erreur
	errorCodes = []int{400, 401, 403, 404, 500}
)

// ERREUR 10: Variables mutables (map) sans type explicite sur la valeur
// Configuration map
var (
	// ConfigDefaults contient les valeurs par défaut
	ConfigDefaults = map[string]string{
		"timeout": "30s",
		"retry":   "3",
	}
	// headerDefaults contient les en-têtes par défaut
	headerDefaults = map[string]string{
		"Content-Type": "application/json",
	}
)

// ERREUR 11: Naming avec underscore (mauvaise pratique)
// HTTP codes
var (
	// HTTP_OK représente le code 200
	HTTP_OK int = 200
	// NOT_FOUND représente le code 404
	NOT_FOUND int = 404
)

// ERREUR 12: Naming en ALL_CAPS (mauvaise pratique)
var MAX_BUFFER_SIZE int = 1024

// ERREUR 13: Multiple variables sur une ligne (difficile à documenter)
// Network settings
var (
	// HostName et Port sont les paramètres réseau
	HostName, Port = "localhost", 8080
)

// ERREUR 14: Variable qui devrait être une constante (valeur immuable)
// Pi value
var (
	// Pi représente la valeur de pi
	Pi float64 = 3.14159265358979323846
)

// ERREUR 15: Zero value non intentionnelle (pas claire)
// Counter variables
var (
	// RequestCount compte les requêtes
	RequestCount int = 0
	// ErrorCount compte les erreurs
	ErrorCount int
	// warningCount compte les avertissements
	warningCount = 0
)

// ERREUR 16: Variables avec fonction d'initialisation mais sans type
// Initialized from function
var (
	// CurrentTime est l'heure actuelle
	CurrentTime = getCurrentTime()
	// defaultLogger est le logger par défaut
	defaultLogger = createLogger()
)

// ERREUR 17: Channel sans buffer size explicite
// Channel variables
var (
	// MessageQueue est la file de messages
	MessageQueue chan string = make(chan string)
	// ErrorQueue est la file d'erreurs
	ErrorQueue chan error
	// doneSignal signale la fin
	doneSignal = make(chan bool)
)

// ERREUR 18: Struct sans type explicite dans l'initialisation
// Config struct
var (
	// DefaultConfig est la configuration par défaut
	DefaultConfig = struct {
		Timeout int
		Retries int
	}{
		Timeout: 30,
		Retries: 3,
	}
)

// ERREUR 19: Variables pointer sans type explicite
// Pointer variables
var (
	// GlobalContext est le contexte global
	GlobalContext = &Context{}
	// defaultUser est l'utilisateur par défaut
	defaultUser = &User{Name: "anonymous"}
)

// ERREUR 20: Mélange de var groupées et non groupées dans le même thème
// Configuration theme - Partie 1
var (
	// ThemeLight est le thème clair
	ThemeLight string = "light"
	// ThemeDark est le thème sombre
	ThemeDark string = "dark"
)

// Configuration theme - Partie 2 (devrait être dans le même groupe)
var ThemeAuto string = "auto"
var ThemeCustom string = "custom"

// ERREUR 21: Variable orpheline (aucune règle respectée)
var orphanVar = 42

// Types factices pour les exemples
func getCurrentTime() string { return "" }
func createLogger() interface{} { return nil }

type Context struct{}
type User struct{ Name string }
