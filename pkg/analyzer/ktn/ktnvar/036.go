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
	// ruleCodeVar036 is the rule code for this analyzer.
	ruleCodeVar036 string = "KTN-VAR-036"
)

// Analyzer036 detects manual index search patterns that can use slices.Index (Go 1.21+).
//
// Detects the pattern:
//
//	for i, v := range s { if v == target { return i } } return -1
var Analyzer036 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar036",
	Doc:      "KTN-VAR-036: Utiliser slices.Index() au lieu du pattern de recherche manuel (Go 1.21+)",
	Run:      runVar036,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar036 runs the KTN-VAR-036 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: any error that occurred
func runVar036(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar036) {
		// Rule disabled
		return nil, nil
	}

	// Get AST inspector
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Node types to analyze
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Traverse all function declarations
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar036, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Cast to function declaration
		funcDecl, ok := n.(*ast.FuncDecl)
		// Check if valid
		if !ok || funcDecl.Body == nil {
			// Not a function declaration or no body
			return
		}

		// Analyze function body for index search pattern
		analyzeFunctionForIndexPattern(pass, funcDecl, cfg)
	})

	// Return analysis result
	return nil, nil
}

// analyzeFunctionForIndexPattern analyzes a function for manual index search pattern.
//
// Params:
//   - pass: analysis context
//   - funcDecl: function declaration to analyze
//   - cfg: configuration
func analyzeFunctionForIndexPattern(pass *analysis.Pass, funcDecl *ast.FuncDecl, cfg *config.Config) {
	// Function must return int
	if !functionReturnsInt(funcDecl) {
		// Not returning int
		return
	}

	// Check if body has the pattern: for range + if == return i + return -1
	body := funcDecl.Body.List
	// Need at least 2 statements (for loop and return -1)
	if len(body) < 2 {
		// Not enough statements
		return
	}

	// Look for the pattern in the statements
	for i, stmt := range body {
		// Check if this is a range statement
		rangeStmt, ok := stmt.(*ast.RangeStmt)
		// Skip if not a range statement
		if !ok {
			// Not a range statement
			continue
		}

		// Check if next statement is return -1
		if i+1 >= len(body) {
			// No next statement
			continue
		}

		// Check if return -1 follows
		if !isReturnMinusOne(body[i+1]) {
			// Not return -1
			continue
		}

		// Check if range has index and value variables
		indexPattern := checkRangeForIndexPattern(rangeStmt)
		// Report if pattern detected
		if indexPattern != nil {
			// Report the issue
			reportIndexPattern(pass, indexPattern.pos, cfg)
		}
	}
}

// indexPatternInfo stores information about a detected index pattern.
type indexPatternInfo struct {
	// pos is the position where to report the issue
	pos token.Pos
}

// functionReturnsInt checks if a function returns a single int.
//
// Params:
//   - funcDecl: function declaration to check
//
// Returns:
//   - bool: true if function returns int
func functionReturnsInt(funcDecl *ast.FuncDecl) bool {
	// Check function has results
	if funcDecl.Type.Results == nil {
		// No results
		return false
	}

	// Check there's exactly one result
	if len(funcDecl.Type.Results.List) != 1 {
		// Not exactly one result
		return false
	}

	// Get the result type
	resultType := funcDecl.Type.Results.List[0].Type
	// Check if it's an ident
	ident, ok := resultType.(*ast.Ident)
	// Check if it's int
	if !ok || ident.Name != "int" {
		// Not int type
		return false
	}

	// Returns int
	return true
}

// isReturnMinusOne checks if a statement is return -1.
//
// Params:
//   - stmt: statement to check
//
// Returns:
//   - bool: true if it's return -1
func isReturnMinusOne(stmt ast.Stmt) bool {
	// Check if it's a return statement
	retStmt, ok := stmt.(*ast.ReturnStmt)
	// Verify it's a return
	if !ok {
		// Not a return statement
		return false
	}

	// Check there's exactly one return value
	if len(retStmt.Results) != 1 {
		// Not exactly one result
		return false
	}

	// Check if it's a unary expression with minus
	unaryExpr, ok := retStmt.Results[0].(*ast.UnaryExpr)
	// Verify it's unary
	if !ok || unaryExpr.Op != token.SUB {
		// Not a minus operator
		return false
	}

	// Check if operand is literal 1
	basicLit, ok := unaryExpr.X.(*ast.BasicLit)
	// Verify it's a basic literal
	if !ok || basicLit.Kind != token.INT || basicLit.Value != "1" {
		// Not literal 1
		return false
	}

	// It's return -1
	return true
}

