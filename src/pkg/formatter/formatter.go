package formatter

import (
	"fmt"
	"go/token"
	"io"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/messageutil"
)

// Colors ANSI
const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"
	Bold    = "\033[1m"
	Reset   = "\033[0m"
)

// DiagnosticGroup regroupe les diagnostics par fichier
type DiagnosticGroup struct {
	Filename    string
	Diagnostics []analysis.Diagnostic
}

// Formatter gère le formatage des diagnostics
type Formatter struct {
	writer     io.Writer
	aiMode     bool
	noColor    bool
	simpleMode bool
}

// NewFormatter crée un nouveau formatter
func NewFormatter(w io.Writer, aiMode bool, noColor bool, simpleMode bool) *Formatter {
	return &Formatter{
		writer:     w,
		aiMode:     aiMode,
		noColor:    noColor,
		simpleMode: simpleMode,
	}
}

// Format affiche les diagnostics de manière lisible
func (f *Formatter) Format(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	if len(diagnostics) == 0 {
		f.printSuccess()
		return
	}

	if f.simpleMode {
		f.formatSimple(fset, diagnostics)
		return
	}

	if f.aiMode {
		f.formatForAI(fset, diagnostics)
		return
	}

	f.formatForHuman(fset, diagnostics)
}

// formatForHuman affiche pour un humain avec couleurs et structure
func (f *Formatter) formatForHuman(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	groups := f.groupByFile(fset, diagnostics)

	f.printHeader(len(diagnostics))

	for _, group := range groups {
		f.printFileHeader(group.Filename, len(group.Diagnostics))

		for i, diag := range group.Diagnostics {
			pos := fset.Position(diag.Pos)
			f.printDiagnostic(i+1, pos, diag)
		}

		fmt.Fprintln(f.writer)
	}

	f.printSummary(len(diagnostics))
}

// formatForAI affiche un format optimisé pour l'IA
func (f *Formatter) formatForAI(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	fmt.Fprintf(f.writer, "# KTN-Linter Report (AI Mode)\n\n")
	fmt.Fprintf(f.writer, "Total issues found: %d\n\n", len(diagnostics))

	groups := f.groupByFile(fset, diagnostics)

	for _, group := range groups {
		fmt.Fprintf(f.writer, "## File: %s (%d issues)\n\n", group.Filename, len(group.Diagnostics))

		for _, diag := range group.Diagnostics {
			pos := fset.Position(diag.Pos)
			code := messageutil.ExtractCode(diag.Message)

			fmt.Fprintf(f.writer, "### Issue at line %d, column %d\n", pos.Line, pos.Column)
			fmt.Fprintf(f.writer, "- **Code**: %s\n", code)
			fmt.Fprintf(f.writer, "- **Message**: %s\n", diag.Message)
			fmt.Fprintf(f.writer, "- **Category**: %s\n", diag.Category)

			if suggestion := messageutil.ExtractSuggestion(diag.Message); suggestion != "" {
				fmt.Fprintf(f.writer, "- **Suggestion**:\n```go\n%s\n```\n", suggestion)
			}

			fmt.Fprintln(f.writer)
		}
	}

	fmt.Fprintf(f.writer, "\n---\n")
	fmt.Fprintf(f.writer, "**Instructions for AI**:\n")
	fmt.Fprintf(f.writer, "- Each issue needs to be fixed according to its code and suggestion\n")
	fmt.Fprintf(f.writer, "- Group related constants together in const () blocks\n")
	fmt.Fprintf(f.writer, "- Add documentation comments for all constants and groups\n")
	fmt.Fprintf(f.writer, "- Always specify explicit types for constants\n")
	fmt.Fprintf(f.writer, "- Refer to target.go for examples of correct implementations\n")
}

// formatSimple affiche un format simple une ligne par erreur (pour IDE)
func (f *Formatter) formatSimple(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	// Trier les diagnostics par position
	sorted := make([]analysis.Diagnostic, len(diagnostics))
	copy(sorted, diagnostics)
	sort.Slice(sorted, func(i, j int) bool {
		posI := fset.Position(sorted[i].Pos)
		posJ := fset.Position(sorted[j].Pos)
		if posI.Filename != posJ.Filename {
			return posI.Filename < posJ.Filename
		}
		if posI.Line != posJ.Line {
			return posI.Line < posJ.Line
		}
		return posI.Column < posJ.Column
	})

	// Afficher chaque diagnostic sur une ligne
	for _, diag := range sorted {
		pos := fset.Position(diag.Pos)
		code := messageutil.ExtractCode(diag.Message)
		message := messageutil.ExtractMessage(diag.Message)

		// Format: file:line:column: [CODE] message
		fmt.Fprintf(f.writer, "%s:%d:%d: [%s] %s\n",
			pos.Filename, pos.Line, pos.Column, code, message)
	}
}

func (f *Formatter) groupByFile(fset *token.FileSet, diagnostics []analysis.Diagnostic) []DiagnosticGroup {
	fileMap := make(map[string][]analysis.Diagnostic)

	for _, diag := range diagnostics {
		pos := fset.Position(diag.Pos)
		filename := pos.Filename
		fileMap[filename] = append(fileMap[filename], diag)
	}

	var groups []DiagnosticGroup
	for filename, diags := range fileMap {
		// Trier par ligne
		sort.Slice(diags, func(i, j int) bool {
			return fset.Position(diags[i].Pos).Line < fset.Position(diags[j].Pos).Line
		})
		groups = append(groups, DiagnosticGroup{
			Filename:    filename,
			Diagnostics: diags,
		})
	}

	// Trier par nom de fichier
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Filename < groups[j].Filename
	})

	return groups
}

