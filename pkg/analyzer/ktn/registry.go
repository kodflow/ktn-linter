package ktn

import (
	"golang.org/x/tools/go/analysis"

	ktnconst "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const"
	ktnfunc "github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func"
)

// GetAllRules retourne toutes les règles KTN disponibles.
//
// Returns:
//   - []*analysis.Analyzer: liste de tous les analyseurs (const + func)
func GetAllRules() []*analysis.Analyzer {
	var all []*analysis.Analyzer
	// Ajoute les analyseurs de constantes
	all = append(all, ktnconst.GetAnalyzers()...)
	// Ajoute les analyseurs de fonctions
	all = append(all, ktnfunc.GetAnalyzers()...)
	// Retourne la liste complète
	return all
}

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
//
// Params:
//   - category: nom de la catégorie ("const" ou "func")
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
	default:
		// Catégorie inconnue
		return nil
	}
}
