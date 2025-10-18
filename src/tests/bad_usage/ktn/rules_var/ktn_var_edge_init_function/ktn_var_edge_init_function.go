package badvarinit

import "os"

// Violations avec init() et variables

var (
	config_path  string         // Snake_case violation
	userSettings map[string]int // Pas d'initialisation explicite
	API_URL      string         // Screaming snake
)

// Variable modifiée dans init sans documentation claire
var globalCounter int

func init() {
	// Modification de variables sans documentation
	config_path = os.Getenv("CONFIG_PATH")
	userSettings = make(map[string]int)
	API_URL = "https://api.example.com"
	globalCounter = 0
}

// Variable avec init complexe
var (
	// serviceRegistry describes this variable.
	serviceRegistry []string
	// isInitialized describes this variable.
	isInitialized bool
)

func init() {
	serviceRegistry = []string{"service1", "service2"}
	isInitialized = true
}
