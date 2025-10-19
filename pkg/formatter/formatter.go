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
	Red     string = "\033[31m"
	Green   string = "\033[32m"
	Yellow  string = "\033[33m"
	Blue    string = "\033[34m"
	Magenta string = "\033[35m"
	Cyan    string = "\033[36m"
	Gray    string = "\033[90m"
	Bold    string = "\033[1m"
	Reset   string = "\033[0m"
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
func (f *formatterImpl) Format(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	if len(diagnostics) == 0 {
		f.printSuccess()
		// Early return from function.
		return
	}

	if f.simpleMode {
		f.formatSimple(fset, diagnostics)
		// Early return from function.
		return
	}

	if f.aiMode {
		f.formatForAI(fset, diagnostics)
		// Early return from function.
		return
	}

	f.formatForHuman(fset, diagnostics)
}

// formatForHuman affiche pour un humain avec couleurs et structure
func (f *formatterImpl) formatForHuman(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	groups := f.groupByFile(fset, diagnostics)

	totalCount := 0
	for _, group := range groups {
		totalCount += len(group.Diagnostics)
	}

	if totalCount == 0 {
		f.printSuccess()
		// Early return from function.
		return
	}

	f.printHeader(totalCount)

	for _, group := range groups {
		f.printFileHeader(group.Filename, len(group.Diagnostics))

		for i, diag := range group.Diagnostics {
			pos := fset.Position(diag.Pos)
			f.printDiagnostic(i+1, pos, diag)
		}

		fmt.Fprintln(f.writer)
	}

	f.printSummary(totalCount)
}

