// Package ktntest implements KTN linter rules.
package ktntest

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs de la catégorie TEST.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs TEST
func Analyzers() []*analysis.Analyzer {
	// Retourne tous les analyseurs de test
	// NOTE: KTN-TEST-001 désactivée car elle force black-box testing (xxx_test)
	// ce qui empêche de tester les fonctions privées (white-box testing)
	// KTN-TEST-008, 009, 010, 011 la remplacent avec une convention hybride internal/external stricte
	return []*analysis.Analyzer{
		// Analyzer001, // Désactivée : remplacée par Analyzer008+009+010+011
		Analyzer002,
		Analyzer003,
		Analyzer004,
		Analyzer005,
		Analyzer006,
		Analyzer007,
		Analyzer008, // Règle 1:2 (chaque .go doit avoir _internal ET _external)
		Analyzer009, // Tests fonctions publiques dans _external uniquement
		Analyzer010, // Tests fonctions privées dans _internal uniquement
		Analyzer011, // Convention package: _internal → xxx, _external → xxx_test
	}
}
