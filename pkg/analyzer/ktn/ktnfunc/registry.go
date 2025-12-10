// Registry of analyzers for the ktnfunc package.
package ktnfunc

import "golang.org/x/tools/go/analysis"

// GetAnalyzers retourne tous les analyseurs relatifs aux fonctions.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de fonctions (001-012)
func GetAnalyzers() []*analysis.Analyzer {
	// Retourne la liste complète des analyseurs de fonctions
	return []*analysis.Analyzer{
		Analyzer001, // Error must be last return
		Analyzer002, // Context must be first parameter
		Analyzer003, // No else after return/continue/break
		Analyzer004, // Private functions must be used
		Analyzer005, // Max 35 lines of pure code
		Analyzer006, // Max 5 parameters
		Analyzer007, // No side effects in getters
		Analyzer008, // Unused parameters must be prefixed with _
		Analyzer009, // No magic numbers
		Analyzer010, // No naked returns (except short functions)
		Analyzer011, // Max cyclomatic complexity 10
		Analyzer012, // Named returns for >3 return values
	}
}