// formatForAI affiche un format optimisé pour l'IA
func (f *formatterImpl) formatForAI(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	groups := f.groupByFile(fset, diagnostics)

	totalCount := 0
	for _, group := range groups {
		totalCount += len(group.Diagnostics)
	}

	fmt.Fprintf(f.writer, "# KTN-Linter Report (AI Mode)\n\n")
	fmt.Fprintf(f.writer, "Total issues found: %d\n\n", totalCount)

	for _, group := range groups {
		fmt.Fprintf(f.writer, "## File: %s (%d issues)\n\n", group.Filename, len(group.Diagnostics))

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
// Format compatible golangci-lint: file:line:col: message (code)
func (f *formatterImpl) formatSimple(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	filtered := f.filterAndSortDiagnostics(fset, diagnostics)

	for _, diag := range filtered {
		pos := fset.Position(diag.Pos)
		code := extractCode(diag.Message)
		message := extractMessage(diag.Message)

		// Format compatible avec golangci-lint et VSCode
		fmt.Fprintf(f.writer, "%s:%d:%d: %s (%s)\n",
			pos.Filename, pos.Line, pos.Column, message, code)
	}
}

// groupByFile regroupe les diagnostics par fichier et les trie
func (f *formatterImpl) groupByFile(fset *token.FileSet, diagnostics []analysis.Diagnostic) []DiagnosticGroupData {
	fileMap := make(map[string][]analysis.Diagnostic)

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
func (f *formatterImpl) filterAndSortDiagnostics(fset *token.FileSet, diagnostics []analysis.Diagnostic) []analysis.Diagnostic {
	// filtered holds the configuration value.

	var filtered []analysis.Diagnostic
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
		if posI.Filename != posJ.Filename {
			// Early return from function.
			return posI.Filename < posJ.Filename
		}
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
func (f *formatterImpl) printHeader(count int) {
	if f.noColor {
		fmt.Fprintf(f.writer, "\n╔════════════════════════════════════════════════════════════╗\n")
		fmt.Fprintf(f.writer, "║  KTN-LINTER REPORT - %d issue(s) found                     ║\n", count)
		fmt.Fprintf(f.writer, "╚════════════════════════════════════════════════════════════╝\n\n")
	} else {
		fmt.Fprintf(f.writer, "\n%s%s╔════════════════════════════════════════════════════════════╗%s\n", Bold, Blue, Reset)
		fmt.Fprintf(f.writer, "%s%s║  KTN-LINTER REPORT - %d issue(s) found                     ║%s\n", Bold, Blue, count, Reset)
		fmt.Fprintf(f.writer, "%s%s╚════════════════════════════════════════════════════════════╝%s\n\n", Bold, Blue, Reset)
	}
}

// printFileHeader affiche l'en-tête pour un fichier
func (f *formatterImpl) printFileHeader(filename string, count int) {
	if f.noColor {
		fmt.Fprintf(f.writer, "📁 File: %s (%d issues)\n", filename, count)
		fmt.Fprintf(f.writer, "────────────────────────────────────────────────────────────\n")
	} else {
		fmt.Fprintf(f.writer, "%s📁 File: %s%s %s(%d issues)%s\n",
			Bold+Cyan, filename, Reset, Gray, count, Reset)
		fmt.Fprintf(f.writer, "%s────────────────────────────────────────────────────────────%s\n", Gray, Reset)
	}
}

// printDiagnostic affiche un diagnostic individuel
func (f *formatterImpl) printDiagnostic(num int, pos token.Position, diag analysis.Diagnostic) {
	code := extractCode(diag.Message)
	message := extractMessage(diag.Message)
	location := fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)

	if f.noColor {
		fmt.Fprintf(f.writer, "\n[%d] %s\n", num, location)
		fmt.Fprintf(f.writer, "  Code: %s\n", code)
		fmt.Fprintf(f.writer, "  Issue: %s\n", message)
	} else {
		codeColor := f.getCodeColor(code)
		fmt.Fprintf(f.writer, "\n%s[%d]%s %s%s%s\n",
			Bold+Yellow, num, Reset,
			Cyan, location, Reset)
		fmt.Fprintf(f.writer, "  %s●%s %sCode:%s %s%s%s\n",
			codeColor, Reset,
			Gray, Reset,
			Bold, code, Reset)
		fmt.Fprintf(f.writer, "  %s▶%s %s\n",
			Blue, Reset, message)
	}
}

// printSuccess affiche un message de succès
func (f *formatterImpl) printSuccess() {
	if f.noColor {
		fmt.Fprintf(f.writer, "\n✅ No issues found! Code is compliant.\n\n")
	} else {
		fmt.Fprintf(f.writer, "\n%s✅ No issues found! Code is compliant.%s\n\n", Bold+Green, Reset)
	}
}

// printSummary affiche le résumé final
func (f *formatterImpl) printSummary(count int) {
	if f.noColor {
		fmt.Fprintf(f.writer, "════════════════════════════════════════════════════════════\n")
		fmt.Fprintf(f.writer, "Total: %d issue(s) to fix\n", count)
		fmt.Fprintf(f.writer, "════════════════════════════════════════════════════════════\n\n")
	} else {
		fmt.Fprintf(f.writer, "%s════════════════════════════════════════════════════════════%s\n", Gray, Reset)
		fmt.Fprintf(f.writer, "%s📊 Total: %s%d%s issue(s) to fix\n",
			Bold, Red, count, Reset)
		fmt.Fprintf(f.writer, "%s════════════════════════════════════════════════════════════%s\n\n", Gray, Reset)
	}
}

// getCodeColor retourne la couleur ANSI appropriée pour un code d'erreur
func (f *formatterImpl) getCodeColor(code string) string {
	if f.noColor {
		// Early return from function.
		return ""
	}

	switch {
	case strings.HasSuffix(code, "-001"):
		// Early return from function.
		return Red
	case strings.HasSuffix(code, "-002"):
		// Early return from function.
		return Yellow
	case strings.HasSuffix(code, "-003"):
		// Early return from function.
		return Magenta
	case strings.HasSuffix(code, "-004"):
		// Early return from function.
		return Cyan
	default:
		// Early return from function.
		return Red
	}
}

// extractCode extrait le code d'erreur du message (ex: "KTN-VAR-001")
func extractCode(message string) string {
	// Cherche le pattern KTN-XXX-XXX avec ou sans crochets
	// Format 1: [KTN-XXX-XXX]
	if start := strings.Index(message, "[KTN-"); start != -1 {
		if end := strings.Index(message[start:], "]"); end != -1 {
			// Early return from function.
			return message[start+1 : start+end]
		}
	}

	// Format 2: KTN-XXX-XXX: (au début du message)
	if strings.HasPrefix(message, "KTN-") {
		if idx := strings.Index(message, ":"); idx != -1 {
			// Early return from function.
			return message[:idx]
		}
	}

	// Early return from function.
	return "UNKNOWN"
}

// extractMessage extrait le message principal en supprimant le code et les exemples
func extractMessage(message string) string {
	// Supprimer le code [KTN-XXX-XXX] ou KTN-XXX-XXX:
	// Format 1: [KTN-XXX-XXX] ...
	if idx := strings.Index(message, "]"); idx != -1 && idx < len(message)-1 {
		message = strings.TrimSpace(message[idx+1:])
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
