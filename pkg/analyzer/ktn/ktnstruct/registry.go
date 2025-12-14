// Registry of analyzers for the ktnstruct package.
package ktnstruct

import "golang.org/x/tools/go/analysis"

// GetAnalyzers retourne tous les analyseurs relatifs aux structures.
// Note: Analyzer001 (KTN-STRUCT-001) est déprécié et remplacé par KTN-API-001.
// KTN-STRUCT-001 demandait des "mirror interfaces" (100% des méthodes), ce qui est un anti-pattern.
// KTN-API-001 impose le bon pattern: interfaces minimales côté consumer (ISP).
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de structures
func GetAnalyzers() []*analysis.Analyzer {
	// Retourne la liste complète des analyseurs de structures
	// Note: Analyzer001 (KTN-STRUCT-001) est déprécié - voir KTN-API-001
	return []*analysis.Analyzer{
		// Analyzer001 déprécié - remplacé par KTN-API-001
		Analyzer002,
		Analyzer003,
		Analyzer004,
		Analyzer005,
		Analyzer006,
		Analyzer007,
	}
}
