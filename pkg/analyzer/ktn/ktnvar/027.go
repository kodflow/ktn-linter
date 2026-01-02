// Package ktnvar provides analyzers for variable-related lint rules.
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
	// ruleCodeVar027 is the rule code for this analyzer
	ruleCodeVar027 string = "KTN-VAR-027"
)

// Analyzer027 checks for classic for loops that can use range over integer (Go 1.22+)
var Analyzer027 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar027",
	Doc:      "KTN-VAR-027: Utiliser 'for i := range n' au lieu de 'for i := 0; i < n; i++'",
	Run:      runVar027,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar027 executes the analysis for KTN-VAR-027.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: potential error
func runVar027(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar027) {
		// Rule disabled
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp, ok := inspAny.(*inspector.Inspector)
	// Defensive: ensure inspector is available
	if !ok || insp == nil {
		return nil, nil
	}

	// Filter for ForStmt nodes
	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		forStmt, ok := n.(*ast.ForStmt)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar027, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Check if this is a convertible for loop
		if isConvertibleToRangeInt(forStmt) {
			msg, ok := messages.Get(ruleCodeVar027)
			// Defensive: avoid panic if message is missing
			if !ok {
				pass.Reportf(forStmt.Pos(), "%s: utiliser 'for i := range n' (Go 1.22+)", ruleCodeVar027)
				return
			}
			pass.Reportf(
				forStmt.Pos(),
				"%s: %s",
				ruleCodeVar027,
				msg.Format(config.Get().Verbose),
			)
		}
	})

	// Return result
	return nil, nil
}

// isConvertibleToRangeInt checks if a for loop can be converted to range over int.
// Pattern: for i := 0; i < n; i++ { ... }
//
// Params:
//   - forStmt: the for statement to check
//
// Returns:
//   - bool: true if convertible to range int
func isConvertibleToRangeInt(forStmt *ast.ForStmt) bool {
	// Must have init, cond, and post
	if forStmt.Init == nil || forStmt.Cond == nil || forStmt.Post == nil {
		// Missing required parts
		return false
	}

	// Check init: must be i := 0
	varName, initOK := checkInitIsZero(forStmt.Init)
	// Init not matching pattern
	if !initOK {
		return false
	}

	// Check post: must be i++ (exactly)
	postOK := checkPostIsIncrement(forStmt.Post, varName)
	// Post not matching pattern
	if !postOK {
		return false
	}

	// Check condition: must be i < n
	condOK := checkConditionIsLessThan(forStmt.Cond, varName)
	// Condition not matching pattern
	if !condOK {
		return false
	}

	// All checks passed
	return true
}

// checkInitIsZero verifies init is "i := 0" and returns the variable name.
//
// Params:
//   - init: the init statement to check
//
// Returns:
//   - string: the variable name if valid
//   - bool: true if init is "varname := 0"
func checkInitIsZero(init ast.Stmt) (string, bool) {
	// Must be an assignment statement
	assignStmt, ok := init.(*ast.AssignStmt)
	// Not an assignment
	if !ok {
		return "", false
	}

	// Must be := (define)
	if assignStmt.Tok != token.DEFINE {
		// Not a short declaration
		return "", false
	}

	// Must have exactly one variable
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		// Multiple variables not supported
		return "", false
	}

	// LHS must be an identifier
	ident, ok := assignStmt.Lhs[0].(*ast.Ident)
	// LHS not an identifier
	if !ok {
		return "", false
	}

	// RHS must be 0 (basic literal)
	basicLit, ok := assignStmt.Rhs[0].(*ast.BasicLit)
	// RHS not a literal
	if !ok {
		return "", false
	}

	// Must be integer 0
	if basicLit.Kind != token.INT || basicLit.Value != "0" {
		// Not integer zero
		return "", false
	}

	// Valid init statement
	return ident.Name, true
}

// checkPostIsIncrement verifies post is "varname++" (exactly).
//
// Params:
//   - post: the post statement to check
//   - varName: the expected variable name
//
// Returns:
//   - bool: true if post is "varname++"
func checkPostIsIncrement(post ast.Stmt, varName string) bool {
	// Must be an increment statement
	incStmt, ok := post.(*ast.IncDecStmt)
	// Not an inc/dec statement
	if !ok {
		return false
	}

	// Must be increment (not decrement)
	if incStmt.Tok != token.INC {
		// Not increment
		return false
	}

	// Must be the same variable
	ident, ok := incStmt.X.(*ast.Ident)
	// X not an identifier
	if !ok {
		return false
	}

	// Check variable name matches
	return ident.Name == varName
}

// checkConditionIsLessThan verifies condition is "varname < expr".
//
// Params:
//   - cond: the condition expression to check
//   - varName: the expected variable name on LHS
//
// Returns:
//   - bool: true if condition is "varname < expr"
func checkConditionIsLessThan(cond ast.Expr, varName string) bool {
	// Must be a binary expression
	binExpr, ok := cond.(*ast.BinaryExpr)
	// Not a binary expression
	if !ok {
		return false
	}

	// Must be less than operator
	if binExpr.Op != token.LSS {
		// Not a less-than comparison
		return false
	}

	// LHS must be the loop variable
	ident, ok := binExpr.X.(*ast.Ident)
	// LHS not an identifier
	if !ok {
		return false
	}

	// Check variable name matches
	return ident.Name == varName
}
