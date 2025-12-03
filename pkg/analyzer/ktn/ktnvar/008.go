// Analyzer 008 for the ktnvar package.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MIN_MAKE_ARGS_VAR008 is the minimum number of arguments for make call
	MIN_MAKE_ARGS_VAR008 int = 2
)

// Analyzer008 checks that make with length > 0 is avoided when append is used
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar008",
	Doc:      "KTN-VAR-008: Vérifie d'éviter make([]T, length) si utilisation avec append",
	Run:      runVar008,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// checkMakeCallVar008 vérifie un appel à make avec length > 0.
//
// Params:
//   - pass: contexte d'analyse
//   - call: appel de fonction à vérifier
func checkMakeCallVar008(pass *analysis.Pass, call *ast.CallExpr) {
	// Vérification que c'est un appel à make
	if !utils.IsMakeCall(call) {
		// Continue traversing AST nodes.
		return
	}

	// Vérification du nombre d'arguments (2 ou 3: type, length, [capacity])
	if len(call.Args) < MIN_MAKE_ARGS_VAR008 {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que le type est un slice
	if !utils.IsSliceTypeWithPass(pass, call.Args[0]) {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que la longueur est > 0
	if !utils.HasPositiveLength(pass, call.Args[1]) {
		// Continue traversing AST nodes.
		return
	}

	// Signalement de l'erreur
	pass.Reportf(
		call.Pos(),
		"KTN-VAR-008: utiliser make([]T, 0, capacity) au lieu de make([]T, length) pour éviter les zéro-values inutiles avant append",
	)
}

// runVar008 exécute l'analyse KTN-VAR-008.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar008(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkMakeCallVar008(pass, call)
	})

	// Retour de la fonction
	return nil, nil
}
