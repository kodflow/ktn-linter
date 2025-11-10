package ktninterface

import "golang.org/x/tools/go/analysis"

// Analyzers returns all interface-related analyzers.
// Params: TODO
// Returns: TODO
func Analyzers() []*analysis.Analyzer {
 // Verification de la condition
	return []*analysis.Analyzer{
		Analyzer001,
	}
}
