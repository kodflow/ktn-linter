package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// testAnalyzerConfig contient la configuration d'un test pour TestAnalyzer.
type testAnalyzerConfig struct {
	name     string
	fileName string
	code     string
	wantErr  bool
	wantMsg  string
}

// runTestAnalyzerTest exécute un test pour le TestAnalyzer.
//
// Params:
//   - t: instance de test
//   - cfg: configuration du test
func runTestAnalyzerTest(t *testing.T, cfg testAnalyzerConfig) {
	t.Helper()
	t.Run(cfg.name, func(t *testing.T) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, cfg.fileName, cfg.code, parser.ParseComments)
		if err != nil {
			t.Fatalf("Failed to parse code: %v", err)
		}

		var diagnostics []analysis.Diagnostic
		pkg := &types.Package{}
		pkg.SetName(file.Name.Name)

		pass := &analysis.Pass{
			Analyzer: analyzer.TestAnalyzer,
			Fset:     fset,
			Files:    []*ast.File{file},
			Pkg:      pkg,
			Report: func(diag analysis.Diagnostic) {
				diagnostics = append(diagnostics, diag)
			},
		}

		_, err = analyzer.TestAnalyzer.Run(pass)
		if err != nil {
			t.Fatalf("Analyzer returned error: %v", err)
		}

		foundExpected := false
		for _, d := range diagnostics {
			if cfg.wantMsg != "" && contains(d.Message, cfg.wantMsg) {
				foundExpected = true
				break
			}
		}

		if cfg.wantErr && !foundExpected {
			t.Errorf("Expected error containing %q, but got: %v", cfg.wantMsg, diagnostics)
		}
		if !cfg.wantErr && len(diagnostics) > 0 {
			t.Errorf("Expected no errors, but got: %v", diagnostics)
		}
	})
}

// TestTestAnalyzerExemptedPackage teste que le package main est exempté.
//
// Params:
//   - t: instance de test
func TestTestAnalyzerExemptedPackage(t *testing.T) {
	runTestAnalyzerTest(t, testAnalyzerConfig{
		name:     "main package exempted",
		fileName: "main.go",
		code: `package main

func main() {
}
`,
		wantErr: false,
	})
}

// TestTestAnalyzerKTNTEST001 teste la vérification des noms de package de test.
//
// Params:
//   - t: instance de test
func TestTestAnalyzerKTNTEST001(t *testing.T) {
	runTestAnalyzerTest(t, testAnalyzerConfig{
		name:     "wrong package name in test file",
		fileName: "mypackage_test.go",
		code: `package mypackage

import "testing"

func TestSomething(t *testing.T) {
}
`,
		wantErr: true,
		wantMsg: "KTN-TEST-001",
	})
}

// TestTestAnalyzerKTNTEST004 teste la détection de fonctions de test dans fichiers non-test.
//
// Params:
//   - t: instance de test
func TestTestAnalyzerKTNTEST004(t *testing.T) {
	runTestAnalyzerTest(t, testAnalyzerConfig{
		name:     "Test function in non-test file",
		fileName: "mypackage.go",
		code: `package mypackage

import "testing"

func TestSomething(t *testing.T) {
}
`,
		wantErr: true,
		wantMsg: "KTN-TEST-004",
	})

	runTestAnalyzerTest(t, testAnalyzerConfig{
		name:     "Benchmark function in non-test file",
		fileName: "mypackage.go",
		code: `package mypackage

import "testing"

func BenchmarkSomething(b *testing.B) {
}
`,
		wantErr: true,
		wantMsg: "KTN-TEST-004",
	})

	runTestAnalyzerTest(t, testAnalyzerConfig{
		name:     "Example function in non-test file",
		fileName: "mypackage.go",
		code: `package mypackage

func ExampleSomething() {
}
`,
		wantErr: true,
		wantMsg: "KTN-TEST-004",
	})
}

// Note: KTN-TEST-002 et KTN-TEST-003 ne peuvent pas être testés facilement
// avec un seul fichier parsé car ils vérifient l'existence de fichiers sur disque.
// Ces règles sont testées via l'InterfaceAnalyzer qui les appelle dans un contexte réel.
