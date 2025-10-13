package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-001 : Variables non groupées dans var ()
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les variables package-level doivent être regroupées dans un bloc var ()
//    au lieu d'être déclarées individuellement avec "var X = ..."
//
//    POURQUOI :
//    - Améliore la lisibilité et l'organisation du code
//    - Facilite la maintenance (variables liées regroupées)
//    - Rend les variables mutables explicites et visibles
//    - Standard Go universel pour variables package-level
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Boolean configuration
//    // Ces variables configurent les fonctionnalités (mutables)
//    var (
//        // EnableFeatureX active la fonctionnalité X
//        EnableFeatureX bool = true
//        // EnableDebug active le mode debug
//        EnableDebug bool = false
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Bool non groupés
// ERREURS : KTN-VAR-001 sur EnableFeatureX, EnableDebug, isProduction
var EnableFeatureX bool = true
var EnableDebug bool = false
var isProduction bool = true

// ❌ CAS INCORRECT 2 : String non groupés
// ERREURS : KTN-VAR-001 sur ThemeAuto, ThemeCustom
var ThemeAuto string = "auto"
var ThemeCustom string = "custom"

// ❌ CAS INCORRECT 3 : Int16 non groupés
// ERREURS : KTN-VAR-001 sur MaxQueueSize, DefaultBufferSize, minCacheSize
var MaxQueueSize int16 = 10000
var DefaultBufferSize int16 = 4096
var minCacheSize int16 = 512

// ❌ CAS INCORRECT 4 : Variables non groupées avec type manquant
// ERREURS : KTN-VAR-001 + KTN-VAR-004 sur MaxUserID, DefaultPoolSize, minWorkers
var MaxUserID = 4294967295
var DefaultPoolSize = 100
var minWorkers = 4

// ❌ CAS INCORRECT 5 : Variable orpheline (toutes les erreurs)
// ERREURS : KTN-VAR-001 + KTN-VAR-004
var orphanVar = 42

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-002 : Groupe sans commentaire de groupe
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Chaque bloc var () doit avoir un commentaire de groupe avant le bloc
//    pour expliquer le contexte global de ces variables mutables.
//
//    POURQUOI :
//    - Documente l'intention du regroupement
//    - Aide à comprendre pourquoi ces variables sont mutables
//    - Requis par les outils de documentation Go (godoc)
//    - Améliore la maintenabilité
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Application metadata
//    // Ces variables contiennent les métadonnées (mutables en production)
//    var (
//        // ApplicationName est le nom de l'application
//        ApplicationName string = "MyApp"
//        // Version est la version actuelle
//        Version string = "1.0.0"
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Pas de commentaire de groupe avec strings
// ERREURS : KTN-VAR-002 sur le groupe + KTN-VAR-003 sur chaque variable
var (
	ApplicationName string = "MyApp"
	Version         string = "1.0.0"
	defaultEncoding string = "UTF-8"
)

