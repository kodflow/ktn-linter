// Package ktntest implements KTN linter rules.
package ktntest

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs de la catégorie TEST.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs TEST (12 règles)
func Analyzers() []*analysis.Analyzer {
	// Retourne tous les analyseurs de test ordonnés par criticité
	return []*analysis.Analyzer{
		Analyzer001, // ERROR: fichier test doit finir par _internal/_external_test.go
		// Analyzer002 désactivée: remplacée par 008+009+010+011
		Analyzer003,
		Analyzer004,
		Analyzer005,
		// Analyzer006 désactivée: doublon de 003 (pattern 1:1)
		Analyzer007,
		Analyzer008,
		Analyzer009,
		Analyzer010,
		Analyzer011,
		Analyzer012,
		Analyzer013,
	}
}
