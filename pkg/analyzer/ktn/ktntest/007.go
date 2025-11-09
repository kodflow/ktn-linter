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
			// Continue traversal
			return
		}

		// Vérifier si la méthode s'appelle "Skip", "Skipf" ou "SkipNow"
		methodName := selExpr.Sel.Name
		// Vérification de la condition
		if methodName != "Skip" && methodName != "Skipf" && methodName != "SkipNow" {
			// Pas un appel à Skip, continuer
			return
		}

		// Vérifier si le receiver est 't' ou un nom de variable de test
		ident, ok := selExpr.X.(*ast.Ident)
		// Pas un identifiant simple, continuer
		if !ok {
			// Continue traversal
			return
		}

		// Vérifier si c'est dans un fichier de test
		filename := pass.Fset.Position(n.Pos()).Filename
		// Vérification de la condition
		if !shared.IsTestFile(filename) {
			// Pas un fichier de test, continuer
			return
		}

		// Reporter l'erreur
		pass.Reportf(
			callExpr.Pos(),
			"KTN-TEST-007: utilisation de %s.%s() interdite. Les tests doivent être corrigés ou supprimés, pas skippés",
			ident.Name,
			methodName,
		)
	})

	// Retour de la fonction
	return nil, nil
}
