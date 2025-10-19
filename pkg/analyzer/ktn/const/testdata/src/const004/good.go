package const004

// Test cases for KTN-CONST-004: Constants with proper comments

// Good: Individual constant with doc comment
// GoodDocComment is a constant with documentation
const GoodDocComment = "value"

// Good: Constant with line comment
const GoodLineComment = "value" // This is a line comment

/* Good: Constant with block comment */
const GoodBlockComment = "value"

// Good: Multiple constants with group doc comment
// This comment documents all constants in this block
const (
	GoodGrouped1 = "value1"
	GoodGrouped2 = "value2"
	GoodGrouped3 = "value3"
)

// Good: Constants with individual doc comments in a block
const (
	// GoodIndividual1 has its own comment
	GoodIndividual1 = "value1"

	// GoodIndividual2 has its own comment
	GoodIndividual2 = "value2"
)

// Good: Constants with inline comments in a block
const (
	GoodInline1 = "value1" // Inline comment for GoodInline1
	GoodInline2 = "value2" // Inline comment for GoodInline2
)

// Good: Typed constant with comment
// GoodTyped is an integer constant
const GoodTyped int = 42

// Good: Multiple names with comment
// GoodMultiNames defines two related constants
const GoodName1, GoodName2 = 1, 2

// Good: iota pattern with group comment
// Status codes using iota
const (
	StatusOK       = iota // 0
	StatusError           // 1
	StatusPending         // 2
)

// Good: iota pattern with individual comments
const (
	// FlagRead represents read permission
	FlagRead = 1 << iota
	// FlagWrite represents write permission
	FlagWrite
	// FlagExecute represents execute permission
	FlagExecute
)

// Good: String constant with doc comment
// AppName is the application name
const AppName = "MyApp"

// Good: Boolean constant with comment
const IsEnabled = true // Feature flag

// Good: Float constant with doc comment
// Pi is the mathematical constant
const Pi = 3.14159

// Good: Complex constant with comment
const ImaginaryUnit = 0 + 1i // The imaginary unit

// Good: Expression constant with doc comment
// MaxValue is calculated at compile time
const MaxValue = 100 * 1024

// Good: Mix of doc and line comments
const (
	// ConfigTimeout has a doc comment
	ConfigTimeout = 30

	ConfigRetries = 3 // Line comment for retries

	/* Block comment for max connections */
	ConfigMaxConn = 100
)

// Good: Exported constants with proper documentation
// ExportedConstant is visible to other packages
const ExportedConstant = "exported"

// Good: Private constant with comment
// privateConstant is internal only
const privateConstant = "private"

// Good: Numeric constant with doc comment
// DefaultPort is the default server port
const DefaultPort = 8080

// Good: Constant with detailed documentation
// DatabaseURL defines the connection string for the database.
// Format: protocol://host:port/database
const DatabaseURL = "postgresql://localhost:5432/mydb"
