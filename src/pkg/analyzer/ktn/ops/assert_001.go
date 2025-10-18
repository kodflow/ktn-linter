package ktn_ops

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleAssert001 = &analysis.Analyzer{
	Name: "KTN_ASSERT_001",
	Doc:  "Détecte les assertions de type sans vérification",
	Run:  runRuleAssert001,
}

func runRuleAssert001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			typeAssert, ok := n.(*ast.TypeAssertExpr)
			if !ok || typeAssert.Type == nil {
				return true
			}

			// Vérifier si c'est dans une assignation à 2 valeurs
			if !isInSafeAssignment(file, typeAssert) {
				pass.Reportf(typeAssert.Pos(),
					"[KTN-ASSERT-001] Assertion de type sans vérification du succès.\n"+
						"Une assertion de type non vérifiée cause un panic si le type est incorrect.\n"+
						"Utilisez toujours la forme à deux valeurs (value, ok) pour vérifier.\n"+
						"Exemple:\n"+
						"  // ❌ MAUVAIS - panic si x n'est pas un int\n"+
						"  v := x.(int)\n"+
						"\n"+
						"  // ✅ CORRECT - vérification du type\n"+
						"  v, ok := x.(int)\n"+
						"  if !ok {\n"+
						"      return errors.New(\"wrong type\")\n"+
						"  }")
			}
			return true
		})
	}
	return nil, nil
}

func isInSafeAssignment(file *ast.File, expr ast.Expr) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for _, rhs := range assignStmt.Rhs {
			if rhs == expr {
				if len(assignStmt.Lhs) == 2 {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}
