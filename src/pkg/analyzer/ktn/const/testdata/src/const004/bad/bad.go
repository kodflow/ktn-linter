package bad

// Mauvais : pas de type explicite
const (
	// MaxRetries is the maximum retries
	MaxRetries = 3 // want `\[KTN_CONST_004\] Constante 'MaxRetries' sans type explicite`
	// Timeout is the timeout duration
	Timeout = 30 // want `\[KTN_CONST_004\] Constante 'Timeout' sans type explicite`
)

// Mauvais : pas de type explicite pour string
const (
	// StatusActive is active
	StatusActive = "active" // want `\[KTN_CONST_004\] Constante 'StatusActive' sans type explicite`
	// StatusInactive is inactive
	StatusInactive = "inactive" // want `\[KTN_CONST_004\] Constante 'StatusInactive' sans type explicite`
)
