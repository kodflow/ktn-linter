package ktnreturn

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer002 detects nil returns for slice and map types.
var Analyzer002 = &analysis.Analyzer{
	Name:     "ktnreturn002",
	Doc:      "KTN-RETURN-002: préférer slice/map vide à nil",
	Run:      runReturn002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runReturn002 analyzes return statements for nil slice/map returns.
// Params:
//   - pass: Analysis pass containing type information
//
// Returns:
//   - any: always nil
//   - error: analysis error if any
func runReturn002(pass *analysis.Pass) (any, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip if function has no return type
		if funcDecl.Type == nil || funcDecl.Type.Results == nil {
			return
		}

		// Check all return types
		for _, result := range funcDecl.Type.Results.List {
			// Verification de la condition
			if isSliceOrMapType(pass, result.Type) {
				// Analyze function body for nil returns
				checkNilReturns(pass, funcDecl)
				break
			}
		}
	})

	return nil, nil
}

// isSliceOrMapType checks if expression is slice or map type.
// Params:
//   - pass: Analysis pass
//   - expr: Expression to check
func isSliceOrMapType(pass *analysis.Pass, expr ast.Expr) bool {
	typeInfo := pass.TypesInfo.TypeOf(expr)
	// Return false if type information is unavailable
	if typeInfo == nil {
		return false
	}

	// Check underlying type
	switch typeInfo.Underlying().(type) {
	// Verification de la condition
	case *types.Slice, *types.Map:
		return true
	}
	return false
}

// checkNilReturns analyzes function body for nil returns.
// Params:
//   - pass: Analysis pass
//   - funcDecl: Function declaration to analyze
func checkNilReturns(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Skip if function has no body
	if funcDecl.Body == nil {
		return
	}

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		retStmt, ok := n.(*ast.ReturnStmt)
		// Continue traversal if not return statement
		if !ok {
			return true
		}

		// Check each return value
		for _, result := range retStmt.Results {
			// Verification de la condition
			if isNilIdent(result) {
				pass.Reportf(
					retStmt.Pos(),
					"KTN-RETURN-002: préférer slice/map vide à nil",
				)
			}
		}

		return true
	})
}

// isNilIdent checks if expression is nil identifier.
// Params:
//   - expr: Expression to check
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	// Return false if not identifier
	if !ok {
		return false
	}
	// Verification de la condition
	return ident.Name == "nil"
}
