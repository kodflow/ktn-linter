// Bad examples for the var003 test case.
package var003

// Bad: Variables with incorrect type visibility (violates KTN-VAR-003)

const (
	// MAX_RETRIES defines maximum retries
	MAX_RETRIES int = 3
	// PORT_VALUE is port value
	PORT_VALUE int = 8080
	// RATIO_VALUE is ratio value
	RATIO_VALUE float64 = 1.5
	// CAP_VALUE for slice capacity
	CAP_VALUE int = 10
	// MAGIC_VALUE for conversion
	MAGIC_VALUE int = 42
)

// Cas 1: Type non visible sans type explicite → ERREUR
// Cas 2: Type redondant (type visible + type explicite) → ERREUR
var (
	// badRetries has inferred type (should be explicit)
	badRetries = MAX_RETRIES // want "KTN-VAR-003: la variable 'badRetries' doit avoir un type explicite"

	// badConfig has inferred type (should be explicit)
	badConfig = "config" // want "KTN-VAR-003: la variable 'badConfig' doit avoir un type explicite"

	// badPort has inferred type
	badPort = PORT_VALUE // want "KTN-VAR-003: la variable 'badPort' doit avoir un type explicite"

	// badHost has inferred type
	badHost = "localhost" // want "KTN-VAR-003: la variable 'badHost' doit avoir un type explicite"

	// badEnabled has inferred type
	badEnabled = true // want "KTN-VAR-003: la variable 'badEnabled' doit avoir un type explicite"

	// badRatio has inferred type
	badRatio = RATIO_VALUE // want "KTN-VAR-003: la variable 'badRatio' doit avoir un type explicite"

	// redundantSlice has redundant type
	redundantSlice []string = []string{"a", "b"} // want "KTN-VAR-003: la variable 'redundantSlice' a un type redondant"

	// redundantMap has redundant type
	redundantMap map[string]int = map[string]int{"x": 1} // want "KTN-VAR-003: la variable 'redundantMap' a un type redondant"

	// redundantMake has redundant type
	redundantMake []byte = make([]byte, 0, CAP_VALUE) // want "KTN-VAR-003: la variable 'redundantMake' a un type redondant"

	// redundantConv has redundant type
	redundantConv int = int(MAGIC_VALUE) // want "KTN-VAR-003: la variable 'redundantConv' a un type redondant"
)
