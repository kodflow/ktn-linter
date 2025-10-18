package const005

// Test cases for KTN-CONST-005: unexported constants (should trigger violations)

// Unexported constant - lowercase start
const maxSize = 100 // want "KTN-CONST-005: la constante 'maxSize' doit être exportée"

// Unexported constant - snake_case
const api_key = "secret" // want "KTN-CONST-005: la constante 'api_key' doit être exportée"

// Multiple unexported constants in group
const (
	defaultTimeout = 30 // want "KTN-CONST-005: la constante 'defaultTimeout' doit être exportée"
	maxRetries     = 3  // want "KTN-CONST-005: la constante 'maxRetries' doit être exportée"
)

// Mixed exported and unexported
const (
	PublicConst  = "public"
	privateConst = "private" // want "KTN-CONST-005: la constante 'privateConst' doit être exportée"
)

// Unexported with explicit type
const bufferSize int = 1024 // want "KTN-CONST-005: la constante 'bufferSize' doit être exportée"

// Unexported typed constants
const (
	red   string = "red"   // want "KTN-CONST-005: la constante 'red' doit être exportée"
	green string = "green" // want "KTN-CONST-005: la constante 'green' doit être exportée"
	blue  string = "blue"  // want "KTN-CONST-005: la constante 'blue' doit être exportée"
)

// Unexported iota pattern
const (
	sunday    = iota // want "KTN-CONST-005: la constante 'sunday' doit être exportée"
	monday           // want "KTN-CONST-005: la constante 'monday' doit être exportée"
	tuesday          // want "KTN-CONST-005: la constante 'tuesday' doit être exportée"
	wednesday        // want "KTN-CONST-005: la constante 'wednesday' doit être exportée"
)

// Unexported with complex type
const complexValue complex64 = 1 + 2i // want "KTN-CONST-005: la constante 'complexValue' doit être exportée"

// Unexported boolean
const isEnabled = true // want "KTN-CONST-005: la constante 'isEnabled' doit être exportée"

// Unexported float
const pi = 3.14159 // want "KTN-CONST-005: la constante 'pi' doit être exportée"

// Unexported with underscore prefix
const _internal = "internal" // want "KTN-CONST-005: la constante '_internal' doit être exportée"

// Multiple unexported on single line (if supported)
const (
	width, height = 800, 600 // want "KTN-CONST-005: la constante 'width' doit être exportée" "KTN-CONST-005: la constante 'height' doit être exportée"
)
