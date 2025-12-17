// Package orchestrator coordinates the linting pipeline.
package orchestrator

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// DiagnosticsProcessor handles filtering and processing diagnostics.
// Provides deduplication, cache file filtering, and modernize prefix addition.
type DiagnosticsProcessor struct{}

// NewDiagnosticsProcessor creates a new DiagnosticsProcessor.
//
// Returns:
//   - *DiagnosticsProcessor: new processor instance
func NewDiagnosticsProcessor() *DiagnosticsProcessor {
	// Return new processor instance
	return &DiagnosticsProcessor{}
}

// Filter filters out diagnostics from cache/tmp files.
//
// Params:
//   - diagnostics: raw diagnostics
//
// Returns:
//   - []DiagnosticResult: filtered diagnostics
func (p *DiagnosticsProcessor) Filter(diagnostics []DiagnosticResult) []DiagnosticResult {
	var filtered []DiagnosticResult
	// Iterate over diagnostics
	for i := range diagnostics {
		pos := diagnostics[i].Position()
		// Skip Go build cache files
		if strings.Contains(pos.Filename, "/.cache/go-build/") ||
			strings.Contains(pos.Filename, "\\cache\\go-build\\") {
			continue
		}
		filtered = append(filtered, diagnostics[i])
	}
	// Return filtered diagnostics
	return filtered
}

// Extract extracts and deduplicates diagnostics.
//
// Params:
//   - diagnostics: raw diagnostics with fset
//
// Returns:
//   - []analysis.Diagnostic: deduplicated diagnostics
func (p *DiagnosticsProcessor) Extract(diagnostics []DiagnosticResult) []analysis.Diagnostic {
	// Deduplicate diagnostics
	seen := make(map[string]bool, len(diagnostics))
	var deduped []DiagnosticResult

	// Iterate over diagnostics
	for i := range diagnostics {
		pos := diagnostics[i].Position()
		key := fmt.Sprintf("%s:%d:%d:%s", pos.Filename, pos.Line, pos.Column, diagnostics[i].Diag.Message)
		// Skip duplicates
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, diagnostics[i])
		}
	}

	// Build result slice
	diags := make([]analysis.Diagnostic, 0, len(deduped))
	// Iterate over deduplicated
	for _, d := range deduped {
		diag := d.Diag
		// Prefix modernize messages
		if p.isModernize(d.AnalyzerName) && !strings.HasPrefix(diag.Message, "KTN-") {
			code := p.formatModernizeCode(d.AnalyzerName)
			diag.Message = code + ": " + diag.Message
		}
		diags = append(diags, diag)
	}

	// Return processed diagnostics
	return diags
}

// isModernize checks if an analyzer is a modernize analyzer.
//
// Params:
//   - name: analyzer name
//
// Returns:
//   - bool: true if modernize analyzer
func (p *DiagnosticsProcessor) isModernize(name string) bool {
	modernizeAnalyzers := map[string]bool{
		"any":              true,
		"bloop":            true,
		"fmtappendf":       true,
		"forvar":           true,
		"mapsloop":         true,
		"minmax":           true,
		"newexpr":          true,
		"omitzero":         true,
		"rangeint":         true,
		"reflecttypefor":   true,
		"slicescontains":   true,
		"slicessort":       true,
		"stditerators":     true,
		"stringscutprefix": true,
		"stringsseq":       true,
		"stringsbuilder":   true,
		"testingcontext":   true,
		"waitgroup":        true,
	}
	// Return lookup result
	return modernizeAnalyzers[name]
}

// formatModernizeCode formats an analyzer name as a KTN-MDRNZ code.
//
// Params:
//   - name: analyzer name
//
// Returns:
//   - string: formatted code
func (p *DiagnosticsProcessor) formatModernizeCode(name string) string {
	// Return formatted code
	return "KTN-MDRNZ-" + strings.ToUpper(name)
}
