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

// Codes de couleurs ANSI pour le formatage terminal
// Ces constantes définissent les séquences d'échappement pour colorer la sortie
const (
	// Red applique la couleur rouge
	Red string = "\033[31m"
	// Green applique la couleur verte
	Green string = "\033[32m"
	// Yellow applique la couleur jaune
	Yellow string = "\033[33m"
	// Blue applique la couleur bleue
	Blue string = "\033[34m"
	// Magenta applique la couleur magenta
	Magenta string = "\033[35m"
	// Cyan applique la couleur cyan
	Cyan string = "\033[36m"
	// Gray applique la couleur grise
	Gray string = "\033[90m"
	// Bold applique le style gras
	Bold string = "\033[1m"
	// Reset réinitialise tous les styles et couleurs
	Reset string = "\033[0m"
)

// DiagnosticGroupData regroupe les diagnostics par fichier (DTO)
type DiagnosticGroupData struct {
	// Filename est le chemin du fichier contenant les diagnostics
	Filename string
	// Diagnostics est la liste des diagnostics trouvés dans ce fichier
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
//   - Formatter: un formatter prêt à utiliser pour afficher les diagnostics
func NewFormatter(w io.Writer, aiMode bool, noColor bool, simpleMode bool) Formatter {
	// Retourne un formatter configuré avec les options spécifiées
	return &formatterImpl{
		writer:     w,
		aiMode:     aiMode,
		noColor:    noColor,
		simpleMode: simpleMode,
	}
}

// Format affiche les diagnostics de manière lisible
//
// Params:
//   - fset: le FileSet contenant les informations de position
//   - diagnostics: la liste des diagnostics à formater
func (f *formatterImpl) Format(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	if len(diagnostics) == 0 {
		f.printSuccess()
		// Retourne après affichage du message de succès car aucun diagnostic à traiter
		return
	}

	if f.simpleMode {
		f.formatSimple(fset, diagnostics)
		// Retourne après formatage en mode simple (une ligne par erreur)
		return
	}

	if f.aiMode {
		f.formatForAI(fset, diagnostics)
		// Retourne après formatage optimisé pour l'IA
		return
	}

	// Format par défaut: affichage pour humain avec couleurs et structure
	f.formatForHuman(fset, diagnostics)
}

// formatForHuman affiche pour un humain avec couleurs et structure
//
// Params:
//   - fset: le FileSet contenant les informations de position
//   - diagnostics: la liste des diagnostics à formater
func (f *formatterImpl) formatForHuman(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	groups := f.groupByFile(fset, diagnostics)

	// Compter le nombre réel de diagnostics après filtrage
	totalCount := 0
	for _, group := range groups {
		totalCount += len(group.Diagnostics)
	}

	// Si tous les diagnostics ont été filtrés, afficher le succès
	if totalCount == 0 {
		f.printSuccess()
		// Retourne après affichage du succès car tous les diagnostics ont été filtrés
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

// countTotalDiagnostics calcule le nombre total de diagnostics après filtrage.
//
// Params:
//   - groups: les groupes de diagnostics par fichier
//
// Returns:
//   - int: le nombre total de diagnostics
func (f *formatterImpl) countTotalDiagnostics(groups []DiagnosticGroupData) int {
	totalCount := 0
	for _, group := range groups {
		totalCount += len(group.Diagnostics)
	}
	// Retourne le nombre total de diagnostics trouvés après filtrage
	return totalCount
}

// printAIFileGroups affiche les groupes de fichiers au format IA.
//
// Params:
//   - groups: les groupes de diagnostics par fichier
//   - fset: le FileSet pour obtenir les positions
func (f *formatterImpl) printAIFileGroups(groups []DiagnosticGroupData, fset *token.FileSet) {
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
}

// printAIInstructions affiche les instructions pour l'IA.
func (f *formatterImpl) printAIInstructions() {
	fmt.Fprintf(f.writer, "\n---\n")
	fmt.Fprintf(f.writer, "**Instructions for AI**:\n")
	fmt.Fprintf(f.writer, "- Each issue needs to be fixed according to its code and suggestion\n")
	fmt.Fprintf(f.writer, "- Group related constants together in const () blocks\n")
	fmt.Fprintf(f.writer, "- Add documentation comments for all constants and groups\n")
	fmt.Fprintf(f.writer, "- Always specify explicit types for constants\n")
	fmt.Fprintf(f.writer, "- Refer to target.go for examples of correct implementations\n")
}

// formatForAI affiche un format optimisé pour l'IA
//
// Params:
//   - fset: le FileSet contenant les informations de position
//   - diagnostics: la liste des diagnostics à formater
func (f *formatterImpl) formatForAI(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	groups := f.groupByFile(fset, diagnostics)
	totalCount := f.countTotalDiagnostics(groups)

	fmt.Fprintf(f.writer, "# KTN-Linter Report (AI Mode)\n\n")
	fmt.Fprintf(f.writer, "Total issues found: %d\n\n", totalCount)

	f.printAIFileGroups(groups, fset)
	f.printAIInstructions()
}

// filterAndSortDiagnostics filtre et trie les diagnostics par position.
//
// Params:
//   - fset: le FileSet contenant les informations de position
//   - diagnostics: la liste des diagnostics à filtrer et trier
//
// Returns:
//   - []analysis.Diagnostic: les diagnostics filtrés et triés
func (f *formatterImpl) filterAndSortDiagnostics(fset *token.FileSet, diagnostics []analysis.Diagnostic) []analysis.Diagnostic {
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
			// Retourne true si le premier fichier vient avant le second dans l'ordre alphabétique
			return posI.Filename < posJ.Filename
		}
		if posI.Line != posJ.Line {
			// Retourne true si la première ligne vient avant la seconde
			return posI.Line < posJ.Line
		}
		// Retourne true si la première colonne vient avant la seconde
		return posI.Column < posJ.Column
	})

	// Retourne les diagnostics filtrés et triés par fichier, ligne et colonne
	return filtered
}

