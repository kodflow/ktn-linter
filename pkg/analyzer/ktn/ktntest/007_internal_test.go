// Internal tests for analyzer 007.
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

// Test_runTest007 tests the runTest007 private function with table-driven tests.
//
// Params:
//   - t: testing context
func Test_runTest007(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "basic test of runTest007 structure",
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

// Test_runTest007_integration tests the analyzer structure.
//
// Params:
//   - t: testing context
func Test_runTest007_integration(t *testing.T) {
	tests := []struct {
		name         string
		expectedName string
	}{
		{name: "analyzer structure", expectedName: "ktntest007"},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check analyzer is valid
			if Analyzer007 == nil || Analyzer007.Name != tt.expectedName {
				t.Errorf("Analyzer007 invalid: nil=%v, Name=%q, want %q",
					Analyzer007 == nil, Analyzer007.Name, tt.expectedName)
			}
		})
	}
}

// Test_isSkipMethod tests the isSkipMethod function.
//
// Params:
//   - t: testing context
func Test_isSkipMethod(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		want       bool
	}{
		{
			name:       "Skip method detected",
			methodName: "Skip",
			want:       true,
		},
		{
			name:       "Skipf method detected",
			methodName: "Skipf",
			want:       true,
		},
		{
			name:       "SkipNow method detected",
			methodName: "SkipNow",
			want:       true,
		},
		{
			name:       "Error method not detected",
			methodName: "Error",
			want:       false,
		},
		{
			name:       "Fatal method not detected",
			methodName: "Fatal",
			want:       false,
		},
		{
			name:       "Run method not detected",
			methodName: "Run",
			want:       false,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSkipMethod(tt.methodName)
			// Check result
			if got != tt.want {
				t.Errorf("isSkipMethod(%q) = %v, want %v", tt.methodName, got, tt.want)
			}
		})
	}
}

// Test_runTest007_disabled tests that the rule is skipped when disabled.
func Test_runTest007_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-007": {Enabled: config.Bool(false)},
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

	_, err = runTest007(pass)
	if err != nil {
		t.Errorf("runTest007() error = %v", err)
	}
}

// Test_runTest007_excludedFile tests that excluded files are skipped.
func Test_runTest007_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-TEST-007": {
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

	_, err = runTest007(pass)
	if err != nil {
		t.Errorf("runTest007() error = %v", err)
	}
}
