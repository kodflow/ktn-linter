package astutil

import (
	"go/ast"
)

// IsConstCompatibleType vérifie si le type est compatible avec const.
//
// Params:
//   - expr: l'expression représentant le type à vérifier
//
// Returns:
//   - bool: true pour les types de base (bool, string, int*, uint*, float*, complex*, byte, rune)
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
			// Retourne true car c'est un type primitif compatible avec const
			return true
		}
	}
	// Retourne false car le type n'est pas un type primitif ou n'est pas compatible avec const
	return false
}

// IsLiteralValue vérifie si l'expression est une valeur littérale.
//
// Params:
//   - expr: l'expression à vérifier
//
// Returns:
//   - bool: true pour les BasicLit (nombres, strings) et Ident (true, false, nil)
func IsLiteralValue(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.BasicLit: // true, false, 123, "string", 3.14, etc.
		// Retourne true car c'est un littéral de base (nombre, string, etc.)
		return true
	case *ast.Ident: // true, false, nil
		// Retourne true car c'est un identifiant (peut être true/false/nil)
		return true
	}
	// Retourne false car l'expression n'est pas une valeur littérale
	return false
}
