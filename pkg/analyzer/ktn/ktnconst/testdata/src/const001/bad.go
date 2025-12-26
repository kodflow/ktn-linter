// Package const001 contains test cases for KTN rules.
package const001

// Bad: All constants WITHOUT explicit types (violates KTN-CONST-001)
// Names use CamelCase to avoid CONST-003 errors
const (
	// === Basic types without explicit type ===

	// BadMaxConnections defines max connections (int)
	BadMaxConnections = 100 // want "KTN-CONST-001"
	// BadPortNumber defines server port (int)
	BadPortNumber = 8080 // want "KTN-CONST-001"
	// BadTimeoutMs defines timeout in ms (int)
	BadTimeoutMs = 5000 // want "KTN-CONST-001"

	// BadApiVersion defines API version (string)
	BadApiVersion = "v1.0" // want "KTN-CONST-001"
	// BadDefaultLang defines default lang (string)
	BadDefaultLang = "en" // want "KTN-CONST-001"

	// BadIsProduction indicates prod mode (bool)
	BadIsProduction = true // want "KTN-CONST-001"
	// BadEnableCache indicates cache on (bool)
	BadEnableCache = false // want "KTN-CONST-001"

	// BadRatio defines calculation ratio (float64)
	BadRatio = 1.5 // want "KTN-CONST-001"

	// === Additional numeric types without explicit type (T1.3) ===

	// BadByteVal represents a byte value
	BadByteVal = 0xFF // want "KTN-CONST-001"
	// BadInt64Val represents a large int
	BadInt64Val = 9223372036854775807 // want "KTN-CONST-001"
	// BadUintVal represents an unsigned int
	BadUintVal = 42 // want "KTN-CONST-001"
	// BadFloat32Val represents a float
	BadFloat32Val = 3.14 // want "KTN-CONST-001"

	// === Rune type without explicit type ===

	// BadNewlineChar represents a newline
	BadNewlineChar = '\n' // want "KTN-CONST-001"

	// === Iota without explicit type ===

	// BadStateA without explicit type
	BadStateA = iota // want "KTN-CONST-001"
	// BadStateB inherits (no error - inherits from previous)
	BadStateB
	// BadStateC inherits (no error - inherits from previous)
	BadStateC

	// === Multi-name without explicit type ===

	// BadMultiA and BadMultiB without explicit type
	// want "KTN-CONST-001"
	// want "KTN-CONST-001"
	BadMultiA, BadMultiB = 10, 20
)
