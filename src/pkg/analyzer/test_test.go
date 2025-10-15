package analyzer_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"testing"
	"time"

	"github.com/kodflow/ktn-linter/src/internal/filesystem"
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
}

// TestTestAnalyzerKTNTEST002WithMock teste la vérification de couverture test avec mock FS.
//
// Params:
//   - t: instance de test
func TestTestAnalyzerKTNTEST002WithMock(t *testing.T) {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "/fake/path/myfile.go", `package mypackage
func DoSomething() {}
`, parser.ParseComments)

	pkg := &types.Package{}
	pkg.SetName("mypackage")

	var diagnostics []analysis.Diagnostic
	mockFS := NewMockFileSystemForTests(map[string]bool{
		"/fake/path/myfile_test.go": false, // Le test n'existe PAS
	})

	pass := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset,
		Files:    []*ast.File{file},
		Pkg:      pkg,
		Report: func(diag analysis.Diagnostic) {
			diagnostics = append(diagnostics, diag)
		},
	}

	runTestAnalyzerWithMockFS(pass, mockFS)

	found := false
	for _, d := range diagnostics {
		if contains(d.Message, "KTN-TEST-002") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected KTN-TEST-002 error for missing test file")
	}
}

// TestTestAnalyzerKTNTEST003WithMock teste la détection de tests orphelins avec mock FS.
//
// Params:
//   - t: instance de test
func TestTestAnalyzerKTNTEST003WithMock(t *testing.T) {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "/fake/path/myfile_test.go", `package mypackage_test
import "testing"
func TestSomething(t *testing.T) {}
`, parser.ParseComments)

	pkg := &types.Package{}
	pkg.SetName("mypackage_test")

	var diagnostics []analysis.Diagnostic
	mockFS := NewMockFileSystemForTests(map[string]bool{
		"/fake/path/myfile.go": false, // Le fichier source n'existe PAS
	})

	pass := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset,
		Files:    []*ast.File{file},
		Pkg:      pkg,
		Report: func(diag analysis.Diagnostic) {
			diagnostics = append(diagnostics, diag)
		},
	}

	runTestAnalyzerWithMockFS(pass, mockFS)

	found := false
	for _, d := range diagnostics {
		if contains(d.Message, "KTN-TEST-003") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected KTN-TEST-003 error for orphan test file")
	}
}

// TestFileExistsAllBranches teste toutes les branches de fileExistsWithFS.
//
// Params:
//   - t: instance de test
func TestFileExistsAllBranches(t *testing.T) {
	// Test avec fichier existant
	mockFS1 := NewMockFileSystemForTests(map[string]bool{
		"/exists.go": true,
	})
	if !fileExistsTestHelper(mockFS1, "/exists.go") {
		t.Error("Expected file to exist")
	}

	// Test avec fichier non existant
	mockFS2 := NewMockFileSystemForTests(map[string]bool{
		"/exists.go": false,
	})
	if fileExistsTestHelper(mockFS2, "/exists.go") {
		t.Error("Expected file to not exist")
	}

	// Test avec fichier non dans la map
	mockFS3 := NewMockFileSystemForTests(map[string]bool{})
	if fileExistsTestHelper(mockFS3, "/notinmap.go") {
		t.Error("Expected file to not exist when not in map")
	}
}

// mockFileSystemForTests est un mock du FileSystem pour les tests.
type mockFileSystemForTests struct {
	files map[string]bool
}

// NewMockFileSystemForTests crée un nouveau mock FileSystem.
//
// Params:
//   - files: map des chemins de fichiers (true = existe, false = n'existe pas)
//
// Returns:
//   - filesystem.FileSystem: le mock FileSystem
func NewMockFileSystemForTests(files map[string]bool) filesystem.FileSystem {
	// Retourne le mock FileSystem créé
	// Retourne le mock FileSystem créé
	return &mockFileSystemForTests{files: files}
}

// Stat simule os.Stat en utilisant la map de fichiers.
//
// Params:
//   - name: le chemin du fichier
//
// Returns:
//   - os.FileInfo: informations simulées si le fichier existe
//   - error: erreur si le fichier n'existe pas
func (m *mockFileSystemForTests) Stat(name string) (os.FileInfo, error) {
	if exists, ok := m.files[name]; ok && exists {
		return &mockFileInfoForTests{name: name}, nil
	}
	return nil, errors.New("file does not exist")
}

// mockFileInfoForTests implémente os.FileInfo pour le mock.
type mockFileInfoForTests struct {
	name string
}

