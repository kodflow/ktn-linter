// Types for the orchestrator package.
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
}
