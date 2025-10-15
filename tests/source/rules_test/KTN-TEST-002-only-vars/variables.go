package onlyvars

import "os"

var (
	// ServiceURL URL du service depuis l'environnement
	ServiceURL string = os.Getenv("SERVICE_URL")

	// DebugMode mode debug activé
	DebugMode bool = false

	// MaxConnections nombre maximum de connexions
	MaxConnections int = 100
)
