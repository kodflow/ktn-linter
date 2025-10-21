package ktnfunc

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer006 checks that error is always the last return value
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc006",
	Doc:      "KTN-FUNC-006: L'erreur doit toujours être en dernière position dans les valeurs de retour",
	Run:      runFunc006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// validateErrorInReturns vérifie que l'erreur est en dernière position.
//
// Params:
//   - pass: contexte d'analyse
//   - funcType: type de la fonction
func validateErrorInReturns(pass *analysis.Pass, funcType *ast.FuncType) {
	// Vérification présence de résultats
	if funcType == nil || funcType.Results == nil {
		// Pas de résultats à vérifier
		return
	}

	results := funcType.Results.List

	// Recherche des positions d'erreur
	var errorPositions []int
	// Itération sur les résultats
	for i, result := range results {
		// Vérification si type error
		if isErrorType(pass, result.Type) {
			errorPositions = append(errorPositions, i)
		}
	}

	// Vérification erreurs mal placées
	if len(errorPositions) > 0 {
		lastPos := len(results) - 1
		// Itération sur les positions d'erreur
		for _, pos := range errorPositions {
			// Vérification position incorrecte
			if pos != lastPos {
				pass.Reportf(
					funcType.Results.Pos(),
					"KTN-FUNC-006: l'erreur doit être en dernière position dans les valeurs de retour (trouvée en position %d sur %d)",
					pos+1,
					len(results),
				)
				// Retour après premier rapport
				return
			}
		}
	}
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
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		var funcType *ast.FuncType
		// Sélection selon la valeur
		switch node := n.(type) {
		// Traitement FuncDecl
		case *ast.FuncDecl:
			funcType = node.Type
		// Traitement FuncLit
		case *ast.FuncLit:
			funcType = node.Type
		}
		validateErrorInReturns(pass, funcType)
	})

	// Retour de la fonction
	return nil, nil
}

// isErrorType checks if a type expression represents the error interface
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si type error
//
func isErrorType(pass *analysis.Pass, expr ast.Expr) bool {
	tv := pass.TypesInfo.Types[expr]

	// Check if it's the error interface
	named, ok := tv.Type.(*types.Named)
 // Vérification de la condition
	if !ok {
  // Retour de la fonction
		return false
	}

	obj := named.Obj()
 // Retour de la fonction
	return obj != nil && obj.Name() == "error" && obj.Pkg() == nil
}
