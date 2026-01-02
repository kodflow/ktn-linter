package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// TestIsCompositeLitZero025 tests the isCompositeLitZero025 function.
func TestIsCompositeLitZero025(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		lit      *ast.CompositeLit
		expected bool
	}{
		{
			name:     "empty composite literal",
			lit:      &ast.CompositeLit{Elts: nil},
			expected: true,
		},
		{
			name:     "empty slice literal",
			lit:      &ast.CompositeLit{Elts: []ast.Expr{}},
			expected: true,
		},
		{
			name: "non-empty composite literal",
			lit: &ast.CompositeLit{
				Elts: []ast.Expr{
					&ast.BasicLit{Value: "1"},
				},
			},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isCompositeLitZero025(tt.lit)
			// Check result
			if result != tt.expected {
				t.Errorf("isCompositeLitZero025() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsZeroConversion tests the isZeroConversion function.
func TestIsZeroConversion(t *testing.T) {
	// Test cases that don't require a valid pass
	tests := []struct {
		name     string
		callExpr *ast.CallExpr
		expected bool
	}{
		{
			name: "multiple args",
			callExpr: &ast.CallExpr{
				Fun:  &ast.Ident{Name: "string"},
				Args: []ast.Expr{&ast.BasicLit{}, &ast.BasicLit{}},
			},
			expected: false,
		},
		{
			name: "non-ident fun",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "foo"},
					Sel: &ast.Ident{Name: "bar"},
				},
				Args: []ast.Expr{&ast.BasicLit{}},
			},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These tests return early before accessing pass
			result := isZeroConversion(nil, tt.callExpr)
			// Check result
			if result != tt.expected {
				t.Errorf("isZeroConversion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsZeroValue tests the isZeroValue function.
func TestIsZeroValue(t *testing.T) {
	// Create a minimal pass for testing
	pass := &analysis.Pass{}

	// Test: unrecognized expression type
	unrecognized := &ast.StarExpr{X: &ast.Ident{Name: "foo"}}
	result := isZeroValue(pass, unrecognized)
	// Expected: false
	if result {
		t.Error("isZeroValue with unrecognized type should return false")
	}
}

// TestIsBasicLitZero025 tests the isBasicLitZero025 function.
func TestIsBasicLitZero025(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		lit      *ast.BasicLit
		expected bool
	}{
		{
			name:     "integer zero",
			lit:      &ast.BasicLit{Value: "0"},
			expected: true,
		},
		{
			name:     "float zero",
			lit:      &ast.BasicLit{Value: "0.0"},
			expected: true,
		},
		{
			name:     "empty string double quotes",
			lit:      &ast.BasicLit{Value: `""`},
			expected: true,
		},
		{
			name:     "empty string backticks",
			lit:      &ast.BasicLit{Value: "``"},
			expected: true,
		},
		{
			name:     "non-zero integer",
			lit:      &ast.BasicLit{Value: "1"},
			expected: false,
		},
		{
			name:     "non-zero float",
			lit:      &ast.BasicLit{Value: "1.5"},
			expected: false,
		},
		{
			name:     "non-empty string",
			lit:      &ast.BasicLit{Value: `"hello"`},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isBasicLitZero025(tt.lit)
			// Check result
			if result != tt.expected {
				t.Errorf("isBasicLitZero025() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsIdentZero025 tests the isIdentZero025 function.
func TestIsIdentZero025(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		ident    *ast.Ident
		expected bool
	}{
		{
			name:     "nil",
			ident:    &ast.Ident{Name: "nil"},
			expected: true,
		},
		{
			name:     "false",
			ident:    &ast.Ident{Name: "false"},
			expected: true,
		},
		{
			name:     "true",
			ident:    &ast.Ident{Name: "true"},
			expected: false,
		},
		{
			name:     "other identifier",
			ident:    &ast.Ident{Name: "x"},
			expected: false,
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isIdentZero025(tt.ident)
			// Check result
			if result != tt.expected {
				t.Errorf("isIdentZero025() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetRangeCollectionIdent tests the getRangeCollectionIdent function.
func TestGetRangeCollectionIdent(t *testing.T) {
	// Test: identifier
	ident := &ast.Ident{Name: "slice"}
	result := getRangeCollectionIdent(ident)
	// Expected: return the identifier
	if result == nil || result.Name != "slice" {
		t.Error("getRangeCollectionIdent should return identifier")
	}

	// Test: parenthesized expression
	parenExpr := &ast.ParenExpr{X: &ast.Ident{Name: "arr"}}
	result = getRangeCollectionIdent(parenExpr)
	// Expected: return the inner identifier
	if result == nil || result.Name != "arr" {
		t.Error("getRangeCollectionIdent should unwrap parenthesized expression")
	}

	// Test: nested parenthesized expression
	nestedParen := &ast.ParenExpr{
		X: &ast.ParenExpr{X: &ast.Ident{Name: "deep"}},
	}
	result = getRangeCollectionIdent(nestedParen)
	// Expected: return the innermost identifier
	if result == nil || result.Name != "deep" {
		t.Error("getRangeCollectionIdent should unwrap nested parenthesized expressions")
	}

	// Test: other expression type
	starExpr := &ast.StarExpr{X: &ast.Ident{Name: "ptr"}}
	result = getRangeCollectionIdent(starExpr)
	// Expected: return nil
	if result != nil {
		t.Error("getRangeCollectionIdent should return nil for other expression types")
	}
}

// TestCheckClearPattern tests the checkClearPattern function.
func TestCheckClearPattern(t *testing.T) {
	// Test: nil body
	rangeStmt := &ast.RangeStmt{
		Body: nil,
	}
	checkClearPattern(nil, rangeStmt)

	// Test: empty body
	rangeStmt2 := &ast.RangeStmt{
		Body: &ast.BlockStmt{List: nil},
	}
	checkClearPattern(nil, rangeStmt2)

	// Test: body with more than one statement
	rangeStmt3 := &ast.RangeStmt{
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "a"}},
				&ast.ExprStmt{X: &ast.Ident{Name: "b"}},
			},
		},
	}
	checkClearPattern(nil, rangeStmt3)
}

// TestCheckMapDeletePattern tests the checkMapDeletePattern function.
func TestCheckMapDeletePattern(t *testing.T) {
	// Test: no key in range
	rangeStmt := &ast.RangeStmt{
		Key:  nil,
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}}},
	}
	result := checkMapDeletePattern(nil, rangeStmt)
	// Expected: false
	if result {
		t.Error("checkMapDeletePattern should return false when key is nil")
	}

	// Test: key is not an identifier
	rangeStmt2 := &ast.RangeStmt{
		Key:  &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}}},
	}
	result = checkMapDeletePattern(nil, rangeStmt2)
	// Expected: false
	if result {
		t.Error("checkMapDeletePattern should return false when key is not an identifier")
	}

	// Test: range expression is not an identifier
	rangeStmt3 := &ast.RangeStmt{
		Key:  &ast.Ident{Name: "k"},
		X:    &ast.CallExpr{Fun: &ast.Ident{Name: "getMap"}},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}}},
	}
	result = checkMapDeletePattern(nil, rangeStmt3)
	// Expected: false
	if result {
		t.Error("checkMapDeletePattern should return false when X is not an identifier")
	}

	// Test: body statement is not an ExprStmt
	rangeStmt4 := &ast.RangeStmt{
		Key:  &ast.Ident{Name: "k"},
		X:    &ast.Ident{Name: "m"},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
	}
	result = checkMapDeletePattern(nil, rangeStmt4)
	// Expected: false
	if result {
		t.Error("checkMapDeletePattern should return false when body is not ExprStmt")
	}

	// Test: ExprStmt.X is not a CallExpr
	rangeStmt5 := &ast.RangeStmt{
		Key: &ast.Ident{Name: "k"},
		X:   &ast.Ident{Name: "m"},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Ident{Name: "something"}},
			},
		},
	}
	result = checkMapDeletePattern(nil, rangeStmt5)
	// Expected: false
	if result {
		t.Error("checkMapDeletePattern should return false when body is not CallExpr")
	}
}

