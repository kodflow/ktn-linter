// Analyzer 016 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/constant"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// minMakeArgsVar016 is minimum arguments for make([]T, N)
	minMakeArgsVar016 int = 2
	// minMakeArgsWithCapVar016 is minimum arguments for make with capacity
	minMakeArgsWithCapVar016 int = 3
	// maxArraySizeVar016 is maximum size for recommending array over slice
	maxArraySizeVar016 int64 = 1024
)

// Analyzer016 checks for make([]T, N) with small constant N
var Analyzer016 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar016",
	Doc:      "KTN-VAR-016: Vérifie l'utilisation de [N]T au lieu de make([]T, N)",
	Run:      runVar016,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar016 exécute l'analyse KTN-VAR-016.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar016(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Check if it's a make call
		if !utils.IsIdentCall(call, "make") {
			// Not a make call
			return
		}

		// Check if it's make([]T, N) with small constant N
		if shouldUseArray(pass, call) {
			reportArraySuggestion(pass, call)
		}
	})

	// Return analysis result
	return nil, nil
}

// shouldUseArray vérifie si make devrait être remplacé par un array.
//
// Params:
//   - pass: contexte d'analyse
//   - call: expression d'appel à make
//
// Returns:
//   - bool: true si un array est préférable
func shouldUseArray(pass *analysis.Pass, call *ast.CallExpr) bool {
	// Need at least 2 args: make([]T, size)
	if len(call.Args) < minMakeArgsVar016 {
		// Not enough arguments
		return false
	}

	// First arg should be a slice type
	if !utils.IsSliceType(call.Args[0]) {
		// Not a slice type
		return false
	}

	// Check if has different capacity (3rd arg)
	if hasDifferentCapacity(call) {
		// Different capacity, needs slice
		return false
	}

	// Second arg should be small constant
	size := getConstantSize(pass, call.Args[1])
	// Return true if size is small constant
	return isSmallConstant(size)
}

// hasDifferentCapacity vérifie si make a une capacité différente.
//
// Params:
//   - call: expression d'appel à make
//
// Returns:
//   - bool: true si capacité différente spécifiée
func hasDifferentCapacity(call *ast.CallExpr) bool {
	// Return true if 3rd argument exists
	return len(call.Args) >= minMakeArgsWithCapVar016
}

// getConstantSize obtient la taille constante d'une expression.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression de taille
//
// Returns:
//   - int64: taille constante ou -1 si non constante
func getConstantSize(pass *analysis.Pass, expr ast.Expr) int64 {
	// Try to get constant value
	tv := pass.TypesInfo.Types[expr]
	// Check if it's a constant
	if tv.Value == nil {
		// Not a constant
		return -1
	}

	// Get int64 value
	if val, ok := constant.Int64Val(tv.Value); ok {
		// Return the constant value
		return val
	}

	// Not an int constant
	return -1
}

// isSmallConstant vérifie si la taille est petite et constante.
//
// Params:
//   - size: taille à vérifier
//
// Returns:
//   - bool: true si petite constante (<= MAX_ARRAY_SIZE_VAR016)
func isSmallConstant(size int64) bool {
	// Check if it's a positive small constant
	return size > 0 && size <= maxArraySizeVar016
}

// reportArraySuggestion rapporte la suggestion d'utiliser un array.
//
// Params:
//   - pass: contexte d'analyse
//   - call: expression d'appel à make
func reportArraySuggestion(pass *analysis.Pass, call *ast.CallExpr) {
	pass.Reportf(
		call.Pos(),
		"KTN-VAR-016: préférer [N]T (array) au lieu de make([]T, N)",
	)
}
