package const002

// Good: All constants grouped in a single block before vars
// Respects all rules: explicit types, proper naming, comments, single block grouping
const (
	// CONFIG_VALUE_1 is the first configuration value
	CONFIG_VALUE_1 string = "config1"
	// CONFIG_VALUE_2 is the second configuration value
	CONFIG_VALUE_2 string = "config2"
	// CONFIG_VALUE_3 is the third configuration value
	CONFIG_VALUE_3 string = "config3"

	// MAX_RETRY defines the maximum retry count
	MAX_RETRY int = 5
	// TIMEOUT_SEC defines the timeout in seconds
	TIMEOUT_SEC int = 30
	// SINGLE_CONST is also in the same block (not scattered)
	SINGLE_CONST string = "single"
)

// Variables come after constants
var (
	GlobalVar1 string = "var1"
	GlobalVar2 string = "var2"
)

var AdditionalVar string = "var3"

// Edge case: File with only variables (no const) - should not trigger analyzer
var OnlyVar1 string = "only"
var OnlyVar2 int = 42

// Helper function (not a const/var declaration - tests non-GenDecl branch)
func helperFunction() string {
	return "helper"
}

var SingleVar string = "after single const"
