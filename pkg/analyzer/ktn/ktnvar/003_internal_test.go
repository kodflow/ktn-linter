package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// TestHasUnderscore003 tests the hasUnderscore003 function.
func TestHasUnderscore003(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "blank identifier",
			input:    "_",
			expected: false,
		},
		{
			name:     "camelCase",
			input:    "myVariable",
			expected: false,
		},
		{
			name:     "PascalCase",
			input:    "MyVariable",
			expected: false,
		},
		{
			name:     "snake_case",
			input:    "my_variable",
			expected: true,
		},
		{
			name:     "SCREAMING_SNAKE_CASE",
			input:    "MY_VARIABLE",
			expected: true,
		},
		{
			name:     "mixed_Case",
			input:    "My_Variable",
			expected: true,
		},
		{
			name:     "single underscore prefix",
			input:    "_private",
			expected: true,
		},
		{
			name:     "acronym",
			input:    "HTTPStatus",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasUnderscore003(tt.input)
			if result != tt.expected {
				t.Errorf("hasUnderscore003(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_runVar003_disabled tests runVar003 with disabled rule.
func Test_runVar003_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-003": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	code := `package test
var my_var int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, err = runVar003(pass)
	if err != nil {
		t.Fatalf("runVar003() error = %v", err)
	}

	// Should not report anything when disabled
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 when disabled", reportCount)
	}
}

// Test_runVar003_fileExcluded tests runVar003 with excluded file.
func Test_runVar003_fileExcluded(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-003": {Exclude: []string{"test.go"}},
		},
	})
	defer config.Reset()

	code := `package test
var my_var int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, err = runVar003(pass)
	if err != nil {
		t.Fatalf("runVar003() error = %v", err)
	}

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 when file excluded", reportCount)
	}
}

// Test_runVar003_nilInspector tests runVar003 with nil inspector.
func Test_runVar003_nilInspector(t *testing.T) {
	config.Reset()
	defer config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: nil},
	}

	result, err := runVar003(pass)
	if err != nil {
		t.Errorf("runVar003() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar003() result = %v, want nil", result)
	}
}

// Test_runVar003_nilFset tests runVar003 with nil Fset.
func Test_runVar003_nilFset(t *testing.T) {
	config.Reset()
	defer config.Reset()

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", "package test", 0)
	insp := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Fset:     nil,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
	}

	result, err := runVar003(pass)
	if err != nil {
		t.Errorf("runVar003() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar003() result = %v, want nil", result)
	}
}

// Test_runVar003_wrongInspectorType tests runVar003 with wrong inspector type.
func Test_runVar003_wrongInspectorType(t *testing.T) {
	config.Reset()
	defer config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: "not an inspector"},
	}

	result, err := runVar003(pass)
	if err != nil {
		t.Errorf("runVar003() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar003() result = %v, want nil", result)
	}
}

// Test_runVar003_snakeCase tests runVar003 detecting snake_case.
func Test_runVar003_snakeCase(t *testing.T) {
	config.Reset()
	defer config.Reset()

	code := `package test
var my_var int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar003(pass)

	// Should report one issue
	if reportCount != 1 {
		t.Errorf("runVar003() reported %d issues, expected 1", reportCount)
	}
}

// Test_checkVar003Names_blankIdent tests checkVar003Names with blank identifier.
func Test_checkVar003Names_blankIdent(t *testing.T) {
	config.Reset()
	defer config.Reset()

	code := `package test
var _ int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar003(pass)

	// Should not report for blank identifier
	if reportCount != 0 {
		t.Errorf("runVar003() reported %d issues, expected 0 for blank identifier", reportCount)
	}
}

// Test_checkVar003Names_fallbackMessage tests with missing message.
func Test_checkVar003Names_fallbackMessage(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Temporarily remove the message to test fallback
	msg, _ := messages.Get(ruleCodeVar003)
	messages.Unregister(ruleCodeVar003)
	defer messages.Register(msg)

	code := `package test
var my_var int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0
	var lastMsg string

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report: func(d analysis.Diagnostic) {
			reportCount++
			lastMsg = d.Message
		},
	}

	_, _ = runVar003(pass)

	// Should report one issue with fallback message
	if reportCount != 1 {
		t.Errorf("runVar003() reported %d issues, expected 1", reportCount)
	}
	// Message should contain fallback text
	if !strings.Contains(lastMsg, "camelCase") {
		t.Errorf("expected fallback message containing camelCase, got: %s", lastMsg)
	}
}
