package rules_const

// ERREUR 1: Constantes sans regroupement (pas de const ())
const EnableFeatureX bool = true
const EnableDebug bool = false
const isProduction bool = true

// ERREUR 2: Pas de commentaire de groupe ET pas de commentaires individuels
const (
	ApplicationName string = "MyApp"
	Version         string = "1.0.0"
	defaultEncoding string = "UTF-8"
)

// ERREUR 3: Type non explicite (inféré)
// Ces constantes n'ont pas de type explicite
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries = 3
)

// ERREUR 4: Commentaire de groupe mais pas de commentaires individuels
// Ces constantes utilisent des entiers 8 bits (-128 à 127)
const (
	MinAge          int8 = 18
	MaxAge          int8 = 120
	defaultPriority int8 = 5
)

// ERREUR 5: Pas de regroupement + pas de commentaire
const MaxQueueSize int16 = 10000
const DefaultBufferSize int16 = 4096
const minCacheSize int16 = 512

// ERREUR 6: Type non explicite + pas de commentaire individuel
// Integer 32-bit constants
const (
	MaxFileSize          = 104857600
	DefaultTimeout       = 30000
	maxRequestsPerMinute = 1000
)

// ERREUR 7: Pas de commentaire du tout (ni groupe ni individuel)
const (
	MaxDiskSpace   int64 = 1099511627776
	UnixEpoch      int64 = 0
	nanosPerSecond int64 = 1000000000
)

// ERREUR 8: Pas de regroupement + type non explicite
const MaxUserID = 4294967295
const DefaultPoolSize = 100
const minWorkers = 4

// ERREUR 9: Commentaire de groupe OK mais types non explicites
// Unsigned 8-bit constants
const (
	// MaxRetryAttempts définit le nombre maximum de tentatives
	MaxRetryAttempts = 10
	// DefaultQuality est la qualité par défaut (0-100)
	DefaultQuality = 85
	// minCompressionLevel est le niveau de compression minimum
	minCompressionLevel = 1
)

// ERREUR 10: Mélange - certaines avec commentaire, d'autres non
// Unsigned 16-bit constants
const (
	// HTTPPort est le port HTTP standard
	HTTPPort   uint16 = 80
	HTTPSPort  uint16 = 443
	customPort uint16 = 3000
)

// ERREUR 11: Pas de regroupement + commentaires individuels mais pas de type
// MaxRecordCount définit le nombre maximum d'enregistrements
const MaxRecordCount = 1000000

// DefaultChunkSize est la taille par défaut d'un chunk en octets
const DefaultChunkSize = 65536
const minBatchSize = 100

// ERREUR 12: Tout correct SAUF type non explicite sur une seule variable
// Unsigned 64-bit constants
const (
	// MaxMemoryBytes définit la mémoire maximale en octets
	MaxMemoryBytes uint64 = 17179869184
	// MaxTransactionID est l'ID de transaction maximum
	MaxTransactionID = 18446744073709551615 // Type manquant !
	// defaultCacheExpiry est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiry uint64 = 3600000000000
)

// ERREUR 13: Pas de commentaire de groupe
const (
	NullByte    byte = 0x00
	NewlineByte byte = 0x0A
	tabByte     byte = 0x09
)

// ERREUR 14: Constantes individuelles sans type
const SpaceRune = ' '
const NewlineRune = '\n'
const heartEmoji = '❤'

// ERREUR 15: Groupe OK, commentaire de groupe OK, mais pas de commentaires individuels
// Float32 constants
const (
	Pi32         float32 = 3.14159265
	DefaultRate  float32 = 1.5
	minThreshold float32 = 0.01
)

// ERREUR 16: Pas de regroupement du tout
const Pi float64 = 3.14159265358979323846
const EulerNumber float64 = 2.71828182845904523536
const goldenRatio float64 = 1.618033988749894848204586

// ERREUR 17: Groupe sans commentaire de groupe ET types manquants
const (
	ImaginaryUnit64 = 0 + 1i
	ComplexZero64   = 0 + 0i
	sampleComplex64 = 3.5 + 2.8i
)

// ERREUR 18: Un seul commentaire pour plusieurs constantes (pas individuel)
const (
	// Ces constantes représentent des valeurs complexes
	ImaginaryUnit     complex128 = 0 + 1i
	ComplexZero       complex128 = 0 + 0i
	eulerIdentityBase complex128 = 2.71828182845904523536 + 0i
)

// ERREUR 19: Constante sans aucune forme de commentaire ni regroupement ni type
const orphanConst = 42

// ERREUR 20: Mélange de const groupées et non groupées dans le même thème
// Configuration theme - Partie 1
const (
	// ThemeLight est le thème clair
	ThemeLight string = "light"
	// ThemeDark est le thème sombre
	ThemeDark string = "dark"
)

// Configuration theme - Partie 2 (devrait être dans le même groupe)
const ThemeAuto string = "auto"
const ThemeCustom string = "custom"
