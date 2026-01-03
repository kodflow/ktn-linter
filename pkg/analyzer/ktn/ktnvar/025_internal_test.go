package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
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
		tt := tt
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
		tt := tt
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

	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "unrecognized expression type",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "foo"}},
			expected: false,
		},
		{
			name:     "zero basic literal",
			expr:     &ast.BasicLit{Value: "0"},
			expected: true,
		},
		{
			name:     "nil identifier",
			expr:     &ast.Ident{Name: "nil"},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Call function
			result := isZeroValue(pass, tt.expr)
			// Check result
			if result != tt.expected {
				t.Errorf("isZeroValue() = %v, want %v", result, tt.expected)
			}
		})
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
		tt := tt
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
		tt := tt
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

// TestIsZeroValueAllBranches tests all branches of isZeroValue.
func TestIsZeroValueAllBranches(t *testing.T) {
	// Create a minimal pass for testing
	pass := &analysis.Pass{}

	// Test: BasicLit branch (zero)
	basicLitZero := &ast.BasicLit{Value: "0"}
	result := isZeroValue(pass, basicLitZero)
	// Expected: true
	if !result {
		t.Error("isZeroValue should return true for BasicLit zero")
	}

	// Test: Ident branch (nil)
	identNil := &ast.Ident{Name: "nil"}
	result = isZeroValue(pass, identNil)
	// Expected: true
	if !result {
		t.Error("isZeroValue should return true for nil identifier")
	}

	// Test: CompositeLit branch (empty)
	compositeLitEmpty := &ast.CompositeLit{Elts: []ast.Expr{}}
	result = isZeroValue(pass, compositeLitEmpty)
	// Expected: true
	if !result {
		t.Error("isZeroValue should return true for empty CompositeLit")
	}

	// Test: CallExpr branch (early return for wrong args)
	callExprMultiArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "string"},
		Args: []ast.Expr{&ast.BasicLit{}, &ast.BasicLit{}},
	}
	result = isZeroValue(pass, callExprMultiArgs)
	// Expected: false
	if result {
		t.Error("isZeroValue should return false for CallExpr with multiple args")
	}
}

// TestIsZeroConversionAllBranches tests all branches of isZeroConversion.
func TestIsZeroConversionAllBranches(t *testing.T) {
	callExpr := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "string"},
		Args: []ast.Expr{&ast.BasicLit{Value: `""`}},
	}

	// Test: pass with nil TypesInfo
	passNilTypesInfo := &analysis.Pass{TypesInfo: nil}
	result := isZeroConversion(passNilTypesInfo, callExpr)
	// Expected: false (TypesInfo is nil)
	if result {
		t.Error("isZeroConversion should return false when TypesInfo is nil")
	}

	// Test: pass with TypesInfo but nil Uses
	passNilUses := &analysis.Pass{TypesInfo: &types.Info{Uses: nil}}
	result = isZeroConversion(passNilUses, callExpr)
	// Expected: false (Uses is nil)
	if result {
		t.Error("isZeroConversion should return false when Uses is nil")
	}

	// Test: pass with TypesInfo and empty Uses (obj is nil)
	passEmptyUses := &analysis.Pass{
		TypesInfo: &types.Info{Uses: make(map[*ast.Ident]types.Object)},
	}
	result = isZeroConversion(passEmptyUses, callExpr)
	// Expected: false (obj is nil - identifier not found in Uses)
	if result {
		t.Error("isZeroConversion should return false when obj is nil")
	}
}

// TestIsZeroConversionNotTypeName tests isZeroConversion when obj is not a TypeName.
func TestIsZeroConversionNotTypeName(t *testing.T) {
	// Create a function identifier
	funIdent := &ast.Ident{Name: "myFunc"}
	callExpr := &ast.CallExpr{
		Fun:  funIdent,
		Args: []ast.Expr{&ast.BasicLit{Value: "0"}},
	}

	// Create a Uses map with a function, not a type
	pkg := types.NewPackage("test", "test")
	funcObj := types.NewFunc(0, pkg, "myFunc", types.NewSignatureType(
		nil, nil, nil, nil, nil, false,
	))

	usesMap := make(map[*ast.Ident]types.Object)
	usesMap[funIdent] = funcObj

	passWithFunc := &analysis.Pass{
		TypesInfo: &types.Info{Uses: usesMap},
	}

	result := isZeroConversion(passWithFunc, callExpr)
	// Expected: false (obj is Func, not TypeName)
	if result {
		t.Error("isZeroConversion should return false when obj is not a TypeName")
	}
}

