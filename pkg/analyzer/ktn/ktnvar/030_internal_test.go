package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// TestIsAppendNilPattern tests isAppendNilPattern function.
func TestIsAppendNilPattern(t *testing.T) {
	// Test: not append function
	notAppend := &ast.CallExpr{
		Fun:      &ast.Ident{Name: "copy"},
		Args:     []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Ellipsis: token.Pos(1),
	}
	result := isAppendNilPattern(notAppend)
	// Expected: false
	if result {
		t.Error("isAppendNilPattern should return false for non-append")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "append"},
		},
		Args:     []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Ellipsis: token.Pos(1),
	}
	result = isAppendNilPattern(funNotIdent)
	// Expected: false
	if result {
		t.Error("isAppendNilPattern should return false for non-ident Fun")
	}

	// Test: wrong number of arguments
	wrongArgs := &ast.CallExpr{
		Fun:      &ast.Ident{Name: "append"},
		Args:     []ast.Expr{&ast.Ident{Name: "a"}},
		Ellipsis: token.Pos(1),
	}
	result = isAppendNilPattern(wrongArgs)
	// Expected: false
	if result {
		t.Error("isAppendNilPattern should return false for wrong arg count")
	}

	// Test: no ellipsis
	noEllipsis := &ast.CallExpr{
		Fun:      &ast.Ident{Name: "append"},
		Args:     []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Ellipsis: token.NoPos,
	}
	result = isAppendNilPattern(noEllipsis)
	// Expected: false
	if result {
		t.Error("isAppendNilPattern should return false for no ellipsis")
	}

	// Test: first arg is not nil slice conversion
	notNilConv := &ast.CallExpr{
		Fun:      &ast.Ident{Name: "append"},
		Args:     []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Ellipsis: token.Pos(1),
	}
	result = isAppendNilPattern(notNilConv)
	// Expected: false
	if result {
		t.Error("isAppendNilPattern should return false for non-nil first arg")
	}
}

// TestIsNilSliceConversion tests isNilSliceConversion function.
func TestIsNilSliceConversion(t *testing.T) {
	// Test: not a call expression
	notCall := &ast.Ident{Name: "x"}
	result := isNilSliceConversion(notCall)
	// Expected: false
	if result {
		t.Error("isNilSliceConversion should return false for non-call")
	}

	// Test: wrong number of arguments
	wrongArgs := &ast.CallExpr{
		Fun:  &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
		Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
	}
	result = isNilSliceConversion(wrongArgs)
	// Expected: false
	if result {
		t.Error("isNilSliceConversion should return false for wrong arg count")
	}

	// Test: arg is not nil
	argNotNil := &ast.CallExpr{
		Fun:  &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
		Args: []ast.Expr{&ast.Ident{Name: "x"}},
	}
	result = isNilSliceConversion(argNotNil)
	// Expected: false
	if result {
		t.Error("isNilSliceConversion should return false for non-nil arg")
	}

	// Test: arg is not an ident
	argNotIdent := &ast.CallExpr{
		Fun:  &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
		Args: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: "getNil"}}},
	}
	result = isNilSliceConversion(argNotIdent)
	// Expected: false
	if result {
		t.Error("isNilSliceConversion should return false for non-ident arg")
	}

	// Test: Fun is not an array type
	funNotArray := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "SomeType"},
		Args: []ast.Expr{&ast.Ident{Name: "nil"}},
	}
	result = isNilSliceConversion(funNotArray)
	// Expected: false
	if result {
		t.Error("isNilSliceConversion should return false for non-array Fun")
	}

	// Test: valid nil slice conversion
	validConv := &ast.CallExpr{
		Fun:  &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
		Args: []ast.Expr{&ast.Ident{Name: "nil"}},
	}
	result = isNilSliceConversion(validConv)
	// Expected: true
	if !result {
		t.Error("isNilSliceConversion should return true for valid conversion")
	}
}

