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
// Les règles sont ordonnées par criticité : ERROR > WARNING > INFO
var rulesSeverity map[string]Level = map[string]Level{
	// COMMENT - Commentaires et documentation (7 règles)
	"KTN-COMMENT-001": SEVERITY_INFO,    // Commentaire inline trop long >80 chars
	"KTN-COMMENT-002": SEVERITY_WARNING, // Description de fichier manquante (ex-PACKAGE-001)
	"KTN-COMMENT-003": SEVERITY_WARNING, // Commentaire requis pour constantes (ex-CONST-002)
	"KTN-COMMENT-004": SEVERITY_WARNING, // Commentaire requis pour var package (ex-VAR-002)
	"KTN-COMMENT-005": SEVERITY_WARNING, // Documentation struct exportée (ex-STRUCT-002)
	"KTN-COMMENT-006": SEVERITY_WARNING, // Documentation fonction complète (ex-FUNC-007)
	"KTN-COMMENT-007": SEVERITY_WARNING, // Commentaires sur blocs de contrôle (ex-FUNC-009)

	// CONST - Constantes (3 règles, ex-002 déplacé vers COMMENT-003)
	"KTN-CONST-001": SEVERITY_WARNING, // Type explicite requis
	"KTN-CONST-002": SEVERITY_INFO,    // Groupement des constantes (ex-003)
	"KTN-CONST-003": SEVERITY_INFO,    // Convention SCREAMING_SNAKE_CASE (ex-004)

	// VAR - Variables (17 règles, ex-002 déplacé vers COMMENT-004)
	"KTN-VAR-001": SEVERITY_ERROR,   // Variables package en camelCase
	"KTN-VAR-002": SEVERITY_WARNING, // Type explicite pour var package (ex-003)
	"KTN-VAR-003": SEVERITY_WARNING, // Utiliser := pour variables locales (ex-004)
	"KTN-VAR-004": SEVERITY_WARNING, // Préallocation slices (ex-005)
	"KTN-VAR-005": SEVERITY_WARNING, // Éviter make([]T, length) avec append (ex-006)
	"KTN-VAR-006": SEVERITY_WARNING, // Préallocation bytes.Buffer/strings.Builder (ex-007)
	"KTN-VAR-007": SEVERITY_WARNING, // Utiliser strings.Builder pour >2 concat (ex-008)
	"KTN-VAR-008": SEVERITY_WARNING, // Éviter allocs dans boucles chaudes (ex-009)
	"KTN-VAR-009": SEVERITY_WARNING, // Pointeurs pour structs >64 bytes (ex-010)
	"KTN-VAR-010": SEVERITY_WARNING, // sync.Pool pour buffers répétés (ex-011)
	"KTN-VAR-011": SEVERITY_WARNING, // Shadowing de variables (ex-012)
	"KTN-VAR-012": SEVERITY_WARNING, // Conversions string() répétées (ex-013)
	"KTN-VAR-013": SEVERITY_INFO,    // Groupement des variables (ex-014)
	"KTN-VAR-014": SEVERITY_INFO,    // Variables après constantes (ex-015)
	"KTN-VAR-015": SEVERITY_INFO,    // Préallocation maps (ex-016)
	"KTN-VAR-016": SEVERITY_INFO,    // Utiliser [N]T au lieu de make([]T, N) (ex-017)
	"KTN-VAR-017": SEVERITY_INFO,    // Copies de mutex (ex-018)

	// FUNC - Fonctions (12 règles, ex-007 et ex-009 déplacés vers COMMENT)
	"KTN-FUNC-001": SEVERITY_ERROR,   // Erreur en dernière position retour
	"KTN-FUNC-002": SEVERITY_ERROR,   // context.Context en premier paramètre
	"KTN-FUNC-003": SEVERITY_ERROR,   // Éviter else après return
	"KTN-FUNC-004": SEVERITY_ERROR,   // Fonctions privées non utilisées
	"KTN-FUNC-005": SEVERITY_WARNING, // Max 35 lignes par fonction
	"KTN-FUNC-006": SEVERITY_WARNING, // Max 5 paramètres
	"KTN-FUNC-007": SEVERITY_WARNING, // Getters sans side effects (ex-008)
	"KTN-FUNC-008": SEVERITY_WARNING, // Paramètres non utilisés préfixés _ (ex-010)
	"KTN-FUNC-009": SEVERITY_INFO,    // Pas de magic numbers (ex-011)
	"KTN-FUNC-010": SEVERITY_INFO,    // Naked returns interdits (ex-012)
	"KTN-FUNC-011": SEVERITY_INFO,    // Complexité cyclomatique ≤10 (ex-013)
	"KTN-FUNC-012": SEVERITY_INFO,    // Named returns si >3 valeurs retour (ex-014)

	// STRUCT - Structures (5 règles, ex-002 déplacé vers COMMENT-005)
	"KTN-STRUCT-001": SEVERITY_WARNING, // Interface pour chaque struct
	"KTN-STRUCT-002": SEVERITY_WARNING, // Constructeur NewX() requis (ex-003)
	"KTN-STRUCT-003": SEVERITY_WARNING, // Pas de préfixe Get pour getters (ex-004)
	"KTN-STRUCT-004": SEVERITY_INFO,    // Une struct par fichier (ex-005)
	"KTN-STRUCT-005": SEVERITY_INFO,    // Champs exportés avant privés (ex-006)

	// TEST - Tests (13 règles)
	"KTN-TEST-001": SEVERITY_ERROR,   // Fichier doit finir par _internal/_external_test.go
	"KTN-TEST-002": SEVERITY_WARNING, // Package xxx_test pour tests (désactivé)
	"KTN-TEST-003": SEVERITY_WARNING, // Fichier .go correspondant requis
	"KTN-TEST-004": SEVERITY_WARNING, // Toutes fonctions doivent avoir tests
	"KTN-TEST-005": SEVERITY_WARNING, // Table-driven tests requis
	"KTN-TEST-006": SEVERITY_WARNING, // Pattern 1:1 fichiers test/source
	"KTN-TEST-007": SEVERITY_WARNING, // Interdiction t.Skip()
	"KTN-TEST-008": SEVERITY_WARNING, // Fichiers test appropriés requis
	"KTN-TEST-009": SEVERITY_WARNING, // Tests publics dans _external_test.go
	"KTN-TEST-010": SEVERITY_WARNING, // Tests privés dans _internal_test.go
	"KTN-TEST-011": SEVERITY_WARNING, // Package correct selon type fichier
	"KTN-TEST-012": SEVERITY_WARNING, // Tests doivent contenir assertions
	"KTN-TEST-013": SEVERITY_INFO,    // Couverture cas d'erreur
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
