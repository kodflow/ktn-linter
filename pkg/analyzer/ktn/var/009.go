package ktnvar

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer009 detects map allocations without capacity hints
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar009",
	Doc:      "KTN-VAR-009: Préallouer maps avec capacité si connue",
	Run:      runVar009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar009 exécute l'analyse KTN-VAR-009.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar009(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Vérification que c'est un appel à "make"
		if !isMakeCall(callExpr) {
			// Continue traversing AST nodes
			return
		}

		// Vérification que le type est une map
		if !isMapType(pass, callExpr) {
			// Continue traversing AST nodes
			return
		}

		// Vérification que make a exactement 1 argument (type seulement)
		if len(callExpr.Args) != 1 {
			// make() avec capacité fournie, conforme
			return
		}

		// Signaler l'erreur
		pass.Reportf(
			callExpr.Pos(),
			"KTN-VAR-009: préallouer la map avec une capacité (make(map[K]V, capacity))",
		)
	})

	// Retour de la fonction
	return nil, nil
}

// isMakeCall vérifie si l'expression est un appel à make.
//
// Params:
//   - callExpr: expression d'appel à vérifier
//
// Returns:
//   - bool: true si c'est un appel à make
func isMakeCall(callExpr *ast.CallExpr) bool {
	ident, ok := callExpr.Fun.(*ast.Ident)
	// Vérification que c'est un identifiant et que c'est "make"
	if !ok || ident.Name != "make" {
		// Retour false si ce n'est pas make
		return false
	}
	// Retour true si c'est make
	return true
}

// isMapType vérifie si le type passé à make est une map.
//
// Params:
//   - pass: contexte d'analyse
//   - callExpr: expression d'appel à vérifier
//
// Returns:
//   - bool: true si le type est une map
func isMapType(pass *analysis.Pass, callExpr *ast.CallExpr) bool {
	// Vérification qu'il y a au moins un argument
	if len(callExpr.Args) == 0 {
		// Retour false si pas d'arguments
		return false
	}

	// Récupération du type depuis TypesInfo
	tv, ok := pass.TypesInfo.Types[callExpr.Args[0]]
	// Vérification de la disponibilité du type
	if !ok {
		// Retour false si type non disponible
		return false
	}

	// Vérification que c'est un type map
	_, isMap := tv.Type.(*types.Map)
	// Retour du résultat
	return isMap
}
