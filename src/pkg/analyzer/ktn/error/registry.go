package ktn_error

import "golang.org/x/tools/go/analysis"

// AllRules contient toutes les règles KTN-ERROR.
var AllRules = []*analysis.Analyzer{
	Rule001, // KTN-ERROR-001: Wrapping d'erreurs avec contexte
}

// GetRules retourne toutes les règles ERROR.
//
// Returns:
//   - []*analysis.Analyzer: la liste de toutes les règles KTN-ERROR
func GetRules() []*analysis.Analyzer {
	return AllRules
}