// checkRangeForIndexPattern checks if a range statement matches the index search pattern.
//
// Params:
//   - rangeStmt: range statement to check
//
// Returns:
//   - *indexPatternInfo: pattern info if detected, nil otherwise
func checkRangeForIndexPattern(rangeStmt *ast.RangeStmt) *indexPatternInfo {
	// Check range has both key and value
	if rangeStmt.Key == nil || rangeStmt.Value == nil {
		// Missing key or value
		return nil
	}

	// Get the index variable name
	indexIdent, ok := rangeStmt.Key.(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return nil
	}

	// Get the value variable name
	valueIdent, ok := rangeStmt.Value.(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return nil
	}

	// Check the body for if == return i pattern
	body := rangeStmt.Body
	// Verify body exists
	if body == nil || len(body.List) == 0 {
		// No body or empty body
		return nil
	}

	// Look for if statement with equality check and return index
	for _, stmt := range body.List {
		// Check if it matches the pattern
		if checkIfEqualReturnIndex(stmt, indexIdent.Name, valueIdent.Name) {
			// Pattern detected - report at the if statement position
			return &indexPatternInfo{pos: stmt.Pos()}
		}
	}

	// No pattern found
	return nil
}

// checkIfEqualReturnIndex checks if a statement is if v == target { return i }.
//
// Params:
//   - stmt: statement to check
//   - indexName: name of the index variable
//   - valueName: name of the value variable
//
// Returns:
//   - bool: true if pattern matches
func checkIfEqualReturnIndex(stmt ast.Stmt, indexName, valueName string) bool {
	// Check if it's an if statement
	ifStmt, ok := stmt.(*ast.IfStmt)
	// Verify it's an if
	if !ok {
		// Not an if statement
		return false
	}

	// Check if condition is an equality binary expression
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	// Verify it's binary
	if !ok || binaryExpr.Op != token.EQL {
		// Not an equality comparison
		return false
	}

	// Check if one side uses the value variable
	if !expressionUsesIdent(binaryExpr.X, valueName) && !expressionUsesIdent(binaryExpr.Y, valueName) {
		// Value variable not used in comparison
		return false
	}

	// Check if body returns the index variable
	return ifBodyReturnsIndex(ifStmt.Body, indexName)
}

// expressionUsesIdent checks if an expression uses a specific identifier.
//
// Params:
//   - expr: expression to check
//   - name: identifier name to look for
//
// Returns:
//   - bool: true if identifier is used
func expressionUsesIdent(expr ast.Expr, name string) bool {
	// Check if it's an identifier
	ident, ok := expr.(*ast.Ident)
	// Return if it matches the name
	return ok && ident.Name == name
}

// ifBodyReturnsIndex checks if an if body contains return with the index variable.
//
// Params:
//   - body: block statement to check
//   - indexName: name of the index variable
//
// Returns:
//   - bool: true if body returns index
func ifBodyReturnsIndex(body *ast.BlockStmt, indexName string) bool {
	// Check body has statements
	if body == nil || len(body.List) == 0 {
		// Empty body
		return false
	}

	// Look for return statement
	for _, stmt := range body.List {
		// Check if it's a return statement
		retStmt, ok := stmt.(*ast.ReturnStmt)
		// Continue if not return
		if !ok {
			// Not a return
			continue
		}

		// Check there's exactly one return value
		if len(retStmt.Results) != 1 {
			// Not exactly one result
			continue
		}

		// Check if return value is the index variable
		ident, ok := retStmt.Results[0].(*ast.Ident)
		// Check if matches index name
		if ok && ident.Name == indexName {
			// Returns index
			return true
		}
	}

	// No return index found
	return false
}

// reportIndexPattern reports the KTN-VAR-036 issue.
//
// Params:
//   - pass: analysis context
//   - pos: position to report at
//   - cfg: configuration
func reportIndexPattern(pass *analysis.Pass, pos token.Pos, cfg *config.Config) {
	// Get the message
	msg, ok := messages.Get(ruleCodeVar036)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(pos, "%s: utiliser slices.Index() (Go 1.21+)", ruleCodeVar036)
		return
	}
	// Report the issue
	pass.Reportf(
		pos,
		"%s: %s",
		ruleCodeVar036,
		msg.Format(cfg.Verbose),
	)
}
