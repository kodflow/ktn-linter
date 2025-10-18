package ktn_ops

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// RulePointer001 analyzer for pointer operations.
var RulePointer001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_POINTER_001",
	Doc:  "Détecte le déréférencement potentiel d'un pointeur nil",
	Run:  runRulePointer001,
}

func runRulePointer001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			unary, ok := n.(*ast.UnaryExpr)
			if !ok || unary.Op.String() != "*" {
				// Continue traversing AST nodes.
				return true
			}

			// Vérifier si c'est *new(...) - sûr
			if call, ok := unary.X.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "new" {
					// Continue traversing AST nodes.
					return true
				}
			}

			// Vérifier si c'est un déréférencement d'une variable potentiellement nil
			if ident, ok := unary.X.(*ast.Ident); ok {
				if isRecentlyAssignedNil(file, ident, unary) {
					pass.Reportf(unary.Pos(),
						"[KTN-POINTER-001] Déréférencement potentiel d'un pointeur nil '%s'.\n"+
							"Déréférencer un pointeur nil cause un panic immédiat.\n"+
							"Vérifiez toujours qu'un pointeur n'est pas nil avant de le déréférencer.\n"+
							"Exemple:\n"+
							"  // ❌ MAUVAIS\n"+
							"  var p *int\n"+
							"  x := *p  // 💥 PANIC\n"+
							"\n"+
							"  // ✅ CORRECT\n"+
							"  var p *int\n"+
							"  if p != nil {\n"+
							"      x := *p\n"+
							"  }",
						ident.Name)
				}
			}
			// Continue traversing AST nodes.
			return true
		})
	}
	// Analysis completed successfully.
	return nil, nil
}

func isRecentlyAssignedNil(file *ast.File, ident *ast.Ident, deref ast.Node) bool {
	nilAssigned := false
	ast.Inspect(file, func(n ast.Node) bool {
		if n == deref {
			// Condition not met, return false.
			return false
		}
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			// Continue traversing AST nodes.
			return true
		}
		for i, lhs := range assignStmt.Lhs {
			if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == ident.Name {
				if i < len(assignStmt.Rhs) {
					if rhsIdent, ok := assignStmt.Rhs[i].(*ast.Ident); ok && rhsIdent.Name == "nil" {
						nilAssigned = true
					}
				}
			}
		}
		// Continue traversing AST nodes.
		return true
	})
	// Early return from function.
	return nilAssigned
}
