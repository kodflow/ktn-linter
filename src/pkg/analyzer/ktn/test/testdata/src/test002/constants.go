package test002

// Fichier qui contient uniquement const/var - devrait être ignoré (pas d'éléments testables)
const (
	MaxRetries = 3
	Timeout    = 30
)

var (
	DefaultHost = "localhost"
	DefaultPort = 8080
)
