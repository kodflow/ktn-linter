// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar037 is the rule code for this analyzer.
	ruleCodeVar037 string = "KTN-VAR-037"
	// collectionTypeKeys indicates map keys collection.
	collectionTypeKeys string = "keys"
	// collectionTypeValues indicates map values collection.
	collectionTypeValues string = "values"
)

// Analyzer037 detects manual map key/value collection patterns.
//
// Since Go 1.23, maps.Keys() and maps.Values() return iterators.
// Combined with slices.Collect(), these provide idiomatic ways to
// collect map keys or values into slices.
var Analyzer037 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar037",
	Doc:      "KTN-VAR-037: Utiliser slices.Collect(maps.Keys/Values()) au lieu de boucle range (Go 1.23+)",
	Run:      runVar037,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar037 runs the KTN-VAR-037 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result (always nil)
//   - error: any error that occurred
func runVar037(pass *analysis.Pass) (any, error) {
	// Get configuration
	cfg := config.Get()

	// Check if rule is enabled
	if !cfg.IsRuleEnabled(ruleCodeVar037) {
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

	// Check for manual map key/value collection patterns
	checkMapCollectionPatterns(pass, insp, cfg)

	// Return analysis result
	return nil, nil
}

// checkMapCollectionPatterns checks for manual map collection patterns.
//
// Params:
//   - pass: analysis context
//   - insp: AST inspector
//   - cfg: configuration
func checkMapCollectionPatterns(pass *analysis.Pass, insp *inspector.Inspector, cfg *config.Config) {
	// Node types to analyze
	nodeFilter := []ast.Node{
		(*ast.RangeStmt)(nil),
	}

	// Traverse all range statements
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar037, pass.Fset.Position(n.Pos()).Filename) {
			// File excluded
			return
		}

		// Cast to range statement
		rangeStmt, ok := n.(*ast.RangeStmt)
		// Check if valid
		if !ok {
			// Not a range statement
			return
		}

		// Check if ranging over a map
		if !isRangingOverMap(pass, rangeStmt) {
			// Not a map range
			return
		}

		// Check if it's a simple append pattern
		collectionType := detectMapCollectionPattern(rangeStmt)
		// Check if pattern was detected
		if collectionType == "" {
			// No pattern detected
			return
		}

		// Report the issue
		reportMapCollectionPattern(pass, rangeStmt, collectionType, cfg)
	})
}

// isRangingOverMap checks if range statement is iterating over a map.
//
// Params:
//   - pass: analysis context
//   - rangeStmt: range statement to check
//
// Returns:
//   - bool: true if ranging over a map
func isRangingOverMap(pass *analysis.Pass, rangeStmt *ast.RangeStmt) bool {
	// Get the type of the range expression
	tv, ok := pass.TypesInfo.Types[rangeStmt.X]
	// Check if type info available
	if !ok {
		// No type info
		return false
	}

	// Check if it's a map type
	_, isMap := tv.Type.Underlying().(*types.Map)
	// Return result
	return isMap
}

// detectMapCollectionPattern detects if range is collecting keys or values.
//
// Params:
//   - rangeStmt: range statement to analyze
//
// Returns:
//   - string: "keys", "values", or "" if no pattern detected
func detectMapCollectionPattern(rangeStmt *ast.RangeStmt) string {
	// Check that body has exactly one statement
	if rangeStmt.Body == nil || len(rangeStmt.Body.List) != 1 {
		// Wrong number of statements
		return ""
	}

	// Get the statement (should be an expression or assignment)
	stmt := rangeStmt.Body.List[0]

	// Check for assignment statement (x = append(x, k))
	assignStmt, ok := stmt.(*ast.AssignStmt)
	// Check if valid
	if !ok {
		// Not an assignment, might be expression statement
		return ""
	}

	// Check that it's a single assignment
	if len(assignStmt.Rhs) != 1 || len(assignStmt.Lhs) != 1 {
		// Not a simple assignment
		return ""
	}

	// Check that RHS is an append call
	appendCall, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	// Check if valid
	if !ok {
		// Not a function call
		return ""
	}

	// Check that it's an append call
	if !isVar037AppendCall(appendCall) {
		// Not an append call
		return ""
	}

	// Check that append has exactly 2 arguments
	if len(appendCall.Args) != 2 {
		// Wrong number of arguments
		return ""
	}

	// Detect if collecting keys or values
	return detectCollectionType(rangeStmt, appendCall)
}

