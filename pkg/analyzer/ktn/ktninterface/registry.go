// Package ktninterface provides analyzers for interface-related lint rules.
package ktninterface

import "golang.org/x/tools/go/analysis"

// Analyzers returns all interface-related analyzers.
//
// Returns:
//   - []*analysis.Analyzer: slice of interface analyzers
func Analyzers() []*analysis.Analyzer {
	// Verification de la condition
	return []*analysis.Analyzer{
		Analyzer001,
		Analyzer003,
		Analyzer004,
	}
}
