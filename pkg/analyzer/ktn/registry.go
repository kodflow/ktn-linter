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

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
//
// Params:
//   - category: nom de la catégorie ("const", "func", "struct", "var", "test", "return", "interface", "comment", "package" ou "modernize")
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
	case "return":
		// Retourne les analyseurs de retours
		return ktnreturn.Analyzers()
	// Traitement
	case "interface":
		// Retourne les analyseurs d'interfaces
		return ktninterface.Analyzers()
	// Traitement
	case "comment":
		// Retourne les analyseurs de commentaires
		return ktncomment.Analyzers()
	// Traitement
	case "package":
		// Retourne les analyseurs de packages
		return ktnpackage.Analyzers()
	// Traitement
	case "modernize":
		// Retourne les analyseurs modernize (golang.org/x/tools)
		return modernize.Analyzers()
	// Traitement
	default:
		// Catégorie inconnue
		return nil
	}
}