// TestExprToString tests exprToString function.
func TestExprToString(t *testing.T) {
	// Test: identifier
	ident := &ast.Ident{Name: "x"}
	result := exprToString(ident)
	// Expected: "x"
	if result != "x" {
		t.Errorf("exprToString should return 'x', got %q", result)
	}

	// Test: selector expression
	selector := &ast.SelectorExpr{
		X:   &ast.Ident{Name: "obj"},
		Sel: &ast.Ident{Name: "field"},
	}
	result = exprToString(selector)
	// Expected: "obj.field"
	if result != "obj.field" {
		t.Errorf("exprToString should return 'obj.field', got %q", result)
	}

	// Test: index expression
	indexExpr := &ast.IndexExpr{
		X:     &ast.Ident{Name: "arr"},
		Index: &ast.Ident{Name: "i"},
	}
	result = exprToString(indexExpr)
	// Expected: "arr[i]"
	if result != "arr[i]" {
		t.Errorf("exprToString should return 'arr[i]', got %q", result)
	}

	// Test: nested selector
	nestedSelector := &ast.SelectorExpr{
		X: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "obj"},
		},
		Sel: &ast.Ident{Name: "field"},
	}
	result = exprToString(nestedSelector)
	// Expected: "pkg.obj.field"
	if result != "pkg.obj.field" {
		t.Errorf("exprToString should return 'pkg.obj.field', got %q", result)
	}

	// Test: unknown expression type
	unknownExpr := &ast.StarExpr{X: &ast.Ident{Name: "ptr"}}
	result = exprToString(unknownExpr)
	// Expected: empty
	if result != "" {
		t.Errorf("exprToString should return empty for unknown type, got %q", result)
	}

	// Test: basic literal
	basicLit := &ast.BasicLit{Value: "0"}
	result = exprToString(basicLit)
	// Expected: empty
	if result != "" {
		t.Errorf("exprToString should return empty for BasicLit, got %q", result)
	}
}

// TestExtractMakeAssign030 tests extractMakeAssign030 function.
func TestExtractMakeAssign030(t *testing.T) {
	// Test: multiple LHS
	multipleLhs := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Rhs: []ast.Expr{&ast.CallExpr{}},
	}
	varName, call := extractMakeAssign030(multipleLhs)
	// Expected: empty, nil
	if varName != "" || call != nil {
		t.Error("extractMakeAssign030 should return empty for multiple LHS")
	}

	// Test: multiple RHS
	multipleRhs := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "a"}},
		Rhs: []ast.Expr{&ast.CallExpr{}, &ast.CallExpr{}},
	}
	varName, call = extractMakeAssign030(multipleRhs)
	// Expected: empty, nil
	if varName != "" || call != nil {
		t.Error("extractMakeAssign030 should return empty for multiple RHS")
	}

	// Test: LHS is not an ident
	lhsNotIdent := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
		Rhs: []ast.Expr{&ast.CallExpr{}},
	}
	varName, call = extractMakeAssign030(lhsNotIdent)
	// Expected: empty, nil
	if varName != "" || call != nil {
		t.Error("extractMakeAssign030 should return empty for non-ident LHS")
	}

	// Test: RHS is not a call
	rhsNotCall := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
	}
	varName, call = extractMakeAssign030(rhsNotCall)
	// Expected: empty, nil
	if varName != "" || call != nil {
		t.Error("extractMakeAssign030 should return empty for non-call RHS")
	}

	// Test: valid assignment
	validAssign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "clone"}},
		Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: "make"}}},
	}
	varName, call = extractMakeAssign030(validAssign)
	// Expected: "clone", non-nil
	if varName != "clone" || call == nil {
		t.Errorf("extractMakeAssign030 should return ('clone', call), got (%q, %v)", varName, call)
	}
}

// TestValidateMakeCall030 tests validateMakeCall030 function.
func TestValidateMakeCall030(t *testing.T) {
	// Test: not make function
	notMake := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "new"},
		Args: []ast.Expr{&ast.ArrayType{}, &ast.Ident{Name: "n"}},
	}
	result := validateMakeCall030(notMake)
	// Expected: empty
	if result != "" {
		t.Error("validateMakeCall030 should return empty for non-make")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "make"},
		},
		Args: []ast.Expr{&ast.ArrayType{}, &ast.Ident{Name: "n"}},
	}
	result = validateMakeCall030(funNotIdent)
	// Expected: empty
	if result != "" {
		t.Error("validateMakeCall030 should return empty for non-ident Fun")
	}

	// Test: wrong number of args
	wrongArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "make"},
		Args: []ast.Expr{&ast.ArrayType{}},
	}
	result = validateMakeCall030(wrongArgs)
	// Expected: empty
	if result != "" {
		t.Error("validateMakeCall030 should return empty for wrong arg count")
	}

	// Test: first arg is not array type
	notArrayType := &ast.CallExpr{
		Fun: &ast.Ident{Name: "make"},
		Args: []ast.Expr{
			&ast.Ident{Name: "map"},
			&ast.CallExpr{
				Fun:  &ast.Ident{Name: "len"},
				Args: []ast.Expr{&ast.Ident{Name: "s"}},
			},
		},
	}
	result = validateMakeCall030(notArrayType)
	// Expected: empty
	if result != "" {
		t.Error("validateMakeCall030 should return empty for non-array type")
	}

	// Test: valid make with len
	validMake := &ast.CallExpr{
		Fun: &ast.Ident{Name: "make"},
		Args: []ast.Expr{
			&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
			&ast.CallExpr{
				Fun:  &ast.Ident{Name: "len"},
				Args: []ast.Expr{&ast.Ident{Name: "source"}},
			},
		},
	}
	result = validateMakeCall030(validMake)
	// Expected: "source"
	if result != "source" {
		t.Errorf("validateMakeCall030 should return 'source', got %q", result)
	}
}

