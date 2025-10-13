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

var (
	aiMode  = flag.Bool("ai", false, "Enable AI-friendly output format")
	noColor = flag.Bool("no-color", false, "Disable colored output")
	simple  = flag.Bool("simple", false, "Simple one-line format for IDE integration")
	verbose = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: ktn-linter [flags] <packages>\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  ktn-linter ./...\n")
		fmt.Fprintf(os.Stderr, "  ktn-linter ./path/to/file.go\n")
		fmt.Fprintf(os.Stderr, "  ktn-linter -ai ./...\n")
		os.Exit(1)
	}

	// Charger les packages
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading packages: %v\n", err)
		os.Exit(1)
	}

	// Vérifier les erreurs de chargement
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

	// Structure pour stocker diagnostic avec son FileSet
	type diagWithFset struct {
		diag analysis.Diagnostic
		fset *token.FileSet
	}

	// Liste des analyseurs à exécuter
	analyzers := []*analysis.Analyzer{
		analyzer.ConstAnalyzer,
		analyzer.VarAnalyzer,
		analyzer.FuncAnalyzer,
	}

	// Exécuter l'analyse
	var allDiagnostics []diagWithFset

	for _, pkg := range pkgs {
		pkgFset := pkg.Fset // Capturer le FileSet du package

		if *verbose {
			fmt.Fprintf(os.Stderr, "Analyzing package: %s\n", pkg.PkgPath)
		}

		// Exécuter chaque analyseur
		for _, a := range analyzers {
			pass := &analysis.Pass{
				Analyzer: a,
				Fset:     pkgFset,
				Files:    pkg.Syntax,
				Pkg:      pkg.Types,
				TypesInfo: pkg.TypesInfo,
				Report: func(diag analysis.Diagnostic) {
					allDiagnostics = append(allDiagnostics, diagWithFset{
						diag: diag,
						fset: pkgFset,
					})
				},
			}

			_, err := a.Run(pass)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running analyzer %s on %s: %v\n", a.Name, pkg.PkgPath, err)
				os.Exit(1)
			}
		}
	}

	// Formater et afficher les résultats
	if len(allDiagnostics) == 0 {
		fmt := formatter.NewFormatter(os.Stdout, *aiMode, *noColor, *simple)
		fmt.Format(nil, nil)
	} else {
		// Utiliser le premier FileSet comme référence (tous devraient pointer vers les mêmes fichiers)
		firstFset := allDiagnostics[0].fset

		// Extraire juste les diagnostics
		diags := make([]analysis.Diagnostic, len(allDiagnostics))
		for i, d := range allDiagnostics {
			diags[i] = d.diag
		}

		fmt := formatter.NewFormatter(os.Stdout, *aiMode, *noColor, *simple)
		fmt.Format(firstFset, diags)
	}

	// Exit code
	if len(allDiagnostics) > 0 {
		os.Exit(1)
	}
}
