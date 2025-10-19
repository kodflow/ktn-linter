package ktnfunc

import "golang.org/x/tools/go/analysis"

// Analyzers returns all func-related analyzers
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer001, // Max 35 lines of pure code
		Analyzer002, // Max 5 parameters
		Analyzer003, // Function names must start with a verb
		Analyzer004, // No naked returns (except short functions)
		Analyzer005, // Max cyclomatic complexity 10
		Analyzer006, // Error must be last
		Analyzer007, // Documentation stricte
		Analyzer008, // Context must be first parameter
		Analyzer009, // No side effects in getters
		Analyzer010, // Named returns for >3 return values
		Analyzer011, // Comments on branches/returns
	}
}