func (f *Formatter) printHeader(count int) {
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

func (f *Formatter) printFileHeader(filename string, count int) {
	if f.noColor {
		fmt.Fprintf(f.writer, "📁 File: %s (%d issues)\n", filename, count)
		fmt.Fprintf(f.writer, "────────────────────────────────────────────────────────────\n")
	} else {
		fmt.Fprintf(f.writer, "%s📁 File: %s%s %s(%d issues)%s\n",
			Bold+Cyan, filename, Reset, Gray, count, Reset)
		fmt.Fprintf(f.writer, "%s────────────────────────────────────────────────────────────%s\n", Gray, Reset)
	}
}

func (f *Formatter) printDiagnostic(num int, pos token.Position, diag analysis.Diagnostic) {
	code := messageutil.ExtractCode(diag.Message)
	message := messageutil.ExtractMessage(diag.Message)
	suggestion := messageutil.ExtractSuggestion(diag.Message)
	example := f.generateExample(code, message, suggestion)

	// Format cliquable : fichier:ligne:colonne
	location := fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)

	if f.noColor {
		fmt.Fprintf(f.writer, "\n[%d] %s\n", num, location)
		fmt.Fprintf(f.writer, "  Code: %s\n", code)
		fmt.Fprintf(f.writer, "  Issue: %s\n", message)
		if example != "" {
			fmt.Fprintf(f.writer, "\n%s\n", example)
		}
	} else {
		// Numéro et location cliquable
		fmt.Fprintf(f.writer, "\n%s[%d]%s %s%s%s\n",
			Bold+Yellow, num, Reset,
			Cyan, location, Reset)

		// Code d'erreur
		codeColor := f.getCodeColor(code)
		fmt.Fprintf(f.writer, "  %s●%s %sCode:%s %s%s%s\n",
			codeColor, Reset,
			Gray, Reset,
			Bold, code, Reset)

		// Message
		fmt.Fprintf(f.writer, "  %s▶%s %s\n",
			Blue, Reset, message)

		// Exemple avant/après
		if example != "" {
			fmt.Fprintf(f.writer, "\n%s\n", example)
		}
	}
}

func (f *Formatter) printSuccess() {
	if f.noColor {
		fmt.Fprintf(f.writer, "\n✅ No issues found! Code is compliant.\n\n")
	} else {
		fmt.Fprintf(f.writer, "\n%s✅ No issues found! Code is compliant.%s\n\n", Bold+Green, Reset)
	}
}

func (f *Formatter) printSummary(count int) {
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

func (f *Formatter) indentCode(code string, indent string) string {
	lines := strings.Split(code, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, indent+line)
		}
	}
	return strings.Join(result, "\n")
}

func (f *Formatter) getCodeColor(code string) string {
	if f.noColor {
		return ""
	}

	switch {
	case strings.HasSuffix(code, "-001"):
		return Red
	case strings.HasSuffix(code, "-002"):
		return Yellow
	case strings.HasSuffix(code, "-003"):
		return Magenta
	case strings.HasSuffix(code, "-004"):
		return Cyan
	default:
		return Red
	}
}

// generateExample crée un exemple avant/après concret
func (f *Formatter) generateExample(code, message, suggestion string) string {
	var before, after string

	switch code {
	case "KTN-CONST-001":
		// Constante non groupée
		constName := messageutil.ExtractConstName(message)
		constType := messageutil.ExtractType(suggestion)
		before = fmt.Sprintf("const %s %s = ...", constName, constType)
		after = fmt.Sprintf(`const (
    %s %s = ...
)`, constName, constType)

	case "KTN-CONST-002":
		// Groupe sans commentaire
		before = `const (
    MaxValue int = 100
)`
		after = `// Configuration constants
// Define application limits
const (
    MaxValue int = 100
)`

	case "KTN-CONST-003":
		// Constante sans commentaire individuel
		constName := messageutil.ExtractConstName(message)
		constType := messageutil.ExtractType(suggestion)
		before = fmt.Sprintf("    %s %s = ...", constName, constType)
		after = fmt.Sprintf(`    // %s defines ...
    %s %s = ...`, constName, constName, constType)

	case "KTN-CONST-004":
		// Constante sans type
		constName := messageutil.ExtractConstName(message)
		before = fmt.Sprintf("    %s = ...", constName)
		after = fmt.Sprintf("    %s int = ...", constName)

	default:
		return ""
	}

	if f.noColor {
		return fmt.Sprintf("  ❌ Avant:\n%s\n\n  ✅ Après:\n%s",
			f.indentCode(before, "    "),
			f.indentCode(after, "    "))
	}

	return fmt.Sprintf("  %s❌ Avant:%s\n%s\n\n  %s✅ Après:%s\n%s",
		Red, Reset,
		f.indentCode(before, "    "+Gray+"│"+Reset+" "),
		Green, Reset,
		f.indentCode(after, "    "+Green+"│"+Reset+" "))
}
