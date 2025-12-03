// Severity levels for linter diagnostics.
package severity

// Level représente le niveau de sévérité d'une règle
type Level int

const (
	// SEVERITY_INFO recommandations et style
	SEVERITY_INFO Level = iota
	// SEVERITY_WARNING problèmes de maintenabilité
	SEVERITY_WARNING
	// SEVERITY_ERROR bugs potentiels et problèmes graves
	SEVERITY_ERROR
)

// String retourne la représentation textuelle du niveau.
//
// Returns:
//   - string: représentation textuelle
func (l Level) String() string {
	// Vérification du niveau
	switch l {
	// Niveau SEVERITY_INFO
	case SEVERITY_INFO:
		// Retour SEVERITY_INFO
		return "INFO"
	// Niveau SEVERITY_WARNING
	case SEVERITY_WARNING:
		// Retour SEVERITY_WARNING
		return "WARNING"
	// Niveau SEVERITY_ERROR
	case SEVERITY_ERROR:
		// Retour SEVERITY_ERROR
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
	"KTN-CONST-001": SEVERITY_WARNING, // Nommage ALL_CAPS (convention importante)
	"KTN-CONST-002": SEVERITY_INFO,    // Groupement iota/valeur (style)
	"KTN-CONST-003": SEVERITY_INFO,    // Groupement logique (organisation)
	"KTN-CONST-004": SEVERITY_WARNING, // Magic numbers (maintenabilité)

	// VAR - Variables
	"KTN-VAR-001": SEVERITY_WARNING, // Type explicite manquant (clarté)
	"KTN-VAR-002": SEVERITY_INFO,    // Groupement (organisation)
	"KTN-VAR-003": SEVERITY_ERROR,   // Réassignation paramètre (bug potentiel)
	"KTN-VAR-004": SEVERITY_ERROR,   // Variables globales mutables (architecture)
	"KTN-VAR-005": SEVERITY_WARNING, // Portée trop large (maintenabilité)
	"KTN-VAR-006": SEVERITY_INFO,    // Utiliser := (style)
	"KTN-VAR-007": SEVERITY_WARNING, // Nommage variables (clarté)
	"KTN-VAR-008": SEVERITY_WARNING, // Déclaration proche utilisation (lisibilité)
	"KTN-VAR-009": SEVERITY_INFO,    // var() pour multiple déclarations (style)
	"KTN-VAR-010": SEVERITY_WARNING, // Nommage receiver (convention)
	"KTN-VAR-011": SEVERITY_WARNING, // Zero values (clarté)
	"KTN-VAR-012": SEVERITY_WARNING, // Shadowing (bug potentiel modéré)
	"KTN-VAR-013": SEVERITY_WARNING, // Unused variables (code mort)
	"KTN-VAR-014": SEVERITY_WARNING, // Variable names context (clarté)
	"KTN-VAR-015": SEVERITY_WARNING, // Error variable naming (convention)
	"KTN-VAR-016": SEVERITY_INFO,    // Boolean prefix (clarté)
	"KTN-VAR-017": SEVERITY_WARNING, // Slice/map initialization (performance)
	"KTN-VAR-018": SEVERITY_WARNING, // Constants for fixed values (maintenabilité)
	"KTN-VAR-019": SEVERITY_INFO,    // Variable grouping by purpose (organisation)

	// FUNC - Fonctions
	"KTN-FUNC-001": SEVERITY_WARNING, // Fonction trop longue (maintenabilité)
	"KTN-FUNC-002": SEVERITY_WARNING, // Trop de paramètres (complexité)
	"KTN-FUNC-003": SEVERITY_INFO,    // Extraire constantes (style léger)
	"KTN-FUNC-004": SEVERITY_INFO,    // Nommage fonctions (style)
	"KTN-FUNC-005": SEVERITY_INFO,    // Éviter effets de bord (architecture légère)
	"KTN-FUNC-006": SEVERITY_ERROR,   // Erreur pas en dernière position (convention Go critique)
	"KTN-FUNC-007": SEVERITY_WARNING, // Documentation manquante (maintenabilité)
	"KTN-FUNC-008": SEVERITY_ERROR,   // Context pas en premier (convention Go critique)
	"KTN-FUNC-009": SEVERITY_WARNING, // Side effects dans getters (architecture)
	"KTN-FUNC-010": SEVERITY_INFO,    // Named returns (style)
	"KTN-FUNC-011": SEVERITY_WARNING, // Commentaires manquants (maintenabilité)
	"KTN-FUNC-012": SEVERITY_ERROR,   // else après return (dead code potentiel)
	"KTN-FUNC-013": SEVERITY_WARNING, // Paramètres non utilisés (clarté)
	"KTN-FUNC-014": SEVERITY_ERROR,   // Fonctions privées non utilisées (code mort)

	// STRUCT - Structures
	"KTN-STRUCT-001": SEVERITY_INFO,    // Un fichier par struct (organisation)
	"KTN-STRUCT-002": SEVERITY_WARNING, // Interface manquante (architecture)
	"KTN-STRUCT-003": SEVERITY_INFO,    // Ordre des champs (style)
	"KTN-STRUCT-004": SEVERITY_WARNING, // Documentation manquante (maintenabilité)
	"KTN-STRUCT-005": SEVERITY_WARNING, // Constructeur manquant (architecture)
	"KTN-STRUCT-006": SEVERITY_WARNING, // Encapsulation manquante (architecture)

	// TEST - Tests
	"KTN-TEST-001": SEVERITY_WARNING, // Package xxx_test (isolation)
	"KTN-TEST-002": SEVERITY_WARNING, // Nommage tests (convention)
	"KTN-TEST-003": SEVERITY_WARNING, // Tests manquants (qualité)
	"KTN-TEST-004": SEVERITY_INFO,    // Coverage cas d'erreur (qualité)
	"KTN-TEST-005": SEVERITY_WARNING, // Table-driven tests (maintenabilité)
	"KTN-TEST-006": SEVERITY_WARNING, // Test helpers (organisation)
	"KTN-TEST-007": SEVERITY_WARNING, // t.Skip() interdit (qualité)
	"KTN-TEST-008": SEVERITY_WARNING, // Convention _internal/_external_test.go (organisation)
	"KTN-TEST-009": SEVERITY_WARNING, // Tests publics dans _external uniquement (organisation)
	"KTN-TEST-010": SEVERITY_WARNING, // Tests privés dans _internal uniquement (organisation)
	"KTN-TEST-011": SEVERITY_WARNING, // Package xxx vs xxx_test (isolation)
	"KTN-TEST-012": SEVERITY_ERROR,   // Fichiers *_test.go interdits (convention stricte)

	// PACKAGE - Documentation de fichiers
	"KTN-PACKAGE-001": SEVERITY_WARNING, // Description de fichier manquante (maintenabilité)
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
	// Par défaut SEVERITY_WARNING
	return SEVERITY_WARNING
}

// ColorCode retourne le code couleur ANSI pour un niveau.
//
// Returns:
//   - string: code couleur ANSI
func (l Level) ColorCode() string {
	// Vérification du niveau
	switch l {
	// Niveau SEVERITY_INFO
	case SEVERITY_INFO:
		// Bleu
		return "\033[34m"
	// Niveau SEVERITY_WARNING
	case SEVERITY_WARNING:
		// Orange (jaune vif)
		return "\033[33m"
	// Niveau SEVERITY_ERROR
	case SEVERITY_ERROR:
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
// Returns:
//   - string: symbole
func (l Level) Symbol() string {
	// Vérification du niveau
	switch l {
	// Niveau SEVERITY_INFO
	case SEVERITY_INFO:
		// Symbole info
		return "ℹ"
	// Niveau SEVERITY_WARNING
	case SEVERITY_WARNING:
		// Symbole warning
		return "⚠"
	// Niveau SEVERITY_ERROR
	case SEVERITY_ERROR:
		// Symbole error
		return "✖"
	// Niveau par défaut
	default:
		// Symbole générique
		return "●"
	}
}
