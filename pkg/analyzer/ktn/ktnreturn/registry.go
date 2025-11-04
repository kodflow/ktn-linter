package ktnreturn

import "golang.org/x/tools/go/analysis"

// Analyzers returns all return-related analyzers.
// Params:
//   - N/A
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer002,
	}
}
