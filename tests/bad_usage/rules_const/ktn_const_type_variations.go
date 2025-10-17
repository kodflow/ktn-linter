package rules_const

// ════════════════════════════════════════════════════════════════════════════
// KTN-CONST : Variations de types numériques pour toutes les règles CONST
// ════════════════════════════════════════════════════════════════════════════
// Ce fichier teste les règles CONST avec tous les types numériques Go:
// int8, int32, int64, uint, uint8, uint16, uint64, float32,
// complex64, complex128, byte, rune, uintptr
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT : int8 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// MaxInt8Value est la valeur maximale pour int8
const MaxInt8Value int8 = 127

// ❌ CAS INCORRECT : int32 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// MaxInt32Value est la valeur maximale pour int32
const MaxInt32Value int32 = 2147483647

// ❌ CAS INCORRECT : int64 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// MaxInt64Value est la valeur maximale pour int64
const MaxInt64Value int64 = 9223372036854775807

// ❌ CAS INCORRECT : uint non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// DefaultPortNumber est le numéro de port par défaut
const DefaultPortNumber uint = 8080

// ❌ CAS INCORRECT : uint8 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// MaxBrightness est la luminosité maximale
const MaxBrightness uint8 = 255

// ❌ CAS INCORRECT : uint16 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// MaxPortNumber est le numéro de port maximum
const MaxPortNumber uint16 = 65535

// ❌ CAS INCORRECT : uint64 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// MaxUint64Value est la valeur maximale pour uint64
const MaxUint64Value uint64 = 18446744073709551615

// ❌ CAS INCORRECT : float32 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// PiFloat32 est une approximation de Pi en float32
const PiFloat32 float32 = 3.14159265

// ❌ CAS INCORRECT : complex64 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// UnitImaginary64 est l'unité imaginaire en complex64
const UnitImaginary64 complex64 = 0 + 1i

// ❌ CAS INCORRECT : complex128 non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// UnitImaginary128 est l'unité imaginaire en complex128
const UnitImaginary128 complex128 = 0 + 1i

// ❌ CAS INCORRECT : byte non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// NullByte est le byte nul
const NullByte byte = 0x00

// ❌ CAS INCORRECT : rune non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// SpaceRune est le caractère espace
const SpaceRune rune = ' '

// ❌ CAS INCORRECT : uintptr non groupé
// ERREUR ATTENDUE : KTN-CONST-001

// NullPointer est une valeur de pointeur nulle
const NullPointer uintptr = 0

// ❌ CAS INCORRECT : Constantes numériques non groupées (types variés)
// ERREURS ATTENDUES : KTN-CONST-001 sur chacune

// MinInt8 est la valeur minimale pour int8
const MinInt8 int8 = -128

// MinInt32 est la valeur minimale pour int32
const MinInt32 int32 = -2147483648

// MinInt64 est la valeur minimale pour int64
const MinInt64 int64 = -9223372036854775808

// DefaultUint8 est une valeur uint8 par défaut
const DefaultUint8 uint8 = 100

// DefaultUint16 est une valeur uint16 par défaut
const DefaultUint16 uint16 = 5000

// DefaultUint64 est une valeur uint64 par défaut
const DefaultUint64 uint64 = 1000000000

// EpsilonFloat32 est une petite valeur pour comparaison float32
const EpsilonFloat32 float32 = 1e-6

// GoldenRatioFloat32 est le nombre d'or en float32
const GoldenRatioFloat32 float32 = 1.618033988

// RealComplexNumber est un nombre complexe avec partie réelle
const RealComplexNumber complex128 = 3.0 + 4.0i

// NewlineRune est le caractère nouvelle ligne
const NewlineRune rune = '\n'

// TabRune est le caractère tabulation
const TabRune rune = '\t'

// MaxByte est la valeur maximale pour byte
const MaxByte byte = 255

// AlphaAByte est le byte pour 'A'
const AlphaAByte byte = 'A'
