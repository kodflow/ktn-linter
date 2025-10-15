package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// TestRunTestAnalyzerWithFS vérifie que RunTestAnalyzerWithFS fonctionne correctement.
//
// Params:
//   - t: instance de test
func TestRunTestAnalyzerWithFS(t *testing.T) {
	code := `package test
func TestSomething(t *testing.T) {
	// Test something
}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test_test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pkg := types.NewPackage("test", "test")
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Pkg:   pkg,
		Report: func(diag analysis.Diagnostic) {
			// Collecter les diagnostics
		},
	}

	mockFS := &mockFileSystem{files: make(map[string]bool)}
	_, err = analyzer.RunTestAnalyzerWithFS(pass, mockFS)
	if err != nil {
		t.Errorf("RunTestAnalyzerWithFS returned error: %v", err)
	}
}

// TestFindASTFileForTest vérifie que FindASTFileForTest trouve les fichiers correctement.
//
// Params:
//   - t: instance de test
func TestFindASTFileForTest(t *testing.T) {
	code := `package test
func TestSomething(t *testing.T) {}`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
	}

	found := analyzer.FindASTFileForTest(pass, fset.File(file.Pos()).Name())
	if found == nil {
		t.Error("FindASTFileForTest should have found the file")
	}

	notFound := analyzer.FindASTFileForTest(pass, "/nonexistent.go")
	if notFound != nil {
		t.Error("FindASTFileForTest should not have found nonexistent file")
	}
}

// mockFileSystem implémente filesystem.FileSystem pour les tests.
type mockFileSystem struct {
	files map[string]bool
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	// Retourne nil pour simuler qu'un fichier n'existe pas
	return nil, os.ErrNotExist
}
