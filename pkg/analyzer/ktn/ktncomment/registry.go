// Package ktncomment provides analyzers for comment formatting rules.
package ktncomment

import "golang.org/x/tools/go/analysis"

// Analyzers returns all comment-related analyzers.
//
// Returns:
//   - []*analysis.Analyzer: list of comment analyzers
func Analyzers() []*analysis.Analyzer {
	// Retourne tous les analyseurs de commentaires
	return []*analysis.Analyzer{
		Analyzer001,
		Analyzer002,
		Analyzer003,
		Analyzer004,
		Analyzer005,
		Analyzer006,
		Analyzer007,
	}
}