// TestExtractLenSource tests extractLenSource function.
func TestExtractLenSource(t *testing.T) {
	// Test: not a call expression
	notCall := &ast.Ident{Name: "n"}
	result := extractLenSource(notCall)
	// Expected: empty
	if result != "" {
		t.Error("extractLenSource should return empty for non-call")
	}

	// Test: not len function
	notLen := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "cap"},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractLenSource(notLen)
	// Expected: empty
	if result != "" {
		t.Error("extractLenSource should return empty for non-len")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "len"},
		},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractLenSource(funNotIdent)
	// Expected: empty
	if result != "" {
		t.Error("extractLenSource should return empty for non-ident Fun")
	}

	// Test: wrong number of args
	wrongArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "len"},
		Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
	}
	result = extractLenSource(wrongArgs)
	// Expected: empty
	if result != "" {
		t.Error("extractLenSource should return empty for wrong arg count")
	}

	// Test: valid len
	validLen := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "len"},
		Args: []ast.Expr{&ast.Ident{Name: "source"}},
	}
	result = extractLenSource(validLen)
	// Expected: "source"
	if result != "source" {
		t.Errorf("extractLenSource should return 'source', got %q", result)
	}
}

// TestExtractCopyCall030 tests extractCopyCall030 function.
func TestExtractCopyCall030(t *testing.T) {
	// Test: X is not a call expression
	notCall := &ast.ExprStmt{X: &ast.Ident{Name: "x"}}
	call, dest, src := extractCopyCall030(notCall)
	// Expected: nil, "", ""
	if call != nil || dest != "" || src != "" {
		t.Error("extractCopyCall030 should return nil for non-call")
	}

	// Test: not copy function
	notCopy := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "append"},
			Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		},
	}
	call, dest, src = extractCopyCall030(notCopy)
	// Expected: nil, "", ""
	if call != nil || dest != "" || src != "" {
		t.Error("extractCopyCall030 should return nil for non-copy")
	}

	// Test: dest is not extractable (unknown expr type)
	unknownDest := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{Name: "copy"},
			Args: []ast.Expr{
				&ast.StarExpr{X: &ast.Ident{Name: "ptr"}},
				&ast.Ident{Name: "src"},
			},
		},
	}
	call, dest, src = extractCopyCall030(unknownDest)
	// Expected: nil, "", ""
	if call != nil || dest != "" || src != "" {
		t.Error("extractCopyCall030 should return nil for unknown dest")
	}

	// Test: src is not extractable (unknown expr type)
	unknownSrc := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{Name: "copy"},
			Args: []ast.Expr{
				&ast.Ident{Name: "dest"},
				&ast.StarExpr{X: &ast.Ident{Name: "ptr"}},
			},
		},
	}
	call, dest, src = extractCopyCall030(unknownSrc)
	// Expected: nil, "", ""
	if call != nil || dest != "" || src != "" {
		t.Error("extractCopyCall030 should return nil for unknown src")
	}

	// Test: valid copy call
	validCopy := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "copy"},
			Args: []ast.Expr{&ast.Ident{Name: "clone"}, &ast.Ident{Name: "source"}},
		},
	}
	call, dest, src = extractCopyCall030(validCopy)
	// Expected: non-nil, "clone", "source"
	if call == nil || dest != "clone" || src != "source" {
		t.Errorf("extractCopyCall030 should return (call, 'clone', 'source'), got (%v, %q, %q)", call, dest, src)
	}
}

