// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer012 checks that functions with >3 return values use named returns
const (
	// maxUnnamedReturns max unnamed returns allowed
	maxUnnamedReturns int = 3
)

// Analyzer012 checks that functions with >3 return values use named returns
var Analyzer012 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc012",
	Doc:      "KTN-FUNC-012: Les fonctions avec plus de 3 valeurs de retour doivent utiliser des named returns",
	Run:      runFunc012,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc012 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc012(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip test functions
		if shared.IsTestFunction(funcDecl) {
			// Retour de la fonction
			return
		}

		funcName := funcDecl.Name.Name

		// Vérification de la condition
		if funcDecl.Type.Results == nil {
			// Retour de la fonction
			return
		}

		// Count total return values
		returnCount := 0
		hasUnnamedReturns := false

		// Itération sur les éléments
		for _, field := range funcDecl.Type.Results.List {
			// Vérification de la condition
			if len(field.Names) == 0 {
				// Unnamed return
				hasUnnamedReturns = true
				// Incrément du compteur
				returnCount++
			} else {
				// Named returns
				// Ajout du nombre de retours nommés
				returnCount += len(field.Names)
			}
		}

		// If more than 3 returns and has unnamed returns, report error
		if returnCount > maxUnnamedReturns && hasUnnamedReturns {
			// Rapport d'erreur pour named returns requis
			pass.Reportf(
				funcDecl.Type.Results.Pos(),
				"KTN-FUNC-012: la fonction '%s' a %d valeurs de retour et doit utiliser des named returns (max %d sans noms)",
				funcName,
				returnCount,
				maxUnnamedReturns,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}
