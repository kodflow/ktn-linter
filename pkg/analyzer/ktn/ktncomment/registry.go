package ktncomment

import "golang.org/x/tools/go/analysis"

// Analyzers returns all comment-related analyzers.
// Params:
//   - N/A
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer001,
		Analyzer002,
	}
}