// TestIsCopyCallWithTwoArgs030 tests isCopyCallWithTwoArgs030 function.
func TestIsCopyCallWithTwoArgs030(t *testing.T) {
	// Test: not copy function
	notCopy := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "append"},
		Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
	}
	result := isCopyCallWithTwoArgs030(notCopy)
	// Expected: false
	if result {
		t.Error("isCopyCallWithTwoArgs030 should return false for non-copy")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "copy"},
		},
		Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
	}
	result = isCopyCallWithTwoArgs030(funNotIdent)
	// Expected: false
	if result {
		t.Error("isCopyCallWithTwoArgs030 should return false for non-ident Fun")
	}

	// Test: wrong number of args
	wrongArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "copy"},
		Args: []ast.Expr{&ast.Ident{Name: "a"}},
	}
	result = isCopyCallWithTwoArgs030(wrongArgs)
	// Expected: false
	if result {
		t.Error("isCopyCallWithTwoArgs030 should return false for wrong arg count")
	}

	// Test: valid copy call
	validCopy := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "copy"},
		Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
	}
	result = isCopyCallWithTwoArgs030(validCopy)
	// Expected: true
	if !result {
		t.Error("isCopyCallWithTwoArgs030 should return true for valid copy")
	}
}

// TestFindMatchingMake030 tests findMatchingMake030 function.
func TestFindMatchingMake030(t *testing.T) {
	makeInfos := map[string]*makeCloneInfo{
		"clone": {
			varName:    "clone",
			sourceExpr: "source",
		},
	}

	// Test: no matching make
	result := findMatchingMake030("other", "source", makeInfos)
	// Expected: nil
	if result != nil {
		t.Error("findMatchingMake030 should return nil for non-matching dest")
	}

	// Test: sources don't match
	result = findMatchingMake030("clone", "other", makeInfos)
	// Expected: nil
	if result != nil {
		t.Error("findMatchingMake030 should return nil for non-matching source")
	}

	// Test: matching make
	result = findMatchingMake030("clone", "source", makeInfos)
	// Expected: non-nil
	if result == nil {
		t.Error("findMatchingMake030 should return matching make info")
	}
}

// TestProcessMakeAssignment tests processMakeAssignment function.
func TestProcessMakeAssignment(t *testing.T) {
	makeInfos := make(map[string]*makeCloneInfo)

	// Test: not a valid make assignment
	notValidAssign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
		Rhs: []ast.Expr{&ast.CallExpr{}},
	}
	processMakeAssignment(notValidAssign, makeInfos)
	// Expected: makeInfos is still empty
	if len(makeInfos) != 0 {
		t.Error("processMakeAssignment should not add for invalid assignment")
	}

	// Test: not a valid make call (not make function)
	notMakeCall := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "clone"}},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  &ast.Ident{Name: "new"},
				Args: []ast.Expr{&ast.ArrayType{}},
			},
		},
	}
	processMakeAssignment(notMakeCall, makeInfos)
	// Expected: makeInfos is still empty
	if len(makeInfos) != 0 {
		t.Error("processMakeAssignment should not add for non-make call")
	}

	// Test: valid make assignment
	validAssign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "clone"}},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.CallExpr{
						Fun:  &ast.Ident{Name: "len"},
						Args: []ast.Expr{&ast.Ident{Name: "source"}},
					},
				},
			},
		},
	}
	processMakeAssignment(validAssign, makeInfos)
	// Expected: makeInfos has entry
	if len(makeInfos) != 1 || makeInfos["clone"] == nil {
		t.Error("processMakeAssignment should add valid make info")
	}
}

// Test_runVar030_ruleDisabled tests runVar030 when rule is disabled.
func Test_runVar030_ruleDisabled(t *testing.T) {
	// Save the current config
	oldCfg := config.Get()

	// Create new config with rule disabled
	newCfg := config.DefaultConfig()
	falseVal := false
	newCfg.Rules[ruleCodeVar030] = &config.RuleConfig{Enabled: &falseVal}
	config.Set(newCfg)
	// Ensure restoration at the end
	defer config.Set(oldCfg)

	// Run analyzer with testhelper on bad.go (would have errors if enabled)
	diags := testhelper.RunAnalyzer(t, Analyzer030, "testdata/src/var030/bad.go")

	// With rule disabled, should have 0 errors
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics when rule disabled, got %d", len(diags))
	}
}

