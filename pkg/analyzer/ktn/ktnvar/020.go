// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar020 is the rule code for this analyzer
	ruleCodeVar020 string = "KTN-VAR-020"
)

// Analyzer020 detects empty slice literals and make([]T, 0) without capacity.
//
// Nil slices are preferred over empty slice literals because they are
// functionally equivalent but more efficient (no allocation).
var Analyzer020 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar020",
	Doc:      "KTN-VAR-020: Preferer nil slice a []T{} ou make([]T, 0)",
	Run:      runVar020,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar020 executes the analysis for KTN-VAR-020.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: always nil
//   - error: any error encountered
func runVar020(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeVar020) {
		// Regle desactivee
		return nil, nil
	}

	// Recuperation de l'inspecteur AST
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Types de noeuds a analyser
	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.CallExpr)(nil),
	}

	// Parcours des noeuds
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar020, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Verification du type de noeud
		switch node := n.(type) {
		// Traitement des composite literals
		case *ast.CompositeLit:
			// Verification des slice literals vides
			checkEmptySliceLiteral(pass, node)
		// Traitement des appels de fonction
		case *ast.CallExpr:
			// Verification des make([]T, 0) sans capacite
			checkMakeSliceZero(pass, node)
		}
	})

	// Traitement
	return nil, nil
}

// checkEmptySliceLiteral checks for empty slice literals []T{}.
//
// Params:
//   - pass: analysis context
//   - lit: composite literal to check
func checkEmptySliceLiteral(pass *analysis.Pass, lit *ast.CompositeLit) {
	// Verifier que c'est un slice type
	if !isSliceType(lit.Type) {
		// Pas un slice
		return
	}

	// Verifier que le literal est vide
	if len(lit.Elts) > 0 {
		// Slice avec elements
		return
	}

	// Recuperation du type pour le message
	typeStr := formatSliceType(lit.Type)

	// Report de l'erreur
	msg, ok := messages.Get(ruleCodeVar020)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(lit.Pos(), "%s: préférer nil slice à %s{}", ruleCodeVar020, typeStr)
		return
	}
	pass.Reportf(
		lit.Pos(),
		"%s: %s",
		ruleCodeVar020,
		msg.Format(config.Get().Verbose, typeStr),
	)
}

// checkMakeSliceZero checks for make([]T, 0) without capacity.
//
// Params:
//   - pass: analysis context
//   - call: call expression to check
func checkMakeSliceZero(pass *analysis.Pass, call *ast.CallExpr) {
	// Verifier que c'est un appel a make
	if !isMakeCall(call) {
		// Pas un appel a make
		return
	}

	// Verifier qu'il y a au moins 2 arguments (type et length)
	if len(call.Args) < 2 {
		// Pas assez d'arguments
		return
	}

	// Verifier que le premier argument est un slice type
	if !isSliceType(call.Args[0]) {
		// Pas un slice
		return
	}

	// Verifier que la longueur est 0
	if !isZeroLiteral(call.Args[1]) {
		// Longueur non nulle
		return
	}

	// Verifier qu'il n'y a pas de capacite specifiee
	if len(call.Args) >= 3 {
		// Capacite specifiee, OK
		return
	}

	// Recuperation du type pour le message
	typeStr := "make(" + formatSliceType(call.Args[0]) + ", 0)"

	// Report de l'erreur
	msg, ok := messages.Get(ruleCodeVar020)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(call.Pos(), "%s: préférer nil slice à %s", ruleCodeVar020, typeStr)
		return
	}
	pass.Reportf(
		call.Pos(),
		"%s: %s",
		ruleCodeVar020,
		msg.Format(config.Get().Verbose, typeStr),
	)
}

// isSliceType checks if an expression is a slice type.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - bool: true if it's a slice type
func isSliceType(expr ast.Expr) bool {
	// Verification du type de noeud
	_, ok := expr.(*ast.ArrayType)
	// Retour du resultat
	return ok
}

// isMakeCall checks if a call expression is a call to make.
//
// Params:
//   - call: call expression to check
//
// Returns:
//   - bool: true if it's a call to make
func isMakeCall(call *ast.CallExpr) bool {
	// Verification que la fonction est un identifiant
	ident, ok := call.Fun.(*ast.Ident)
	// Verification de la condition
	if !ok {
		// Pas un identifiant
		return false
	}

	// Verification que c'est make
	return ident.Name == "make"
}

// isZeroLiteral checks if an expression is a zero literal.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - bool: true if it's a zero literal
func isZeroLiteral(expr ast.Expr) bool {
	// Verification du type de noeud
	lit, ok := expr.(*ast.BasicLit)
	// Verification de la condition
	if !ok {
		// Pas un literal basique
		return false
	}

	// Verification de la valeur
	return lit.Value == "0"
}

// formatSliceType formats a slice type expression as a string.
//
// Params:
//   - expr: slice type expression
//
// Returns:
//   - string: formatted type string
func formatSliceType(expr ast.Expr) string {
	// Verification du type de noeud
	arrType, ok := expr.(*ast.ArrayType)
	// Verification de la condition
	if !ok {
		// Pas un array type
		return "[]T"
	}

	// Formatage du type d'element
	return "[]" + formatElementType(arrType.Elt)
}

// formatElementType formats an element type expression as a string.
//
// Params:
//   - expr: element type expression
//
// Returns:
//   - string: formatted type string
func formatElementType(expr ast.Expr) string {
	// Verification du type de noeud
	switch e := expr.(type) {
	// Identifiant simple
	case *ast.Ident:
		// Retour du nom
		return e.Name
	// Selector expression (pkg.Type)
	case *ast.SelectorExpr:
		// Formatage avec package
		return formatSelectorType(e)
	// Pointer type
	case *ast.StarExpr:
		// Formatage avec asterisque
		return "*" + formatElementType(e.X)
	// Array type
	case *ast.ArrayType:
		// Formatage du slice
		return formatSliceType(e)
	// Map type
	case *ast.MapType:
		// Formatage de la map
		return formatMapType(e)
	// Interface type
	case *ast.InterfaceType:
		// Interface vide
		return "interface{}"
	// Default
	default:
		// Type inconnu
		return "T"
	}
}

// formatSelectorType formats a selector expression as a string.
//
// Params:
//   - sel: selector expression
//
// Returns:
//   - string: formatted type string
func formatSelectorType(sel *ast.SelectorExpr) string {
	// Verification du type X
	if ident, ok := sel.X.(*ast.Ident); ok {
		// Retour du type qualifie
		return ident.Name + "." + sel.Sel.Name
	}

	// Default
	return sel.Sel.Name
}

// formatMapType formats a map type expression as a string.
//
// Params:
//   - m: map type expression
//
// Returns:
//   - string: formatted type string
func formatMapType(m *ast.MapType) string {
	// Formatage de la map
	return "map[" + formatElementType(m.Key) + "]" + formatElementType(m.Value)
}
