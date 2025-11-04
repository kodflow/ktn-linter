package ktncomment

import "golang.org/x/tools/go/analysis"

// Analyzers returns all comment-related analyzers.
//
// Returns:
//   - []*analysis.Analyzer: list of comment analyzers
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer002,
	}
}
