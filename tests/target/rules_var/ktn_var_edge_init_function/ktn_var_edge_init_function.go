package rules_var

import "os"

// Variables correctement documentées et initialisées avec init().

var (
	// ConfigPath chemin du fichier de configuration, initialisé depuis l'environnement.
	ConfigPath string
	// UserSettings paramètres utilisateur, initialisé dans init().
	UserSettings map[string]int
	// APIURL URL de l'API externe.
	APIURL string
)

// GlobalCounter compteur global initialisé à zéro dans init().
var GlobalCounter int

func init() {
	// Initialisation des variables de configuration depuis l'environnement
	ConfigPath = os.Getenv("CONFIG_PATH")
	UserSettings = make(map[string]int)
	APIURL = "https://api.example.com"
	GlobalCounter = 0
}

// Variables de services initialisées au démarrage.
var (
	// ServiceRegistry liste des services enregistrés.
	ServiceRegistry []string
	// IsInitialized indique si l'initialisation est complète.
	IsInitialized bool
)

func init() {
	// Enregistrement des services par défaut
	ServiceRegistry = []string{"service1", "service2"}
	IsInitialized = true
}
