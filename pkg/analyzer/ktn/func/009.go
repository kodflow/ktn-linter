package ktnfunc

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer009 checks that getter functions don't have side effects
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc009",
	Doc:      "KTN-FUNC-009: Les getters (Get*/Is*/Has*) ne doivent pas avoir de side effects (assignations, appels de fonctions modifiant l'état)",
	Run:      runFunc009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc009 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc009(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		funcName := funcDecl.Name.Name

		// Skip test functions
		if isTestFunction(funcName) {
			// Retour de la fonction
			return
		}

		// Skip if not a getter (Get*, Is*, Has*)
		if !isGetter(funcName) {
			// Retour de la fonction
			return
		}

		// Skip if no body (external functions)
		if funcDecl.Body == nil {
			// Retour de la fonction
			return
		}

		// Check for side effects
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			// Sélection selon la valeur
			switch stmt := node.(type) {
			// Traitement
			case *ast.AssignStmt:
				// Check if it's assigning to a field or external variable
				// Assignments to local variables (created in the function) are OK
				for _, lhs := range stmt.Lhs {
					// Vérification de la condition
					if hasSideEffect(lhs) {
						pass.Reportf(
							stmt.Pos(),
							"KTN-FUNC-009: le getter '%s' ne doit pas modifier l'état (assignation détectée)",
							funcName,
						)
					}
				}
			// Traitement
			case *ast.IncDecStmt:
				// ++ or -- on fields
				if hasSideEffect(stmt.X) {
					pass.Reportf(
						stmt.Pos(),
						"KTN-FUNC-009: le getter '%s' ne doit pas modifier l'état (incrémentation/décrémentation détectée)",
						funcName,
					)
				}
			}
			// Retour de la fonction
			return true
		})
	})

	// Retour de la fonction
	return nil, nil
}

// isGetter checks if a function name suggests it's a getter
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si fonction getter
func isGetter(name string) bool {
	// Retour de la fonction
	return strings.HasPrefix(name, "Get") ||
		strings.HasPrefix(name, "Is") ||
		strings.HasPrefix(name, "Has")
}

// hasSideEffect checks if an expression modifies external state
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si effet de bord détecté
func hasSideEffect(expr ast.Expr) bool {
	// Sélection selon la valeur
	switch e := expr.(type) {
	// Traitement
	case *ast.SelectorExpr:
		// Modifying a field is a side effect
		return true
	// Traitement
	case *ast.IndexExpr:
		// Modifying an index (array/map/slice element) could be a side effect
		// Check if the base is a selector
		if _, ok := e.X.(*ast.SelectorExpr); ok {
			// Retour de la fonction
			return true
		}
	}
	// Retour de la fonction
	return false
}
