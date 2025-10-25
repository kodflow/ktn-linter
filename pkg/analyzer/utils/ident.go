package utils

import (
	"go/ast"
)

// IsIdentCall vérifie si un CallExpr est un appel à une fonction avec un nom spécifique.
//
// Params:
//   - call: expression d'appel à vérifier
//   - name: nom de la fonction attendue
//
// Returns:
//   - bool: true si c'est un appel à la fonction nommée
func IsIdentCall(call *ast.CallExpr, name string) bool {
	// Check if Fun is an identifier
	ident, ok := call.Fun.(*ast.Ident)
	// Return true if identifier matches name
	return ok && ident.Name == name
}

// IsBuiltinCall vérifie si c'est un appel à une fonction builtin.
//
// Params:
//   - call: expression d'appel à vérifier
//   - builtinName: nom de la fonction builtin
//
// Returns:
//   - bool: true si c'est un appel au builtin spécifié
func IsBuiltinCall(call *ast.CallExpr, builtinName string) bool {
	// Builtins are also identifiers
	return IsIdentCall(call, builtinName)
}

// GetIdentName extrait le nom d'un identifiant depuis une expression.
//
// Params:
//   - expr: expression à analyser
//
// Returns:
//   - string: nom de l'identifiant ou chaîne vide
func GetIdentName(expr ast.Expr) string {
	// Check if expression is identifier
	if ident, ok := expr.(*ast.Ident); ok {
		// Return identifier name
		return ident.Name
	}
	// Not an identifier
	return ""
}

// ExtractVarName extrait le nom d'une variable depuis une expression complexe.
//
// Params:
//   - expr: expression à analyser
//
// Returns:
//   - string: nom de la variable ou représentation
func ExtractVarName(expr ast.Expr) string {
	// Switch on expression type
	switch e := expr.(type) {
	// Simple identifier
	case *ast.Ident:
		// Return name
		return e.Name
	// Array/slice index: items[i]
	case *ast.IndexExpr:
		// Get base name and add index indicator
		if base := ExtractVarName(e.X); base != "" {
			// Return with index notation
			return base + "[...]"
		}
		// Cannot extract
		return ""
	// Struct field: obj.field
	case *ast.SelectorExpr:
		// Get base and field
		if base := ExtractVarName(e.X); base != "" {
			// Return with field
			return base + "." + e.Sel.Name
		}
		// Just field name
		return e.Sel.Name
	// Star expression: *ptr
	case *ast.StarExpr:
		// Get pointed expression
		if base := ExtractVarName(e.X); base != "" {
			// Return with dereference
			return "*" + base
		}
		// Cannot extract
		return ""
	// Parenthesized expression
	case *ast.ParenExpr:
		// Extract from inner
		return ExtractVarName(e.X)
	// Default case
	default:
		// Cannot extract variable name
		return ""
	}
}
