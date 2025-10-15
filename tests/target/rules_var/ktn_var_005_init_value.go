package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-005 : Utilisation correcte de var vs const
// ════════════════════════════════════════════════════════════════════════════

// Application metadata
// Ces constantes définissent les métadonnées de l'application (immuables)
const (
	// AppVersionV005Good est la version actuelle de l'application
	AppVersionV005Good string = "1.0.0"
	// AppNameV005Good est le nom de l'application
	AppNameV005Good string = "MyApp"
)

// Retry configuration
// Ces constantes configurent les tentatives (valeurs fixes)
const (
	// MaxRetriesV005Good définit le nombre maximum de tentatives
	MaxRetriesV005Good int = 3
	// TimeoutSecondsV005Good définit le timeout en secondes
	TimeoutSecondsV005Good int = 30
)

// Mathematical constants
// Ces constantes représentent des valeurs mathématiques (immuables)
const (
	// PiV005Good représente la valeur de pi
	PiV005Good float64 = 3.14159265358979323846
	// EulerV005Good représente le nombre d'Euler
	EulerV005Good float64 = 2.71828182845904523536
)

// Feature flags
// Ces constantes activent/désactivent les fonctionnalités (configuration fixe)
const (
	// DebugModeV005Good active le mode debug
	DebugModeV005Good bool = false
	// VerboseLoggingV005Good active les logs verbeux
	VerboseLoggingV005Good bool = true
)

// Counter variables
// Ces variables comptent les événements (modifiées à runtime)
var (
	// requestCountV005Good compte le nombre de requêtes (modifiable)
	requestCountV005Good int = 0
	// errorCountV005Good compte le nombre d'erreurs (modifiable)
	errorCountV005Good int = 0
)

// incrementCounterV005Good incrémente le compteur de requêtes
func incrementCounterV005Good() {
	requestCountV005Good = requestCountV005Good + 1
}

// incrementErrorsV005Good incrémente le compteur d'erreurs
func incrementErrorsV005Good() {
	errorCountV005Good = errorCountV005Good + 1
}
