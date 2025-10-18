package const001

// Bad: constants without explicit types
const BadPi = 3.14 // want "KTN-CONST-001: la constante 'BadPi' doit avoir un type explicite"

const BadMaxSize = 100 // want "KTN-CONST-001: la constante 'BadMaxSize' doit avoir un type explicite"

const BadMessage = "hello" // want "KTN-CONST-001: la constante 'BadMessage' doit avoir un type explicite"

const BadEnabled = true // want "KTN-CONST-001: la constante 'BadEnabled' doit avoir un type explicite"

// Grouped constants without types
const (
	BadFirst  = 1 // want "KTN-CONST-001: la constante 'BadFirst' doit avoir un type explicite"
	BadSecond = 2 // want "KTN-CONST-001: la constante 'BadSecond' doit avoir un type explicite"
	BadThird  = 3 // want "KTN-CONST-001: la constante 'BadThird' doit avoir un type explicite"
)

// Iota without explicit type
const (
	BadStatusPending   = iota // want "KTN-CONST-001: la constante 'BadStatusPending' doit avoir un type explicite"
	BadStatusRunning          // OK: inherits type from iota pattern
	BadStatusCompleted        // OK: inherits type from iota pattern
)

// Good: constants with explicit types
const GoodPi float64 = 3.14

const GoodMaxSize int = 100

const GoodMessage string = "hello"

const GoodEnabled bool = true

// Grouped constants with types
const (
	GoodFirst  int = 1
	GoodSecond int = 2
	GoodThird  int = 3
)

// Custom type for status
type Status int

// Iota with explicit type
const (
	GoodStatusPending Status = iota
	GoodStatusRunning
	GoodStatusCompleted
)

// Different types
const (
	Timeout      int     = 30
	RetryCount   int     = 3
	DefaultRatio float64 = 1.5
	AppName      string  = "myapp"
)
