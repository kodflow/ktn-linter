package ktnreturn

import "golang.org/x/tools/go/analysis"

// Analyzers returns all return-related analyzers.
// Params: TODO
// Returns: TODO
func Analyzers() []*analysis.Analyzer {
 // Verification de la condition
	return []*analysis.Analyzer{
		Analyzer002,
	}
}
