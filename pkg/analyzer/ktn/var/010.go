package ktnvar

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer010 detects empty slice literals that should use nil
var Analyzer010 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar010",
	Doc:      "KTN-VAR-010: Préférer var slice []T (nil) pour zero-value",
	Run:      runVar010,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar010 exécute l'analyse KTN-VAR-010.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar010(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		compositeLit := n.(*ast.CompositeLit)

		// Vérification que c'est un slice type
		if !isEmptySliceLiteral(compositeLit) {
			// Continue traversing AST nodes
			return
		}

		// Vérification que le slice est vide
		if len(compositeLit.Elts) != 0 {
			// Slice non vide, pas d'erreur
			return
		}

		// Signaler l'erreur
		pass.Reportf(
			compositeLit.Pos(),
			"KTN-VAR-010: préférer 'var slice []T' (nil) au lieu de '[]T{}' (économise allocation)",
		)
	})

	// Retour de la fonction
	return nil, nil
}

// isEmptySliceLiteral vérifie si le CompositeLit est un slice.
//
// Params:
//   - lit: composite literal à vérifier
//
// Returns:
//   - bool: true si c'est un slice type
func isEmptySliceLiteral(lit *ast.CompositeLit) bool {
	// Vérification que le type n'est pas nil
	if lit.Type == nil {
		// Retour false si type nil
		return false
	}

	// Vérification que c'est un ArrayType
	arrayType, ok := lit.Type.(*ast.ArrayType)
	// Vérification du type array
	if !ok {
		// Retour false si ce n'est pas un array
		return false
	}

	// Un slice a Len == nil, un array a une taille explicite
	// Retour true si c'est un slice (Len == nil)
	return arrayType.Len == nil
}
