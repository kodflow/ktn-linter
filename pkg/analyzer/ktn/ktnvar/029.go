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
	// ruleCodeVar029 is the rule code for this analyzer.
	ruleCodeVar029 string = "KTN-VAR-029"
)

// Analyzer029 detects manual slice grow patterns that can use slices.Grow (Go 1.21+).
//
// Detects:
//   - if cap(s)-len(s) < n { newSlice := make([]T, ...); copy(newSlice, s); s = newSlice }
var Analyzer029 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar029",
	Doc:      "KTN-VAR-029: Utiliser slices.Grow() au lieu du pattern manuel de grow (Go 1.21+)",
	Run:      runVar029,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// growPatternInfo stores information about a potential grow pattern.
type growPatternInfo struct {
	// sliceName is the name of the slice being grown
	sliceName string
	// hasMake indicates if make was found
	hasMake bool
	// hasCopy indicates if copy was found
	hasCopy bool
	// hasReassign indicates if reassignment to original slice was found
	hasReassign bool
	// pos is the position of the if statement
	pos token.Pos
}

// runVar029 runs the KTN-VAR-029 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: any error that occurred
func runVar029(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar029) {
		// Rule disabled
		return nil, nil
	}

	// Get AST inspector
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Node types to analyze
	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}

	// Traverse all if statements
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar029, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Cast to if statement
		ifStmt, ok := n.(*ast.IfStmt)
		// Check if valid
		if !ok {
			// Not an if statement
			return
		}

		// Check if this is a grow pattern
		checkGrowPattern(pass, ifStmt, cfg)
	})

	// Return analysis result
	return nil, nil
}

// checkGrowPattern checks if an if statement is a manual grow pattern.
//
// Params:
//   - pass: analysis context
//   - ifStmt: if statement to check
//   - cfg: configuration
func checkGrowPattern(pass *analysis.Pass, ifStmt *ast.IfStmt, cfg *config.Config) {
	// Check if condition matches cap(s)-len(s) < n pattern
	sliceName := extractSliceFromCapLenCondition(ifStmt.Cond)
	// Verify slice name was found
	if sliceName == "" {
		// Not the pattern we're looking for
		return
	}

	// Check if body contains make + copy + reassign pattern
	if hasGrowPatternInBody(ifStmt.Body, sliceName) {
		// Report the issue
		reportGrowPattern(pass, ifStmt, cfg)
	}
}

// extractSliceFromCapLenCondition extracts slice name from cap(s)-len(s) < n condition.
//
// Params:
//   - cond: condition expression to check
//
// Returns:
//   - string: slice name or empty if not matching
func extractSliceFromCapLenCondition(cond ast.Expr) string {
	// Check if it's a binary expression
	binExpr, ok := cond.(*ast.BinaryExpr)
	// Verify it's a comparison
	if !ok {
		// Not a binary expression
		return ""
	}

	// Check if it's a < comparison
	if binExpr.Op != token.LSS {
		// Not a less-than comparison
		return ""
	}

	// Left side should be cap(s)-len(s)
	return extractSliceFromCapMinusLen(binExpr.X)
}

// extractSliceFromCapMinusLen extracts slice name from cap(s)-len(s) expression.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - string: slice name or empty if not matching
func extractSliceFromCapMinusLen(expr ast.Expr) string {
	// Check if it's a binary expression with subtraction
	binExpr, ok := expr.(*ast.BinaryExpr)
	// Verify it's a subtraction
	if !ok || binExpr.Op != token.SUB {
		// Not a subtraction expression
		return ""
	}

	// Extract cap(s) and len(s)
	capSlice := extractSliceFromCapCall(binExpr.X)
	lenSlice := extractSliceFromLenCall(binExpr.Y)

	// Check if both match the same slice
	if capSlice == "" || lenSlice == "" || capSlice != lenSlice {
		// Slices don't match
		return ""
	}

	// Return the slice name
	return capSlice
}

// extractSliceFromCapCall extracts slice name from cap(s) call.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - string: slice name or empty if not a cap call
func extractSliceFromCapCall(expr ast.Expr) string {
	// Check if it's a call expression
	call, ok := expr.(*ast.CallExpr)
	// Verify it's a call
	if !ok {
		// Not a call expression
		return ""
	}

	// Check if it's a cap call
	ident, ok := call.Fun.(*ast.Ident)
	// Verify it's cap
	if !ok || ident.Name != "cap" {
		// Not a cap call
		return ""
	}

	// Check that there's exactly one argument
	if len(call.Args) != 1 {
		// Wrong number of arguments
		return ""
	}

	// Extract slice name from argument
	return extractSliceName(call.Args[0])
}

// extractSliceFromLenCall extracts slice name from len(s) call.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - string: slice name or empty if not a len call
func extractSliceFromLenCall(expr ast.Expr) string {
	// Check if it's a call expression
	call, ok := expr.(*ast.CallExpr)
	// Verify it's a call
	if !ok {
		// Not a call expression
		return ""
	}

	// Check if it's a len call
	ident, ok := call.Fun.(*ast.Ident)
	// Verify it's len
	if !ok || ident.Name != "len" {
		// Not a len call
		return ""
	}

	// Check that there's exactly one argument
	if len(call.Args) != 1 {
		// Wrong number of arguments
		return ""
	}

	// Extract slice name from argument
	return extractSliceName(call.Args[0])
}

