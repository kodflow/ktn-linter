// Severity levels for linter diagnostics.
package severity

const (
	// SeverityInfo recommandations et style
	SeverityInfo Level = iota
	// SeverityWarning problèmes de maintenabilité
	SeverityWarning
	// SeverityError bugs potentiels et problèmes graves
	SeverityError
)

// Level représente le niveau de sévérité d'une règle
type Level int

// String retourne la représentation textuelle du niveau.
//
// Returns:
//   - string: représentation textuelle
func (l Level) String() string {
	// Vérification du niveau
	switch l {
	// Niveau SeverityInfo
	case SeverityInfo:
		// Retour SeverityInfo
		return "INFO"
	// Niveau SeverityWarning
	case SeverityWarning:
		// Retour SeverityWarning
		return "WARNING"
	// Niveau SeverityError
	case SeverityError:
		// Retour SeverityError
		return "ERROR"
	// Niveau par défaut
	default:
		// Retour Unknown
		return "UNKNOWN"
	}
}

// rulesSeverity mappe chaque règle à son niveau de sévérité
// Les règles sont ordonnées par criticité : ERROR > WARNING > INFO
var rulesSeverity map[string]Level = map[string]Level{
	// COMMENT - Commentaires et documentation (7 règles)
	"KTN-COMMENT-001": SeverityInfo,    // Commentaire inline trop long >80 chars
	"KTN-COMMENT-002": SeverityWarning, // Description de fichier manquante (ex-PACKAGE-001)
	"KTN-COMMENT-003": SeverityWarning, // Commentaire requis pour constantes (ex-CONST-002)
	"KTN-COMMENT-004": SeverityWarning, // Commentaire requis pour var package (ex-VAR-002)
	"KTN-COMMENT-005": SeverityWarning, // Documentation struct exportée (ex-STRUCT-002)
	"KTN-COMMENT-006": SeverityWarning, // Documentation fonction complète (ex-FUNC-007)
	"KTN-COMMENT-007": SeverityWarning, // Commentaires sur blocs de contrôle (ex-FUNC-009)

	// CONST - Constantes (3 règles, ex-002 déplacé vers COMMENT-003)
	"KTN-CONST-001": SeverityWarning, // Type explicite requis
	"KTN-CONST-002": SeverityInfo,    // Groupement des constantes (ex-003)
	"KTN-CONST-003": SeverityInfo,    // Convention SCREAMING_SNAKE_CASE (ex-004)

	// VAR - Variables (17 règles, ex-002 déplacé vers COMMENT-004)
	"KTN-VAR-001": SeverityError,   // Variables package en camelCase
	"KTN-VAR-002": SeverityWarning, // Type explicite pour var package (ex-003)
	"KTN-VAR-003": SeverityWarning, // Utiliser := pour variables locales (ex-004)
	"KTN-VAR-004": SeverityWarning, // Préallocation slices (ex-005)
	"KTN-VAR-005": SeverityWarning, // Éviter make([]T, length) avec append (ex-006)
	"KTN-VAR-006": SeverityWarning, // Préallocation bytes.Buffer/strings.Builder (ex-007)
	"KTN-VAR-007": SeverityWarning, // Utiliser strings.Builder pour >2 concat (ex-008)
	"KTN-VAR-008": SeverityWarning, // Éviter allocs dans boucles chaudes (ex-009)
	"KTN-VAR-009": SeverityWarning, // Pointeurs pour structs >64 bytes (ex-010)
	"KTN-VAR-010": SeverityWarning, // sync.Pool pour buffers répétés (ex-011)
	"KTN-VAR-011": SeverityWarning, // Shadowing de variables (ex-012)
	"KTN-VAR-012": SeverityWarning, // Conversions string() répétées (ex-013)
	"KTN-VAR-013": SeverityInfo,    // Groupement des variables (ex-014)
	"KTN-VAR-014": SeverityInfo,    // Variables après constantes (ex-015)
	"KTN-VAR-015": SeverityInfo,    // Préallocation maps (ex-016)
	"KTN-VAR-016": SeverityInfo,    // Utiliser [N]T au lieu de make([]T, N) (ex-017)
	"KTN-VAR-017": SeverityInfo,    // Copies de mutex (ex-018)

	// FUNC - Fonctions (12 règles, ex-007 et ex-009 déplacés vers COMMENT)
	"KTN-FUNC-001": SeverityError,   // Erreur en dernière position retour
	"KTN-FUNC-002": SeverityError,   // context.Context en premier paramètre
	"KTN-FUNC-003": SeverityError,   // Éviter else après return
	"KTN-FUNC-004": SeverityError,   // Fonctions privées non utilisées
	"KTN-FUNC-005": SeverityWarning, // Max 35 lignes par fonction
	"KTN-FUNC-006": SeverityWarning, // Max 5 paramètres
	"KTN-FUNC-007": SeverityWarning, // Getters sans side effects (ex-008)
	"KTN-FUNC-008": SeverityWarning, // Paramètres non utilisés préfixés _ (ex-010)
	"KTN-FUNC-009": SeverityInfo,    // Pas de magic numbers (ex-011)
	"KTN-FUNC-010": SeverityInfo,    // Naked returns interdits (ex-012)
	"KTN-FUNC-011": SeverityInfo,    // Complexité cyclomatique ≤10 (ex-013)
	"KTN-FUNC-012": SeverityInfo,    // Named returns si >3 valeurs retour (ex-014)

	// INTERFACE - Interfaces (2 règles)
	"KTN-INTERFACE-001": SeverityWarning, // Interface privée non utilisée
	"KTN-INTERFACE-002": SeverityInfo,    // Interface exportée design-first

	// STRUCT - Structures (5 règles, ex-002 déplacé vers COMMENT-005)
	"KTN-STRUCT-001": SeverityWarning, // Interface pour chaque struct
	"KTN-STRUCT-002": SeverityWarning, // Constructeur NewX() requis (ex-003)
	"KTN-STRUCT-003": SeverityWarning, // Pas de préfixe Get pour getters (ex-004)
	"KTN-STRUCT-004": SeverityInfo,    // Une struct par fichier (ex-005)
	"KTN-STRUCT-005": SeverityInfo,    // Champs exportés avant privés (ex-006)

	// TEST - Tests (13 règles)
	"KTN-TEST-001": SeverityError,   // Fichier doit finir par _internal/_external_test.go
	"KTN-TEST-002": SeverityWarning, // Package xxx_test pour tests (désactivé)
	"KTN-TEST-003": SeverityWarning, // Fichier .go correspondant requis
	"KTN-TEST-004": SeverityWarning, // Toutes fonctions doivent avoir tests
	"KTN-TEST-005": SeverityWarning, // Table-driven tests requis
	"KTN-TEST-006": SeverityWarning, // Pattern 1:1 fichiers test/source
	"KTN-TEST-007": SeverityWarning, // Interdiction t.Skip()
	"KTN-TEST-008": SeverityWarning, // Fichiers test appropriés requis
	"KTN-TEST-009": SeverityWarning, // Tests publics dans _external_test.go
	"KTN-TEST-010": SeverityWarning, // Tests privés dans _internal_test.go
	"KTN-TEST-011": SeverityWarning, // Package correct selon type fichier
	"KTN-TEST-012": SeverityWarning, // Tests doivent contenir assertions
	"KTN-TEST-013": SeverityInfo,    // Couverture cas d'erreur
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
	// Par défaut SeverityWarning
	return SeverityWarning
}

// ColorCode retourne le code couleur ANSI pour un niveau.
//
// Returns:
//   - string: code couleur ANSI
func (l Level) ColorCode() string {
	// Vérification du niveau
	switch l {
	// Niveau SeverityInfo
	case SeverityInfo:
		// Bleu
		return "\033[34m"
	// Niveau SeverityWarning
	case SeverityWarning:
		// Orange (jaune vif)
		return "\033[33m"
	// Niveau SeverityError
	case SeverityError:
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
	// Niveau SeverityInfo
	case SeverityInfo:
		// Symbole info
		return "ℹ"
	// Niveau SeverityWarning
	case SeverityWarning:
		// Symbole warning
		return "⚠"
	// Niveau SeverityError
	case SeverityError:
		// Symbole error
		return "✖"
	// Niveau par défaut
	default:
		// Symbole générique
		return "●"
	}
}
