package ktnfunc

import "golang.org/x/tools/go/analysis"

// GetAnalyzers retourne tous les analyseurs relatifs aux fonctions.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de fonctions (001-012)
func GetAnalyzers() []*analysis.Analyzer {
	// Retourne la liste complète des analyseurs de fonctions
	return []*analysis.Analyzer{
		Analyzer001, // Max 35 lines of pure code
		Analyzer002, // Max 5 parameters
		Analyzer003, // No magic numbers
		Analyzer004, // No naked returns (except short functions)
		Analyzer005, // Max cyclomatic complexity 10
		Analyzer006, // Error must be last
		Analyzer007, // Documentation stricte
		Analyzer008, // Context must be first parameter
		Analyzer009, // No side effects in getters
		Analyzer010, // Named returns for >3 return values
		Analyzer011, // Comments on branches/returns
		Analyzer012, // No else after return/continue/break
	}
}
