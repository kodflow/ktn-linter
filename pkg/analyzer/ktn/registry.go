// Registry of analyzers for the ktn package.
package ktn

import (
	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnpackage"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/modernize"
)

// GetAllRules retourne toutes les règles KTN disponibles.
//
// Returns:
//   - []*analysis.Analyzer: liste de tous les analyseurs (const + func + struct + var + test + package + modernize)
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
	// Ajoute les analyseurs de retours
	all = append(all, ktnreturn.Analyzers()...)
	// Ajoute les analyseurs d'interfaces
	all = append(all, ktninterface.Analyzers()...)
	// Ajoute les analyseurs de commentaires
	all = append(all, ktncomment.Analyzers()...)
	// Ajoute les analyseurs de packages
	all = append(all, ktnpackage.Analyzers()...)
	// Ajoute les analyseurs modernize (golang.org/x/tools)
	all = append(all, modernize.Analyzers()...)
	// Retourne la liste complète
	return all
}

// categoryAnalyzers retourne la map des catégories vers leurs analyseurs.
//
// Returns:
//   - map[string]func() []*analysis.Analyzer: map des fonctions d'analyseurs par catégorie
func categoryAnalyzers() map[string]func() []*analysis.Analyzer {
	// Retour de la map des catégories
	return map[string]func() []*analysis.Analyzer{
		"const":     ktnconst.GetAnalyzers,
		"func":      ktnfunc.GetAnalyzers,
		"struct":    ktnstruct.GetAnalyzers,
		"var":       ktnvar.Analyzers,
		"test":      ktntest.Analyzers,
		"return":    ktnreturn.Analyzers,
		"interface": ktninterface.Analyzers,
		"comment":   ktncomment.Analyzers,
		"package":   ktnpackage.Analyzers,
		"modernize": modernize.Analyzers,
	}
}

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
//
// Params:
//   - category: nom de la catégorie ("const", "func", "struct", "var", "test", "return", "interface", "comment", "package" ou "modernize")
//
// Returns:
//   - []*analysis.Analyzer: liste des analyseurs de la catégorie demandée
func GetRulesByCategory(category string) []*analysis.Analyzer {
	// Récupérer la map des catégories
	categories := categoryAnalyzers()

	// Rechercher la fonction d'analyseurs pour cette catégorie
	analyzerFunc, exists := categories[category]
	// Vérification de la condition
	if !exists {
		// Catégorie inconnue - retourner slice vide
		return []*analysis.Analyzer{}
	}

	// Retour des analyseurs de la catégorie
	return analyzerFunc()
}
