package ktnvar

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

// Test_runVar006_disabled tests runVar006 with disabled rule.
func Test_runVar006_disabled(t *testing.T) {
	// Setup config with rule disabled
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-006": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	// Parse simple code with builtin shadowing
	code := `package test
var string = "test"
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	result, err := runVar006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar006() error = %v", err)
	}
	// Should return nil
	if result != nil {
		t.Errorf("runVar006() result = %v, expected nil", result)
	}
	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar006() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar006_fileExcluded tests runVar006 with excluded file.
func Test_runVar006_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-006": {
				Exclude: []string{"test.go"},
			},
		},
	})
	defer config.Reset()

	// Parse code with builtin shadowing
	code := `package test
var string = "test"
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar006() error = %v", err)
	}
	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar006() reported %d issues, expected 0 when file excluded", reportCount)
	}
}

// Test_runVar006_nonVarDecl tests runVar006 skips non-var declarations.
func Test_runVar006_nonVarDecl(t *testing.T) {
	// Reset config to enable rule
	config.Reset()

	// Parse code with const declaration only
	code := `package test
const string = "test"
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = runVar006(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar006() error = %v", err)
	}
	// Should not report anything for const
	if reportCount != 0 {
		t.Errorf("runVar006() reported %d issues, expected 0 for const decl", reportCount)
	}
}

// Test_isBuiltinIdentifier006 tests the isBuiltinIdentifier006 function.
func Test_isBuiltinIdentifier006(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"blank identifier", "_", false},
		{"builtin type bool", "bool", true},
		{"builtin type string", "string", true},
		{"builtin func len", "len", true},
		{"builtin const nil", "nil", true},
		{"regular name", "myVar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBuiltinIdentifier006(tt.input)
			// Check result
			if result != tt.expected {
				t.Errorf("isBuiltinIdentifier006(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
