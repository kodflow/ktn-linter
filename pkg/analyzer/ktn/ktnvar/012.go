package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer012 checks for string concatenation in loops
var Analyzer012 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar012",
	Doc:      "KTN-VAR-012: Utiliser strings.Builder pour >2 concaténations",
	Run:      runVar012,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar012 exécute l'analyse KTN-VAR-012.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar012(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		checkStringConcatInLoop(pass, n)
	})

	// Retour de la fonction
	return nil, nil
}

// checkStringConcatInLoop checks for string += in loops.
//
// Params:
//   - pass: analysis pass context
//   - n: AST node (ForStmt or RangeStmt)
func checkStringConcatInLoop(pass *analysis.Pass, n ast.Node) {
	// Body of the loop
	var loopBody *ast.BlockStmt

	// Check type of loop
	switch node := n.(type) {
	// Case: for loop
	case *ast.ForStmt:
		loopBody = node.Body
	// Case: range loop
	case *ast.RangeStmt:
		loopBody = node.Body
	}

	// Check if loop body exists
	if loopBody == nil {
		// Return early if no body
		return
	}

	// Traverse statements in loop body
	ast.Inspect(loopBody, func(child ast.Node) bool {
		// Check if assignment statement
		if assign, ok := child.(*ast.AssignStmt); ok {
			// Check if += operator
			if assign.Tok == token.ADD_ASSIGN {
				// Check if string concatenation
				if isStringConcatenation(pass, assign) {
					pass.Reportf(
						assign.Pos(),
						"KTN-VAR-012: utiliser strings.Builder au lieu de += pour concaténer des strings dans une boucle",
					)
				}
			}
		}
		// Continue traversing
		return true
	})
}

// isStringConcatenation checks if += operates on string.
//
// Params:
//   - pass: analysis pass context
//   - assign: assignment statement to check
//
// Returns:
//   - bool: true if string concatenation
func isStringConcatenation(pass *analysis.Pass, assign *ast.AssignStmt) bool {
	// Check if left-hand side exists
	if len(assign.Lhs) == 0 {
		// Return false if no left-hand side
		return false
	}

	// Get first left-hand side expression
	lhs := assign.Lhs[0]

	// Get type information
	if tv, ok := pass.TypesInfo.Types[lhs]; ok {
		var basic *types.Basic
		// Check if type is string
		if basic, ok = tv.Type.Underlying().(*types.Basic); ok {
			// Return true if string type
			return basic.Kind() == types.String
		}
	}

	// Return false if not string
	return false
}
