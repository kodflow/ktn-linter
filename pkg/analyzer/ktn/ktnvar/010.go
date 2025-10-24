package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
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
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		compositeLit := n.(*ast.CompositeLit)

		// Vérification que c'est un slice vide
		if !utils.IsEmptySliceLiteral(compositeLit) {
			// Continue traversing AST nodes
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

