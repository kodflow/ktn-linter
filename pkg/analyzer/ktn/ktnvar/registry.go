// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import "golang.org/x/tools/go/analysis"

// Analyzers retourne tous les analyseurs de la catégorie VAR.
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs VAR
func Analyzers() []*analysis.Analyzer {
	// Retour de la liste complete des analyseurs VAR (36 regles)
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
		Analyzer020, // Nil slice preferred (NEW)
		Analyzer021, // Receiver consistency (NEW)
		Analyzer022, // Pointer to interface
		Analyzer023, // crypto/rand for secrets
		Analyzer024, // any vs interface{}
		Analyzer025, // clear() built-in (Go 1.21+)
		Analyzer026, // min()/max() built-in (Go 1.21+)
		Analyzer027, // range over int (Go 1.22+)
		Analyzer028, // loop var copy obsolete (Go 1.22+)
		Analyzer029, // slices.Grow (Go 1.21+)
		Analyzer030, // slices.Clone (Go 1.21+)
		Analyzer031, // maps.Clone (Go 1.21+)
		Analyzer033, // cmp.Or (Go 1.22+)
		Analyzer034, // WaitGroup.Go (Go 1.25+)
		Analyzer035, // slices.Contains (Go 1.21+)
		Analyzer036, // slices.Index (Go 1.21+)
		Analyzer037, // maps.Keys/Values (Go 1.23+)
	}
}