// Test_checkAppendNilPattern_fileExcluded tests file exclusion in checkAppendNilPattern.
func Test_checkAppendNilPattern_fileExcluded(t *testing.T) {
	// Source with append([]T(nil), s...) pattern
	src := `package test

func badAppend() []int {
	original := []int{1, 2, 3}
	return append([]int(nil), original...)
}
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded_append.go", src, 0)
	// Check parse error
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Save current config
	oldCfg := config.Get()
	defer config.Set(oldCfg)

	// Create config excluding the file
	newCfg := config.DefaultConfig()
	newCfg.Rules[ruleCodeVar030] = &config.RuleConfig{
		Exclude: []string{"excluded_append.go"},
	}
	config.Set(newCfg)

	// Create mock pass with diagnostics capture
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer030,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Create inspector and add to ResultOf
	insp := inspector.New([]*ast.File{file})
	pass.ResultOf[inspect.Analyzer] = insp

	// Call checkAppendNilPattern directly
	checkAppendNilPattern(pass, insp, newCfg)

	// With file excluded, should have 0 errors
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics when file excluded, got %d", len(diagnostics))
	}
}

// Test_checkMakeCopyPattern_fileExcluded tests file exclusion in checkMakeCopyPattern.
func Test_checkMakeCopyPattern_fileExcluded(t *testing.T) {
	// Source with make+copy pattern
	src := `package test

func badMakeCopy() []int {
	original := []int{1, 2, 3}
	clone := make([]int, len(original))
	copy(clone, original)
	return clone
}
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded_makecopy.go", src, 0)
	// Check parse error
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Save current config
	oldCfg := config.Get()
	defer config.Set(oldCfg)

	// Create config excluding the file
	newCfg := config.DefaultConfig()
	newCfg.Rules[ruleCodeVar030] = &config.RuleConfig{
		Exclude: []string{"excluded_makecopy.go"},
	}
	config.Set(newCfg)

	// Create mock pass with diagnostics capture
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer030,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Create inspector and add to ResultOf
	insp := inspector.New([]*ast.File{file})
	pass.ResultOf[inspect.Analyzer] = insp

	// Call checkMakeCopyPattern directly
	checkMakeCopyPattern(pass, insp, newCfg)

	// With file excluded, should have 0 errors
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics when file excluded, got %d", len(diagnostics))
	}
}

