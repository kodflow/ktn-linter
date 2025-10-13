package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-001 : Variables groupées dans var ()
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

// Boolean configuration
// Ces variables configurent les fonctionnalités de l'application (mutables)
var (
	// EnableFeatureX active la fonctionnalité X
	EnableFeatureX bool = true
	// EnableDebug active le mode debug
	EnableDebug bool = false
	// isProduction indique si l'environnement est en production
	isProduction bool = true
)

// String configuration
// Ces variables configurent les thèmes de l'application
var (
	// ThemeAuto est l'identifiant du thème automatique
	ThemeAuto string = "auto"
	// ThemeCustom est l'identifiant du thème personnalisé
	ThemeCustom string = "custom"
)

// Integer configuration
// Ces variables configurent les limites entières (ajustables à runtime)
var (
	// MaxQueueSize est la taille maximale de la queue
	MaxQueueSize int16 = 10000
	// DefaultBufferSize est la taille par défaut du buffer
	DefaultBufferSize int16 = 4096
	// minCacheSize est la taille minimale du cache
	minCacheSize int16 = 512
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-002 : Groupe avec commentaire de groupe
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

// Application metadata
// Ces variables contiennent les métadonnées (peuvent être modifiées à runtime)
var (
	// ApplicationName est le nom de l'application
	ApplicationName string = "MyApp"
	// Version est la version actuelle de l'application
	Version string = "1.0.0"
	// defaultEncoding est l'encodage par défaut utilisé
	defaultEncoding string = "UTF-8"
)

// Disk and time values
// Ces variables utilisent int64 pour les grandes valeurs
var (
	// MaxDiskSpace est l'espace disque maximum en octets
	MaxDiskSpace int64 = 1099511627776
	// UnixEpoch représente le timestamp Unix epoch (intentionnellement 0)
	UnixEpoch int64 = 0
	// nanosPerSecond est le nombre de nanosecondes par seconde
	nanosPerSecond int64 = 1000000000
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-003 : Variable avec commentaire individuel
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
//    // Age limits
//    // Ces variables définissent les limites d'âge (configurables)
//    var (
//        // MinAge est l'âge minimum autorisé
//        MinAge int8 = 18
//        // MaxAge est l'âge maximum autorisé
//        MaxAge int8 = 120
//        // defaultPriority est la priorité par défaut
//        defaultPriority int8 = 5
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Age limits
// Ces variables définissent les limites d'âge (configurables)
var (
	// MinAge est l'âge minimum autorisé
	MinAge int8 = 18
	// MaxAge est l'âge maximum autorisé
	MaxAge int8 = 120
	// defaultPriority est la priorité par défaut
	defaultPriority int8 = 5
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-004 : Variable avec type explicite
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

// Integer configuration
// Ces variables configurent les limites entières (ajustables à runtime)
var (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections int = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort int = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries int = 3
)

// File and timeout limits
// Ces variables définissent les limites de fichier et timeout
var (
	// MaxFileSize est la taille maximale d'un fichier en octets
	MaxFileSize int32 = 104857600
	// DefaultTimeout est le timeout par défaut en millisecondes
	DefaultTimeout int32 = 30000
	// maxRequestsPerMinute limite le nombre de requêtes par minute
	maxRequestsPerMinute int32 = 1000
)

// Default tags
// Ces variables définissent les tags par défaut (modifiables)
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

// Initialized variables
// Ces variables sont initialisées via des fonctions à l'initialisation du package
var (
	// CurrentTime contient l'heure de démarrage de l'application
	CurrentTime string = getCurrentTime()
	// defaultLogger est l'instance de logger par défaut
	defaultLogger interface{} = createLogger()
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
// Ces variables sont des pointeurs vers des structures (partagées)
var (
	// GlobalContext est le contexte global de l'application
	GlobalContext *Context = &Context{}
	// defaultUser est l'utilisateur par défaut (anonyme)
	defaultUser *User = &User{Name: "anonymous"}
)

// Counter variables with explicit zero value
// Ces variables comptent les événements (zero value intentionnelle)
var (
	// RequestCount compte le nombre total de requêtes
	RequestCount int = 0
	// ErrorCount compte le nombre total d'erreurs
	ErrorCount int = 0
	// warningCount compte le nombre total d'avertissements
	warningCount int = 0
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-005 : Variable vs Constante (utiliser const quand approprié)
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
// ✅ CAS PARFAIT (utiliser var correctement) :
//
//    Variables qui PEUVENT changer à runtime :
//
//    // User and pool limits
//    // Ces variables définissent les limites (ajustables à runtime)
//    var (
//        // MaxUserID est l'ID utilisateur maximum (peut être augmenté)
//        MaxUserID uint = 4294967295
//        // DefaultPoolSize est la taille par défaut du pool (configurable)
//        DefaultPoolSize uint = 100
//    )
//
//    Valeurs immuables qui devraient être const, pas var :
//
//    // Mathematical constants
//    const (
//        // Pi est la valeur de pi (immuable)
//        Pi float64 = 3.14159265358979323846
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// User and pool limits
// Ces variables définissent les limites utilisateur et pool (ajustables)
var (
	// MaxUserID est l'ID utilisateur maximum (peut être modifié dynamiquement)
	MaxUserID uint = 4294967295
	// DefaultPoolSize est la taille par défaut du pool (configurable)
	DefaultPoolSize uint = 100
	// minWorkers est le nombre minimum de workers (peut varier selon charge)
	minWorkers uint = 4
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-006 : Une variable par ligne (pas de déclaration multiple)
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

// Network settings
// Ces variables configurent les paramètres réseau
var (
	// HostName est le nom d'hôte par défaut
	HostName string = "localhost"
	// Port est le port réseau par défaut
	Port int = 8080
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-007 : Channel avec buffer size explicite
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
//    // Channel variables
//    // Ces variables gèrent les messages inter-goroutines
//    var (
//        // MessageQueue canal de messages (buffer=100)
//        MessageQueue chan string = make(chan string, 100)
//        // ErrorQueue canal d'erreurs (buffer=50)
//        ErrorQueue chan error = make(chan error, 50)
//        // DoneSignal signale la fin (unbuffered intentionnel)
//        DoneSignal chan bool = make(chan bool)
//    )
//
// ════════════════════════════════════════════════════════════════════════════

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

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-008 : Nom avec MixedCaps (pas d'underscore)
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
//    // HTTP status codes
//    // Ces variables contiennent les codes HTTP standards
//    var (
//        // HTTPOK représente le code 200
//        HTTPOK int = 200
//        // HTTPNotFound représente le code 404
//        HTTPNotFound int = 404
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// HTTP status codes
// Ces variables représentent les codes de statut HTTP standards
var (
	// HTTPOK représente le code HTTP 200
	HTTPOK int = 200
	// HTTPNotFound représente le code HTTP 404
	HTTPNotFound int = 404
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-009 : Nom en MixedCaps (pas ALL_CAPS)
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
//        // DefaultTimeout est le timeout par défaut
//        DefaultTimeout int = 30
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Buffer configuration
// Ces variables configurent les buffers et timeouts
var (
	// MaxBufferSize est la taille maximale du buffer
	MaxBufferSize int = 1024
	// DefaultTimeoutSeconds est le timeout par défaut en secondes
	DefaultTimeoutSeconds int = 30
)

// ════════════════════════════════════════════════════════════════════════════
// CAS MIXTES : Cumul de bonnes pratiques
// ════════════════════════════════════════════════════════════════════════════

// Theme configuration - Toutes les variables du même thème regroupées
// Ces variables définissent les thèmes disponibles (tous configurables ensemble)
var (
	// ThemeLight est l'identifiant du thème clair
	ThemeLight string = "light"
	// ThemeDark est l'identifiant du thème sombre
	ThemeDark string = "dark"
	// ThemeHighContrast est l'identifiant du thème à haut contraste
	ThemeHighContrast string = "high-contrast"
	// ThemeSepia est l'identifiant du thème sépia
	ThemeSepia string = "sepia"
)

// ════════════════════════════════════════════════════════════════════════════
// Types factices pour les exemples
// ════════════════════════════════════════════════════════════════════════════
func getCurrentTime() string    { return "" }
func createLogger() interface{} { return nil }

type Context struct{}
type User struct{ Name string }
