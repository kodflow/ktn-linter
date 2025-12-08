// Bad examples for the var002 test case.
package var002

// Bad: Variables without explicit type OR without value (violates KTN-VAR-002)

const (
	// MAX_RETRIES defines maximum retries
	MAX_RETRIES int = 3
	// PORT_VALUE is port value
	PORT_VALUE int = 8080
	// RATIO_VALUE is ratio value
	RATIO_VALUE float64 = 1.5
)

// Cas 1: Pas de type explicite = ERREUR
// Cas 2: Pas de valeur d'initialisation = ERREUR
var (
	// badRetries has no explicit type
	badRetries = MAX_RETRIES // want "KTN-VAR-002: la variable 'badRetries' doit avoir un type explicite"

	// badConfig has no explicit type
	badConfig = "config" // want "KTN-VAR-002: la variable 'badConfig' doit avoir un type explicite"

	// badPort has no explicit type
	badPort = PORT_VALUE // want "KTN-VAR-002: la variable 'badPort' doit avoir un type explicite"

	// badHost has no explicit type
	badHost = "localhost" // want "KTN-VAR-002: la variable 'badHost' doit avoir un type explicite"

	// badEnabled has no explicit type
	badEnabled = true // want "KTN-VAR-002: la variable 'badEnabled' doit avoir un type explicite"

	// badRatio has no explicit type
	badRatio = RATIO_VALUE // want "KTN-VAR-002: la variable 'badRatio' doit avoir un type explicite"

	// badSlice has no explicit type
	badSlice = []string{"a", "b"} // want "KTN-VAR-002: la variable 'badSlice' doit avoir un type explicite"

	// badMap has no explicit type
	badMap = map[string]int{"x": 1} // want "KTN-VAR-002: la variable 'badMap' doit avoir un type explicite"

	// uninitializedVar has type but no value
	uninitializedVar int // want "KTN-VAR-002: la variable 'uninitializedVar' doit être initialisée"

	// uninitializedString has type but no value
	uninitializedString string // want "KTN-VAR-002: la variable 'uninitializedString' doit être initialisée"
)
