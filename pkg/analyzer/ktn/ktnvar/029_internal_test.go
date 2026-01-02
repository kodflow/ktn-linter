package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper"
	"github.com/kodflow/ktn-linter/pkg/config"
)

// TestExtractSliceFromCapLenCondition tests extractSliceFromCapLenCondition function.
func TestExtractSliceFromCapLenCondition(t *testing.T) {
	// Test: not a binary expression
	notBinary := &ast.Ident{Name: "x"}
	result := extractSliceFromCapLenCondition(notBinary)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapLenCondition should return empty for non-binary")
	}

	// Test: not LSS operator
	notLss := &ast.BinaryExpr{
		X:  &ast.Ident{Name: "x"},
		Op: token.GTR,
		Y:  &ast.Ident{Name: "n"},
	}
	result = extractSliceFromCapLenCondition(notLss)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapLenCondition should return empty for non-LSS")
	}
}

// TestExtractSliceFromCapMinusLen tests extractSliceFromCapMinusLen function.
func TestExtractSliceFromCapMinusLen(t *testing.T) {
	// Test: not a binary expression
	notBinary := &ast.Ident{Name: "x"}
	result := extractSliceFromCapMinusLen(notBinary)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapMinusLen should return empty for non-binary")
	}

	// Test: not SUB operator
	notSub := &ast.BinaryExpr{
		X:  &ast.Ident{Name: "x"},
		Op: token.ADD,
		Y:  &ast.Ident{Name: "y"},
	}
	result = extractSliceFromCapMinusLen(notSub)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapMinusLen should return empty for non-SUB")
	}

	// Test: cap returns empty
	noCapResult := &ast.BinaryExpr{
		X:  &ast.Ident{Name: "x"}, // Not a cap call
		Op: token.SUB,
		Y: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "len"},
			Args: []ast.Expr{&ast.Ident{Name: "s"}},
		},
	}
	result = extractSliceFromCapMinusLen(noCapResult)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapMinusLen should return empty when cap is invalid")
	}

	// Test: len returns empty
	noLenResult := &ast.BinaryExpr{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "cap"},
			Args: []ast.Expr{&ast.Ident{Name: "s"}},
		},
		Op: token.SUB,
		Y:  &ast.Ident{Name: "y"}, // Not a len call
	}
	result = extractSliceFromCapMinusLen(noLenResult)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapMinusLen should return empty when len is invalid")
	}

	// Test: slices don't match
	slicesDontMatch := &ast.BinaryExpr{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "cap"},
			Args: []ast.Expr{&ast.Ident{Name: "s1"}},
		},
		Op: token.SUB,
		Y: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "len"},
			Args: []ast.Expr{&ast.Ident{Name: "s2"}},
		},
	}
	result = extractSliceFromCapMinusLen(slicesDontMatch)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapMinusLen should return empty when slices don't match")
	}

	// Test: valid cap(s) - len(s)
	validExpr := &ast.BinaryExpr{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "cap"},
			Args: []ast.Expr{&ast.Ident{Name: "s"}},
		},
		Op: token.SUB,
		Y: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "len"},
			Args: []ast.Expr{&ast.Ident{Name: "s"}},
		},
	}
	result = extractSliceFromCapMinusLen(validExpr)
	// Expected: "s"
	if result != "s" {
		t.Errorf("extractSliceFromCapMinusLen should return 's', got %q", result)
	}
}

// TestExtractSliceFromCapCall tests extractSliceFromCapCall function.
func TestExtractSliceFromCapCall(t *testing.T) {
	// Test: not a call expression
	notCall := &ast.Ident{Name: "x"}
	result := extractSliceFromCapCall(notCall)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapCall should return empty for non-call")
	}

	// Test: not cap function
	notCap := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "len"},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractSliceFromCapCall(notCap)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapCall should return empty for non-cap function")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "cap"},
		},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractSliceFromCapCall(funNotIdent)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapCall should return empty for non-ident Fun")
	}

	// Test: wrong number of arguments
	wrongArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "cap"},
		Args: []ast.Expr{&ast.Ident{Name: "s1"}, &ast.Ident{Name: "s2"}},
	}
	result = extractSliceFromCapCall(wrongArgs)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromCapCall should return empty for wrong arg count")
	}

	// Test: valid cap(s)
	validCap := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "cap"},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractSliceFromCapCall(validCap)
	// Expected: "s"
	if result != "s" {
		t.Errorf("extractSliceFromCapCall should return 's', got %q", result)
	}
}

