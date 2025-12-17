// Analyzer selector for choosing which analyzers to run.
package orchestrator

import (
	"fmt"
	"io"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerSelector handles selecting analyzers based on options.
// Supports filtering by category, single rule, or returning all rules.
type AnalyzerSelector struct {
	stderr  io.Writer
	verbose bool
}

// NewAnalyzerSelector creates a new AnalyzerSelector.
//
// Params:
//   - stderr: writer for verbose output
//   - verbose: enable verbose logging
//
// Returns:
//   - *AnalyzerSelector: new selector instance
func NewAnalyzerSelector(stderr io.Writer, verbose bool) *AnalyzerSelector {
	// Return new selector instance
	return &AnalyzerSelector{
		stderr:  stderr,
		verbose: verbose,
	}
}

// Select selects analyzers based on the given options.
//
// Params:
//   - opts: selection options
//
// Returns:
//   - []*analysis.Analyzer: selected analyzers
//   - error: selection error if any
func (s *AnalyzerSelector) Select(opts Options) ([]*analysis.Analyzer, error) {
	// Select by single rule
	if opts.OnlyRule != "" {
		// Return single rule
		return s.selectSingleRule(opts.OnlyRule)
	}

	// Select by category
	if opts.Category != "" {
		// Return category rules
		return s.selectByCategory(opts.Category)
	}

	// Default: all rules
	return s.selectAll()
}

// selectSingleRule selects a single rule by code.
//
// Params:
//   - code: rule code to select
//
// Returns:
//   - []*analysis.Analyzer: single analyzer
//   - error: if rule not found
func (s *AnalyzerSelector) selectSingleRule(code string) ([]*analysis.Analyzer, error) {
	analyzer := ktn.GetRuleByCode(code)
	// Check if rule exists
	if analyzer == nil {
		// Return error
		return []*analysis.Analyzer{}, fmt.Errorf("unknown rule code: %s", code)
	}

	// Log if verbose
	if s.verbose {
		fmt.Fprintf(s.stderr, "Running only rule '%s'\n", code)
	}

	// Return single analyzer
	return []*analysis.Analyzer{analyzer}, nil
}

// selectByCategory selects all rules from a category.
//
// Params:
//   - category: category name
//
// Returns:
//   - []*analysis.Analyzer: category analyzers
//   - error: if category not found
func (s *AnalyzerSelector) selectByCategory(category string) ([]*analysis.Analyzer, error) {
	analyzers := ktn.GetRulesByCategory(category)
	// Check if category exists
	if len(analyzers) == 0 {
		// Return error
		return []*analysis.Analyzer{}, fmt.Errorf("unknown category: %s", category)
	}

	// Log if verbose
	if s.verbose {
		fmt.Fprintf(s.stderr, "Running %d rules from category '%s'\n", len(analyzers), category)
	}

	// Return category analyzers
	return analyzers, nil
}

// selectAll selects all available rules.
//
// Returns:
//   - []*analysis.Analyzer: all analyzers
//   - error: always nil
func (s *AnalyzerSelector) selectAll() ([]*analysis.Analyzer, error) {
	analyzers := ktn.GetAllRules()

	// Log if verbose
	if s.verbose {
		fmt.Fprintf(s.stderr, "Running all %d KTN rules\n", len(analyzers))
	}

	// Return all analyzers
	return analyzers, nil
}
