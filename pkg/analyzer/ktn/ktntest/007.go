// Analyzer 007 for the ktntest package.
package ktntest

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
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
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Vérifier si c'est un appel de méthode
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		// Pas un appel de méthode, continuer
		if !ok {
			return
		}

		// Vérifier si la méthode s'appelle "Skip", "Skipf" ou "SkipNow"
		methodName := selExpr.Sel.Name
		// Vérification de la condition
		if !isSkipMethod(methodName) {
			return
		}

		// Vérifier si le receiver est un identifiant simple
		ident, ok := selExpr.X.(*ast.Ident)
		// Pas un identifiant simple, continuer
		if !ok {
			return
		}

		// Vérifier si c'est dans un fichier de test
		filename := pass.Fset.Position(n.Pos()).Filename
		// Vérification de la condition
		if !shared.IsTestFile(filename) {
			return
		}

		// Reporter l'erreur - AUCUNE exception
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
