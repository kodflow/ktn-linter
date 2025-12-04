// Registry of analyzers for the ktnvar package.
package ktnvar

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs de la catégorie VAR.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs VAR
func Analyzers() []*analysis.Analyzer {
	// Retour de la liste complète des analyseurs VAR (17 règles)
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
		Analyzer010,
		Analyzer011,
		Analyzer012,
		Analyzer013,
		Analyzer014,
		Analyzer015,
		Analyzer016,
		Analyzer017,
	}
}