// TestExtractSliceFromLenCall tests extractSliceFromLenCall function.
func TestExtractSliceFromLenCall(t *testing.T) {
	// Test: not a call expression
	notCall := &ast.Ident{Name: "x"}
	result := extractSliceFromLenCall(notCall)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromLenCall should return empty for non-call")
	}

	// Test: not len function
	notLen := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "cap"},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractSliceFromLenCall(notLen)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromLenCall should return empty for non-len function")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "len"},
		},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractSliceFromLenCall(funNotIdent)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromLenCall should return empty for non-ident Fun")
	}

	// Test: wrong number of arguments
	wrongArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "len"},
		Args: []ast.Expr{&ast.Ident{Name: "s1"}, &ast.Ident{Name: "s2"}},
	}
	result = extractSliceFromLenCall(wrongArgs)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceFromLenCall should return empty for wrong arg count")
	}

	// Test: valid len(s)
	validLen := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "len"},
		Args: []ast.Expr{&ast.Ident{Name: "s"}},
	}
	result = extractSliceFromLenCall(validLen)
	// Expected: "s"
	if result != "s" {
		t.Errorf("extractSliceFromLenCall should return 's', got %q", result)
	}
}

// TestExtractSliceName tests extractSliceName function.
func TestExtractSliceName(t *testing.T) {
	// Test: valid identifier
	validIdent := &ast.Ident{Name: "s"}
	result := extractSliceName(validIdent)
	// Expected: "s"
	if result != "s" {
		t.Errorf("extractSliceName should return 's', got %q", result)
	}

	// Test: not an identifier
	notIdent := &ast.CallExpr{Fun: &ast.Ident{Name: "getSlice"}}
	result = extractSliceName(notIdent)
	// Expected: empty
	if result != "" {
		t.Error("extractSliceName should return empty for non-ident")
	}
}

// TestHasGrowPatternInBody tests hasGrowPatternInBody function.
func TestHasGrowPatternInBody(t *testing.T) {
	// Test: nil body
	result := hasGrowPatternInBody(nil, "s")
	// Expected: false
	if result {
		t.Error("hasGrowPatternInBody should return false for nil body")
	}

	// Test: empty body
	emptyBody := &ast.BlockStmt{List: nil}
	result = hasGrowPatternInBody(emptyBody, "s")
	// Expected: false
	if result {
		t.Error("hasGrowPatternInBody should return false for empty body")
	}

	// Test: incomplete pattern (only make)
	makeOnlyBody := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{&ast.Ident{Name: "newS"}},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{Name: "make"},
						Args: []ast.Expr{
							&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
							&ast.Ident{Name: "n"},
						},
					},
				},
			},
		},
	}
	result = hasGrowPatternInBody(makeOnlyBody, "s")
	// Expected: false
	if result {
		t.Error("hasGrowPatternInBody should return false for incomplete pattern")
	}
}

// TestAnalyzeGrowStatement tests analyzeGrowStatement function.
func TestAnalyzeGrowStatement(t *testing.T) {
	info := &growPatternInfo{sliceName: "s"}

	// Test: other statement types (not Assign or Expr)
	returnStmt := &ast.ReturnStmt{}
	analyzeGrowStatement(returnStmt, info)
	// No panic expected
}

// TestAnalyzeAssignForGrow tests analyzeAssignForGrow function.
func TestAnalyzeAssignForGrow(t *testing.T) {
	info := &growPatternInfo{sliceName: "s"}

	// Test: multiple RHS
	multipleRhs := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "x"}, &ast.Ident{Name: "y"}},
	}
	analyzeAssignForGrow(multipleRhs, info)

	// Test: RHS is make
	makeAssign := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "newS"}},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.Ident{Name: "make"},
				Args: []ast.Expr{
					&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
					&ast.Ident{Name: "n"},
				},
			},
		},
	}
	info2 := &growPatternInfo{sliceName: "s"}
	analyzeAssignForGrow(makeAssign, info2)
	// Expected: hasMake = true
	if !info2.hasMake {
		t.Error("analyzeAssignForGrow should set hasMake for make call")
	}
}

