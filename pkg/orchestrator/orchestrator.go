// Package orchestrator provides the main orchestration for linting Go packages.
package orchestrator

import (
	"fmt"
	"io"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// Orchestrator coordinates the linting process.
// It manages package loading, analyzer selection, execution, and diagnostics processing.
type Orchestrator struct {
	loader    *PackageLoader
	selector  *AnalyzerSelector
	runner    *AnalysisRunner
	processor *DiagnosticsProcessor
	discovery *ModuleDiscovery
	stderr    io.Writer
	verbose   bool
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
		discovery: NewModuleDiscovery(),
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

// DiscoverModules finds all Go modules in paths.
//
// Params:
//   - paths: paths to search
//
// Returns:
//   - []string: module root directories
//   - error: discovery error if any
func (o *Orchestrator) DiscoverModules(paths []string) ([]string, error) {
	// Delegate to discovery
	return o.discovery.FindModules(paths)
}

// LoadPackagesFromDir loads packages from a specific directory.
//
// Params:
//   - dir: module root directory
//   - patterns: package patterns
//
// Returns:
//   - []*packages.Package: loaded packages
//   - error: loading error if any
func (o *Orchestrator) LoadPackagesFromDir(dir string, patterns []string) ([]*packages.Package, error) {
	// Delegate to loader
	return o.loader.LoadFromDir(dir, patterns)
}

// RunMultiModule runs analysis across multiple modules.
//
// Params:
//   - paths: paths to analyze (may contain multiple modules)
//   - opts: linting options
//
// Returns:
//   - []DiagnosticResult: aggregated diagnostics
//   - error: pipeline error if any
func (o *Orchestrator) RunMultiModule(paths []string, opts Options) ([]DiagnosticResult, error) {
	// Discover modules
	modules, err := o.DiscoverModules(paths)
	// Check for error
	if err != nil {
		return nil, fmt.Errorf("discovering modules: %w", err)
	}

	// Check if no modules found
	if len(modules) == 0 {
		// Fall back to standard load from current directory
		return o.runSingleModule("", paths, opts)
	}

	// Log if verbose
	if o.verbose {
		fmt.Fprintf(o.stderr, "Found %d Go module(s)\n", len(modules))
	}

	// Select analyzers once
	analyzers, err := o.SelectAnalyzers(opts)
	// Check for error
	if err != nil {
		return nil, err
	}

	// Aggregate results
	var allDiags []DiagnosticResult

	// Process each module
	for _, moduleRoot := range modules {
		// Log if verbose
		if o.verbose {
			fmt.Fprintf(o.stderr, "Analyzing module: %s\n", moduleRoot)
		}

		// Get patterns for this module
		patterns := o.discovery.ResolvePatterns(moduleRoot, paths)

		// Load packages from module directory
		pkgs, err := o.LoadPackagesFromDir(moduleRoot, patterns)
		// Check for error
		if err != nil {
			// Log warning and continue
			if o.verbose {
				fmt.Fprintf(o.stderr, "Warning: %v\n", err)
			}
			continue
		}

		// Run analyzers
		diags := o.RunAnalyzers(pkgs, analyzers)
		allDiags = append(allDiags, diags...)
	}

	return allDiags, nil
}

// runSingleModule runs analysis for a single module.
//
// Params:
//   - dir: module directory (empty for current)
//   - patterns: package patterns
//   - opts: linting options
//
// Returns:
//   - []DiagnosticResult: collected diagnostics
//   - error: pipeline error if any
func (o *Orchestrator) runSingleModule(dir string, patterns []string, opts Options) ([]DiagnosticResult, error) {
	// Load packages
	pkgs, err := o.LoadPackagesFromDir(dir, patterns)
	// Check for error
	if err != nil {
		return nil, err
	}

	// Select analyzers
	analyzers, err := o.SelectAnalyzers(opts)
	// Check for error
	if err != nil {
		return nil, err
	}

	// Run analyzers
	return o.RunAnalyzers(pkgs, analyzers), nil
}
