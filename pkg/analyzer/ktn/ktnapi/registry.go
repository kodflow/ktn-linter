// Package ktnapi provides analyzers for API/dependency coupling rules.
// These rules help ensure loose coupling and testability by promoting
// consumer-side interfaces over concrete external dependencies.
package ktnapi

import "golang.org/x/tools/go/analysis"

// Analyzers returns all analyzers in the ktnapi package.
//
// Returns:
//   - []*analysis.Analyzer: all API analyzers
func Analyzers() []*analysis.Analyzer {
	// Retour de la liste des analyseurs API
	return []*analysis.Analyzer{
		Analyzer001,
	}
}
