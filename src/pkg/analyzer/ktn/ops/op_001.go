package ktn_ops

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var RuleOp001 = &analysis.Analyzer{
	Name: "KTN_OP_001",
	Doc:  "Détecte la division ou modulo par zéro",
	Run:  runRuleOp001,
}

func runRuleOp001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			binary, ok := n.(*ast.BinaryExpr)
			if !ok {
				return true
			}

			// Vérifier si c'est division ou modulo
			if binary.Op != token.QUO && binary.Op != token.REM {
				return true
			}

			// Vérifier si right est zéro
			if isZeroLiteral(binary.Y) {
				pass.Reportf(binary.Pos(),
					"[KTN-OP-001] Division ou modulo par zéro.\n"+
						"Division/modulo par zéro cause un panic immédiat en Go.\n"+
						"Vérifier que le diviseur n'est pas zéro avant l'opération.\n"+
						"Exemple:\n"+
						"  // ❌ MAUVAIS - panic\n"+
						"  result := x / 0\n"+
						"\n"+
						"  // ✅ CORRECT - vérifier avant\n"+
						"  if divisor == 0 {\n"+
						"      return errors.New(\"division by zero\")\n"+
						"  }\n"+
						"  result := x / divisor")
			}
			return true
		})
	}
	return nil, nil
}

func isZeroLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	return lit.Value == "0" || lit.Value == "0.0"
}
