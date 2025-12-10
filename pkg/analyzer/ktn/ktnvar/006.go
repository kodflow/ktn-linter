// Analyzer 006 for the ktnvar package.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar006 is the rule code for this analyzer
	ruleCodeVar006 string = "KTN-VAR-006"
)

// Analyzer006 checks for strings.Builder/bytes.Buffer without Grow preallocate
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar006",
	Doc:      "KTN-VAR-006: Préallouer bytes.Buffer/strings.Builder avec Grow",
	Run:      runVar006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar006 exécute l'analyse KTN-VAR-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar006(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar006) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ValueSpec)(nil),
		(*ast.AssignStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar006, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		checkBuilderWithoutGrow(pass, n)
	})

	// Retour de la fonction
	return nil, nil
}

// checkBuilderWithoutGrow checks if Builder/Buffer lacks Grow call.
//
// Params:
//   - pass: analysis pass context
//   - n: AST node to check
func checkBuilderWithoutGrow(pass *analysis.Pass, n ast.Node) {
	// Composite literal found
	var compositeLit *ast.CompositeLit
	// Position for reporting
	var pos ast.Node

	// Check type of node
	switch node := n.(type) {
	// Case: variable specification
	case *ast.ValueSpec:
		compositeLit, pos = checkValueSpec(node)
	// Case: assignment statement
	case *ast.AssignStmt:
		compositeLit, pos = checkAssignStmt(node)
	}

	// Report if composite literal found
	if compositeLit != nil && pos != nil {
		reportMissingGrow(pass, pos)
	}
}

// checkValueSpec checks var declaration for Builder/Buffer.
//
// Params:
//   - node: variable specification node
//
// Returns:
//   - *ast.CompositeLit: composite literal if found
//   - ast.Node: position node for reporting
func checkValueSpec(node *ast.ValueSpec) (*ast.CompositeLit, ast.Node) {
	// Only check var sb = strings.Builder{}
	if len(node.Values) == 0 {
		// Return nil if no values
		return nil, nil
	}

	// Check var sb = strings.Builder{}
	for _, val := range node.Values {
		// Check if value is composite literal
		if lit, ok := val.(*ast.CompositeLit); ok {
			// Check if it's Builder/Buffer type
			if isBuilderCompositeLit(lit) {
				// Return found composite literal
				return lit, node
			}
		}
	}

	// Return nil if not found
	return nil, nil
}

// checkAssignStmt checks assignment for Builder/Buffer.
//
// Params:
//   - node: assignment statement node
//
// Returns:
//   - *ast.CompositeLit: composite literal if found
//   - ast.Node: position node for reporting
func checkAssignStmt(node *ast.AssignStmt) (*ast.CompositeLit, ast.Node) {
	// Check short declaration (sb := strings.Builder{})
	for _, rhs := range node.Rhs {
		// Check if right-hand side is composite literal
		if lit, ok := rhs.(*ast.CompositeLit); ok {
			// Check if it's Builder/Buffer type
			if isBuilderCompositeLit(lit) {
				// Return found composite literal
				return lit, node
			}
		}
	}

	// Return nil if not found
	return nil, nil
}

// isBuilderCompositeLit checks if composite is Builder/Buffer.
//
// Params:
//   - lit: composite literal to check
//
// Returns:
//   - bool: true if Builder/Buffer composite literal
func isBuilderCompositeLit(lit *ast.CompositeLit) bool {
	// Check if type is selector expression
	if sel, ok := lit.Type.(*ast.SelectorExpr); ok {
		var pkg *ast.Ident
		// Check package and type names
		if pkg, ok = sel.X.(*ast.Ident); ok {
			pkgName := pkg.Name
			typeName := sel.Sel.Name
			// Check for strings.Builder or bytes.Buffer
			return (pkgName == "strings" && typeName == "Builder") ||
				(pkgName == "bytes" && typeName == "Buffer")
		}
	}
	// Return false if not Builder/Buffer
	return false
}

// reportMissingGrow reports missing Grow() call.
//
// Params:
//   - pass: analysis pass context
//   - node: AST node position for reporting
func reportMissingGrow(pass *analysis.Pass, node ast.Node) {
	// Get type information
	var typeStr string
	// Check type of node
	switch n := node.(type) {
	// Case: variable specification
	case *ast.ValueSpec:
		typeStr = extractTypeString(n.Type, n.Values)
	// Case: assignment statement
	case *ast.AssignStmt:
		typeStr = extractAssignTypeString(n)
	}

	// Check if type string is valid
	if typeStr != "" {
		pass.Reportf(
			node.Pos(),
			"KTN-VAR-006: préallouer %s avec Grow() avant boucle pour optimiser les allocations",
			typeStr,
		)
	}
}

// extractTypeString extracts type name from ValueSpec.
//
// Params:
//   - typeExpr: type expression
//   - values: initialization values
//
// Returns:
//   - string: type name as string
func extractTypeString(typeExpr ast.Expr, values []ast.Expr) string {
	// Check if type expression exists
	if typeExpr != nil {
		// Check if selector expression
		if sel, ok := typeExpr.(*ast.SelectorExpr); ok {
			var pkg *ast.Ident
			// Check package name
			if pkg, ok = sel.X.(*ast.Ident); ok {
				// Return package.Type format
				return pkg.Name + "." + sel.Sel.Name
			}
		}
	}

	// Check values for composite literal
	if len(values) > 0 {
		// Check first value
		if lit, ok := values[0].(*ast.CompositeLit); ok {
			var sel *ast.SelectorExpr
			// Check if type is selector expression
			if sel, ok = lit.Type.(*ast.SelectorExpr); ok {
				var pkg *ast.Ident
				// Check package name
				if pkg, ok = sel.X.(*ast.Ident); ok {
					// Return package.Type format
					return pkg.Name + "." + sel.Sel.Name
				}
			}
		}
	}

	// Return empty if type cannot be extracted
	return ""
}

// extractAssignTypeString extracts type from assignment.
//
// Params:
//   - assign: assignment statement
//
// Returns:
//   - string: type name as string
func extractAssignTypeString(assign *ast.AssignStmt) string {
	// Iteration over right-hand side expressions
	for _, rhs := range assign.Rhs {
		// Check if composite literal
		if lit, ok := rhs.(*ast.CompositeLit); ok {
			var sel *ast.SelectorExpr
			// Check if type is selector expression
			if sel, ok = lit.Type.(*ast.SelectorExpr); ok {
				var pkg *ast.Ident
				// Check package name
				if pkg, ok = sel.X.(*ast.Ident); ok {
					// Return package.Type format
					return pkg.Name + "." + sel.Sel.Name
				}
			}
		}
	}

	// Return empty if type cannot be extracted
	return ""
}
