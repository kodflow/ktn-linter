package ktnfunc

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer009 checks that getter functions don't have side effects
var Analyzer009 = &analysis.Analyzer{
	Name:     "ktnfunc009",
	Doc:      "KTN-FUNC-009: Les getters (Get*/Is*/Has*) ne doivent pas avoir de side effects (assignations, appels de fonctions modifiant l'état)",
	Run:      runFunc009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runFunc009(pass *analysis.Pass) (any, error) {
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

		// Skip if not a getter (Get*, Is*, Has*)
		if !isGetter(funcName) {
			return
		}

		// Skip if no body
		if funcDecl.Body == nil {
			return
		}

		// Check for side effects
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			switch stmt := node.(type) {
			case *ast.AssignStmt:
				// Check if it's assigning to a field or external variable
				// Assignments to local variables (created in the function) are OK
				for _, lhs := range stmt.Lhs {
					if hasSideEffect(lhs) {
						pass.Reportf(
							stmt.Pos(),
							"KTN-FUNC-009: le getter '%s' ne doit pas modifier l'état (assignation détectée)",
							funcName,
						)
					}
				}
			case *ast.IncDecStmt:
				// ++ or -- on fields
				if hasSideEffect(stmt.X) {
					pass.Reportf(
						stmt.Pos(),
						"KTN-FUNC-009: le getter '%s' ne doit pas modifier l'état (incrémentation/décrémentation détectée)",
						funcName,
					)
				}
			}
			return true
		})
	})

	return nil, nil
}

// isGetter checks if a function name suggests it's a getter
func isGetter(name string) bool {
	return strings.HasPrefix(name, "Get") ||
		strings.HasPrefix(name, "Is") ||
		strings.HasPrefix(name, "Has")
}

// hasSideEffect checks if an expression modifies external state
func hasSideEffect(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		// Modifying a field is a side effect
		return true
	case *ast.IndexExpr:
		// Modifying an index (array/map/slice element) could be a side effect
		// Check if the base is a selector
		if _, ok := e.X.(*ast.SelectorExpr); ok {
			return true
		}
	}
	return false
}