// formatSimple affiche un format simple une ligne par erreur (pour IDE)
//
// Params:
//   - fset: le FileSet contenant les informations de position
//   - diagnostics: la liste des diagnostics à formater
func (f *formatterImpl) formatSimple(fset *token.FileSet, diagnostics []analysis.Diagnostic) {
	filtered := f.filterAndSortDiagnostics(fset, diagnostics)

	// Afficher chaque diagnostic sur une ligne
	for _, diag := range filtered {
		pos := fset.Position(diag.Pos)
		code := messageutil.ExtractCode(diag.Message)
		message := messageutil.ExtractMessage(diag.Message)

		// Format: file:line:column: [CODE] message
		fmt.Fprintf(f.writer, "%s:%d:%d: [%s] %s\n",
			pos.Filename, pos.Line, pos.Column, code, message)
	}
}

// buildFileMapFiltered construit une map de diagnostics filtrés par fichier.
//
// Params:
//   - fset: le FileSet pour obtenir les positions des diagnostics
//   - diagnostics: la liste des diagnostics à grouper
//
// Returns:
//   - map[string][]analysis.Diagnostic: map des diagnostics par fichier
func (f *formatterImpl) buildFileMapFiltered(fset *token.FileSet, diagnostics []analysis.Diagnostic) map[string][]analysis.Diagnostic {
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

	// Retourne la map des diagnostics groupés par fichier après filtrage
	return fileMap
}

// sortGroupsFromMap crée et trie les groupes de diagnostics à partir d'une map.
//
// Params:
//   - fset: le FileSet pour obtenir les positions des diagnostics
//   - fileMap: la map des diagnostics par fichier
//
// Returns:
//   - []DiagnosticGroupData: les groupes de diagnostics triés
func (f *formatterImpl) sortGroupsFromMap(fset *token.FileSet, fileMap map[string][]analysis.Diagnostic) []DiagnosticGroupData {
	var groups []DiagnosticGroupData
	for filename, diags := range fileMap {
		// Trier par ligne
		sort.Slice(diags, func(i, j int) bool {
			// Retourne true si le premier diagnostic est situé avant le second dans le fichier
			return fset.Position(diags[i].Pos).Line < fset.Position(diags[j].Pos).Line
		})
		groups = append(groups, DiagnosticGroupData{
			Filename:    filename,
			Diagnostics: diags,
		})
	}

	// Trier par nom de fichier
	sort.Slice(groups, func(i, j int) bool {
		// Retourne true si le premier groupe vient avant le second dans l'ordre alphabétique
		return groups[i].Filename < groups[j].Filename
	})

	// Retourne les groupes de diagnostics triés par fichier et par ligne
	return groups
}

// groupByFile regroupe les diagnostics par fichier et les trie.
//
// Params:
//   - fset: le FileSet pour obtenir les positions des diagnostics
//   - diagnostics: la liste des diagnostics à grouper
//
// Returns:
//   - []DiagnosticGroupData: les diagnostics regroupés et triés par fichier
func (f *formatterImpl) groupByFile(fset *token.FileSet, diagnostics []analysis.Diagnostic) []DiagnosticGroupData {
	fileMap := f.buildFileMapFiltered(fset, diagnostics)
	// Retourne les diagnostics groupés et triés par fichier et par position
	return f.sortGroupsFromMap(fset, fileMap)
}

