// Analyzer 009 for the ktnvar package.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer009 detects map allocations without capacity hints
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar009",
	Doc:      "KTN-VAR-009: Préallouer maps avec capacité si connue",
	Run:      runVar009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar009 exécute l'analyse KTN-VAR-009.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar009(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Vérification que c'est un appel à "make"
		if !utils.IsMakeCall(callExpr) {
			// Continue traversing AST nodes
			return
		}

		// Vérification que le type est une map
		if len(callExpr.Args) == 0 || !utils.IsMapTypeWithPass(pass, callExpr.Args[0]) {
			// Continue traversing AST nodes
			return
		}

		// Vérification que make a exactement 1 argument (type seulement)
		if len(callExpr.Args) != 1 {
			// make() avec capacité fournie, conforme
			return
		}

		// Signaler l'erreur
		pass.Reportf(
			callExpr.Pos(),
			"KTN-VAR-009: préallouer la map avec une capacité (make(map[K]V, capacity))",
		)
	})

	// Retour de la fonction
	return nil, nil
}
