package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"sort"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
	"github.com/kodflow/ktn-linter/pkg/formatter"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// diagWithFset associe un diagnostic avec son FileSet et son analyseur
type diagWithFset struct {
	diag         analysis.Diagnostic
	fset         *token.FileSet
	analyzerName string
}

// lintCmd represents the lint command
var lintCmd *cobra.Command = &cobra.Command{
	Use:   "lint [packages...]",
	Short: "Lint Go packages using KTN rules",
	Long: `Lint analyzes Go packages and reports issues based on KTN conventions.

Examples:
  ktn-linter lint ./...
  ktn-linter lint -category=error ./...
  ktn-linter lint -ai ./path/to/file.go
  ktn-linter lint --fix ./...`,
	Args: cobra.MinimumNArgs(1),
	Run:  runLint,
}

// init enregistre la commande lint auprès de la commande root.
//
// Returns: aucun
//
// Params: aucun
func init() {
	rootCmd.AddCommand(lintCmd)
}

// runLint exécute l'analyse du linter.
//
// Params:
//   - cmd: commande Cobra
//   - args: arguments de la ligne de commande
//
// Returns: aucun
func runLint(cmd *cobra.Command, args []string) {
	pkgs := loadPackages(args)
	diagnostics := runAnalyzers(pkgs)

	// Filter out diagnostics from cache/tmp files (same logic as formatter)
	filteredDiags := filterDiagnostics(diagnostics)

	// Apply fixes if --fix flag is set
	if Fix {
		fixCount := applyFixes(filteredDiags)
		// Vérification de la condition
		if fixCount > 0 {
			fmt.Fprintf(os.Stderr, "Applied fixes to %d file(s)\n", fixCount)
  // Alternative path handling
		} else {
			fmt.Fprintf(os.Stderr, "No fixes to apply\n")
		}
		// Success - exit with 0
		OsExit(0)
	}

	formatAndDisplay(filteredDiags)

	// Vérification de la condition
	if len(filteredDiags) > 0 {
		OsExit(1)
	}

	// Success - exit with 0
	OsExit(0)
}

// loadPackages charge les packages Go à analyser.
//
// Params:
//   - patterns: liste des patterns de packages à charger
//
// Returns:
//   - []*packages.Package: packages chargés
func loadPackages(patterns []string) []*packages.Package {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests:      true,
		BuildFlags: []string{"-buildvcs=false"},
	}

	pkgs, err := packages.Load(cfg, patterns...)
	// Vérification de la condition
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading packages: %v\n", err)
		OsExit(1)
	}

	checkLoadErrors(pkgs)
	// Early return from function.
	return pkgs
}

// checkLoadErrors vérifie les erreurs de chargement des packages.
// Params:
//   - pass: contexte d'analyse
//
// Returns: aucun
//
//   - pkgs: liste des packages chargés
//
// Params:
func checkLoadErrors(pkgs []*packages.Package) {
	// hasLoadErrors holds the configuration value.

	var hasLoadErrors bool
	// Itération sur les éléments
	for _, pkg := range pkgs {
		// Itération sur les éléments
		for _, err := range pkg.Errors {
			// Only warn about VCS errors, don't exit
			if strings.Contains(err.Error(), "VCS status") {
				// Skip VCS errors
				continue
			}
			fmt.Fprintf(os.Stderr, "%v\n", err)
			hasLoadErrors = true
		}
	}
	// Vérification de la condition
	if hasLoadErrors {
		OsExit(1)
	}
}

// runAnalyzers exécute tous les analyseurs sur les packages.
//
// Params:
//   - pkgs: packages à analyser
//
// Returns:
//   - []diagWithFset: diagnostics trouvés
func runAnalyzers(pkgs []*packages.Package) []diagWithFset {
	// analyzers holds the configuration value.

	var analyzers []*analysis.Analyzer

	// Sélectionner les analyseurs selon la catégorie
	if Category != "" {
		analyzers = ktn.GetRulesByCategory(Category)
		// Vérification de la condition
		if analyzers == nil {
			fmt.Fprintf(os.Stderr, "Unknown category: %s\n", Category)
			OsExit(1)
		}
		// Vérification de la condition
		if Verbose {
			fmt.Fprintf(os.Stderr, "Running %d rules from category '%s'\n", len(analyzers), Category)
		}
		// Cas alternatif
	} else {
		analyzers = ktn.GetAllRules()
		// Vérification de la condition
		if Verbose {
			fmt.Fprintf(os.Stderr, "Running all %d KTN rules\n", len(analyzers))
		}
	}

	// allDiagnostics holds the configuration value.

	var allDiagnostics []diagWithFset

	// Store results of required analyzers (reused across packages)
	results := make(map[*analysis.Analyzer]any, len(analyzers))

	// Itération sur les éléments
	for _, pkg := range pkgs {
		pkgFset := pkg.Fset

		// Vérification de la condition
		if Verbose {
			fmt.Fprintf(os.Stderr, "Analyzing package: %s\n", pkg.PkgPath)
		}

		// Clear results for this package
		for k := range results {
			delete(results, k)
		}

		// Itération sur les éléments
		for _, a := range analyzers {
			pass := createAnalysisPass(a, pkg, pkgFset, &allDiagnostics, results)

			result, err := a.Run(pass)
			// Vérification de la condition
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running analyzer %s on %s: %v\n", a.Name, pkg.PkgPath, err)
			}
			results[a] = result
		}
	}

	// Early return from function.
	return allDiagnostics
}

