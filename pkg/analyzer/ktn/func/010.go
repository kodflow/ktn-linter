package ktnfunc

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer010 checks that functions with >3 return values use named returns
const (
	// MAX_UNNAMED_RETURNS définit le nombre maximum de retours non-nommés autorisés
	MAX_UNNAMED_RETURNS int = 3
)

// Analyzer010 checks that functions with >3 return values use named returns
var Analyzer010 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc010",
	Doc:      "KTN-FUNC-010: Les fonctions avec plus de 3 valeurs de retour doivent utiliser des named returns",
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
		funcName := funcDecl.Name.Name

		// Skip test functions
		if isTestFunction(funcName) {
   // Retour de la fonction
			return
		}

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
				returnCount++
			} else {
				// Named returns
				returnCount += len(field.Names)
			}
		}

		// If more than 3 returns and has unnamed returns, report error
		if returnCount > MAX_UNNAMED_RETURNS && hasUnnamedReturns {
			pass.Reportf(
				funcDecl.Type.Results.Pos(),
				"KTN-FUNC-010: la fonction '%s' a %d valeurs de retour et doit utiliser des named returns (max %d sans noms)",
				funcName,
				returnCount,
				MAX_UNNAMED_RETURNS,
			)
		}
	})

 // Retour de la fonction
	return nil, nil
}
