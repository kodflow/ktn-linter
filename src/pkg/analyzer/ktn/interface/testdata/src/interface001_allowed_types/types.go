package interface001_allowed_types

// Test des types publics autorisés (pas besoin de interfaces.go)
// UserID represents the struct.
type UserID struct {
	value string
}
// ErrorType represents the struct.

type ErrorType struct {
	code int
// OrderStatus represents the struct.
}

type OrderStatus struct {
// AppConfig represents the struct.
	status string
}

// MetaData represents the struct.
type AppConfig struct {
	host string
}

type MetaData struct {
	info string
}
