package ktn_ops

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var RuleComp001 = &analysis.Analyzer{
	Name: "KTN_COMP_001",
	Doc:  "Détecte les comparaisons booléennes redondantes",
	Run:  runRuleComp001,
}

func runRuleComp001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			binaryExpr, ok := n.(*ast.BinaryExpr)
			if !ok {
				return true
			}

			// Vérifier si c'est == ou !=
			if binaryExpr.Op != token.EQL && binaryExpr.Op != token.NEQ {
				return true
			}

			// Vérifier si un des côtés est true ou false
			xBool := isBooleanLiteral(binaryExpr.X)
			yBool := isBooleanLiteral(binaryExpr.Y)

			if !xBool && !yBool {
				return true
			}

			var boolSide ast.Expr
			if xBool {
				boolSide = binaryExpr.X
			} else {
				boolSide = binaryExpr.Y
			}

			boolIdent, ok := boolSide.(*ast.Ident)
			if !ok {
				return true
			}

			boolValue := boolIdent.Name
			var suggestion string
			if binaryExpr.Op == token.EQL {
				if boolValue == "true" {
					suggestion = "utilisez directement l'expression"
				} else {
					suggestion = "utilisez !expression"
				}
			} else {
				if boolValue == "true" {
					suggestion = "utilisez !expression"
				} else {
					suggestion = "utilisez directement l'expression"
				}
			}

			pass.Reportf(binaryExpr.Pos(),
				"[KTN-COMP-001] Comparaison booléenne redondante avec %s.\n"+
					"Comparer un booléen à true/false est inutile et nuit à la lisibilité.\n"+
					"Utilisez directement l'expression ou sa négation.\n"+
					"Suggestion: %s\n"+
					"Exemple:\n"+
					"  // ❌ MAUVAIS\n"+
					"  if isValid == true { }\n"+
					"  if isValid == false { }\n"+
					"\n"+
					"  // ✅ CORRECT\n"+
					"  if isValid { }\n"+
					"  if !isValid { }",
				boolValue, suggestion)
			return true
		})
	}
	return nil, nil
}

func isBooleanLiteral(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "true" || ident.Name == "false"
}
