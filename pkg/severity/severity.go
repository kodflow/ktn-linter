package severity

// Level représente le niveau de sévérité d'une règle
type Level int

const (
	// Info recommandations et style
	Info Level = iota
	// Warning problèmes de maintenabilité
	Warning
	// Error bugs potentiels et problèmes graves
	Error
)

// String retourne la représentation textuelle du niveau.
//
// Returns:
//   - string: représentation textuelle
func (l Level) String() string {
	// Vérification du niveau
	switch l {
	// Niveau Info
	case Info:
		// Retour Info
		return "INFO"
	// Niveau Warning
	case Warning:
		// Retour Warning
		return "WARNING"
	// Niveau Error
	case Error:
		// Retour Error
		return "ERROR"
	// Niveau par défaut
	default:
		// Retour Unknown
		return "UNKNOWN"
	}
}

// rulesSeverity mappe chaque règle à son niveau de sévérité
var rulesSeverity = map[string]Level{
	// CONST - Constantes
	"KTN-CONST-001": Warning, // Nommage ALL_CAPS (convention importante)
	"KTN-CONST-002": Info,    // Groupement iota/valeur (style)
	"KTN-CONST-003": Info,    // Groupement logique (organisation)
	"KTN-CONST-004": Warning, // Magic numbers (maintenabilité)

	// VAR - Variables
	"KTN-VAR-001": Warning, // Type explicite manquant (clarté)
	"KTN-VAR-002": Info,    // Groupement (organisation)
	"KTN-VAR-003": Error,   // Réassignation paramètre (bug potentiel)
	"KTN-VAR-004": Error,   // Variables globales mutables (architecture)
	"KTN-VAR-005": Warning, // Portée trop large (maintenabilité)
	"KTN-VAR-006": Info,    // Utiliser := (style)
	"KTN-VAR-007": Warning, // Nommage variables (clarté)
	"KTN-VAR-008": Warning, // Déclaration proche utilisation (lisibilité)
	"KTN-VAR-009": Info,    // var() pour multiple déclarations (style)
	"KTN-VAR-010": Warning, // Nommage receiver (convention)
	"KTN-VAR-011": Warning, // Zero values (clarté)
	"KTN-VAR-012": Warning, // Shadowing (bug potentiel modéré)
	"KTN-VAR-013": Warning, // Unused variables (code mort)
	"KTN-VAR-014": Warning, // Variable names context (clarté)
	"KTN-VAR-015": Warning, // Error variable naming (convention)
	"KTN-VAR-016": Info,    // Boolean prefix (clarté)
	"KTN-VAR-017": Warning, // Slice/map initialization (performance)
	"KTN-VAR-018": Warning, // Constants for fixed values (maintenabilité)
	"KTN-VAR-019": Info,    // Variable grouping by purpose (organisation)

	// FUNC - Fonctions
	"KTN-FUNC-001": Warning, // Fonction trop longue (maintenabilité)
	"KTN-FUNC-002": Warning, // Trop de paramètres (complexité)
	"KTN-FUNC-003": Info,    // Extraire constantes (style léger)
	"KTN-FUNC-004": Info,    // Nommage fonctions (style)
	"KTN-FUNC-005": Info,    // Éviter effets de bord (architecture légère)
	"KTN-FUNC-006": Error,   // Erreur pas en dernière position (convention Go critique)
	"KTN-FUNC-007": Warning, // Documentation manquante (maintenabilité)
	"KTN-FUNC-008": Error,   // Context pas en premier (convention Go critique)
	"KTN-FUNC-009": Warning, // Side effects dans getters (architecture)
	"KTN-FUNC-010": Info,    // Named returns (style)
	"KTN-FUNC-011": Warning, // Commentaires manquants (maintenabilité)
	"KTN-FUNC-012": Error,   // else après return (dead code potentiel)

	// STRUCT - Structures
	"KTN-STRUCT-001": Info,    // Un fichier par struct (organisation)
	"KTN-STRUCT-002": Warning, // Interface manquante (architecture)
	"KTN-STRUCT-003": Info,    // Ordre des champs (style)
	"KTN-STRUCT-004": Warning, // Documentation manquante (maintenabilité)
	"KTN-STRUCT-005": Warning, // Constructeur manquant (architecture)
	"KTN-STRUCT-006": Warning, // Encapsulation manquante (architecture)

	// TEST - Tests
	"KTN-TEST-001": Warning, // Package xxx_test (isolation)
	"KTN-TEST-002": Warning, // Nommage tests (convention)
	"KTN-TEST-003": Warning, // Tests manquants (qualité)
	"KTN-TEST-004": Info,    // Coverage cas d'erreur (qualité)
	"KTN-TEST-005": Warning, // Table-driven tests (maintenabilité)
	"KTN-TEST-006": Warning, // Test helpers (organisation)
}

// GetSeverity retourne le niveau de sévérité d'une règle.
//
// Params:
//   - ruleCode: code de la règle (ex: "KTN-VAR-001")
//
// Returns:
//   - Level: niveau de sévérité
func GetSeverity(ruleCode string) Level {
	// Vérification si la règle existe
	if level, ok := rulesSeverity[ruleCode]; ok {
		// Retour du niveau
		return level
	}
	// Par défaut Warning
	return Warning
}

// ColorCode retourne le code couleur ANSI pour un niveau.
//
// Params:
//   - l: niveau de sévérité
//
// Returns:
//   - string: code couleur ANSI
func (l Level) ColorCode() string {
	// Vérification du niveau
	switch l {
	// Niveau Info
	case Info:
		// Bleu
		return "\033[34m"
	// Niveau Warning
	case Warning:
		// Orange (jaune vif)
		return "\033[33m"
	// Niveau Error
	case Error:
		// Rouge
		return "\033[31m"
	// Niveau par défaut
	default:
		// Blanc
		return "\033[37m"
	}
}

// Symbol retourne le symbole associé au niveau.
//
// Params:
//   - l: niveau de sévérité
//
// Returns:
//   - string: symbole
func (l Level) Symbol() string {
	// Vérification du niveau
	switch l {
	// Niveau Info
	case Info:
		// Symbole info
		return "ℹ"
	// Niveau Warning
	case Warning:
		// Symbole warning
		return "⚠"
	// Niveau Error
	case Error:
		// Symbole error
		return "✖"
	// Niveau par défaut
	default:
		// Symbole générique
		return "●"
	}
}
