package mock004

// Fichier avec des déclarations mais pas d'interfaces

const (
	MaxRetries = 3
	Timeout    = 30
)

var (
	DefaultValue = 42
	ConfigPath   = "/etc/config"
)

type Alias = string

type Config struct {
	Name  string
	Value int
}
