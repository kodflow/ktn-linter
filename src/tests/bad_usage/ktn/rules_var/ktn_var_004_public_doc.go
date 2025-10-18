package rules_var

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
// ERREURS : KTN-VAR-004 sur MaxConnectionsV004, DefaultPortV004, maxRetriesV004
// Ces variables n'ont pas de type explicite
var (
	// MaxConnectionsV004 définit le nombre maximum de connexions simultanées
	MaxConnectionsV004 = 1000
	// DefaultPortV004 est le port par défaut de l'application
	DefaultPortV004 = 8080
	// maxRetriesV004 définit le nombre maximum de tentatives
	maxRetriesV004 = 3
)

// ❌ CAS INCORRECT 2 : Int32 sans type explicite
// NOTE : Tout est parfait (groupe + commentaire groupe + commentaires individuels) SAUF types manquants
// ERREURS ATTENDUES : KTN-VAR-004 UNIQUEMENT sur MaxFileSizeV004, DefaultTimeoutV004, maxRequestsPerMinuteV004
// Integer 32-bit variables
// Ces variables utilisent des entiers 32 bits
var (
	// MaxFileSizeV004 définit la taille maximale d'un fichier en octets
	MaxFileSizeV004 = 104857600
	// DefaultTimeoutV004 est le timeout par défaut en millisecondes
	DefaultTimeoutV004 = 30000
	// maxRequestsPerMinuteV004 définit le nombre maximum de requêtes par minute
	maxRequestsPerMinuteV004 = 1000
)

// ❌ CAS INCORRECT 3 : Slice sans type explicite complet
// ERREURS : KTN-VAR-004 sur DefaultTagsV004, AllowedMethodsV004, errorCodesV004
// Slice variables
var (
	// DefaultTagsV004 est la liste des tags par défaut
	DefaultTagsV004 = []string{"production", "main"}
	// AllowedMethodsV004 liste les méthodes HTTP autorisées
	AllowedMethodsV004 = []string{"GET", "POST"}
	// errorCodesV004 liste les codes d'erreur
	errorCodesV004 = []int{400, 401, 403, 404, 500}
)

// ❌ CAS INCORRECT 4 : Map sans type explicite complet
// ERREURS : KTN-VAR-004 sur ConfigDefaultsV004, headerDefaultsV004
// Configuration map
var (
	// ConfigDefaultsV004 contient les valeurs par défaut
	ConfigDefaultsV004 = map[string]string{
		"timeout": "30s",
		"retry":   "3",
	}
	// headerDefaultsV004 contient les en-têtes par défaut
	headerDefaultsV004 = map[string]string{
		"Content-Type": "application/json",
	}
)

// ❌ CAS INCORRECT 5 : Variables avec fonction d'initialisation mais sans type
// ERREURS : KTN-VAR-004 sur CurrentTimeV004, defaultLoggerV004
// Initialized from function
var (
	// CurrentTimeV004 est l'heure actuelle
	CurrentTimeV004 = getCurrentTimeV004()
	// defaultLoggerV004 est le logger par défaut
	defaultLoggerV004 = createLoggerV004()
)

// ❌ CAS INCORRECT 6 : Struct anonyme sans type explicite
// ERREURS : KTN-VAR-004 sur DefaultConfigV004
// Config struct
var (
	// DefaultConfigV004 est la configuration par défaut
	DefaultConfigV004 = struct {
		Timeout int
		Retries int
	}{
		Timeout: 30,
		Retries: 3,
	}
)

// ❌ CAS INCORRECT 7 : Pointer sans type explicite
// ERREURS : KTN-VAR-004 sur GlobalContextV004, defaultUserV004
// Pointer variables
var (
	// GlobalContextV004 est le contexte global
	GlobalContextV004 = &ContextV004{}
	// defaultUserV004 est l'utilisateur par défaut
	defaultUserV004 = &UserV004{Name: "anonymous"}
)

// ❌ CAS INCORRECT 8 : Zero value non claire (mélange avec/sans initialisation)
// ERREURS : KTN-VAR-004 sur warningCountV004 (type manquant)
// Counter variables
var (
	// RequestCountV004 compte les requêtes
	RequestCountV004 int = 0
	// ErrorCountV004 compte les erreurs
	ErrorCountV004 int
	// warningCountV004 compte les avertissements
	warningCountV004 = 0
)

// Types factices pour les exemples
func getCurrentTimeV004() string    { return "" }
func createLoggerV004() interface{} { return nil }

// ContextV004 represents the struct.
// UserV004 represents the struct.
type ContextV004 struct{}
// UserV004 represents the struct.
type UserV004 struct{ Name string }