// TestIsDeleteCallWithKeyAndMap tests isDeleteCallWithKeyAndMap function.
func TestIsDeleteCallWithKeyAndMap(t *testing.T) {
	keyIdent := &ast.Ident{Name: "k"}
	mapIdent := &ast.Ident{Name: "m"}

	// Test: not a delete call (different function name)
	callNotDelete := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "remove"},
		Args: []ast.Expr{&ast.Ident{Name: "m"}, &ast.Ident{Name: "k"}},
	}
	result := isDeleteCallWithKeyAndMap(callNotDelete, keyIdent, mapIdent, nil)
	// Expected: false
	if result {
		t.Error("isDeleteCallWithKeyAndMap should return false for non-delete function")
	}

	// Test: delete with wrong number of arguments
	callWrongArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "delete"},
		Args: []ast.Expr{&ast.Ident{Name: "m"}},
	}
	result = isDeleteCallWithKeyAndMap(callWrongArgs, keyIdent, mapIdent, nil)
	// Expected: false
	if result {
		t.Error("isDeleteCallWithKeyAndMap should return false for wrong arg count")
	}

	// Test: delete with wrong map argument
	callWrongMap := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "delete"},
		Args: []ast.Expr{&ast.Ident{Name: "other"}, &ast.Ident{Name: "k"}},
	}
	result = isDeleteCallWithKeyAndMap(callWrongMap, keyIdent, mapIdent, nil)
	// Expected: false
	if result {
		t.Error("isDeleteCallWithKeyAndMap should return false for wrong map")
	}

	// Test: delete with wrong key argument
	callWrongKey := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "delete"},
		Args: []ast.Expr{&ast.Ident{Name: "m"}, &ast.Ident{Name: "other"}},
	}
	result = isDeleteCallWithKeyAndMap(callWrongKey, keyIdent, mapIdent, nil)
	// Expected: false
	if result {
		t.Error("isDeleteCallWithKeyAndMap should return false for wrong key")
	}

	// Test: delete with non-ident second arg
	callNonIdentKey := &ast.CallExpr{
		Fun: &ast.Ident{Name: "delete"},
		Args: []ast.Expr{
			&ast.Ident{Name: "m"},
			&ast.CallExpr{Fun: &ast.Ident{Name: "getKey"}},
		},
	}
	result = isDeleteCallWithKeyAndMap(callNonIdentKey, keyIdent, mapIdent, nil)
	// Expected: false
	if result {
		t.Error("isDeleteCallWithKeyAndMap should return false for non-ident key arg")
	}

	// Test: Fun is not an ident
	callNonIdentFun := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "maps"},
			Sel: &ast.Ident{Name: "Delete"},
		},
		Args: []ast.Expr{&ast.Ident{Name: "m"}, &ast.Ident{Name: "k"}},
	}
	result = isDeleteCallWithKeyAndMap(callNonIdentFun, keyIdent, mapIdent, nil)
	// Expected: false
	if result {
		t.Error("isDeleteCallWithKeyAndMap should return false for non-ident fun")
	}

	// Test: first arg is not an ident
	callNonIdentMapArg := &ast.CallExpr{
		Fun: &ast.Ident{Name: "delete"},
		Args: []ast.Expr{
			&ast.CallExpr{Fun: &ast.Ident{Name: "getMap"}},
			&ast.Ident{Name: "k"},
		},
	}
	result = isDeleteCallWithKeyAndMap(callNonIdentMapArg, keyIdent, mapIdent, nil)
	// Expected: false
	if result {
		t.Error("isDeleteCallWithKeyAndMap should return false for non-ident map arg")
	}
}