// printHeader affiche l'en-tête du rapport avec le nombre total de problèmes.
//
// Params:
//   - count: le nombre total de problèmes trouvés
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

// printFileHeader affiche l'en-tête pour un fichier avec son nombre de problèmes.
//
// Params:
//   - filename: le chemin du fichier
//   - count: le nombre de problèmes dans ce fichier
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

// printDiagnosticNoColor affiche un diagnostic sans couleurs.
//
// Params:
//   - num: le numéro séquentiel du diagnostic dans le fichier
//   - location: la localisation formatée (fichier:ligne:colonne)
//   - code: le code d'erreur
//   - message: le message d'erreur
//   - example: l'exemple avant/après
func (f *formatterImpl) printDiagnosticNoColor(num int, location, code, message, example string) {
	fmt.Fprintf(f.writer, "\n[%d] %s\n", num, location)
	fmt.Fprintf(f.writer, "  Code: %s\n", code)
	fmt.Fprintf(f.writer, "  Issue: %s\n", message)
	if example != "" {
		fmt.Fprintf(f.writer, "\n%s\n", example)
	}
}

// printDiagnosticWithColor affiche un diagnostic avec couleurs.
//
// Params:
//   - num: le numéro séquentiel du diagnostic dans le fichier
//   - location: la localisation formatée (fichier:ligne:colonne)
//   - code: le code d'erreur
//   - message: le message d'erreur
//   - example: l'exemple avant/après
func (f *formatterImpl) printDiagnosticWithColor(num int, location, code, message, example string) {
	codeColor := f.getCodeColor(code)

	// Numéro et location cliquable
	fmt.Fprintf(f.writer, "\n%s[%d]%s %s%s%s\n",
		Bold+Yellow, num, Reset,
		Cyan, location, Reset)

	// Code d'erreur
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

// printDiagnostic affiche un diagnostic individuel avec son contexte et ses exemples.
//
// Params:
//   - num: le numéro séquentiel du diagnostic dans le fichier
//   - pos: la position du diagnostic dans le code source
//   - diag: le diagnostic à afficher
func (f *formatterImpl) printDiagnostic(num int, pos token.Position, diag analysis.Diagnostic) {
	code := messageutil.ExtractCode(diag.Message)
	message := messageutil.ExtractMessage(diag.Message)
	suggestion := messageutil.ExtractSuggestion(diag.Message)
	example := f.generateExample(code, message, suggestion)

	// Format cliquable : fichier:ligne:colonne
	location := fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)

	if f.noColor {
		f.printDiagnosticNoColor(num, location, code, message, example)
	} else {
		f.printDiagnosticWithColor(num, location, code, message, example)
	}
}

// printSuccess affiche un message de succès quand aucun problème n'est trouvé.
func (f *formatterImpl) printSuccess() {
	if f.noColor {
		fmt.Fprintf(f.writer, "\n✅ No issues found! Code is compliant.\n\n")
	} else {
		fmt.Fprintf(f.writer, "\n%s✅ No issues found! Code is compliant.%s\n\n", Bold+Green, Reset)
	}
}

// printSummary affiche le résumé final avec le nombre total de problèmes à corriger.
//
// Params:
//   - count: le nombre total de problèmes à corriger
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

// indentCode ajoute une indentation à chaque ligne non-vide d'un bloc de code.
//
// Params:
//   - code: le code source à indenter
//   - indent: la chaîne d'indentation à ajouter au début de chaque ligne
//
// Returns:
//   - string: le code source indenté
func (f *formatterImpl) indentCode(code string, indent string) string {
	lines := strings.Split(code, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, indent+line)
		}
	}
	// Retourne le code avec toutes les lignes non-vides indentées
	return strings.Join(result, "\n")
}

// getCodeColor retourne la couleur ANSI appropriée pour un code d'erreur donné.
//
// Params:
//   - code: le code d'erreur (ex: "KTN-CONST-001")
//
// Returns:
//   - string: la séquence d'échappement ANSI pour la couleur, ou chaîne vide si noColor
func (f *formatterImpl) getCodeColor(code string) string {
	if f.noColor {
		// Retourne une chaîne vide si les couleurs sont désactivées
		return ""
	}

	switch {
	case strings.HasSuffix(code, "-001"):
		// Retourne rouge pour les codes -001 (problèmes critiques)
		return Red
	case strings.HasSuffix(code, "-002"):
		// Retourne jaune pour les codes -002 (avertissements)
		return Yellow
	case strings.HasSuffix(code, "-003"):
		// Retourne magenta pour les codes -003
		return Magenta
	case strings.HasSuffix(code, "-004"):
		// Retourne cyan pour les codes -004
		return Cyan
	default:
		// Retourne rouge par défaut pour les autres codes
		return Red
	}
}

