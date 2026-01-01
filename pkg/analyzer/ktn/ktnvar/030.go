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
	// ruleCodeVar030 is the rule code for this analyzer.
	ruleCodeVar030 string = "KTN-VAR-030"
)

// Analyzer030 detects slice cloning patterns that can use slices.Clone (Go 1.21+).
//
// Detects:
//   - make([]T, len(s)) followed by copy(clone, s)
//   - append([]T(nil), s...)
var Analyzer030 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar030",
	Doc:      "KTN-VAR-030: Utiliser slices.Clone() au lieu des patterns de copie manuels (Go 1.21+)",
	Run:      runVar030,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// makeCloneInfo stores information about a make call that could be a clone.
type makeCloneInfo struct {
	// varName is the name of the variable being assigned
	varName string
	// sourceExpr is the source slice expression (for len(source))
	sourceExpr string
	// pos is the position of the make call
	pos token.Pos
}

// runVar030 runs the KTN-VAR-030 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: any error that occurred
func runVar030(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar030) {
		// Rule disabled
		return nil, nil
	}

	// Get AST inspector
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Check append([]T(nil), s...) pattern
	checkAppendNilPattern(pass, insp, cfg)

	// Check make + copy pattern
	checkMakeCopyPattern(pass, insp, cfg)

	// Return analysis result
	return nil, nil
}

// checkAppendNilPattern checks for append([]T(nil), s...) pattern.
//
// Params:
//   - pass: analysis context
//   - insp: AST inspector
//   - cfg: configuration
func checkAppendNilPattern(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Node types to analyze
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	// Traverse all call expressions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar030, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Cast to call expression
		call, ok := n.(*ast.CallExpr)
		// Check if valid
		if !ok {
			// Not a call expression
			return
		}

		// Check if it's an append call with ellipsis
		if !isAppendNilPattern(call) {
			// Not the pattern we're looking for
			return
		}

		// Report the issue
		msg, _ := messages.Get(ruleCodeVar030)
		pass.Reportf(
			call.Pos(),
			"%s: %s",
			ruleCodeVar030,
			msg.Format(cfg.Verbose, "append([]T(nil), s...)"),
		)
	})
}

// isAppendNilPattern checks if a call is append([]T(nil), s...).
//
// Params:
//   - call: call expression to check
//
// Returns:
//   - bool: true if it matches the pattern
func isAppendNilPattern(call *ast.CallExpr) bool {
	// Check if it's an append call
	ident, ok := call.Fun.(*ast.Ident)
	// Verify it's append
	if !ok || ident.Name != "append" {
		// Not an append call
		return false
	}

	// Check that there are exactly 2 arguments
	if len(call.Args) != 2 {
		// Not matching pattern
		return false
	}

	// Check if it has ellipsis
	if !call.Ellipsis.IsValid() {
		// No ellipsis
		return false
	}

	// Check first argument is []T(nil)
	return isNilSliceConversion(call.Args[0])
}

// isNilSliceConversion checks if an expression is []T(nil).
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - bool: true if it's []T(nil)
func isNilSliceConversion(expr ast.Expr) bool {
	// Check if it's a call expression (type conversion)
	call, ok := expr.(*ast.CallExpr)
	// Verify it's a call
	if !ok {
		// Not a call expression
		return false
	}

	// Check that there's exactly one argument
	if len(call.Args) != 1 {
		// Not matching pattern
		return false
	}

	// Check if the argument is nil
	ident, ok := call.Args[0].(*ast.Ident)
	// Verify it's nil
	if !ok || ident.Name != "nil" {
		// Not nil
		return false
	}

	// Check if the function is a slice type
	_, ok = call.Fun.(*ast.ArrayType)
	// Return result
	return ok
}

// checkMakeCopyPattern checks for make([]T, len(s)) + copy(clone, s) pattern.
//
// Params:
//   - pass: analysis context
//   - insp: AST inspector
//   - cfg: configuration
func checkMakeCopyPattern(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Node types to analyze
	nodeFilter := []ast.Node{
		(*ast.BlockStmt)(nil),
	}

	// Traverse all block statements
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar030, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Cast to block statement
		block, ok := n.(*ast.BlockStmt)
		// Check if valid
		if !ok {
			// Not a block statement
			return
		}

		// Analyze the block for make+copy patterns
		analyzeBlockForMakeCopy(pass, block, cfg)
	})
}

