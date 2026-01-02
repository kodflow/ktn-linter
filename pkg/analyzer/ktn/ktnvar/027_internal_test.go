package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"
)

// TestIsConvertibleToRangeInt tests isConvertibleToRangeInt function.
func TestIsConvertibleToRangeInt(t *testing.T) {
	// Test: missing init
	forNoInit := &ast.ForStmt{
		Init: nil,
		Cond: &ast.BinaryExpr{Op: token.LSS},
		Post: &ast.IncDecStmt{Tok: token.INC},
	}
	result := isConvertibleToRangeInt(forNoInit)
	// Expected: false
	if result {
		t.Error("isConvertibleToRangeInt should return false when init is nil")
	}

	// Test: missing cond
	forNoCond := &ast.ForStmt{
		Init: &ast.AssignStmt{Tok: token.DEFINE},
		Cond: nil,
		Post: &ast.IncDecStmt{Tok: token.INC},
	}
	result = isConvertibleToRangeInt(forNoCond)
	// Expected: false
	if result {
		t.Error("isConvertibleToRangeInt should return false when cond is nil")
	}

	// Test: missing post
	forNoPost := &ast.ForStmt{
		Init: &ast.AssignStmt{Tok: token.DEFINE},
		Cond: &ast.BinaryExpr{Op: token.LSS},
		Post: nil,
	}
	result = isConvertibleToRangeInt(forNoPost)
	// Expected: false
	if result {
		t.Error("isConvertibleToRangeInt should return false when post is nil")
	}
}

// TestCheckInitIsZero tests checkInitIsZero function.
func TestCheckInitIsZero(t *testing.T) {
	// Test: valid init
	validInit := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
	}
	varName, ok := checkInitIsZero(validInit)
	// Expected: true
	if !ok || varName != "i" {
		t.Errorf("checkInitIsZero should return ('i', true), got (%q, %v)", varName, ok)
	}

	// Test: not an assignment
	notAssign := &ast.ExprStmt{X: &ast.Ident{Name: "x"}}
	varName, ok = checkInitIsZero(notAssign)
	// Expected: false
	if ok {
		t.Error("checkInitIsZero should return false for non-assignment")
	}

	// Test: assignment with = instead of :=
	assignNotDefine := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
	}
	varName, ok = checkInitIsZero(assignNotDefine)
	// Expected: false
	if ok {
		t.Error("checkInitIsZero should return false for non-DEFINE assignment")
	}

	// Test: multiple LHS
	multipleLhs := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}, &ast.Ident{Name: "j"}},
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
	}
	varName, ok = checkInitIsZero(multipleLhs)
	// Expected: false
	if ok {
		t.Error("checkInitIsZero should return false for multiple LHS")
	}

	// Test: LHS is not an identifier
	lhsNotIdent := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
	}
	varName, ok = checkInitIsZero(lhsNotIdent)
	// Expected: false
	if ok {
		t.Error("checkInitIsZero should return false when LHS is not an identifier")
	}

	// Test: RHS is not a literal
	rhsNotLit := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "start"}},
	}
	varName, ok = checkInitIsZero(rhsNotLit)
	// Expected: false
	if ok {
		t.Error("checkInitIsZero should return false when RHS is not a literal")
	}

	// Test: RHS is not zero
	rhsNotZero := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "1"}},
	}
	varName, ok = checkInitIsZero(rhsNotZero)
	// Expected: false
	if ok {
		t.Error("checkInitIsZero should return false when RHS is not zero")
	}

	// Test: RHS is not INT
	rhsNotInt := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.FLOAT, Value: "0.0"}},
	}
	varName, ok = checkInitIsZero(rhsNotInt)
	// Expected: false
	if ok {
		t.Error("checkInitIsZero should return false when RHS is not INT")
	}
}

// TestExtractAssignFromInit027 tests extractAssignFromInit027 function.
func TestExtractAssignFromInit027(t *testing.T) {
	// Test: valid assignment
	validAssign := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}},
	}
	result, ok := extractAssignFromInit027(validAssign)
	// Expected: true
	if !ok || result == nil {
		t.Error("extractAssignFromInit027 should return valid assignment")
	}

	// Test: not an assignment
	notAssign := &ast.ExprStmt{X: &ast.Ident{Name: "x"}}
	result, ok = extractAssignFromInit027(notAssign)
	// Expected: false
	if ok {
		t.Error("extractAssignFromInit027 should return false for non-assignment")
	}

	// Test: not DEFINE
	notDefine := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}},
	}
	result, ok = extractAssignFromInit027(notDefine)
	// Expected: false
	if ok {
		t.Error("extractAssignFromInit027 should return false for non-DEFINE")
	}

	// Test: multiple LHS
	multipleLhs := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}, &ast.Ident{Name: "j"}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}, &ast.BasicLit{Value: "1"}},
	}
	result, ok = extractAssignFromInit027(multipleLhs)
	// Expected: false
	if ok {
		t.Error("extractAssignFromInit027 should return false for multiple LHS")
	}

	// Test: multiple RHS
	multipleRhs := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}, &ast.BasicLit{Value: "1"}},
	}
	result, ok = extractAssignFromInit027(multipleRhs)
	// Expected: false
	if ok {
		t.Error("extractAssignFromInit027 should return false for multiple RHS")
	}
}

