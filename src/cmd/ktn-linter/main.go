package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn"
	"github.com/kodflow/ktn-linter/src/pkg/formatter"
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

// diagWithFset associe un diagnostic avec son FileSet
type diagWithFset struct {
	diag analysis.Diagnostic
	fset *token.FileSet
}

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

func parseFlags() {
	flag.BoolVar(&aiMode, "ai", false, "Enable AI-friendly output format")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&simple, "simple", false, "Simple one-line format for IDE integration")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.StringVar(&category, "category", "", "Run only rules from specific category (func, var, error, etc.)")
	flag.Parse()
}

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

func loadPackages(patterns []string) []*packages.Package {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading packages: %v\n", err)
		os.Exit(1)
	}

	checkLoadErrors(pkgs)
	// Early return from function.
	return pkgs
}

func checkLoadErrors(pkgs []*packages.Package) {
	// hasLoadErrors holds the configuration value.

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

func runAnalyzers(pkgs []*packages.Package) []diagWithFset {
	// analyzers holds the configuration value.

	var analyzers []*analysis.Analyzer

	// Sélectionner les analyseurs selon la catégorie
	if category != "" {
		analyzers = ktn.GetRulesByCategory(category)
		if analyzers == nil {
			fmt.Fprintf(os.Stderr, "Unknown category: %s\n", category)
			os.Exit(1)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Running %d rules from category '%s'\n", len(analyzers), category)
		}
	} else {
		analyzers = ktn.GetAllRules()
		if verbose {
			fmt.Fprintf(os.Stderr, "Running all %d KTN rules\n", len(analyzers))
		}
	}

	// allDiagnostics holds the configuration value.

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
			}
		}
	}

	// Early return from function.
	return allDiagnostics
}

func createAnalysisPass(a *analysis.Analyzer, pkg *packages.Package, fset *token.FileSet, diagnostics *[]diagWithFset) *analysis.Pass {
	// Early return from function.
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

func formatAndDisplay(diagnostics []diagWithFset) {
	fmt := formatter.NewFormatter(os.Stdout, aiMode, noColor, simple)

	if len(diagnostics) == 0 {
		fmt.Format(nil, nil)
		// Early return from function.
		return
	}

	firstFset := diagnostics[0].fset
	diags := extractDiagnostics(diagnostics)
	fmt.Format(firstFset, diags)
}

func extractDiagnostics(diagnostics []diagWithFset) []analysis.Diagnostic {
	diags := make([]analysis.Diagnostic, len(diagnostics))
	for i, d := range diagnostics {
		diags[i] = d.diag
	}
	// Early return from function.
	return diags
}
