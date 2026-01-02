package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"
)

// TestCheckLoopVarCopyPattern tests checkLoopVarCopyPattern function.
func TestCheckLoopVarCopyPattern(t *testing.T) {
	// Test: no range variables
	rangeNoVars := &ast.RangeStmt{
		Key:   nil,
		Value: nil,
		Body:  &ast.BlockStmt{},
	}
	checkLoopVarCopyPattern(nil, rangeNoVars)

	// Test: nil body
	rangeNilBody := &ast.RangeStmt{
		Key:  &ast.Ident{Name: "i"},
		Body: nil,
	}
	checkLoopVarCopyPattern(nil, rangeNilBody)

	// Test: empty body
	rangeEmptyBody := &ast.RangeStmt{
		Key:  &ast.Ident{Name: "i"},
		Body: &ast.BlockStmt{List: nil},
	}
	checkLoopVarCopyPattern(nil, rangeEmptyBody)

	// Test: key is blank identifier
	rangeBlankKey := &ast.RangeStmt{
		Key:  &ast.Ident{Name: "_"},
		Body: &ast.BlockStmt{},
	}
	checkLoopVarCopyPattern(nil, rangeBlankKey)

	// Test: value is blank identifier
	rangeBlankValue := &ast.RangeStmt{
		Key:   &ast.Ident{Name: "i"},
		Value: &ast.Ident{Name: "_"},
		Body:  &ast.BlockStmt{},
	}
	checkLoopVarCopyPattern(nil, rangeBlankValue)
}

// TestGetRangeVariableNames tests getRangeVariableNames function.
func TestGetRangeVariableNames(t *testing.T) {
	// Test: no variables
	rangeNoVars := &ast.RangeStmt{
		Key:   nil,
		Value: nil,
	}
	result := getRangeVariableNames(rangeNoVars)
	// Expected: empty
	if len(result) != 0 {
		t.Error("getRangeVariableNames should return empty map for no vars")
	}

	// Test: only key
	rangeKeyOnly := &ast.RangeStmt{
		Key:   &ast.Ident{Name: "i"},
		Value: nil,
	}
	result = getRangeVariableNames(rangeKeyOnly)
	// Expected: one entry
	if len(result) != 1 || !result["i"] {
		t.Error("getRangeVariableNames should include key")
	}

	// Test: only value
	rangeValueOnly := &ast.RangeStmt{
		Key:   nil,
		Value: &ast.Ident{Name: "v"},
	}
	result = getRangeVariableNames(rangeValueOnly)
	// Expected: one entry
	if len(result) != 1 || !result["v"] {
		t.Error("getRangeVariableNames should include value")
	}

	// Test: both key and value
	rangeBoth := &ast.RangeStmt{
		Key:   &ast.Ident{Name: "i"},
		Value: &ast.Ident{Name: "v"},
	}
	result = getRangeVariableNames(rangeBoth)
	// Expected: two entries
	if len(result) != 2 || !result["i"] || !result["v"] {
		t.Error("getRangeVariableNames should include both key and value")
	}

	// Test: blank key
	rangeBlankKey := &ast.RangeStmt{
		Key:   &ast.Ident{Name: "_"},
		Value: &ast.Ident{Name: "v"},
	}
	result = getRangeVariableNames(rangeBlankKey)
	// Expected: only value
	if len(result) != 1 || !result["v"] {
		t.Error("getRangeVariableNames should not include blank key")
	}

	// Test: key is not an ident
	rangeNonIdentKey := &ast.RangeStmt{
		Key:   &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
		Value: &ast.Ident{Name: "v"},
	}
	result = getRangeVariableNames(rangeNonIdentKey)
	// Expected: only value
	if len(result) != 1 || !result["v"] {
		t.Error("getRangeVariableNames should not include non-ident key")
	}

	// Test: value is not an ident
	rangeNonIdentValue := &ast.RangeStmt{
		Key:   &ast.Ident{Name: "i"},
		Value: &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
	}
	result = getRangeVariableNames(rangeNonIdentValue)
	// Expected: only key
	if len(result) != 1 || !result["i"] {
		t.Error("getRangeVariableNames should not include non-ident value")
	}
}

// TestCheckShortVarDecl tests checkShortVarDecl function.
func TestCheckShortVarDecl(t *testing.T) {
	rangeVars := map[string]bool{"v": true}

	// Test: not an assignment
	notAssign := &ast.ExprStmt{X: &ast.Ident{Name: "x"}}
	checkShortVarDecl(nil, notAssign, rangeVars)

	// Test: assignment with = instead of :=
	assignNotDefine := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "v"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "v"}},
	}
	checkShortVarDecl(nil, assignNotDefine, rangeVars)
}

// TestCheckAssignmentPair tests checkAssignmentPair function.
func TestCheckAssignmentPair(t *testing.T) {
	rangeVars := map[string]bool{"v": true}

	// Test: index out of bounds
	assignOutOfBounds := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "v"}, &ast.Ident{Name: "x"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "v"}},
	}
	checkAssignmentPair(nil, assignOutOfBounds, 1, rangeVars)

	// Test: LHS is not an ident
	assignLhsNotIdent := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
		Rhs: []ast.Expr{&ast.Ident{Name: "v"}},
	}
	checkAssignmentPair(nil, assignLhsNotIdent, 0, rangeVars)

	// Test: RHS is not an ident
	assignRhsNotIdent := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "v"}},
		Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: "getV"}}},
	}
	checkAssignmentPair(nil, assignRhsNotIdent, 0, rangeVars)

	// Test: names don't match
	assignNamesDiffer := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "v"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "x"}},
	}
	checkAssignmentPair(nil, assignNamesDiffer, 0, rangeVars)

	// Test: variable is not a range variable
	notRangeVar := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "x"}},
	}
	checkAssignmentPair(nil, notRangeVar, 0, rangeVars)
}
