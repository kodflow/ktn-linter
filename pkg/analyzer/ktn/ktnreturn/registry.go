package ktnreturn

import "golang.org/x/tools/go/analysis"

// Analyzers returns all return-related analyzers.
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer002,
	}
}
