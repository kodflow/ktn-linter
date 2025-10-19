package ktnfunc

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer006 checks that error is always the last return value
var Analyzer006 = &analysis.Analyzer{
	Name:     "ktnfunc006",
	Doc:      "KTN-FUNC-006: L'erreur doit toujours être en dernière position dans les valeurs de retour",
	Run:      runFunc006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runFunc006(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		var funcType *ast.FuncType

		switch node := n.(type) {
		case *ast.FuncDecl:
			funcType = node.Type
		case *ast.FuncLit:
			funcType = node.Type
		}

		if funcType == nil || funcType.Results == nil {
			return
		}

		results := funcType.Results.List
		if len(results) == 0 {
			return
		}

		// Check each return value to find error types
		var errorPositions []int
		for i, result := range results {
			if isErrorType(pass, result.Type) {
				errorPositions = append(errorPositions, i)
			}
		}

		// If there are errors and any is not in the last position, report
		if len(errorPositions) > 0 {
			lastPos := len(results) - 1
			for _, pos := range errorPositions {
				if pos != lastPos {
					pass.Reportf(
						funcType.Results.Pos(),
						"KTN-FUNC-006: l'erreur doit être en dernière position dans les valeurs de retour (trouvée en position %d sur %d)",
						pos+1,
						len(results),
					)
					return
				}
			}
		}
	})

	return nil, nil
}

// isErrorType checks if a type expression represents the error interface
func isErrorType(pass *analysis.Pass, expr ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok {
		return false
	}

	// Check if it's the error interface
	named, ok := tv.Type.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	return obj != nil && obj.Name() == "error" && obj.Pkg() == nil
}
