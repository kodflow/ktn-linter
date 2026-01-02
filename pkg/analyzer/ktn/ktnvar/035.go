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
	// ruleCodeVar035 is the rule code for this analyzer.
	ruleCodeVar035 string = "KTN-VAR-035"
)

// Analyzer035 detects manual contains patterns that can use slices.Contains (Go 1.21+).
//
// Detects:
//   - for _, x := range s { if x == target { return true } } return false
var Analyzer035 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar035",
	Doc:      "KTN-VAR-035: Utiliser slices.Contains() au lieu du pattern for-range manuel (Go 1.21+)",
	Run:      runVar035,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar035 runs the KTN-VAR-035 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: any error that occurred
func runVar035(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar035) {
		// Rule disabled
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp := inspAny.(*inspector.Inspector)

	// Check for manual contains pattern in functions
	checkContainsPattern(pass, insp, cfg)

	// Return analysis result
	return nil, nil
}

// checkContainsPattern checks for manual contains pattern in functions.
//
// Params:
//   - pass: analysis context
//   - insp: AST inspector
//   - cfg: configuration
func checkContainsPattern(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Node types to analyze
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	// Traverse all function declarations and literals
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar035, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Get function body
		var body *ast.BlockStmt
		// Check function type
		switch fn := n.(type) {
		// Function declaration case
		case *ast.FuncDecl:
			// Skip functions without body
			if fn.Body == nil {
				// No body to analyze
				return
			}
			// Get body from declaration
			body = fn.Body
		// Function literal case
		case *ast.FuncLit:
			// Get body from literal
			body = fn.Body
		// Default case
		default:
			// Unknown function type
			return
		}

		// Analyze function body for contains pattern
		analyzeBodyForContainsPattern(pass, body, cfg)
	})
}

// analyzeBodyForContainsPattern analyzes a function body for contains pattern.
//
// Params:
//   - pass: analysis context
//   - body: function body to analyze
//   - cfg: configuration
func analyzeBodyForContainsPattern(pass *analysis.Pass, body *ast.BlockStmt, cfg *config.Config) {
	// Check that body has at least 2 statements
	if len(body.List) < 2 {
		// Not enough statements for pattern
		return
	}

	// Check each pair of consecutive statements
	for idx := 0; idx < len(body.List)-1; idx++ {
		// Get the range statement
		rangeStmt, ok := body.List[idx].(*ast.RangeStmt)
		// Check if valid range
		if !ok {
			// Not a range statement
			continue
		}

		// Get the return statement after range
		returnStmt, ok := body.List[idx+1].(*ast.ReturnStmt)
		// Check if valid return
		if !ok {
			// Not a return statement
			continue
		}

		// Check if it matches the contains pattern
		if matchesContainsPattern(rangeStmt, returnStmt) {
			// Report the issue at the if statement position
			reportContainsPattern(pass, rangeStmt, cfg)
		}
	}
}

// matchesContainsPattern checks if range+return matches contains pattern.
//
// Params:
//   - rangeStmt: range statement to check
//   - returnStmt: return statement to check
//
// Returns:
//   - bool: true if it matches the pattern
func matchesContainsPattern(rangeStmt *ast.RangeStmt, returnStmt *ast.ReturnStmt) bool {
	// Check that range uses blank identifier for key
	if !isBlankOrNil(rangeStmt.Key) {
		// Key is used, might not be simple contains
		return false
	}

	// Check that range has a value variable
	if rangeStmt.Value == nil {
		// No value variable
		return false
	}

	// Check that range body has exactly one statement
	if rangeStmt.Body == nil || len(rangeStmt.Body.List) != 1 {
		// Wrong number of statements
		return false
	}

	// Get the if statement from range body
	ifStmt, ok := rangeStmt.Body.List[0].(*ast.IfStmt)
	// Check if valid
	if !ok {
		// Not an if statement
		return false
	}

	// Check the if statement matches the pattern
	if !matchesIfReturnTruePattern(ifStmt, rangeStmt.Value) {
		// If doesn't match pattern
		return false
	}

	// Check that return statement returns false
	return returnsLiteralBool(returnStmt, false)
}

