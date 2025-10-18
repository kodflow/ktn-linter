package interface002_allowed

// Types autorisés - pas d'erreur attendue
type UserID struct {
	value string
}

type AppConfig struct {
	host string
	port int
}

type ItemCount struct {
	total int
}
