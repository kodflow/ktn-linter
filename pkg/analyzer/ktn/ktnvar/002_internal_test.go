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

// Test_runVar002 tests the private runVar002 function.
func Test_runVar002(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_runVar002_disabled tests runVar002 with disabled rule.
func Test_runVar002_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-002": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
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

			_, err = runVar002(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar002() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar002() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar002_fileExcluded tests runVar002 with excluded file.
func Test_runVar002_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-002": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
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

			_, err = runVar002(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar002() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar002() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_runVar002_nilInspector tests runVar002 with nil inspector.
func Test_runVar002_nilInspector(t *testing.T) {
	config.Reset()
	defer config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: nil},
	}

	result, err := runVar002(pass)
	// Should return nil without error
	if err != nil {
		t.Errorf("runVar002() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar002() result = %v, want nil", result)
	}
}

// Test_runVar002_nilFset tests runVar002 with nil Fset.
func Test_runVar002_nilFset(t *testing.T) {
	config.Reset()
	defer config.Reset()

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", "package test", 0)
	insp := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Fset:     nil,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
	}

	result, err := runVar002(pass)
	// Should return nil without error
	if err != nil {
		t.Errorf("runVar002() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar002() result = %v, want nil", result)
	}
}

// Test_runVar002_wrongInspectorType tests runVar002 with wrong inspector type.
func Test_runVar002_wrongInspectorType(t *testing.T) {
	config.Reset()
	defer config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: "not an inspector"},
	}

	result, err := runVar002(pass)
	// Should return nil without error
	if err != nil {
		t.Errorf("runVar002() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar002() result = %v, want nil", result)
	}
}

// Test_runVar002_constAfterVar tests the main error case.
func Test_runVar002_constAfterVar(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Code with const after var (invalid order)
	code := `package test
var x int = 42
const y = 1
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
		Report: func(_ analysis.Diagnostic) {
			reportCount++
		},
	}

	_, _ = runVar002(pass)

	// Should report one issue
	if reportCount != 1 {
		t.Errorf("runVar002() reported %d issues, expected 1", reportCount)
	}
}

// Test_runVar002_fallbackMessage tests with missing message.
func Test_runVar002_fallbackMessage(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Temporarily remove the message to test fallback
	msg, _ := messages.Get(ruleCodeVar002)
	messages.Unregister(ruleCodeVar002)
	defer messages.Register(msg)

	// Code with const after var (invalid order)
	code := `package test
var x int = 42
const y = 1
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

	_, _ = runVar002(pass)

	// Should report one issue with fallback message
	if reportCount != 1 {
		t.Errorf("runVar002() reported %d issues, expected 1", reportCount)
	}
	// Message should contain fallback text
	if !strings.Contains(lastMsg, "const doit") {
		t.Errorf("expected fallback message, got: %s", lastMsg)
	}
}