// TestCheckSliceZeroPattern tests checkSliceZeroPattern function.
func TestCheckSliceZeroPattern(t *testing.T) {
	// Test: no key in range
	rangeStmt := &ast.RangeStmt{
		Key:  nil,
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.AssignStmt{}}},
	}
	checkSliceZeroPattern(nil, rangeStmt)

	// Test: key is not an identifier
	rangeStmt2 := &ast.RangeStmt{
		Key:  &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.AssignStmt{}}},
	}
	checkSliceZeroPattern(nil, rangeStmt2)

	// Test: X is not an identifier
	rangeStmt3 := &ast.RangeStmt{
		Key:  &ast.Ident{Name: "i"},
		X:    &ast.CallExpr{Fun: &ast.Ident{Name: "getSlice"}},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.AssignStmt{}}},
	}
	checkSliceZeroPattern(nil, rangeStmt3)

	// Test: body statement is not an AssignStmt
	rangeStmt4 := &ast.RangeStmt{
		Key:  &ast.Ident{Name: "i"},
		X:    &ast.Ident{Name: "s"},
		Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}}},
	}
	checkSliceZeroPattern(nil, rangeStmt4)

	// Test: assign token is not =
	rangeStmt5 := &ast.RangeStmt{
		Key: &ast.Ident{Name: "i"},
		X:   &ast.Ident{Name: "s"},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{Tok: token.DEFINE},
			},
		},
	}
	checkSliceZeroPattern(nil, rangeStmt5)
}

