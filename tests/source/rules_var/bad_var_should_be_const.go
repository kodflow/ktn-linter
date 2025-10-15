package rules_var

// ANTI-PATTERN: Variables jamais réassignées avec valeurs littérales
// Viole KTN-VAR-005 - devraient être des constantes !

// Variables du package.
var (
	// MaxRetries jamais modifié - devrait être const !
	MaxRetries int = 3

	// DefaultPort jamais modifié - devrait être const !
	DefaultPort int = 8080

	// AppName jamais modifié - devrait être const !
	AppName string = "MyApp"

	// Version jamais modifié - devrait être const !
	Version string = "1.0.0"

	// EnableFeatureX jamais modifié - devrait être const !
	EnableFeatureX bool = true

	// MaxBufferSize jamais modifié - devrait être const !
	MaxBufferSize int = 4096

	// DefaultLocale jamais modifié - devrait être const !
	DefaultLocale string = "en-US"

	// PiValue jamais modifié - devrait être const !
	PiValue float64 = 3.14159

	// MaxUsers jamais modifié - devrait être const !
	MaxUsers int = 1000

	// APIVersion jamais modifié - devrait être const !
	APIVersion string = "v2"
)
