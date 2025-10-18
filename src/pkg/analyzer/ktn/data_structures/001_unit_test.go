package ktn_data_structures_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_data_structures "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/data_structures"
	"golang.org/x/tools/go/analysis"
)

// Tests unitaires pour getArraySize

// TestGetArraySize_NilLen tests getArraySize with nil Len.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestGetArraySize_NilLen(t *testing.T) {
	// Note: getArraySize is not exported, so we test it indirectly through RuleArray001
	// This test is kept for documentation but will be refactored to test the public API
	t.Skip("getArraySize is not exported - test indirectly through RuleArray001")
}

// TestGetArraySize_NotBasicLit tests getArraySize with non-BasicLit Len.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestGetArraySize_NotBasicLit(t *testing.T) {
	t.Skip("getArraySize is not exported - test indirectly through RuleArray001")
}

// TestGetArraySize_InvalidFormat tests getArraySize with invalid format.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestGetArraySize_InvalidFormat(t *testing.T) {
	t.Skip("getArraySize is not exported - test indirectly through RuleArray001")
}

// TestGetArraySize_ValidSize tests getArraySize with valid size.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestGetArraySize_ValidSize(t *testing.T) {
	t.Skip("getArraySize is not exported - test indirectly through RuleArray001")
}

// Tests unitaires pour runRuleArray001 - cas edge

// TestRunRuleArray001_NotCompositeLit tests runRuleArray001 with non-CompositeLit.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleArray001_NotCompositeLit(t *testing.T) {
	src := `
package test
func test() {
	x := 42
}
`
	testRunRule(t, src, ktn_data_structures.RuleArray001, 0)
}

// TestRunRuleArray001_NotArrayType tests runRuleArray001 with non-ArrayType.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleArray001_NotArrayType(t *testing.T) {
	src := `
package test
func test() {
	s := []int{1, 2, 3}
}
`
	testRunRule(t, src, ktn_data_structures.RuleArray001, 0)
}

// TestRunRuleArray001_InvalidArraySize tests runRuleArray001 with invalid array size.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleArray001_InvalidArraySize(t *testing.T) {
	// Cas où getArraySize retourne -1
	src := `
package test
const N = 5
func test() {
	arr := [N]int{1, 2, 3}
}
`
	testRunRule(t, src, ktn_data_structures.RuleArray001, 0)
}

// TestRunRuleArray001_ValidArray tests runRuleArray001 with valid array.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleArray001_ValidArray(t *testing.T) {
	src := `
package test
func test() {
	arr := [3]int{1, 2, 3}
}
`
	testRunRule(t, src, ktn_data_structures.RuleArray001, 0)
}

// TestRunRuleArray001_TooManyElements tests runRuleArray001 with too many elements.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleArray001_TooManyElements(t *testing.T) {
	src := `
package test
func test() {
	arr := [2]int{1, 2, 3}
}
`
	testRunRule(t, src, ktn_data_structures.RuleArray001, 1)
}

// Tests unitaires pour runRuleMap001 - cas edge

// TestRunRuleMap001_NotAssignStmt tests runRuleMap001 with non-AssignStmt.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleMap001_NotAssignStmt(t *testing.T) {
	src := `
package test
func test() {
	x := 42
}
`
	testRunRule(t, src, ktn_data_structures.RuleMap001, 0)
}

// TestRunRuleMap001_LhsNotIndexExpr tests runRuleMap001 with non-IndexExpr LHS.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleMap001_LhsNotIndexExpr(t *testing.T) {
	src := `
package test
func test() {
	x := 42
	x = 10
}
`
	testRunRule(t, src, ktn_data_structures.RuleMap001, 0)
}

// TestRunRuleMap001_IndexExprXNotIdent tests runRuleMap001 with IndexExpr X not Ident.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleMap001_IndexExprXNotIdent(t *testing.T) {
	// Cas où indexExpr.X n'est pas un Ident simple
	src := `
package test
func getMap() map[string]int {
	return make(map[string]int)
}
func test() {
	getMap()["key"] = 42
}
`
	testRunRule(t, src, ktn_data_structures.RuleMap001, 0)
}

// TestRunRuleMap001_SafeMap tests runRuleMap001 with safe map.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleMap001_SafeMap(t *testing.T) {
	src := `
package test
func test() {
	m := make(map[string]int)
	m["key"] = 42
}
`
	testRunRule(t, src, ktn_data_structures.RuleMap001, 0)
}

