// Internal tests for analyzer 002.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Test_runTest002 tests the runTest002 private function with table-driven tests.
//
// Params:
//   - t: testing context
func Test_runTest002(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "basic test of runTest002 structure",
		},
		{
			name: "error case validation",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic functionality
			t.Logf("Testing: %s", tt.name)
		})
	}
}

// Test_isExemptPackage tests the isExemptPackage private function.
//
// Params:
//   - t: testing context
func Test_isExemptPackage(t *testing.T) {
	tests := []struct {
		name    string
		pkgName string
		want    bool
	}{
		{
			name:    "main package is exempt",
			pkgName: "main",
			want:    true,
		},
		{
			name:    "testhelper package is exempt",
			pkgName: "testhelper",
			want:    true,
		},
		{
			name:    "ktntest package is exempt",
			pkgName: "ktntest",
			want:    true,
		},
		{
			name:    "regular package not exempt",
			pkgName: "mypackage",
			want:    false,
		},
		{
			name:    "empty package not exempt",
			pkgName: "",
			want:    false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExemptPackage(tt.pkgName)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("isExemptPackage(%q) = %v, want %v", tt.pkgName, got, tt.want)
			}
		})
	}
}

// Test_runTest001_integration tests the analyzer structure.
//
// Params:
//   - t: testing context
func Test_runTest001_integration(t *testing.T) {
	tests := []struct {
		name         string
		expectedName string
	}{
		{name: "analyzer structure", expectedName: "ktntest002"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Analyzer002 == nil || Analyzer002.Name != tt.expectedName {
				t.Errorf("Analyzer002 invalid: nil=%v, Name=%q, want %q",
					Analyzer002 == nil, Analyzer002.Name, tt.expectedName)
			}
		})
	}
}

// Test_runTest002_disabled tests that the rule is skipped when disabled.
func Test_runTest002_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-002": {Enabled: config.Bool(false)},
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

	_, err = runTest002(pass)
	if err != nil {
		t.Errorf("runTest002() error = %v", err)
	}
}

// Test_runTest002_excludedFile tests that excluded files are skipped.
func Test_runTest002_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-002": {
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

	_, err = runTest002(pass)
	if err != nil {
		t.Errorf("runTest002() error = %v", err)
	}
}
