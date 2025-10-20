package formatter

import (
	"fmt"
	"go/token"
	"io"
	"sort"
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
)

// DiagnosticGroupData regroupe les diagnostics par fichier
type DiagnosticGroupData struct {
	Filename    string
	Diagnostics []analysis.Diagnostic
}

// formatterImpl implémente l'interface Formatter
type formatterImpl struct {
	writer     io.Writer
	aiMode     bool
	noColor    bool
	simpleMode bool
}

// NewFormatter crée un nouveau formatter avec les options spécifiées.
//
// Params:
//   - w: le writer où écrire la sortie
//   - aiMode: true pour activer le mode optimisé pour l'IA
//   - noColor: true pour désactiver les couleurs
//   - simpleMode: true pour activer le format simple une-ligne
//
// Returns:
//   - Formatter: un formatter prêt à utiliser
func NewFormatter(w io.Writer, aiMode bool, noColor bool, simpleMode bool) Formatter {
	// Early return from function.
	return &formatterImpl{
		writer:     w,
		aiMode:     aiMode,
		noColor:    noColor,
		simpleMode: simpleMode,
	}
}

// Format affiche les diagnostics de manière lisible
// Params:
//   - pass: contexte d'analyse
//
func (f *formatterImpl) Format(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
 // Vérification de la condition
	if len(diagnostics) == 0 {
		f.printSuccess()
		// Early return from function.
		return
	}

 // Vérification de la condition
	if f.simpleMode {
		f.formatSimple(fset, diagnostics)
		// Early return from function.
		return
	}

 // Vérification de la condition
	if f.aiMode {
		f.formatForAI(fset, diagnostics)
		// Early return from function.
		return
	}

	f.formatForHuman(fset, diagnostics)
}

// formatForHuman affiche pour un humain avec couleurs et structure
// Params:
//   - pass: contexte d'analyse
//
func (f *formatterImpl) formatForHuman(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	groups := f.groupByFile(fset, diagnostics)

	totalCount := 0
 // Itération sur les éléments
	for _, group := range groups {
		totalCount += len(group.Diagnostics)
	}

 // Vérification de la condition
	if totalCount == 0 {
		f.printSuccess()
		// Early return from function.
		return
	}

	f.printHeader(totalCount)

 // Itération sur les éléments
	for _, group := range groups {
		f.printFileHeader(group.Filename, len(group.Diagnostics))

  // Itération sur les éléments
		for i, diag := range group.Diagnostics {
			pos := fset.Position(diag.Pos)
			f.printDiagnostic(i+1, pos, diag)
		}

		fmt.Fprintln(f.writer)
	}

	f.printSummary(totalCount)
}

// formatForAI affiche un format optimisé pour l'IA
// Params:
//   - pass: contexte d'analyse
//
func (f *formatterImpl) formatForAI(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	groups := f.groupByFile(fset, diagnostics)

	totalCount := 0
 // Itération sur les éléments
	for _, group := range groups {
		totalCount += len(group.Diagnostics)
	}

	fmt.Fprintf(f.writer, "# KTN-Linter Report (AI Mode)\n\n")
	fmt.Fprintf(f.writer, "Total issues found: %d\n\n", totalCount)

 // Itération sur les éléments
	for _, group := range groups {
		fmt.Fprintf(f.writer, "## File: %s (%d issues)\n\n", group.Filename, len(group.Diagnostics))

  // Itération sur les éléments
		for _, diag := range group.Diagnostics {
			pos := fset.Position(diag.Pos)
			code := extractCode(diag.Message)

			fmt.Fprintf(f.writer, "### Issue at line %d, column %d\n", pos.Line, pos.Column)
			fmt.Fprintf(f.writer, "- **Code**: %s\n", code)
			fmt.Fprintf(f.writer, "- **Message**: %s\n", diag.Message)
			fmt.Fprintf(f.writer, "- **Category**: %s\n", diag.Category)
			fmt.Fprintln(f.writer)
		}
	}
}

// formatSimple affiche un format simple une ligne par erreur (pour IDE)
// Params:
//   - pass: contexte d'analyse
//
// Format compatible golangci-lint: file:line:col: message (code)
func (f *formatterImpl) formatSimple(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	filtered := f.filterAndSortDiagnostics(fset, diagnostics)

 // Itération sur les éléments
	for _, diag := range filtered {
		pos := fset.Position(diag.Pos)
		code := extractCode(diag.Message)
		message := extractMessage(diag.Message)

		// Format compatible avec golangci-lint et VSCode : code en premier
		fmt.Fprintf(f.writer, "%s:%d:%d: [%s] %s\n",
			pos.Filename, pos.Line, pos.Column, code, message)
	}
}

