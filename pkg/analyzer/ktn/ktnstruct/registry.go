// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import "golang.org/x/tools/go/analysis"

// GetAnalyzers retourne tous les analyseurs relatifs aux structures.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de structures
func GetAnalyzers() []*analysis.Analyzer {
	// Retourne la liste complète des analyseurs de structures
	return []*analysis.Analyzer{
		Analyzer001,
		Analyzer002,
		Analyzer003,
		Analyzer004,
		Analyzer005,
		Analyzer006,
		Analyzer007,
		Analyzer008,
		Analyzer009,
	}
}
