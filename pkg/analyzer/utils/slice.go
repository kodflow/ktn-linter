// Utility functions for slice operations.
package utils

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// IsSliceTypeWithPass vérifie si une expression est un type slice en utilisant TypesInfo.
//
// Params:
//   - pass: contexte d'analyse avec TypesInfo
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si c'est un type slice
func IsSliceTypeWithPass(pass *analysis.Pass, expr ast.Expr) bool {
	// Vérification via TypesInfo
	if tv, ok := pass.TypesInfo.Types[expr]; ok {
		// Check if underlying type is slice
		_, isSlice := tv.Type.Underlying().(*types.Slice)
		// Return result
		return isSlice
	}
	// Fallback to AST checking if TypesInfo not available
	return IsSliceType(expr)
}

// IsMapTypeWithPass vérifie si une expression est un type map en utilisant TypesInfo.
//
// Params:
//   - pass: contexte d'analyse avec TypesInfo
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si c'est un type map
func IsMapTypeWithPass(pass *analysis.Pass, expr ast.Expr) bool {
	// Check via TypesInfo
	if tv, ok := pass.TypesInfo.Types[expr]; ok {
		// Check if underlying type is map
		_, isMap := tv.Type.Underlying().(*types.Map)
		// Return result
		return isMap
	}
	// Fallback to AST checking
	return IsMapType(expr)
}

// IsMapType vérifie si une expression est un type map en utilisant l'AST.
//
// Params:
//   - expr: expression AST à vérifier
//
// Returns:
//   - bool: true si c'est un type map basé sur l'AST
func IsMapType(expr ast.Expr) bool {
	// Check if it's a map type in AST
	_, ok := expr.(*ast.MapType)
	// Return result
	return ok
}

// IsSliceOrMapTypeWithPass vérifie si une expression est un slice ou une map.
//
// Params:
//   - pass: contexte d'analyse avec TypesInfo
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si c'est un slice ou une map
func IsSliceOrMapTypeWithPass(pass *analysis.Pass, expr ast.Expr) bool {
	// Check for slice or map
	return IsSliceTypeWithPass(pass, expr) || IsMapTypeWithPass(pass, expr)
}

// IsSliceOrMapType vérifie si une expression est un slice ou une map via AST.
//
// Params:
//   - expr: expression AST à vérifier
//
// Returns:
//   - bool: true si c'est un slice ou une map basé sur l'AST
func IsSliceOrMapType(expr ast.Expr) bool {
	// Check for slice or map using AST
	return IsSliceType(expr) || IsMapType(expr)
}

// IsEmptySliceLiteral vérifie si un CompositeLit est un slice literal vide.
//
// Params:
//   - lit: composite literal à vérifier
//
// Returns:
//   - bool: true si c'est un slice literal vide ([]T{})
func IsEmptySliceLiteral(lit *ast.CompositeLit) bool {
	// Check if has no elements
	if len(lit.Elts) > 0 {
		// Has elements, not empty
		return false
	}
	// Check if type is a slice
	return IsSliceType(lit.Type)
}

// IsByteSliceWithPass vérifie si un type est []byte.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression de type à vérifier
//
// Returns:
//   - bool: true si c'est un []byte
func IsByteSliceWithPass(pass *analysis.Pass, expr ast.Expr) bool {
	var slice *types.Slice
	var basic *types.Basic
	var sliceOk bool
	var basicOk bool
	// Check via TypesInfo
	if tv, ok := pass.TypesInfo.Types[expr]; ok {
		// Check if it's a slice
		slice, sliceOk = tv.Type.Underlying().(*types.Slice)
		// Vérification de la condition
		if sliceOk {
			// Check if element type is byte
			basic, basicOk = slice.Elem().(*types.Basic)
			// Vérification de la condition
			if basicOk {
				// Byte is uint8
				return basic.Kind() == types.Byte || basic.Kind() == types.Uint8
			}
		}
	}
	// Fallback to AST checking
	return IsByteSlice(expr)
}

// IsByteSlice vérifie si un type est []byte via l'AST.
//
// Params:
//   - expr: expression AST à vérifier
//
// Returns:
//   - bool: true si c'est un []byte basé sur l'AST
func IsByteSlice(expr ast.Expr) bool {
	var ident *ast.Ident
	var identOk bool
	// Check if it's an array type (slice)
	arrayType, ok := expr.(*ast.ArrayType)
	// Vérification de la condition
	if !ok || arrayType.Len != nil {
		// Not a slice
		return false
	}
	// Check if element type is byte
	ident, identOk = arrayType.Elt.(*ast.Ident)
	// Vérification de la condition
	if identOk {
		// Check for byte or uint8
		return ident.Name == "byte" || ident.Name == "uint8"
	}
	// Not a byte slice
	return false
}
