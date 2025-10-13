package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-001 : Constantes groupées dans const ()
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Les constantes package-level doivent être regroupées dans un bloc const ()
//    au lieu d'être déclarées individuellement avec "const X = ..."
//
//    POURQUOI :
//    - Améliore la lisibilité et l'organisation du code
//    - Facilite la maintenance (constantes liées regroupées)
//    - Rend le code plus compact et structuré
//    - Standard Go universel pour constantes package-level
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Boolean configuration
//    // Ces constantes représentent des valeurs booléennes
//    const (
//        // EnableFeatureX active la fonctionnalité X
//        EnableFeatureX bool = true
//        // EnableDebug active le mode debug
//        EnableDebug bool = false
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Boolean configuration
// Ces constantes configurent les fonctionnalités de l'application
const (
	// EnableFeatureX active la fonctionnalité X
	EnableFeatureX bool = true
	// EnableDebug active le mode debug
	EnableDebug bool = false
	// isProduction indique si l'environnement est en production
	isProduction bool = true
)

// String configuration
// Ces constantes définissent les thèmes de l'application
const (
	// ThemeAuto est l'identifiant du thème automatique
	ThemeAuto string = "auto"
	// ThemeCustom est l'identifiant du thème personnalisé
	ThemeCustom string = "custom"
)

