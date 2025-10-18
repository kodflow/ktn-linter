package ktnconst

import "golang.org/x/tools/go/analysis"

// Analyzers returns all const-related analyzers
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer001,
		Analyzer002,
		Analyzer003,
		Analyzer004,
		Analyzer005,
	}
}
