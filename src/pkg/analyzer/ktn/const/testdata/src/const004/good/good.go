package good

// Bon : type explicite
const (
	// MaxRetries is the maximum retries
	MaxRetries int = 3
	// Timeout is the timeout duration
	Timeout int = 30
)

// Bon : type explicite pour string
const (
	// StatusActive is active
	StatusActive string = "active"
	// StatusInactive is inactive
	StatusInactive string = "inactive"
)

// Exception : iota n'a pas besoin de type explicite
const (
	// First value with iota
	First = iota
	// Second value with iota
	Second
	// Third value with iota
	Third
)

// Exception : iota avec opérations binaires
const (
	// KB uses iota with binary operation
	KB = 1 << (10 * iota)
	// MB value
	MB
	// GB value
	GB
)

// Exception : iota avec addition
const (
	// Offset uses iota + 100
	Offset = iota + 100
	// Next value
	Next
)

// Exception : iota avec unary
const (
	// Negative uses -iota
	Negative = -iota
)

// Exception : iota avec parenthèses
const (
	// Wrapped uses (iota)
	Wrapped = (iota)
)

// Exception : iota dans un appel de fonction
const (
	// Converted uses byte(iota)
	Converted = byte(iota)
	// ConvertedNext value
	ConvertedNext
)
