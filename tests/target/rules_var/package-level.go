package rules_var

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : BOOL
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Regroupement dans var ()
//    ✓ Commentaire de groupe (2 lignes : titre + description)
//    ✓ Chaque variable a son commentaire individuel
//    ✓ Type bool explicite pour toutes
//    ✓ Naming MixedCaps (publiques et privées)
//    ✓ Variables mutables (peuvent changer à runtime)
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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : STRING
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type string explicite
//    ✓ Cohésion thématique (métadonnées ensemble, thèmes ensemble)
//    ✓ Commentaires expliquent pourquoi variables (mutables en production)
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

// Theme configuration
// Ces variables définissent les thèmes disponibles (configurables)
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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : INT
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int explicite
//    ✓ Mix de variables publiques et privées
//    ✓ Commentaires expliquent le rôle de chaque variable
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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : INT8
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int8 explicite (valeurs -128 à 127)
//    ✓ Approprié pour âges, priorités
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
// VARIABLES TYPE : INT16
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int16 explicite (valeurs -32768 à 32767)
//    ✓ Approprié pour tailles de queue, buffers
// ════════════════════════════════════════════════════════════════════════════

// Queue configuration
// Ces variables configurent les tailles de queue (ajustables)
var (
	// MaxQueueSize est la taille maximale de la queue
	MaxQueueSize int16 = 10000
	// DefaultBufferSize est la taille par défaut du buffer
	DefaultBufferSize int16 = 4096
	// minCacheSize est la taille minimale du cache
	minCacheSize int16 = 512
)

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : INT32
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int32 explicite
//    ✓ Commentaires incluent les unités (octets, millisecondes)
// ════════════════════════════════════════════════════════════════════════════

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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : INT64
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int64 explicite
//    ✓ Valeur 0 documentée comme "intentionnelle" (UnixEpoch)
// ════════════════════════════════════════════════════════════════════════════

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
// VARIABLES TYPE : UINT
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type uint explicite
//    ✓ Approprié pour valeurs garanties positives (IDs, compteurs)
// ════════════════════════════════════════════════════════════════════════════

// User and pool limits
// Ces variables définissent les limites utilisateur et pool
var (
	// MaxUserID est l'ID utilisateur maximum
	MaxUserID uint = 4294967295
	// DefaultPoolSize est la taille par défaut du pool
	DefaultPoolSize uint = 100
	// minWorkers est le nombre minimum de workers
	minWorkers uint = 4
)

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : SLICE
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type slice complet explicite : []string, []int
//    ✓ Initialisation claire avec valeurs
//    ✓ Commentaires expliquent le contenu
//
// ⚠️  IMPORTANT : Pour les slices, TOUJOURS spécifier le type complet
//                 []string = []string{...}, pas = []string{...}
// ════════════════════════════════════════════════════════════════════════════

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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : MAP
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type map complet explicite : map[string]string
//    ✓ Initialisation avec make ou littéral
//    ✓ Commentaires expliquent le contenu
//
// ⚠️  IMPORTANT : Pour les maps, TOUJOURS spécifier le type complet
//                 map[string]string = map[string]string{...}
// ════════════════════════════════════════════════════════════════════════════

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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : CHANNEL
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type channel explicite : chan string, chan error
//    ✓ Buffer size explicite dans le commentaire
//    ✓ "unbuffered" mentionné quand pertinent
//
// ⚠️  IMPORTANT : TOUJOURS préciser dans le commentaire :
//                 - (buffer=N) pour buffered channels
//                 - (unbuffered) pour channels synchrones
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
// VARIABLES TYPE : INT (COMPTEURS AVEC ZERO VALUE)
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type int explicite
//    ✓ Initialisation à 0 explicite (zero value intentionnelle)
//    ✓ Commentaire mentionne "zero value intentionnelle"
//
// 📝 NOTE : Pour les compteurs, toujours initialiser explicitement à 0
//           et mentionner dans le commentaire que c'est intentionnel
// ════════════════════════════════════════════════════════════════════════════

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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : HTTP STATUS CODES
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Naming avec initialismes corrects (HTTPOK, HTTPNotFound)
//    ✓ Type int explicite
//    ✓ Commentaires clairs
//
// 📝 NAMING : HTTPOK (pas Http_OK, HTTP_OK, ou HttpOk)
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
// VARIABLES TYPE : NETWORK SETTINGS
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Une variable par ligne (pas HostName, Port = ...)
//    ✓ Chaque variable a son propre commentaire
//    ✓ Types explicites
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
// VARIABLES TYPE : INITIALISÉES PAR FONCTION
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type explicite (string, interface{})
//    ✓ Commentaires expliquent la source d'initialisation
//    ✓ Même si initialisées par fonction, le type doit être explicite
// ════════════════════════════════════════════════════════════════════════════

