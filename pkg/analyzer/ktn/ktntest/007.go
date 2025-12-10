// Analyzer 007 for the ktntest package.
package ktntest

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	ruleCodeTest007 string = "KTN-TEST-007"
)

// Analyzer007 checks that tests don't use t.Skip()
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest007",
	Doc:      "KTN-TEST-007: Interdiction d'utiliser t.Skip() dans les tests",
	Run:      runTest007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest007 exécute l'analyse KTN-TEST-007.
// AUCUN t.Skip() n'est autorisé - les tests doivent être corrigés ou supprimés.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest007(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest007) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	// Parcourir les appels
	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(n.Pos()).Filename
		if cfg.IsFileExcluded(ruleCodeTest007, filename) {
			// Fichier exclu
			return
		}

		// Vérifier si appel de méthode
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		// Vérification sélecteur
		if !ok {
			// Retour si pas sélecteur
			return
		}

		// Vérifier nom méthode
		methodName := selExpr.Sel.Name
		// Vérification méthode Skip
		if !isSkipMethod(methodName) {
			// Retour si pas Skip
			return
		}

		// Vérifier receiver
		ident, ok := selExpr.X.(*ast.Ident)
		// Vérification identifiant
		if !ok {
			// Retour si pas identifiant
			return
		}

		// Vérification fichier test
		if !shared.IsTestFile(filename) {
			// Retour si pas test
			return
		}

		// Signaler l'erreur
		pass.Reportf(
			callExpr.Pos(),
			"KTN-TEST-007: utilisation de %s.%s() interdite. "+
				"Les tests doivent être corrigés ou supprimés, jamais skippés. "+
				"Un test skippé est une dette technique cachée",
			ident.Name,
			methodName,
		)
	})

	// Retour de la fonction
	return nil, nil
}

// isSkipMethod vérifie si le nom de méthode est une méthode Skip.
//
// Params:
//   - methodName: nom de la méthode à vérifier
//
// Returns:
//   - bool: true si c'est Skip, Skipf ou SkipNow
func isSkipMethod(methodName string) bool {
	// Vérifier les noms de méthodes Skip
	return methodName == "Skip" || methodName == "Skipf" || methodName == "SkipNow"
}
