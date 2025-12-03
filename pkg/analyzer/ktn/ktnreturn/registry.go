// Registry of analyzers for the ktnreturn package.
package ktnreturn

import "golang.org/x/tools/go/analysis"

// Analyzers returns all return-related analyzers.
//
// Returns:
//   - []*analysis.Analyzer: slice of return analyzers
func Analyzers() []*analysis.Analyzer {
	// Verification de la condition
	return []*analysis.Analyzer{
		Analyzer002,
	}
}
