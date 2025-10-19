package ktnfunc

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer005 checks that functions don't exceed cyclomatic complexity of 10
var Analyzer005 = &analysis.Analyzer{
	Name:     "ktnfunc005",
	Doc:      "KTN-FUNC-005: La complexité cyclomatique ne doit pas dépasser 10",
	Run:      runFunc005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

const maxCyclomaticComplexity = 10

func runFunc005(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip if no body
		if funcDecl.Body == nil {
			return
		}

		// Skip test functions
		funcName := funcDecl.Name.Name
		if isTestFunction(funcName) {
			return
		}

		// Calculate cyclomatic complexity
		complexity := calculateComplexity(funcDecl.Body)

		if complexity > maxCyclomaticComplexity {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-005: la fonction '%s' a une complexité cyclomatique de %d (max: %d)",
				funcName,
				complexity,
				maxCyclomaticComplexity,
			)
		}
	})

	return nil, nil
}

// calculateComplexity calculates the cyclomatic complexity of a function
func calculateComplexity(body *ast.BlockStmt) int {
	// Start with complexity of 1 (the function itself)
	complexity := 1

	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			// +1 for if
			complexity++
		case *ast.ForStmt, *ast.RangeStmt:
			// +1 for each loop
			complexity++
		case *ast.CaseClause:
			// +1 for each case (except default)
			if node.List != nil {
				complexity++
			}
		case *ast.CommClause:
			// +1 for each comm case in select
			if node.Comm != nil {
				complexity++
			}
		case *ast.BinaryExpr:
			// +1 for && and ||
			if node.Op.String() == "&&" || node.Op.String() == "||" {
				complexity++
			}
		}
		return true
	})

	return complexity
}