// TestCheckIndexAssignZero tests checkIndexAssignZero function.
func TestCheckIndexAssignZero(t *testing.T) {
	indexIdent := &ast.Ident{Name: "i"}
	sliceIdent := &ast.Ident{Name: "s"}

	// Test: multiple LHS
	assignMultipleLhs := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}},
	}
	checkIndexAssignZero(nil, assignMultipleLhs, indexIdent, sliceIdent)

	// Test: multiple RHS
	assignMultipleRhs := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "a"}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}, &ast.BasicLit{Value: "1"}},
	}
	checkIndexAssignZero(nil, assignMultipleRhs, indexIdent, sliceIdent)

	// Test: LHS is not an IndexExpr
	assignNotIndex := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}},
	}
	checkIndexAssignZero(nil, assignNotIndex, indexIdent, sliceIdent)

	// Test: RHS is not a zero value
	assignNonZero := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.IndexExpr{
				X:     &ast.Ident{Name: "s"},
				Index: &ast.Ident{Name: "i"},
			},
		},
		Rhs: []ast.Expr{&ast.BasicLit{Value: "1"}},
	}
	checkIndexAssignZero(nil, assignNonZero, indexIdent, sliceIdent)
}

// TestIsMatchingSliceIndex tests isMatchingSliceIndex function.
func TestIsMatchingSliceIndex(t *testing.T) {
	indexIdent := &ast.Ident{Name: "i"}
	sliceIdent := &ast.Ident{Name: "s"}

	// Test: slice doesn't match
	indexExprWrongSlice := &ast.IndexExpr{
		X:     &ast.Ident{Name: "other"},
		Index: &ast.Ident{Name: "i"},
	}
	result := isMatchingSliceIndex(indexExprWrongSlice, indexIdent, sliceIdent)
	// Expected: false
	if result {
		t.Error("isMatchingSliceIndex should return false for wrong slice")
	}

	// Test: X is not an identifier
	indexExprNonIdentX := &ast.IndexExpr{
		X:     &ast.CallExpr{Fun: &ast.Ident{Name: "getSlice"}},
		Index: &ast.Ident{Name: "i"},
	}
	result = isMatchingSliceIndex(indexExprNonIdentX, indexIdent, sliceIdent)
	// Expected: false
	if result {
		t.Error("isMatchingSliceIndex should return false for non-ident X")
	}

	// Test: index doesn't match
	indexExprWrongIndex := &ast.IndexExpr{
		X:     &ast.Ident{Name: "s"},
		Index: &ast.Ident{Name: "j"},
	}
	result = isMatchingSliceIndex(indexExprWrongIndex, indexIdent, sliceIdent)
	// Expected: false
	if result {
		t.Error("isMatchingSliceIndex should return false for wrong index")
	}

	// Test: index is not an identifier
	indexExprNonIdentIdx := &ast.IndexExpr{
		X:     &ast.Ident{Name: "s"},
		Index: &ast.BasicLit{Value: "0"},
	}
	result = isMatchingSliceIndex(indexExprNonIdentIdx, indexIdent, sliceIdent)
	// Expected: false
	if result {
		t.Error("isMatchingSliceIndex should return false for non-ident index")
	}

	// Test: valid match
	indexExprValid := &ast.IndexExpr{
		X:     &ast.Ident{Name: "s"},
		Index: &ast.Ident{Name: "i"},
	}
	result = isMatchingSliceIndex(indexExprValid, indexIdent, sliceIdent)
	// Expected: true
	if !result {
		t.Error("isMatchingSliceIndex should return true for valid match")
	}
}