// groupByFile regroupe les diagnostics par fichier et les trie
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - []DiagnosticGroupData: groupes de diagnostics
//
func (f *formatterImpl) groupByFile(fset *token.FileSet, diagnostics []analysis.Diagnostic) []DiagnosticGroupData {
	fileMap := make(map[string][]analysis.Diagnostic)

 // Itération sur les éléments
	for _, diag := range diagnostics {
		pos := fset.Position(diag.Pos)
		filename := pos.Filename

		// Ignorer les fichiers du cache Go et les fichiers temporaires
		if strings.Contains(filename, "/.cache/go-build/") ||
			strings.Contains(filename, "/tmp/") ||
			strings.Contains(filename, "\\cache\\go-build\\") {
			continue
		}

		fileMap[filename] = append(fileMap[filename], diag)
	}

	// groups holds the configuration value.

	var groups []DiagnosticGroupData
 // Itération sur les éléments
	for filename, diags := range fileMap {
		// Trier par ligne
		sort.Slice(diags, func(i, j int) bool {
			// Early return from function.
			return fset.Position(diags[i].Pos).Line < fset.Position(diags[j].Pos).Line
		})
		groups = append(groups, DiagnosticGroupData{
			Filename:    filename,
			Diagnostics: diags,
		})
	}

	// Trier par nom de fichier
	sort.Slice(groups, func(i, j int) bool {
		// Early return from function.
		return groups[i].Filename < groups[j].Filename
	})

	// Early return from function.
	return groups
}

// filterAndSortDiagnostics filtre et trie les diagnostics par position
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - []analysis.Diagnostic: diagnostics filtrés
//
func (f *formatterImpl) filterAndSortDiagnostics(fset *token.FileSet, diagnostics []analysis.Diagnostic) []analysis.Diagnostic {
	// filtered holds the configuration value.

	var filtered []analysis.Diagnostic
 // Itération sur les éléments
	for _, diag := range diagnostics {
		pos := fset.Position(diag.Pos)
		// Ignorer les fichiers du cache Go et les fichiers temporaires
		if strings.Contains(pos.Filename, "/.cache/go-build/") ||
			strings.Contains(pos.Filename, "/tmp/") ||
			strings.Contains(pos.Filename, "\\cache\\go-build\\") {
			continue
		}
		filtered = append(filtered, diag)
	}

	sort.Slice(filtered, func(i, j int) bool {
		posI := fset.Position(filtered[i].Pos)
		posJ := fset.Position(filtered[j].Pos)
  // Vérification de la condition
		if posI.Filename != posJ.Filename {
			// Early return from function.
			return posI.Filename < posJ.Filename
		}
  // Vérification de la condition
		if posI.Line != posJ.Line {
			// Early return from function.
			return posI.Line < posJ.Line
		}
		// Early return from function.
		return posI.Column < posJ.Column
	})

	// Early return from function.
	return filtered
}

// printHeader affiche l'en-tête du rapport
// Params:
//   - pass: contexte d'analyse
//
func (f *formatterImpl) printHeader(count int) {
 // Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "\n╔════════════════════════════════════════════════════════════╗\n")
		fmt.Fprintf(f.writer, "║  KTN-LINTER REPORT - %d issue(s) found                     ║\n", count)
		fmt.Fprintf(f.writer, "╚════════════════════════════════════════════════════════════╝\n\n")
 // Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "\n%s%s╔════════════════════════════════════════════════════════════╗%s\n", BOLD, BLUE, RESET)
		fmt.Fprintf(f.writer, "%s%s║  KTN-LINTER REPORT - %d issue(s) found                     ║%s\n", BOLD, BLUE, count, RESET)
		fmt.Fprintf(f.writer, "%s%s╚════════════════════════════════════════════════════════════╝%s\n\n", BOLD, BLUE, RESET)
	}
}

