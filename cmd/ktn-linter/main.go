package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
	"github.com/kodflow/ktn-linter/pkg/formatter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// Options de ligne de commande
var (
	// aiMode enables AI-friendly output format.
	aiMode bool
	// noColor disables colored output.
	noColor bool
	// simple enables simple one-line format for IDE integration.
	simple bool
	// verbose enables verbose output.
	verbose bool
	// category filters rules by specific category.
	category string
)

// osExit est une variable pour permettre le mocking dans les tests.
// Par défaut, elle pointe vers os.Exit, mais peut être remplacée par un mock.
var osExit = os.Exit

// diagWithFset associe un diagnostic avec son FileSet
type diagWithFset struct {
	diag analysis.Diagnostic
	fset *token.FileSet
}

// main est le point d'entrée du linter KTN.
// Returns: aucun
//
// Params: aucun
func main() {
	parseFlags()

	// Vérification de la condition
	if len(flag.Args()) == 0 {
		printUsage()
		osExit(1)
	}

	pkgs := loadPackages(flag.Args())
	diagnostics := runAnalyzers(pkgs)

	// Filter out diagnostics from cache/tmp files (same logic as formatter)
	filteredDiags := filterDiagnostics(diagnostics)

	formatAndDisplay(filteredDiags)

	// Vérification de la condition
	if len(filteredDiags) > 0 {
		osExit(1)
	}

	// Success - exit with 0
	osExit(0)
}

// parseFlags analyse les drapeaux de ligne de commande.
// Returns: aucun
//
// Params: aucun
func parseFlags() {
	flag.BoolVar(&aiMode, "ai", false, "Enable AI-friendly output format")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&simple, "simple", false, "Simple one-line format for IDE integration")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.StringVar(&category, "category", "", "Run only rules from specific category (func, var, error, etc.)")
	flag.Parse()
}

// printUsage affiche l'aide d'utilisation du linter.
// Returns: aucun
//
// Params: aucun
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: ktn-linter [flags] <packages>\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nCategories disponibles:\n")
	fmt.Fprintf(os.Stderr, "  func, var, struct, interface, const, error, test\n")
	fmt.Fprintf(os.Stderr, "  alloc, goroutine, pool, mock, method, package\n")
	fmt.Fprintf(os.Stderr, "  control_flow, data_structures, ops\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  ktn-linter ./...\n")
	fmt.Fprintf(os.Stderr, "  ktn-linter -category=error ./...\n")
	fmt.Fprintf(os.Stderr, "  ktn-linter -ai ./path/to/file.go\n")
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
		osExit(1)
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
		osExit(1)
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
	if category != "" {
		analyzers = ktn.GetRulesByCategory(category)
		// Vérification de la condition
		if analyzers == nil {
			fmt.Fprintf(os.Stderr, "Unknown category: %s\n", category)
			osExit(1)
		}
		// Vérification de la condition
		if verbose {
			fmt.Fprintf(os.Stderr, "Running %d rules from category '%s'\n", len(analyzers), category)
		}
		// Cas alternatif
	} else {
		analyzers = ktn.GetAllRules()
		// Vérification de la condition
		if verbose {
			fmt.Fprintf(os.Stderr, "Running all %d KTN rules\n", len(analyzers))
		}
	}

	// allDiagnostics holds the configuration value.

	var allDiagnostics []diagWithFset

	// Itération sur les éléments
	for _, pkg := range pkgs {
		pkgFset := pkg.Fset

		// Vérification de la condition
		if verbose {
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
	fmt := formatter.NewFormatter(os.Stdout, aiMode, noColor, simple)

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