// Integer configuration
// Ces constantes configurent les limites entières
const (
	// MaxQueueSize définit la taille maximale de la queue
	MaxQueueSize int16 = 10000
	// DefaultBufferSize est la taille du buffer par défaut
	DefaultBufferSize int16 = 4096
	// minCacheSize est la taille minimale du cache
	minCacheSize int16 = 512
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-002 : Groupe avec commentaire de groupe
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Chaque bloc const () doit avoir un commentaire de groupe avant le bloc
//    pour expliquer le contexte global de ces constantes.
//
//    POURQUOI :
//    - Documente l'intention du regroupement
//    - Aide à comprendre le rôle global des constantes
//    - Requis par les outils de documentation Go (godoc)
//    - Améliore la maintenabilité
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Application metadata
//    // Ces constantes contiennent les métadonnées de l'application
//    const (
//        // ApplicationName est le nom de l'application
//        ApplicationName string = "MyApp"
//        // Version est la version actuelle
//        Version string = "1.0.0"
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Application metadata
// Ces constantes contiennent les métadonnées de l'application
const (
	// ApplicationName est le nom de l'application
	ApplicationName string = "MyApp"
	// Version est la version actuelle de l'application
	Version string = "1.0.0"
	// defaultEncoding est l'encodage par défaut utilisé
	defaultEncoding string = "UTF-8"
)

// Integer 64-bit constants
// Ces constantes utilisent des entiers 64 bits (très grandes valeurs)
const (
	// MaxDiskSpace définit l'espace disque maximum en octets
	MaxDiskSpace int64 = 1099511627776 // 1 TB
	// UnixEpoch représente le timestamp Unix de référence
	UnixEpoch int64 = 0
	// nanosPerSecond est le nombre de nanosecondes dans une seconde
	nanosPerSecond int64 = 1000000000
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-003 : Constante avec commentaire individuel
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
//    - Facilite la compréhension du code
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Integer 8-bit constants
//    // Ces constantes utilisent des entiers 8 bits (-128 à 127)
//    const (
//        // MinAge est l'âge minimum requis
//        MinAge int8 = 18
//        // MaxAge est l'âge maximum accepté
//        MaxAge int8 = 120
//        // defaultPriority est la priorité par défaut
//        defaultPriority int8 = 5
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Integer 8-bit constants
// Ces constantes utilisent des entiers 8 bits (-128 à 127)
const (
	// MinAge est l'âge minimum requis
	MinAge int8 = 18
	// MaxAge est l'âge maximum accepté
	MaxAge int8 = 120
	// defaultPriority est la priorité par défaut
	defaultPriority int8 = 5
)

// Unsigned 16-bit constants
// Ces constantes utilisent des entiers non signés 16 bits (0 à 65535)
const (
	// HTTPPort est le port HTTP standard
	HTTPPort uint16 = 80
	// HTTPSPort est le port HTTPS standard
	HTTPSPort uint16 = 443
	// customPort est un port personnalisé
	customPort uint16 = 3000
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST-004 : Constante avec type explicite
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    TOUTES les constantes doivent avoir un type explicite spécifié.
//    Ne jamais laisser le compilateur inférer le type, même si c'est évident.
//
//    POURQUOI :
//    - Élimine l'ambiguïté (int ? int32 ? int64 ?)
//    - Rend le contrat explicite et clair
//    - Évite les surprises de conversion de types
//    - Facilite la relecture et la maintenance
//    - Documentation auto-générée plus précise
//
// ✅ CAS PARFAIT (pas d'erreur) :
//
//    // Integer constants (int)
//    // Ces constantes utilisent le type int
//    const (
//        // MaxConnections définit le nombre maximum de connexions
//        MaxConnections int = 1000
//        // DefaultPort est le port par défaut
//        DefaultPort int = 8080
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// Integer constants (int)
// Ces constantes utilisent le type int (taille dépend de l'architecture)
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections int = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort int = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries int = 3
)

// Integer 32-bit constants
// Ces constantes utilisent des entiers 32 bits (-2147483648 à 2147483647)
const (
	// MaxFileSize définit la taille maximale d'un fichier en octets
	MaxFileSize int32 = 104857600 // 100 MB
	// DefaultTimeout est le timeout par défaut en millisecondes
	DefaultTimeout int32 = 30000
	// maxRequestsPerMinute limite le nombre de requêtes par minute
	maxRequestsPerMinute int32 = 1000
)

// Unsigned integer constants (uint)
// Ces constantes utilisent des entiers non signés (taille dépend de l'architecture)
const (
	// MaxUserID est l'ID utilisateur maximum
	MaxUserID uint = 4294967295
	// DefaultPoolSize est la taille par défaut du pool
	DefaultPoolSize uint = 100
	// minWorkers est le nombre minimum de workers
	minWorkers uint = 4
)

// Float32 constants
// Ces constantes utilisent des nombres à virgule flottante 32 bits (précision simple)
const (
	// Pi32 est une approximation de Pi en float32
	Pi32 float32 = 3.14159265
	// DefaultRate est le taux par défaut
	DefaultRate float32 = 1.5
	// minThreshold est le seuil minimum
	minThreshold float32 = 0.01
)

// Float64 constants
// Ces constantes utilisent des nombres à virgule flottante 64 bits (double précision)
const (
	// Pi est une approximation de Pi en haute précision
	Pi float64 = 3.14159265358979323846
	// EulerNumber est le nombre d'Euler (e)
	EulerNumber float64 = 2.71828182845904523536
	// goldenRatio est le nombre d'or (phi)
	goldenRatio float64 = 1.618033988749894848204586
)

// Byte constants
// Ces constantes représentent des octets individuels pour encodages et protocoles
const (
	// NullByte représente l'octet null
	NullByte byte = 0x00
	// NewlineByte représente le caractère newline
	NewlineByte byte = 0x0A
	// tabByte représente le caractère tabulation
	tabByte byte = 0x09
)

// Rune constants
// Ces constantes représentent des caractères Unicode (code points)
const (
	// SpaceRune représente le caractère espace
	SpaceRune rune = ' '
	// NewlineRune représente le caractère retour à la ligne
	NewlineRune rune = '\n'
	// heartEmoji représente l'emoji cœur
	heartEmoji rune = '❤'
)

// Complex64 constants
// Ces constantes représentent des nombres complexes en précision simple (float32 + float32)
const (
	// ImaginaryUnit64 représente l'unité imaginaire i en complex64
	ImaginaryUnit64 complex64 = 0 + 1i
	// ComplexZero64 représente zéro en complex64
	ComplexZero64 complex64 = 0 + 0i
	// sampleComplex64 est un exemple de nombre complexe
	sampleComplex64 complex64 = 3.5 + 2.8i
)

// Complex128 constants
// Ces constantes représentent des nombres complexes haute précision (float64 + float64)
const (
	// ImaginaryUnit représente l'unité imaginaire i
	ImaginaryUnit complex128 = 0 + 1i
	// ComplexZero représente zéro en complex128
	ComplexZero complex128 = 0 + 0i
	// eulerIdentityBase est la base pour l'identité d'Euler
	eulerIdentityBase complex128 = 2.71828182845904523536 + 0i
)

// Unsigned 32-bit constants
// Ces constantes utilisent des entiers non signés 32 bits (0 à 4294967295)
const (
	// MaxRecordCount définit le nombre maximum d'enregistrements
	MaxRecordCount uint32 = 1000000
	// DefaultChunkSize est la taille par défaut d'un chunk en octets
	DefaultChunkSize uint32 = 65536
	// minBatchSize est la taille minimale d'un batch
	minBatchSize uint32 = 100
)

// Unsigned 8-bit constants
// Ces constantes utilisent des entiers non signés 8 bits (0 à 255)
const (
	// MaxRetryAttempts définit le nombre maximum de tentatives
	MaxRetryAttempts uint8 = 10
	// DefaultQuality est la qualité par défaut (0-100)
	DefaultQuality uint8 = 85
	// minCompressionLevel est le niveau de compression minimum
	minCompressionLevel uint8 = 1
)

// Unsigned 64-bit constants
// Ces constantes utilisent des entiers non signés 64 bits (très grandes valeurs positives)
const (
	// MaxMemoryBytes définit la mémoire maximale en octets
	MaxMemoryBytes uint64 = 17179869184 // 16 GB
	// MaxTransactionID est l'ID de transaction maximum
	MaxTransactionID uint64 = 18446744073709551615
	// defaultCacheExpiry est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiry uint64 = 3600000000000
)

// ════════════════════════════════════════════════════════════════════════════
// CAS MIXTES : Cumul de bonnes pratiques
// ════════════════════════════════════════════════════════════════════════════

// Theme configuration - Toutes les constantes du même thème regroupées
// Ces constantes définissent les thèmes disponibles dans l'interface
const (
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
// 📚 RÉSUMÉ DES BONNES PRATIQUES POUR CONSTANTES
// ════════════════════════════════════════════════════════════════════════════
//
// 1. REGROUPEMENT (KTN-CONST-001) :
//    ✓ Toujours utiliser const () pour regrouper
//    ✓ Grouper les constantes par thème/domaine fonctionnel
//    ✓ Ne jamais déclarer const X = ... individuellement
//
// 2. COMMENTAIRE DE GROUPE (KTN-CONST-002) :
//    ✓ Chaque bloc const () doit avoir un commentaire avant
//    ✓ Format : 2 lignes (titre + description)
//    ✓ Exemple : // HTTP ports / Ces constantes définissent...
//
// 3. COMMENTAIRE INDIVIDUEL (KTN-CONST-003) :
//    ✓ CHAQUE constante (publique ET privée) a son commentaire
//    ✓ Commentaire sur la ligne juste au-dessus
//    ✓ Format : // NomConstante description de son rôle
//    ✓ Exemple : // MaxRetries définit le nombre maximum de tentatives
//
// 4. TYPE EXPLICITE (KTN-CONST-004) :
//    ✓ TOUJOURS spécifier le type : bool, string, int, int8-64, uint, float, etc.
//    ✓ Ne jamais écrire : const X = 1
//    ✓ Toujours écrire : const X int = 1
//    ✓ Choisir le bon type selon la plage de valeurs
//
// 5. NAMING :
//    ✓ MixedCaps : MaxConnections, defaultPort
//    ✓ Jamais underscore : max_connections ❌
//    ✓ Jamais ALL_CAPS : MAX_CONNECTIONS ❌
//    ✓ Initialismes en majuscules : HTTPPort, URLMaxLength
//
// 6. DOCUMENTATION :
//    ✓ Mentionner les unités (octets, millisecondes, etc.)
//    ✓ Mentionner les plages valides si pertinent (0-100)
//    ✓ Expliquer le rôle, pas juste répéter le nom
//
// 7. ORGANISATION :
//    ✓ Constantes du même domaine ensemble
//    ✓ Ordre logique par type ou par fonctionnalité
//    ✓ Séparation visuelle avec commentaires de section
//
// ════════════════════════════════════════════════════════════════════════════
