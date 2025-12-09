// Utility functions for types operations.
package utils

import (
	"go/ast"
	"go/constant"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Constantes pour les nombres utilisés dans les fonctions
const (
	makeArgsWithLength   int = 2 // make([]T, length)
	makeArgsWithCapacity int = 3 // make([]T, length, capacity)
	secondArgIndex       int = 1 // index du deuxième argument
	thirdArgIndex        int = 2 // index du troisième argument
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
	// Vérification de la condition
	if !ok {
		// Condition not met, return false.
		return false
	}
	// Early return from function.
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
	// Sélection selon la valeur
	switch t := expr.(type) {
	// Traitement
	case *ast.ArrayType:
		// Slice (ArrayType sans longueur)
		return t.Len == nil
	// Traitement
	case *ast.MapType:
		// Continue traversing AST nodes.
		return true
	// Traitement
	case *ast.ChanType:
		// Continue traversing AST nodes.
		return true
	// Traitement
	case *ast.Ident:
		// Vérifier les types natifs par nom
		return strings.Contains(t.Name, "map") ||
			strings.Contains(t.Name, "chan") ||
			strings.Contains(t.Name, "slice")
	}
	// Condition not met, return false.
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
	// Sélection selon la valeur
	switch expr.(type) {
	// Traitement
	case *ast.StructType:
		// Continue traversing AST nodes.
		return true
	// Traitement
	case *ast.Ident:
		// Identifiant (type nommé) - potentiellement une struct
		return true
	// Traitement
	case *ast.SelectorExpr:
		// Type importé (ex: pkg.MyStruct)
		return true
	}
	// Condition not met, return false.
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
	// Early return from function.
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
	// Sélection selon la valeur
	switch t := expr.(type) {
	// Traitement
	case *ast.ArrayType:
		elemType := GetTypeName(t.Elt)
		// Early return from function.
		return "[]" + elemType
	// Traitement
	case *ast.MapType:
		keyType := GetTypeName(t.Key)
		valueType := GetTypeName(t.Value)
		// Early return from function.
		return "map[" + keyType + "]" + valueType
	// Traitement
	case *ast.ChanType:
		elemType := GetTypeName(t.Value)
		// Early return from function.
		return "chan " + elemType
	// Traitement
	case *ast.Ident:
		// Early return from function.
		return t.Name
	// Traitement
	case *ast.SelectorExpr:
		pkg := GetTypeName(t.X)
		// Early return from function.
		return pkg + "." + t.Sel.Name
	// Traitement
	case *ast.StarExpr:
		base := GetTypeName(t.X)
		// Early return from function.
		return "*" + base
	}
	// Early return from function.
	return "T"
}

// HasPositiveLength vérifie si un argument de longueur est > 0.
//
// Params:
//   - pass: contexte d'analyse
//   - lengthArg: argument de longueur
//
// Returns:
//   - bool: true si la longueur est positive
func HasPositiveLength(pass *analysis.Pass, lengthArg ast.Expr) bool {
	// Check via TypesInfo
	if pass != nil && pass.TypesInfo != nil {
		// Vérification du type de l'argument
		if tv, ok := pass.TypesInfo.Types[lengthArg]; ok {
			// Check if it's a constant
			if tv.Value != nil && tv.Value.Kind() == constant.Int {
				var val int64
				var isInt bool
				// Get integer value
				val, isInt = constant.Int64Val(tv.Value)
				// Vérification de la condition
				if isInt {
					// Return true if positive
					return val > 0
				}
			}
		}
	}
	// Check via AST for basic literals
	if lit, ok := lengthArg.(*ast.BasicLit); ok && lit.Kind == token.INT {
		// Return true if not "0"
		return lit.Value != "0"
	}
	// Default to true (assume positive for variables)
	return true
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
	// Vérification de la condition
	if !ok {
		// Condition not met, return false.
		return false
	}

	// Vérifier si c'est make()
	ident, ok := callExpr.Fun.(*ast.Ident)
	// Vérification de la condition
	if !ok || ident.Name != "make" {
		// Condition not met, return false.
		return false
	}

	// Vérifier si c'est un slice
	if len(callExpr.Args) < 1 || !IsSliceType(callExpr.Args[0]) {
		// Condition not met, return false.
		return false
	}

	// Vérifier les arguments
	if len(callExpr.Args) == makeArgsWithLength {
		// make([]T, length) - vérifier que length est 0
		return IsZeroLiteral(callExpr.Args[secondArgIndex])
	}

	// Vérification de la condition
	if len(callExpr.Args) == makeArgsWithCapacity {
		// make([]T, length, capacity) - vérifier que length et capacity sont 0
		return IsZeroLiteral(callExpr.Args[secondArgIndex]) && IsZeroLiteral(callExpr.Args[thirdArgIndex])
	}

	// Condition not met, return false.
	return false
}
