// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer002 checks that functions don't have too many parameters
const (
	// MAX_PARAMS max params allowed in a function
	MAX_PARAMS int = 5
)

// Analyzer002 checks that functions don't have more than MAX_PARAMS parameters
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc002",
	Doc:      "KTN-FUNC-002: Les fonctions ne doivent pas dépasser 5 paramètres",
	Run:      runFunc002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc002 exécute l'analyse KTN-FUNC-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc002(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		var funcType *ast.FuncType
		var pos ast.Node
		var name string

		// Sélection selon la valeur
		switch fn := n.(type) {
		// Traitement
		case *ast.FuncDecl:
			funcType = fn.Type
			pos = fn.Name
			name = fn.Name.Name

			// Skip test functions
			if shared.IsTestFunction(fn) {
				// Retour de la fonction
				return
			}
		// Traitement
		case *ast.FuncLit:
			funcType = fn.Type
			pos = fn
			name = "function literal"
		}

		// Count total parameters
		paramCount := 0
		// Itération sur les éléments
		for _, field := range funcType.Params.List {
			// Each field can declare multiple params: func(a, b, c int)
			if len(field.Names) > 0 {
				paramCount += len(field.Names)
			} else {
				// Unnamed parameter (e.g., in interface or func literal)
				paramCount++
			}
		}

		// Vérification de la condition
		if paramCount > MAX_PARAMS {
			pass.Reportf(
				pos.Pos(),
				"KTN-FUNC-002: la fonction '%s' a %d paramètres (max: %d)",
				name,
				paramCount,
				MAX_PARAMS,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}
