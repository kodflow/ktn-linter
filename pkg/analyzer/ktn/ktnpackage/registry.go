package ktnpackage

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs de la catégorie package.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs
func Analyzers() []*analysis.Analyzer {
	// Retourne la liste des analyseurs
	return []*analysis.Analyzer{
		Analyzer001, // Description obligatoire avant package
	}
}
