package ktnfunc

import (
	"go/ast"

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

func runFunc008(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		funcName := funcDecl.Name.Name

		// Skip test functions
		if isTestFunction(funcName) {
			return
		}

		if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
			return
		}

		// Check each parameter
		contextParamIndex := -1
		for i, field := range funcDecl.Type.Params.List {
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

	return nil, nil
}

// isContextType checks if a type is context.Context
func isContextType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "context" && sel.Sel.Name == "Context"
}
