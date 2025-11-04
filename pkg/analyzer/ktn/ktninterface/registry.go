package ktninterface

import "golang.org/x/tools/go/analysis"

// Analyzers returns all interface-related analyzers.
// Params:
//   - N/A
func Analyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		Analyzer001,
	}
}
