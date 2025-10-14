package rules_const

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
	// MaxConnectionsC004Good définit le nombre maximum de connexions simultanées
	MaxConnectionsC004Good int = 1000
	// DefaultPortC004Good est le port par défaut de l'application
	DefaultPortC004Good int = 8080
	// maxRetriesC004Good définit le nombre maximum de tentatives
	maxRetriesC004Good int = 3
)

// Integer 32-bit constants
// Ces constantes utilisent des entiers 32 bits (-2147483648 à 2147483647)
const (
	// MaxFileSizeC004Good définit la taille maximale d'un fichier en octets
	MaxFileSizeC004Good int32 = 104857600 // 100 MB
	// DefaultTimeoutC004Good est le timeout par défaut en millisecondes
	DefaultTimeoutC004Good int32 = 30000
	// maxRequestsPerMinuteC004Good limite le nombre de requêtes par minute
	maxRequestsPerMinuteC004Good int32 = 1000
)

// Unsigned integer constants (uint)
// Ces constantes utilisent des entiers non signés (taille dépend de l'architecture)
const (
	// MaxUserIDC004Good est l'ID utilisateur maximum
	MaxUserIDC004Good uint = 4294967295
	// DefaultPoolSizeC004Good est la taille par défaut du pool
	DefaultPoolSizeC004Good uint = 100
	// minWorkersC004Good est le nombre minimum de workers
	minWorkersC004Good uint = 4
)

// Float32 constants
// Ces constantes utilisent des nombres à virgule flottante 32 bits (précision simple)
const (
	// Pi32C004Good est une approximation de Pi en float32
	Pi32C004Good float32 = 3.14159265
	// DefaultRateC004Good est le taux par défaut
	DefaultRateC004Good float32 = 1.5
	// minThresholdC004Good est le seuil minimum
	minThresholdC004Good float32 = 0.01
)

// Float64 constants
// Ces constantes utilisent des nombres à virgule flottante 64 bits (double précision)
const (
	// PiC004Good est une approximation de Pi en haute précision
	PiC004Good float64 = 3.14159265358979323846
	// EulerNumberC004Good est le nombre d'Euler (e)
	EulerNumberC004Good float64 = 2.71828182845904523536
	// goldenRatioC004Good est le nombre d'or (phi)
	goldenRatioC004Good float64 = 1.618033988749894848204586
)

// Byte constants
// Ces constantes représentent des octets individuels pour encodages et protocoles
const (
	// NullByteC004Good représente l'octet null
	NullByteC004Good byte = 0x00
	// NewlineByteC004Good représente le caractère newline
	NewlineByteC004Good byte = 0x0A
	// tabByteC004Good représente le caractère tabulation
	tabByteC004Good byte = 0x09
)

// Rune constants
// Ces constantes représentent des caractères Unicode (code points)
const (
	// SpaceRuneC004Good représente le caractère espace
	SpaceRuneC004Good rune = ' '
	// NewlineRuneC004Good représente le caractère retour à la ligne
	NewlineRuneC004Good rune = '\n'
	// heartEmojiC004Good représente l'emoji cœur
	heartEmojiC004Good rune = '❤'
)

// Complex64 constants
// Ces constantes représentent des nombres complexes en précision simple (float32 + float32)
const (
	// ImaginaryUnit64C004Good représente l'unité imaginaire i en complex64
	ImaginaryUnit64C004Good complex64 = 0 + 1i
	// ComplexZero64C004Good représente zéro en complex64
	ComplexZero64C004Good complex64 = 0 + 0i
	// sampleComplex64C004Good est un exemple de nombre complexe
	sampleComplex64C004Good complex64 = 3.5 + 2.8i
)

// Complex128 constants
// Ces constantes représentent des nombres complexes haute précision (float64 + float64)
const (
	// ImaginaryUnitC004Good représente l'unité imaginaire i
	ImaginaryUnitC004Good complex128 = 0 + 1i
	// ComplexZeroC004Good représente zéro en complex128
	ComplexZeroC004Good complex128 = 0 + 0i
	// eulerIdentityBaseC004Good est la base pour l'identité d'Euler
	eulerIdentityBaseC004Good complex128 = 2.71828182845904523536 + 0i
)

// Unsigned 32-bit constants
// Ces constantes utilisent des entiers non signés 32 bits (0 à 4294967295)
const (
	// MaxRecordCountC004Good définit le nombre maximum d'enregistrements
	MaxRecordCountC004Good uint32 = 1000000
	// DefaultChunkSizeC004Good est la taille par défaut d'un chunk en octets
	DefaultChunkSizeC004Good uint32 = 65536
	// minBatchSizeC004Good est la taille minimale d'un batch
	minBatchSizeC004Good uint32 = 100
)

// Unsigned 8-bit constants
// Ces constantes utilisent des entiers non signés 8 bits (0 à 255)
const (
	// MaxRetryAttemptsC004Good définit le nombre maximum de tentatives
	MaxRetryAttemptsC004Good uint8 = 10
	// DefaultQualityC004Good est la qualité par défaut (0-100)
	DefaultQualityC004Good uint8 = 85
	// minCompressionLevelC004Good est le niveau de compression minimum
	minCompressionLevelC004Good uint8 = 1
)

// Unsigned 64-bit constants
// Ces constantes utilisent des entiers non signés 64 bits (très grandes valeurs positives)
const (
	// MaxMemoryBytesC004Good définit la mémoire maximale en octets
	MaxMemoryBytesC004Good uint64 = 17179869184 // 16 GB
	// MaxTransactionIDC004Good est l'ID de transaction maximum
	MaxTransactionIDC004Good uint64 = 18446744073709551615
	// defaultCacheExpiryC004Good est le délai d'expiration du cache en nanosecondes
	defaultCacheExpiryC004Good uint64 = 3600000000000
)

// Theme configuration - Toutes les constantes du même thème regroupées
// Ces constantes définissent les thèmes disponibles dans l'interface
const (
	// ThemeLightC004Good est l'identifiant du thème clair
	ThemeLightC004Good string = "light"
	// ThemeDarkC004Good est l'identifiant du thème sombre
	ThemeDarkC004Good string = "dark"
	// ThemeHighContrastC004Good est l'identifiant du thème à haut contraste
	ThemeHighContrastC004Good string = "high-contrast"
	// ThemeSepiaC004Good est l'identifiant du thème sépia
	ThemeSepiaC004Good string = "sepia"
)