// TestIsZeroConversionWithTypeName tests isZeroConversion with a valid TypeName.
func TestIsZeroConversionWithTypeName(t *testing.T) {
	// Create a type identifier
	typeIdent := &ast.Ident{Name: "MyType"}
	callExpr := &ast.CallExpr{
		Fun:  typeIdent,
		Args: []ast.Expr{&ast.BasicLit{Value: "0"}},
	}

	// Create a Uses map with a type
	pkg := types.NewPackage("test", "test")
	typeObj := types.NewTypeName(0, pkg, "MyType", types.Typ[types.Int])

	usesMap := make(map[*ast.Ident]types.Object)
	usesMap[typeIdent] = typeObj

	passWithType := &analysis.Pass{
		TypesInfo: &types.Info{Uses: usesMap},
	}

	result := isZeroConversion(passWithType, callExpr)
	// Expected: true (type conversion with zero value)
	if !result {
		t.Error("isZeroConversion should return true for valid type conversion with zero")
	}
}

// TestRunVar025RuleDisabled tests runVar025 when rule is disabled.
func TestRunVar025RuleDisabled(t *testing.T) {
	// Save current config
	oldCfg := config.Get()
	defer config.Set(oldCfg)

	// Disable the rule
	newCfg := config.DefaultConfig()
	newCfg.Rules[ruleCodeVar025] = &config.RuleConfig{Enabled: config.Bool(false)}
	config.Set(newCfg)

	// Parse code with clear pattern
	src := `package test

func clearMap() {
	m := make(map[string]int)
	for k := range m {
		delete(m, k)
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Analyzer: Analyzer025,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(_ analysis.Diagnostic) {
			reportCount++
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	result, err := runVar025(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar025() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when rule is disabled
	if reportCount != 0 {
		t.Errorf("expected 0 reports when rule disabled, got %d", reportCount)
	}
}

// TestRunVar025NilTypesInfo tests runVar025 when TypesInfo is nil.
func TestRunVar025NilTypesInfo(t *testing.T) {
	// Reset config to enable rule
	config.Reset()
	defer config.Reset()

	// Parse code with clear pattern
	src := `package test

func clearSlice() {
	s := []int{1, 2, 3}
	for i := range s {
		s[i] = 0
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Analyzer:  Analyzer025,
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: nil, // nil TypesInfo
		Report: func(_ analysis.Diagnostic) {
			reportCount++
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	result, err := runVar025(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar025() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when TypesInfo is nil (early return)
	if reportCount != 0 {
		t.Errorf("expected 0 reports when TypesInfo is nil, got %d", reportCount)
	}
}

// TestRunVar025FileExcluded tests runVar025 when file is excluded.
func TestRunVar025FileExcluded(t *testing.T) {
	// Save current config
	oldCfg := config.Get()
	defer config.Set(oldCfg)

	// Configure rule with file exclusion
	newCfg := config.DefaultConfig()
	newCfg.Rules[ruleCodeVar025] = &config.RuleConfig{
		Enabled: config.Bool(true),
		Exclude: []string{"excluded.go"},
	}
	config.Set(newCfg)

	// Parse code with clear pattern
	src := `package test

func clearMap() {
	m := make(map[string]int)
	for k := range m {
		delete(m, k)
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", src, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Analyzer:  Analyzer025,
		Fset:      fset,
		Files:     []*ast.File{file},
		TypesInfo: &types.Info{},
		Report: func(_ analysis.Diagnostic) {
			reportCount++
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspector.New([]*ast.File{file}),
		},
	}

	result, err := runVar025(pass)
	// Check result
	if err != nil || result != nil {
		t.Errorf("runVar025() = (%v, %v), want (nil, nil)", result, err)
	}
	// Check no reports when file is excluded
	if reportCount != 0 {
		t.Errorf("expected 0 reports when file excluded, got %d", reportCount)
	}
}
