// Package utils provides AST utility functions for analyzers.
package utils

import (
	"go/ast"
)

// GetExprAsString convertit une expression AST en sa représentation textuelle.
//
// Params:
//   - expr: l'expression AST à convertir
//
// Returns:
//   - string: la représentation textuelle (gère identifiants, sélecteurs, tableaux, maps, pointeurs, etc.)
func GetExprAsString(expr ast.Expr) string {
	// Sélection selon la valeur
	switch exprType := expr.(type) {
	// Traitement
	case *ast.Ident:
		// Early return from function.
		return exprType.Name
	// Traitement
	case *ast.SelectorExpr:
		// Early return from function.
		return GetExprAsString(exprType.X) + "." + exprType.Sel.Name
	// Traitement
	case *ast.ArrayType:
		// Early return from function.
		return "[]" + GetExprAsString(exprType.Elt)
	// Traitement
	case *ast.MapType:
		// Early return from function.
		return "map[" + GetExprAsString(exprType.Key) + "]" + GetExprAsString(exprType.Value)
	// Traitement
	case *ast.StarExpr:
		// Early return from function.
		return "*" + GetExprAsString(exprType.X)
	// Traitement
	case *ast.ChanType:
		// Sélection selon la valeur
		switch exprType.Dir {
		// Traitement
		case ast.SEND:
			// Early return from function.
			return "chan<- " + GetExprAsString(exprType.Value)
		// Traitement
		case ast.RECV:
			// Early return from function.
			return "<-chan " + GetExprAsString(exprType.Value)
		// Traitement
		default:
			// Early return from function.
			return "chan " + GetExprAsString(exprType.Value)
		}
	// Traitement
	default:
		// Early return from function.
		return "unknown"
	}
}

// GetTypeString extrait la représentation textuelle du type d'une ValueSpec.
//
// Params:
//   - spec: la ValueSpec dont on veut extraire le type
//
// Returns:
//   - string: la représentation textuelle du type, ou "<type>" si non spécifié
func GetTypeString(spec *ast.ValueSpec) string {
	// Vérification de la condition
	if spec.Type != nil {
		// Early return from function.
		return GetExprAsString(spec.Type)
	}
	// Early return from function.
	return "<type>"
}
