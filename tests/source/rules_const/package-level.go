package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-001 : Constantes non groupées dans const ()
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les constantes doivent être regroupées dans un bloc const () au lieu
//    d'être déclarées individuellement avec "const X = ..."
//
//    POURQUOI :
//    - Améliore la lisibilité en regroupant les constantes liées
//    - Facilite la maintenance (une section = un thème)
//    - Évite la pollution du namespace package-level
//    - Standard Go universellement accepté
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Boolean configuration
//    // Ces constantes configurent les fonctionnalités
//    const (
//        // EnableFeatureX active la fonctionnalité X
//        EnableFeatureX bool = true
//        // EnableDebug active le mode debug
//        EnableDebug bool = false
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Bool non groupés
// ERREURS : KTN-CONST-001 sur EnableFeatureX, EnableDebug, isProduction
const EnableFeatureX bool = true
const EnableDebug bool = false
const isProduction bool = true

// ❌ CAS INCORRECT 2 : String non groupés
// ERREURS : KTN-CONST-001 sur ThemeAuto, ThemeCustom
const ThemeAuto string = "auto"
const ThemeCustom string = "custom"

// ❌ CAS INCORRECT 3 : Int non groupés avec type manquant
// ERREURS : KTN-CONST-001 + KTN-CONST-004 sur MaxUserID, DefaultPoolSize, minWorkers
const MaxUserID = 4294967295
const DefaultPoolSize = 100
const minWorkers = 4

// ❌ CAS INCORRECT 4 : Int16 non groupés
// ERREURS : KTN-CONST-001 sur MaxQueueSize, DefaultBufferSize, minCacheSize
const MaxQueueSize int16 = 10000
const DefaultBufferSize int16 = 4096
const minCacheSize int16 = 512

// ❌ CAS INCORRECT 5 : Uint32 non groupés avec type manquant
// ERREURS : KTN-CONST-001 + KTN-CONST-004 sur MaxRecordCount, DefaultChunkSize, minBatchSize
// MaxRecordCount définit le nombre maximum d'enregistrements
const MaxRecordCount = 1000000

// DefaultChunkSize est la taille par défaut d'un chunk en octets
const DefaultChunkSize = 65536
const minBatchSize = 100

// ❌ CAS INCORRECT 6 : Float64 non groupés
// ERREURS : KTN-CONST-001 sur Pi, EulerNumber, goldenRatio
const Pi float64 = 3.14159265358979323846
const EulerNumber float64 = 2.71828182845904523536
const goldenRatio float64 = 1.618033988749894848204586

// ❌ CAS INCORRECT 7 : Rune non groupés avec type manquant
// ERREURS : KTN-CONST-001 + KTN-CONST-004 sur SpaceRune, NewlineRune, heartEmoji
const SpaceRune = ' '
const NewlineRune = '\n'
const heartEmoji = '❤'

// ❌ CAS INCORRECT 8 : Constante orpheline (toutes les erreurs)
// ERREURS : KTN-CONST-001 + KTN-CONST-004
const orphanConst = 42

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-002 : Groupe sans commentaire de groupe
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Chaque bloc const () doit avoir un commentaire de groupe avant le bloc
//    pour expliquer le contexte global de ce groupe de constantes.
//
//    POURQUOI :
//    - Documente l'intention du regroupement (pourquoi ces constantes ensemble ?)
//    - Aide les développeurs à comprendre le contexte global
//    - Requis par les outils de documentation Go (godoc)
//    - Améliore la maintenabilité long terme
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Application metadata
//    // Ces constantes définissent les métadonnées de l'application
//    const (
//        // ApplicationName est le nom de l'application
//        ApplicationName string = "MyApp"
//        // Version est la version actuelle
//        Version string = "1.0.0"
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Pas de commentaire de groupe avec strings (+ pas de commentaires individuels)
// ERREURS : KTN-CONST-002 sur le groupe + KTN-CONST-003 sur chaque constante
const (
	ApplicationName string = "MyApp"
	Version         string = "1.0.0"
	defaultEncoding string = "UTF-8"
)

// ❌ CAS INCORRECT 2 : Pas de commentaire de groupe avec int64 (+ pas de commentaires individuels)
// ERREURS : KTN-CONST-002 sur le groupe + KTN-CONST-003 sur chaque constante
const (
	MaxDiskSpace   int64 = 1099511627776
	UnixEpoch      int64 = 0
	nanosPerSecond int64 = 1000000000
)

// ❌ CAS INCORRECT 3 : Pas de commentaire de groupe avec byte
// ERREURS : KTN-CONST-002 sur le groupe + KTN-CONST-003 sur chaque constante
const (
	NullByte    byte = 0x00
	NewlineByte byte = 0x0A
	tabByte     byte = 0x09
)