// filterTestFiles filtre les fichiers de test.
//
// Params:
//   - files: fichiers à filtrer
//   - fset: fileset pour position
//
// Returns:
//   - []*ast.File: fichiers filtrés
func filterTestFiles(files []*ast.File, fset *token.FileSet) []*ast.File {
	filtered := make([]*ast.File, 0, len(files))
	// Itération sur les éléments
	for _, file := range files {
		pos := fset.Position(file.Pos())
		// Vérification de la condition
		if !strings.HasSuffix(pos.Filename, "_test.go") {
			filtered = append(filtered, file)
		}
	}
	// Retour de la fonction
	return filtered
}

// selectFilesForAnalyzer détermine les fichiers à analyser.
//
// Params:
//   - a: analyseur
//   - pkg: package
//   - fset: fileset
//
// Returns:
//   - []*ast.File: fichiers à analyser
func selectFilesForAnalyzer(a *analysis.Analyzer, pkg *packages.Package, fset *token.FileSet) []*ast.File {
	// TEST analyzers need all files (including test files)
	if strings.HasPrefix(a.Name, "ktntest") {
		// Return all files
		return pkg.Syntax
	}
	// Other analyzers get only non-test files
	return filterTestFiles(pkg.Syntax, fset)
}

// runRequiredAnalyzers exécute les analyseurs requis.
//
// Params:
//   - a: analyseur
//   - files: fichiers
//   - pkg: package
//   - fset: fileset
//   - results: map des résultats
func runRequiredAnalyzers(a *analysis.Analyzer, files []*ast.File, pkg *packages.Package, fset *token.FileSet, results map[*analysis.Analyzer]any) {
	// Run required analyzers first
	for _, req := range a.Requires {
		// IMPORTANT: Always run inspect.Analyzer with the correct file set
		// Different analyzers need different files (test vs non-test)
		// So we can't reuse inspect results across analyzers with different file sets
		reqPass := &analysis.Pass{
			Analyzer:  req,
			Fset:      fset,
			Files:     files,
			Pkg:       pkg.Types,
			TypesInfo: pkg.TypesInfo,
			ResultOf:  results,
			Report:    func(analysis.Diagnostic) {},
			ReadFile: func(filename string) ([]byte, error) {
				// Lit le contenu du fichier
				return os.ReadFile(filename)
			},
		}
		result, _ := req.Run(reqPass)
		results[req] = result
	}
}

// createAnalysisPass crée un pass d'analyse pour un package.
//
// Params:
//   - a: analyseur à exécuter
//   - pkg: package à analyser
//   - fset: fileset pour positions
//   - diagnostics: slice pour collecter diagnostics
//   - results: résultats des analyseurs requis
//
// Returns:
//   - *analysis.Pass: pass d'analyse créé
func createAnalysisPass(a *analysis.Analyzer, pkg *packages.Package, fset *token.FileSet, diagnostics *[]diagWithFset, results map[*analysis.Analyzer]any) *analysis.Pass {
	filesToAnalyze := selectFilesForAnalyzer(a, pkg, fset)
	runRequiredAnalyzers(a, filesToAnalyze, pkg, fset, results)

	// Early return from function.
	return &analysis.Pass{
		Analyzer:  a,
		Fset:      fset,
		Files:     filesToAnalyze,
		Pkg:       pkg.Types,
		TypesInfo: pkg.TypesInfo,
		ResultOf:  results,
		Report: func(diag analysis.Diagnostic) {
			*diagnostics = append(*diagnostics, diagWithFset{
				diag:         diag,
				fset:         fset,
				analyzerName: a.Name,
			})
		},
		ReadFile: func(filename string) ([]byte, error) {
			// Retour du contenu du fichier
			return os.ReadFile(filename)
		},
	}
}

