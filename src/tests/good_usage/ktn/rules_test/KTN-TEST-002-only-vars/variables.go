package onlyvars

import "os"

// Runtime configuration variables loaded from environment.
var (
	// ServiceURL URL du service depuis l'environnement
	ServiceURL string = os.Getenv("SERVICE_URL")

	// EnableLogging active les logs détaillés
	EnableLogging bool = os.Getenv("ENABLE_LOGGING") == "true"

	// ConnectionPoolSize taille du pool de connexions
	ConnectionPoolSize int = getPoolSize()
)

// getPoolSize retourne la taille du pool depuis l'environnement.
//
// Returns:
//   - int: taille du pool, 100 par défaut
func getPoolSize() int {
	// Cette fonction pourrait lire depuis l'environnement
	return 100
}
