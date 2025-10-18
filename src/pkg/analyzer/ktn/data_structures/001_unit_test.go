package ktn_data_structures

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// Tests unitaires pour getArraySize

func TestGetArraySize_NilLen(t *testing.T) {
	arrayType := &ast.ArrayType{
		Len: nil,
	}
	size := getArraySize(arrayType)
	if size != -1 {
		t.Errorf("Expected -1 for nil Len, got %d", size)
	}
}

func TestGetArraySize_NotBasicLit(t *testing.T) {
	// Cas où Len n'est pas un BasicLit (ex: une expression)
	arrayType := &ast.ArrayType{
		Len: &ast.Ident{Name: "N"},
	}
	size := getArraySize(arrayType)
	if size != -1 {
		t.Errorf("Expected -1 for non-BasicLit, got %d", size)
	}
}

func TestGetArraySize_InvalidFormat(t *testing.T) {
	// Cas où le BasicLit n'est pas un nombre valide
	arrayType := &ast.ArrayType{
		Len: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"abc"`,
		},
	}
	size := getArraySize(arrayType)
	if size != -1 {
		t.Errorf("Expected -1 for invalid format, got %d", size)
	}
}

func TestGetArraySize_ValidSize(t *testing.T) {
	arrayType := &ast.ArrayType{
		Len: &ast.BasicLit{
			Kind:  token.INT,
			Value: "10",
		},
	}
	size := getArraySize(arrayType)
	if size != 10 {
		t.Errorf("Expected 10, got %d", size)
	}
}

// Tests unitaires pour runRuleArray001 - cas edge

func TestRunRuleArray001_NotCompositeLit(t *testing.T) {
	src := `
package test
func test() {
	x := 42
}
`
	testRunRule(t, src, RuleArray001, 0)
}

func TestRunRuleArray001_NotArrayType(t *testing.T) {
	src := `
package test
func test() {
	s := []int{1, 2, 3}
}
`
	testRunRule(t, src, RuleArray001, 0)
}

func TestRunRuleArray001_InvalidArraySize(t *testing.T) {
	// Cas où getArraySize retourne -1
	src := `
package test
const N = 5
func test() {
	arr := [N]int{1, 2, 3}
}
`
	testRunRule(t, src, RuleArray001, 0)
}

func TestRunRuleArray001_ValidArray(t *testing.T) {
	src := `
package test
func test() {
	arr := [3]int{1, 2, 3}
}
`
	testRunRule(t, src, RuleArray001, 0)
}

func TestRunRuleArray001_TooManyElements(t *testing.T) {
	src := `
package test
func test() {
	arr := [2]int{1, 2, 3}
}
`
	testRunRule(t, src, RuleArray001, 1)
}

// Tests unitaires pour runRuleMap001 - cas edge

func TestRunRuleMap001_NotAssignStmt(t *testing.T) {
	src := `
package test
func test() {
	x := 42
}
`
	testRunRule(t, src, RuleMap001, 0)
}

func TestRunRuleMap001_LhsNotIndexExpr(t *testing.T) {
	src := `
package test
func test() {
	x := 42
	x = 10
}
`
	testRunRule(t, src, RuleMap001, 0)
}

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
	testRunRule(t, src, RuleMap001, 0)
}

func TestRunRuleMap001_SafeMap(t *testing.T) {
	src := `
package test
func test() {
	m := make(map[string]int)
	m["key"] = 42
}
`
	testRunRule(t, src, RuleMap001, 0)
}

func TestRunRuleMap001_UnsafeMap(t *testing.T) {
	src := `
package test
func test() {
	var m map[string]int
	m["key"] = 42
}
`
	testRunRule(t, src, RuleMap001, 1)
}

// Tests unitaires pour runRuleSlice001 - cas edge

func TestRunRuleSlice001_NotIndexExpr(t *testing.T) {
	src := `
package test
func test() {
	x := 42
}
`
	testRunRule(t, src, RuleSlice001, 0)
}

func TestRunRuleSlice001_BasicLitIndex(t *testing.T) {
	src := `
package test
func test() {
	s := []int{1, 2, 3}
	_ = s[0]
}
`
	testRunRule(t, src, RuleSlice001, 0)
}

func TestRunRuleSlice001_IndexNotIdent(t *testing.T) {
	// Index est une expression complexe
	src := `
package test
func test() {
	s := []int{1, 2, 3}
	_ = s[1+1]
}
`
	testRunRule(t, src, RuleSlice001, 0)
}

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
	testRunRule(t, src, RuleSlice001, 0)
}

func TestRunRuleSlice001_UncheckedIndex(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	_ = items[i]
}
`
	testRunRule(t, src, RuleSlice001, 1)
}

func TestRunRuleSlice001_CheckedIndex(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	if i < len(items) {
		_ = items[i]
	}
}
`
	testRunRule(t, src, RuleSlice001, 0)
}

// Tests unitaires pour isIndexChecked - cas edge

func TestIsIndexChecked_NotIfStmt(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	_ = items[i]
}
`
	testRunRule(t, src, RuleSlice001, 1)
}

func TestIsIndexChecked_NotBinaryExpr(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	if true {
		_ = items[i]
	}
}
`
	testRunRule(t, src, RuleSlice001, 1)
}

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
	testRunRule(t, src, RuleSlice001, 1)
}

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
	testRunRule(t, src, RuleSlice001, 1)
}

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
	testRunRule(t, src, RuleSlice001, 1)
}

func TestIsIndexChecked_CorrectCheck(t *testing.T) {
	src := `
package test
func test(items []int, i int) {
	if i < len(items) {
		_ = items[i]
	}
}
`
	testRunRule(t, src, RuleSlice001, 0)
}

// Fonction utilitaire pour tester les règles

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
