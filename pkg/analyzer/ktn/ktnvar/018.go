// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/constant"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar018 is the rule code for this analyzer
	ruleCodeVar018 string = "KTN-VAR-018"
	// minMakeArgsVar016 is minimum arguments for make([]T, N)
	minMakeArgsVar016 int = 2
	// minMakeArgsWithCapVar016 is minimum arguments for make with capacity
	minMakeArgsWithCapVar016 int = 3
	// maxArraySizeBytes is maximum total bytes for recommending array over slice
	maxArraySizeBytes int64 = 64
)

// Analyzer018 checks for make([]T, N) with small constant N
var Analyzer018 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar018",
	Doc:      "KTN-VAR-018: Vérifie l'utilisation de [N]T au lieu de make([]T, N)",
	Run:      runVar018,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar018 exécute l'analyse KTN-VAR-018.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar018(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar018) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp, ok := inspAny.(*inspector.Inspector)
	// Defensive: ensure inspector is available
	if !ok || insp == nil {
		return nil, nil
	}
	// Defensive: avoid nil dereference when resolving positions
	if pass.Fset == nil {
		return nil, nil
	}
	// Defensive: avoid nil dereference when resolving types
	if pass.TypesInfo == nil {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar018, pass.Fset.Position(call.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Check if it's a make call
		if !utils.IsIdentCall(call, "make") {
			// Not a make call
			return
		}

		// Check if it's make([]T, N) with small constant N and ≤64 bytes
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

	// Second arg should be constant
	size := getConstantSize(pass, call.Args[1])
	// Not a constant size
	if size <= 0 {
		// Dynamic or invalid size, can't use array
		return false
	}

	// Check total size in bytes
	return isTotalSizeSmall(pass, call.Args[0], size)
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

// isTotalSizeSmall vérifie si la taille totale en bytes est petite.
//
// Params:
//   - pass: contexte d'analyse
//   - sliceType: type de slice
//   - elementCount: nombre d'éléments
//
// Returns:
//   - bool: true si taille totale <= 64 bytes
func isTotalSizeSmall(pass *analysis.Pass, sliceType ast.Expr, elementCount int64) bool {
	// Get element type from []T
	arrayType, ok := sliceType.(*ast.ArrayType)
	// Not a valid array type
	if !ok {
		// Invalid slice type expression
		return false
	}

	// Get element type
	elemType := pass.TypesInfo.TypeOf(arrayType.Elt)
	// Could not determine element type
	if elemType == nil {
		// Type information not available
		return false
	}

	// Get size of element type in bytes using Sizes API
	sizes := pass.TypesSizes
	// Use default sizes for amd64 if not available
	if sizes == nil {
		sizes = types.SizesFor("gc", "amd64")
	}
	// Still no sizes information available
	if sizes == nil {
		// Cannot determine type sizes
		return false
	}

	// Get element size
	elemSize := sizes.Sizeof(elemType)

	// Calculate total size
	totalSize := elemSize * elementCount

	// Check if total size <= 64 bytes
	return totalSize <= maxArraySizeBytes
}

// reportArraySuggestion rapporte la suggestion d'utiliser un array.
//
// Params:
//   - pass: contexte d'analyse
//   - call: expression d'appel à make
func reportArraySuggestion(pass *analysis.Pass, call *ast.CallExpr) {
	msg, ok := messages.Get(ruleCodeVar018)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(call.Pos(), "%s: utiliser [N]T au lieu de make([]T, N) pour ≤64 bytes", ruleCodeVar018)
		return
	}
	pass.Reportf(
		call.Pos(),
		"%s: %s",
		ruleCodeVar018,
		msg.Format(config.Get().Verbose),
	)
}
