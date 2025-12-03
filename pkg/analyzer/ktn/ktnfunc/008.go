// Analyzer 008 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer008 checks that context.Context is always the first parameter
var Analyzer008 = &analysis.Analyzer{
	Name:     "ktnfunc008",
	Doc:      "KTN-FUNC-008: context.Context doit toujours être le premier paramètre (après le receiver pour les méthodes)",
	Run:      runFunc008,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc008 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc008(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip test functions
		if shared.IsTestFunction(funcDecl) {
			// Retour de la fonction
			return
		}

		funcName := funcDecl.Name.Name

		// Vérification de la condition
		if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
			// Retour de la fonction
			return
		}

		// Check each parameter
		contextParamIndex := -1
		// Itération sur les éléments
		for i, field := range funcDecl.Type.Params.List {
			// Vérification de la condition
			if isContextType(field.Type) {
				contextParamIndex = i
				break
			}
		}

		// If there's a context parameter and it's not first, report error
		if contextParamIndex > 0 {
			pass.Reportf(
				funcDecl.Type.Params.List[contextParamIndex].Pos(),
				"KTN-FUNC-008: context.Context doit être le premier paramètre de la fonction '%s'",
				funcName,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// isContextType checks if a type is context.Context
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si type context
func isContextType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	// Vérification de la condition
	if !ok {
		// Retour de la fonction
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	// Retour de la fonction
	return ok && ident.Name == "context" && sel.Sel.Name == "Context"
}
