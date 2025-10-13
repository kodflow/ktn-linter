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

// ❌ CAS INCORRECT 1 : Bool non groupé (SEULE ERREUR : KTN-CONST-001)
// NOTE : Tout est parfait (commentaire + type) SAUF pas de ()
// ERREUR ATTENDUE : KTN-CONST-001 sur EnableFeatureX

// EnableFeatureX active la fonctionnalité X
const EnableFeatureX bool = true

// ❌ CAS INCORRECT 2 : String non groupé (SEULE ERREUR : KTN-CONST-001)
// NOTE : Tout est parfait (commentaire + type) SAUF pas de ()
// ERREUR ATTENDUE : KTN-CONST-001 sur ThemeAuto

// ThemeAuto est l'identifiant du thème automatique
const ThemeAuto string = "auto"

// ❌ CAS INCORRECT 3 : Int non groupé (SEULE ERREUR : KTN-CONST-001)
// NOTE : Tout est parfait (commentaire + type) SAUF pas de ()
// ERREUR ATTENDUE : KTN-CONST-001 sur MaxUserID

// MaxUserID définit l'ID utilisateur maximum
const MaxUserID uint32 = 4294967295

// ❌ CAS INCORRECT 4 : Int16 non groupés (SEULE ERREUR : KTN-CONST-001 x3)
// NOTE : Tout est parfait (commentaires + types) SAUF pas de ()
// ERREURS ATTENDUES : KTN-CONST-001 sur MaxQueueSize, DefaultBufferSize, minCacheSize

// MaxQueueSize définit la taille maximale de la queue
const MaxQueueSize int16 = 10000

// DefaultBufferSize est la taille du buffer par défaut
const DefaultBufferSize int16 = 4096

// minCacheSize est la taille minimale du cache
const minCacheSize int16 = 512

// ❌ CAS INCORRECT 5 : Float64 non groupés (SEULE ERREUR : KTN-CONST-001 x3)
// NOTE : Tout est parfait (commentaires + types) SAUF pas de ()
// ERREURS ATTENDUES : KTN-CONST-001 sur Pi, EulerNumber, goldenRatio

// Pi est une approximation de Pi en haute précision
const Pi float64 = 3.14159265358979323846

// EulerNumber est le nombre d'Euler (e)
const EulerNumber float64 = 2.71828182845904523536

// goldenRatio est le nombre d'or (phi)
const goldenRatio float64 = 1.618033988749894848204586

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

// ❌ CAS INCORRECT 1 : Pas de commentaire de groupe (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// ApplicationName est le nom de l'application
	ApplicationName string = "MyApp"
	// Version est la version actuelle
	Version string = "1.0.0"
	// defaultEncoding est l'encodage par défaut
	defaultEncoding string = "UTF-8"
)

// ❌ CAS INCORRECT 2 : Pas de commentaire de groupe avec int64 (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// MaxDiskSpace définit l'espace disque maximum en octets
	MaxDiskSpace int64 = 1099511627776
	// UnixEpoch représente le timestamp Unix de référence
	UnixEpoch int64 = 0
	// nanosPerSecond est le nombre de nanosecondes dans une seconde
	nanosPerSecond int64 = 1000000000
)

// ❌ CAS INCORRECT 3 : Pas de commentaire de groupe avec byte (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// NullByte représente l'octet null
	NullByte byte = 0x00
	// NewlineByte représente l'octet newline
	NewlineByte byte = 0x0A
	// tabByte représente l'octet tabulation
	tabByte byte = 0x09
)

// ❌ CAS INCORRECT 4 : Pas de commentaire de groupe avec complex64 (SEULE ERREUR : KTN-CONST-002)
// NOTE : Groupe OK, commentaires individuels OK, types OK, MAIS pas de commentaire de groupe
// ERREUR ATTENDUE : KTN-CONST-002 sur le groupe

