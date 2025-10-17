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

// TestErrorAnalyzer teste l'analyseur error (KTN-ERROR-001).
//
// Params:
//   - t: l'instance de test
func TestErrorAnalyzer(t *testing.T) {
	// ❌ MAUVAIS: return err sans wrapping
	runErrorTest(t, "unwrapped error", `package test
import "errors"
func badError() error {
	err := errors.New("oops")
	if err != nil {
		return err
	}
	return nil
}
`, true, "KTN-ERROR-001")

	// ✅ BON: return nil
	runErrorTest(t, "return nil", `package test
func goodNil() error {
	return nil
}
`, false, "")

	// ✅ BON: return avec fmt.Errorf
	runErrorTest(t, "wrapped error", `package test
import "fmt"
func goodError() error {
	err := doSomething()
	if err != nil {
		return fmt.Errorf("context: %w", err)
	}
	return nil
}
func doSomething() error {
	return nil
}
`, false, "")

	// ✅ BON: fonction sans error return
	runErrorTest(t, "no error return", `package test
func noError() int {
	return 42
}
`, false, "")
}

// runErrorTest exécute un test pour ErrorAnalyzer.
//
// Params:
//   - t: instance de test
//   - name: nom du test
//   - code: code source à tester
//   - wantDiag: true si on attend un diagnostic
//   - wantMsg: message attendu dans le diagnostic
func runErrorTest(t *testing.T, name, code string, wantDiag bool, wantMsg string) {
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
			Analyzer:  analyzer.ErrorAnalyzer,
			Fset:      fset,
			Files:     []*ast.File{file},
			TypesInfo: info,
			Report: func(diag analysis.Diagnostic) {
				diagnostics = append(diagnostics, diag)
			},
		}

		_, err = analyzer.ErrorAnalyzer.Run(pass)
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
