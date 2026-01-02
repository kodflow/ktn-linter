// Package ktnvar implements KTN linter rules.
package ktnvar

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar026 is the rule code for this analyzer
	ruleCodeVar026 string = "KTN-VAR-026"
)

// Analyzer026 detects manual min/max patterns that can use Go 1.21+ built-ins.
//
// Since Go 1.21, min() and max() are built-in functions.
// This analyzer detects:
// - math.Min and math.Max calls
// - if/else patterns for min/max (if a < b { return a } return b)
var Analyzer026 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar026",
	Doc:      "KTN-VAR-026: Utiliser min()/max() built-in (Go 1.21+)",
	Run:      runVar026,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar026 executes the analysis for KTN-VAR-026.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: always nil
//   - error: any error encountered
func runVar026(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeVar026) {
		// Regle desactivee
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp := inspAny.(*inspector.Inspector)

	// Types de noeuds a analyser
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.IfStmt)(nil),
	}

	// Parcours des noeuds
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar026, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Verification du type de noeud
		switch node := n.(type) {
		// Case CallExpr pour math.Min/math.Max
		case *ast.CallExpr:
			// Verification des appels math.Min/math.Max
			checkMathMinMax(pass, node)
		// Case IfStmt pour pattern if a < b
		case *ast.IfStmt:
			// Verification du pattern if/else min/max
			checkIfMinMaxPattern(pass, node)
		}
	})

	// Traitement
	return nil, nil
}

// checkMathMinMax checks for math.Min or math.Max calls.
//
// Params:
//   - pass: analysis context
//   - call: call expression to check
func checkMathMinMax(pass *analysis.Pass, call *ast.CallExpr) {
	// Verification si c'est un selector (math.Min ou math.Max)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	// Verification de la condition
	if !ok {
		// Pas un selector
		return
	}

	// Recuperation du package
	pkgIdent, ok := sel.X.(*ast.Ident)
	// Verification de la condition
	if !ok {
		// Pas un identifiant
		return
	}

	// Verification si c'est le package math
	if pkgIdent.Name != "math" {
		// Pas le package math
		return
	}

	// Verification si c'est Min ou Max
	if sel.Sel.Name != "Min" && sel.Sel.Name != "Max" {
		// Pas Min ou Max
		return
	}

	// Determiner le nom du built-in equivalent
	builtinName := getBuiltinName(sel.Sel.Name)

	// Report de l'erreur
	msg, _ := messages.Get(ruleCodeVar026)
	pass.Reportf(
		call.Pos(),
		"%s: %s",
		ruleCodeVar026,
		msg.Format(config.Get().Verbose, builtinName, sel.Sel.Name),
	)
}

// getBuiltinName returns the built-in name for a math function.
//
// Params:
//   - mathFuncName: name of the math function (Min or Max)
//
// Returns:
//   - string: the built-in function name (min or max)
func getBuiltinName(mathFuncName string) string {
	// Conversion en minuscule pour le built-in
	if mathFuncName == "Min" {
		// Built-in min
		return "min"
	}
	// Built-in max
	return "max"
}

// checkIfMinMaxPattern checks for if/else min/max patterns.
//
// Params:
//   - pass: analysis context
//   - ifStmt: if statement to check
func checkIfMinMaxPattern(pass *analysis.Pass, ifStmt *ast.IfStmt) {
	// Verification de la condition binaire
	if !isMinMaxCondition(ifStmt.Cond) {
		// Pas une condition min/max
		return
	}

	// Verification du body (doit contenir un return)
	if !hasReturnInBody(ifStmt.Body) {
		// Pas de return dans le body
		return
	}

	// Verification du else ou du statement suivant
	if !hasMatchingReturn(ifStmt) {
		// Pas de return correspondant
		return
	}

	// Report de l'erreur
	msg, _ := messages.Get(ruleCodeVar026)
	pass.Reportf(
		ifStmt.Pos(),
		"%s: %s",
		ruleCodeVar026,
		msg.Format(config.Get().Verbose, "pattern"),
	)
}

// isMinMaxCondition checks if condition is a < b or a > b.
//
// Params:
//   - cond: condition expression to check
//
// Returns:
//   - bool: true if it's a min/max condition
func isMinMaxCondition(cond ast.Expr) bool {
	// Verification si c'est une expression binaire
	binary, ok := cond.(*ast.BinaryExpr)
	// Verification de la condition
	if !ok {
		// Pas une expression binaire
		return false
	}

	// Verification si c'est < ou >
	return binary.Op == token.LSS || binary.Op == token.GTR
}

// hasReturnInBody checks if the body contains a return statement.
//
// Params:
//   - body: block statement to check
//
// Returns:
//   - bool: true if body has a return
func hasReturnInBody(body *ast.BlockStmt) bool {
	// Verification du body
	if body == nil || len(body.List) == 0 {
		// Body vide
		return false
	}

	// Verification si le premier statement est un return
	_, ok := body.List[0].(*ast.ReturnStmt)
	// Retour de la condition
	return ok
}

// hasMatchingReturn checks if there's a matching return after the if.
//
// Params:
//   - ifStmt: if statement to check
//
// Returns:
//   - bool: true if there's a matching return
func hasMatchingReturn(ifStmt *ast.IfStmt) bool {
	// Verification du else
	if ifStmt.Else != nil {
		// Verification si else contient un return
		return hasReturnInElse(ifStmt.Else)
	}

	// Pas de else, on ne peut pas verifier sans contexte parent
	// Pour simplifier, on detecte seulement les patterns explicites
	return false
}

// hasReturnInElse checks if the else branch contains a return.
//
// Params:
//   - elseStmt: else statement to check
//
// Returns:
//   - bool: true if else has a return
func hasReturnInElse(elseStmt ast.Stmt) bool {
	// Verification du type de else
	switch stmt := elseStmt.(type) {
	// Case BlockStmt
	case *ast.BlockStmt:
		// Verification du body
		if stmt == nil || len(stmt.List) == 0 {
			// Body vide
			return false
		}
		// Verification si contient un return
		_, ok := stmt.List[0].(*ast.ReturnStmt)
		// Retour de la condition
		return ok
	// Autres cas
	default:
		// Pas un block
		return false
	}
}
