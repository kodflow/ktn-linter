package interface001_allowed_types

// Test des types publics autorisés (pas besoin de interfaces.go)
type UserID struct {
	value string
}

type ErrorType struct {
	code int
}

type OrderStatus struct {
	status string
}

type AppConfig struct {
	host string
}

type MetaData struct {
	info string
}