// extractSliceName extracts the slice name from an expression.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - string: slice name or empty if not an identifier
func extractSliceName(expr ast.Expr) string {
	// Check if it's an identifier
	ident, ok := expr.(*ast.Ident)
	// Verify it's an identifier
	if !ok {
		// Not an identifier
		return ""
	}

	// Return the identifier name
	return ident.Name
}

// hasGrowPatternInBody checks if a block contains make+copy+reassign pattern.
//
// Params:
//   - body: block statement to check
//   - sliceName: name of the slice being grown
//
// Returns:
//   - bool: true if pattern found
func hasGrowPatternInBody(body *ast.BlockStmt, sliceName string) bool {
	// Defensive: malformed ASTs may have a nil body
	if body == nil {
		return false
	}

	// Track pattern components
	info := &growPatternInfo{
		sliceName: sliceName,
	}

	// Iterate through statements in the body
	for _, stmt := range body.List {
		// Check for different statement types
		analyzeGrowStatement(stmt, info)
	}

	// Check if all components are present
	return info.hasMake && info.hasCopy && info.hasReassign
}

// analyzeGrowStatement analyzes a statement for grow pattern components.
//
// Params:
//   - stmt: statement to analyze
//   - info: pattern info to update
func analyzeGrowStatement(stmt ast.Stmt, info *growPatternInfo) {
	// Check different statement types
	switch s := stmt.(type) {
	// Assignment or declaration statement
	case *ast.AssignStmt:
		// Check for make or reassignment
		analyzeAssignForGrow(s, info)
	// Expression statement (copy call)
	case *ast.ExprStmt:
		// Check for copy call
		analyzeCopyForGrow(s, info)
	}
}

// analyzeAssignForGrow analyzes an assignment for make or reassignment.
//
// Params:
//   - assignStmt: assignment to analyze
//   - info: pattern info to update
func analyzeAssignForGrow(assignStmt *ast.AssignStmt, info *growPatternInfo) {
	// Check for single assignment
	if len(assignStmt.Rhs) != 1 {
		// Multiple assignments
		return
	}

	// Check if RHS is a make call
	if isMakeSliceCall(assignStmt.Rhs[0]) {
		// Mark make as found
		info.hasMake = true
		// Return early
		return
	}

	// Check for reassignment to original slice
	checkSliceReassign(assignStmt, info)
}

// isMakeSliceCall checks if an expression is a make([]T, ...) call.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - bool: true if it's a make call for a slice
func isMakeSliceCall(expr ast.Expr) bool {
	// Check if it's a call expression
	call, ok := expr.(*ast.CallExpr)
	// Verify it's a call
	if !ok {
		// Not a call expression
		return false
	}

	// Check if it's a make call
	ident, ok := call.Fun.(*ast.Ident)
	// Verify it's make
	if !ok || ident.Name != "make" {
		// Not a make call
		return false
	}

	// Check that first arg is a slice type
	if len(call.Args) < 1 {
		// No arguments
		return false
	}

	// Check if first arg is a slice type
	_, ok = call.Args[0].(*ast.ArrayType)
	// Return result
	return ok
}

// checkSliceReassign checks if an assignment reassigns to the original slice.
//
// Params:
//   - assignStmt: assignment to check
//   - info: pattern info to update
func checkSliceReassign(assignStmt *ast.AssignStmt, info *growPatternInfo) {
	// Check for single LHS
	if len(assignStmt.Lhs) != 1 {
		// Multiple targets
		return
	}

	// Get the target identifier
	targetIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	// Verify it's an identifier
	if !ok {
		// Not an identifier
		return
	}

	// Check if it matches the slice name
	if targetIdent.Name == info.sliceName {
		// Check that it's a simple assignment (=), not a declaration (:=)
		if assignStmt.Tok == token.ASSIGN {
			// Mark reassignment as found
			info.hasReassign = true
		}
	}
}

// analyzeCopyForGrow analyzes an expression statement for a copy call.
//
// Params:
//   - exprStmt: expression statement to analyze
//   - info: pattern info to update
func analyzeCopyForGrow(exprStmt *ast.ExprStmt, info *growPatternInfo) {
	// Check if it's a call expression
	call, ok := exprStmt.X.(*ast.CallExpr)
	// Verify it's a call
	if !ok {
		// Not a call expression
		return
	}

	// Check if it's a copy call
	ident, ok := call.Fun.(*ast.Ident)
	// Verify it's copy
	if !ok || ident.Name != "copy" {
		// Not a copy call
		return
	}

	// Check that there are 2 arguments
	if len(call.Args) != 2 {
		// Wrong number of arguments
		return
	}

	// Check if source is the original slice
	sourceIdent, ok := call.Args[1].(*ast.Ident)
	// Verify source matches
	if ok && sourceIdent.Name == info.sliceName {
		// Mark copy as found
		info.hasCopy = true
	}
}

// reportGrowPattern reports a manual grow pattern.
//
// Params:
//   - pass: analysis context
//   - ifStmt: if statement containing the pattern
//   - cfg: configuration
func reportGrowPattern(pass *analysis.Pass, ifStmt *ast.IfStmt, cfg *config.Config) {
	// Get the message
	msg, ok := messages.Get(ruleCodeVar029)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(ifStmt.Pos(), "%s: utiliser slices.Grow() au lieu du pattern manuel",
			ruleCodeVar029)
		return
	}
	// Report the issue
	pass.Reportf(
		ifStmt.Pos(),
		"%s: %s",
		ruleCodeVar029,
		msg.Format(cfg.Verbose),
	)
}