// TestIsMakeSliceCall tests isMakeSliceCall function.
func TestIsMakeSliceCall(t *testing.T) {
	// Test: not a call expression
	notCall := &ast.Ident{Name: "x"}
	result := isMakeSliceCall(notCall)
	// Expected: false
	if result {
		t.Error("isMakeSliceCall should return false for non-call")
	}

	// Test: not make function
	notMake := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "new"},
		Args: []ast.Expr{&ast.Ident{Name: "int"}},
	}
	result = isMakeSliceCall(notMake)
	// Expected: false
	if result {
		t.Error("isMakeSliceCall should return false for non-make function")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "pkg"},
			Sel: &ast.Ident{Name: "make"},
		},
		Args: []ast.Expr{&ast.ArrayType{Elt: &ast.Ident{Name: "int"}}},
	}
	result = isMakeSliceCall(funNotIdent)
	// Expected: false
	if result {
		t.Error("isMakeSliceCall should return false for non-ident Fun")
	}

	// Test: no arguments
	noArgs := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "make"},
		Args: nil,
	}
	result = isMakeSliceCall(noArgs)
	// Expected: false
	if result {
		t.Error("isMakeSliceCall should return false for no args")
	}

	// Test: first arg is not array type
	notArrayType := &ast.CallExpr{
		Fun:  &ast.Ident{Name: "make"},
		Args: []ast.Expr{&ast.Ident{Name: "map"}},
	}
	result = isMakeSliceCall(notArrayType)
	// Expected: false
	if result {
		t.Error("isMakeSliceCall should return false for non-array type arg")
	}

	// Test: valid make slice
	validMake := &ast.CallExpr{
		Fun: &ast.Ident{Name: "make"},
		Args: []ast.Expr{
			&ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
			&ast.Ident{Name: "n"},
		},
	}
	result = isMakeSliceCall(validMake)
	// Expected: true
	if !result {
		t.Error("isMakeSliceCall should return true for valid make slice")
	}
}

// TestCheckSliceReassign tests checkSliceReassign function.
func TestCheckSliceReassign(t *testing.T) {
	info := &growPatternInfo{sliceName: "s"}

	// Test: multiple LHS
	multipleLhs := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "newA"}, &ast.Ident{Name: "newB"}},
	}
	checkSliceReassign(multipleLhs, info)
	// Expected: hasReassign = false
	if info.hasReassign {
		t.Error("checkSliceReassign should not set hasReassign for multiple LHS")
	}

	// Test: LHS is not an ident
	lhsNotIdent := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
		Rhs: []ast.Expr{&ast.Ident{Name: "newArr"}},
	}
	info2 := &growPatternInfo{sliceName: "s"}
	checkSliceReassign(lhsNotIdent, info2)
	// Expected: hasReassign = false
	if info2.hasReassign {
		t.Error("checkSliceReassign should not set hasReassign for non-ident LHS")
	}

	// Test: target doesn't match slice name
	wrongTarget := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "other"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "newS"}},
	}
	info3 := &growPatternInfo{sliceName: "s"}
	checkSliceReassign(wrongTarget, info3)
	// Expected: hasReassign = false
	if info3.hasReassign {
		t.Error("checkSliceReassign should not set hasReassign for wrong target")
	}

	// Test: DEFINE instead of ASSIGN
	defineStmt := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{&ast.Ident{Name: "s"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "newS"}},
	}
	info4 := &growPatternInfo{sliceName: "s"}
	checkSliceReassign(defineStmt, info4)
	// Expected: hasReassign = false
	if info4.hasReassign {
		t.Error("checkSliceReassign should not set hasReassign for DEFINE")
	}

	// Test: valid reassignment
	validReassign := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{&ast.Ident{Name: "s"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "newS"}},
	}
	info5 := &growPatternInfo{sliceName: "s"}
	checkSliceReassign(validReassign, info5)
	// Expected: hasReassign = true
	if !info5.hasReassign {
		t.Error("checkSliceReassign should set hasReassign for valid reassignment")
	}
}

