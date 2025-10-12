package rules_const

// Boolean constants
// Ces constantes représentent des valeurs booléennes pour la configuration
const (
	// EnableFeatureX active la fonctionnalité X
	EnableFeatureX bool = true
	// EnableDebug active le mode debug
	EnableDebug bool = false
	// isProduction indique si l'environnement est en production
	isProduction bool = true
)

// String constants
// Ces constantes définissent des valeurs textuelles
const (
	// ApplicationName est le nom de l'application
	ApplicationName string = "MyApp"
	// Version est la version actuelle de l'application
	Version string = "1.0.0"
	// defaultEncoding est l'encodage par défaut
	defaultEncoding string = "UTF-8"
)

// Integer constants (int)
// Ces constantes utilisent le type int (taille dépend de l'architecture: 32 ou 64 bits)
const (
	// MaxConnections définit le nombre maximum de connexions simultanées
	MaxConnections int = 1000
	// DefaultPort est le port par défaut de l'application
	DefaultPort int = 8080
	// maxRetries définit le nombre maximum de tentatives
	maxRetries int = 3
)

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

// Integer 16-bit constants
// Ces constantes utilisent des entiers 16 bits (-32768 à 32767)
const (
	// MaxQueueSize définit la taille maximale de la queue
	MaxQueueSize int16 = 10000
	// DefaultBufferSize est la taille du buffer par défaut
	DefaultBufferSize int16 = 4096
	// minCacheSize est la taille minimale du cache
	minCacheSize int16 = 512
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

// Byte constants (alias de uint8)
// Ces constantes représentent des octets individuels (0 à 255)
const (
	// NullByte représente l'octet null
	NullByte byte = 0x00
	// NewlineByte représente le caractère newline
	NewlineByte byte = 0x0A
	// tabByte représente le caractère tabulation
	tabByte byte = 0x09
)

// Rune constants (alias de int32)
// Ces constantes représentent des caractères Unicode
const (
	// SpaceRune représente le caractère espace
	SpaceRune rune = ' '
	// NewlineRune représente le caractère retour à la ligne
	NewlineRune rune = '\n'
	// heartEmoji représente l'emoji cœur
	heartEmoji rune = '❤'
)

// Float32 constants
// Ces constantes utilisent des nombres à virgule flottante 32 bits
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
	// goldenRatio est le nombre d'or
	goldenRatio float64 = 1.618033988749894848204586
)

// Complex64 constants
// Ces constantes représentent des nombres complexes (float32 + float32)
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
