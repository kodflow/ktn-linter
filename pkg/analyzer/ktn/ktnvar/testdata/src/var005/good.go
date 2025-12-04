// Good examples for the var005 test case.
package var005

const (
	// MAX_RETRIES_VALUE is the max retries value
	MAX_RETRIES_VALUE int = 3
	// ANSWER is the answer
	ANSWER int = 42
	// AGE is example age
	AGE int = 30
	// X_VALUE is example x value
	X_VALUE int = 10
	// Y_VALUE is example y value
	Y_VALUE int = 20
	// INT64_VALUE is int64 example
	INT64_VALUE int64 = 5
	// FLOAT_VALUE is float example
	FLOAT_VALUE float64 = 3.14
)

// Package-level variables can use var with explicit type
var (
	// GlobalConfig is the global configuration
	GlobalConfig string = "default"
	// MaxRetries is the maximum retries
	MaxRetries int = MAX_RETRIES_VALUE
)

// goodShortDeclaration shows correct use of short declaration
func goodShortDeclaration() {
	// Use := for local variables
	name := "Alice"
	age := AGE
	isActive := true

	_ = name
	_ = age
	_ = isActive
}

// goodMultipleAssignment shows correct use of := with multiple variables
func goodMultipleAssignment() {
	// Use := for multiple variables
	x, y := X_VALUE, Y_VALUE

	_ = x
	_ = y
}

// goodExplicitTypeWhenNeeded shows when var is acceptable with explicit type
func goodExplicitTypeWhenNeeded() {
	// This is OK: explicit type conversion is necessary
	// x is an int64 value
	var x int64 = INT64_VALUE
	// y is a float64 value
	var y float64 = FLOAT_VALUE

	_ = x
	_ = y
}

// goodVarWithoutInit shows var without initialization is OK
func goodVarWithoutInit() {
	// This is OK: no initialization
	var result int
	result = ANSWER

	_ = result
}

// goodMultipleVarWithoutInit shows multiple var without init is OK
func goodMultipleVarWithoutInit() {
	const (
		// A_VALUE is value for a
		A_VALUE int = 1
		// B_VALUE is value for b
		B_VALUE int = 2
		// C_VALUE is value for c
		C_VALUE int = 3
	)
	// This is OK: no initialization
	var a, b, c int
	a, b, c = A_VALUE, B_VALUE, C_VALUE

	_ = a
	_ = b
	_ = c
}

// goodReassignment shows that := is used for new variables
func goodReassignment() {
	const (
		// INITIAL_COUNT is the initial count
		INITIAL_COUNT int = 0
		// NEW_COUNT is the new count
		NEW_COUNT int = 5
	)
	count := INITIAL_COUNT
	count = NEW_COUNT // Regular assignment is fine

	_ = count
}

// init utilise les fonctions privées
func init() {
	// Appel de goodShortDeclaration
	goodShortDeclaration()
	// Appel de goodMultipleAssignment
	goodMultipleAssignment()
	// Appel de goodExplicitTypeWhenNeeded
	goodExplicitTypeWhenNeeded()
	// Appel de goodVarWithoutInit
	goodVarWithoutInit()
	// Appel de goodMultipleVarWithoutInit
	goodMultipleVarWithoutInit()
	// Appel de goodReassignment
	goodReassignment()
}