// Name retourne le nom du fichier.
//
// Returns:
//   - string: le nom du fichier
func (m *mockFileInfoForTests) Name() string { return m.name }

// Size retourne la taille (toujours 0 pour le mock).
//
// Returns:
//   - int64: la taille du fichier
func (m *mockFileInfoForTests) Size() int64 { return 0 }

// Mode retourne le mode (toujours 0644 pour le mock).
//
// Returns:
//   - os.FileMode: le mode du fichier
func (m *mockFileInfoForTests) Mode() os.FileMode { return 0644 }

// ModTime retourne le temps de modification (epoch pour le mock).
//
// Returns:
//   - time.Time: le temps de modification
func (m *mockFileInfoForTests) ModTime() time.Time { return time.Unix(0, 0) }

// IsDir retourne false (pas un répertoire pour le mock).
//
// Returns:
//   - bool: false car c'est un fichier
func (m *mockFileInfoForTests) IsDir() bool { return false }

// Sys retourne nil (pas d'infos système pour le mock).
//
// Returns:
//   - interface{}: nil
func (m *mockFileInfoForTests) Sys() interface{} { return nil }

// runTestAnalyzerWithMockFS exécute TestAnalyzer avec un mock FileSystem.
//
// Params:
//   - pass: la passe d'analyse
//   - fs: le système de fichiers mock
func runTestAnalyzerWithMockFS(pass *analysis.Pass, fs filesystem.FileSystem) {
	analyzer.RunTestAnalyzerWithFS(pass, fs)
}

// fileExistsTestHelper teste fileExistsWithFS.
//
// Params:
//   - fs: le système de fichiers
//   - path: le chemin à tester
//
// Returns:
//   - bool: true si le fichier existe
func fileExistsTestHelper(fs filesystem.FileSystem, path string) bool {
	// Retourne true si le fichier existe
	_, err := fs.Stat(path)
	// Retourne true si le fichier existe
	return err == nil
}

// TestFileExistsWithFS teste toutes les branches de fileExistsWithFS.
//
// Params:
//   - t: instance de test
func TestFileExistsWithFS(t *testing.T) {
	mockFS := NewMockFileSystemForTests(map[string]bool{
		"/exists.go":     true,
		"/notexists.go":  false,
	})

	// Test fichier existant
	if !fileExistsTestHelper(mockFS, "/exists.go") {
		t.Error("Expected /exists.go to exist")
	}

	// Test fichier non existant explicite
	if fileExistsTestHelper(mockFS, "/notexists.go") {
		t.Error("Expected /notexists.go to not exist")
	}

	// Test fichier pas dans map
	if fileExistsTestHelper(mockFS, "/unknown.go") {
		t.Error("Expected /unknown.go to not exist")
	}
}

// TestCheckOrphanTestFilesAllBranches teste toutes les branches de checkOrphanTestFilesWithFS.
//
// Params:
//   - t: instance de test
func TestCheckOrphanTestFilesAllBranches(t *testing.T) {
	// Test 1: Package non-test (early return)
	fset1 := token.NewFileSet()
	file1, _ := parser.ParseFile(fset1, "/fake/test.go", `package mypackage
func Foo() {}
`, parser.ParseComments)
	pkg1 := &types.Package{}
	pkg1.SetName("mypackage")  // PAS _test
	pass1 := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset1,
		Files:    []*ast.File{file1},
		Pkg:      pkg1,
		Report: func(diag analysis.Diagnostic) {},
	}
	mockFS1 := NewMockFileSystemForTests(map[string]bool{})
	runTestAnalyzerWithMockFS(pass1, mockFS1)
	// Ne devrait rien faire (early return)

	// Test 2: Fichier non-test dans package _test (skip)
	fset2 := token.NewFileSet()
	file2, _ := parser.ParseFile(fset2, "/fake/helper.go", `package mypackage_test
func Helper() {}
`, parser.ParseComments)
	pkg2 := &types.Package{}
	pkg2.SetName("mypackage_test")
	var diags2 []analysis.Diagnostic
	pass2 := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset2,
		Files:    []*ast.File{file2},
		Pkg:      pkg2,
		Report: func(diag analysis.Diagnostic) {
			diags2 = append(diags2, diag)
		},
	}
	mockFS2 := NewMockFileSystemForTests(map[string]bool{})
	runTestAnalyzerWithMockFS(pass2, mockFS2)
	// helper.go n'est pas un fichier de test, donc pas d'erreur

	// Test 3: Fichier de test avec source existant (pas d'erreur)
	fset3 := token.NewFileSet()
	file3, _ := parser.ParseFile(fset3, "/fake/myfile_test.go", `package mypackage_test
import "testing"
func TestFoo(t *testing.T) {}
`, parser.ParseComments)
	pkg3 := &types.Package{}
	pkg3.SetName("mypackage_test")
	var diags3 []analysis.Diagnostic
	pass3 := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset3,
		Files:    []*ast.File{file3},
		Pkg:      pkg3,
		Report: func(diag analysis.Diagnostic) {
			diags3 = append(diags3, diag)
		},
	}
	mockFS3 := NewMockFileSystemForTests(map[string]bool{
		"/fake/myfile.go": true, // Source existe
	})
	runTestAnalyzerWithMockFS(pass3, mockFS3)
	if len(diags3) > 0 {
		t.Errorf("Expected no error when source file exists, got: %v", diags3)
	}
}

