package rules_var

import "os"

// Variables correctement documentées et initialisées avec init().

// Variables du package.
var (
	// ConfigPath chemin du fichier de configuration, initialisé depuis l'environnement.
	ConfigPath string
	// UserSettings paramètres utilisateur, initialisé dans init().
	UserSettings map[string]int
	// APIURL URL de l'API externe.
	APIURL string
	// GlobalCounter compteur global initialisé à zéro dans init().
	GlobalCounter int
)

// init initialise le package.
func init() {
	// Initialisation des variables de configuration depuis l'environnement
	ConfigPath = os.Getenv("CONFIG_PATH")
	UserSettings = make(map[string]int)
	APIURL = "https://api.example.com"
	GlobalCounter = 0
}

// Variables de services initialisées au démarrage.
// Variables du package.
var (
	// ServiceRegistry liste des services enregistrés.
	ServiceRegistry []string
	// IsInitialized indique si l'initialisation est complète.
	IsInitialized bool
)

// init initialise le package.
func init() {
	// Enregistrement des services par défaut
	ServiceRegistry = []string{"service1", "service2"}
	IsInitialized = true
}
