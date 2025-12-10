// Bad examples for the const002 test case.
package const002

// First const block (OK - at the top)
const (
	// FirstConst is the first constant
	FirstConst string = "first"
)

// Second scattered const block (ERROR - scattered)
const ( // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
	// ScatteredConst is scattered
	ScatteredConst string = "scattered"
)

// Variable declaration
var globalVar string = "var"

// Third const block after var (ERROR - after var + scattered)
const ( // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc" "KTN-CONST-002: les constantes doivent être placées avant les déclarations var"
	// ConstAfterVar is after var
	ConstAfterVar string = "after_var"
)

// MyType is a type for testing const order.
// This struct demonstrates type declarations in file order.
type MyType struct {
	// Field is a field
	Field string
}

// Fourth const block after type (ERROR - after type + after var + scattered)
const ( // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc" "KTN-CONST-002: les constantes doivent être placées avant les déclarations var" "KTN-CONST-002: les constantes doivent être placées avant les déclarations type"
	// ConstAfterType is after type
	ConstAfterType string = "after_type"
)

// init uses the declarations to avoid unused errors
func init() {
	// Use all declarations
	_ = globalVar
	_ = MyType{}
}

// Fifth const block after func (ERROR - all violations)
const ( // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc" "KTN-CONST-002: les constantes doivent être placées avant les déclarations var" "KTN-CONST-002: les constantes doivent être placées avant les déclarations type" "KTN-CONST-002: les constantes doivent être placées avant les déclarations func"
	// ConstAfterFunc is after func
	ConstAfterFunc string = "after_func"
)