// TestAnalyzeCopyForGrow tests analyzeCopyForGrow function.
func TestAnalyzeCopyForGrow(t *testing.T) {
	info := &growPatternInfo{sliceName: "s"}

	// Test: X is not a call expression
	notCall := &ast.ExprStmt{X: &ast.Ident{Name: "x"}}
	analyzeCopyForGrow(notCall, info)
	// Expected: hasCopy = false
	if info.hasCopy {
		t.Error("analyzeCopyForGrow should not set hasCopy for non-call")
	}

	// Test: not copy function
	notCopy := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "append"},
			Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		},
	}
	info2 := &growPatternInfo{sliceName: "s"}
	analyzeCopyForGrow(notCopy, info2)
	// Expected: hasCopy = false
	if info2.hasCopy {
		t.Error("analyzeCopyForGrow should not set hasCopy for non-copy function")
	}

	// Test: Fun is not an ident
	funNotIdent := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "pkg"},
				Sel: &ast.Ident{Name: "copy"},
			},
			Args: []ast.Expr{&ast.Ident{Name: "a"}, &ast.Ident{Name: "b"}},
		},
	}
	info3 := &growPatternInfo{sliceName: "s"}
	analyzeCopyForGrow(funNotIdent, info3)
	// Expected: hasCopy = false
	if info3.hasCopy {
		t.Error("analyzeCopyForGrow should not set hasCopy for non-ident Fun")
	}

	// Test: wrong number of arguments
	wrongArgs := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "copy"},
			Args: []ast.Expr{&ast.Ident{Name: "a"}},
		},
	}
	info4 := &growPatternInfo{sliceName: "s"}
	analyzeCopyForGrow(wrongArgs, info4)
	// Expected: hasCopy = false
	if info4.hasCopy {
		t.Error("analyzeCopyForGrow should not set hasCopy for wrong arg count")
	}

	// Test: source doesn't match slice name
	wrongSource := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "copy"},
			Args: []ast.Expr{&ast.Ident{Name: "newS"}, &ast.Ident{Name: "other"}},
		},
	}
	info5 := &growPatternInfo{sliceName: "s"}
	analyzeCopyForGrow(wrongSource, info5)
	// Expected: hasCopy = false
	if info5.hasCopy {
		t.Error("analyzeCopyForGrow should not set hasCopy for wrong source")
	}

	// Test: source is not an ident
	sourceNotIdent := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{Name: "copy"},
			Args: []ast.Expr{
				&ast.Ident{Name: "newS"},
				&ast.CallExpr{Fun: &ast.Ident{Name: "getS"}},
			},
		},
	}
	info6 := &growPatternInfo{sliceName: "s"}
	analyzeCopyForGrow(sourceNotIdent, info6)
	// Expected: hasCopy = false
	if info6.hasCopy {
		t.Error("analyzeCopyForGrow should not set hasCopy for non-ident source")
	}

	// Test: valid copy
	validCopy := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  &ast.Ident{Name: "copy"},
			Args: []ast.Expr{&ast.Ident{Name: "newS"}, &ast.Ident{Name: "s"}},
		},
	}
	info7 := &growPatternInfo{sliceName: "s"}
	analyzeCopyForGrow(validCopy, info7)
	// Expected: hasCopy = true
	if !info7.hasCopy {
		t.Error("analyzeCopyForGrow should set hasCopy for valid copy")
	}
}

// TestRunVar029_RuleDisabled tests runVar029 when rule is disabled.
func TestRunVar029_RuleDisabled(t *testing.T) {
	// Save original config
	originalCfg := config.Get()
	defer config.Set(originalCfg)

	// Create config with rule disabled
	falseVal := false
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			ruleCodeVar029: {
				Enabled: &falseVal,
			},
		},
	}
	config.Set(cfg)

	// Run analyzer - should have 0 diagnostics when disabled
	diags := testhelper.RunAnalyzer(t, Analyzer029, "testdata/src/var029/bad.go")

	// Verify no diagnostics when rule disabled
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics when rule disabled, got %d", len(diags))
	}
}

// TestRunVar029_FileExcluded tests runVar029 when file is excluded.
func TestRunVar029_FileExcluded(t *testing.T) {
	// Save original config
	originalCfg := config.Get()
	defer config.Set(originalCfg)

	// Create config with file exclusion pattern
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			ruleCodeVar029: {
				Exclude: []string{"**/bad.go"},
			},
		},
	}
	config.Set(cfg)

	// Run analyzer - should have 0 diagnostics when file excluded
	diags := testhelper.RunAnalyzer(t, Analyzer029, "testdata/src/var029/bad.go")

	// Verify no diagnostics when file excluded
	if len(diags) != 0 {
		t.Errorf("Expected 0 diagnostics when file excluded, got %d", len(diags))
	}
}
