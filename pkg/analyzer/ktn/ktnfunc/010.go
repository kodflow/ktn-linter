// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer010 checks that functions don't use naked returns (except for very short functions)
const (
	// MAX_LINES_FOR_NAKED_RETURN max lines for naked return
	MAX_LINES_FOR_NAKED_RETURN int = 5
)

// Analyzer010 checks that naked returns are only used in very short functions
var Analyzer010 = &analysis.Analyzer{
	Name:     "ktnfunc010",
	Doc:      "KTN-FUNC-010: Les naked returns sont interdits sauf pour les fonctions très courtes (<5 lignes)",
	Run:      runFunc010,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc010 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc010(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip if no body (external functions)
		if funcDecl.Body == nil {
			// Retour de la fonction
			return
		}

		// Skip test functions
		if shared.IsTestFunction(funcDecl) {
			// Retour de la fonction
			return
		}

		funcName := funcDecl.Name.Name

		// Skip if function doesn't have named return values
		if !hasNamedReturns(funcDecl.Type.Results) {
			// Retour de la fonction
			return
		}

		// Count the lines of the function
		pureLines := countPureCodeLines(pass, funcDecl.Body)

		// Check for naked returns
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			ret, ok := node.(*ast.ReturnStmt)
			// Vérification de la condition
			if !ok {
				// Retour de la fonction
				return true
			}

			// Naked return has no results specified
			if len(ret.Results) == 0 {
				// Allow naked returns in very short functions
				if pureLines >= MAX_LINES_FOR_NAKED_RETURN {
					pass.Reportf(
						ret.Pos(),
						"KTN-FUNC-010: naked return interdit dans la fonction '%s' (%d lignes, max: %d pour naked return)",
						funcName,
						pureLines,
						MAX_LINES_FOR_NAKED_RETURN-1,
					)
				}
			}

			// Retour de la fonction
			return true
		})
	})

	// Retour de la fonction
	return nil, nil
}

// hasNamedReturns checks if the function has named return values
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si retours nommés
func hasNamedReturns(results *ast.FieldList) bool {
	// Vérification de la condition
	if results == nil || len(results.List) == 0 {
		// Retour de la fonction
		return false
	}

	// Itération sur les éléments
	for _, field := range results.List {
		// Vérification de la condition
		if len(field.Names) > 0 {
			// Retour de la fonction
			return true
		}
	}

	// Retour de la fonction
	return false
}
