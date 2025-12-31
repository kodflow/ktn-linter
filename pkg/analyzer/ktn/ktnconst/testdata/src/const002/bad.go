// Package const002 contains test cases for KTN rules.
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

// Third scattered const block (ERROR - scattered)
const ( // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc"
	// ThirdConst is another scattered const
	ThirdConst string = "third"
)

// MyType is a type for testing const order.
// This struct demonstrates type declarations in file order.
type MyType struct {
	// Field is a field
	Field string
}

// Fourth const block after type (ERROR - after type + scattered)
const ( // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc" "KTN-CONST-002: les constantes doivent être placées avant les déclarations type"
	// ConstAfterType is after type
	ConstAfterType string = "after_type"
)

// init uses the declarations to avoid unused errors.
//
// Params: none
//
// Returns: none
func init() {
	// Use all declarations
	_ = MyType{}
}

// Fifth const block after func (ERROR - after func + after type + scattered)
const ( // want "KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc" "KTN-CONST-002: les constantes doivent être placées avant les déclarations type" "KTN-CONST-002: les constantes doivent être placées avant les déclarations func"
	// ConstAfterFunc is after func
	ConstAfterFunc string = "after_func"
)
