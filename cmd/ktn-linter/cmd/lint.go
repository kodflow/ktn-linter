// Lint command implementation for analyzing Go code.
package cmd

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/formatter"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// lintOptions holds all options for the lint command.
type lintOptions struct {
	verbose    bool
	category   string
	onlyRule   string
	configPath string
}

// lintCmd represents the lint command.
var lintCmd *cobra.Command = &cobra.Command{
	Use:   "lint [packages...]",
	Short: "Lint Go packages using KTN rules",
	Long:  `Lint analyzes Go packages and reports issues based on KTN conventions.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runLint,
}

// init enregistre la commande lint auprès de la commande root.
//
// Returns: aucun
//
// Params: aucun
func init() {
	rootCmd.AddCommand(lintCmd)
}

// loadConfiguration charge la configuration du linter.
//
// Params:
//   - opts: options de lint
//
// Returns: aucun
func loadConfiguration(opts lintOptions) {
	// Vérification si un fichier de config est spécifié
	if opts.configPath != "" {
		// Charger depuis le fichier spécifié
		if err := config.LoadAndSet(opts.configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file %s: %v\n", opts.configPath, err)
			OsExit(1)
		}
		// Log si verbose
		if opts.verbose {
			fmt.Fprintf(os.Stderr, "Loaded configuration from %s\n", opts.configPath)
		}
		// Retour de la fonction
		return
	}

	// Tenter de charger depuis les emplacements par défaut
	if err := config.LoadAndSet(""); err == nil {
		// Log si verbose et config trouvée
		if opts.verbose {
			fmt.Fprintf(os.Stderr, "Loaded configuration from default location\n")
		}
	}
	// Si pas de config par défaut, utiliser les valeurs par défaut (pas d'erreur)
}

// runLint exécute l'analyse du linter.
//
// Params:
//   - cmd: commande Cobra (utilisé pour récupérer les flags)
//   - args: arguments de la ligne de commande
//
// Returns: aucun
func runLint(cmd *cobra.Command, args []string) {
	// Récupérer les options depuis les flags Cobra
	opts := parseLintOptions(cmd)

	// Charger la configuration si spécifiée
	loadConfiguration(opts)

	// Propager le flag verbose dans la config pour les règles
	config.Get().Verbose = opts.verbose

	pkgs := loadPackages(args)
	diagnostics := runAnalyzers(pkgs, opts)

	// Filter out diagnostics from cache/tmp files (same logic as formatter)
	filteredDiags := filterDiagnostics(diagnostics)

	formatAndDisplay(filteredDiags, opts)

	// Vérification de la condition
	if len(filteredDiags) > 0 {
		OsExit(1)
	}

	// Success - exit with 0
	OsExit(0)
}

// parseLintOptions extrait les options depuis les flags Cobra.
//
// Params:
//   - cmd: commande Cobra (unused, kept for API compatibility)
//
// Returns:
//   - lintOptions: options extraites
func parseLintOptions(_ *cobra.Command) lintOptions {
	// Use rootCmd.PersistentFlags() directly since persistent flags are defined there
	flags := rootCmd.PersistentFlags()

	verbose, _ := flags.GetBool(flagVerbose)
	category, _ := flags.GetString(flagCategory)
	onlyRule, _ := flags.GetString(flagOnlyRule)
	configPath, _ := flags.GetString(flagConfig)

	// Retour des options parsées
	return lintOptions{
		verbose:    verbose,
		category:   category,
		onlyRule:   onlyRule,
		configPath: configPath,
	}
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

// selectAnalyzers sélectionne les analyseurs selon les filtres.
//
// Params:
//   - opts: options de lint
//
// Returns:
//   - []*analysis.Analyzer: analyseurs sélectionnés
func selectAnalyzers(opts lintOptions) []*analysis.Analyzer {
	// Sélection selon --only-rule
	if opts.onlyRule != "" {
		// Retourner l'analyseur unique
		return selectSingleRule(opts)
	}

	// Sélection selon --category
	if opts.category != "" {
		// Retourner les analyseurs de la catégorie
		return selectByCategory(opts)
	}

	// Mode par défaut: toutes les règles
	analyzers := ktn.GetAllRules()

	// Log si verbose
	if opts.verbose {
		fmt.Fprintf(os.Stderr, "Running all %d KTN rules\n", len(analyzers))
	}

	// Retourner tous les analyseurs
	return analyzers
}

// selectSingleRule sélectionne une seule règle spécifique.
//
// Params:
//   - opts: options de lint
//
// Returns:
//   - []*analysis.Analyzer: analyseur unique
func selectSingleRule(opts lintOptions) []*analysis.Analyzer {
	// Récupérer la règle
	analyzer := ktn.GetRuleByCode(opts.onlyRule)

	// Vérifier si la règle existe
	if analyzer == nil {
		fmt.Fprintf(os.Stderr, "Unknown rule code: %s\n", opts.onlyRule)
		OsExit(1)
	}

	// Log si verbose
	if opts.verbose {
		fmt.Fprintf(os.Stderr, "Running only rule '%s'\n", opts.onlyRule)
	}

	// Retourner l'analyseur unique
	return []*analysis.Analyzer{analyzer}
}

// selectByCategory sélectionne les règles d'une catégorie.
//
// Params:
//   - opts: options de lint
//
// Returns:
//   - []*analysis.Analyzer: analyseurs de la catégorie
func selectByCategory(opts lintOptions) []*analysis.Analyzer {
	// Récupérer les analyseurs de la catégorie
	analyzers := ktn.GetRulesByCategory(opts.category)

	// Vérifier si la catégorie existe
	if len(analyzers) == 0 {
		fmt.Fprintf(os.Stderr, "Unknown category: %s\n", opts.category)
		OsExit(1)
	}

	// Log si verbose
	if opts.verbose {
		fmt.Fprintf(os.Stderr, "Running %d rules from category '%s'\n", len(analyzers), opts.category)
	}

	// Retourner les analyseurs
	return analyzers
}

// runAnalyzers exécute tous les analyseurs sur les packages.
//
// Params:
//   - pkgs: packages à analyser
//   - opts: options de lint
//
// Returns:
//   - []diagWithFset: diagnostics trouvés
func runAnalyzers(pkgs []*packages.Package, opts lintOptions) []diagWithFset {
	// Sélectionner les analyseurs
	analyzers := selectAnalyzers(opts)

	// allDiagnostics collecte les diagnostics
	var allDiagnostics []diagWithFset

	// Store results of required analyzers (reused across packages)
	results := make(map[*analysis.Analyzer]any, len(analyzers))

	// Analyser chaque package
	for _, pkg := range pkgs {
		// Analyser le package
		analyzePackage(pkg, analyzers, results, &allDiagnostics, opts)
	}

	// Retourner les diagnostics
	return allDiagnostics
}

// analyzePackage analyse un package avec les analyseurs donnés.
//
// Params:
//   - pkg: package à analyser
//   - analyzers: analyseurs à exécuter
//   - results: résultats des analyseurs (modifié in-place)
//   - allDiagnostics: diagnostics collectés (modifié in-place)
//   - opts: options de lint
func analyzePackage(pkg *packages.Package, analyzers []*analysis.Analyzer, results map[*analysis.Analyzer]any, allDiagnostics *[]diagWithFset, opts lintOptions) {
	pkgFset := pkg.Fset

	// Log si verbose
	if opts.verbose {
		fmt.Fprintf(os.Stderr, "Analyzing package: %s\n", pkg.PkgPath)
	}

	// Clear results for this package
	for k := range results {
		delete(results, k)
	}

	// Exécuter chaque analyseur
	for _, a := range analyzers {
		// Créer le pass et exécuter
		pass := createAnalysisPass(a, pkg, pkgFset, allDiagnostics, results)
		result, err := a.Run(pass)

		// Gérer les erreurs
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running analyzer %s on %s: %v\n", a.Name, pkg.PkgPath, err)
		}

		// Stocker le résultat
		results[a] = result
	}
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

	// Check if force mode is enabled (run all rules on test files)
	if config.Get().ForceAllRulesOnTests {
		// Return all files in force mode
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
//
// Params:
//   - diagnostics: diagnostics à afficher
//   - opts: options de lint
//
// Returns: aucun
func formatAndDisplay(diagnostics []diagWithFset, opts lintOptions) {
	// Simple format is now the default (no color, one-line per diagnostic)
	fmtr := formatter.NewFormatter(os.Stdout, false, true, true, opts.verbose)

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
		// Ignorer uniquement les fichiers du cache Go (pas les projets utilisateur dans /tmp)
		// Vérification de la condition
		if strings.Contains(pos.Filename, "/.cache/go-build/") ||
			strings.Contains(pos.Filename, "\\cache\\go-build\\") {
			continue
		}
		filtered = append(filtered, d)
	}
	// Early return from function.
	return filtered
}

// isModernizeAnalyzer vérifie si un analyseur est un analyseur modernize.
//
// Params:
//   - name: nom de l'analyseur
//
// Returns:
//   - bool: true si c'est un analyseur modernize
func isModernizeAnalyzer(name string) bool {
	// Liste des analyseurs modernize de golang.org/x/tools (mise à jour)
	modernizeAnalyzers := map[string]bool{
		"any":              true,
		"bloop":            true,
		"fmtappendf":       true,
		"forvar":           true,
		"mapsloop":         true,
		"minmax":           true,
		"newexpr":          true,
		"omitzero":         true,
		"rangeint":         true,
		"reflecttypefor":   true,
		"slicescontains":   true,
		"slicessort":       true,
		"stditerators":     true,
		"stringscutprefix": true,
		"stringsseq":       true,
		"stringsbuilder":   true,
		"testingcontext":   true,
		"waitgroup":        true,
	}
	// Retour du résultat
	return modernizeAnalyzers[name]
}

// formatModernizeCode formate un nom d'analyseur en code KTN-MDRNZ.
//
// Params:
//   - name: nom de l'analyseur
//
// Returns:
//   - string: code formaté (ex: "KTN-MDRNZ-ANY")
func formatModernizeCode(name string) string {
	// Retour du code formaté en majuscules
	return "KTN-MDRNZ-" + strings.ToUpper(name)
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
		diag := d.diag
		// Préfixer les messages modernize avec le code KTN-MDRNZ
		if isModernizeAnalyzer(d.analyzerName) && !strings.HasPrefix(diag.Message, "KTN-") {
			code := formatModernizeCode(d.analyzerName)
			diag.Message = code + ": " + diag.Message
		}
		diags = append(diags, diag)
	}
	// Early return from function.
	return diags
}

