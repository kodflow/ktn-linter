package cmd

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
	"github.com/kodflow/ktn-linter/pkg/formatter"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// diagWithFset associe un diagnostic avec son FileSet
type diagWithFset struct {
	diag analysis.Diagnostic
	fset *token.FileSet
}

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint [packages...]",
	Short: "Lint Go packages using KTN rules",
	Long: `Lint analyzes Go packages and reports issues based on KTN conventions.

Examples:
  ktn-linter lint ./...
  ktn-linter lint -category=error ./...
  ktn-linter lint -ai ./path/to/file.go`,
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
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true,
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

	// Itération sur les éléments
	for _, pkg := range pkgs {
		pkgFset := pkg.Fset

		// Vérification de la condition
		if Verbose {
			fmt.Fprintf(os.Stderr, "Analyzing package: %s\n", pkg.PkgPath)
		}

		// Store results of required analyzers
		results := make(map[*analysis.Analyzer]interface{})

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
func createAnalysisPass(a *analysis.Analyzer, pkg *packages.Package, fset *token.FileSet, diagnostics *[]diagWithFset, results map[*analysis.Analyzer]interface{}) *analysis.Pass {
	// Filter out test files
	nonTestFiles := filterTestFiles(pkg.Syntax, fset)

	// Run required analyzers first
	for _, req := range a.Requires {
		// Vérification de la condition
		if _, ok := results[req]; !ok {
			reqPass := &analysis.Pass{
				Analyzer:  req,
				Fset:      fset,
				Files:     nonTestFiles,
				Pkg:       pkg.Types,
				TypesInfo: pkg.TypesInfo,
				ResultOf:  results,
				Report:    func(analysis.Diagnostic) {},
				ReadFile: func(filename string) ([]byte, error) {
					// Retour du contenu du fichier
					return os.ReadFile(filename)
				},
			}
			result, _ := req.Run(reqPass)
			results[req] = result
		}
	}

	// Early return from function.
	return &analysis.Pass{
		Analyzer:  a,
		Fset:      fset,
		Files:     nonTestFiles,
		Pkg:       pkg.Types,
		TypesInfo: pkg.TypesInfo,
		ResultOf:  results,
		Report: func(diag analysis.Diagnostic) {
			*diagnostics = append(*diagnostics, diagWithFset{
				diag: diag,
				fset: fset,
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
	fmt := formatter.NewFormatter(os.Stdout, AIMode, NoColor, Simple)

	// Vérification de la condition
	if len(diagnostics) == 0 {
		fmt.Format(nil, nil)
		// Early return from function.
		return
	}

	firstFset := diagnostics[0].fset
	diags := extractDiagnostics(diagnostics)
	fmt.Format(firstFset, diags)
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
	seen := make(map[string]bool)
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

	diags := make([]analysis.Diagnostic, len(deduped))
	// Itération sur les éléments
	for i, d := range deduped {
		diags[i] = d.diag
	}
	// Early return from function.
	return diags
}
