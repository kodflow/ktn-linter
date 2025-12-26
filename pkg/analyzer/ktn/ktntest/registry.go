// Package ktntest implements KTN linter rules.
package ktntest

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs de la catégorie TEST.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs TEST (11 règles)
func Analyzers() []*analysis.Analyzer {
	// Retourne tous les analyseurs de test ordonnés par criticité
	return []*analysis.Analyzer{
		Analyzer001, // ERROR: fichier test doit finir par _internal/_external_test.go
		Analyzer002, // Orphan test detection
		Analyzer003, // Function test coverage
		Analyzer004, // Table-driven pattern required
		Analyzer005, // t.Run() subtest usage
		Analyzer006, // Public/private test file separation
		Analyzer007, // Test file naming convention
		Analyzer008, // Exported function test coverage
		Analyzer009, // Internal/external test file pairing
		Analyzer010, // Mock file detection
		Analyzer011, // Mock function detection
	}
}
