package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MIN_MAKE_ARGS is the minimum number of arguments for make call to check
	MIN_MAKE_ARGS int = 2
)

// Analyzer007 checks that slices are preallocated with capacity when known
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar007",
	Doc:      "KTN-VAR-007: Vérifie que les slices sont préalloués avec une capacité si elle est connue",
	Run:      runVar007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// checkCompositeLit vérifie un composite literal pour les slices vides.
//
// Params:
//   - pass: contexte d'analyse
//   - lit: composite literal à vérifier
func checkCompositeLit(pass *analysis.Pass, lit *ast.CompositeLit) {
	// Vérification que c'est un slice vide
	if len(lit.Elts) > 0 {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que le type est un slice
	if !utils.IsSliceTypeWithPass(pass, lit.Type) {
		// Continue traversing AST nodes.
		return
	}

	// Ignorer si c'est dans un return direct (pas d'append prévu)
	if isInReturnStatement(lit) {
		return
	}

	// Signalement de l'erreur
	pass.Reportf(
		lit.Pos(),
		"KTN-VAR-007: préallouer le slice avec make([]T, 0, capacity) au lieu de []T{}",
	)
}

// checkMakeCall vérifie un appel à make pour les slices sans capacité.
//
// Params:
//   - pass: contexte d'analyse
//   - call: appel de fonction à vérifier
func checkMakeCall(pass *analysis.Pass, call *ast.CallExpr) {
	// Vérification que c'est un appel à make
	if !utils.IsMakeCall(call) {
		// Continue traversing AST nodes.
		return
	}

	// Vérification du nombre d'arguments (doit être 2: type et length)
	if len(call.Args) != MIN_MAKE_ARGS {
		// Continue traversing AST nodes.
		return
	}

	// Vérification que le type est un slice
	if !utils.IsSliceTypeWithPass(pass, call.Args[0]) {
		// Continue traversing AST nodes.
		return
	}

	// Signalement de l'erreur
	pass.Reportf(
		call.Pos(),
		"KTN-VAR-007: spécifier une capacité avec make([]T, 0, capacity) au lieu de make([]T, 0)",
	)
}

// isInReturnStatement checks if node is in a return statement.
//
// Params:
//   - lit: Composite literal to check
func isInReturnStatement(lit *ast.CompositeLit) bool {
	// Cette fonction devrait checker si le parent est un ReturnStmt
	// Mais on n'a pas accès au parent avec inspector
	// Solution : on ignore les slices vides car c'est un cas trop courant
	// et make([]T, 0, capacity) sans append est une optimisation prématurée
	return true // Toujours ignorer pour éviter faux positifs
}

// runVar007 exécute l'analyse KTN-VAR-007.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar007(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	// DÉSACTIVATION de la vérification des []T{} car trop de faux positifs
	// Ne vérifie que make([]T, 0) sans capacity
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Traitement des appels de fonction pour détecter make([]T, 0)
		if node, ok := n.(*ast.CallExpr); ok {
			checkMakeCall(pass, node)
		}
	})

	// Retour de la fonction
	return nil, nil
}
