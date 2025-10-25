package ktn

import (
	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
)

// GetAllRules retourne toutes les règles KTN disponibles.
//
// Returns:
//   - []*analysis.Analyzer: liste de tous les analyseurs (const + func + struct + var + test)
func GetAllRules() []*analysis.Analyzer {
	var all []*analysis.Analyzer
	// Ajoute les analyseurs de constantes
	all = append(all, ktnconst.GetAnalyzers()...)
	// Ajoute les analyseurs de fonctions
	all = append(all, ktnfunc.GetAnalyzers()...)
	// Ajoute les analyseurs de structures
	all = append(all, ktnstruct.GetAnalyzers()...)
	// Ajoute les analyseurs de variables
	all = append(all, ktnvar.Analyzers()...)
	// Ajoute les analyseurs de tests
	all = append(all, ktntest.Analyzers()...)
	// Retourne la liste complète
	return all
}

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
//
// Params:
//   - category: nom de la catégorie ("const", "func", "struct", "var" ou "test")
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
	case "struct":
		// Retourne les analyseurs de structures
		return ktnstruct.GetAnalyzers()
	// Traitement
	case "var":
		// Retourne les analyseurs de variables
		return ktnvar.Analyzers()
	// Traitement
	case "test":
		// Retourne les analyseurs de tests
		return ktntest.Analyzers()
	// Traitement
	default:
		// Catégorie inconnue
		return nil
	}
}
