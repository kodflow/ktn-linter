package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-005 : Variable jamais réassignée devrait être const
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Une variable qui n'est jamais réassignée après sa déclaration initiale
//    devrait être déclarée comme const pour indiquer son immuabilité.
//
//    POURQUOI :
//    - const est thread-safe par nature (immuable)
//    - Indique clairement l'intention (ne changera jamais)
//    - Évite les modifications accidentelles
//    - Optimisations possibles par le compilateur
//
//    DÉTECTION :
//    - Analyse statique du package entier
//    - Détecte si une variable est jamais réassignée (write vs read)
//    - Ne signale que les types compatibles avec const (bool, string, int*, float*, etc.)
//
// ✅ CAS PARFAIT (utiliser const) :
//
//    // Application metadata
//    // Ces constantes définissent les métadonnées de l'application
//    const (
//        // AppVersion est la version actuelle de l'application
//        AppVersion string = "1.0.0"
//        // MaxRetries est le nombre maximum de tentatives
//        MaxRetries int = 3
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ✅ CAS CORRECT 1 : Constantes string (corrigé)
// CONFORME : KTN-VAR-005 (utilise const pour valeurs immuables)
// Application configuration
// Ces constantes définissent la configuration de l'application
const (
	// AppVersionV005 est la version de l'application
	AppVersionV005 string = "1.0.0"
	// AppNameV005 est le nom de l'application
	AppNameV005 string = "MyApp"
)

// ✅ CAS CORRECT 2 : Constantes int (corrigé)
// CONFORME : KTN-VAR-005 (utilise const pour valeurs immuables)
// Retry configuration
// Ces constantes configurent les tentatives
const (
	// MaxRetriesV005 définit le nombre maximum de tentatives
	MaxRetriesV005 int = 3
	// TimeoutSecondsV005 définit le timeout en secondes
	TimeoutSecondsV005 int = 30
)

// ✅ CAS CORRECT 3 : Constantes float64 (corrigé)
// CONFORME : KTN-VAR-005 (utilise const pour valeurs immuables)
// Mathematical constants
// Ces constantes représentent des valeurs mathématiques
const (
	// PiV005 représente la valeur de pi
	PiV005 float64 = 3.14159265358979323846
	// EulerV005 représente le nombre d'Euler
	EulerV005 float64 = 2.71828182845904523536
)

// ✅ CAS CORRECT 4 : Constantes bool (corrigé)
// CONFORME : KTN-VAR-005 (utilise const pour valeurs immuables)
// Feature flags
// Ces constantes activent/désactivent les fonctionnalités
const (
	// DebugModeV005 active le mode debug
	DebugModeV005 bool = false
	// VerboseLoggingV005 active les logs verbeux
	VerboseLoggingV005 bool = true
)

// ✅ CAS CORRECT : Variable réassignée (pas d'erreur, var est approprié)
// Counter variables
// Ces variables comptent les événements (modifiées à runtime)
var (
	// requestCountV005 compte le nombre de requêtes (réassigné plus bas)
	requestCountV005 int = 0
)

// Fonction qui réassigne requestCountV005 (donc var est correct)
func incrementCounterV005() {
	requestCountV005 = requestCountV005 + 1 // Réassignation ici
}
