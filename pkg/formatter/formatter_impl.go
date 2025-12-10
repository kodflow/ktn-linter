// Implementation of the formatter interface.
package formatter

import (
	"fmt"
	"go/token"
	"io"
	"sort"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/severity"
	"golang.org/x/tools/go/analysis"
)

// Formatter définit l'interface pour formater et afficher les diagnostics.
type Formatter interface {
	// Format affiche les diagnostics de manière lisible
	//
	// Params:
	//   - fset: le FileSet contenant les informations de position
	//   - diagnostics: la liste des diagnostics à formater
	Format(fset *token.FileSet, diagnostics []analysis.Diagnostic)
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
//   - fset: ensemble de fichiers
//   - diagnostics: liste des diagnostics
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
//   - fset: ensemble de fichiers
//   - diagnostics: liste des diagnostics
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
//   - fset: ensemble de fichiers
//   - diagnostics: liste des diagnostics
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
//   - fset: ensemble de fichiers
//   - diagnostics: liste des diagnostics
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
//   - fset: ensemble de fichiers
//   - diagnostics: liste des diagnostics
//
// Returns:
//   - []DiagnosticGroupData: groupes de diagnostics
func (f *formatterImpl) groupByFile(fset *token.FileSet, diagnostics []analysis.Diagnostic) []DiagnosticGroupData {
	fileMap := make(map[string][]analysis.Diagnostic, InitialFileMapCap)

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
//   - fset: ensemble de fichiers
//   - diagnostics: liste des diagnostics
//
// Returns:
//   - []analysis.Diagnostic: diagnostics filtrés
func (f *formatterImpl) filterAndSortDiagnostics(fset *token.FileSet, diagnostics []analysis.Diagnostic) []analysis.Diagnostic {
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
//   - count: nombre total de diagnostics
func (f *formatterImpl) printHeader(count int) {
	// Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "\n╔════════════════════════════════════════════════════════════╗\n")
		fmt.Fprintf(f.writer, "║  KTN-LINTER REPORT - %d issue(s) found                     ║\n", count)
		fmt.Fprintf(f.writer, "╚════════════════════════════════════════════════════════════╝\n\n")
		// Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "\n%s%s╔════════════════════════════════════════════════════════════╗%s\n", Bold, Blue, Reset)
		fmt.Fprintf(f.writer, "%s%s║  KTN-LINTER REPORT - %d issue(s) found                     ║%s\n", Bold, Blue, count, Reset)
		fmt.Fprintf(f.writer, "%s%s╚════════════════════════════════════════════════════════════╝%s\n\n", Bold, Blue, Reset)
	}
}

// printFileHeader affiche l'en-tête pour un fichier
// Params:
//   - filename: nom du fichier
//   - count: nombre de diagnostics
func (f *formatterImpl) printFileHeader(filename string, count int) {
	// Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "📁 File: %s (%d issues)\n", filename, count)
		fmt.Fprintf(f.writer, "────────────────────────────────────────────────────────────\n")
		// Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "%s📁 File: %s%s %s(%d issues)%s\n",
			Bold+Cyan, filename, Reset, Gray, count, Reset)
		fmt.Fprintf(f.writer, "%s────────────────────────────────────────────────────────────%s\n", Gray, Reset)
	}
}

// printDiagnostic affiche un diagnostic individuel
// Params:
//   - num: numéro du diagnostic
//   - pos: position du diagnostic
//   - diag: diagnostic à afficher
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
		symbol := f.getSymbol(code)
		fmt.Fprintf(f.writer, "\n%s[%d]%s %s%s%s\n",
			Bold+Yellow, num, Reset,
			Cyan, location, Reset)
		fmt.Fprintf(f.writer, "  %s%s%s %sCode:%s %s%s%s\n",
			codeColor, symbol, Reset,
			Gray, Reset,
			Bold, code, Reset)
		fmt.Fprintf(f.writer, "  %s▶%s %s\n",
			Blue, Reset, message)
	}
}

// printSuccess affiche un message de succès
func (f *formatterImpl) printSuccess() {
	// Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "\n✅ No issues found! Code is compliant.\n\n")
		// Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "\n%s✅ No issues found! Code is compliant.%s\n\n", Bold+Green, Reset)
	}
}

// printSummary affiche le résumé final
// Params:
//   - count: nombre total de diagnostics
func (f *formatterImpl) printSummary(count int) {
	// Vérification de la condition
	if f.noColor {
		fmt.Fprintf(f.writer, "════════════════════════════════════════════════════════════\n")
		fmt.Fprintf(f.writer, "Total: %d issue(s) to fix\n", count)
		fmt.Fprintf(f.writer, "════════════════════════════════════════════════════════════\n\n")
		// Cas alternatif
	} else {
		fmt.Fprintf(f.writer, "%s════════════════════════════════════════════════════════════%s\n", Gray, Reset)
		fmt.Fprintf(f.writer, "%s📊 Total: %s%d%s issue(s) to fix (WARNING)\n",
			Bold, Yellow, count, Reset)
		fmt.Fprintf(f.writer, "%s════════════════════════════════════════════════════════════%s\n\n", Gray, Reset)
	}
}

// getCodeColor retourne la couleur ANSI appropriée pour un code d'erreur basée sur la sévérité
// Params:
//   - code: code d'erreur (ex: "KTN-VAR-001")
//
// Returns:
//   - string: couleur ANSI selon la sévérité (rouge/orange/bleu)
func (f *formatterImpl) getCodeColor(code string) string {
	// Vérification de la condition
	if f.noColor {
		// Early return from function.
		return ""
	}

	// Obtenir la sévérité du code
	level := severity.GetSeverity(code)
	// Retour de la couleur selon le niveau
	return level.ColorCode()
}

// getSymbol retourne le symbole approprié pour un code d'erreur basé sur la sévérité
// Params:
//   - code: code d'erreur (ex: "KTN-VAR-001")
//
// Returns:
//   - string: symbole (✖/⚠/ℹ)
func (f *formatterImpl) getSymbol(code string) string {
	// Obtenir la sévérité du code
	level := severity.GetSeverity(code)
	// Retour du symbole selon le niveau
	return level.Symbol()
}
