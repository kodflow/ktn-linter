package formatter

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Codes de couleurs ANSI pour le formatage terminal
const (
	// RED représente le code ANSI pour la couleur rouge
	RED string = "\033[31m"
	// GREEN représente le code ANSI pour la couleur verte
	GREEN string = "\033[32m"
	// YELLOW représente le code ANSI pour la couleur jaune
	YELLOW string = "\033[33m"
	// BLUE représente le code ANSI pour la couleur bleue
	BLUE string = "\033[34m"
	// MAGENTA représente le code ANSI pour la couleur magenta
	MAGENTA string = "\033[35m"
	// CYAN représente le code ANSI pour la couleur cyan
	CYAN string = "\033[36m"
	// GRAY représente le code ANSI pour la couleur grise
	GRAY string = "\033[90m"
	// BOLD représente le code ANSI pour le texte en gras
	BOLD string = "\033[1m"
	// RESET représente le code ANSI pour réinitialiser le formatage
	RESET string = "\033[0m"

	// INITIAL_FILE_MAP_CAP définit la capacité estimée pour groupement par fichier
	INITIAL_FILE_MAP_CAP int = 16
)

// DiagnosticGroupData regroupe les diagnostics par fichier.
// Structure utilisée pour organiser les violations détectées par fichier lors du formatage de la sortie.
type DiagnosticGroupData struct {
	Filename    string
	Diagnostics []analysis.Diagnostic
}

// extractCode extrait le code d'erreur du message (ex: "KTN-VAR-001")
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - string: message extrait
func extractCode(message string) string {
	// Cherche le pattern KTN-XXX-XXX avec ou sans crochets
	// Format 1: [KTN-XXX-XXX]
	if start := strings.Index(message, "[KTN-"); start != -1 {
		// Vérification de la condition
		if end := strings.Index(message[start:], "]"); end != -1 {
			// Early return from function.
			return message[start+1 : start+end]
		}
	}

	// Format 2: KTN-XXX-XXX: (au début du message)
	if strings.HasPrefix(message, "KTN-") {
		// Vérification de la condition
		if idx := strings.Index(message, ":"); idx != -1 {
			// Early return from function.
			return message[:idx]
		}
	}

	// Early return from function.
	return "UNKNOWN"
}

// extractMessage extrait le message principal en supprimant le code et les exemples
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - string: couleur ANSI
func extractMessage(message string) string {
	var idx int
	// Supprimer le code [KTN-XXX-XXX] ou KTN-XXX-XXX:
	// Format 1: [KTN-XXX-XXX] ...
	if idx = strings.Index(message, "]"); idx != -1 && idx < len(message)-1 {
		message = strings.TrimSpace(message[idx+1:])
		// Traitement
		// Vérification de la condition
	} else if strings.HasPrefix(message, "KTN-") {
		// Format 2: KTN-XXX-XXX: ...
		if idx = strings.Index(message, ":"); idx != -1 && idx < len(message)-1 {
			message = strings.TrimSpace(message[idx+1:])
		}
	}

	// Tronquer au premier \n pour avoir juste la première ligne
	if idx = strings.Index(message, "\n"); idx != -1 {
		message = message[:idx]
	}

	// Early return from function.
	return message
}
