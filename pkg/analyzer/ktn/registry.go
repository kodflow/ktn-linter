package ktn

import (
	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"
	ktnvar "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var"
)

// GetAllRules retourne toutes les règles KTN disponibles.
//
// Returns:
//   - []*analysis.Analyzer: liste de tous les analyseurs (const + func + var)
func GetAllRules() []*analysis.Analyzer {
	var all []*analysis.Analyzer
	// Ajoute les analyseurs de constantes
	all = append(all, ktnconst.GetAnalyzers()...)
	// Ajoute les analyseurs de fonctions
	all = append(all, ktnfunc.GetAnalyzers()...)
	// Ajoute les analyseurs de variables
	all = append(all, ktnvar.Analyzers()...)
	// Retourne la liste complète
	return all
}

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
//
// Params:
//   - category: nom de la catégorie ("const", "func" ou "var")
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de la catégorie demandée
func GetRulesByCategory(category string) []*analysis.Analyzer {
	// Sélection de la catégorie
	switch category {
 // Traitement
	case "const":
		// Retourne les analyseurs de constantes
		return ktnconst.GetAnalyzers()
 // Traitement
	case "func":
		// Retourne les analyseurs de fonctions
		return ktnfunc.GetAnalyzers()
 // Traitement
	case "var":
		// Retourne les analyseurs de variables
		return ktnvar.Analyzers()
 // Traitement
	default:
		// Catégorie inconnue
		return nil
	}
}