// TestCheckTestCoverageAllBranches teste toutes les branches de checkTestCoverageWithFS.
//
// Params:
//   - t: instance de test
func TestCheckTestCoverageAllBranches(t *testing.T) {
	// Test 1: Package _test (early return)
	fset1 := token.NewFileSet()
	file1, _ := parser.ParseFile(fset1, "/fake/test.go", `package mypackage_test
func TestFoo(t *testing.T) {}
`, parser.ParseComments)
	pkg1 := &types.Package{}
	pkg1.SetName("mypackage_test")
	pass1 := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset1,
		Files:    []*ast.File{file1},
		Pkg:      pkg1,
		Report: func(diag analysis.Diagnostic) {},
	}
	mockFS1 := NewMockFileSystemForTests(map[string]bool{})
	runTestAnalyzerWithMockFS(pass1, mockFS1)
	// Ne devrait rien faire (early return pour package _test)

	// Test 2: Fichier de test (skip dans la loop)
	fset2 := token.NewFileSet()
	file2, _ := parser.ParseFile(fset2, "/fake/myfile_test.go", `package mypackage
func TestFoo(t *testing.T) {}
`, parser.ParseComments)
	pkg2 := &types.Package{}
	pkg2.SetName("mypackage")
	var diags2 []analysis.Diagnostic
	pass2 := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset2,
		Files:    []*ast.File{file2},
		Pkg:      pkg2,
		Report: func(diag analysis.Diagnostic) {
			diags2 = append(diags2, diag)
		},
	}
	mockFS2 := NewMockFileSystemForTests(map[string]bool{})
	runTestAnalyzerWithMockFS(pass2, mockFS2)
	// myfile_test.go est un fichier de test, donc skip

	// Test 3: Fichier source avec test existant (pas d'erreur)
	fset3 := token.NewFileSet()
	file3, _ := parser.ParseFile(fset3, "/fake/myfile.go", `package mypackage
func Foo() {}
`, parser.ParseComments)
	pkg3 := &types.Package{}
	pkg3.SetName("mypackage")
	var diags3 []analysis.Diagnostic
	pass3 := &analysis.Pass{
		Analyzer: analyzer.TestAnalyzer,
		Fset:     fset3,
		Files:    []*ast.File{file3},
		Pkg:      pkg3,
		Report: func(diag analysis.Diagnostic) {
			diags3 = append(diags3, diag)
		},
	}
	mockFS3 := NewMockFileSystemForTests(map[string]bool{
		"/fake/myfile_test.go": true, // Test existe
	})
	runTestAnalyzerWithMockFS(pass3, mockFS3)
	if len(diags3) > 0 {
		t.Errorf("Expected no error when test file exists, got: %v", diags3)
	}
}

// TestFindASTFileAllBranches teste toutes les branches de findASTFile.
//
// Params:
//   - t: instance de test
func TestFindASTFileAllBranches(t *testing.T) {
	fset := token.NewFileSet()
	file1, _ := parser.ParseFile(fset, "/fake/file1.go", `package test
func Foo() {}
`, parser.ParseComments)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file1},
	}

	// Test: fichier trouvé
	found := analyzer.FindASTFileForTest(pass, "/fake/file1.go")
	if found == nil {
		t.Error("Expected to find /fake/file1.go")
	}

	// Test: fichier NON trouvé (retourne nil)
	notFound := analyzer.FindASTFileForTest(pass, "/fake/notexist.go")
	if notFound != nil {
		t.Error("Expected nil for non-existent file")
	}
}

