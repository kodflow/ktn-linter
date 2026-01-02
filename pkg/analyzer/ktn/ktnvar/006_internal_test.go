package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
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

// Test_runVar006_nilInspector tests runVar006 with nil inspector.
func Test_runVar006_nilInspector(t *testing.T) {
	// Reset config to enable rule
	config.Reset()

	fset := token.NewFileSet()
	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: nil,
		},
		Report: func(_d analysis.Diagnostic) {},
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
}

// Test_runVar006_invalidInspector tests runVar006 with wrong type.
func Test_runVar006_invalidInspector(t *testing.T) {
	// Reset config to enable rule
	config.Reset()

	fset := token.NewFileSet()
	pass := &analysis.Pass{
		Fset: fset,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: "not an inspector",
		},
		Report: func(_d analysis.Diagnostic) {},
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
}

// Test_runVar006_nilFset tests runVar006 with nil Fset.
func Test_runVar006_nilFset(t *testing.T) {
	// Reset config to enable rule
	config.Reset()

	code := `package test
var x int
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	pass := &analysis.Pass{
		Fset: nil, // nil Fset
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {},
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
}

// Test_runVar006_wrongInspectorType tests runVar006 with wrong inspector type.
func Test_runVar006_wrongInspectorType(t *testing.T) {
	config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: "not an inspector"},
		Report:   func(_ analysis.Diagnostic) {},
	}

	result, err := runVar006(pass)
	if err != nil {
		t.Errorf("runVar006() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar006() result = %v, want nil", result)
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

// Test_runVar006_missingMessage tests runVar006 when message is missing.
func Test_runVar006_missingMessage(t *testing.T) {
	// Reset config to enable rule
	config.Reset()

	// Store original messages and clear
	originalMsg, hasOriginal := messages.Get("KTN-VAR-006")
	messages.Unregister("KTN-VAR-006")
	defer func() {
		// Restore original message
		if hasOriginal {
			messages.Register(originalMsg)
		}
	}()

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
	// Should report fallback message
	if reportCount != 1 {
		t.Errorf("runVar006() reported %d issues, expected 1 with fallback message", reportCount)
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