// analyzeBlockForMakeCopy analyzes a block for make+copy patterns.
//
// Params:
//   - pass: analysis context
//   - block: block statement to analyze
//   - cfg: configuration
func analyzeBlockForMakeCopy(pass *analysis.Pass, block *ast.BlockStmt, cfg *config.Config) {
	// Map to track make calls that could be clone patterns
	makeInfos := make(map[string]*makeCloneInfo)

	// Iterate through statements
	for _, stmt := range block.List {
		// Check for assignment with make
		if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			// Process assignment for make pattern
			processMakeAssignment(assignStmt, makeInfos)
		}

		// Check for expression statement with copy
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			// Check for copy call
			checkCopyCall(pass, exprStmt, makeInfos, cfg)
		}
	}
}

// processMakeAssignment processes an assignment that might contain make.
//
// Params:
//   - assignStmt: assignment statement to check
//   - makeInfos: map to store make info
func processMakeAssignment(assignStmt *ast.AssignStmt, makeInfos map[string]*makeCloneInfo) {
	// Check that there's exactly one LHS and RHS
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		// Not matching pattern
		return
	}

	// Get the variable name
	varIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return
	}

	// Check if RHS is a make call
	makeCall, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	// Verify it's a call
	if !ok {
		// Not a call expression
		return
	}

	// Check if it's a make call
	makeIdent, ok := makeCall.Fun.(*ast.Ident)
	// Verify it's make
	if !ok || makeIdent.Name != "make" {
		// Not a make call
		return
	}

	// Check that there are exactly 2 arguments (type and length)
	if len(makeCall.Args) != 2 {
		// Not matching pattern (might have capacity)
		return
	}

	// Check if first arg is a slice type
	if _, ok := makeCall.Args[0].(*ast.ArrayType); !ok {
		// Not a slice type
		return
	}

	// Check if second arg is len(source)
	sourceExpr := extractLenSource(makeCall.Args[1])
	// Verify it's valid
	if sourceExpr == "" {
		// Not a len() call
		return
	}

	// Store the make info
	makeInfos[varIdent.Name] = &makeCloneInfo{
		varName:    varIdent.Name,
		sourceExpr: sourceExpr,
		pos:        assignStmt.Pos(),
	}
}

// extractLenSource extracts the source from a len(source) call.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - string: source expression name or empty if not a len call
func extractLenSource(expr ast.Expr) string {
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

	// Return the source expression as string
	return exprToString(call.Args[0])
}

// exprToString converts an expression to a string representation.
//
// Params:
//   - expr: expression to convert
//
// Returns:
//   - string: string representation
func exprToString(expr ast.Expr) string {
	// Handle different expression types
	switch e := expr.(type) {
	// Identifier case
	case *ast.Ident:
		// Return identifier name
		return e.Name
	// Selector case (e.g., obj.field)
	case *ast.SelectorExpr:
		// Build selector string
		return exprToString(e.X) + "." + e.Sel.Name
	// Index case (e.g., arr[i])
	case *ast.IndexExpr:
		// Build index string
		return exprToString(e.X) + "[" + exprToString(e.Index) + "]"
	// Default case
	default:
		// Unknown expression type
		return ""
	}
}

// checkCopyCall checks if an expression statement is a copy call matching a make.
//
// Params:
//   - pass: analysis context
//   - exprStmt: expression statement to check
//   - makeInfos: map of make infos
//   - cfg: configuration
func checkCopyCall(pass *analysis.Pass, exprStmt *ast.ExprStmt, makeInfos map[string]*makeCloneInfo, cfg *config.Config) {
	// Get the call expression
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

	// Check that there are exactly 2 arguments
	if len(call.Args) != 2 {
		// Wrong number of arguments
		return
	}

	// Get the destination variable name
	destName := exprToString(call.Args[0])
	// Check if valid
	if destName == "" {
		// Invalid destination
		return
	}

	// Get the source expression
	sourceName := exprToString(call.Args[1])
	// Check if valid
	if sourceName == "" {
		// Invalid source
		return
	}

	// Check if there's a matching make
	makeInfo, exists := makeInfos[destName]
	// Verify it exists
	if !exists {
		// No matching make
		return
	}

	// Check if sources match
	if makeInfo.sourceExpr != sourceName {
		// Sources don't match
		return
	}

	// Report the issue
	msg, _ := messages.Get(ruleCodeVar030)
	pass.Reportf(
		call.Pos(),
		"%s: %s",
		ruleCodeVar030,
		msg.Format(cfg.Verbose, "make+copy"),
	)

	// Remove from map to avoid duplicate reports
	delete(makeInfos, destName)
}
