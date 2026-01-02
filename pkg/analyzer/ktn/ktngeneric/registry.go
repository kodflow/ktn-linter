// Package ktngeneric provides analyzers for generic function lint rules.
package ktngeneric

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs relatifs aux fonctions generiques.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de fonctions generiques (001, 002, 003, 005, 006)
func Analyzers() []*analysis.Analyzer {
	// Retourne la liste complete des analyseurs de fonctions generiques
	return []*analysis.Analyzer{
		Analyzer001,
		Analyzer002,
		Analyzer003,
		Analyzer005,
		Analyzer006,
	}
}