// formatAndDisplay formate et affiche les diagnostics.
// Params:
//   - pass: contexte d'analyse
//
// Returns: aucun
//
//   - diagnostics: diagnostics à afficher
//
// Params:
func formatAndDisplay(diagnostics []diagWithFset) {
	fmtr := formatter.NewFormatter(os.Stdout, AIMode, NoColor, Simple)

	// Vérification de la condition
	if len(diagnostics) == 0 {
		fmtr.Format(nil, nil)
		// Early return from function.
		return
	}

	firstFset := diagnostics[0].fset
	diags := extractDiagnostics(diagnostics)
	fmtr.Format(firstFset, diags)
}

// filterDiagnostics filtre les diagnostics des fichiers cache/tmp.
//
// Params:
//   - diagnostics: diagnostics bruts avec fset
//
// Returns:
//   - []diagWithFset: diagnostics filtrés
func filterDiagnostics(diagnostics []diagWithFset) []diagWithFset {
	// filtered holds the configuration value.

	var filtered []diagWithFset
	// Itération sur les éléments
	for _, d := range diagnostics {
		pos := d.fset.Position(d.diag.Pos)
		// Ignorer les fichiers du cache Go et les fichiers temporaires
		// Vérification de la condition
		if strings.Contains(pos.Filename, "/.cache/go-build/") ||
			strings.Contains(pos.Filename, "/tmp/") ||
			strings.Contains(pos.Filename, "\\cache\\go-build\\") {
			continue
		}
		filtered = append(filtered, d)
	}
	// Early return from function.
	return filtered
}

// extractDiagnostics extrait et déduplique les diagnostics.
//
// Params:
//   - diagnostics: diagnostics bruts avec fset
//
// Returns:
//   - []analysis.Diagnostic: diagnostics dédupliqués
func extractDiagnostics(diagnostics []diagWithFset) []analysis.Diagnostic {
	// Dédupliquer les diagnostics (même position + même message)
	seen := make(map[string]bool, len(diagnostics))
	// deduped holds the configuration value.

	var deduped []diagWithFset
	// Itération sur les éléments
	for _, d := range diagnostics {
		pos := d.fset.Position(d.diag.Pos)
		key := fmt.Sprintf("%s:%d:%d:%s", pos.Filename, pos.Line, pos.Column, d.diag.Message)
		// Vérification de la condition
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, d)
		}
	}

	diags := make([]analysis.Diagnostic, 0, len(deduped))
	// Itération sur les éléments
	for _, d := range deduped {
		diags = append(diags, d.diag)
	}
	// Early return from function.
	return diags
}

// textEdit représente une modification de texte avec position.
type textEdit struct {
	start   int
	end     int
	newText []byte
}

// applyFixes applique les fixes suggérés aux fichiers source.
//
// Params:
//   - diagnostics: diagnostics avec fixes suggérés
//
// Returns:
//   - int: nombre de fichiers modifiés
func applyFixes(diagnostics []diagWithFset) int {
	// Liste blanche des analyseurs dont les fixes sont sûrs
	safeAnalyzers := map[string]bool{
		"any": true, // interface{} → any
	}

	// Collecter les éditions par fichier
	fileEdits, skippedCount := collectSafeEdits(diagnostics, safeAnalyzers)

	// Afficher le nombre de fixes skippés
	// Vérification de la condition
	if skippedCount > 0 {
		fmt.Fprintf(os.Stderr, "Skipped %d unsafe fixes (use modernize tool for complex transformations)\n", skippedCount)
	}

	// Appliquer les éditions collectées
	// Early return from function.
	return applyCollectedEdits(fileEdits)
}

// collectSafeEdits collecte les éditions de texte sûres depuis les diagnostics.
//
// Params:
//   - diagnostics: diagnostics avec fixes suggérés
//   - safeAnalyzers: map des analyseurs sûrs
//
// Returns:
//   - map[string][]textEdit: éditions groupées par fichier
//   - int: nombre de fixes skippés
func collectSafeEdits(diagnostics []diagWithFset, safeAnalyzers map[string]bool) (map[string][]textEdit, int) {
	fileEdits := make(map[string][]textEdit)
	skippedCount := 0

	// Parcourir tous les diagnostics
	// Itération sur les éléments
	for _, d := range diagnostics {
		// Skip diagnostics without suggested fixes
		if len(d.diag.SuggestedFixes) == 0 {
			continue
		}

		// Skip analyzers not in safe list
		// Vérification de la condition
		if !safeAnalyzers[d.analyzerName] {
			skippedCount++
			// Log si verbose
			// Vérification de la condition
			if Verbose {
				pos := d.fset.Position(d.diag.Pos)
				fmt.Fprintf(os.Stderr, "Skipping unsafe fix at %s:%d (analyzer: %s)\n",
					pos.Filename, pos.Line, d.analyzerName)
			}
			continue
		}

		// Extraire les éditions de texte
		extractTextEdits(d, &fileEdits)
	}

	// Early return from function.
	return fileEdits, skippedCount
}

