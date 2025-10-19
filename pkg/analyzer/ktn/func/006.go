package ktnfunc

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer006 checks that error is always the last return value
var Analyzer006 = &analysis.Analyzer{
	Name:     "ktnfunc006",
	Doc:      "KTN-FUNC-006: L'erreur doit toujours être en dernière position dans les valeurs de retour",
	Run:      runFunc006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

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
  // Traitement
		case *ast.FuncDecl:
			funcType = node.Type
  // Traitement
		case *ast.FuncLit:
			funcType = node.Type
		}

  // Vérification de la condition
		if funcType == nil || funcType.Results == nil {
   // Retour de la fonction
			return
		}

		results := funcType.Results.List
  // Vérification de la condition
		if len(results) == 0 {
   // Retour de la fonction
			return
		}

		// Check each return value to find error types
		var errorPositions []int
  // Itération sur les éléments
		for i, result := range results {
   // Vérification de la condition
			if isErrorType(pass, result.Type) {
				errorPositions = append(errorPositions, i)
			}
		}

		// If there are errors and any is not in the last position, report
		if len(errorPositions) > 0 {
			lastPos := len(results) - 1
   // Itération sur les éléments
			for _, pos := range errorPositions {
    // Vérification de la condition
				if pos != lastPos {
					pass.Reportf(
						funcType.Results.Pos(),
						"KTN-FUNC-006: l'erreur doit être en dernière position dans les valeurs de retour (trouvée en position %d sur %d)",
						pos+1,
						len(results),
					)
     // Retour de la fonction
					return
				}
			}
		}
	})

 // Retour de la fonction
	return nil, nil
}

// isErrorType checks if a type expression represents the error interface
func isErrorType(pass *analysis.Pass, expr ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[expr]
 // Vérification de la condition
	if !ok {
  // Retour de la fonction
		return false
	}

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
