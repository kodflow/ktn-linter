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
	// ruleCodeVar031 is the rule code for this analyzer
	ruleCodeVar031 string = "KTN-VAR-031"
)

// Analyzer031 detects manual map cloning that can be replaced with maps.Clone
var Analyzer031 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar031",
	Doc:      "KTN-VAR-031: Utiliser maps.Clone() au lieu du pattern make+range (Go 1.21+)",
	Run:      runVar031,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar031 executes the KTN-VAR-031 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: possible error
func runVar031(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar031) {
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

	// Filter for function declarations and function literals
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar031, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Get function body
		var body *ast.BlockStmt
		// Check node type
		switch fn := n.(type) {
		// Function declaration
		case *ast.FuncDecl:
			body = fn.Body
		// Function literal
		case *ast.FuncLit:
			body = fn.Body
		}

		// Skip if no body
		if body == nil {
			// No body to analyze
			return
		}

		// Analyze the function body for map clone patterns
		analyzeBlockForMapClone(pass, body)
	})

	// Return result
	return nil, nil
}

// analyzeBlockForMapClone analyzes a block for map cloning patterns.
//
// Params:
//   - pass: analysis context
//   - block: block statement to analyze
func analyzeBlockForMapClone(pass *analysis.Pass, block *ast.BlockStmt) {
	// Track make(map) assignments: varName -> position
	mapMakes := make(map[string]token.Pos)

	// Iterate through statements
	for i, stmt := range block.List {
		// Check for make(map) assignment
		if varName, pos := extractMakeMapAssign(stmt); varName != "" {
			// Store the make position
			mapMakes[varName] = pos
			// Continue to next statement
			continue
		}

		// Check for range statement that clones a map
		rangeStmt, ok := stmt.(*ast.RangeStmt)
		// Skip if not a range statement
		if !ok {
			// Continue to next statement
			continue
		}

		// Check if this is a simple map clone pattern
		if isSimpleMapClone(rangeStmt, mapMakes, i, block.List) {
			// Report the issue
			msg, ok := messages.Get(ruleCodeVar031)
			// Defensive: avoid panic if message is missing
			if !ok {
				pass.Reportf(rangeStmt.Pos(), "%s: utiliser maps.Clone() (Go 1.21+)", ruleCodeVar031)
				continue
			}
			pass.Reportf(
				rangeStmt.Pos(),
				"%s: %s",
				ruleCodeVar031,
				msg.Format(config.Get().Verbose),
			)
		}
	}
}

// extractMakeMapAssign extracts variable name from make(map) assignment.
//
// Params:
//   - stmt: statement to check
//
// Returns:
//   - string: variable name or empty if not a make(map) assignment
//   - token.Pos: position of the statement
func extractMakeMapAssign(stmt ast.Stmt) (string, token.Pos) {
	// Check for assignment statement
	assignStmt, ok := stmt.(*ast.AssignStmt)
	// Skip if not assignment
	if !ok {
		// Return empty
		return "", token.NoPos
	}

	// Must be simple assignment with one LHS and one RHS
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		// Return empty
		return "", token.NoPos
	}

	// LHS must be identifier
	ident, ok := assignStmt.Lhs[0].(*ast.Ident)
	// Skip if not identifier
	if !ok {
		// Return empty
		return "", token.NoPos
	}

	// RHS must be make(map[K]V, ...)
	callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	// Skip if not call expression
	if !ok {
		// Return empty
		return "", token.NoPos
	}

	// Check if it's a make call
	if !isMakeMapCallExpr(callExpr) {
		// Return empty
		return "", token.NoPos
	}

	// Return variable name and position
	return ident.Name, assignStmt.Pos()
}

// isMakeMapCallExpr checks if a call expression is make(map[K]V, ...).
//
// Params:
//   - call: call expression to check
//
// Returns:
//   - bool: true if it's a make(map[K]V, ...) call
func isMakeMapCallExpr(call *ast.CallExpr) bool {
	// Check if function is "make"
	funIdent, ok := call.Fun.(*ast.Ident)
	// Skip if not identifier
	if !ok || funIdent.Name != "make" {
		// Return false
		return false
	}

	// Must have at least one argument (the map type)
	if len(call.Args) < 1 {
		// Return false
		return false
	}

	// First argument must be map type
	_, ok = call.Args[0].(*ast.MapType)
	// Return result
	return ok
}

// isSimpleMapClone checks if a range statement is a simple map clone.
//
// Params:
//   - rangeStmt: range statement to check
//   - mapMakes: tracked make(map) assignments
//   - stmtIndex: index of current statement
//   - stmts: list of statements in the block
//
// Returns:
//   - bool: true if it's a simple map clone pattern
func isSimpleMapClone(rangeStmt *ast.RangeStmt, mapMakes map[string]token.Pos, stmtIndex int, stmts []ast.Stmt) bool {
	// Range must have key and value
	if rangeStmt.Key == nil || rangeStmt.Value == nil {
		// Not a k, v range
		return false
	}

	// Get key and value names
	keyIdent, keyOk := rangeStmt.Key.(*ast.Ident)
	valIdent, valOk := rangeStmt.Value.(*ast.Ident)
	// Both must be identifiers
	if !keyOk || !valOk {
		// Return false
		return false
	}

	// Range body must have exactly one statement
	if rangeStmt.Body == nil || len(rangeStmt.Body.List) != 1 {
		// Not a simple clone
		return false
	}

	// The statement must be an assignment: clone[k] = v
	assignStmt, ok := rangeStmt.Body.List[0].(*ast.AssignStmt)
	// Skip if not assignment
	if !ok {
		// Return false
		return false
	}

	// Must be simple assignment (=) not define (:=)
	if assignStmt.Tok != token.ASSIGN {
		// Return false
		return false
	}

	// Must have one LHS and one RHS
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		// Return false
		return false
	}

	// LHS must be index expression: clone[k]
	indexExpr, ok := assignStmt.Lhs[0].(*ast.IndexExpr)
	// Skip if not index expression
	if !ok {
		// Return false
		return false
	}

	// X must be identifier (the clone map name)
	cloneIdent, ok := indexExpr.X.(*ast.Ident)
	// Skip if not identifier
	if !ok {
		// Return false
		return false
	}

	// The clone map must have been created with make before
	if _, exists := mapMakes[cloneIdent.Name]; !exists {
		// Clone not from make
		return false
	}

	// Index must be the key variable
	indexIdent, ok := indexExpr.Index.(*ast.Ident)
	// Skip if not identifier
	if !ok || indexIdent.Name != keyIdent.Name {
		// Return false
		return false
	}

	// RHS must be the value variable exactly (no transformation)
	rhsIdent, ok := assignStmt.Rhs[0].(*ast.Ident)
	// Skip if not identifier
	if !ok || rhsIdent.Name != valIdent.Name {
		// This is a transformation, not a simple clone
		return false
	}

	// This is a simple map clone pattern
	return true
}
