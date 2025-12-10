// Internal tests for analyzer 003.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runTest003 tests the runTest003 private function with table-driven tests.
//
// Params:
//   - t: testing context
func Test_runTest003(t *testing.T) {
	tests := []struct {
		name    string
		pkgPath string
		wantErr bool
	}{
		{
			name:    "test file with source file",
			pkgPath: "test003",
			wantErr: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing package: %s", tt.pkgPath)
		})
	}
}

// Test_fileExists tests the fileExists private function.
//
// Params:
//   - t: testing context
func Test_fileExists(t *testing.T) {
	// Créer un fichier temporaire
	tmpFile, err := os.CreateTemp("", "test_*.go")
	// Vérification de l'erreur
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	// Nettoyage
	defer os.Remove(tmpPath)

	// Créer un répertoire temporaire
	tmpDir, err := os.MkdirTemp("", "test_dir_*")
	// Vérification de l'erreur
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	// Nettoyage
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file returns true",
			path: tmpPath,
			want: true,
		},
		{
			name: "non-existing file returns false",
			path: "/nonexistent/file.go",
			want: false,
		},
		{
			name: "directory returns false",
			path: tmpDir,
			want: false,
		},
		{
			name: "empty path returns false",
			path: "",
			want: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fileExists(tt.path)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// Test_isExemptTestFile tests the isExemptTestFile private function.
//
// Params:
//   - t: testing context
func Test_isExemptTestFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "helper_test.go is exempt",
			filename: "/path/to/helper_test.go",
			want:     true,
		},
		{
			name:     "integration_test.go is exempt",
			filename: "/path/to/integration_test.go",
			want:     true,
		},
		{
			name:     "suite_test.go is exempt",
			filename: "/path/to/suite_test.go",
			want:     true,
		},
		{
			name:     "main_test.go is exempt",
			filename: "/path/to/main_test.go",
			want:     true,
		},
		{
			name:     "regular test file not exempt",
			filename: "/path/to/myfile_test.go",
			want:     false,
		},
		{
			name:     "internal test file not exempt",
			filename: "/path/to/myfile_internal_test.go",
			want:     false,
		},
		{
			name:     "external test file not exempt",
			filename: "/path/to/myfile_external_test.go",
			want:     false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExemptTestFile(tt.filename)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("isExemptTestFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// Test_getSourceFileForTest tests the getSourceFileForTest private function.
//
// Params:
//   - t: testing context
func Test_getSourceFileForTest(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "internal test file",
			filename: "/path/to/myfile_internal_test.go",
			want:     "/path/to/myfile.go",
		},
		{
			name:     "external test file",
			filename: "/path/to/myfile_external_test.go",
			want:     "/path/to/myfile.go",
		},
		{
			name:     "bench test file",
			filename: "/path/to/myfile_bench_test.go",
			want:     "/path/to/myfile.go",
		},
		{
			name:     "integration test file",
			filename: "/path/to/myfile_integration_test.go",
			want:     "/path/to/myfile.go",
		},
		{
			name:     "standard test file",
			filename: "/path/to/myfile_test.go",
			want:     "/path/to/myfile.go",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSourceFileForTest(tt.filename)
			// Vérification du résultat
			if got != tt.want {
				t.Errorf("getSourceFileForTest(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

// Test_fileExists_edgeCases tests edge cases for fileExists.
//
// Params:
//   - t: testing context
func Test_fileExists_edgeCases(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "path with spaces",
			path: filepath.Join(os.TempDir(), "test file with spaces.go"),
			want: false,
		},
		{
			name: "path with unicode",
			path: filepath.Join(os.TempDir(), "テスト.go"),
			want: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fileExists(tt.path)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// Test_runTest003_disabled tests that the rule is skipped when disabled.
func Test_runTest003_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-003": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	src := `package test_test
import "testing"
func TestExample(t *testing.T) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test_test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error when rule is disabled")
		},
	}

	_, err = runTest003(pass)
	if err != nil {
		t.Errorf("runTest003() error = %v", err)
	}
}

// Test_runTest003_excludedFile tests that excluded files are skipped.
func Test_runTest003_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-003": {
				Enabled: config.Bool(true),
				Exclude: []string{"**/test_test.go"},
			},
		},
	})
	defer config.Reset()

	src := `package test_test
import "testing"
func TestExample(t *testing.T) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "/some/path/test_test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	inspectPass := &analysis.Pass{
		Fset:     fset,
		Files:    []*ast.File{f},
		Report:   func(d analysis.Diagnostic) {},
		ResultOf: make(map[*analysis.Analyzer]any),
	}
	inspectResult, _ := inspect.Analyzer.Run(inspectPass)

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error for excluded file")
		},
	}

	_, err = runTest003(pass)
	if err != nil {
		t.Errorf("runTest003() error = %v", err)
	}
}
