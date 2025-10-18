package const005

// Test cases for KTN-CONST-005: exported constants (valid, no violations)

// Exported constant - uppercase start
const MaxSize = 100

// Exported constant - uppercase with underscore
const API_KEY = "secret"

// Multiple exported constants in group
const (
	DefaultTimeout = 30
	MaxRetries     = 3
)

// Exported with explicit type
const BufferSize int = 1024

// Exported typed constants
const (
	Red   string = "red"
	Green string = "green"
	Blue  string = "blue"
)

// Exported iota pattern
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
)

// Exported with complex type
const ComplexValue complex64 = 1 + 2i

// Exported boolean
const IsEnabled = true

// Exported float
const Pi = 3.14159

// Exported with uppercase letters throughout
const HTTP_STATUS_OK = 200

// Exported CamelCase
const MaxConnectionPoolSize = 100

// Exported PascalCase
const DatabaseConnectionString = "postgres://localhost"

// Blank identifier (should be ignored)
const _ = "ignored"

// Mixed case but starts with uppercase
const HTTPSPort = 443

// Single letter uppercase
const X = 10
const Y = 20

// Exported constants with various types
const (
	MaxInt    int     = 2147483647
	MinFloat  float64 = -1.7976931348623157e+308
	TrueValue bool    = true
	Name      string  = "Application"
)

// Multiple exported constants
const (
	Width, Height = 800, 600
)
