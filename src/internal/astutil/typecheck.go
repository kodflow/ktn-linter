package astutil

import (
	"go/ast"
)

// IsConstCompatibleType vérifie si le type est compatible avec const
// Retourne true pour les types de base : bool, string, int*, uint*, float*, complex*, byte, rune
func IsConstCompatibleType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Types de base compatibles avec const
		switch e.Name {
		case "bool", "string",
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64",
			"complex64", "complex128",
			"byte", "rune":
			return true
		}
	}
	return false
}

// IsLiteralValue vérifie si l'expression est une valeur littérale
// Retourne true pour les BasicLit (nombres, strings) et Ident (true, false, nil)
func IsLiteralValue(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.BasicLit: // true, false, 123, "string", 3.14, etc.
		return true
	case *ast.Ident: // true, false, nil
		return true
	}
	return false
}