// isVar037AppendCall checks if a call expression is a call to append.
//
// Params:
//   - call: call expression to check
//
// Returns:
//   - bool: true if it's an append call
func isVar037AppendCall(call *ast.CallExpr) bool {
	// Get the function identifier
	ident, ok := call.Fun.(*ast.Ident)
	// Check if valid
	if !ok {
		// Not an identifier
		return false
	}

	// Check if it's append
	return ident.Name == "append"
}

// detectCollectionType detects whether collecting keys or values.
//
// Params:
//   - rangeStmt: range statement
//   - appendCall: the append call expression
//
// Returns:
//   - string: "keys", "values", or "" if not a simple collection
func detectCollectionType(rangeStmt *ast.RangeStmt, appendCall *ast.CallExpr) string {
	// Get the second argument of append (the value being appended)
	appendedValue := appendCall.Args[1]

	// Check if appending the key
	if isKeyCollection(rangeStmt, appendedValue) {
		// Collecting keys
		return collectionTypeKeys
	}

	// Check if appending the value
	if isValueCollection(rangeStmt, appendedValue) {
		// Collecting values
		return collectionTypeValues
	}

	// Not a simple collection pattern
	return ""
}

// isKeyCollection checks if appending the map key without transformation.
//
// Params:
//   - rangeStmt: range statement
//   - appendedValue: value being appended
//
// Returns:
//   - bool: true if appending the key directly
func isKeyCollection(rangeStmt *ast.RangeStmt, appendedValue ast.Expr) bool {
	// Range must have a key and blank/nil value
	if rangeStmt.Key == nil {
		// No key in range
		return false
	}

	// Value should be blank or nil (for k := range m)
	if !isVar037BlankOrNil(rangeStmt.Value) {
		// Value is used
		return false
	}

	// Check that appended value is the key
	return isSameIdent(rangeStmt.Key, appendedValue)
}

// isValueCollection checks if appending the map value without transformation.
//
// Params:
//   - rangeStmt: range statement
//   - appendedValue: value being appended
//
// Returns:
//   - bool: true if appending the value directly
func isValueCollection(rangeStmt *ast.RangeStmt, appendedValue ast.Expr) bool {
	// Range must have a value
	if rangeStmt.Value == nil {
		// No value in range
		return false
	}

	// Key should be blank or nil (for _, v := range m)
	if !isVar037BlankOrNil(rangeStmt.Key) {
		// Key is used for something else
		return false
	}

	// Check that appended value is the range value
	return isSameIdent(rangeStmt.Value, appendedValue)
}

// isVar037BlankOrNil checks if expression is blank identifier or nil.
//
// Params:
//   - expr: expression to check
//
// Returns:
//   - bool: true if blank or nil
func isVar037BlankOrNil(expr ast.Expr) bool {
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

// isSameIdent checks if two expressions are the same identifier.
//
// Params:
//   - expr1: first expression
//   - expr2: second expression
//
// Returns:
//   - bool: true if same identifier
func isSameIdent(expr1, expr2 ast.Expr) bool {
	// Get first identifier
	ident1, ok1 := expr1.(*ast.Ident)
	// Check if valid
	if !ok1 {
		// Not an identifier
		return false
	}

	// Get second identifier
	ident2, ok2 := expr2.(*ast.Ident)
	// Check if valid
	if !ok2 {
		// Not an identifier
		return false
	}

	// Compare names
	return ident1.Name == ident2.Name
}

// reportMapCollectionPattern reports the map collection pattern issue.
//
// Params:
//   - pass: analysis context
//   - rangeStmt: range statement to report on
//   - collectionType: "keys" or "values"
//   - cfg: configuration
func reportMapCollectionPattern(
	pass *analysis.Pass,
	rangeStmt *ast.RangeStmt,
	collectionType string,
	cfg *config.Config,
) {
	// Get message and format based on collection type
	msg, ok := messages.Get(ruleCodeVar037)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(rangeStmt.Body.List[0].Pos(), "%s: utiliser slices.Collect(maps.%s()) (Go 1.23+)", ruleCodeVar037, collectionType)
		return
	}

	// Report the issue on the append statement
	pass.Reportf(
		rangeStmt.Body.List[0].Pos(),
		"%s: %s",
		ruleCodeVar037,
		msg.Format(cfg.Verbose, collectionType),
	)
}
