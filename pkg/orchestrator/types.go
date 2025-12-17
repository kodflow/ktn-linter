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
