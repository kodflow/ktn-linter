package rules_var

// ════════════════════════════════════════════════════════════════════════════
// ✅ Code conforme : Variations de types numériques pour règles VAR
// ════════════════════════════════════════════════════════════════════════════
// Ce fichier démontre le code correct pour tous les types numériques Go:
// int8, int32, int64, uint, uint8, uint16, uint64, float32,
// complex64, complex128, byte, rune, uintptr
// ════════════════════════════════════════════════════════════════════════════

// Integer counters and limits
// Compteurs et limites entiers avec types appropriés
var (
	// counterInt8 est un compteur de petites valeurs
	counterInt8 int8 = 0
	// limitInt32 est une limite moyenne
	limitInt32 int32 = 1000
	// timestampInt64 est un timestamp Unix
	timestampInt64 int64 = 1234567890
)

// Unsigned integer values
// Valeurs entières non signées pour ports, tailles, etc.
var (
	// portNumber est le numéro de port du serveur
	portNumber uint = 8080
	// alphaValue est la valeur de transparence alpha
	alphaValue uint8 = 255
	// maxConnections est le nombre maximum de connexions
	maxConnections uint16 = 5000
	// fileSize est la taille du fichier en bytes
	fileSize uint64 = 1024000
)

// Floating-point values
// Valeurs à virgule flottante pour mesures physiques
var (
	// temperature est la température en degrés Celsius
	temperature float32 = 23.5
)

// Complex numbers
// Nombres complexes pour calculs scientifiques
var (
	// impedance est l'impédance électrique
	impedance complex64 = 50 + 0i
	// waveFunction est la fonction d'onde quantique
	waveFunction complex128 = 1 + 1i
)

// Character types
// Types caractères pour manipulation de texte
var (
	// nullChar est le caractère nul
	nullChar byte = 0
	// unicodePoint est un point de code Unicode
	unicodePoint rune = 'A'
)

// Pointer addresses
// Adresses pointeur pour manipulation bas niveau
var (
	// pointerAddress est une adresse mémoire
	pointerAddress uintptr = 0
)

// ✅ Variables avec type inference correct (sans type explicite redondant)
var (
	// inferredInt8 utilise l'inférence de type
	inferredInt8 = int8(10)
	// inferredUint16 utilise l'inférence de type
	inferredUint16 = uint16(1000)
	// inferredFloat32 utilise l'inférence de type
	inferredFloat32 = float32(3.14)
	// inferredComplex64 utilise l'inférence de type
	inferredComplex64 = complex64(1 + 2i)
)

// ✅ Noms MixedCaps corrects (pas ALL_CAPS)
var (
	// MaxInt8Value utilise MixedCaps (correct)
	MaxInt8Value int8 = 100
	// MaxUint16Value utilise MixedCaps (correct)
	MaxUint16Value uint16 = 60000
	// DefaultFloat32 utilise MixedCaps (correct)
	DefaultFloat32 float32 = 1.0
	// ComplexUnit utilise MixedCaps (correct)
	ComplexUnit complex64 = 1 + 0i
	// NullByte utilise MixedCaps (correct)
	NullByte byte = 0x00
	// SpaceRune utilise MixedCaps (correct)
	SpaceRune rune = ' '
)

// ✅ Variables groupées correctement (pas de déclarations individuelles)
var (
	// firstInt8 est la première valeur int8
	firstInt8 int8 = 1
	// secondInt8 est la deuxième valeur int8
	secondInt8 int8 = 2
	// thirdInt8 est la troisième valeur int8
	thirdInt8 int8 = 3
)

var (
	// firstUint16 est la première valeur uint16
	firstUint16 uint16 = 100
	// secondUint16 est la deuxième valeur uint16
	secondUint16 uint16 = 200
)

var (
	// firstFloat32 est la première valeur float32
	firstFloat32 float32 = 1.1
	// secondFloat32 est la deuxième valeur float32
	secondFloat32 float32 = 2.2
)

// ✅ Variables avec valeurs initiales explicites
var (
	// globalInt8Counter est initialisé à zéro explicitement
	globalInt8Counter int8 = 0
	// globalInt32Counter est initialisé à zéro explicitement
	globalInt32Counter int32 = 0
	// globalInt64Counter est initialisé à zéro explicitement
	globalInt64Counter int64 = 0
	// globalUintCounter est initialisé à zéro explicitement
	globalUintCounter uint = 0
	// globalUint64Counter est initialisé à zéro explicitement
	globalUint64Counter uint64 = 0
	// globalFloat32Sum est initialisé à zéro explicitement
	globalFloat32Sum float32 = 0.0
	// globalComplex64Val est initialisé à zéro explicitement
	globalComplex64Val complex64 = 0 + 0i
)

// ✅ Toutes les variables avec commentaires complets
var (
	// smallInt8 est une petite valeur entière signée sur 8 bits
	smallInt8 int8 = 1
	// largeInt64 est une grande valeur entière signée sur 64 bits
	largeInt64 int64 = 999999
	// defaultUint est une valeur entière non signée par défaut
	defaultUint uint = 100
	// preciseFloat32 est une valeur flottante précise sur 32 bits
	preciseFloat32 float32 = 0.001
	// phaseComplex128 est un nombre complexe représentant une phase
	phaseComplex128 complex128 = 0.707 + 0.707i
	// controlByte est un byte de contrôle
	controlByte byte = 0x1F
	// escapeRune est la rune d'échappement backslash
	escapeRune rune = '\\'
)

// ✅ Types numériques dans configuration réseau
var (
	// httpPort est le port HTTP standard
	httpPort uint16 = 80
	// httpsPort est le port HTTPS standard
	httpsPort uint16 = 443
	// maxRetries est le nombre maximum de tentatives
	maxRetries int8 = 3
	// timeoutSeconds est le timeout en secondes
	timeoutSeconds int32 = 30
)

// ✅ Types numériques dans calculs scientifiques
var (
	// gravitationalConstant est la constante gravitationnelle
	gravitationalConstant float32 = 6.674e-11
	// avogadroNumber est le nombre d'Avogadro (approximation)
	avogadroNumber float32 = 6.022e23
	// quantumState est un état quantique complexe
	quantumState complex128 = 0.707 + 0.707i
)

// ✅ Types byte et rune pour manipulation de texte
var (
	// asciiLetters contient les bytes des lettres ASCII
	asciiStart byte = 'A'
	// unicodeSnowman est la rune du bonhomme de neige
	unicodeSnowman rune = '☃'
	// emojiRocket est la rune de la fusée
	emojiRocket rune = '🚀'
)

// ✅ Uintptr pour manipulation bas niveau
var (
	// pageSize est la taille d'une page mémoire
	pageSize uintptr = 4096
	// alignmentBoundary est la limite d'alignement mémoire
	alignmentBoundary uintptr = 8
)
