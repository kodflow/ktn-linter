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
	insp := inspAny.(*inspector.Inspector)

	// Filter for ForStmt nodes
	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		forStmt := n.(*ast.ForStmt)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar027, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Check if this is a convertible for loop
		if isConvertibleToRangeInt(forStmt) {
			msg, _ := messages.Get(ruleCodeVar027)
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
	// Extract and validate assignment statement
	assignStmt, ok := extractAssignFromInit027(init)
	// Not a valid short declaration
	if !ok {
		return "", false
	}

	// Extract variable name from LHS
	varName, ok := extractVarNameFromAssign027(assignStmt)
	// LHS not a single identifier
	if !ok {
		return "", false
	}

	// Validate RHS is integer zero
	if !validateInitZero027(assignStmt) {
		// RHS not integer zero
		return "", false
	}

	// Valid init statement
	return varName, true
}

// extractAssignFromInit027 extracts and validates the assignment statement from init.
//
// Params:
//   - init: the init statement to check
//
// Returns:
//   - *ast.AssignStmt: the assignment statement if valid
//   - bool: true if init is a short variable declaration with single variable
func extractAssignFromInit027(init ast.Stmt) (*ast.AssignStmt, bool) {
	// Must be an assignment statement
	assignStmt, ok := init.(*ast.AssignStmt)
	// Not an assignment
	if !ok {
		return nil, false
	}

	// Must be := (define)
	if assignStmt.Tok != token.DEFINE {
		// Not a short declaration
		return nil, false
	}

	// Must have exactly one variable on each side
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		// Multiple variables not supported
		return nil, false
	}

	// Valid assignment statement
	return assignStmt, true
}

// extractVarNameFromAssign027 extracts the variable name from the LHS of an assignment.
//
// Params:
//   - assignStmt: the assignment statement to check
//
// Returns:
//   - string: the variable name if valid
//   - bool: true if LHS is a single identifier
func extractVarNameFromAssign027(assignStmt *ast.AssignStmt) (string, bool) {
	// LHS must be an identifier
	ident, ok := assignStmt.Lhs[0].(*ast.Ident)
	// LHS not an identifier
	if !ok {
		return "", false
	}

	// Valid identifier
	return ident.Name, true
}

// validateInitZero027 validates that the RHS of the assignment is integer zero.
//
// Params:
//   - assignStmt: the assignment statement to check
//
// Returns:
//   - bool: true if RHS is the integer literal 0
func validateInitZero027(assignStmt *ast.AssignStmt) bool {
	// RHS must be a basic literal
	basicLit, ok := assignStmt.Rhs[0].(*ast.BasicLit)
	// RHS not a literal
	if !ok {
		return false
	}

	// Must be integer 0
	return basicLit.Kind == token.INT && basicLit.Value == "0"
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
