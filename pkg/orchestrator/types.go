// Package orchestrator coordinates the linting pipeline.
package orchestrator

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// DiagnosticResult associates a diagnostic with its FileSet and analyzer.
// Used to track the source of each diagnostic for proper formatting.
type DiagnosticResult struct {
	Diag         analysis.Diagnostic
	Fset         *token.FileSet
	AnalyzerName string
	cachedPos    *token.Position // Cached position to avoid repeated lookups
}

// NewDiagnosticResult creates a new DiagnosticResult.
//
// Params:
//   - diag: the diagnostic
//   - fset: the fileset for position lookup
//   - analyzerName: name of the analyzer that reported the diagnostic
//
// Returns:
//   - DiagnosticResult: new diagnostic result
func NewDiagnosticResult(diag analysis.Diagnostic, fset *token.FileSet, analyzerName string) DiagnosticResult {
	// Return new diagnostic result
	return DiagnosticResult{
		Diag:         diag,
		Fset:         fset,
		AnalyzerName: analyzerName,
	}
}

// Position returns the token.Position for this diagnostic.
// Caches the result to avoid repeated Fset.Position() calls.
//
// Returns:
//   - token.Position: position of the diagnostic
func (d *DiagnosticResult) Position() token.Position {
	// Check if position was already computed and cached
	if d.cachedPos != nil {
		// Return cached position to avoid redundant FileSet lookup
		return *d.cachedPos
	}
	// Compute position from FileSet and cache for future calls
	pos := d.Fset.Position(d.Diag.Pos)
	d.cachedPos = &pos
	// Return the newly computed position
	return pos
}