// isBlankOrNil checks if expression is blank identifier or nil.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - bool: true if blank or nil
func isBlankOrNil(expr ast.Expr) bool {
	// Check for nil
	if expr == nil {
		// Nil is ok
		return true
	}

	// Check for blank identifier
	ident, ok := expr.(*ast.Ident)
	// Return result
	return ok && ident.Name == "_"
}

// matchesIfReturnTruePattern checks if an if statement matches equality check with return true.
//
// Params:
//   - ifStmt: if statement to check
//   - rangeValue: the range value variable
//
// Returns:
//   - bool: true if it matches the pattern
func matchesIfReturnTruePattern(ifStmt *ast.IfStmt, rangeValue ast.Expr) bool {
	// Check that if has no else
	if ifStmt.Else != nil {
		// Has else, might be more complex
		return false
	}

	// Check the condition is an equality comparison
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	// Verify it's a binary expression
	if !ok {
		// Not a binary expression
		return false
	}

	// Check that operator is ==
	if binaryExpr.Op != token.EQL {
		// Not an equality comparison
		return false
	}

	// Check that one side is the range value
	if !containsRangeValue(binaryExpr, rangeValue) {
		// Range value not in comparison
		return false
	}

	// Check that if body has exactly one statement
	if ifStmt.Body == nil || len(ifStmt.Body.List) != 1 {
		// Wrong number of statements
		return false
	}

	// Get the return statement from if body
	returnStmt, ok := ifStmt.Body.List[0].(*ast.ReturnStmt)
	// Check if valid
	if !ok {
		// Not a return statement
		return false
	}

	// Check that it returns true
	return returnsLiteralBool(returnStmt, true)
}

// containsRangeValue checks if binary expr uses the range value.
//
// Params:
//   - binaryExpr: binary expression to check
//   - rangeValue: the range value to look for
//
// Returns:
//   - bool: true if range value is used
func containsRangeValue(binaryExpr *ast.BinaryExpr, rangeValue ast.Expr) bool {
	// Get range value name
	rangeIdent, ok := rangeValue.(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return false
	}

	// Check left side
	if leftIdent, ok := binaryExpr.X.(*ast.Ident); ok {
		// Compare names
		if leftIdent.Name == rangeIdent.Name {
			// Found on left
			return true
		}
	}

	// Check right side
	if rightIdent, ok := binaryExpr.Y.(*ast.Ident); ok {
		// Compare names
		if rightIdent.Name == rangeIdent.Name {
			// Found on right
			return true
		}
	}

	// Not found
	return false
}

// returnsLiteralBool checks if return statement returns a boolean literal.
//
// Params:
//   - returnStmt: return statement to check
//   - expected: expected boolean value
//
// Returns:
//   - bool: true if returns expected boolean
func returnsLiteralBool(returnStmt *ast.ReturnStmt, expected bool) bool {
	// Check that there's exactly one return value
	if len(returnStmt.Results) != 1 {
		// Wrong number of return values
		return false
	}

	// Get the return value
	ident, ok := returnStmt.Results[0].(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return false
	}

	// Check the value
	if expected {
		// Check for true
		return ident.Name == "true"
	}

	// Check for false
	return ident.Name == "false"
}

// reportContainsPattern reports the contains pattern issue.
//
// Params:
//   - pass: analysis context
//   - rangeStmt: range statement to report on
//   - cfg: configuration
func reportContainsPattern(pass *analysis.Pass, rangeStmt *ast.RangeStmt, cfg *config.Config) {
	// Get the if statement position for more precise reporting
	ifStmt := rangeStmt.Body.List[0].(*ast.IfStmt)

	// Report the issue
	msg, _ := messages.Get(ruleCodeVar035)
	pass.Reportf(
		ifStmt.Cond.Pos(),
		"%s: %s",
		ruleCodeVar035,
		msg.Format(cfg.Verbose),
	)
}
