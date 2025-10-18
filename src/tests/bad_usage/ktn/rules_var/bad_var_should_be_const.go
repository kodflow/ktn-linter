package rules_var

// BEST PRACTICE: Constantes pour valeurs immuables
// Conforme à KTN-VAR-005 - utilise const pour valeurs jamais réassignées

// Variables du package.
const (
	// MaxRetries définit le nombre maximum de tentatives
	MaxRetries int = 3

	// DefaultPort définit le port par défaut
	DefaultPort int = 8080

	// AppName définit le nom de l'application
	AppName string = "MyApp"

	// Version définit la version de l'application
	Version string = "1.0.0"

	// EnableFeatureX active la fonctionnalité X
	EnableFeatureX bool = true

	// MaxBufferSize définit la taille maximale du buffer
	MaxBufferSize int = 4096

	// DefaultLocale définit la locale par défaut
	DefaultLocale string = "en-US"

	// PiValue définit la valeur de pi
	PiValue float64 = 3.14159

	// MaxUsers définit le nombre maximum d'utilisateurs
	MaxUsers int = 1000

	// APIVersion définit la version de l'API
	APIVersion string = "v2"
)