// printFileHeader affiche l'en-tête pour un fichier
// Params:
//   - pass: contexte d'analyse
//
func (f *formatterImpl) printFileHeader(filename string, count int) {
 // Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "📁 File: %s (%d issues)\n", filename, count)
		fmt.Fprintf(f.writer, "────────────────────────────────────────────────────────────\n")
 // Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "%s📁 File: %s%s %s(%d issues)%s\n",
			BOLD+CYAN, filename, RESET, GRAY, count, RESET)
		fmt.Fprintf(f.writer, "%s────────────────────────────────────────────────────────────%s\n", GRAY, RESET)
	}
}

// printDiagnostic affiche un diagnostic individuel
// Params:
//   - pass: contexte d'analyse
//
func (f *formatterImpl) printDiagnostic(num int, pos token.Position, diag analysis.Diagnostic) {
	code := extractCode(diag.Message)
	message := extractMessage(diag.Message)
	location := fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)

 // Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "\n[%d] %s\n", num, location)
		fmt.Fprintf(f.writer, "  Code: %s\n", code)
		fmt.Fprintf(f.writer, "  Issue: %s\n", message)
 // Cas alternatif
	} else {
		codeColor := f.getCodeColor(code)
		fmt.Fprintf(f.writer, "\n%s[%d]%s %s%s%s\n",
			BOLD+YELLOW, num, RESET,
			CYAN, location, RESET)
		fmt.Fprintf(f.writer, "  %s●%s %sCode:%s %s%s%s\n",
			codeColor, RESET,
			GRAY, RESET,
			BOLD, code, RESET)
		fmt.Fprintf(f.writer, "  %s▶%s %s\n",
			BLUE, RESET, message)
	}
}

// printSuccess affiche un message de succès
func (f *formatterImpl) printSuccess() {
 // Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "\n✅ No issues found! Code is compliant.\n\n")
 // Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "\n%s✅ No issues found! Code is compliant.%s\n\n", BOLD+GREEN, RESET)
	}
}

// printSummary affiche le résumé final
// Params:
//   - pass: contexte d'analyse
//
func (f *formatterImpl) printSummary(count int) {
 // Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "════════════════════════════════════════════════════════════\n")
		fmt.Fprintf(f.writer, "Total: %d issue(s) to fix\n", count)
		fmt.Fprintf(f.writer, "════════════════════════════════════════════════════════════\n\n")
 // Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "%s════════════════════════════════════════════════════════════%s\n", GRAY, RESET)
		fmt.Fprintf(f.writer, "%s📊 Total: %s%d%s issue(s) to fix\n",
			BOLD, RED, count, RESET)
		fmt.Fprintf(f.writer, "%s════════════════════════════════════════════════════════════%s\n\n", GRAY, RESET)
	}
}

// getCodeColor retourne la couleur ANSI appropriée pour un code d'erreur
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - string: code extrait
//
func (f *formatterImpl) getCodeColor(code string) string {
 // Vérification de la condition
	if f.noColor {
		// Early return from function.
		return ""
	}

 // Sélection selon la valeur
	switch {
 // Traitement
	case strings.HasSuffix(code, "-001"):
		// Early return from function.
		return RED
 // Traitement
	case strings.HasSuffix(code, "-002"):
		// Early return from function.
		return YELLOW
 // Traitement
	case strings.HasSuffix(code, "-003"):
		// Early return from function.
		return MAGENTA
 // Traitement
	case strings.HasSuffix(code, "-004"):
		// Early return from function.
		return CYAN
 // Traitement
	default:
		// Early return from function.
		return RED
	}
}

// extractCode extrait le code d'erreur du message (ex: "KTN-VAR-001")
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - string: message extrait
//
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
//
func extractMessage(message string) string {
	// Supprimer le code [KTN-XXX-XXX] ou KTN-XXX-XXX:
	// Format 1: [KTN-XXX-XXX] ...
	if idx := strings.Index(message, "]"); idx != -1 && idx < len(message)-1 {
		message = strings.TrimSpace(message[idx+1:])
 // Traitement
 // Vérification de la condition
	} else if strings.HasPrefix(message, "KTN-") {
		// Format 2: KTN-XXX-XXX: ...
		if idx := strings.Index(message, ":"); idx != -1 && idx < len(message)-1 {
			message = strings.TrimSpace(message[idx+1:])
		}
	}

	// Tronquer au premier \n pour avoir juste la première ligne
	if idx := strings.Index(message, "\n"); idx != -1 {
		message = message[:idx]
	}

	// Early return from function.
	return message
}
