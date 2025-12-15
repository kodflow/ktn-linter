// Package const001 contains test cases for KTN rules.
package const001

// Bad: All constants WITHOUT explicit types (violates KTN-CONST-001)
// Names use CamelCase to avoid CONST-003 errors
const (
	// BadMaxConnections defines max connections
	BadMaxConnections = 100 // want "KTN-CONST-001"
	// BadPortNumber defines server port
	BadPortNumber = 8080 // want "KTN-CONST-001"
	// BadTimeoutMs defines timeout in ms
	BadTimeoutMs = 5000 // want "KTN-CONST-001"

	// BadHttpOk represents HTTP 200 status
	BadHttpOk = 200 // want "KTN-CONST-001"
	// BadHttpNotFound represents 404
	BadHttpNotFound = 404 // want "KTN-CONST-001"

	// BadApiVersion defines API version
	BadApiVersion = "v1.0" // want "KTN-CONST-001"
	// BadDefaultLang defines default lang
	BadDefaultLang = "en" // want "KTN-CONST-001"

	// BadIsProduction indicates prod mode
	BadIsProduction = true // want "KTN-CONST-001"
	// BadEnableCache indicates cache on
	BadEnableCache = false // want "KTN-CONST-001"

	// BadRatio defines calculation ratio
	BadRatio = 1.5 // want "KTN-CONST-001"

	// BadStateA without explicit type
	BadStateA = iota // want "KTN-CONST-001"
	// BadStateB inherits (no error)
	BadStateB
	// BadStateC inherits (no error)
	BadStateC

	// BadMultiA without explicit type
	// want "KTN-CONST-001"
	// want "KTN-CONST-001"
	BadMultiA, BadMultiB = 10, 20
)

// BadStatus is a type declared after const (for testing purposes only)
type BadStatus int
