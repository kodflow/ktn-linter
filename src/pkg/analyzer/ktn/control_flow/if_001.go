package ktn_control_flow

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleIf001 = &analysis.Analyzer{
	Name: "KTN_IF_004",
	Doc:  "Détecte les expressions booléennes simplifiables (Staticcheck S1008)",
	Run:  runRuleIf001,
}

func runRuleIf001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}

			// Vérifier que le if n'a pas d'initialisation
			if ifStmt.Init != nil || ifStmt.Body == nil || len(ifStmt.Body.List) != 1 {
				return true
			}

			// Vérifier que le body est un return avec 1 valeur booléenne
			bodyStmt, ok := ifStmt.Body.List[0].(*ast.ReturnStmt)
			if !ok || len(bodyStmt.Results) != 1 {
				return true
			}

			bodyValue, bodyBool := getBooleanLiteral(bodyStmt.Results[0])
			if !bodyBool {
				return true
			}

			// Cas 1: if/else
			if ifStmt.Else != nil {
				elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
				if !ok || len(elseBlock.List) != 1 {
					return true
				}

				elseStmt, ok := elseBlock.List[0].(*ast.ReturnStmt)
				if !ok || len(elseStmt.Results) != 1 {
					return true
				}

				elseValue, elseBool := getBooleanLiteral(elseStmt.Results[0])
				if !elseBool || bodyValue == elseValue {
					return true
				}

				suggestion := "return <condition>"
				if !bodyValue {
					suggestion = "return !<condition>"
				}

				pass.Reportf(ifStmt.Pos(),
					"[KTN-IF-001] Expression booléenne simplifiable (Staticcheck S1008).\n"+
						"Un if qui retourne des littéraux booléens peut être simplifié.\n"+
						"Suggestion: %s\n"+
						"Exemple:\n"+
						"  // ❌ MAUVAIS\n"+
						"  if isValid { return true } else { return false }\n"+
						"\n"+
						"  // ✅ CORRECT\n"+
						"  return isValid",
					suggestion)
			}
			return true
		})
	}
	return nil, nil
}

func getBooleanLiteral(expr ast.Expr) (bool, bool) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false, false
	}
	if ident.Name == "true" {
		return true, true
	}
	if ident.Name == "false" {
		return false, true
	}
	return false, false
}
