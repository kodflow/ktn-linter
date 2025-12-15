// Registry of analyzers for the ktn package.
package ktn

import (
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnapi"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest"
	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar"
	"github.com/kodflow/ktn-linter/pkg/analyzer/modernize"
)

const (
	// ktnPrefixLen is the length of "KTN-" prefix.
	ktnPrefixLen int = 4
	// codePartsCount is the expected number of parts in rule code (CATEGORY-NNN).
	codePartsCount int = 2
)

// GetAllRules retourne toutes les règles KTN disponibles.
//
// Returns:
//   - []*analysis.Analyzer: liste de tous les analyseurs (api + const + func + struct + var + test + package + modernize)
func GetAllRules() []*analysis.Analyzer {
	var all []*analysis.Analyzer
	// Ajoute les analyseurs d'API/dépendances
	all = append(all, ktnapi.Analyzers()...)
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
		"api":       ktnapi.Analyzers,
		"const":     ktnconst.GetAnalyzers,
		"func":      ktnfunc.GetAnalyzers,
		"struct":    ktnstruct.GetAnalyzers,
		"var":       ktnvar.Analyzers,
		"test":      ktntest.Analyzers,
		"return":    ktnreturn.Analyzers,
		"interface": ktninterface.Analyzers,
		"comment":   ktncomment.Analyzers,
		"modernize": modernize.Analyzers,
	}
}

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
//
// Params:
//   - category: nom de la catégorie ("const", "func", "struct", "var", "test", "return", "interface", "comment" ou "modernize")
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

// GetRuleByCode retourne un analyseur par son code (ex: KTN-FUNC-001).
//
// Params:
//   - code: code de la règle (ex: "KTN-FUNC-001", "KTN-VAR-002")
//
// Returns:
//   - *analysis.Analyzer: l'analyseur correspondant ou nil si non trouvé
func GetRuleByCode(code string) *analysis.Analyzer {
	// Convertir le code en nom d'analyseur
	// KTN-FUNC-001 -> ktnfunc001
	// KTN-VAR-002 -> ktnvar002
	analyzerName := codeToAnalyzerName(code)
	// Vérifier si le nom est vide
	if analyzerName == "" {
		// Code invalide
		return nil
	}

	// Chercher l'analyseur dans toutes les règles
	for _, a := range GetAllRules() {
		// Comparer les noms
		if a.Name == analyzerName {
			// Retourne l'analyseur trouvé
			return a
		}
	}

	// Analyseur non trouvé
	return nil
}

// codeToAnalyzerName convertit un code de règle en nom d'analyseur.
//
// Params:
//   - code: code de la règle (ex: "KTN-FUNC-001")
//
// Returns:
//   - string: nom de l'analyseur (ex: "ktnfunc001") ou vide si invalide
func codeToAnalyzerName(code string) string {
	// Format attendu: KTN-CATEGORY-NNN
	// Vérifier le préfixe KTN-
	if !strings.HasPrefix(code, "KTN-") {
		// Format invalide
		return ""
	}

	// Supprimer le préfixe KTN-
	rest := code[ktnPrefixLen:]

	// Séparer la catégorie et le numéro
	parts := strings.Split(rest, "-")
	// Vérifier le format
	if len(parts) != codePartsCount {
		// Format invalide
		return ""
	}

	category := strings.ToLower(parts[0])
	number := parts[1]

	// Construire le nom de l'analyseur: ktn<category><number>
	return "ktn" + category + number
}
