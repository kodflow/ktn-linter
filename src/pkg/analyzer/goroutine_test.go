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

// TestGoroutineAnalyzer teste l'analyseur goroutine (KTN-GOROUTINE-001 et KTN-GOROUTINE-002).
//
// Params:
//   - t: l'instance de test
func TestGoroutineAnalyzer(t *testing.T) {
	t.Run("KTN-GOROUTINE-001", testGoroutineKTNGOROUTINE001)
	t.Run("KTN-GOROUTINE-002", testGoroutineKTNGOROUTINE002)
}

// testGoroutineKTNGOROUTINE001 teste les goroutines dans une boucle sans limitation.
//
// Params:
//   - t: instance de test
func testGoroutineKTNGOROUTINE001(t *testing.T) {
	// ❌ MAUVAIS: go dans boucle sans sync
	runGoroutineTest(t, "go in loop without sync", `package test
func badLoop() {
	items := []int{1, 2, 3}
	for _, item := range items {
		go process(item)
	}
}
func process(i int) {}
`, true, "KTN-GOROUTINE-001")

	// ✅ BON: go avec WaitGroup
	runGoroutineTest(t, "go with waitgroup", `package test
import "sync"
func goodLoop() {
	var wg sync.WaitGroup
	items := []int{1, 2, 3}
	for _, item := range items {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			process(i)
		}(item)
	}
	wg.Wait()
}
func process(i int) {}
`, false, "")

	// ✅ BON: go hors boucle avec synchronisation
	runGoroutineTest(t, "go outside loop", `package test
import "sync"
func goodNoLoop() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		process(42)
	}()
	wg.Wait()
}
func process(i int) {}
`, false, "")
}

// testGoroutineKTNGOROUTINE002 teste les goroutines sans mécanisme de synchronisation.
//
// Params:
//   - t: instance de test
func testGoroutineKTNGOROUTINE002(t *testing.T) {
	// ❌ MAUVAIS: go sans synchronisation
	runGoroutineTest(t, "go without sync", `package test
func badNoSync() {
	go func() {
		doWork()
	}()
}
func doWork() {}
`, true, "KTN-GOROUTINE-002")

	// ✅ BON: go avec channel
	runGoroutineTest(t, "go with channel", `package test
func goodChannel() {
	ch := make(chan int)
	go func() {
		ch <- 42
	}()
	<-ch
}
`, false, "")

	// ✅ BON: go avec context
	runGoroutineTest(t, "go with context", `package test
import "context"
func goodContext(ctx context.Context) {
	go func() {
		<-ctx.Done()
	}()
}
`, false, "")
}

// runGoroutineTest exécute un test pour GoroutineAnalyzer.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - wantDiag: true si on attend un diagnostic
//   - wantMsg: message attendu dans le diagnostic
func runGoroutineTest(t *testing.T, name, code string, wantDiag bool, wantMsg string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		var diagnostics []analysis.Diagnostic
		pass := &analysis.Pass{
			Analyzer: analyzer.GoroutineAnalyzer,
			Fset:     fset,
			Files:    []*ast.File{file},
			Report: func(diag analysis.Diagnostic) {
				diagnostics = append(diagnostics, diag)
			},
		}

		_, err = analyzer.GoroutineAnalyzer.Run(pass)
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
