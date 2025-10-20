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
)

// Variables come after constants
var (
	GlobalVar1 string = "var1"
	GlobalVar2 string = "var2"
)

var AdditionalVar string = "var3"
