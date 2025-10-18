package utils

import (
	"go/ast"
	"go/token"
	"strings"
)

// IsZeroLiteral vérifie si une expression est le literal zéro (0).
//
// Params:
//   - expr: l'expression à vérifier
//
// Returns:
//   - bool: true si c'est le literal 0
func IsZeroLiteral(expr ast.Expr) bool {
	basicLit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	return basicLit.Kind == token.INT && basicLit.Value == "0"
}

// IsReferenceType vérifie si un type est un type référence (slice/map/chan).
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - bool: true si c'est un slice, map ou channel
func IsReferenceType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.ArrayType:
		// Slice (ArrayType sans longueur)
		return t.Len == nil
	case *ast.MapType:
		return true
	case *ast.ChanType:
		return true
	case *ast.Ident:
		// Vérifier les types natifs par nom
		return strings.Contains(t.Name, "map") ||
			strings.Contains(t.Name, "chan") ||
			strings.Contains(t.Name, "slice")
	}
	return false
}

// IsStructType vérifie si un type est une struct.
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - bool: true si c'est une struct
func IsStructType(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.StructType:
		return true
	case *ast.Ident:
		// Identifiant (type nommé) - potentiellement une struct
		return true
	case *ast.SelectorExpr:
		// Type importé (ex: pkg.MyStruct)
		return true
	}
	return false
}

// IsSliceType vérifie si un type est un slice.
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - bool: true si c'est un slice
func IsSliceType(expr ast.Expr) bool {
	arrayType, ok := expr.(*ast.ArrayType)
	return ok && arrayType.Len == nil
}

// GetTypeName extrait le nom d'un type depuis une expression.
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - string: le nom du type (ex: "map[string]int", "[]int", "chan int")
func GetTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.ArrayType:
		elemType := GetTypeName(t.Elt)
		return "[]" + elemType
	case *ast.MapType:
		keyType := GetTypeName(t.Key)
		valueType := GetTypeName(t.Value)
		return "map[" + keyType + "]" + valueType
	case *ast.ChanType:
		elemType := GetTypeName(t.Value)
		return "chan " + elemType
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		pkg := GetTypeName(t.X)
		return pkg + "." + t.Sel.Name
	case *ast.StarExpr:
		base := GetTypeName(t.X)
		return "*" + base
	}
	return "T"
}

// IsMakeSliceZero vérifie si une expression est make([]T, 0) ou make([]T, 0, 0).
//
// Params:
//   - expr: l'expression à vérifier
//
// Returns:
//   - bool: true si c'est make([]T, 0) ou make([]T, 0, 0)
func IsMakeSliceZero(expr ast.Expr) bool {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	// Vérifier si c'est make()
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok || ident.Name != "make" {
		return false
	}

	// Vérifier si c'est un slice
	if len(callExpr.Args) < 1 || !IsSliceType(callExpr.Args[0]) {
		return false
	}

	// Vérifier les arguments
	if len(callExpr.Args) == 2 {
		// make([]T, length)
		return IsZeroLiteral(callExpr.Args[1])
	} else if len(callExpr.Args) == 3 {
		// make([]T, length, capacity)
		return IsZeroLiteral(callExpr.Args[1]) && IsZeroLiteral(callExpr.Args[2])
	}

	return false
}
