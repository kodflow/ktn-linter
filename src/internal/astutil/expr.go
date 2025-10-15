// Package astutil fournit des utilitaires pour manipuler l'AST Go
package astutil

import (
	"go/ast"
)

// ExprToString convertit une expression AST en sa représentation textuelle.
//
// Params:
//   - expr: l'expression AST à convertir
//
// Returns:
//   - string: la représentation textuelle (gère identifiants, sélecteurs, tableaux, maps, pointeurs, etc.)
func ExprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		// Retourne le nom de l'identifiant (ex: "int", "myVar")
		return e.Name
	case *ast.SelectorExpr:
		// Retourne la sélection complète (ex: "pkg.Type")
		return ExprToString(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		// Retourne la représentation du tableau (ex: "[]string")
		return "[]" + ExprToString(e.Elt)
	case *ast.MapType:
		// Retourne la représentation de la map (ex: "map[string]int")
		return "map[" + ExprToString(e.Key) + "]" + ExprToString(e.Value)
	case *ast.StarExpr:
		// Retourne la représentation du pointeur (ex: "*MyType")
		return "*" + ExprToString(e.X)
	case *ast.ChanType:
		// Gérer les différents types de channels
		switch e.Dir {
		case ast.SEND:
			// Retourne un channel send-only (ex: "chan<- int")
			return "chan<- " + ExprToString(e.Value)
		case ast.RECV:
			// Retourne un channel receive-only (ex: "<-chan int")
			return "<-chan " + ExprToString(e.Value)
		default:
			// Retourne un channel bidirectionnel (ex: "chan int")
			return "chan " + ExprToString(e.Value)
		}
	default:
		// Retourne "unknown" pour les types non gérés
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
	if spec.Type != nil {
		// Retourne la représentation textuelle du type spécifié
		return ExprToString(spec.Type)
	}
	// Retourne "<type>" quand le type n'est pas spécifié (inférence de type)
	return "<type>"
}
