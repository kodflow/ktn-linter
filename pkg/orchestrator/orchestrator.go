// Package orchestrator provides the main orchestration for linting Go packages.
package orchestrator

import (
	"io"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// Orchestrator coordinates the linting process.
// It manages package loading, analyzer selection, execution, and diagnostics processing.
type Orchestrator struct {
	loader     *PackageLoader
	selector   *AnalyzerSelector
	runner     *AnalysisRunner
	processor  *DiagnosticsProcessor
	stderr     io.Writer
	verbose    bool
}

// NewOrchestrator creates a new Orchestrator.
//
// Params:
//   - stderr: writer for error/verbose output
//   - verbose: enable verbose logging
//
// Returns:
//   - *Orchestrator: new orchestrator instance
func NewOrchestrator(stderr io.Writer, verbose bool) *Orchestrator {
	// Return new orchestrator
	return &Orchestrator{
		loader:    NewPackageLoader(stderr),
		selector:  NewAnalyzerSelector(stderr, verbose),
		runner:    NewAnalysisRunner(stderr, verbose),
		processor: NewDiagnosticsProcessor(),
		stderr:    stderr,
		verbose:   verbose,
	}
}

// LoadPackages loads Go packages from patterns.
//
// Params:
//   - patterns: package patterns to load
//
// Returns:
//   - []*packages.Package: loaded packages
//   - error: loading error if any
func (o *Orchestrator) LoadPackages(patterns []string) ([]*packages.Package, error) {
	// Delegate to loader
	return o.loader.Load(patterns)
}

// SelectAnalyzers selects analyzers based on options.
//
// Params:
//   - opts: selection options
//
// Returns:
//   - []*analysis.Analyzer: selected analyzers
//   - error: selection error if any
func (o *Orchestrator) SelectAnalyzers(opts Options) ([]*analysis.Analyzer, error) {
	// Delegate to selector
	return o.selector.Select(opts)
}

// RunAnalyzers runs analyzers on packages.
//
// Params:
//   - pkgs: packages to analyze
//   - analyzers: analyzers to run
//
// Returns:
//   - []DiagnosticResult: collected diagnostics
func (o *Orchestrator) RunAnalyzers(pkgs []*packages.Package, analyzers []*analysis.Analyzer) []DiagnosticResult {
	// Delegate to runner
	return o.runner.Run(pkgs, analyzers)
}

// FilterDiagnostics filters out cache/tmp diagnostics.
//
// Params:
//   - diagnostics: raw diagnostics
//
// Returns:
//   - []DiagnosticResult: filtered diagnostics
func (o *Orchestrator) FilterDiagnostics(diagnostics []DiagnosticResult) []DiagnosticResult {
	// Delegate to processor
	return o.processor.Filter(diagnostics)
}

// ExtractDiagnostics extracts and deduplicates diagnostics.
//
// Params:
//   - diagnostics: raw diagnostics
//
// Returns:
//   - []analysis.Diagnostic: processed diagnostics
func (o *Orchestrator) ExtractDiagnostics(diagnostics []DiagnosticResult) []analysis.Diagnostic {
	// Delegate to processor
	return o.processor.Extract(diagnostics)
}

// Run executes the full linting pipeline.
//
// Params:
//   - patterns: package patterns to analyze
//   - opts: linting options
//
// Returns:
//   - []analysis.Diagnostic: found issues
//   - error: pipeline error if any
func (o *Orchestrator) Run(patterns []string, opts Options) ([]analysis.Diagnostic, error) {
	// Load packages
	pkgs, err := o.LoadPackages(patterns)
	// Check for error
	if err != nil {
		// Return error
		return []analysis.Diagnostic{}, err
	}

	// Select analyzers
	analyzers, err := o.SelectAnalyzers(opts)
	// Check for error
	if err != nil {
		// Return error
		return []analysis.Diagnostic{}, err
	}

	// Run analyzers
	rawDiags := o.RunAnalyzers(pkgs, analyzers)

	// Filter diagnostics
	filtered := o.FilterDiagnostics(rawDiags)

	// Extract and deduplicate
	diags := o.ExtractDiagnostics(filtered)

	// Return diagnostics
	return diags, nil
}

// GetFirstFset returns the FileSet from the first diagnostic.
//
// Params:
//   - diagnostics: diagnostics with fsets
//
// Returns:
//   - any: first fset or nil
func GetFirstFset(diagnostics []DiagnosticResult) any {
	// Check if empty
	if len(diagnostics) == 0 {
		// Return nil
		return nil
	}
	// Return first fset
	return diagnostics[0].Fset
}
