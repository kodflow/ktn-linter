package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR : Variations de types numériques pour toutes les règles VAR
// ════════════════════════════════════════════════════════════════════════════
// Ce fichier teste les règles VAR avec tous les types numériques Go:
// int8, int32, int64, uint, uint8, uint16, uint64, float32,
// complex64, complex128, byte, rune, uintptr
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT : var non groupées (KTN-VAR-007)
// ERREURS ATTENDUES : KTN-VAR-007 sur chaque variable

// counterInt8 est un compteur int8 non groupé
var counterInt8 int8 = 0

// limitInt32 est une limite int32 non groupée
var limitInt32 int32 = 1000

// timestampInt64 est un timestamp int64 non groupé
var timestampInt64 int64 = 1234567890

// portNumber est un numéro de port uint non groupé
var portNumber uint = 8080

// alphaValue est une valeur alpha uint8 non groupée
var alphaValue uint8 = 255

// maxConnections est le nombre de connexions uint16 non groupé
var maxConnections uint16 = 5000

// fileSize est une taille de fichier uint64 non groupée
var fileSize uint64 = 1024000

// temperature est une température float32 non groupée
var temperature float32 = 23.5

// impedance est une impédance complex64 non groupée
var impedance complex64 = 50 + 0i

// waveFunction est une fonction d'onde complex128 non groupée
var waveFunction complex128 = 1 + 1i

// nullChar est un caractère nul byte non groupé
var nullChar byte = 0

// unicodePoint est un point Unicode rune non groupé
var unicodePoint rune = 'A'

// pointerAddress est une adresse pointeur uintptr non groupée
var pointerAddress uintptr = 0

// ❌ CAS INCORRECT : Variables avec type explicite redondant (KTN-VAR-006)
// ERREURS : KTN-VAR-006 + KTN-VAR-007 (type inference + non groupé)

// explicitInt8 a un type explicite redondant
var explicitInt8 int8 = int8(10)

// explicitUint16 a un type explicite redondant
var explicitUint16 uint16 = uint16(1000)

// explicitFloat32 a un type explicite redondant
var explicitFloat32 float32 = float32(3.14)

// explicitComplex64 a un type explicite redondant
var explicitComplex64 complex64 = complex64(1 + 2i)

// ✅ CAS CORRECT : Constantes pour valeurs immuables (KTN-VAR-005)
// CONFORME : KTN-VAR-005 (utilise const pour valeurs jamais modifiées)

// Buffer configuration
const (
	// fixedMaxInt8 ne change jamais
	fixedMaxInt8 int8 = 127

	// fixedPortUint16 ne change jamais
	fixedPortUint16 uint16 = 443

	// fixedPiFloat32 ne change jamais
	fixedPiFloat32 float32 = 3.14159

	// fixedImaginaryUnit ne change jamais
	fixedImaginaryUnit complex128 = 0 + 1i

	// fixedByteValue ne change jamais
	fixedByteValue byte = 0xFF

	// fixedRuneValue ne change jamais
	fixedRuneValue rune = '\n'
)

// ❌ CAS INCORRECT : Noms ALL_CAPS avec types numériques (KTN-VAR-009)
// ERREURS : KTN-VAR-009 sur chaque variable

// Numeric configuration
var (
	// MAX_INT8_VALUE viole le naming (ALL_CAPS)
	MAX_INT8_VALUE int8 = 100

	// MAX_UINT16_VALUE viole le naming (ALL_CAPS)
	MAX_UINT16_VALUE uint16 = 60000

	// DEFAULT_FLOAT32 viole le naming (ALL_CAPS)
	DEFAULT_FLOAT32 float32 = 1.0

	// COMPLEX_UNIT viole le naming (ALL_CAPS)
	COMPLEX_UNIT complex64 = 1 + 0i

	// NULL_BYTE viole le naming (ALL_CAPS)
	NULL_BYTE byte = 0x00

	// SPACE_RUNE viole le naming (ALL_CAPS)
	SPACE_RUNE rune = ' '
)

// ❌ CAS INCORRECT : Variables sans commentaires (KTN-VAR-003)
// ERREURS : KTN-VAR-003 sur chaque variable

// Numeric values
var (
	// smallInt8 describes this variable.
	smallInt8 int8 = 1
	// largeInt64 describes this variable.
	largeInt64 int64 = 999999
	// defaultUint describes this variable.
	defaultUint uint = 100
	// preciseFloat32 describes this variable.
	preciseFloat32 float32 = 0.001
	// phaseComplex128 describes this variable.
	phaseComplex128 complex128 = 0.707 + 0.707i
	// controlByte describes this variable.
	controlByte byte = 0x1F
	// escapeRune describes this variable.
	escapeRune rune = '\\'
)

// ❌ CAS INCORRECT : Déclaration sans valeur initiale pour types numériques (KTN-VAR-004)
// NOTE: Acceptable pour var à l'intérieur de fonctions, mais au niveau package c'est suspect

// counters est un groupe de compteurs sans initialisation
var (
	// globalInt8Counter describes this variable.
	globalInt8Counter int8
	// globalInt32Counter describes this variable.
	globalInt32Counter int32
	// globalInt64Counter describes this variable.
	globalInt64Counter int64
	// globalUintCounter describes this variable.
	globalUintCounter uint
	// globalUint64Counter describes this variable.
	globalUint64Counter uint64
	// globalFloat32Sum describes this variable.
	globalFloat32Sum float32
	// globalComplex64Val describes this variable.
	globalComplex64Val complex64
)

// ❌ CAS INCORRECT : Multiple déclarations pour types similaires (KTN-VAR-007)
// Au lieu d'un seul bloc, plusieurs déclarations séparées

// firstInt8 est la première valeur int8
var firstInt8 int8 = 1

// secondInt8 est la deuxième valeur int8
var secondInt8 int8 = 2

// thirdInt8 est la troisième valeur int8
var thirdInt8 int8 = 3

// firstUint16 est la première valeur uint16
var firstUint16 uint16 = 100

// secondUint16 est la deuxième valeur uint16
var secondUint16 uint16 = 200

// firstFloat32 est la première valeur float32
var firstFloat32 float32 = 1.1

// secondFloat32 est la deuxième valeur float32
var secondFloat32 float32 = 2.2
