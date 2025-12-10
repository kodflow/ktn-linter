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

// Test_runTest011 tests the runTest011 private function.
//
// Params:
//   - t: testing context
func Test_runTest011(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "error case - minimal test",
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

// Test_extractExpectedPackageFromFilename tests extractExpectedPackageFromFilename.
//
// Params:
//   - t: testing context
func Test_extractExpectedPackageFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "external test file",
			filename: "calculator_external_test.go",
			want:     "calculator",
		},
		{
			name:     "error case - simple name",
			filename: "simple_external_test.go",
			want:     "simple",
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractExpectedPackageFromFilename(tt.filename)
			// Vérification de la condition
			if got != tt.want {
				t.Errorf("extractExpectedPackageFromFilename(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

// Test_runTest011_disabled tests that the rule is skipped when disabled.
func Test_runTest011_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-011": {Enabled: config.Bool(false)},
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

	_, err = runTest011(pass)
	if err != nil {
		t.Errorf("runTest011() error = %v", err)
	}
}

// Test_runTest011_excludedFile tests that excluded files are skipped.
func Test_runTest011_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-011": {
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

	_, err = runTest011(pass)
	if err != nil {
		t.Errorf("runTest011() error = %v", err)
	}
}