// TestRunRuleMap001_UnsafeMap tests runRuleMap001 with unsafe map.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleMap001_UnsafeMap(t *testing.T) {
	src := `
package test
func test() {
	var m map[string]int
	m["key"] = 42
}
`
	testRunRule(t, src, ktn_data_structures.RuleMap001, 1)
}

// Tests unitaires pour runRuleSlice001 - cas edge

// TestRunRuleSlice001_NotIndexExpr tests runRuleSlice001 with non-IndexExpr.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleSlice001_NotIndexExpr(t *testing.T) {
	src := `
package test
func test() {
	x := 42
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 0)
}

// TestRunRuleSlice001_BasicLitIndex tests runRuleSlice001 with BasicLit index.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleSlice001_BasicLitIndex(t *testing.T) {
	src := `
package test
func test() {
	s := []int{1, 2, 3}
	_ = s[0]
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 0)
}

// TestRunRuleSlice001_IndexNotIdent tests runRuleSlice001 with index not Ident.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleSlice001_IndexNotIdent(t *testing.T) {
	// Index est une expression complexe
	src := `
package test
func test() {
	s := []int{1, 2, 3}
	_ = s[1+1]
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 0)
}

// TestRunRuleSlice001_IndexFromRange tests runRuleSlice001 with index from range.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleSlice001_IndexFromRange(t *testing.T) {
	src := `
package test
func test() {
	s := []int{1, 2, 3}
	for i := range s {
		_ = s[i]
	}
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 0)
}

// TestRunRuleSlice001_UncheckedIndex tests runRuleSlice001 with unchecked index.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleSlice001_UncheckedIndex(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	_ = items[i]
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 1)
}

// TestRunRuleSlice001_CheckedIndex tests runRuleSlice001 with checked index.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestRunRuleSlice001_CheckedIndex(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	if i < len(items) {
		_ = items[i]
	}
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 0)
}

// Tests unitaires pour isIndexChecked - cas edge

// TestIsIndexChecked_NotIfStmt tests isIndexChecked with non-IfStmt.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestIsIndexChecked_NotIfStmt(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	_ = items[i]
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 1)
}

// TestIsIndexChecked_NotBinaryExpr tests isIndexChecked with non-BinaryExpr.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestIsIndexChecked_NotBinaryExpr(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	if true {
		_ = items[i]
	}
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 1)
}

// TestIsIndexChecked_WrongOperator tests isIndexChecked with wrong operator.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestIsIndexChecked_WrongOperator(t *testing.T) {
	// Opérateur différent de <
	src := `
package test
func test(items []int, i int) {
	if i > len(items) {
		_ = items[i]
	}
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 1)
}

// TestIsIndexChecked_IdentMismatch tests isIndexChecked with ident mismatch.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestIsIndexChecked_IdentMismatch(t *testing.T) {
	// L'identifiant vérifié n'est pas le bon
	src := `
package test
func test(items []int, i, j int) {
	if j < len(items) {
		_ = items[i]
	}
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 1)
}

// TestIsIndexChecked_NotLenCall tests isIndexChecked with non-len call.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestIsIndexChecked_NotLenCall(t *testing.T) {
	// La fonction appelée n'est pas len
	src := `
package test
func getSize(s []int) int { return len(s) }
func test(items []int, i int) {
	if i < getSize(items) {
		_ = items[i]
	}
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 1)
}

// TestIsIndexChecked_CorrectCheck tests isIndexChecked with correct check.
//
// Params:
//   - t: testing instance
//
// nolint:KTN-FUNC-001
func TestIsIndexChecked_CorrectCheck(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	if i < len(items) {
		_ = items[i]
	}
}
`
	testRunRule(t, src, ktn_data_structures.RuleSlice001, 0)
}

// Fonction utilitaire pour tester les règles

// testRunRule runs a rule analyzer on source code and checks diagnostic count.
//
// Params:
//   - t: testing instance
//   - src: source code to analyze
//   - analyzer: the analyzer to run
//   - expectedDiags: expected number of diagnostics
func testRunRule(t *testing.T, src string, analyzer *analysis.Analyzer, expectedDiags int) {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		Report: func(d analysis.Diagnostic) {
			// Count diagnostics
		},
		ResultOf: make(map[*analysis.Analyzer]interface{}),
	}

	// Count diagnostics
	diagCount := 0
	originalReport := pass.Report
	pass.Report = func(d analysis.Diagnostic) {
		diagCount++
		originalReport(d)
	}

	_, err = analyzer.Run(pass)
	if err != nil {
		t.Fatalf("Analyzer failed: %v", err)
	}

	if diagCount != expectedDiags {
		t.Errorf("Expected %d diagnostics, got %d", expectedDiags, diagCount)
	}
}
