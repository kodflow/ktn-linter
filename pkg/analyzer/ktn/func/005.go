package ktnfunc

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer005 checks that functions don't exceed cyclomatic complexity of 10
const (
	// MAX_CYCLOMATIC_COMPLEXITY définit la complexité cyclomatique maximale autorisée
	MAX_CYCLOMATIC_COMPLEXITY int = 10
)

var Analyzer005 = &analysis.Analyzer{
	Name:     "ktnfunc005",
	Doc:      "KTN-FUNC-005: La complexité cyclomatique ne doit pas dépasser 10",
	Run:      runFunc005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}


func runFunc005(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip if no body
		if funcDecl.Body == nil {
   // Retour de la fonction
			return
		}

		// Skip test functions
		funcName := funcDecl.Name.Name
  // Vérification de la condition
		if isTestFunction(funcName) {
   // Retour de la fonction
			return
		}

		// Calculate cyclomatic complexity
		complexity := calculateComplexity(funcDecl.Body)

  // Vérification de la condition
		if complexity > MAX_CYCLOMATIC_COMPLEXITY {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-005: la fonction '%s' a une complexité cyclomatique de %d (max: %d)",
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
