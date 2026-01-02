// Package var002 contains test cases for KTN rules.
package var002

// Bad: Variables without explicit type (violates KTN-VAR-001)
// Note: Zero-values (type without init) are now valid

const (
	// MaxRetries defines maximum retries
	MaxRetries int = 3
	// PortValue is port value
	PortValue int = 8080
	// RatioValue is ratio value
	RatioValue float64 = 1.5
)

// Cas: Pas de type explicite = ERREUR
var (
	// badRetries has no explicit type
	badRetries = MaxRetries // want "KTN-VAR-001"

	// badConfig has no explicit type
	badConfig = "config" // want "KTN-VAR-001"

	// badPort has no explicit type
	badPort = PortValue // want "KTN-VAR-001"

	// badHost has no explicit type
	badHost = "localhost" // want "KTN-VAR-001"

	// badEnabled has no explicit type
	badEnabled = true // want "KTN-VAR-001"

	// badRatio has no explicit type
	badRatio = RatioValue // want "KTN-VAR-001"

	// badSlice has no explicit type
	badSlice = []string{"a", "b"} // want "KTN-VAR-001"

	// badMap has no explicit type
	badMap = map[string]int{"x": 1} // want "KTN-VAR-001"
)