// generateExampleConst001 génère l'exemple pour KTN-CONST-001 (constante non groupée).
//
// Params:
//   - message: le message d'erreur contenant le nom de la constante
//   - suggestion: la suggestion contenant le type
//
// Returns:
//   - before: le code avant correction
//   - after: le code après correction
func (f *formatterImpl) generateExampleConst001(message, suggestion string) (before, after string) {
	constName := messageutil.ExtractConstName(message)
	constType := messageutil.ExtractType(suggestion)
	before = fmt.Sprintf("const %s %s = ...", constName, constType)
	after = fmt.Sprintf(`const (
    %s %s = ...
)`, constName, constType)
	// Retourne l'exemple avant/après pour une constante non groupée
	return before, after
}

// generateExampleConst002 génère l'exemple pour KTN-CONST-002 (groupe sans commentaire).
//
// Returns:
//   - before: le code avant correction
//   - after: le code après correction
func (f *formatterImpl) generateExampleConst002() (before, after string) {
	before = `const (
    MaxValue int = 100
)`
	after = `// Configuration constants
// Define application limits
const (
    MaxValue int = 100
)`
	// Retourne l'exemple avant/après pour un groupe sans commentaire
	return before, after
}

// generateExampleConst003 génère l'exemple pour KTN-CONST-003 (constante sans commentaire).
//
// Params:
//   - message: le message d'erreur contenant le nom de la constante
//   - suggestion: la suggestion contenant le type
//
// Returns:
//   - before: le code avant correction
//   - after: le code après correction
func (f *formatterImpl) generateExampleConst003(message, suggestion string) (before, after string) {
	constName := messageutil.ExtractConstName(message)
	constType := messageutil.ExtractType(suggestion)
	before = fmt.Sprintf("    %s %s = ...", constName, constType)
	after = fmt.Sprintf(`    // %s defines ...
    %s %s = ...`, constName, constName, constType)
	// Retourne l'exemple avant/après pour une constante sans commentaire
	return before, after
}

// generateExampleConst004 génère l'exemple pour KTN-CONST-004 (constante sans type).
//
// Params:
//   - message: le message d'erreur contenant le nom de la constante
//
// Returns:
//   - before: le code avant correction
//   - after: le code après correction
func (f *formatterImpl) generateExampleConst004(message string) (before, after string) {
	constName := messageutil.ExtractConstName(message)
	before = fmt.Sprintf("    %s = ...", constName)
	after = fmt.Sprintf("    %s int = ...", constName)
	// Retourne l'exemple avant/après pour une constante sans type
	return before, after
}

// formatExampleBeforeAfter formate un exemple avant/après avec ou sans couleurs.
//
// Params:
//   - before: le code avant correction
//   - after: le code après correction
//
// Returns:
//   - string: l'exemple formaté avec le style approprié
func (f *formatterImpl) formatExampleBeforeAfter(before, after string) string {
	if f.noColor {
		// Retourne l'exemple formaté sans couleurs
		return fmt.Sprintf("  ❌ Avant:\n%s\n\n  ✅ Après:\n%s",
			f.indentCode(before, "    "),
			f.indentCode(after, "    "))
	}

	// Retourne l'exemple formaté avec couleurs et séparateurs visuels
	return fmt.Sprintf("  %s❌ Avant:%s\n%s\n\n  %s✅ Après:%s\n%s",
		Red, Reset,
		f.indentCode(before, "    "+Gray+"│"+Reset+" "),
		Green, Reset,
		f.indentCode(after, "    "+Green+"│"+Reset+" "))
}

// generateExample crée un exemple avant/après concret
//
// Params:
//   - code: le code d'erreur (ex: "KTN-CONST-001")
//   - message: le message d'erreur contenant les détails
//   - suggestion: la suggestion de correction
//
// Returns:
//   - string: l'exemple formaté avec le code avant/après, ou chaîne vide si non applicable
func (f *formatterImpl) generateExample(code, message, suggestion string) string {
	var before, after string

	switch code {
	case "KTN-CONST-001":
		before, after = f.generateExampleConst001(message, suggestion)
	case "KTN-CONST-002":
		before, after = f.generateExampleConst002()
	case "KTN-CONST-003":
		before, after = f.generateExampleConst003(message, suggestion)
	case "KTN-CONST-004":
		before, after = f.generateExampleConst004(message)
	default:
		// Retourne une chaîne vide pour les codes non reconnus
		return ""
	}

	// Retourne l'exemple formaté avec le style approprié (couleurs ou pas)
	return f.formatExampleBeforeAfter(before, after)
}
