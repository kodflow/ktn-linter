package analyzer

import (
	"go/ast"
	"go/token"
)

// isZeroLiteral vérifie si une expression est le littéral 0.
//
// Params:
//   - expr: l'expression à vérifier
//
// Returns:
//   - bool: true si c'est le littéral "0"
func isZeroLiteral(expr ast.Expr) bool {
	basicLit, ok := expr.(*ast.BasicLit)
	// Retourne true si c'est le littéral entier "0"
	return ok && basicLit.Kind == token.INT && basicLit.Value == "0"
}
