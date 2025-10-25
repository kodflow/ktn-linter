package ktnconst

import "golang.org/x/tools/go/analysis"

// GetAnalyzers retourne tous les analyseurs relatifs aux constantes.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de constantes (001-004)
func GetAnalyzers() []*analysis.Analyzer {
	// Retourne la liste complète des analyseurs de constantes
	return []*analysis.Analyzer{
		Analyzer001,
		Analyzer002,
		Analyzer003,
		Analyzer004,
	}
}