// ❌ CAS INCORRECT 4 : Pas de commentaire de groupe avec complex64 (+ types manquants)
// ERREURS : KTN-CONST-002 + KTN-CONST-003 sur chaque + KTN-CONST-004 sur chaque
const (
	ImaginaryUnit64 = 0 + 1i
	ComplexZero64   = 0 + 0i
	sampleComplex64 = 3.5 + 2.8i
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-003 : Constante sans commentaire individuel
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    CHAQUE constante (publique ET privée) doit avoir son propre commentaire
//    individuel qui explique son rôle spécifique. Le commentaire doit être
//    sur la ligne juste au-dessus de la constante.
//
//    POURQUOI :
//    - Documente précisément le rôle de CETTE constante
//    - Obligatoire pour les constantes publiques (godoc)
//    - Recommandé aussi pour les privées (maintenabilité)
//    - Facilite la compréhension sans avoir à lire le code
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // HTTP configuration
//    // Ces constantes définissent les ports HTTP standards
//    const (
//        // HTTPPort est le port HTTP standard
//        HTTPPort uint16 = 80
//        // HTTPSPort est le port HTTPS standard
//        HTTPSPort uint16 = 443
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Int8 avec commentaire de groupe mais pas de commentaires individuels
// ERREURS : KTN-CONST-003 sur MinAge, MaxAge, defaultPriority
// Ces constantes utilisent des entiers 8 bits (-128 à 127)
const (
	MinAge          int8 = 18
	MaxAge          int8 = 120
	defaultPriority int8 = 5
)

// ❌ CAS INCORRECT 2 : Uint16 avec commentaire groupe + partiellement commentées
// ERREURS : KTN-CONST-003 sur HTTPSPort, customPort (HTTPPort est OK)
// Unsigned 16-bit constants
const (
	// HTTPPort est le port HTTP standard
	HTTPPort   uint16 = 80
	HTTPSPort  uint16 = 443
	customPort uint16 = 3000
)

// ❌ CAS INCORRECT 3 : Float32 avec commentaire groupe mais pas individuels
// ERREURS : KTN-CONST-003 sur Pi32, DefaultRate, minThreshold
// Float32 constants
const (
	Pi32         float32 = 3.14159265
	DefaultRate  float32 = 1.5
	minThreshold float32 = 0.01
)

// ❌ CAS INCORRECT 4 : Complex128 avec un seul commentaire partagé
// ERREURS : KTN-CONST-003 sur ComplexZero, eulerIdentityBase
const (
	// Ces constantes représentent des valeurs complexes
	ImaginaryUnit     complex128 = 0 + 1i
	ComplexZero       complex128 = 0 + 0i
	eulerIdentityBase complex128 = 2.71828182845904523536 + 0i
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-004 : Constante sans type explicite
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    TOUTES les constantes doivent avoir un type explicite spécifié.
//    Ne jamais laisser le compilateur inférer le type, même si c'est évident.
//
//    POURQUOI :
//    - Élimine l'ambiguïté (int ? int32 ? int64 ?)
//    - Rend le contrat explicite (importante pour APIs)
//    - Évite les surprises de conversion de types
//    - Facilite la relecture et la maintenance
//    - Standard pour code production
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Integer constants
//    // Ces constantes utilisent le type int explicite
//    const (
//        // MaxConnections définit le nombre maximum de connexions
//        MaxConnections int = 1000
//        // DefaultPort est le port par défaut
//        DefaultPort int = 8080
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Int sans type explicite
// ERREURS : KTN-CONST-004 sur MaxConnections, DefaultPort, maxRetries
// Ces constantes n'ont pas de type explicite
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries = 3
)

// ❌ CAS INCORRECT 2 : Int32 sans type explicite (avec commentaires individuels manquants)
// ERREURS : KTN-CONST-003 + KTN-CONST-004 sur chaque constante
// Integer 32-bit constants
const (
	MaxFileSize          = 104857600
	DefaultTimeout       = 30000
	maxRequestsPerMinute = 1000
)

// ❌ CAS INCORRECT 3 : Uint8 sans type explicite
// ERREURS : KTN-CONST-004 sur MaxRetryAttempts, DefaultQuality, minCompressionLevel
// Unsigned 8-bit constants
const (
	// MaxRetryAttempts définit le nombre maximum de tentatives
	MaxRetryAttempts = 10
	// DefaultQuality est la qualité par défaut (0-100)
	DefaultQuality = 85
	// minCompressionLevel est le niveau de compression minimum
	minCompressionLevel = 1
)

// ❌ CAS INCORRECT 4 : Uint64 avec une seule constante sans type dans un groupe presque parfait
// ERREURS : KTN-CONST-004 sur MaxTransactionID uniquement
// Unsigned 64-bit constants
const (
	// MaxMemoryBytes définit la mémoire maximale en octets
	MaxMemoryBytes uint64 = 17179869184
	// MaxTransactionID est l'ID de transaction maximum
	MaxTransactionID = 18446744073709551615 // Type manquant !
	// defaultCacheExpiry est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiry uint64 = 3600000000000
)

// ════════════════════════════════════════════════════════════════════════════
// CAS MIXTES : Cumul de plusieurs erreurs
// ════════════════════════════════════════════════════════════════════════════

// Configuration theme - Partie 1 (groupe valide)
const (
	// ThemeLight est le thème clair
	ThemeLight string = "light"
	// ThemeDark est le thème sombre
	ThemeDark string = "dark"
)

// ❌ ERREUR : Mélange groupé/non-groupé sur le même thème
// Les constantes ci-dessous devraient être dans le groupe au-dessus
// ERREURS : KTN-CONST-001 sur ThemeHighContrast, ThemeSepia
const ThemeHighContrast string = "high-contrast"
const ThemeSepia string = "sepia"
