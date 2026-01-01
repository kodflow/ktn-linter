// Package const002 contains test cases for KTN-CONST-002.
// This file contains ALL cases that MUST NOT trigger KTN-CONST-002 errors.
package const002

// =============================================================================
// SECTION 1: Single inline const at the top (valid - first const)
// =============================================================================

// GoodSingleConst is a single inline const at the very top
const GoodSingleConst string = "single"

// =============================================================================
// SECTION 2: All constants grouped in a single block
// Correct order: const → var → type → func
// =============================================================================

const (
	// GoodConfigA is the first configuration value
	GoodConfigA string = "config_a"

	// GoodConfigB is the second configuration value
	GoodConfigB string = "config_b"

	// GoodConfigC is the third configuration value
	GoodConfigC string = "config_c"

	// GoodMaxRetry defines the maximum retry count
	GoodMaxRetry int = 5

	// GoodTimeout defines the timeout in seconds
	GoodTimeout int = 30

	// GoodEnabled is a boolean constant
	GoodEnabled bool = true

	// GoodRatio is a float constant
	GoodRatio float64 = 3.14
)

// =============================================================================
// SECTION 3: Iota patterns in main const block
// =============================================================================

const (
	// GoodStatusPending demonstrates iota with standard type
	GoodStatusPending int = iota
	// GoodStatusActive inherits iota
	GoodStatusActive
	// GoodStatusDone inherits iota
	GoodStatusDone
)

const (
	// GoodFlagA demonstrates bit flags with iota
	GoodFlagA int = 1 << iota
	// GoodFlagB inherits bit shift pattern
	GoodFlagB
	// GoodFlagC inherits bit shift pattern
	GoodFlagC
)

// =============================================================================
// SECTION 4: Expression list shorthand (omitting expression)
// =============================================================================

const (
	// GoodRepeatA is the first value
	GoodRepeatA int = 100
	// GoodRepeatB reuses previous expression (implicit = 100)
	GoodRepeatB
	// GoodRepeatC also reuses (implicit = 100)
	GoodRepeatC
)

// =============================================================================
// SECTION 5: Multiple iota in same ConstSpec
// =============================================================================

const (
	// GoodBit0 and GoodMask0 use iota in same line
	GoodBit0, GoodMask0 int = 1 << iota, 1<<iota - 1 // bit0=1, mask0=0
	// GoodBit1 and GoodMask1 inherit pattern
	GoodBit1, GoodMask1 // bit1=2, mask1=1
	// GoodBit2 and GoodMask2 inherit pattern
	GoodBit2, GoodMask2 // bit2=4, mask2=3
)

// =============================================================================
// SECTION 6: Blank identifier in const (should be ignored)
// =============================================================================

const (
	// Blank identifier to skip first iota value
	_ int = iota
	// GoodAfterBlank starts at 1
	GoodAfterBlank
	// GoodAfterBlank2 is 2
	GoodAfterBlank2
)

// =============================================================================
// SECTION 7: Variables come after constants (correct order)
// =============================================================================

var (
	// goodGlobalVar1 is the first global variable
	goodGlobalVar1 string = "var1"
	// goodGlobalVar2 is the second global variable
	goodGlobalVar2 string = "var2"
	// goodGlobalVar3 is the third global variable
	goodGlobalVar3 int = 42
)

// =============================================================================
// SECTION 8: Type declarations come after variables (correct order)
// =============================================================================

// GoodStruct is a struct type.
type GoodStruct struct {
	// Name is the name field
	Name string
	// Value is the value field
	Value int
}

// GoodAlias is a type alias.
type GoodAlias = string

// GoodNewType is a new type based on int.
type GoodNewType int

// Status is a custom type for status values (used for iota exception).
type Status int

// Priority is another custom type for priority values.
type Priority int

// Level is a third custom type for testing multiple type exceptions.
type Level int

// =============================================================================
// SECTION 9: Exception - Const with iota using custom type AFTER type declaration
// This is allowed because the const block uses a type defined in the same file
// =============================================================================

const (
	// StatusUnknown is the unknown status - uses custom type Status
	StatusUnknown Status = iota
	// StatusPending2 is the pending status (renamed to avoid conflict)
	StatusPending2
	// StatusRunning is the running status
	StatusRunning
	// StatusCompleted is the completed status
	StatusCompleted
	// StatusFailed is the failed status
	StatusFailed
)

const (
	// PriorityLow is low priority - uses custom type Priority
	PriorityLow Priority = iota
	// PriorityMedium is medium priority
	PriorityMedium
	// PriorityHigh is high priority
	PriorityHigh
	// PriorityCritical is critical priority
	PriorityCritical
)

const (
	// LevelDebug is debug level - third custom type
	LevelDebug Level = iota
	// LevelInfo is info level
	LevelInfo
	// LevelWarn is warn level
	LevelWarn
	// LevelError is error level
	LevelError
)

// =============================================================================
// SECTION 10: Functions come last (correct order)
// =============================================================================

// goodInit demonstrates func declarations after all other declarations.
//
// Returns: none
func goodInit() {
	// Use declarations to avoid unused errors
	_ = GoodSingleConst
	_ = GoodConfigA
	_ = goodGlobalVar1
	_ = GoodStruct{}
	_ = GoodAlias("test")
	_ = GoodNewType(0)
	_ = StatusUnknown
	_ = PriorityLow
	_ = LevelDebug
	_ = GoodBit0
	_ = GoodMask0
	_ = GoodAfterBlank
}

// goodHelper is another function at the end.
//
// Returns: the helper result
func goodHelper() string {
	// Return helper result
	return GoodConfigB
}
