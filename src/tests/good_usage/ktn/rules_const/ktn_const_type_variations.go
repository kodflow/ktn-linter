package rules_const

// ════════════════════════════════════════════════════════════════════════════
// ✅ Code conforme : Variations de types numériques pour règles CONST
// ════════════════════════════════════════════════════════════════════════════
// Ce fichier démontre le code correct pour tous les types numériques Go:
// int8, int32, int64, uint, uint8, uint16, uint64, float32,
// complex64, complex128, byte, rune, uintptr
// ════════════════════════════════════════════════════════════════════════════

// Integer types (signed)
// Constantes pour les types entiers signés
const (
	// MaxInt8Value est la valeur maximale pour int8
	MaxInt8Value int8 = 127
	// MinInt8 est la valeur minimale pour int8
	MinInt8 int8 = -128

	// MaxInt32Value est la valeur maximale pour int32
	MaxInt32Value int32 = 2147483647
	// MinInt32 est la valeur minimale pour int32
	MinInt32 int32 = -2147483648

	// MaxInt64Value est la valeur maximale pour int64
	MaxInt64Value int64 = 9223372036854775807
	// MinInt64 est la valeur minimale pour int64
	MinInt64 int64 = -9223372036854775808
)

// Integer types (unsigned)
// Constantes pour les types entiers non signés
const (
	// DefaultPortNumber est le numéro de port par défaut
	DefaultPortNumber uint = 8080

	// MaxBrightness est la luminosité maximale
	MaxBrightness uint8 = 255
	// DefaultUint8 est une valeur uint8 par défaut
	DefaultUint8 uint8 = 100

	// MaxPortNumber est le numéro de port maximum
	MaxPortNumber uint16 = 65535
	// DefaultUint16 est une valeur uint16 par défaut
	DefaultUint16 uint16 = 5000

	// MaxUint64Value est la valeur maximale pour uint64
	MaxUint64Value uint64 = 18446744073709551615
	// DefaultUint64 est une valeur uint64 par défaut
	DefaultUint64 uint64 = 1000000000
)

// Floating-point types
// Constantes pour les types à virgule flottante
const (
	// PiFloat32 est une approximation de Pi en float32
	PiFloat32 float32 = 3.14159265
	// EpsilonFloat32 est une petite valeur pour comparaison float32
	EpsilonFloat32 float32 = 1e-6
	// GoldenRatioFloat32 est le nombre d'or en float32
	GoldenRatioFloat32 float32 = 1.618033988
)

// Complex types
// Constantes pour les types complexes
const (
	// UnitImaginary64 est l'unité imaginaire en complex64
	UnitImaginary64 complex64 = 0 + 1i

	// UnitImaginary128 est l'unité imaginaire en complex128
	UnitImaginary128 complex128 = 0 + 1i
	// RealComplexNumber est un nombre complexe avec partie réelle
	RealComplexNumber complex128 = 3.0 + 4.0i
)

// Character types
// Constantes pour les caractères
const (
	// NullByte est le byte nul
	NullByte byte = 0x00
	// MaxByte est la valeur maximale pour byte
	MaxByte byte = 255
	// AlphaAByte est le byte pour 'A'
	AlphaAByte byte = 'A'
)

// Rune types
// Constantes pour les runes (caractères Unicode)
const (
	// SpaceRune est le caractère espace
	SpaceRune rune = ' '
	// NewlineRune est le caractère nouvelle ligne
	NewlineRune rune = '\n'
	// TabRune est le caractère tabulation
	TabRune rune = '\t'
)

// Pointer types
// Constantes pour uintptr
const (
	// NullPointer est une valeur de pointeur nulle
	NullPointer uintptr = 0
)

// Mixed numeric constants
// Constantes numériques mixtes pour cas d'usage réels
const (
	// HTTPPort est le port HTTP standard
	HTTPPort uint16 = 80
	// HTTPSPort est le port HTTPS standard
	HTTPSPort uint16 = 443
	// MaxConnections est le nombre maximum de connexions
	MaxConnections int32 = 10000
	// TimeoutSeconds est le timeout en secondes
	TimeoutSeconds int8 = 30
	// BufferSizeKB est la taille du buffer en KB
	BufferSizeKB uint16 = 4096
	// MaxFileSizeBytes est la taille maximale de fichier
	MaxFileSizeBytes uint64 = 10737418240
	// PrecisionThreshold est le seuil de précision
	PrecisionThreshold float32 = 0.0001
	// PhaseShift est le déphasage complexe
	PhaseShift complex64 = 1 + 0i
)
