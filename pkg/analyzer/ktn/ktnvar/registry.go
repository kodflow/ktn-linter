// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs de la catégorie VAR.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs VAR
func Analyzers() []*analysis.Analyzer {
	// Retour de la liste complète des analyseurs VAR (19 règles)
	return []*analysis.Analyzer{
		Analyzer001, // Types explicites (ex-VAR-002)
		Analyzer002, // Ordre déclaration (ex-VAR-014)
		Analyzer003, // CamelCase (fusion VAR-001+018)
		Analyzer004, // Longueur min (NEW)
		Analyzer005, // Longueur max (NEW)
		Analyzer006, // Shadowing (ex-VAR-011)
		Analyzer007, // := vs var (ex-VAR-003)
		Analyzer008, // Slices préalloc (ex-VAR-004)
		Analyzer009, // make+append (ex-VAR-005)
		Analyzer010, // Buffer.Grow (ex-VAR-006)
		Analyzer011, // strings.Builder (ex-VAR-007)
		Analyzer012, // Alloc loops (ex-VAR-008)
		Analyzer013, // Struct size (ex-VAR-009)
		Analyzer014, // sync.Pool (ex-VAR-010)
		Analyzer015, // string() (ex-VAR-012)
		Analyzer016, // Groupement (ex-VAR-013)
		Analyzer017, // Map prealloc (ex-VAR-015)
		Analyzer018, // Array vs slice (ex-VAR-016)
		Analyzer019, // Mutex copies (ex-VAR-017)
	}
}
