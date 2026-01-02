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
	// ruleCodeVar033 is the rule code for this analyzer.
	ruleCodeVar033 string = "KTN-VAR-033"
)

// Analyzer033 detects patterns that can use cmp.Or (Go 1.22+).
//
// Detects:
//   - if x != 0 { return x } return default (int)
//   - if x != "" { return x } return default (string)
//   - if x != nil { return x } return default (pointer/slice/map)
var Analyzer033 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar033",
	Doc:      "KTN-VAR-033: Utiliser cmp.Or() au lieu du pattern if x != zeroValue (Go 1.22+)",
	Run:      runVar033,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar033 runs the KTN-VAR-033 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: any error that occurred
func runVar033(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar033) {
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
	// Defensive: avoid nil dereference when resolving positions
	if pass.Fset == nil {
		return nil, nil
	}

	// Check for cmp.Or patterns in functions
	checkCmpOrPattern(pass, insp, cfg)

	// Return analysis result
	return nil, nil
}

// checkCmpOrPattern checks for patterns that can use cmp.Or.
//
// Params:
//   - pass: analysis context
//   - insp: AST inspector
//   - cfg: configuration
func checkCmpOrPattern(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Node types to analyze
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	// Traverse all function declarations and literals
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar033, pass.Fset.Position(n.Pos()).Filename) {
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

		// Analyze function body for cmp.Or pattern
		analyzeBodyForCmpOrPattern(pass, body, cfg)
	})
}

// analyzeBodyForCmpOrPattern analyzes a function body for cmp.Or patterns.
//
// Params:
//   - pass: analysis context
//   - body: function body to analyze
//   - cfg: configuration
func analyzeBodyForCmpOrPattern(pass *analysis.Pass, body *ast.BlockStmt, cfg *config.Config) {
	// Check that body has at least 2 statements
	if len(body.List) < 2 {
		// Not enough statements for pattern
		return
	}

	// Check each pair of consecutive statements
	for idx := 0; idx < len(body.List)-1; idx++ {
		// Get the if statement
		ifStmt, ok := body.List[idx].(*ast.IfStmt)
		// Check if valid if statement
		if !ok {
			// Not an if statement
			continue
		}

		// Get the return statement after if
		returnStmt, ok := body.List[idx+1].(*ast.ReturnStmt)
		// Check if valid return
		if !ok {
			// Not a return statement
			continue
		}

		// Check if it matches the cmp.Or pattern
		if matchesCmpOrPattern(ifStmt, returnStmt) {
			// Report the issue at the if condition position
			reportCmpOrPattern(pass, ifStmt, cfg)
		}
	}
}

// matchesCmpOrPattern checks if if+return matches cmp.Or pattern.
//
// Params:
//   - ifStmt: if statement to check
//   - returnStmt: return statement to check
//
// Returns:
//   - bool: true if it matches the pattern
func matchesCmpOrPattern(ifStmt *ast.IfStmt, returnStmt *ast.ReturnStmt) bool {
	// Check that if has no init and no else
	if ifStmt.Init != nil || ifStmt.Else != nil {
		// Has init or else, too complex
		return false
	}

	// Check the condition is a not-equal comparison
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	// Verify it's a binary expression
	if !ok {
		// Not a binary expression
		return false
	}

	// Check that operator is !=
	if binaryExpr.Op != token.NEQ {
		// Not a not-equal comparison
		return false
	}

	// Get the variable being compared
	varExpr, zeroValue := extractVarAndZeroFor033(binaryExpr)
	// Check if valid
	if varExpr == nil || !isSimpleZeroValueFor033(zeroValue) {
		// Not a valid pattern
		return false
	}

	// Check that if body has exactly one statement
	if ifStmt.Body == nil || len(ifStmt.Body.List) != 1 {
		// Wrong number of statements
		return false
	}

	// Get the return statement from if body
	ifReturn, ok := ifStmt.Body.List[0].(*ast.ReturnStmt)
	// Check if valid
	if !ok {
		// Not a return statement
		return false
	}

	// Check that if return returns the same variable
	if !returnsVariableFor033(ifReturn, varExpr) {
		// Doesn't return the variable
		return false
	}

	// Check that after-if return has at least one result
	if len(returnStmt.Results) < 1 {
		// No return value
		return false
	}

	// Pattern matches
	return true
}

// extractVarAndZeroFor033 extracts variable and zero value from comparison.
//
// Params:
//   - binaryExpr: binary expression to analyze
//
// Returns:
//   - ast.Expr: the variable being compared
//   - ast.Expr: the zero value being compared against
func extractVarAndZeroFor033(binaryExpr *ast.BinaryExpr) (ast.Expr, ast.Expr) {
	// Check if left side is a zero value
	if isSimpleZeroValueFor033(binaryExpr.Y) {
		// Right is zero value, left is variable
		return binaryExpr.X, binaryExpr.Y
	}

	// Check if right side is a zero value
	if isSimpleZeroValueFor033(binaryExpr.X) {
		// Left is zero value, right is variable
		return binaryExpr.Y, binaryExpr.X
	}

	// No valid pattern found
	return nil, nil
}

// isSimpleZeroValueFor033 checks if expression is a zero value (0, "", nil).
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - bool: true if it's a zero value
func isSimpleZeroValueFor033(expr ast.Expr) bool {
	// Check for different types of zero values
	switch val := expr.(type) {
	// Basic literal case (0 or "")
	case *ast.BasicLit:
		// Check for integer zero
		if val.Kind == token.INT && val.Value == "0" {
			// Integer zero
			return true
		}
		// Check for empty string
		if val.Kind == token.STRING && (val.Value == `""` || val.Value == "``") {
			// Empty string
			return true
		}
		// Not a zero value literal
		return false
	// Identifier case (nil)
	case *ast.Ident:
		// Check for nil
		return val.Name == "nil"
	// Default case
	default:
		// Not a recognized zero value
		return false
	}
}

// returnsVariableFor033 checks if return statement returns the given variable.
//
// Params:
//   - returnStmt: return statement to check
//   - varExpr: variable expression to match
//
// Returns:
//   - bool: true if return statement returns the variable
func returnsVariableFor033(returnStmt *ast.ReturnStmt, varExpr ast.Expr) bool {
	// Check that there's exactly one return value
	if len(returnStmt.Results) != 1 {
		// Wrong number of return values
		return false
	}

	// Get the variable name
	varIdent, ok := varExpr.(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return false
	}

	// Get the return value
	retIdent, ok := returnStmt.Results[0].(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return false
	}

	// Compare names
	return varIdent.Name == retIdent.Name
}

// reportCmpOrPattern reports the cmp.Or pattern issue.
//
// Params:
//   - pass: analysis context
//   - ifStmt: if statement to report on
//   - cfg: configuration
func reportCmpOrPattern(pass *analysis.Pass, ifStmt *ast.IfStmt, cfg *config.Config) {
	// Report the issue at the if condition position
	msg, ok := messages.Get(ruleCodeVar033)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(ifStmt.Cond.Pos(), "%s: utiliser cmp.Or() (Go 1.22+)", ruleCodeVar033)
		return
	}
	pass.Reportf(
		ifStmt.Cond.Pos(),
		"%s: %s",
		ruleCodeVar033,
		msg.Format(cfg.Verbose),
	)
}
