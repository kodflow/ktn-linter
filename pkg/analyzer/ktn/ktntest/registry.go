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
	return []*analysis.Analyzer{
		// Analyzer001, // Désactivée : incompatible avec white-box testing
		Analyzer002,
		Analyzer003,
		Analyzer004,
		Analyzer005,
		Analyzer006,
		Analyzer007,
	}
}
