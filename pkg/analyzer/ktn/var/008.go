package ktnvar

import (
	"go/ast"
	"go/constant"
	"go/types"

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

// isSliceTypeVar008 vérifie si une expression est un type slice.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si c'est un type slice
func isSliceTypeVar008(pass *analysis.Pass, expr ast.Expr) bool {
	// Récupération du type de l'expression
	if tv, ok := pass.TypesInfo.Types[expr]; ok {
		_, isSlice := tv.Type.Underlying().(*types.Slice)
		// Retour du résultat
		return isSlice
	}
	// Vérification via analyse de l'expression
	if arrType, ok := expr.(*ast.ArrayType); ok {
		// Retour si c'est un slice (pas de longueur)
		return arrType.Len == nil
	}
	// Retour par défaut
	return false
}

// hasPositiveLength vérifie si l'argument de longueur est > 0.
//
// Params:
//   - pass: contexte d'analyse
//   - lengthArg: argument de longueur
//
// Returns:
//   - bool: true si length > 0
func hasPositiveLength(pass *analysis.Pass, lengthArg ast.Expr) bool {
	// Récupération de la valeur de l'expression
	if tv, ok := pass.TypesInfo.Types[lengthArg]; ok {
		// Vérification si c'est une constante
		if tv.Value != nil && tv.Value.Kind() == constant.Int {
			val, _ := constant.Int64Val(tv.Value)
			// Retour si valeur > 0
			return val > 0
		}
	}

	// Vérification via analyse de l'expression comme BasicLit
	if lit, ok := lengthArg.(*ast.BasicLit); ok {
		// Retour si la valeur n'est pas "0"
		return lit.Value != "0"
	}

	// Par défaut, considérer comme positif pour détecter les cas suspects
	return true
}

// checkMakeCallVar008 vérifie un appel à make avec length > 0.
//
// Params:
//   - pass: contexte d'analyse
//   - call: appel de fonction à vérifier
func checkMakeCallVar008(pass *analysis.Pass, call *ast.CallExpr) {
	// Vérification que c'est un appel à make
	if ident, ok := call.Fun.(*ast.Ident); !ok || ident.Name != "make" {
		// Continue traversing AST nodes.
		return
	}

	// Vérification du nombre d'arguments (2 ou 3: type, length, [capacity])
	if len(call.Args) < MIN_MAKE_ARGS_VAR008 {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que le type est un slice
	if !isSliceTypeVar008(pass, call.Args[0]) {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que la longueur est > 0
	if !hasPositiveLength(pass, call.Args[1]) {
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
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkMakeCallVar008(pass, call)
	})

	// Retour de la fonction
	return nil, nil
}
