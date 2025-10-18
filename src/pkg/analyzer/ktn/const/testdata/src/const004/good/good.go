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
