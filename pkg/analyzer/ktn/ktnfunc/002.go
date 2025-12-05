// Analyzer 002 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer002 checks that context.Context is always the first parameter
var Analyzer002 = &analysis.Analyzer{
	Name:     "ktnfunc002",
	Doc:      "KTN-FUNC-002: context.Context doit toujours être le premier paramètre (après le receiver pour les méthodes)",
	Run:      runFunc002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc002 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc002(pass *analysis.Pass) (any, error) {
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
		if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
			// Retour de la fonction
			return
		}

		// Check each parameter
		contextParamIndex := -1
		// Itération sur les éléments
		for i, field := range funcDecl.Type.Params.List {
			// Vérification de la condition
			if isContextTypeWithPass(pass, field.Type) {
				contextParamIndex = i
				break
			}
		}

		// If there's a context parameter and it's not first, report error
		if contextParamIndex > 0 {
			pass.Reportf(
				funcDecl.Type.Params.List[contextParamIndex].Pos(),
				"KTN-FUNC-002: context.Context doit être le premier paramètre de la fonction '%s'",
				funcName,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// isContextTypeWithPass checks if a type is context.Context using type info.
// Handles both direct context.Context and type aliases.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression de type à vérifier
//
// Returns:
//   - bool: true si type context.Context ou alias
func isContextTypeWithPass(pass *analysis.Pass, expr ast.Expr) bool {
	// Try using type info first
	tv := pass.TypesInfo.Types[expr]
	if tv.Type != nil {
		return isContextTypeByType(tv.Type)
	}
	// Fallback to AST-based check
	return isContextType(expr)
}

// isContextTypeByType checks if a types.Type is context.Context.
//
// Params:
//   - t: type à vérifier
//
// Returns:
//   - bool: true si context.Context
func isContextTypeByType(t types.Type) bool {
	// Get the named type
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj == nil {
		return false
	}
	// Check if it's from context package
	if obj.Pkg() != nil && obj.Pkg().Path() == "context" && obj.Name() == "Context" {
		return true
	}
	// Check underlying type for aliases
	if obj.Pkg() != nil {
		underlying := t.Underlying()
		// Vérifier si le type sous-jacent est nommé
		if underNamed, ok := underlying.(*types.Named); ok {
			underObj := underNamed.Obj()
			// Vérifier context.Context sous-jacent
			if underObj != nil && underObj.Pkg() != nil {
				// Retour si match avec context.Context
				return underObj.Pkg().Path() == "context" && underObj.Name() == "Context"
			}
		}
	}
	// Retour false par défaut
	return false
}

// isContextType checks if a type is context.Context (AST-based fallback).
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si type context.Context
func isContextType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	// Vérification de la condition
	if !ok {
		// Retour de la fonction
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	// Retour de la fonction
	return ok && ident.Name == "context" && sel.Sel.Name == "Context"
}
