// Formatter interface for linter output formatting.
package formatter

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Codes de couleurs ANSI pour le formatage terminal
const (
	// Red représente le code ANSI pour la couleur rouge
	Red string = "\033[31m"
	// Green représente le code ANSI pour la couleur verte
	Green string = "\033[32m"
	// Yellow représente le code ANSI pour la couleur jaune
	Yellow string = "\033[33m"
	// Blue représente le code ANSI pour la couleur bleue
	Blue string = "\033[34m"
	// Magenta représente le code ANSI pour la couleur magenta
	Magenta string = "\033[35m"
	// Cyan représente le code ANSI pour la couleur cyan
	Cyan string = "\033[36m"
	// Gray représente le code ANSI pour la couleur grise
	Gray string = "\033[90m"
	// Bold représente le code ANSI pour le texte en gras
	Bold string = "\033[1m"
	// Reset représente le code ANSI pour réinitialiser le formatage
	Reset string = "\033[0m"

	// InitialFileMapCap définit la capacité estimée pour groupement par fichier
	InitialFileMapCap int = 16
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

// extractMessage extrait le message principal en supprimant le code et les exemples.
//
// Params:
//   - message: message brut du diagnostic
//
// Returns:
//   - string: message nettoyé sans le code KTN
func extractMessage(message string) string {
	// Extraction avec troncature (mode normal)
	return extractMessageWithOptions(message, true)
}

// extractMessageWithOptions extrait le message avec options.
//
// Params:
//   - message: message brut du diagnostic
//   - truncate: true pour tronquer au premier \n
//
// Returns:
//   - string: message nettoyé sans le code KTN
func extractMessageWithOptions(message string, truncate bool) string {
	var idx int
	// Supprimer le code [KTN-XXX-XXX] ou KTN-XXX-XXX:
	// Traitement Format 1: [KTN-XXX-XXX] ...
	if strings.HasPrefix(message, "[KTN-") {
		// Chercher le ] correspondant au code KTN
		if idx = strings.Index(message, "]"); idx != -1 && idx < len(message)-1 {
			message = strings.TrimSpace(message[idx+1:])
		}
		// Sinon traitement Format 2: KTN-XXX-XXX: ...
	} else if strings.HasPrefix(message, "KTN-") {
		// Chercher le : séparateur
		if idx = strings.Index(message, ":"); idx != -1 && idx < len(message)-1 {
			message = strings.TrimSpace(message[idx+1:])
		}
	}

	// Tronquer au premier \n si demandé (mode normal)
	if truncate {
		// Chercher le premier \n
		if idx = strings.Index(message, "\n"); idx != -1 {
			message = message[:idx]
		}
	}

	// Retour du message nettoyé
	return message
}
