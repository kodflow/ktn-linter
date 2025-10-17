package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestAllocAnalyzer teste l'analyseur d'allocation (KTN-ALLOC-*).
//
// Params:
//   - t: l'instance de test
func TestAllocAnalyzer(t *testing.T) {
	t.Run("KTN-ALLOC-001", testAllocKTNALLOC001)
	t.Run("KTN-ALLOC-002", testAllocKTNALLOC002)
	t.Run("KTN-ALLOC-004", testAllocKTNALLOC004)
}

// testAllocKTNALLOC001 teste l'interdiction de new() avec types référence.
//
// Params:
//   - t: instance de test
func testAllocKTNALLOC001(t *testing.T) {
	// ❌ MAUVAIS: new() avec slice
	runAllocTest(t, "new with slice", `package test
func badSlice() *[]int {
	return new([]int)
}
`, true, "KTN-ALLOC-001")

	// ❌ MAUVAIS: new() avec map
	runAllocTest(t, "new with map", `package test
func badMap() *map[string]int {
	return new(map[string]int)
}
`, true, "KTN-ALLOC-001")

	// ❌ MAUVAIS: new() avec chan
	runAllocTest(t, "new with chan", `package test
func badChan() *chan int {
	return new(chan int)
}
`, true, "KTN-ALLOC-001")

	// ✅ BON: make() avec slice
	runAllocTest(t, "make with slice", `package test
func goodSlice() []int {
	return make([]int, 0)
}
`, false, "")
}

// testAllocKTNALLOC002 teste le pattern make([]T, 0) suivi d'append.
//
// Params:
//   - t: instance de test
func testAllocKTNALLOC002(t *testing.T) {
	// ❌ MAUVAIS: make([]T, 0) puis append
	runAllocTest(t, "make zero then append", `package test
func badMakeAppend() []int {
	items := make([]int, 0)
	items = append(items, 1)
	return items
}
`, true, "KTN-ALLOC-002")

	// ✅ BON: make avec capacité
	runAllocTest(t, "make with capacity", `package test
func goodMakeAppend() []int {
	items := make([]int, 0, 10)
	items = append(items, 1)
	return items
}
`, false, "")

	// ✅ BON: make sans append
	runAllocTest(t, "make without append", `package test
func goodMake() []int {
	items := make([]int, 0)
	return items
}
`, false, "")
}

// testAllocKTNALLOC004 teste la préférence pour &struct{} au lieu de new(struct).
//
// Params:
//   - t: instance de test
func testAllocKTNALLOC004(t *testing.T) {
	// ❌ MAUVAIS: new(MyStruct)
	runAllocTest(t, "new with struct", `package test
type User struct {
	Name string
}
func badNew() *User {
	return new(User)
}
`, true, "KTN-ALLOC-004")

	// ✅ BON: &MyStruct{}
	runAllocTest(t, "composite literal", `package test
type User struct {
	Name string
}
func goodLiteral() *User {
	return &User{}
}
`, false, "")
}

// runAllocTest exécute un test pour AllocAnalyzer.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - wantDiag: true si on attend un diagnostic
//   - wantMsg: message attendu dans le diagnostic
func runAllocTest(t *testing.T, name, code string, wantDiag bool, wantMsg string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		var diagnostics []analysis.Diagnostic
		pass := &analysis.Pass{
			Analyzer: analyzer.AllocAnalyzer,
			Fset:     fset,
			Files:    []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {
				diagnostics = append(diagnostics, diag)
			},
		}

		_, err = analyzer.AllocAnalyzer.Run(pass)
		if err != nil {
			t.Fatalf("analyzer failed: %v", err)
		}

		hasExpectedDiag := false
		for _, d := range diagnostics {
			if wantMsg != "" && strings.Contains(d.Message, wantMsg) {
				hasExpectedDiag = true
				break
			}
		}

		if wantDiag && !hasExpectedDiag {
			t.Errorf("expected diagnostic %q but got none. Diagnostics: %v", wantMsg, diagnostics)
		}
		if !wantDiag && len(diagnostics) > 0 {
			t.Errorf("expected no diagnostic but got: %v", diagnostics)
		}
	})
}
