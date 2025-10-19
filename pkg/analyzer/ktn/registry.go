package ktn

import (
	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"
)

// GetAllRules retourne toutes les règles KTN disponibles.
func GetAllRules() []*analysis.Analyzer {
	var all []*analysis.Analyzer
	all = append(all, ktnconst.Analyzers()...)
	all = append(all, ktnfunc.Analyzers()...)
	return all
}

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
func GetRulesByCategory(category string) []*analysis.Analyzer {
	switch category {
	case "const":
		return ktnconst.Analyzers()
	case "func":
		return ktnfunc.Analyzers()
	default:
		return nil
	}
}
