// Package const002 contains test cases for KTN-CONST-002.
// This file contains ALL cases that MUST trigger KTN-CONST-002 errors.
package const002

// =============================================================================
// SECTION 1: Multiple const blocks at top (now VALID - all before var/type/func)
// These do NOT trigger errors because they're all at the top of the file.
// =============================================================================

// FirstConstBlock is at the top (OK - first const is always OK)
const FirstConstBlock string = "first"

// SecondConstBlock is at the top (OK - before any var/type/func)
const SecondConstBlock string = "second"

// ThirdConstBlock is also at the top (OK - before any var/type/func)
const ThirdConstBlock string = "third"

// FourthConstBlock group is at the top (OK - before any var/type/func)
const (
	// FourthA is constant A
	FourthA string = "a"
	// FourthB is constant B
	FourthB string = "b"
)

// =============================================================================
// SECTION 2: Const after var declarations
// =============================================================================

// badGlobalVar is a variable
var badGlobalVar string = "var"

// ConstAfterVar appears after var - ERROR
const ConstAfterVar string = "after_var" // want "KTN-CONST-002"

// =============================================================================
// SECTION 3: Const intercalated between var and type
// =============================================================================

// badSecondVar is another variable
var badSecondVar int = 42

// ConstBetweenVarAndType is intercalated - ERROR
const ConstBetweenVarAndType string = "between" // want "KTN-CONST-002"

// =============================================================================
// SECTION 4: Const after type declarations
// =============================================================================

// BadType is a type declaration
type BadType struct {
	// Field is a field
	Field string
}

// ConstAfterType appears after type - ERROR
const ConstAfterType string = "after_type" // want "KTN-CONST-002"

// =============================================================================
// SECTION 5: Const with iota using UNDEFINED type (not valid exception)
// =============================================================================

// ConstWithUndefinedType uses a type not defined in this file - ERROR
const ( // want "KTN-CONST-002"
	// BadIotaA uses UndefinedType which is not in this file
	BadIotaA int = iota // Using built-in type, not custom = scattered
	// BadIotaB inherits
	BadIotaB
)

// =============================================================================
// SECTION 6: Const after func declarations
// =============================================================================

// badInit uses the declarations to avoid unused errors.
//
// Returns: none
func badInit() {
	// Use all declarations
	_ = badGlobalVar
	_ = badSecondVar
	_ = BadType{}
}

// ConstAfterFunc appears after func - ERROR (all violations)
const ConstAfterFunc string = "after_func" // want "KTN-CONST-002"

// =============================================================================
// SECTION 7: Const block after func with multiple constants - ERROR
// =============================================================================

// MultiConstAfterFunc is a group after func - ERROR
const ( // want "KTN-CONST-002"
	// AfterFuncA is after func
	AfterFuncA string = "func_a"
	// AfterFuncB is after func
	AfterFuncB string = "func_b"
	// AfterFuncC is after func
	AfterFuncC string = "func_c"
)

// =============================================================================
// SECTION 8: Const with custom type AFTER func (no exception for after func)
// Even with custom type iota pattern, const after func is always error
// =============================================================================

// BadCustomType is a custom type defined before func
type BadCustomType int

// badSecondFunc is another function
//
// Returns: none
func badSecondFunc() {
	// Use type
	_ = BadCustomType(0)
}

// ConstWithCustomTypeAfterFunc uses custom type but is after func - ERROR
// The iota exception only applies to const after TYPE, not after FUNC
const ( // want "KTN-CONST-002"
	// BadCustomA uses custom type but after func
	BadCustomA BadCustomType = iota
	// BadCustomB inherits
	BadCustomB
)
