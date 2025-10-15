package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"github.com/kodflow/ktn-linter/src/pkg/formatter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// Options de ligne de commande pour configurer le comportement du linter
var (
	// aiMode active le format de sortie optimisé pour les IA
	aiMode bool
	// noColor désactive les couleurs dans la sortie
	noColor bool
	// simple active le format une-ligne pour l'intégration IDE
	simple bool
	// verbose active les logs détaillés
	verbose bool
)

// diagWithFset associe un diagnostic avec son FileSet
type diagWithFset struct {
	diag analysis.Diagnostic
	fset *token.FileSet
}

// main est le point d'entrée du linter KTN
func main() {
	parseFlags()

	if len(flag.Args()) == 0 {
		printUsage()
		os.Exit(1)
	}

	pkgs := loadPackages(flag.Args())
	diagnostics := runAnalyzers(pkgs)
	formatAndDisplay(diagnostics)

	if len(diagnostics) > 0 {
		os.Exit(1)
	}
}

// parseFlags configure les flags de ligne de commande
func parseFlags() {
	flag.BoolVar(&aiMode, "ai", false, "Enable AI-friendly output format")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&simple, "simple", false, "Simple one-line format for IDE integration")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()
}

// printUsage affiche l'aide d'utilisation du linter
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: ktn-linter [flags] <packages>\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  ktn-linter ./...\n")
	fmt.Fprintf(os.Stderr, "  ktn-linter ./path/to/file.go\n")
	fmt.Fprintf(os.Stderr, "  ktn-linter -ai ./...\n")
}

// loadPackages charge les packages Go depuis les patterns fournis.
//
// Params:
//   - patterns: les patterns de packages à charger
//
// Returns:
//   - []*packages.Package: la liste des packages chargés ou quitte en cas d'erreur
func loadPackages(patterns []string) []*packages.Package {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true, // Charger les fichiers *_test.go pour KTN-TEST rules
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading packages: %v\n", err)
		os.Exit(1)
	}

	checkLoadErrors(pkgs)
	// Retourne les packages chargés avec succès
	return pkgs
}

// checkLoadErrors vérifie et affiche les erreurs de chargement.
//
// Params:
//   - pkgs: la liste des packages à vérifier
func checkLoadErrors(pkgs []*packages.Package) {
	var hasLoadErrors bool
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			hasLoadErrors = true
		}
	}
	if hasLoadErrors {
		os.Exit(1)
	}
}

// runAnalyzers exécute tous les analyseurs sur les packages.
//
// Params:
//   - pkgs: la liste des packages à analyser
//
// Returns:
//   - []diagWithFset: la liste des diagnostics trouvés avec leurs FileSets
func runAnalyzers(pkgs []*packages.Package) []diagWithFset {
	analyzers := []*analysis.Analyzer{
		analyzer.ConstAnalyzer,
		analyzer.VarAnalyzer,
		analyzer.FuncAnalyzer,
		analyzer.InterfaceAnalyzer,
		analyzer.InterfaceStrictAnalyzer,
		analyzer.TestAnalyzer,
	}

	var allDiagnostics []diagWithFset

	for _, pkg := range pkgs {
		pkgFset := pkg.Fset

		if verbose {
			fmt.Fprintf(os.Stderr, "Analyzing package: %s\n", pkg.PkgPath)
		}

		for _, a := range analyzers {
			pass := createAnalysisPass(a, pkg, pkgFset, &allDiagnostics)

			if _, err := a.Run(pass); err != nil {
				fmt.Fprintf(os.Stderr, "Error running analyzer %s on %s: %v\n", a.Name, pkg.PkgPath, err)
				os.Exit(1)
			}
		}
	}

	// Retourne tous les diagnostics collectés depuis tous les packages
	return allDiagnostics
}

// createAnalysisPass crée une passe d'analyse pour un analyseur.
//
// Params:
//   - a: l'analyseur à exécuter
//   - pkg: le package à analyser
//   - fset: le FileSet pour les positions
//   - diagnostics: le pointeur vers la liste des diagnostics collectés
//
// Returns:
//   - *analysis.Pass: la passe d'analyse configurée
func createAnalysisPass(a *analysis.Analyzer, pkg *packages.Package, fset *token.FileSet, diagnostics *[]diagWithFset) *analysis.Pass {
	// Retourne une passe d'analyse configurée avec le callback Report
	return &analysis.Pass{
		Analyzer:  a,
		Fset:      fset,
		Files:     pkg.Syntax,
		Pkg:       pkg.Types,
		TypesInfo: pkg.TypesInfo,
		Report: func(diag analysis.Diagnostic) {
			*diagnostics = append(*diagnostics, diagWithFset{
				diag: diag,
				fset: fset,
			})
		},
	}
}

// formatAndDisplay formate et affiche les diagnostics.
//
// Params:
//   - diagnostics: la liste des diagnostics à formater et afficher
func formatAndDisplay(diagnostics []diagWithFset) {
	fmt := formatter.NewFormatter(os.Stdout, aiMode, noColor, simple)

	if len(diagnostics) == 0 {
		fmt.Format(nil, nil)
		// Retourne car il n'y a aucun diagnostic à afficher
		return
	}

	firstFset := diagnostics[0].fset
	diags := extractDiagnostics(diagnostics)
	fmt.Format(firstFset, diags)
}

// extractDiagnostics extrait la liste des diagnostics depuis diagWithFset.
//
// Params:
//   - diagnostics: la liste des diagWithFset à extraire
//
// Returns:
//   - []analysis.Diagnostic: la liste des diagnostics extraits
func extractDiagnostics(diagnostics []diagWithFset) []analysis.Diagnostic {
	diags := make([]analysis.Diagnostic, len(diagnostics))
	for i, d := range diagnostics {
		diags[i] = d.diag
	}
	// Retourne la liste des diagnostics extraits
	return diags
}
