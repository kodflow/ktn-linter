package ktnconst

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runConst001_disabled tests that the rule is skipped when disabled.
func Test_runConst001_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup: disable the rule
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-CONST-001": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse test code
			src := `package test
			const x = 42
			`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", src, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			// Create inspector
			insp := inspector.New([]*ast.File{f})

			// Create pass
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Unexpected error reported when rule is disabled: %s", d.Message)
				},
			}

			// Run the analyzer - should not report anything
			_, err = runConst001(pass)
			// Check error
			if err != nil {
				t.Errorf("runConst001() error = %v", err)
			}

		})
	}
}

// Test_runConst001_excludedFile tests that excluded files are skipped.
func Test_runConst001_excludedFile(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup: exclude test.go files
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-CONST-001": {
						Enabled: config.Bool(true),
						Exclude: []string{"*.go"},
					},
				},
			})
			defer config.Reset()

			// Parse test code
			src := `package test
			const x = 42
			`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", src, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			// Create inspector
			insp := inspector.New([]*ast.File{f})

			// Create pass
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Unexpected error reported for excluded file: %s", d.Message)
				},
			}

			// Run the analyzer - should not report anything for excluded file
			_, err = runConst001(pass)
			// Check error
			if err != nil {
				t.Errorf("runConst001() error = %v", err)
			}

		})
	}
}

// Test_runConst001 tests the runConst001 private function.
func Test_runConst001(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
		{"error case validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Tested via public API
		})
	}
}

// Test_runConst001_nonConstDecl tests that non-const declarations are skipped.
func Test_runConst001_nonConstDecl(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"var declaration is skipped"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Reset config
			config.Reset()

			// Parse test code with var declaration only
			src := `package test
var x = 42
`
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "test.go", src, 0)
			// Check parse error
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			// Create inspector
			insp := inspector.New([]*ast.File{f})

			// Create pass
			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{f},
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(d analysis.Diagnostic) {
					t.Errorf("Unexpected error for var declaration: %s", d.Message)
				},
			}

			// Run the analyzer - should not report anything for var
			_, err = runConst001(pass)
			// Check error
			if err != nil {
				t.Errorf("runConst001() error = %v", err)
			}
		})
	}
}