// Initialized variables
// Ces variables sont initialisées via des fonctions à l'initialisation du package
var (
	// CurrentTime contient l'heure de démarrage de l'application
	CurrentTime string = getCurrentTime()
	// defaultLogger est l'instance de logger par défaut
	defaultLogger interface{} = createLogger()
)

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : STRUCT
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type struct complet explicite
//    ✓ Définition du type ET initialisation
//    ✓ Commentaire explique la structure
// ════════════════════════════════════════════════════════════════════════════

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

// ════════════════════════════════════════════════════════════════════════════
// VARIABLES TYPE : POINTER
// ════════════════════════════════════════════════════════════════════════════
// ✅ CE QUI EST CORRECT :
//    ✓ Type pointer explicite : *Context, *User
//    ✓ Initialisation claire avec &Type{}
//    ✓ Commentaires expliquent l'usage
// ════════════════════════════════════════════════════════════════════════════

// Pointer variables
// Ces variables sont des pointeurs vers des structures (partagées)
var (
	// GlobalContext est le contexte global de l'application
	GlobalContext *Context = &Context{}
	// defaultUser est l'utilisateur par défaut (anonyme)
	defaultUser *User = &User{Name: "anonymous"}
)

// ════════════════════════════════════════════════════════════════════════════
// 📚 RÉSUMÉ DES BONNES PRATIQUES POUR VARIABLES
// ════════════════════════════════════════════════════════════════════════════
//
// 1. REGROUPEMENT :
//    ✓ Toujours utiliser var () pour regrouper
//    ✓ Grouper les variables par thème/domaine fonctionnel
//    ✓ Ne jamais déclarer var X = ... individuellement
//
// 2. COMMENTAIRES :
//    ✓ Commentaire de groupe : 2 lignes (titre + description)
//    ✓ Commentaire individuel : 1 ligne par variable
//    ✓ Mentionner si la valeur est "mutable" ou pourquoi c'est une var
//
// 3. TYPES :
//    ✓ TOUJOURS spécifier le type explicitement
//    ✓ Même pour slices : []string = []string{...}
//    ✓ Même pour maps : map[K]V = map[K]V{...}
//    ✓ Même initialisées par fonction : var X Type = func()
//
// 4. CHANNELS :
//    ✓ TOUJOURS préciser buffer size dans commentaire
//    ✓ Exemple : // Queue canal (buffer=100)
//    ✓ Ou : // Done signal (unbuffered)
//
// 5. ZERO VALUES :
//    ✓ Toujours initialiser explicitement : int = 0
//    ✓ Mentionner "zero value intentionnelle" dans commentaire
//
// 6. CONST vs VAR :
//    ✓ Si la valeur ne change JAMAIS → utiliser const
//    ✓ var est pour les valeurs MUTABLES uniquement
//
// 7. NAMING :
//    ✓ MixedCaps : MaxConnections, defaultPort
//    ✓ Jamais underscore : max_connections ❌
//    ✓ Jamais ALL_CAPS : MAX_CONNECTIONS ❌
//
// 8. ORGANISATION :
//    ✓ Variables du même domaine ensemble
//    ✓ Ordre logique par type (simple → complexe)
//    ✓ Séparation visuelle avec commentaires de section
//
// ════════════════════════════════════════════════════════════════════════════

// Types factices pour les exemples
func getCurrentTime() string    { return "" }
func createLogger() interface{} { return nil }

type Context struct{}
type User struct{ Name string }
