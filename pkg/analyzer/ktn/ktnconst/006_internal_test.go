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

// Test_runConst006 tests the runConst006 function.
func Test_runConst006(t *testing.T) {
	tests := []struct {
		name      string
		fn        func(*testing.T)
		expectErr bool
	}{
		{"disabled", Test_runConst006_disabled, false},
		{"excludedFile", Test_runConst006_excludedFile, false},
		{"nonConstDecl", Test_runConst006_nonConstDecl, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Call the sub-test function
			// Sub-tests handle their own error checking
			tt.fn(t)
			// Verify no unexpected errors occurred
			if t.Failed() && !tt.expectErr {
				// Test failed unexpectedly
				t.Errorf("Test failed when no error was expected")
			}
		})
	}
}

// Test_runConst006_disabled tests that the rule is skipped when disabled.
func Test_runConst006_disabled(t *testing.T) {
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
					"KTN-CONST-006": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse test code with builtin name (should not be reported)
			src := `package test
const true = 42
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
			_, err = runConst006(pass)
			// Check error
			if err != nil {
				t.Errorf("runConst006() error = %v", err)
			}

		})
	}
}

// Test_runConst006_excludedFile tests that excluded files are skipped.
func Test_runConst006_excludedFile(t *testing.T) {
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
					"KTN-CONST-006": {
						Enabled: config.Bool(true),
						Exclude: []string{"*.go"},
					},
				},
			})
			defer config.Reset()

			// Parse test code with builtin name
			src := `package test
const true = 42
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
			_, err = runConst006(pass)
			// Check error
			if err != nil {
				t.Errorf("runConst006() error = %v", err)
			}

		})
	}
}

// Test_isBuiltinIdentifier tests the private isBuiltinIdentifier function.
func Test_isBuiltinIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"blank identifier", "_", false},
		{"builtin type bool", "bool", true},
		{"builtin type int", "int", true},
		{"builtin type string", "string", true},
		{"builtin constant true", "true", true},
		{"builtin constant false", "false", true},
		{"builtin constant iota", "iota", true},
		{"builtin nil", "nil", true},
		{"builtin function len", "len", true},
		{"builtin function make", "make", true},
		{"builtin function append", "append", true},
		{"non-builtin name", "myVar", false},
		{"non-builtin MaxSize", "MaxSize", false},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isBuiltinIdentifier(tt.input)
			// Verify result
			if result != tt.expected {
				t.Errorf("isBuiltinIdentifier(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_runConst006_nonConstDecl tests that non-const declarations are skipped.
func Test_runConst006_nonConstDecl(t *testing.T) {
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
var true = 42
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
			_, err = runConst006(pass)
			// Check error
			if err != nil {
				t.Errorf("runConst006() error = %v", err)
			}
		})
	}
}