// Test_checkCopyCall_matchingMake tests checkCopyCall with matching make info.
func Test_checkCopyCall_matchingMake(t *testing.T) {
	// Source with make+copy pattern inside a function
	src := `package test

func cloneSlice() {
	original := []int{1, 2, 3}
	clone := make([]int, len(original))
	copy(clone, original)
	_ = clone
}
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test_copy.go", src, 0)
	// Check parse error
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Save current config
	oldCfg := config.Get()
	defer config.Set(oldCfg)

	// Create default config (rule enabled)
	newCfg := config.DefaultConfig()
	config.Set(newCfg)

	// Create mock pass with diagnostics capture
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer030,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Create inspector and add to ResultOf
	insp := inspector.New([]*ast.File{file})
	pass.ResultOf[inspect.Analyzer] = insp

	// Call runVar030 directly
	_, err = runVar030(pass)
	// Check run error
	if err != nil {
		t.Fatalf("runVar030() error = %v", err)
	}

	// Should detect the make+copy pattern
	if len(diagnostics) < 1 {
		t.Errorf("Expected at least 1 diagnostic for make+copy pattern, got %d", len(diagnostics))
	}
}

// Test_checkCopyCall_noMatchingMake tests checkCopyCall when no matching make exists.
func Test_checkCopyCall_noMatchingMake(t *testing.T) {
	// Create makeInfos with a different variable name
	makeInfos := map[string]*makeCloneInfo{
		"differentVar": {
			varName:    "differentVar",
			sourceExpr: "source",
		},
	}

	// Create a copy call that doesn't match the make
	exprStmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "copy"},
			Args: []ast.Expr{&ast.Ident{Name: "clone"}, &ast.Ident{Name: "source"}},
		},
	}

	// Create minimal pass
	fset := token.NewFileSet()
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset: fset,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
	}

	// Create config
	cfg := config.DefaultConfig()

	// Call checkCopyCall - should not report anything
	checkCopyCall(pass, exprStmt, makeInfos, cfg)

	// Should have 0 diagnostics since dest doesn't match
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics when no matching make, got %d", len(diagnostics))
	}

	// makeInfos should still have the entry
	if len(makeInfos) != 1 {
		t.Errorf("makeInfos should still have 1 entry, got %d", len(makeInfos))
	}
}

// Test_checkCopyCall_invalidCall tests checkCopyCall when extractCopyCall030 returns nil.
func Test_checkCopyCall_invalidCall(t *testing.T) {
	// Create makeInfos
	makeInfos := map[string]*makeCloneInfo{
		"clone": {
			varName:    "clone",
			sourceExpr: "source",
		},
	}

	// Create an exprStmt with non-copy call (will cause extractCopyCall030 to return nil)
	exprStmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "append"},
			Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		},
	}

	// Create minimal pass
	fset := token.NewFileSet()
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset: fset,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
	}

	// Create config
	cfg := config.DefaultConfig()

	// Call checkCopyCall - should return early since extractCopyCall030 returns nil
	checkCopyCall(pass, exprStmt, makeInfos, cfg)

	// Should have 0 diagnostics
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics for invalid call, got %d", len(diagnostics))
	}

	// makeInfos should still have the entry (not deleted)
	if len(makeInfos) != 1 {
		t.Errorf("makeInfos should still have 1 entry, got %d", len(makeInfos))
	}
}

// Test_analyzeBlockForMakeCopy tests analyzeBlockForMakeCopy with various statements.
func Test_analyzeBlockForMakeCopy(t *testing.T) {
	// Source with various statement types in a block
	src := `package test

func mixedStatements() {
	original := []int{1, 2, 3}
	clone := make([]int, len(original))
	x := 10
	_ = x
	copy(clone, original)
	y := 20
	_ = y
}
`
	// Parse the source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test_block.go", src, 0)
	// Check parse error
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Save current config
	oldCfg := config.Get()
	defer config.Set(oldCfg)

	// Create default config
	newCfg := config.DefaultConfig()
	config.Set(newCfg)

	// Create mock pass with diagnostics capture
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer: Analyzer030,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
		ResultOf: make(map[*analysis.Analyzer]any),
	}

	// Create inspector and add to ResultOf
	insp := inspector.New([]*ast.File{file})
	pass.ResultOf[inspect.Analyzer] = insp

	// Call runVar030 directly
	_, err = runVar030(pass)
	// Check run error
	if err != nil {
		t.Fatalf("runVar030() error = %v", err)
	}

	// Should detect the make+copy pattern even with intervening statements
	if len(diagnostics) < 1 {
		t.Errorf("Expected at least 1 diagnostic, got %d", len(diagnostics))
	}
}

// Test_checkCopyCall_notCopyCall tests checkCopyCall with non-copy expression statement.
func Test_checkCopyCall_notCopyCall(t *testing.T) {
	// Create an expression statement that is NOT a copy call (e.g., println)
	exprStmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "println"},
			Args: []ast.Expr{&ast.BasicLit{Value: `"hello"`}},
		},
	}

	// Create makeInfos
	makeInfos := map[string]*makeCloneInfo{
		"clone": {
			varName:    "clone",
			sourceExpr: "source",
		},
	}

	// Create minimal pass
	fset := token.NewFileSet()
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset: fset,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
	}

	// Create config
	cfg := config.DefaultConfig()

	// Call checkCopyCall - should return early since it's not a copy call
	checkCopyCall(pass, exprStmt, makeInfos, cfg)

	// Should have 0 diagnostics since it's not a copy call
	if len(diagnostics) != 0 {
		t.Errorf("Expected 0 diagnostics for non-copy call, got %d", len(diagnostics))
	}

	// makeInfos should still have the entry (not deleted)
	if len(makeInfos) != 1 {
		t.Errorf("makeInfos should still have 1 entry, got %d", len(makeInfos))
	}
}

// Test_reportMakeCopy030 tests reportMakeCopy030 function directly.
func Test_reportMakeCopy030(t *testing.T) {
	// Create a simple call expression
	call := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "copy"},
		Args: []ast.Expr{&ast.Ident{Name: "dest"}, &ast.Ident{Name: "src"}},
	}

	// Create minimal pass with diagnostics capture
	fset := token.NewFileSet()
	var diagnostics []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset: fset,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
	}

	// Create config
	cfg := config.DefaultConfig()

	// Call reportMakeCopy030
	reportMakeCopy030(pass, call, cfg)

	// Should have 1 diagnostic
	if len(diagnostics) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(diagnostics))
	}
}