const (
	// ImaginaryUnit64 est l'unité imaginaire en complex64
	ImaginaryUnit64 complex64 = 0 + 1i
	// ComplexZero64 est zéro en complex64
	ComplexZero64 complex64 = 0 + 0i
	// sampleComplex64 est un exemple de nombre complexe
	sampleComplex64 complex64 = 3.5 + 2.8i
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

// ❌ CAS INCORRECT 1 : Int8 sans commentaires individuels (SEULE ERREUR : KTN-CONST-003 x3)
// NOTE : Groupe OK, commentaire de groupe OK, types OK, MAIS pas de commentaires individuels
// ERREURS ATTENDUES : KTN-CONST-003 sur MinAge, MaxAge, defaultPriority
// Age configuration
// Ces constantes utilisent des entiers 8 bits pour les âges
const (
	MinAge          int8 = 18
	MaxAge          int8 = 120
	defaultPriority int8 = 5
)

// ❌ CAS INCORRECT 2 : Uint16 avec partiellement commentées (SEULE ERREUR : KTN-CONST-003 x2)
// NOTE : Groupe OK, commentaire groupe OK, types OK, HTTPPort commenté, MAIS HTTPSPort et customPort non commentés
// ERREURS ATTENDUES : KTN-CONST-003 sur HTTPSPort, customPort
// Port configuration
// Ces constantes définissent les ports réseau standards
const (
	// HTTPPort est le port HTTP standard
	HTTPPort   uint16 = 80
	HTTPSPort  uint16 = 443
	customPort uint16 = 3000
)

// ❌ CAS INCORRECT 3 : Float32 sans commentaires individuels (SEULE ERREUR : KTN-CONST-003 x3)
// NOTE : Groupe OK, commentaire groupe OK, types OK, MAIS pas de commentaires individuels
// ERREURS ATTENDUES : KTN-CONST-003 sur Pi32, DefaultRate, minThreshold
// Mathematical constants
// Ces constantes représentent des valeurs mathématiques en float32
const (
	Pi32         float32 = 3.14159265
	DefaultRate  float32 = 1.5
	minThreshold float32 = 0.01
)

// ❌ CAS INCORRECT 4 : Complex128 avec première constante non commentée (SEULE ERREUR : KTN-CONST-003 x1)
// NOTE : Groupe OK, commentaire groupe OK, types OK, MAIS ImaginaryUnit sans commentaire individuel
// ERREUR ATTENDUE : KTN-CONST-003 sur ImaginaryUnit
// Complex number constants
// Ces constantes représentent des nombres complexes en complex128
const (
	ImaginaryUnit complex128 = 0 + 1i
	// ComplexZero est zéro en complex128
	ComplexZero complex128 = 0 + 0i
	// eulerIdentityBase est la base de l'identité d'Euler
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

// ❌ CAS INCORRECT 1 : Int sans type explicite (SEULE ERREUR : KTN-CONST-004 x3)
// NOTE : Groupe OK, commentaire groupe OK, commentaires individuels OK, MAIS types manquants
// ERREURS ATTENDUES : KTN-CONST-004 sur MaxConnections, DefaultPort, maxRetries
// Connection limits
// Ces constantes définissent les limites de connexion
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries = 3
)

// ❌ CAS INCORRECT 2 : Int32 sans type explicite (SEULE ERREUR : KTN-CONST-004 x3)
// NOTE : Groupe OK, commentaire groupe OK, commentaires individuels OK, MAIS types manquants
// ERREURS ATTENDUES : KTN-CONST-004 sur MaxFileSize, DefaultTimeout, maxRequestsPerMinute
// File and timeout settings
// Ces constantes définissent les limites de fichiers et timeouts
const (
	// MaxFileSize définit la taille maximale d'un fichier en octets
	MaxFileSize = 104857600
	// DefaultTimeout est le timeout par défaut en millisecondes
	DefaultTimeout = 30000
	// maxRequestsPerMinute définit le nombre maximum de requêtes par minute
	maxRequestsPerMinute = 1000
)

// ❌ CAS INCORRECT 3 : Uint8 sans type explicite (SEULE ERREUR : KTN-CONST-004 x3)
// NOTE : Groupe OK, commentaire groupe OK, commentaires individuels OK, MAIS types manquants
// ERREURS ATTENDUES : KTN-CONST-004 sur MaxRetryAttempts, DefaultQuality, minCompressionLevel
// Quality settings
// Ces constantes définissent les paramètres de qualité
const (
	// MaxRetryAttempts définit le nombre maximum de tentatives
	MaxRetryAttempts = 10
	// DefaultQuality est la qualité par défaut (0-100)
	DefaultQuality = 85
	// minCompressionLevel est le niveau de compression minimum
	minCompressionLevel = 1
)

// ❌ CAS INCORRECT 4 : Une seule constante sans type dans un groupe presque parfait (SEULE ERREUR : KTN-CONST-004 x1)
// NOTE : Groupe OK, commentaires OK, 2 constantes avec types, MAIS MaxTransactionID sans type
// ERREUR ATTENDUE : KTN-CONST-004 sur MaxTransactionID uniquement
// Transaction limits
// Ces constantes définissent les limites de transactions
const (
	// MaxMemoryBytes définit la mémoire maximale en octets
	MaxMemoryBytes uint64 = 17179869184
	// MaxTransactionID est l'ID de transaction maximum
	MaxTransactionID = 18446744073709551615
	// defaultCacheExpiry est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiry uint64 = 3600000000000
)
