// Analysis runner for executing analyzers on packages.
package orchestrator

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// AnalysisRunner handles running analyzers on packages.
// Manages file selection, required analyzer execution, and pass creation.
type AnalysisRunner struct {
	stderr  io.Writer
	verbose bool
}

// NewAnalysisRunner creates a new AnalysisRunner.
//
// Params:
//   - stderr: writer for verbose output
//   - verbose: enable verbose logging
//
// Returns:
//   - *AnalysisRunner: new runner instance
func NewAnalysisRunner(stderr io.Writer, verbose bool) *AnalysisRunner {
	// Return new runner instance
	return &AnalysisRunner{
		stderr:  stderr,
		verbose: verbose,
	}
}

// Run executes all analyzers on the given packages.
//
// Params:
//   - pkgs: packages to analyze
//   - analyzers: analyzers to run
//
// Returns:
//   - []DiagnosticResult: collected diagnostics
func (r *AnalysisRunner) Run(pkgs []*packages.Package, analyzers []*analysis.Analyzer) []DiagnosticResult {
	var allDiagnostics []DiagnosticResult
	results := make(map[*analysis.Analyzer]any, len(analyzers))

	// Analyze each package
	for _, pkg := range pkgs {
		r.analyzePackage(pkg, analyzers, results, &allDiagnostics)
	}

	// Return collected diagnostics
	return allDiagnostics
}

// analyzePackage analyzes a single package with the given analyzers.
//
// Params:
//   - pkg: package to analyze
//   - analyzers: analyzers to run
//   - results: analyzer results map (modified in-place)
//   - diagnostics: diagnostics slice (modified in-place)
func (r *AnalysisRunner) analyzePackage(
	pkg *packages.Package,
	analyzers []*analysis.Analyzer,
	results map[*analysis.Analyzer]any,
	diagnostics *[]DiagnosticResult,
) {
	pkgFset := pkg.Fset

	// Log if verbose
	if r.verbose {
		fmt.Fprintf(r.stderr, "Analyzing package: %s\n", pkg.PkgPath)
	}

	// Clear results for this package
	for k := range results {
		delete(results, k)
	}

	// Run each analyzer
	for _, a := range analyzers {
		pass := r.createPass(a, pkg, pkgFset, diagnostics, results)
		result, err := a.Run(pass)

		// Handle errors
		if err != nil {
			fmt.Fprintf(r.stderr, "Error running analyzer %s on %s: %v\n", a.Name, pkg.PkgPath, err)
		}

		// Store result
		results[a] = result
	}
}

// createPass creates an analysis pass for a package.
//
// Params:
//   - a: analyzer to run
//   - pkg: package to analyze
//   - fset: fileset for positions
//   - diagnostics: diagnostics slice for collecting results
//   - results: required analyzer results
//
// Returns:
//   - *analysis.Pass: created pass
func (r *AnalysisRunner) createPass(
	a *analysis.Analyzer,
	pkg *packages.Package,
	fset *token.FileSet,
	diagnostics *[]DiagnosticResult,
	results map[*analysis.Analyzer]any,
) *analysis.Pass {
	files := r.selectFiles(a, pkg, fset)
	r.runRequired(a, files, pkg, fset, results)

	// Return created pass
	return &analysis.Pass{
		Analyzer:  a,
		Fset:      fset,
		Files:     files,
		Pkg:       pkg.Types,
		TypesInfo: pkg.TypesInfo,
		ResultOf:  results,
		Report: func(diag analysis.Diagnostic) {
			*diagnostics = append(*diagnostics, DiagnosticResult{
				Diag:         diag,
				Fset:         fset,
				AnalyzerName: a.Name,
			})
		},
		ReadFile: func(filename string) ([]byte, error) {
			// Return file content
			return os.ReadFile(filename)
		},
	}
}

// selectFiles determines which files to analyze for an analyzer.
//
// Params:
//   - a: analyzer
//   - pkg: package
//   - fset: fileset
//
// Returns:
//   - []*ast.File: files to analyze
func (r *AnalysisRunner) selectFiles(a *analysis.Analyzer, pkg *packages.Package, fset *token.FileSet) []*ast.File {
	// Test analyzers need all files
	if strings.HasPrefix(a.Name, "ktntest") {
		// Return all files
		return pkg.Syntax
	}

	// Check force mode
	if config.Get().ForceAllRulesOnTests {
		// Return all files
		return pkg.Syntax
	}

	// Filter test files for other analyzers
	return r.filterTestFiles(pkg.Syntax, fset)
}

// filterTestFiles filters out test files.
//
// Params:
//   - files: files to filter
//   - fset: fileset for position
//
// Returns:
//   - []*ast.File: filtered files
func (r *AnalysisRunner) filterTestFiles(files []*ast.File, fset *token.FileSet) []*ast.File {
	filtered := make([]*ast.File, 0, len(files))
	// Iterate over files
	for _, file := range files {
		pos := fset.Position(file.Pos())
		// Skip test files
		if !strings.HasSuffix(pos.Filename, "_test.go") {
			filtered = append(filtered, file)
		}
	}
	// Return filtered files
	return filtered
}

// runRequired runs required analyzers first.
//
// Params:
//   - a: analyzer with requirements
//   - files: files to analyze
//   - pkg: package
//   - fset: fileset
//   - results: results map (modified in-place)
func (r *AnalysisRunner) runRequired(
	a *analysis.Analyzer,
	files []*ast.File,
	pkg *packages.Package,
	fset *token.FileSet,
	results map[*analysis.Analyzer]any,
) {
	// Run required analyzers
	for _, req := range a.Requires {
		reqPass := &analysis.Pass{
			Analyzer:  req,
			Fset:      fset,
			Files:     files,
			Pkg:       pkg.Types,
			TypesInfo: pkg.TypesInfo,
			ResultOf:  results,
			Report:    func(analysis.Diagnostic) {},
			ReadFile: func(filename string) ([]byte, error) {
				// Return file content
				return os.ReadFile(filename)
			},
		}
		result, _ := req.Run(reqPass)
		results[req] = result
	}
}