// ❌ CAS INCORRECT 2 : Pas de commentaire de groupe avec int64
// ERREURS : KTN-VAR-002 + KTN-VAR-003 sur chaque variable
var (
	MaxDiskSpace   int64 = 1099511627776
	UnixEpoch      int64 = 0
	nanosPerSecond int64 = 1000000000
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-003 : Variable sans commentaire individuel
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    CHAQUE variable (publique ET privée) doit avoir son propre commentaire
//    individuel qui explique son rôle spécifique. Le commentaire doit être
//    sur la ligne juste au-dessus de la variable.
//
//    POURQUOI :
//    - Documente précisément le rôle de CETTE variable
//    - Obligatoire pour les variables publiques (godoc)
//    - Recommandé aussi pour les privées (maintenabilité)
//    - Variables mutables nécessitent plus de documentation
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // HTTP configuration
//    // Ces variables configurent le serveur HTTP (mutables)
//    var (
//        // HTTPPort est le port HTTP à utiliser
//        HTTPPort uint16 = 80
//        // HTTPSPort est le port HTTPS à utiliser
//        HTTPSPort uint16 = 443
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Int8 avec commentaire de groupe mais pas individuels
// ERREURS : KTN-VAR-003 sur MinAge, MaxAge, defaultPriority
// Ces variables utilisent des entiers 8 bits (-128 à 127)
var (
	MinAge          int8 = 18
	MaxAge          int8 = 120
	defaultPriority int8 = 5
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-004 : Variable sans type explicite
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    TOUTES les variables doivent avoir un type explicite spécifié.
//    Ne jamais laisser le compilateur inférer le type, même si c'est évident.
//
//    POURQUOI :
//    - Élimine l'ambiguïté (int ? int32 ? int64 ?)
//    - Rend le contrat explicite (important pour variables mutables)
//    - Évite les surprises de conversion de types
//    - Facilite la relecture et la maintenance
//    - Plus critique que pour const car variables mutables
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Integer configuration
//    // Ces variables configurent les limites (mutables)
//    var (
//        // MaxConnections définit le nombre maximum de connexions
//        MaxConnections int = 1000
//        // DefaultPort est le port par défaut
//        DefaultPort int = 8080
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Int sans type explicite
// ERREURS : KTN-VAR-004 sur MaxConnections, DefaultPort, maxRetries
// Ces variables n'ont pas de type explicite
var (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries = 3
)

// ❌ CAS INCORRECT 2 : Int32 sans type explicite + pas de commentaires individuels
// ERREURS : KTN-VAR-003 + KTN-VAR-004 sur chaque variable
// Integer 32-bit variables
var (
	MaxFileSize          = 104857600
	DefaultTimeout       = 30000
	maxRequestsPerMinute = 1000
)

// ❌ CAS INCORRECT 3 : Slice sans type explicite complet
// ERREURS : KTN-VAR-004 sur DefaultTags, AllowedMethods, errorCodes
// Slice variables
var (
	// DefaultTags est la liste des tags par défaut
	DefaultTags = []string{"production", "main"}
	// AllowedMethods liste les méthodes HTTP autorisées
	AllowedMethods = []string{"GET", "POST"}
	// errorCodes liste les codes d'erreur
	errorCodes = []int{400, 401, 403, 404, 500}
)

// ❌ CAS INCORRECT 4 : Map sans type explicite complet
// ERREURS : KTN-VAR-004 sur ConfigDefaults, headerDefaults
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

// ❌ CAS INCORRECT 5 : Variables avec fonction d'initialisation mais sans type
// ERREURS : KTN-VAR-004 sur CurrentTime, defaultLogger
// Initialized from function
var (
	// CurrentTime est l'heure actuelle
	CurrentTime = getCurrentTime()
	// defaultLogger est le logger par défaut
	defaultLogger = createLogger()
)

// ❌ CAS INCORRECT 6 : Struct anonyme sans type explicite
// ERREURS : KTN-VAR-004 sur DefaultConfig
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

// ❌ CAS INCORRECT 7 : Pointer sans type explicite
// ERREURS : KTN-VAR-004 sur GlobalContext, defaultUser
// Pointer variables
var (
	// GlobalContext est le contexte global
	GlobalContext = &Context{}
	// defaultUser est l'utilisateur par défaut
	defaultUser = &User{Name: "anonymous"}
)

// ❌ CAS INCORRECT 8 : Zero value non claire (mélange avec/sans initialisation)
// ERREURS : KTN-VAR-004 sur warningCount (type manquant)
// Counter variables
var (
	// RequestCount compte les requêtes
	RequestCount int = 0
	// ErrorCount compte les erreurs
	ErrorCount int
	// warningCount compte les avertissements
	warningCount = 0
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-005 : Variable devrait être une constante
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Si une valeur ne change jamais (comme Pi, Version fixe, etc.),
//    elle devrait être déclarée comme const, pas var.
//
//    POURQUOI :
//    - const est thread-safe par nature (immuable)
//    - Indique clairement l'intention (ne changera jamais)
//    - Évite les modifications accidentelles
//    - Optimisations possibles par le compilateur
//
// ✅ CAS PARFAIT (utiliser const) :
//
//    // Mathematical constants
//    // Ces constantes représentent des valeurs mathématiques (immuables)
//    const (
//        // Pi est la valeur de pi
//        Pi float64 = 3.14159265358979323846
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Valeur mathématique immuable en var
// ERREURS : KTN-VAR-005 (devrait être const)
// Pi value
var (
	// Pi représente la valeur de pi
	Pi float64 = 3.14159265358979323846
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-006 : Multiple variables sur une ligne
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Plusieurs variables déclarées sur une ligne (HostName, Port = "localhost", 8080)
//    rendent impossible la documentation individuelle de chaque variable.
//
//    POURQUOI :
//    - Impossible de mettre un commentaire par variable
//    - Difficile à lire et à maintenir
//    - Contraire aux bonnes pratiques de documentation
//
// ✅ CAS PARFAIT (une variable par ligne) :
//
//    // Network settings
//    // Ces variables configurent la connexion réseau
//    var (
//        // HostName est le nom d'hôte
//        HostName string = "localhost"
//        // Port est le port réseau
//        Port int = 8080
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Multiple variables sur une ligne
// ERREURS : KTN-VAR-006 sur HostName, Port
// Network settings
var (
	// HostName et Port sont les paramètres réseau
	HostName, Port = "localhost", 8080
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-007 : Channel sans buffer size explicite
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les channels doivent avoir le buffer size explicite dans le commentaire
//    ou préciser "unbuffered" si intentionnel.
//
//    POURQUOI :
//    - Clarté sur la sémantique (synchrone vs asynchrone)
//    - Aide à comprendre les performances attendues
//    - Évite les deadlocks non intentionnels
//    - Important pour la concurrence
//
// ✅ CAS PARFAIT (buffer size explicite) :
//
//    // Message channels
//    // Ces variables gèrent les messages inter-goroutines
//    var (
//        // MessageQueue canal de messages (buffer=100)
//        MessageQueue chan string = make(chan string, 100)
//        // DoneSignal signale la fin (unbuffered intentionnel)
//        DoneSignal chan bool = make(chan bool)
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Channels sans buffer info dans commentaire
// ERREURS : KTN-VAR-007 sur MessageQueue, doneSignal
// ERREURS : KTN-VAR-004 sur doneSignal (type manquant)
// Channel variables
var (
	// MessageQueue est la file de messages
	MessageQueue chan string = make(chan string)
	// ErrorQueue est la file d'erreurs
	ErrorQueue chan error
	// doneSignal signale la fin
	doneSignal = make(chan bool)
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-008 : Nom avec underscore (utiliser MixedCaps)
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les noms de variables doivent utiliser MixedCaps, pas underscore.
//    Go style : HTTPPort, maxRetries (pas HTTP_PORT, max_retries)
//
//    POURQUOI :
//    - Convention Go standard (Effective Go)
//    - Cohérence avec la stdlib Go
//    - Facilite la lecture (style uniforme)
//
// ✅ CAS PARFAIT (MixedCaps) :
//
//    // HTTP codes
//    // Ces variables contiennent les codes HTTP standards
//    var (
//        // HTTPOK représente le code 200
//        HTTPOK int = 200
//        // NotFound représente le code 404
//        NotFound int = 404
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Noms avec underscore (mais PAS en ALL_CAPS)
// ERREURS : KTN-VAR-008 UNIQUEMENT sur max_size, buffer_Size, http_client
// Buffer settings
var (
	// max_size définit la taille maximale (snake_case)
	max_size int = 1024
	// buffer_Size définit la taille du buffer (mixte avec underscore)
	buffer_Size int = 512
	// http_client est le client HTTP (snake_case)
	http_client string = "default"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-009 : Nom en ALL_CAPS (utiliser MixedCaps)
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les noms en ALL_CAPS sont réservés aux constantes dans d'autres langages
//    (C, Java) mais Go utilise MixedCaps pour tout.
//
//    POURQUOI :
//    - Convention Go standard (pas ALL_CAPS)
//    - Évite confusion avec conventions d'autres langages
//    - MixedCaps est le style unifié Go
//
// ✅ CAS PARFAIT (MixedCaps) :
//
//    // Buffer configuration
//    // Cette variable configure la taille du buffer
//    var (
//        // MaxBufferSize est la taille maximale du buffer
//        MaxBufferSize int = 1024
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Noms en ALL_CAPS (sans underscore)
// ERREURS : KTN-VAR-009 UNIQUEMENT sur MAXSIZE, TIMEOUT, BUFFERSIZE
// Size configuration
var (
	// MAXSIZE est la taille maximale (ALL_CAPS sans underscore)
	MAXSIZE int = 1024
	// TIMEOUT est le timeout par défaut (ALL_CAPS sans underscore)
	TIMEOUT int = 30
	// BUFFERSIZE est la taille du buffer (ALL_CAPS sans underscore)
	BUFFERSIZE int = 512
)

// ════════════════════════════════════════════════════════════════════════════
// CAS MIXTES : Cumul de plusieurs erreurs
// ════════════════════════════════════════════════════════════════════════════

// Configuration theme - Partie 1 (groupe valide)
var (
	// ThemeLight est le thème clair
	ThemeLight string = "light"
	// ThemeDark est le thème sombre
	ThemeDark string = "dark"
)

// ❌ ERREUR : Mélange groupé/non-groupé sur le même thème
// Les variables ci-dessous devraient être dans le groupe au-dessus
// ERREURS : KTN-VAR-001 sur ThemeHighContrast, ThemeSepia
var ThemeHighContrast string = "high-contrast"
var ThemeSepia string = "sepia"

// ❌ ERREUR : Cumul VAR-008 + VAR-009 (underscore ET ALL_CAPS)
// ERREURS : KTN-VAR-008 + KTN-VAR-009 sur MAX_BUFFER_SIZE, DEFAULT_TIMEOUT
// Buffer configuration
var (
	// MAX_BUFFER_SIZE viole les deux règles (underscore + ALL_CAPS)
	MAX_BUFFER_SIZE int = 2048
	// DEFAULT_TIMEOUT viole les deux règles (underscore + ALL_CAPS)
	DEFAULT_TIMEOUT int = 60
)

// ════════════════════════════════════════════════════════════════════════════
// Types factices pour les exemples
// ════════════════════════════════════════════════════════════════════════════
func getCurrentTime() string    { return "" }
func createLogger() interface{} { return nil }

type Context struct{}
type User struct{ Name string }
