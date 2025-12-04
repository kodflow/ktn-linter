// Analyzer 016 for the ktnvar package.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer016 detects map allocations without capacity hints
var Analyzer016 = &analysis.Analyzer{
	Name:     "ktnvar016",
	Doc:      "KTN-VAR-016: Préallouer maps avec capacité si connue",
	Run:      runVar016,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar016 exécute l'analyse KTN-VAR-016.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar016(pass *analysis.Pass) (any, error) {
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
			"KTN-VAR-016: préallouer la map avec une capacité (make(map[K]V, capacity))",
		)
	})

	// Retour de la fonction
	return nil, nil
}
