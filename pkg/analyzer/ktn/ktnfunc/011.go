// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer011 checks that functions don't exceed cyclomatic complexity of 10
const (
	// MAX_CYCLOMATIC_COMPLEXITY max cyclomatic complexity
	MAX_CYCLOMATIC_COMPLEXITY int = 10
)

// Analyzer011 checks that functions don't exceed maximum cyclomatic complexity
var Analyzer011 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc011",
	Doc:      "KTN-FUNC-011: La complexité cyclomatique ne doit pas dépasser 10",
	Run:      runFunc011,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc011 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc011(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// runFunc011 exécute l'analyse KTN-FUNC-011.
		//
		// Params:
		//   - pass: contexte d'analyse
		//
		// Returns:
		//   - any: résultat de l'analyse
		//   - error: erreur éventuelle
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

		// Calculate cyclomatic complexity
		complexity := calculateComplexity(funcDecl.Body)

		// Vérification de la condition
		if complexity > MAX_CYCLOMATIC_COMPLEXITY {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-011: la fonction '%s' a une complexité cyclomatique de %d (max: %d)",
				funcName,
				complexity,
				MAX_CYCLOMATIC_COMPLEXITY,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// calculateComplexity calculates the cyclomatic complexity of a function
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - int: complexité calculée
func calculateComplexity(body *ast.BlockStmt) int {
	// Start with complexity of 1 (the function itself)
	complexity := 1

	ast.Inspect(body, func(n ast.Node) bool {
		// Sélection selon la valeur
		switch node := n.(type) {
		// Traitement
		case *ast.IfStmt:
			// +1 for if
			complexity++
		// Traitement
		case *ast.ForStmt, *ast.RangeStmt:
			// +1 for each loop
			complexity++
		// Traitement
		case *ast.CaseClause:
			// +1 for each case (except default)
			if node.List != nil {
				complexity++
			}
		// Traitement
		case *ast.CommClause:
			// +1 for each comm case in select
			if node.Comm != nil {
				complexity++
			}
		// Traitement
		case *ast.BinaryExpr:
			// +1 for && and ||
			if node.Op.String() == "&&" || node.Op.String() == "||" {
				complexity++
			}
		}
		// Retour de la fonction
		return true
	})

	// Retour de la fonction
	return complexity
}