// TestExtractVarNameFromAssign027 tests extractVarNameFromAssign027 function.
func TestExtractVarNameFromAssign027(t *testing.T) {
	// Test: valid identifier
	validAssign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "i"}},
	}
	result, ok := extractVarNameFromAssign027(validAssign)
	// Expected: true
	if !ok || result != "i" {
		t.Errorf("extractVarNameFromAssign027 should return ('i', true), got (%q, %v)", result, ok)
	}

	// Test: LHS is not an identifier
	notIdent := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
	}
	result, ok = extractVarNameFromAssign027(notIdent)
	// Expected: false
	if ok {
		t.Error("extractVarNameFromAssign027 should return false for non-ident LHS")
	}
}

// TestValidateInitZero027 tests validateInitZero027 function.
func TestValidateInitZero027(t *testing.T) {
	// Test: valid zero
	validZero := &ast.AssignStmt{
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
	}
	result := validateInitZero027(validZero)
	// Expected: true
	if !result {
		t.Error("validateInitZero027 should return true for zero")
	}

	// Test: not a literal
	notLit := &ast.AssignStmt{
		Rhs: []ast.Expr{&ast.Ident{Name: "x"}},
	}
	result = validateInitZero027(notLit)
	// Expected: false
	if result {
		t.Error("validateInitZero027 should return false for non-literal")
	}

	// Test: not zero value
	notZero := &ast.AssignStmt{
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "1"}},
	}
	result = validateInitZero027(notZero)
	// Expected: false
	if result {
		t.Error("validateInitZero027 should return false for non-zero value")
	}

	// Test: not INT kind
	notInt := &ast.AssignStmt{
		Rhs: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: `"0"`}},
	}
	result = validateInitZero027(notInt)
	// Expected: false
	if result {
		t.Error("validateInitZero027 should return false for non-INT kind")
	}
}

// TestCheckPostIsIncrement tests checkPostIsIncrement function.
func TestCheckPostIsIncrement(t *testing.T) {
	// Test: valid increment
	validInc := &ast.IncDecStmt{
		X:   &ast.Ident{Name: "i"},
		Tok: token.INC,
	}
	result := checkPostIsIncrement(validInc, "i")
	// Expected: true
	if !result {
		t.Error("checkPostIsIncrement should return true for valid increment")
	}

	// Test: not an IncDecStmt
	notIncDec := &ast.ExprStmt{X: &ast.Ident{Name: "x"}}
	result = checkPostIsIncrement(notIncDec, "i")
	// Expected: false
	if result {
		t.Error("checkPostIsIncrement should return false for non-IncDecStmt")
	}

	// Test: decrement instead of increment
	decrement := &ast.IncDecStmt{
		X:   &ast.Ident{Name: "i"},
		Tok: token.DEC,
	}
	result = checkPostIsIncrement(decrement, "i")
	// Expected: false
	if result {
		t.Error("checkPostIsIncrement should return false for decrement")
	}

	// Test: X is not an identifier
	notIdent := &ast.IncDecStmt{
		X:   &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
		Tok: token.INC,
	}
	result = checkPostIsIncrement(notIdent, "i")
	// Expected: false
	if result {
		t.Error("checkPostIsIncrement should return false when X is not an ident")
	}

	// Test: wrong variable name
	wrongVar := &ast.IncDecStmt{
		X:   &ast.Ident{Name: "j"},
		Tok: token.INC,
	}
	result = checkPostIsIncrement(wrongVar, "i")
	// Expected: false
	if result {
		t.Error("checkPostIsIncrement should return false for wrong variable name")
	}
}

// TestCheckConditionIsLessThan tests checkConditionIsLessThan function.
func TestCheckConditionIsLessThan(t *testing.T) {
	// Test: valid less than
	validLss := &ast.BinaryExpr{
		X:  &ast.Ident{Name: "i"},
		Op: token.LSS,
		Y:  &ast.Ident{Name: "n"},
	}
	result := checkConditionIsLessThan(validLss, "i")
	// Expected: true
	if !result {
		t.Error("checkConditionIsLessThan should return true for valid LSS")
	}

	// Test: not a binary expression
	notBinary := &ast.Ident{Name: "x"}
	result = checkConditionIsLessThan(notBinary, "i")
	// Expected: false
	if result {
		t.Error("checkConditionIsLessThan should return false for non-binary expr")
	}

	// Test: not LSS operator
	notLss := &ast.BinaryExpr{
		X:  &ast.Ident{Name: "i"},
		Op: token.LEQ,
		Y:  &ast.Ident{Name: "n"},
	}
	result = checkConditionIsLessThan(notLss, "i")
	// Expected: false
	if result {
		t.Error("checkConditionIsLessThan should return false for non-LSS operator")
	}

	// Test: X is not an identifier
	xNotIdent := &ast.BinaryExpr{
		X:  &ast.CallExpr{Fun: &ast.Ident{Name: "getI"}},
		Op: token.LSS,
		Y:  &ast.Ident{Name: "n"},
	}
	result = checkConditionIsLessThan(xNotIdent, "i")
	// Expected: false
	if result {
		t.Error("checkConditionIsLessThan should return false when X is not an ident")
	}

	// Test: wrong variable name
	wrongVar := &ast.BinaryExpr{
		X:  &ast.Ident{Name: "j"},
		Op: token.LSS,
		Y:  &ast.Ident{Name: "n"},
	}
	result = checkConditionIsLessThan(wrongVar, "i")
	// Expected: false
	if result {
		t.Error("checkConditionIsLessThan should return false for wrong variable name")
	}
}
