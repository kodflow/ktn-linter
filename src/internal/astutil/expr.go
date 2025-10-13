// Package astutil fournit des utilitaires pour manipuler l'AST Go
package astutil

import (
	"go/ast"
)

// ExprToString convertit une expression AST en sa représentation textuelle
// Gère les identifiants, sélecteurs, tableaux, maps, pointeurs, etc.
func ExprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return ExprToString(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		return "[]" + ExprToString(e.Elt)
	case *ast.MapType:
		return "map[" + ExprToString(e.Key) + "]" + ExprToString(e.Value)
	case *ast.StarExpr:
		return "*" + ExprToString(e.X)
	case *ast.ChanType:
		// Gérer les différents types de channels
		switch e.Dir {
		case ast.SEND:
			return "chan<- " + ExprToString(e.Value)
		case ast.RECV:
			return "<-chan " + ExprToString(e.Value)
		default:
			return "chan " + ExprToString(e.Value)
		}
	default:
		return "unknown"
	}
}

// GetTypeString extrait la représentation textuelle du type d'une ValueSpec
// Si le type n'est pas spécifié, retourne "<type>"
func GetTypeString(spec *ast.ValueSpec) string {
	if spec.Type != nil {
		return ExprToString(spec.Type)
	}
	return "<type>"
}
