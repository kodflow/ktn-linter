// Package orchestrator coordinates the linting pipeline.
package orchestrator

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// diagChannelBufferMultiplier is the buffer multiplier for diagnostic channels.
const diagChannelBufferMultiplier int = 10

// waitGroup defines the interface for wait group operations.
type waitGroup interface {
	Done()
}

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

// Run executes all analyzers on the given packages in parallel.
// Uses worker pool limited by GOMAXPROCS for concurrency control.
//
// Params:
//   - pkgs: packages to analyze
//   - analyzers: analyzers to run
//
// Returns:
//   - []DiagnosticResult: collected diagnostics
func (r *AnalysisRunner) Run(pkgs []*packages.Package, analyzers []*analysis.Analyzer) []DiagnosticResult {
	// Use channel for concurrent-safe diagnostic collection
	diagChan := make(chan DiagnosticResult, len(pkgs)*diagChannelBufferMultiplier)
	var wg sync.WaitGroup

	// Limit concurrent workers to GOMAXPROCS
	workerCount := runtime.GOMAXPROCS(0)
	pkgChan := make(chan *packages.Package, len(pkgs))

	// Start workers (one goroutine per available CPU)
	for range workerCount {
		wg.Add(1)
		go r.worker(analyzers, pkgChan, diagChan, &wg)
	}

	// Send packages to workers
	for _, pkg := range pkgs {
		pkgChan <- pkg
	}
	close(pkgChan)

	// Wait for workers and close results channel
	go func() {
		wg.Wait()
		close(diagChan)
	}()

	// Collect all diagnostics in result slice
	var allDiagnostics []DiagnosticResult
	// Range over channel until closed by background goroutine
	for diag := range diagChan {
		allDiagnostics = append(allDiagnostics, diag)
	}

	// Return collected diagnostics
	return allDiagnostics
}

// worker processes packages from pkgChan and sends diagnostics to diagChan.
//
// Params:
//   - analyzers: analyzers to run
//   - pkgChan: channel receiving packages to analyze
//   - diagChan: channel for sending diagnostics
//   - wg: wait group to signal completion
func (r *AnalysisRunner) worker(
	analyzers []*analysis.Analyzer,
	pkgChan <-chan *packages.Package,
	diagChan chan<- DiagnosticResult,
	wg waitGroup,
) {
	defer wg.Done()
	// Each worker has its own results map
	results := make(map[*analysis.Analyzer]any, len(analyzers))

	// Process packages from channel
	for pkg := range pkgChan {
		r.analyzePackageParallel(pkg, analyzers, results, diagChan)
	}
}

// analyzePackageParallel analyzes a package and sends diagnostics to a channel.
//
// Params:
//   - pkg: package to analyze
//   - analyzers: analyzers to run
//   - results: analyzer results map (modified in-place)
//   - diagChan: channel for sending diagnostics
func (r *AnalysisRunner) analyzePackageParallel(
	pkg *packages.Package,
	analyzers []*analysis.Analyzer,
	results map[*analysis.Analyzer]any,
	diagChan chan<- DiagnosticResult,
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
		pass := r.createPassParallel(a, pkg, pkgFset, diagChan, results)
		result, err := a.Run(pass)

		// Handle errors
		if err != nil {
			fmt.Fprintf(r.stderr, "Error running analyzer %s on %s: %v\n", a.Name, pkg.PkgPath, err)
		}

		// Store result
		results[a] = result
	}
}

// createPassParallel creates an analysis pass for parallel execution.
// Sends diagnostics to a channel instead of appending to a slice.
//
// Params:
//   - a: analyzer to run
//   - pkg: package to analyze
//   - fset: fileset for positions
//   - diagChan: channel for sending diagnostics
//   - results: required analyzer results
//
// Returns:
//   - *analysis.Pass: created pass
func (r *AnalysisRunner) createPassParallel(
	a *analysis.Analyzer,
	pkg *packages.Package,
	fset *token.FileSet,
	diagChan chan<- DiagnosticResult,
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
			diagChan <- DiagnosticResult{
				Diag:         diag,
				Fset:         fset,
				AnalyzerName: a.Name,
			}
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
// Caches results to avoid re-running the same analyzer multiple times per package.
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
		// Skip if already computed for this package
		if _, exists := results[req]; exists {
			continue
		}

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
