package ktninterface

import "golang.org/x/tools/go/analysis"

// Analyzers returns all interface-related analyzers.
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer001,
	}
}
