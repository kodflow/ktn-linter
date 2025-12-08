// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer006 checks that functions don't have too many parameters
const (
	// MAX_PARAMS max params allowed in a function (context.Context excluded)
	MAX_PARAMS int = 5
)

// Analyzer006 checks that functions don't have more than MAX_PARAMS parameters
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc006",
	Doc:      "KTN-FUNC-006: Les fonctions ne doivent pas dépasser 5 paramètres",
	Run:      runFunc006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc006 exécute l'analyse KTN-FUNC-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc006(pass *analysis.Pass) (any, error) {
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

		// Count total parameters (excluding context.Context)
		paramCount := countEffectiveParams(pass, funcType.Params)

		// Vérification de la condition
		if paramCount > MAX_PARAMS {
			pass.Reportf(
				pos.Pos(),
				"KTN-FUNC-006: la fonction '%s' a %d paramètres (max: %d)",
				name,
				paramCount,
				MAX_PARAMS,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// countEffectiveParams counts parameters excluding context.Context.
// context.Context is excluded because KTN-FUNC-002 already mandates it as first param.
//
// Params:
//   - pass: contexte d'analyse
//   - params: field list of function parameters
//
// Returns:
//   - int: effective parameter count
func countEffectiveParams(pass *analysis.Pass, params *ast.FieldList) int {
	// Check for nil params
	if params == nil {
		// Retour 0 si pas de paramètres
		return 0
	}

	count := 0
	// Iterate over parameter fields
	for _, field := range params.List {
		// Skip context.Context parameters (including aliases)
		if isContextTypeWithPass(pass, field.Type) {
			continue
		}

		// Each field can declare multiple params: func(a, b, c int)
		if len(field.Names) > 0 {
			count += len(field.Names)
		} else {
			// Unnamed parameter (e.g., in interface or func literal)
			count++
		}
	}

	// Retour du compte effectif
	return count
}
