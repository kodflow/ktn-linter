package analyzer_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestPoolAnalyzer teste l'analyseur pool (KTN-POOL-001 et KTN-STRUCT-004).
//
// Params:
//   - t: l'instance de test
func TestPoolAnalyzer(t *testing.T) {
	t.Run("KTN-POOL-001", testPoolKTNPOOL001)
	t.Run("KTN-STRUCT-004", testPoolKTNSTRUCT004)
}

// testPoolKTNPOOL001 teste pool.Get() sans defer pool.Put().
//
// Params:
//   - t: instance de test
func testPoolKTNPOOL001(t *testing.T) {
	// ❌ MAUVAIS: pool.Get() sans defer Put()
	runPoolTest(t, "get without defer put", `package test
import "sync"
var bufferPool = sync.Pool{}
func badPool() {
	buf := bufferPool.Get()
	process(buf)
}
func process(v interface{}) {}
`, true, "KTN-POOL-001")

	// ✅ BON: pool.Get() avec defer Put()
	runPoolTest(t, "get with defer put", `package test
import "sync"
var bufferPool = sync.Pool{}
func goodPool() {
	buf := bufferPool.Get()
	defer bufferPool.Put(buf)
	process(buf)
}
func process(v interface{}) {}
`, false, "")

	// ✅ BON: pas de pool
	runPoolTest(t, "no pool", `package test
func noPool() {
	x := 42
	process(x)
}
func process(v interface{}) {}
`, false, "")
}

// testPoolKTNSTRUCT004 teste les grandes structs passées par valeur.
//
// Params:
//   - t: instance de test
func testPoolKTNSTRUCT004(t *testing.T) {
	// ❌ MAUVAIS: grande struct par valeur
	runPoolTest(t, "large struct by value", `package test
type LargeStruct struct {
	A, B, C, D, E, F, G, H int64
	I, J, K, L, M, N, O, P int64
	Q, R, S, T, U, V, W, X int64
}
func badLargeStruct(data LargeStruct) {
	process(data)
}
func process(v interface{}) {}
`, true, "KTN-STRUCT-004")

	// ✅ BON: grande struct par pointeur
	runPoolTest(t, "large struct by pointer", `package test
type LargeStruct struct {
	A, B, C, D, E, F, G, H int64
	I, J, K, L, M, N, O, P int64
	Q, R, S, T, U, V, W, X int64
}
func goodLargeStruct(data *LargeStruct) {
	process(data)
}
func process(v interface{}) {}
`, false, "")

	// ✅ BON: petite struct par valeur
	runPoolTest(t, "small struct by value", `package test
type SmallStruct struct {
	X int
}
func goodSmallStruct(data SmallStruct) {
	process(data)
}
func process(v interface{}) {}
`, false, "")
}

// runPoolTest exécute un test pour PoolAnalyzer.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - wantDiag: true si on attend un diagnostic
//   - wantMsg: message attendu dans le diagnostic
func runPoolTest(t *testing.T, name, code string, wantDiag bool, wantMsg string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		// Type check le code
		conf := types.Config{Importer: importer.Default()}
		info := &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Uses:  make(map[*ast.Ident]types.Object),
			Defs:  make(map[*ast.Ident]types.Object),
		}
		_, err = conf.Check("test", fset, []*ast.File{file}, info)
		if err != nil {
			// Ignorer les erreurs de type check pour les tests simples
		}

		var diagnostics []analysis.Diagnostic
		pass := &analysis.Pass{
			Analyzer:  analyzer.PoolAnalyzer,
			Fset:      fset,
			Files:     []*ast.File{file},
			TypesInfo: info,
			Report: func(diag analysis.Diagnostic) {
				diagnostics = append(diagnostics, diag)
			},
		}

		_, err = analyzer.PoolAnalyzer.Run(pass)
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