// extractTextEdits extrait les éditions de texte depuis un diagnostic.
//
// Params:
//   - d: diagnostic avec fixes
//   - fileEdits: map pour stocker les éditions par fichier
func extractTextEdits(d diagWithFset, fileEdits *map[string][]textEdit) {
	// Process only the first suggested fix
	if len(d.diag.SuggestedFixes) == 0 {
		// Early return from function.
		return
	}

	fix := d.diag.SuggestedFixes[0]

	// Process all text edits in this fix
	// Itération sur les éléments
	for _, edit := range fix.TextEdits {
		// Convert token positions to byte offsets
		file := d.fset.File(edit.Pos)
		// Vérification de la condition
		if file == nil {
			continue
		}

		startOffset := file.Offset(edit.Pos)
		endOffset := file.Offset(edit.End)
		pos := d.fset.Position(edit.Pos)

		// Store edit for this file
		(*fileEdits)[pos.Filename] = append((*fileEdits)[pos.Filename], textEdit{
			start:   startOffset,
			end:     endOffset,
			newText: edit.NewText,
		})
	}
}

// applyCollectedEdits applique les éditions collectées à chaque fichier.
//
// Params:
//   - fileEdits: éditions groupées par fichier
//
// Returns:
//   - int: nombre de fichiers modifiés
func applyCollectedEdits(fileEdits map[string][]textEdit) int {
	fixCount := 0

	// Itération sur les éléments
	for filename, edits := range fileEdits {
		// Vérification de la condition
		if applyEditsToFile(filename, edits) {
			fixCount++
			// Log si verbose
			// Vérification de la condition
			if Verbose {
				fmt.Fprintf(os.Stderr, "Applied %d edits to %s\n", len(edits), filename)
			}
		}
	}

	// Early return from function.
	return fixCount
}

// applyEditsToFile applique les modifications à un fichier.
//
// Params:
//   - filename: chemin du fichier
//   - edits: modifications à appliquer
//
// Returns:
//   - bool: true si succès, false sinon
func applyEditsToFile(filename string, edits []textEdit) bool {
	// Read file content
	content, err := os.ReadFile(filename)
	// Vérification de la condition
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		// Early return from function.
		return false
	}

	// Sort edits by start position (descending) to apply from end to start
	sort.Slice(edits, func(i, j int) bool {
		// Sort by start position in reverse order
		return edits[i].start > edits[j].start
	})

	// Remove overlapping edits (keep only first in reverse order = last in file)
	nonOverlapping := filterOverlappingEdits(edits)

	// Apply each edit from end to start
	result := content
	// Itération sur les éléments
	for _, edit := range nonOverlapping {
		// Validate edit positions
		if edit.start < 0 || edit.end > len(result) || edit.start > edit.end {
			fmt.Fprintf(os.Stderr, "Invalid edit in %s: start=%d, end=%d, len=%d\n",
				filename, edit.start, edit.end, len(result))
			continue
		}

		// Apply the edit: before + newText + after
		var buf bytes.Buffer
		buf.Write(result[:edit.start])
		buf.Write(edit.newText)
		buf.Write(result[edit.end:])
		result = buf.Bytes()
	}

	// Write back to file with same permissions
	err = os.WriteFile(filename, result, 0644)
	// Vérification de la condition
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", filename, err)
		// Early return from function.
		return false
	}

	// Early return from function.
	return true
}

// filterOverlappingEdits supprime les éditions qui se chevauchent.
//
// Params:
//   - edits: éditions triées par position décroissante
//
// Returns:
//   - []textEdit: éditions non-chevauchantes
func filterOverlappingEdits(edits []textEdit) []textEdit {
	// Vérification de la condition
	if len(edits) == 0 {
		// Early return from function.
		return edits
	}

	// Keep track of the last edit's start position
	filtered := []textEdit{edits[0]}
	lastStart := edits[0].start

	// Check each subsequent edit
	for i := 1; i < len(edits); i++ {
		edit := edits[i]
		// No overlap if this edit ends before the last edit starts
		if edit.end <= lastStart {
			filtered = append(filtered, edit)
			lastStart = edit.start
		} else {
			// Skip overlapping edit
			if Verbose {
				fmt.Fprintf(os.Stderr, "Skipping overlapping edit at [%d:%d]\n", edit.start, edit.end)
			}
		}
	}

	// Early return from function.
	return filtered
}
